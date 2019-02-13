package api

import (
	"net/http"
	"github.com/gorilla/websocket"

	"github.com/wtg/shuttletracker/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func (api *API) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("unable to upgrade connection")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.WithError(err).Error("unable to read message")
			return
		}
		log.Infof("received message %s", string(p))
	}
}
