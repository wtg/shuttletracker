package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type fusionMessageEnvelope struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type fusionPosition struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     *float64  `json:"speed"` // meters per second
	Heading   *float64  `json:"heading"`
	Track     string    `json:"track"`
	Time      time.Time `json:"time"`
}

type fusionBusButton struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"latitude"`
}

type fusionClient struct {
	conn *websocket.Conn
}

type fusionManager struct {
	clients          []*fusionClient
	tracks           map[string][]fusionPosition
	newConnChan      chan *websocket.Conn
	removeClientChan chan *fusionClient
	positionChan     chan fusionPosition
	busButtonChan    chan fusionBusButton
}

func newFusionManager() *fusionManager {
	fm := &fusionManager{
		newConnChan:      make(chan *websocket.Conn),
		removeClientChan: make(chan *fusionClient),
		positionChan:     make(chan fusionPosition),
		busButtonChan:    make(chan fusionBusButton),
		tracks:           map[string][]fusionPosition{},
	}
	go fm.clientsLoop()
	go fm.messagesLoop()
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

func (fm *fusionManager) removeClient(client *fusionClient) error {
	for i, c := range fm.clients {
		if client == c {
			fm.clients = append(fm.clients[:i], fm.clients[i+1:]...)
			return nil
		}
	}
	return errors.New("client not found")
}

// clientsLoop selects over the channels related to clients appearing or going away:
// newConnChan and removeClientChan.
func (fm *fusionManager) clientsLoop() {
	for {
		select {
		case conn := <-fm.newConnChan:
			err := fm.addClient(conn)
			if err != nil {
				log.WithError(err).Error("unable to add client")
			}
		case client := <-fm.removeClientChan:
			err := fm.removeClient(client)
			if err != nil {
				log.WithError(err).Error("umable to remove client")
			}
		}
	}
}

func decodeFusionMessage(r io.Reader) (string, json.RawMessage, error) {
	var message json.RawMessage
	fm := fusionMessageEnvelope{
		Message: &message,
	}
	dec := json.NewDecoder(r)
	err := dec.Decode(&fm)
	if err != nil {
		return "", message, err
	}
	return fm.Type, message, nil
}

func (fm *fusionManager) handleClient(client *fusionClient) {
	for {
		_, r, err := client.conn.NextReader()
		if err != nil {
			log.WithError(err).Error("unable to get reader")
			break
		}
		messageType, message, err := decodeFusionMessage(r)
		if err != nil {
			log.WithError(err).Error("unable to decode message")
			break
		}
		switch messageType {
		case "position":
			fp := fusionPosition{}
			err = json.Unmarshal(message, &fp)
			if err != nil {
				log.WithError(err).Error("unable to decode fusionPosition")
				break
			}
			fp.Time = time.Now()
			fm.positionChan <- fp
		case "bus_button":
			fbb := fusionBusButton{}
			err = json.Unmarshal(message, &fbb)
			if err != nil {
				log.WithError(err).Error("unable to decode fusionBusButton")
				break
			}
			fm.busButtonChan <- fbb
		default:
			log.WithError(err).Errorf("unknown message type \"%s\"", messageType)
		}
	}

	// remove client since the connection is dead
	fm.removeClientChan <- client
}

func (fm *fusionManager) messagesLoop() {
	for {
		select {
		case pos := <-fm.positionChan:
			log.Debugf("new position: %+v", pos)
			fm.tracks[pos.Track] = append(fm.tracks[pos.Track], pos)
		case bb := <-fm.busButtonChan:
			log.Debugf("new bus button: %+v", bb)
		}
	}
}

func (fm *fusionManager) debugHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "fusionManager debug\n\n")
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}

	_, err = fmt.Fprintf(w, "%d tracks\n", len(fm.tracks))
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}

	numPositions := 0
	for _, track := range fm.tracks {
		numPositions += len(track)
	}
	_, err = fmt.Fprintf(w, "%d positions\n\n", numPositions)
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}

	_, err = fmt.Fprintf(w, "%d clients:\n", len(fm.clients))
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}
	for _, client := range fm.clients {
		_, err = fmt.Fprintf(w, "%+v\n", client)
		if err != nil {
			log.WithError(err).Error("unable to write response")
			return
		}
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
