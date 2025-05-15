package wsserver

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Username string `json:"username"`
	Text string `json:"text"`
}


type Client struct{
	Conn *websocket.Conn
	Send chan Message
}

type WSSrv interface {
	Start() error
}

type wssrv struct {
	log *slog.Logger

	mutex sync.Mutex
	clients map[*Client]bool
}

func New(log *slog.Logger) WSSrv{
	return &wssrv{
		log: log,
		clients: make(map[*Client]bool),
	}
}

func (ws *wssrv) Start() error {
	router := chi.NewRouter()

	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ws.log.Error("Произошла ошибка при улучшении запроса до websocket")
			return
		}
		defer conn.Close()

		client := &Client{
			Conn: conn,
			Send: make(chan Message),
		}

		ws.mutex.Lock()
		ws.clients[client] = true
		ws.mutex.Unlock()

		defer func(){
			ws.mutex.Lock()
			defer ws.mutex.Unlock()

			delete(ws.clients, client)
			close(client.Send)
		}()

	})

	return server.ListenAndServe()
}
