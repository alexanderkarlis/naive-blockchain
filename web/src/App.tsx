import React, { useState, useEffect } from "react";
import { TextField, Button, Box, Typography } from "@material-ui/core";
import Snackbar from '@material-ui/core/Snackbar';
import MuiAlert, { AlertProps } from '@material-ui/lab/Alert';

import { BlockComponent } from "./block";
import "./App.css";

let socket = new WebSocket("ws://localhost:8080/broadcast")

const Alert = (props: AlertProps) => {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

const ConnectionStatAlert = (props) => {
    let socketConn = props.socketConnection
    const [open, setOpen] = useState(true);

    const handleClose = (event?: React.SyntheticEvent, reason?: string) => {
      if (reason === 'clickaway') {
        return;
      }
  
      setOpen(false);
    };
    return (
      <Box>
        <Snackbar 
            open={open} 
            autoHideDuration={6000} 
            onClose={handleClose}
            anchorOrigin={{ vertical: 'top', horizontal: 'center' }}>
        {socketConn 
        ? 
          <Alert severity="success">
            Successfully connected to websocket!
          </Alert>
        :
          <Alert severity="error">
            There was a problem connecting to websocket!
          </Alert>
        }
        </Snackbar>
      </Box>
    )
}

const App = () => {
  const [text, setText] = useState("");
  const [blockchain, setBlockChain] = useState([]);
  const [okConn, setOkConn] = useState(false)

  let connect = () => {
    socket.onopen = () => {
      console.log("Successfully Connected"); 
      setOkConn(true)
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
      setOkConn(false)
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
      setBlockChain(JSON.parse(msg.data));
    };
  });

  return (
    <Box margin="10px">
      <header>
        <ConnectionStatAlert socketConnection={okConn} />
        <Box>
          <Typography variant="h4">Enter a new block</Typography>
        </Box>
        <div>
          <TextField
            style={{ paddingBottom: "10px" }}
            id="standard-basic"
            variant="outlined"
            label="Block Data"
            onChange={handleInput}
          />
        </div>
        <div>
          <Button variant="contained" onClick={sendMsg}>
            send
          </Button>
        </div>
        <div>
          <Box borderTop="solid 1px black" margin="10px">
            <Typography variant="h5">Blockchain ({blockchain.length})</Typography>
          </Box>
          {blockchain.map((block) => {
            return <BlockComponent block={block} />;
          })}
        </div>
      </header>
    </Box>
  );
}

export default App;
