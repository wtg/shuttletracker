package main

import (
	"net/http"
	"shuttle_tracking_2/tracking"

	"gopkg.in/cas.v1"
	"strings"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)
var users []User

type User struct{
		Name   string
}

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
		fmt.Printf("%u",App.Config.Authenticate);
		if  App.Config.Authenticate && !cas.IsAuthenticated(r) {
			cas.RedirectToLogin(w, r)
			return
		}else{
			valid := false;
			for _, u := range users{
				if(u.Name == strings.ToLower(cas.Username(r))){
					valid = true;
				}
			}
			if(App.Config.Authenticate == false){
				valid = true;
				fmt.Printf("not authenticating");
			}
			if valid{
				http.Redirect(w,r,"/admin/success/",301);
			}else{
				http.Redirect(w,r,"/admin/logout/",301);
			}
		}

}

func AdminPageServer(w http.ResponseWriter, r *http.Request) {

		if App.Config.Authenticate && !cas.IsAuthenticated(r) {
			http.Redirect(w,r,"/admin/",301);
			return
		}else{
			http.ServeFile(w, r, "admin.html")
		}

}

func AdminLogout(w http.ResponseWriter, r *http.Request) {

		if cas.IsAuthenticated(r){
			cas.RedirectToLogout(w,r);
		}

}

func main() {
	// close Mongo session when server terminates
	defer App.Session.Close()


	err := App.Users.Find(bson.M{}).All(&users)
	_ = err;

	// Start auto updater
	go App.UpdateShuttles(Config.DataFeed, Config.UpdateInterval)
	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.Handle("/admin/", App.CasAUTH.HandleFunc(AdminHandler)).Methods("GET")
	r.Handle("/admin", App.CasAUTH.HandleFunc(AdminHandler)).Methods("GET")
	r.Handle("/admin/success/", App.CasAUTH.HandleFunc(AdminPageServer)).Methods("GET")
	r.Handle("/admin/success", App.CasAUTH.HandleFunc(AdminPageServer)).Methods("GET")
	r.Handle("/admin/logout/", App.CasAUTH.HandleFunc(AdminLogout)).Methods("GET")
	r.Handle("/admin/logout", App.CasAUTH.HandleFunc(AdminLogout)).Methods("GET")
	//Has to do with r not being a client?
	r.HandleFunc("/vehicles", App.VehiclesHandler).Methods("GET")
	r.Handle("/vehicles/create", App.CasAUTH.HandleFunc(App.VehiclesCreateHandler)).Methods("POST")
	r.Handle("/vehicles/edit", App.CasAUTH.HandleFunc(App.VehiclesEditHandler)).Methods("POST")
	r.Handle("/vehicles/{id:[0-9]+}", App.CasAUTH.HandleFunc(App.VehiclesDeleteHandler)).Methods("DELETE")
	r.HandleFunc("/updates", App.UpdatesHandler).Methods("GET")
	r.HandleFunc("/updates/message", App.UpdateMessageHandler).Methods("GET")
	r.HandleFunc("/routes", App.RoutesHandler).Methods("GET")
	r.Handle("/routes/create", App.CasAUTH.HandleFunc(App.RoutesCreateHandler)).Methods("POST")
	r.Handle("/routes/{id:.+}", App.CasAUTH.HandleFunc(App.RoutesDeleteHandler)).Methods("DELETE")
	r.HandleFunc("/stops", App.StopsHandler).Methods("GET")
	r.Handle("/stops/create", App.CasAUTH.HandleFunc(App.StopsCreateHandler)).Methods("POST")
	r.Handle("/stops/{id:.+}", App.CasAUTH.HandleFunc(App.StopsDeleteHandler)).Methods("DELETE")
	//r.HandleFunc("/import", App.ImportHandler).Methods("GET")
	// Static files
	r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	// Serve requests
	hand := App.CasAUTH.Handle(r)
	if err := http.ListenAndServe(":8080", hand); err != nil {
		log.Fatalf("Unable to ListenAndServe: %v", err)
	}
}
