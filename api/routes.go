package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/cas.v1"
	"math"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"io/ioutil"

	"github.com/wtg/shuttletracker/model"
	"gopkg.in/mgo.v2/bson"
)

// RoutesHandler finds all of the routes in the database
func (App *API) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all routes in database
	var routes []model.Route
	err := App.db.Routes.Find(bson.M{}).All(&routes)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each route to client as JSON
	WriteJSON(w, routes)
}

// StopsHandler finds all of the route stops in the database
func (App *API) StopsHandler(w http.ResponseWriter, r *http.Request) {
	// Find all stops in databases
	var stops []model.Stop
	err := App.db.Stops.Find(bson.M{}).All(&stops)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each stop to client as JSON
	WriteJSON(w, stops)
}

// Interpolate do interpolation using user input coordinates
func Interpolate(coords []model.Coord, key string) []model.Coord {
	// make request
	prefix := "https://roads.googleapis.com/v1/snapToRoads?"
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString("path=")
	for i, coord := range coords {
		buffer.WriteString(strconv.FormatFloat(coord.Lat, 'f', 10, 64))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatFloat(coord.Lng, 'f', 10, 64))
		if i < len(coords)-1 {
			buffer.WriteString("|")
		}
	}
	// add the first point to be evaluated
	if len(coords) > 1 {
		buffer.WriteString("|" + strconv.FormatFloat(coords[0].Lat, 'f', 10, 64) + "," + strconv.FormatFloat(coords[1].Lng, 'f', 10, 64))
	}
	buffer.WriteString("&interpolate=true&key=")
	buffer.WriteString(key)
	fmt.Println(buffer.String())
	// send request
	resp, err := http.Get(buffer.String())
	if err != nil {
		fmt.Errorf("Error Not valid response from Google API")
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	mapResponse := model.MapResponse{}
	json.Unmarshal(body, &mapResponse)
	// read response
	result := []model.Coord{}
	for _, location := range mapResponse.SnappedPoints {
		currentLocation := model.Coord{
			Lat: float64(location.Location.Latitude),
			Lng: float64(location.Location.Longitude),
		}
		result = append(result, currentLocation)
	}
	return result
}

