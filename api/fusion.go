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

// Messages from clients must be in this envelope. Depending on Type, fusionManager
// unmarshals Message into the associated type of struct. fusionManager also uses
// this struct to send messages to clients.
type fusionMessageEnvelope struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type fusionMessageSubscribe struct {
	Topic string `json:"topic"`
}

type fusionPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	// Meters per second. Yes, this is different from shuttletracker.Location,
	// which is in miles per hour...
	// It's a pointer because it's often unknown and therefore nil.
	Speed *float64 `json:"speed"`

	// Pointer because it may be unknown.
	Heading *float64 `json:"heading"`

	// Client-provided UUID that associates a list of positions to form a track.
	Track string `json:"track"`

	// Time is when fusionManager receives the position. We don't want to trust
	// the client's timestamp.
	Time time.Time `json:"time"`
}

type fusionBusButton struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type fusionClient struct {
	conn            *websocket.Conn
	lastMessageTime time.Time
	userAgent       string
}

type fusionManagerDebug struct {
	// subscriptions    []sub
	clients        []fusionClient
	tracks         [][]fusionPosition
	busButtonCount uint64
}

type fusionManager struct {
	addClientChan    chan *fusionClient
	removeClientChan chan *fusionClient

	clientMsg chan interface{}

	// This is a little gnarly... basically we can ask fusionManager to send some
	// information about itself to a channel so that we don't have to put its internal
	// state behind a mutex to inspect it. No locks around maps or slices required.
	debug chan chan *fusionManagerDebug

	// Everything after this is considered internal state. Only fm.run will read
	// or modify these fields, and it is considered the owner of this state.

	// clients can subscribe to topics that they're interested in
	subs map[string][]*fusionClient

	clients        []*fusionClient
	tracks         map[string][]fusionPosition
	busButtonCount uint64
}

func newFusionManager() *fusionManager {
	fm := &fusionManager{
		addClientChan:    make(chan *fusionClient),
		removeClientChan: make(chan *fusionClient),
		clientMsg:        make(chan interface{}),
		debug:            make(chan chan *fusionManagerDebug),
		tracks:           map[string][]fusionPosition{},
		subs:             map[string][]*fusionClient{},
	}
	go fm.run()
	return fm
}

