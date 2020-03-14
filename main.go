package main

import (
	// "fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)

var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{}

type Message struct {
	Message string `json:message`
}

func HandleClients(w http.ResponseWriter, r *http.Request) {
	go broadcastMessageToClients()

	websocket, error := upgrader.Upgrade(w, r, nil)

	if error != nil {
		log.Fatal("error", error)
		defer websocket.Close()
	}

	clients[websocket] = true

	for {
		var message Message

		error := websocket.ReadJSON(&message)

		if error != nil {
			log.Printf("error", error)
			delete(clients, websocket)
			break
		}

		broadcast <- message
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "index.html") })

	http.HandleFunc("/chat", HandleClients)
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		log.Fatal("error", error)
	}
}

func broadcastMessageToClients() {
	for {
		message := <-broadcast

		for client := range clients {
			error := client.WriteJSON(message)
			if error != nil {
				log.Fatal("error", error)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
