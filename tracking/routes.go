package tracking

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"fmt"
	"io/ioutil"

	"gopkg.in/mgo.v2/bson"
)

// Coord represents a single lat/lng point used to draw routes
type Coord struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// Route represents a set of coordinates to draw a path on our tracking map
type Route struct {
	ID             string    `json:"id"             bson:"id"`
	Name           string    `json:"name"           bson:"name"`
	Description    string    `json:"description"    bson:"description"`
	StartTime      string    `json:"startTime"      bson:"startTime"`
	EndTime        string    `json:"endTime"        bson:"endTime"`
	Enabled        bool      `json:"enabled,string" bson:"enabled"`
	Color          string    `json:"color"          bson:"color"`
	Width          int       `json:"width,string"   bson:"width"`
	Coords         []Coord   `json:"coords"         bson:"coords"`
	Duration       []Segment `json:"duration"       bson:"duration"`
	StopsID        []string  `json:"stopsid"        bson:"stopsid"`
	AvailableRoute int       `json:"availableroute" bson:"availableroute"`
	Created        time.Time `json:"created"        bson:"created"`
	Updated        time.Time `json:"updated"        bson:"updated"`
}

// Stop indicates where a tracked object is scheduled to arrive
type Stop struct {
	ID           string  `json:"id"             bson:"id"`
	Name         string  `json:"name"           bson:"name"`
	Description  string  `json:"description"    bson:"description"`
	Lat          float64 `json:"lat,string"     bson:"lat"`
	Lng          float64 `json:"lng,string"     bson:"lng"`
	Address      string  `json:"address"        bson:"address"`
	StartTime    string  `json:"startTime"      bson:"startTime"`
	EndTime      string  `json:"endTime"        bson:"endTime"`
	Enabled      bool    `json:"enabled,string" bson:"enabled"`
	RouteID      string  `json:"routeId"        bson:"routeId"`
	SegmentIndex int     `json:"segmentindex"   bson:"segmentindex"`
}

type MapPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type MapResponsePoint struct {
	Location      MapPoint `json:"location"`
	OriginalIndex int      `json:"originalIndex,omitempty"`
	PlaceID       string   `json:"placeId"`
}
type MapResponse struct {
	SnappedPoints []MapResponsePoint
}

type MapDistanceMatrixDuration struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type MapDistanceMatrixDistance struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

type MapDistanceMatrixElement struct {
	Status   string                    `json:"status"`
	Duration MapDistanceMatrixDuration `json:"duration"`
	Distance MapDistanceMatrixDistance `json:"distance"`
}

type MapDistanceMatrixElements struct {
	Elements []MapDistanceMatrixElement `json:"elements"`
}
type MapDistanceMatrixResponse struct {
	Status               string                      `json:"status"`
	OriginAddresses      []string                    `json:"origin_addresses"`
	DestinationAddresses []string                    `json:"destination_addresses"`
	Rows                 []MapDistanceMatrixElements `json:"rows"`
}

type Segment struct {
	ID       string   `json:"id"`
	Start    MapPoint `json:"origin"`
	End      MapPoint `json:"destination"`
	Distance float64  `json:"distance"`
	Duration float64  `json:"duration"`
}

// RoutesHandler finds all of the routes in the database
func (App *App) RoutesHandler(w http.ResponseWriter, r *http.Request) {
	// Find all routes in database
	var routes []Route
	err := App.Routes.Find(bson.M{}).All(&routes)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each route to client as JSON
	WriteJSON(w, routes)
}

// StopsHandler finds all of the route stops in the database
func (App *App) StopsHandler(w http.ResponseWriter, r *http.Request) {
	// Find all stops in database
	var stops []Stop
	err := App.Stops.Find(bson.M{}).All(&stops)
	// Handle query errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Send each stop to client as JSON
	WriteJSON(w, stops)
}

