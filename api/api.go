package api

import (
	"encoding/json"

	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/database"
	"github.com/wtg/shuttletracker/log"
)

// Config holds API settings.
type Config struct {
	GoogleMapAPIKey      string
	GoogleMapMinDistance int
	CasURL               string
	Authenticate         bool
	ListenURL            string
	MapboxAPIKey         string
}

// API is responsible for configuring handlers for HTTP endpoints.
type API struct {
	cfg     Config
	db      database.Database
	handler http.Handler
	ms      shuttletracker.ModelService
}

// New initializes the application given a config and connects to backends.
// It also seeds any needed information to the database.
func New(cfg Config, db database.Database, ms shuttletracker.ModelService) (*API, error) {
	// Set up CAS authentication
	url, err := url.Parse(cfg.CasURL)
	if err != nil {
		return nil, err
	}

	// Create API instance to store database session and collections
	api := API{
		cfg: cfg,
		db:  db,
		ms:  ms,
	}

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)
	r.Use(etag)

	cli := CreateCASClient(url, db)

	// Vehicles
	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", api.VehiclesHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/create", api.VehiclesCreateHandler)
			r.Post("/edit", api.VehiclesEditHandler)
			r.Delete("/", api.VehiclesDeleteHandler)
		})
	})

	// Updates
	r.Route("/updates", func(r chi.Router) {
		r.Get("/", api.UpdatesHandler)
	})

	// Admin message
	r.Route("/adminMessage", func(r chi.Router) {
		r.Get("/", api.AdminMessageHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/", api.SetAdminMessage)
		})
	})

	// Routes
	r.Route("/routes", func(r chi.Router) {
		r.Get("/", api.RoutesHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/create", api.RoutesCreateHandler)
			r.Post("/schedule", api.RoutesScheduler)
			r.Post("/edit", api.RoutesEditHandler)
			r.Delete("/", api.RoutesDeleteHandler)
		})
	})

	// Stops
	r.Route("/stops", func(r chi.Router) {
		r.Get("/", api.StopsHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/create", api.StopsCreateHandler)
			r.Delete("/{id:.+}", api.StopsDeleteHandler)
		})
	})

	r.Get("/logout/", cli.logout)
	// Admin
	r.Route("/admin", func(r chi.Router) {
		r.Use(cli.casauth)
		r.Get("/", api.AdminHandler)
		r.Get("/login", api.AdminHandler)
		r.Get("/logout", cli.logout)
	})

	r.Group(func(r chi.Router) {
		r.Use(cli.casauth)
		r.Get("/getKey/", api.KeyHandler)
	})

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
	id := r.URL.Query()
	//authenticated with a ticket, redirect the request
	if len(id["ticket"]) != 0 {
		http.Redirect(w, r, "/admin", 301)
	}
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, "admin.html")

}

//KeyHandler sends Mapbox api key to authenticated user
func (api *API) KeyHandler(w http.ResponseWriter, r *http.Request) {
	err := WriteJSON(w, api.cfg.MapboxAPIKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
