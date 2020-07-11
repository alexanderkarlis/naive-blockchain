import React, { useState, useEffect } from "react";
import { TextField, Button } from "@material-ui/core";
import { BlockComponent } from "./block";
import "./App.css";

const socket = new WebSocket("ws://localhost:8080/echo");

function App() {
    const [text, setText] = useState("");
    const [blockchain, setBlockChain] = useState([]);
    
    
    
    let connect = () => {
        socket.onopen = () => {
            console.log("Successfully Connected");
        };
        
        socket.onmessage = (msg) => {
            console.log("msg");
            console.log(JSON.parse(msg.data));
            setBlockChain(JSON.parse(msg.data));
        };
        
        socket.onclose = (event) => {
            console.log("Socket Closed Connection: ", event);
        };
        
        socket.onerror = (error) => {
            console.log("Socket Error: ", error);
        };
    };
    
    let sendMsg = () => {
        console.log("sending msg: ", text);
        socket.send(text);
    };
    
    const handleInput = (e) => {
        setText(e.target.value);
    };
    
    useEffect(() => {
        connect();
    }, []);
    
    useEffect(() => {
        socket.onmessage = (msg) => {
            console.log("msg");
            setBlockChain(JSON.parse(msg.data))
        };

    })
    
    return (
        <div>
            <header>
                <div>
                <TextField
                    style={{ paddingBottom: "10px" }}
                    id="standard-basic"
                    variant="outlined"
                    label="Standard"
                    onChange={handleInput}
                />
                </div>
                <div>
                    <Button variant="contained" onClick={sendMsg}>
                        send msg
                    </Button>
                </div>
                <div>
                    {blockchain.map((block) => {
                        return <BlockComponent block={block} />;
                    })}
                </div>
            </header>
        </div>
        );
    }
    
    export default App;
    