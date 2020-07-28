# Naive Blockchain
#### Simple blockchain implementation with Go running the server and React running the socket connection to 'peers'. Based on [this](https://lhartikk.github.io/) 

### Run
#### Start the server
```sh
go run main.go
```
#### Then `start 'peer to peer' socket connection
```sh
npm install --prefix web
npm run start --prefix web
```

#### Add a new block to the blockchain!
```sh
curl -X POST http://localhost:8080/blockdata -d "{\"data\": \"New block post!\"}"
```

#### Proof of work done by hard coded difficulty to 2.
** leading hash must be lead by `difficulty * '0'` (padded number of times) ** 


## Todo's!
**lot's of TODO's left.**
