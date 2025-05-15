package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var clientsMutex sync.Mutex
var broadcast = make(chan Message)

var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
	AddSource: true,
}))

type Message struct {
	Username string `json:"username"`
	Text string `json:"text"`
}

var messages []Message

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	r := chi.NewRouter()

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Произошла ошибка при улучшении запроса до websocket")
			return
		}
		defer conn.Close()

		clientsMutex.Lock()
		clients[conn] = true
		clientsMutex.Unlock()

		for _, message := range messages {
			err := conn.WriteJSON(message)
			if err != nil {
				clientsMutex.Lock()
				delete(clients, conn)
				clientsMutex.Unlock()
				log.Error("Ошибка при отправке сообщений новому пользователю")
				return
			}
		}

		for {
			var message Message

			err := conn.ReadJSON(&message)
			if err != nil {
      	log.Error("Ошибка чтения")
				clientsMutex.Lock()
        delete(clients, conn)
        clientsMutex.Unlock()
				break
      }

			messages = append(messages, message)

			broadcast <- message
		}
	})

	go broadcastMessages()
	log.Info("Сервер запущен")
	http.ListenAndServe(":8080", r)
}

func broadcastMessages() {
	for {
		message := <-broadcast

		clientsMutex.Lock()
		for conn := range clients {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Error("Ошибка при отправке сообщения")
				conn.Close()
				delete(clients, conn)
			}
		}
		clientsMutex.Unlock()
	}
}
