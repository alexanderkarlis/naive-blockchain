package main

import (
	"github.com/alexanderkarlis/naive-blockchain/block"
	server "github.com/alexanderkarlis/naive-blockchain/socket"
)

var blockchain chan []block.Block

func main() {
	server.SocketServe()
}
