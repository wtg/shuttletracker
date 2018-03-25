package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	"gopkg.in/cas.v1"

	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
)

// Configuration holds the settings for connecting to outside resources.
type Config struct {
	GoogleMapAPIKey      string
	GoogleMapMinDistance int
	CasURL               string
	Authenticate         bool
	ListenURL            string
	MapboxAPIKey         string
}

// App holds references to Mongo resources.
type API struct {
	cfg     Config
	CasAUTH *cas.Client
	CasMEM  *cas.MemoryStore
	db      database.Database
	handler http.Handler
}

// InitApp initializes the application given a config and connects to backends.
// It also seeds any needed information to the database.
func New(cfg Config, db database.Database) (*API, error) {
	// Set up CAS authentication
	url, err := url.Parse(cfg.CasURL)
	if err != nil {
		return nil, err
	}
	var tickets *cas.MemoryStore

	client := cas.NewClient(&cas.Options{
		URL:   url,
		Store: nil,
	})

	// Create API instance to store database session and collections
	api := API{
		cfg:     cfg,
		CasAUTH: client,
		CasMEM:  tickets,
		db:      db,
	}

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)
	r.Use(etag)

	// Vehicles
	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", api.VehiclesHandler)
		r.Method("POST", "/create", api.CasAUTH.HandleFunc(api.VehiclesCreateHandler))
		r.Method("POST", "/edit", api.CasAUTH.HandleFunc(api.VehiclesEditHandler))
		r.Method("DELETE", "/{id:[0-9]+}", api.CasAUTH.HandleFunc(api.VehiclesDeleteHandler))
	})

	// Updates
	r.Route("/updates", func(r chi.Router) {
		r.Get("/", api.UpdatesHandler)
		r.Get("/message", api.UpdateMessageHandler)
	})

	// Admin message
	r.Route("/adminMessage", func(r chi.Router) {
		r.Get("/", api.AdminMessageHandler)
		r.Post("/", api.SetAdminMessage)
	})

	// Routes
	r.Route("/routes", func(r chi.Router) {
		r.Get("/", api.RoutesHandler)
		r.Method("POST", "/create", api.CasAUTH.HandleFunc(api.RoutesCreateHandler))
		r.Method("POST", "/schedule", api.CasAUTH.HandleFunc(api.RoutesScheduler))
		r.Method("POST", "/edit", api.CasAUTH.HandleFunc(api.RoutesEditHandler))
		r.Method("DELETE", "/{id:.+}", api.CasAUTH.HandleFunc(api.RoutesDeleteHandler))
	})

	// Stops
	r.Route("/stops", func(r chi.Router) {
		r.Get("/", api.StopsHandler)
		r.Method("POST", "/create", api.CasAUTH.HandleFunc(api.StopsCreateHandler))
		r.Method("DELETE", "/{id:.+}", api.CasAUTH.HandleFunc(api.StopsDeleteHandler))
	})

	// Admin
	r.Route("/admin", func(r chi.Router) {
		r.Method("GET", "/", api.CasAUTH.HandleFunc(api.AdminHandler))
		r.Method("GET", "/success/", api.CasAUTH.HandleFunc(api.AdminPageServer))
		r.Method("GET", "/logout/", api.CasAUTH.HandleFunc(api.AdminLogout))

	})

	r.Method("GET", "/getKey/", api.CasAUTH.HandleFunc(api.KeyHandler))

	// Static files
	r.Get("/", IndexHandler)
	r.Method("GET", "/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	api.handler = r

	return &api, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		ListenURL:    "0.0.0.0:8080",
		Authenticate: true,
	}
	v.SetDefault("api.listenurl", cfg.ListenURL)
	v.SetDefault("api.casurl", cfg.CasURL)
	v.SetDefault("api.authenticate", cfg.Authenticate)
	return cfg
}

func (api *API) Run() {
	if err := http.ListenAndServe(api.cfg.ListenURL, api.handler); err != nil {
		log.WithError(err).Error("Unable to serve.")
	}
}

// IndexHandler serves the index page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// AdminHandler serves the admin page.
func (api *API) AdminHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	} else {
		valid := false
		users, _ := api.db.GetUsers()
		for _, u := range users {
			if u.Name == strings.ToLower(cas.Username(r)) {
				valid = true
			}
		}
		if api.cfg.Authenticate == false {
			valid = true
			fmt.Printf("not authenticating")
		}
		if valid {
			http.Redirect(w, r, "/admin/success/", 301)
		} else {
			http.Redirect(w, r, "/admin/logout/", 301)
		}
	}

}

//KeyHandler sends Mapbox api key to authenticated user
func (api *API) KeyHandler(w http.ResponseWriter, r *http.Request) {
	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		http.Redirect(w, r, "/admin/", 301)
	} else {
		WriteJSON(w, api.cfg.MapboxAPIKey)
	}
}

func (api *API) AdminPageServer(w http.ResponseWriter, r *http.Request) {

	if api.cfg.Authenticate && !cas.IsAuthenticated(r) {
		http.Redirect(w, r, "/admin/", 301)
		return
	} else {
		http.ServeFile(w, r, "admin.html")
	}

}

func (api *API) AdminLogout(w http.ResponseWriter, r *http.Request) {

	if cas.IsAuthenticated(r) {
		cas.RedirectToLogout(w, r)
	}

}

// WriteJSON writes the data as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}
