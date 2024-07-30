package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"github.com/sirupsen/logrus"
)

var Manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}
var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		Logger.Error(err)
		http.Error(res, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		Logger.Error(err)
		http.Error(res, "Failed to generate UUID", http.StatusInternalServerError)
		return
	}
	client := &Client{id: id.String(), socket: conn, send: make(chan []byte)}
	Logger.Infof("New client connected: %s", client.id)

	Manager.register <- client
	Logger.Infof("Client registered: %s", client.id)

	go client.read()
	go client.write()
}

func main() {
	go Manager.start()
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12312", nil)
}
