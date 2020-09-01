package interactiveserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	register <- conn
}

var reload = make(chan bool)
var register = make(chan *websocket.Conn)
var conns = make(map[*websocket.Conn]bool)

func eventLoop() {
	for {
		select {
		case conn := <-register:
			conns[conn] = true
		case <-reload:
			for conn := range conns {
				err := conn.WriteMessage(websocket.TextMessage, []byte("reload"))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// TriggerReload will trigger a reload of all clients connected to the interactive server
func TriggerReload() {
	reload <- true
}

// Start the interactive web server
func Start() error {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", index)
	r.HandleFunc("/ws", websocketHandler)
	r.HandleFunc("/reports/", index)
	r.HandleFunc("/reports/{name}/", report)

	http.Handle("/", r)

	go eventLoop()

	fmt.Println("Starting interactive server on :16772")
	err := http.ListenAndServe(":16772", nil)
	return err
}
