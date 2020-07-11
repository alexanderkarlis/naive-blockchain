package server

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/alexanderkarlis/naive-blockchain/block"
	"github.com/gorilla/websocket"
)

// Serve function serves a basic form of the listening websocket.
func Serve() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	messages := []block.Block{}
	genesis := block.CreateGenesisBlock()
	messages = append(messages, *genesis)

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
		blk := block.New(string(message))

		messages = append(messages, *blk)
		bSlice, _ := json.Marshal(messages)
		log.Printf("\n%+v", string(bSlice))
		// log.Printf("recv: %+v", messages)
		err = c.WriteMessage(mt, bSlice)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
    <head>
    <meta charset="utf-8">
    <script>  
        window.addEventListener("load", function(evt) {
			var output = document.getElementById("output");
            var data = document.getElementById("data");
            var ws;
            var print = function(message) {
				output.textContent = ""
                var d = document.createElement("div");
                d.textContent = message;
                output.appendChild(d);
            };
            document.getElementById("open").onclick = function(evt) {
                if (ws) {
                    return false;
                }
                ws = new WebSocket("{{.}}");
                ws.onopen = function(evt) {
                    print("OPEN");
                }
                ws.onclose = function(evt) {
                    print("CLOSE");
                    ws = null;
                }
                ws.onmessage = function(evt) {
                    print(evt.data);
                }
                ws.onerror = function(evt) {
                    print("ERROR: " + evt.data);
                }
                return false;
            };
            document.getElementById("send").onclick = function(evt) {
                if (!ws) {
                    return false;
                }
                ws.send(data.value);
                return false;
            };
            document.getElementById("close").onclick = function(evt) {
                if (!ws) {
                    return false;
                }
                ws.close();
                return false;
            };
        });
    </script>
    </head>
<body>
    <table>
        <tr><td valign="top" width="50%">
        <p>Click "Open" to create a connection to the server, 
        "Send" to send a message to the server and "Close" to close the connection. 
        You can change the message and send multiple times.
        <p>
            <form>
                <button id="open">Open</button>
                <button id="close">Close</button>
				<p>
				<label for="data">Data:</label>
				<input id="data" type="text" value="hello">
                <button id="send">Send</button>
            </form>
        </td><td valign="top" width="50%">
        <div id="output"></div>
        </td></tr>
    </table>
    <div>
        Blockchain
    </div>
    </body>
</html>
`))