func GoogleSegmentCompute(from model.Coord, to model.Coord, key string) model.Segment {
	prefix := "https://maps.googleapis.com/maps/api/distancematrix/json?units=imperial&"
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	origin := fmt.Sprintf("origins=%f,%f", from.Lat, from.Lng)
	destination := fmt.Sprintf("destinations=%f,%f", to.Lat, to.Lng)
	buffer.WriteString(origin + "&" + destination + "&key=" + key)
	fmt.Println(buffer.String())
	resp, err := http.Get(buffer.String())
	if err != nil {
		fmt.Errorf("Error Not valid response from Google API")
		return model.Segment{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	mapResponse := model.MapDistanceMatrixResponse{}
	json.Unmarshal(body, &mapResponse)
	fmt.Println(mapResponse)
	result := model.Segment{
		Start: model.MapPoint{
			Latitude:  float64(from.Lat),
			Longitude: float64(from.Lng),
		},
		End: model.MapPoint{
			Latitude:  float64(to.Lat),
			Longitude: float64(to.Lng),
		},
		Distance: float64(mapResponse.Rows[0].Elements[0].Distance.Value),
		Duration: float64(mapResponse.Rows[0].Elements[0].Duration.Value),
	}
	fmt.Println(result)
	return result
}

// compute distance between two coordinates and return a value
func ComputeDistance(c1 model.Coord, c2 model.Coord) float64 {
	return float64(math.Sqrt(math.Pow(c1.Lat-c2.Lat, 2) + math.Pow(c1.Lng-c2.Lng, 2)))
}

func ComputeDistanceMapPoint(c1 model.MapPoint, c2 model.MapPoint) float64 {
	return float64(math.Sqrt(math.Pow(c1.Latitude-c2.Latitude, 2) + math.Pow(c1.Longitude-c2.Longitude, 2)))
}

// Compute the Segment for each segment of the coordinates
func ComputeSegments(coords []model.Coord, key string, threshold int) []model.Segment {
	result := []model.Segment{}
	// only compute the distance greater than some theshold distance and assume all in between has the same Segment
	prev := 0
	index := 1
	// This part could be improved by rewriting with asynchronized call
	for index = 1; index < len(coords); index++ {
		if index%threshold == 0 {
			v := GoogleSegmentCompute(coords[prev], coords[index], key)
			for inner := prev + 1; inner <= index; inner++ {
				result = append(result, model.Segment{
					Distance: v.Distance / float64(index-prev),
					Duration: v.Duration / float64(index-prev),
					Start:    model.MapPoint{Latitude: float64(coords[inner-1].Lat), Longitude: float64(coords[inner-1].Lng)},
					End:      model.MapPoint{Latitude: float64(coords[inner].Lat), Longitude: float64(coords[inner].Lng)},
				})
			}
			prev = index
		}
	}
	return result
}

//This is really a temporary funcion to import the old database, only supports adding two routes
func (App *API) ImportHandler(w http.ResponseWriter, r *http.Request) {
	var count int
	count = 0
	db, err := sql.Open("mysql", "root:pass@/shuttle_tracking")
	//Begin connecting to database
	if err != nil {
		log.Fatalf("Couldnt connect to mysql")
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	if db == nil {
		log.Fatalf("db empty")
	}
	//Begin grabbing information we need
	rows, err := db.Query("SELECT * FROM coords")
	if err != nil {
		log.Fatalf("bad query")
	}
	//iterate through rows
	coords := []model.Coord{}
	var oldId int
	oldId = 1
	for rows.Next() {
		var id int
		var lat float64
		var long float64
		var position int
		var route_id int
		var created_at string
		var updated_at string

		if err := rows.Scan(&id, &lat, &long, &position, &route_id, &created_at, &updated_at); err != nil {
			log.Fatal(err)
		}
		//We're done with the first route, update it and put it in the database.
		if oldId == 1 && route_id == 2 {

			coords = Interpolate(coords, App.cfg.GoogleMapAPIKey)
			segments := ComputeSegments(coords, App.cfg.GoogleMapAPIKey, App.cfg.GoogleMapMinDistance)

			route := db.QueryRow("SELECT name,description,start_time,end_time,color FROM routes where id = 1")
			var name string
			var desc string
			var color string
			var start_time string
			var end_time string
			err := route.Scan(&name, &desc, &start_time, &end_time, &color)
			if err != nil {
				log.Fatal(err)
			}

			newRoute := model.Route{
				ID:          bson.NewObjectId().Hex(),
				Name:        name,
				Description: desc,
				StartTime:   start_time,
				EndTime:     end_time,
				Enabled:     true,
				Color:       color,
				Width:       4,
				Coords:      coords,
				Duration:    segments,
				Created:     time.Now(),
				Updated:     time.Now()}
			_ = newRoute
			fmt.Printf(name)
			coords = nil
			err = App.db.Routes.Insert(&newRoute)
		}

		oldId = route_id
		myCoord := model.Coord{lat, long}
		if route_id == 1 {
			coords = append(coords, myCoord)
		} else {
			count += 1
			if count%10 == 0 {
				coords = append(coords, myCoord)
			}
		}
	}

	coords = Interpolate(coords, App.cfg.GoogleMapAPIKey)
	segments := ComputeSegments(coords, App.cfg.GoogleMapAPIKey, App.cfg.GoogleMapMinDistance)

	route := db.QueryRow("SELECT name,description,start_time,end_time,color FROM routes where id = 2")
	var name string
	var desc string
	var color string
	var start_time string
	var end_time string
	err = route.Scan(&name, &desc, &start_time, &end_time, &color)
	if err != nil {
		log.Fatal(err)
	}

	newRoute := model.Route{
		ID:          bson.NewObjectId().Hex(),
		Name:        name,
		Description: desc,
		StartTime:   start_time,
		EndTime:     end_time,
		Enabled:     true,
		Color:       color,
		Width:       4,
		Coords:      coords,
		Duration:    segments,
		Created:     time.Now(),
		Updated:     time.Now()}
	_ = newRoute
	fmt.Printf(name)
	coords = []model.Coord{}
	err = App.db.Routes.Insert(&newRoute)

}

// RoutesCreateHandler adds a new route to the database
func (App *API) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new route object using request fields
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}
	var routeData map[string]string
	var coordsData []map[string]float64
	// Decode route details
	err := json.NewDecoder(r.Body).Decode(&routeData)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Unmarshal route coordinates
	err = json.Unmarshal([]byte(routeData["coords"]), &coordsData)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Create a Coord from each set of input coordinates
	coords := []model.Coord{}
	for _, c := range coordsData {
		coord := model.Coord{c["lat"], c["lng"]}
		coords = append(coords, coord)
	}
	// Here do the interpolation
	coords = Interpolate(coords, App.cfg.GoogleMapAPIKey)
	segments := ComputeSegments(coords, App.cfg.GoogleMapAPIKey, App.cfg.GoogleMapMinDistance)
	// now we get the Segment for each segment ( this should be stored in database, just store it inside route for god sake)

	fmt.Printf("Size of coordinates = %d", len(coords))
	// Type conversions
	enabled, _ := strconv.ParseBool(routeData["enabled"])
	width, _ := strconv.Atoi(routeData["width"])
	currentTime := time.Now()
	// Create a new route
	route := model.Route{
		ID:          bson.NewObjectId().Hex(),
		Name:        routeData["name"],
		Description: routeData["description"],
		StartTime:   routeData["startTime"],
		EndTime:     routeData["endTime"],
		Enabled:     enabled,
		Color:       routeData["color"],
		Width:       width,
		Coords:      coords,
		Duration:    segments,
		Created:     currentTime,
		Updated:     currentTime}
	// Store new route under routes collection
	err = App.db.Routes.Insert(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//Deletes route from database
func (App *API) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	vars := mux.Vars(r)
	fmt.Printf(vars["id"])
	log.Debugf("deleting", vars["id"])
	err := App.db.Routes.Remove(bson.M{"id": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this could be improved
func GetSegment(stop model.Stop, startIndex int, segments []model.Segment) int {
	// choose the segment with lowest distance
	fmt.Println(len(segments))
	x0 := float64(stop.Lat)
	y0 := float64(stop.Lng)
	minimumLen := 1000.0
	minimumIndex := startIndex
	for i := startIndex; i < len(segments); i++ {
		x1 := segments[i].Start.Latitude
		y1 := segments[i].Start.Longitude
		x2 := segments[i].End.Latitude
		y2 := segments[i].End.Longitude
		// compute the distance between a point and a line
		length := math.Abs((x2-x1)*(y1-y0)-(x1-x0)*(y2-y1)) / math.Sqrt(math.Pow(x2-x1, 2)+math.Pow(y2-y1, 2))
		if length < minimumLen {
			minimumLen = length
			minimumIndex = i
		}
	}
	return minimumIndex
}

// StopsCreateHandler adds a new route stop to the database
func (App *API) StopsCreateHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	fmt.Print("Create Stop Handler called")
	// Create a new stop object using request fields
	stop := model.Stop{}
	err := json.NewDecoder(r.Body).Decode(&stop)
	stop.ID = bson.NewObjectId().Hex()
	route := model.Route{}
	err1 := App.db.Routes.Find(bson.M{"id": stop.RouteID}).One(&route)
	// Error handling

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// We have to know the order of the stop and store a velocity vector into duration for the prediction
	route.StopsID = append(route.StopsID, stop.ID) // THIS REQUIRES the front end to have correct order << to be improved
	fmt.Println(route.StopsID)
	stop.SegmentIndex = GetSegment(stop, route.AvailableRoute, route.Duration)

	// Store new stop under stops collection
	err = App.db.Stops.Insert(&stop)
	// Error handling
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query := bson.M{"id": stop.RouteID}
	change := bson.M{"$set": bson.M{"availableroute": stop.SegmentIndex + 1, "stopsid": route.StopsID}}

	err = App.db.Routes.Update(query, change)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Closest Segment ID = " + strconv.Itoa(stop.SegmentIndex))
	WriteJSON(w, stop)
}

func (App *API) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if App.cfg.Authenticate && !cas.IsAuthenticated(r) {
		return
	}

	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	fmt.Printf(vars["id"])
	err := App.db.Stops.Remove(bson.M{"id": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
