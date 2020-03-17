# Crypto-server

A simple Crypto server that supports below URI's
```bash
curl http://localhost:8080/currency/btcusd
curl http://localhost:8080/currency/ethbtc
curl http://localhost:8080/currency/all

```


## Build and run

Run below commands

```bash
go get https://github.com/basicsbeauty/crypto-server
cd crypto-server
go run main.go
```

## Checklist of items

- [x] Server configuration is driven by a JSON file, `config.json`
- [x] Support for parameters {port number, source API URL, supported tickers}
- [x] Websocket implementation to keep the data in sync
- [x] Used local cache to store and serve requests
- [x] Used native HTTP library for handling REST API calls
- [x] Handled corner cases, the server only responds to GET calls for /currency/ 
- [x] Developed separate modules for {config, pricing data, websocket communication} for modularity
- [ ] Unit tests


## Libraries
Used gorilla/websocket library for handling web socket network IO for better reliability and its battle-tested.

