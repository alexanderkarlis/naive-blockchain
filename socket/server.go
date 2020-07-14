package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/alexanderkarlis/naive-blockchain/block"
	"github.com/gorilla/websocket"
)

// Serve function serves a basic form of the listening websocket.
func Serve() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/broadcast", broadcast)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func broadcast(w http.ResponseWriter, r *http.Request) {
	var blockchain block.Blockchain
	genesis := block.CreateGenesisBlock()
	blockchain = append(blockchain, *genesis)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		blockchain.AddNewBlockToBlockChain(string(message))
		valid, index := blockchain.IsValidateBlockChain()
		fmt.Printf("Is valid blockchain?: %t, index: %v\n", valid, index)

		bSlice, _ := json.Marshal(blockchain)
		log.Printf("\n%+v", string(bSlice))

		err = c.WriteMessage(mt, bSlice)
		if err != nil {
			log.Println("err:", err)
			break
		}
	}
}
