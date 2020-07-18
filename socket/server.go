package server

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexanderkarlis/naive-blockchain/block"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline      = []byte{'\n'}
	space        = []byte{' '}
	upgrader     = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	addr         = flag.String("addr", "localhost:8080", "http service address")
	genesisBlock = block.CreateGenesisBlock()
	blockchain   = block.Blockchain{}
	msgs         = []string{}
)

// Client is the middleman between client and server.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// SocketServe function serves a basic form of the listening websocket.
func SocketServe() {
	blockchain = append(blockchain, *genesisBlock)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	flag.Parse()
	log.SetFlags(0)
	hub := newHub()
	go hub.run()
	go func() {
		http.HandleFunc("/helloServer", helloServer)
		log.Fatal(http.ListenAndServe("localhost:8081", nil))
	}()
	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		go func() {
			msgs = append(msgs, string(message))
			log.Println("\nmsgs⏎")
			for x := range msgs {
				log.Printf("%v\n", msgs[x])
			}
		}()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		blockchain.AddNewBlockToBlockChain(string(message))
		valid, index := blockchain.IsValidBlockChain()
		if !valid {
			panic(fmt.Sprintf("WARNING: Blockchain not valid at %+d", index))
		}
		bSlice, _ := json.Marshal(blockchain)
		c.hub.broadcast <- bSlice
	}
}

func (c *Client) broadcast() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.broadcast()
	go client.readPump()
}
