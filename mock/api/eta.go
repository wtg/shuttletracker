package api

import (
	_"encoding/json"
	"net/http"
	_"strconv"
	"fmt"
	_"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
    "reflect"
    "math"
)

// EtaHandler finds all of the routes in the database
func (api *API) EtaHandler(w http.ResponseWriter, r *http.Request) {
	routes, err := api.ms.Routes()
	if err != nil {
		log.WithError(err).Error("unable to get eta")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(reflect.TypeOf(routes))
	// fmt.Println(len(routes))
	// points:=routes[1].Points[1].Latitude
	distances:= make([]float64, len(routes)) 
	for i:=0; i<len(routes);i++{
	 	// fmt.Println(routes[i].ID,routes[i].Name)
	 	points:=routes[i].Points
	 	var dis float64=0
	 	var lat,lon float64
	 	for j:=0; j<len(points);j++{
	 		if j==0{
		 		lat=points[j].Latitude/360*2*math.Pi
		 		lon=points[j].Longitude/360*2*math.Pi
		 	}else{
		 		lat1:=points[j].Latitude/360*2*math.Pi
		 		lon1:=points[j].Longitude/360*2*math.Pi
		 		dlon:=lon1-lon
		 		dlat:=lat1-lat
		 		a:=math.Pow(math.Sin(dlat/2),2) + math.Cos(lat1)*math.Cos(lat)*math.Pow(math.Sin(dlon/2),2)
		 		c := 2*math.Asin(math.Sqrt(a))
		 		dis+=c*6371/1.6
		 		lat=lat1
		 		lon=lon1
		 	}
	 		// fmt.Println(lat,lon)
	 	}
	 	fmt.Println(dis)
	 	distances[i]=dis
	}
	fmt.Println()
	// fmt.Println(distances)
	// fmt.Println(reflect.TypeOf(points))
	WriteJSON(w, distances)
}
