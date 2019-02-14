package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"

	"github.com/wtg/shuttletracker/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type fusionPosition struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     *float64  `json:"speed"` // meters per second
	Heading   *float64  `json:"heading"`
	Track     string    `json:"track"`
	Time      time.Time `json:"time"`
}

type fusionClient struct {
	conn *websocket.Conn
}

type fusionManager struct {
	clients         []*fusionClient
	tracks          map[string][]fusionPosition
	newConnChan     chan *websocket.Conn
	newPositionChan chan fusionPosition
}

func newFusionManager() *fusionManager {
	fm := &fusionManager{
		newConnChan:     make(chan *websocket.Conn),
		newPositionChan: make(chan fusionPosition),
		tracks:          map[string][]fusionPosition{},
	}
	go fm.handleNewConns()
	go fm.handleNewPositions()
	return fm
}

func (fm *fusionManager) addClient(conn *websocket.Conn) error {
	client := &fusionClient{
		conn: conn,
	}

	fm.clients = append(fm.clients, client)
	go fm.handleClient(client)
	return nil
}

func (fm *fusionManager) handleNewConns() {
	for conn := range fm.newConnChan {
		err := fm.addClient(conn)
		if err != nil {
			log.WithError(err).Error("unable to add client")
		}
	}
}

func (fm *fusionManager) handleClient(client *fusionClient) {
	for {
		_, r, err := client.conn.NextReader()
		if err != nil {
			log.WithError(err).Error("unable to get reader")
			return
		}
		dec := json.NewDecoder(r)
		fp := fusionPosition{}
		err = dec.Decode(&fp)
		if err != nil {
			log.WithError(err).Error("unable to decode message")
		}
		fp.Time = time.Now()
		// log.Infof("received message %+v", fp)
		fm.newPositionChan <- fp
		// client.positions = append(client.positions, fp)
	}
}

func (fm *fusionManager) handleNewPositions() {
	for pos := range fm.newPositionChan {
		log.Debugf("new position: %+v", pos)
		fm.tracks[pos.Track] = append(fm.tracks[pos.Track], pos)
	}
}

func (fm *fusionManager) debugHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fusionManager debug\n\n")

	fmt.Fprintf(w, "%d tracks\n", len(fm.tracks))

	numPositions := 0
	for _, track := range fm.tracks {
		numPositions += len(track)
	}
	fmt.Fprintf(w, "%d positions\n\n", numPositions)

	fmt.Fprintf(w, "%d clients:\n", len(fm.clients))
	for _, client := range fm.clients {
		fmt.Fprintf(w, "%+v\n", client)
	}
}

func (fm *fusionManager) exportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err := enc.Encode(fm.tracks)
	if err != nil {
		log.WithError(err).Error("unable to encode")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (fm *fusionManager) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("unable to upgrade connection")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fm.newConnChan <- conn
}
func (fm *fusionManager) router(auth func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.HandleFunc("/", fm.webSocketHandler)
	r.With(auth).Get("/debug", fm.debugHandler)
	r.With(auth).Get("/export", fm.exportHandler)
	return r
}