func Interpolate(coords []Coord, key string) []Coord {
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
	buffer.WriteString("&interpolate=true&key=")
	buffer.WriteString(key)
	resp, err := http.Get(buffer.String())
	if err != nil {
		fmt.Errorf("Error Not valid response from Google API")
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	mapResponse := MapResponse{}
	json.Unmarshal(body, &mapResponse)
	result := []Coord{}
	for _, location := range mapResponse.SnappedPoints {
		currentLocation := Coord{
			Lat: float64(location.Location.Latitude),
			Lng: float64(location.Location.Longitude),
		}
		result = append(result, currentLocation)
	}
	return result
}

func GoogleSegmentCompute(from Coord, to Coord, key string) Segment {
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
		return Segment{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	mapResponse := MapDistanceMatrixResponse{}
	json.Unmarshal(body, &mapResponse)
	fmt.Println(mapResponse)
	result := Segment{
		Start: MapPoint{
			Latitude:  float64(from.Lat),
			Longitude: float64(from.Lng),
		},
		End: MapPoint{
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
func ComputeDistance(c1 Coord, c2 Coord) float64 {
	return float64(math.Sqrt(math.Pow(c1.Lat-c2.Lat, 2) + math.Pow(c1.Lng-c2.Lng, 2)))
}

func ComputeDistanceMapPoint(c1 MapPoint, c2 MapPoint) float64 {
	return float64(math.Sqrt(math.Pow(c1.Latitude-c2.Latitude, 2) + math.Pow(c1.Longitude-c2.Longitude, 2)))
}

// Compute the Segment for each segment of the coordinates
func ComputeSegments(coords []Coord, key string, threshold int) []Segment {
	result := []Segment{}
	// only compute the distance greater than some theshold distance and assume all in between has the same Segment
	prev := 0
	index := 1
	// This part could be improved by rewriting with asynchronized call
	for index = 1; index < len(coords); index++ {
		if index%threshold == 0 {
			v := GoogleSegmentCompute(coords[prev], coords[index], key)
			for inner := prev + 1; inner <= index; inner++ {
				result = append(result, Segment{
					Distance: v.Distance / float64(index-prev),
					Duration: v.Duration / float64(index-prev),
					Start:    MapPoint{Latitude: float64(coords[inner-1].Lat), Longitude: float64(coords[inner-1].Lng)},
					End:      MapPoint{Latitude: float64(coords[inner].Lat), Longitude: float64(coords[inner].Lng)},
				})
			}
			prev = index
		}
	}
	// compute the last segment
	result = append(result, GoogleSegmentCompute(coords[prev], coords[0], key))
	return result
}

// RoutesCreateHandler adds a new route to the database
func (App *App) RoutesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new route object using request fields
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
	coords := []Coord{}
	for _, c := range coordsData {
		coord := Coord{c["lat"], c["lng"]}
		coords = append(coords, coord)
	}
	// Here do the interpolation
	coords = Interpolate(coords, App.Config.GoogleMapAPIKey)
	segments := ComputeSegments(coords, App.Config.GoogleMapAPIKey, App.Config.GoogleMapMinDistance)
	// now we get the Segment for each segment ( this should be stored in database, just store it inside route for god sake)

	fmt.Printf("Size of coordinates = %d", len(coords))
	// Type conversions
	enabled, _ := strconv.ParseBool(routeData["enabled"])
	width, _ := strconv.Atoi(routeData["width"])
	currentTime := time.Now()
	// Create a new route
	route := Route{
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
	err = App.Routes.Insert(&route)
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Deletes route from database
func (App *App) RoutesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	err := App.Routes.Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"])})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this could be improved
func GetSegment(stop Stop, startIndex int, segments []Segment) int {
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
func (App *App) StopsCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Create Stop Handler called")
	// Create a new stop object using request fields
	stop := Stop{}
	err := json.NewDecoder(r.Body).Decode(&stop)
	stop.ID = bson.NewObjectId().Hex()
	route := Route{}
	err1 := App.Routes.Find(bson.M{"id": stop.RouteID}).One(&route)
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
	err = App.Stops.Insert(&stop)
	// Error handling
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query := bson.M{"id": stop.RouteID}
	change := bson.M{"$set": bson.M{"availableroute": stop.SegmentIndex + 1, "stopsid": route.StopsID}}

	err = App.Routes.Update(query, change)
	if err != nil {
		fmt.Println(err.Error())
	}
	// check if update success
	fmt.Println("Closest Segment ID = " + strconv.Itoa(stop.SegmentIndex))
	fmt.Println(route.Duration[stop.SegmentIndex])
	test := Route{}
	App.Routes.Find(bson.M{"id": stop.RouteID}).One(&test)
	fmt.Println(test.StopsID)
	fmt.Println(test.AvailableRoute)
	// When creating a Stop, it actually changes the segments by adding a segment in a route, which the segment will take a different duration
	// select the route pointed by the stop
}

func (App *App) StopsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debugf("deleting", vars["id"])
	err := App.Stops.Remove(bson.M{"name": vars["id"]})
	// Error handling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
