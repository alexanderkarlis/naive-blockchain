package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	newline       = []byte{'\n'}
	space         = []byte{' '}
	upgrader      = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	addr          = flag.String("addr", ":8080", "http service address")
	genesisBlock  = block.CreateGenesisBlock()
	blockchain    = block.Blockchain{}
	msgs          = []string{}
	msgBlockchain = []block.Block{}
	a             = make(chan []string)
	clientList    = []Client{}
)

// Client is the middleman between client and server.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// BlockRequestData is a http request to add a block to the blockchain.
type BlockRequestData struct {
	Data string `json:"data"`
}

// Serve function serves a basic form of the listening websocket.
func Serve() {
	blockchain = append(blockchain, *genesisBlock)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	flag.Parse()
	log.SetFlags(0)
	hub := newHub()
	go hub.run()
	go func() {
		http.HandleFunc("/blockdata", postBlockData)
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func postBlockData(w http.ResponseWriter, r *http.Request) {
	var blockData BlockRequestData
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "could not process data correctly")
	}

	json.Unmarshal(reqBody, &blockData)
	if blockData.Data == "" {
		fmt.Fprintf(w, "data field is required in order to make a block request.")
	}
	blockchain.AddNewBlockToBlockChain(string(blockData.Data))
	valid, index := blockchain.IsValidBlockChain()
	if !valid {
		panic(fmt.Sprintf("WARNING: Blockchain not valid at %+d", index))
	}

	go func(data string) {
		bSlice, _ := json.Marshal(blockchain)
		for _, c := range clientList {
			c.send <- bSlice
		}
	}(blockData.Data)
}

func (c *Client) broadcastAll() {
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
	client.send <- blockchain.BlockchainToBytes()
	clientList = append(clientList, *client)
	go client.broadcastAll()
}
