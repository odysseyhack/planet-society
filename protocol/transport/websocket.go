package transport

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/odysseyhack/planet-society/protocol/protocol"
)

type Websocket struct {
	upgrader   *websocket.Upgrader
	connection chan protocol.Conn
}

func NewWebsocket(connection chan protocol.Conn) *Websocket {
	return &Websocket{
		connection: connection,
		upgrader:   &websocket.Upgrader{},
	}
}

func (ws *Websocket) Listen(addr string) error {
	router := mux.NewRouter()
	router.HandleFunc("/", ws.client)

	server := &http.Server{
		Addr: addr,
		// todo: review timeouts
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (ws *Websocket) client(w http.ResponseWriter, r *http.Request) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("websocket upgrade failed:", err)
		return
	}
	ws.connection <- &Conn{Conn: c}
}

type Conn struct {
	Conn *websocket.Conn
}

func (c *Conn) Read() (*protocol.Message, error) {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var protocolMessage protocol.Message

	if err := gob.NewDecoder(bytes.NewReader(msg)).Decode(&protocolMessage); err != nil {
		return nil, err
	}

	return &protocolMessage, nil

}

func (c *Conn) Write(msg *protocol.Message) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(msg); err != nil {
		return err
	}

	return c.Conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}
