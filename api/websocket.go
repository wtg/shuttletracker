package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/wtg/shuttletracker/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type fusionPosition struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Speed     *float64 `json:"speed"` // meters per second
	Heading   *float64 `json:"heading"`
	Time      time.Time
}

type fusionClient struct {
	socket    *websocket.Conn
	positions []fusionPosition
}

type fusionManager struct {
	clients []*fusionClient
}

func (fm *fusionManager) addClient(conn *websocket.Conn) {
	client := &fusionClient{
		socket: conn,
	}
	fm.clients = append(fm.clients, client)

	go fm.handleClient(client)
}

func (fm *fusionManager) handleClient(client *fusionClient) {
	conn := client.socket
	for {
		_, r, err := conn.NextReader()
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
		log.Infof("received message %+v", fp)
		client.positions = append(client.positions, fp)
	}
}

func (fm *fusionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fusionManager debug")

	fmt.Fprintf(w, "\n\n%d clients:\n", len(fm.clients))
	for _, client := range fm.clients {
		fmt.Fprintf(w, "%+v\n", client)
	}
}

func (api *API) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("unable to upgrade connection")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.fm.addClient(conn)

}
