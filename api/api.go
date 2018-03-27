package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"

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

	// Create API instance to store database session and collections
	api := API{
		cfg: cfg,
		db:  db,
	}

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)
	r.Use(etag)
	cli := CasClient{}
	cli.Create(url, api.db)
	// Vehicles
	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", api.VehiclesHandler)
		r.Group(func(r chi.Router) {
			r.Use(cli.casauth)
			r.Post("/create", api.VehiclesCreateHandler)
			r.Post("/edit", api.VehiclesEditHandler)
			r.Delete("/{id:[0-9]+}", api.VehiclesDeleteHandler)
		})
	})

	// Updates
	r.Route("/updates", func(r chi.Router) {
		r.Get("/", api.UpdatesHandler)
		r.Get("/message", api.UpdateMessageHandler)
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
			r.Delete("/{id:.+}", api.RoutesDeleteHandler)
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
	http.ServeFile(w, r, "admin.html")

}

//KeyHandler sends Mapbox api key to authenticated user
func (api *API) KeyHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, api.cfg.MapboxAPIKey)

}

func (api *API) AdminLogout(w http.ResponseWriter, r *http.Request) {

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