func (fm *fusionManager) addClient(client *fusionClient) error {
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

// Select handle client connections, disconnections, and messages.
// Responsible (along with any methods it calls) for managing fusionManager state.
// Anything run calls should obtain the lock on fusionManager state.
func (fm *fusionManager) run() {
	for {
		select {
		case c := <-fm.addClientChan:
			err := fm.addClient(c)
			if err != nil {
				log.WithError(err).Error("unable to add client")
			}
		case client := <-fm.removeClientChan:
			err := fm.removeClient(client)
			if err != nil {
				log.WithError(err).Error("umable to remove client")
			}
		case msg := <-fm.clientMsg:
			fm.processMessage(msg)
		case debugChan := <-fm.debug:
			fm.processDebug(debugChan)
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

// handleClient is expected to be called inside of a goroutine associated with a client.
// It does not directly manipulate fusionManager state—this is done by sending messages
// through a chan that is read elsewhere. We do as much JSON parsing here as possible
// since each connection is handled concurrently.
func (fm *fusionManager) handleClient(client *fusionClient) {
	for {
		_, r, err := client.conn.NextReader()
		if err != nil {
			// did the client e.g. close the tab? then we expect a normal error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.WithError(err).Error("unable to get reader")
			}
			break
		}
		client.lastMessageTime = time.Now()
		messageType, message, err := decodeFusionMessage(r)
		if err != nil {
			log.WithError(err).Error("unable to decode message")
			continue
		}

		switch messageType {
		case "subscribe":
			fms := fusionMessageSubscribe{}
			err = json.Unmarshal(message, &fms)
			if err != nil {
				log.WithError(err).Error("unable to decode fusionMessageSubscribe")
				break
			}
			fm.clientMsg <- fms
		case "position":
			fp := fusionPosition{}
			err = json.Unmarshal(message, &fp)
			if err != nil {
				log.WithError(err).Error("unable to decode fusionPosition")
				break
			}
			fp.Time = time.Now()
			fm.clientMsg <- fp
		case "bus_button":
			fbb := fusionBusButton{}
			err = json.Unmarshal(message, &fbb)
			if err != nil {
				log.WithError(err).Error("unable to decode fusionBusButton")
				break
			}
			fm.clientMsg <- fbb
		default:
			// This is just a warning and not an error since messageType comes straight
			// from the client. We can't trust it.
			log.WithError(err).Warnf("unknown message type \"%s\"", messageType)
		}
	}

	// remove client since the connection is dead
	fm.removeClientChan <- client
}

// processMessage handles messages from clients after they are parsed. it does not
// need any locks or mutexes since it is only called from the goroutine that "owns"
// the state inside of fusionManager.
func (fm *fusionManager) processMessage(msg interface{}) {
	switch t := msg.(type) {
	case fusionMessageSubscribe:
		fms := msg.(fusionMessageSubscribe)
		fm.handleMsgSubscribe(fms)
	case fusionPosition:
		fp := msg.(fusionPosition)
		fm.handleMsgPosition(fp)
	case fusionBusButton:
		fbb := msg.(fusionBusButton)
		fm.handleMsgBusButton(fbb)
	default:
		// This is an error since it means that an unhandled message type was sent to
		// the channel, probably by handleClient. This shouldn't happen, so please fix
		// it if it does (make sure all possible message types are being handled).
		log.Errorf("unknown message type \"%s\"", t)
	}
}

func (fm *fusionManager) handleMsgSubscribe(fms fusionMessageSubscribe) {
	log.Debugf("new subscribe: %+v", fms)
}

func (fm *fusionManager) handleMsgPosition(fp fusionPosition) {
	fp.Time = time.Now()
	log.Debugf("new position: %+v", fp)
	fm.tracks[fp.Track] = append(fm.tracks[fp.Track], fp)
}

func (fm *fusionManager) handleMsgBusButton(fbb fusionBusButton) {
	log.Debugf("new bus button: %+v", fbb)
	fm.busButtonCount++
	fme := fusionMessageEnvelope{
		Type:    "bus_button",
		Message: fbb,
	}
	b, err := json.Marshal(fme)
	if err != nil {
		log.WithError(err).Error("unable to marshal")
		return
	}
	for _, client := range fm.clients {
		err = client.conn.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.WithError(err).Error("unable to write")
		}
	}
}

func (fm *fusionManager) processDebug(ch chan *fusionManagerDebug) {
	// assemble the data...
	debug := &fusionManagerDebug{
		clients:        make([]fusionClient, len(fm.clients)),
		tracks:         make([][]fusionPosition, 0, len(fm.tracks)),
		busButtonCount: fm.busButtonCount,
	}

	for i, c := range fm.clients {
		debug.clients[i] = fusionClient{
			// don't copy the websocket conn
			lastMessageTime: c.lastMessageTime,
			userAgent:       c.userAgent,
		}
	}

	for _, v := range fm.tracks {
		newTrack := make([]fusionPosition, len(v))
		copy(newTrack, v)
		debug.tracks = append(debug.tracks, newTrack)
	}

	// send it 📬
	ch <- debug
}

func (fm *fusionManager) debugInfo() *fusionManagerDebug {
	copyChan := make(chan *fusionManagerDebug)
	fm.debug <- copyChan
	return <-copyChan
}

// debugHandler gets a copy of fusionManager's state and then writes some interesting
// information to an HTTP request. This handler (and any modifications you're thinking
// of making to it) MUST NOT perform any operations on fusionManager's state. In order
// to avoid data races, use the copy.
func (fm *fusionManager) debugHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "fusionManager debug\n\n")
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}

	// ask fusionManager for debug info
	fmDebug := fm.debugInfo()

	_, err = fmt.Fprintf(w, "%d tracks\n", len(fmDebug.tracks))
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}

	numPositions := 0
	for _, track := range fmDebug.tracks {
		numPositions += len(track)
	}
	_, err = fmt.Fprintf(w, "%d positions\n", numPositions)
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}

	_, err = fmt.Fprintf(w, "%d bus buttons\n\n", fmDebug.busButtonCount)
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}

	_, err = fmt.Fprintf(w, "%d clients:\n", len(fmDebug.clients))
	if err != nil {
		log.WithError(err).Error("unable to write response")
		return
	}
	for _, client := range fmDebug.clients {
		_, err = fmt.Fprintf(w, "%s\t%s\n", client.lastMessageTime.Format(time.RFC3339), client.userAgent)
		if err != nil {
			log.WithError(err).Error("unable to write response")
			return
		}
	}
}

func (fm *fusionManager) exportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	fmDebug := fm.debugInfo()
	err := enc.Encode(fmDebug.tracks)
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

	c := &fusionClient{
		conn:            conn,
		lastMessageTime: time.Now(),
		userAgent:       r.UserAgent(),
	}
	fm.addClientChan <- c
}
func (fm *fusionManager) router(auth func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.HandleFunc("/", fm.webSocketHandler)
	r.With(auth).Get("/debug", fm.debugHandler)
	r.With(auth).Get("/export", fm.exportHandler)
	return r
}
