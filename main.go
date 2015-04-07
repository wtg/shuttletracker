package main

import (
  "fmt"
  "net/http"
  "gopkg.in/mgo.v2"
  "github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
}

func VehiclesHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Vehicles")
}

func main() {
  // Connect to MongoDB
  session, err := mgo.Dial("localhost:27017")
  if err != nil {
    panic(err)
  }
  // close Mongo session when server terminates
  defer session.Close()

  // Routing 
  r := mux.NewRouter()
  r.HandleFunc("/", IndexHandler).Methods("GET")
  r.HandleFunc("/admin", IndexHandler).Methods("GET")
  r.HandleFunc("/vehicles", VehiclesHandler).Methods("GET")
  // Static files
  r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
  // Serve requests
  http.Handle("/", r)
  http.ListenAndServe(":8080", r)
}