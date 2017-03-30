package main

import (
	"net/http"
	"shuttle_tracking_2/tracking"
  _"io"
  "encoding/json"
  "os"
  "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

)
/*This serves collected shuttle data in a loop*/
var (
	// Config holds the global app settings.
	Config = tracking.InitConfig()
	// App holds the application itself.
	App = tracking.InitApp(Config)
)

type Dump struct{
	Id				int				`json:id						bson:"id"`
	Data			string		`json:"Data"				bson:"Data"`
}
var data []Dump;
func readFromFile(id int) string{
  if(len(data) == 0){
    in, err := os.Open("dummy_data.json")
    if err != nil {
        fmt.Printf("you failed")
    }
      jsonParser := json.NewDecoder(in)
      if err = jsonParser.Decode(&data); err != nil {
        fmt.Printf("you failed")
     }
  }
  return data[id].Data;
}

var count = 0;

// IndexHandler serves the index page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
  byteArray := []byte(readFromFile(count));
  w.Write(byteArray)
  count += 1
  if(count >= len(data)){
    count = 0;
  }
}
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	//r.HandleFunc("/import", App.ImportHandler).Methods("GET")
	// Static files
	r.PathPrefix("/bower_components/").Handler(http.StripPrefix("/bower_components/", http.FileServer(http.Dir("bower_components/"))))
	// Serve requests
	http.Handle("/", r)

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("Unable to ListenAndServe: %v", err)
	}
}
