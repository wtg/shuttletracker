package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"shuttle_tracking_2/tracking"
)

var (
	// Config holds the global app settings.
	Config = tracking.InitConfig()
	// App holds the application itself.
	App = tracking.InitApp(Config)
)

// IndexHandler serves the index page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// AdminHandler serves the admin page.
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "admin.html")
}

func main() {
	// close Mongo session when server terminates
	defer App.Session.Close()

	// Start auto updater
	go App.UpdateShuttles(Config.DataFeed, Config.UpdateInterval)

	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/admin", AdminHandler).Methods("GET")
	r.HandleFunc("/admin/{*}", AdminHandler).Methods("GET")
	r.HandleFunc("/vehicles", App.VehiclesHandler).Methods("GET")
	r.HandleFunc("/vehicles/create", App.VehiclesCreateHandler).Methods("POST")
	r.HandleFunc("/updates", App.UpdatesHandler).Methods("GET")
	r.HandleFunc("/routes", App.RoutesHandler).Methods("GET")
	r.HandleFunc("/routes/create", App.RoutesCreateHandler).Methods("POST")
	r.HandleFunc("/stops", App.StopsHandler).Methods("GET")
	r.HandleFunc("/stops/create", App.StopsCreateHandler).Methods("POST")
	// Static files
	r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	// Serve requests
	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Unable to ListenAndServe: %v", err)
	}
}
