package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/basicsbeauty/crypto-server/config"
	config2 "github.com/basicsbeauty/crypto-server/config"
	"github.com/basicsbeauty/crypto-server/price"
	"github.com/basicsbeauty/crypto-server/wsclient"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getIdFromPath(path string) (string, error) {
	const URLSeparator = "/"
	const MinURLWordCount = 1
	tokens := strings.Split(path, URLSeparator)
	if len(tokens) > MinURLWordCount {
		return tokens[len(tokens)-1], nil
	} else {
		return "", errors.New("invalid request URL")
	}
}

func handleGetCurrency(w http.ResponseWriter, r *http.Request) {

	id, err := getIdFromPath(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
	}

	// Convert ticker/id to uppercase
	id = strings.ToUpper(id)

	// Check if the ticker is valid
	log.Printf("handleGetCurrency: Check if ticker: %s is valid\n", id)
	if price.IsSymbolSupported(id) {
		pricing, err := price.GetPricingBySymbol(id)
		if err != nil {
			log.Printf("handleGetCurrency: Failed to get pricing for ticker: %s error: %s\n", id, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Failed to get pricing for ticker: %s error: %s"}`, id, err.Error())))
			return
		}

		payload, err := pricing.MarshalJSON()
		if err != nil {
			log.Printf("handleGetCurrency: Failed to get pricing for ticker: %s error: %s\n", id, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Failed to marshal pricing for ticker: %s error: %s"}`, id, err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	} else if id == config.GetAllTickerValue {
		prices := price.GetAllPricing()

		payload, err := json.Marshal(prices)
		if err != nil {
			log.Printf("handleGetCurrency: Failed to get pricing for ticker: %s error: %s\n", id, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message": "Failed to marshal pricing for ticker: %s error: %s"}`, id, err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Ticker: %s is not supported"}`, id)))
	}
}

func handleCurrency(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		handleGetCurrency(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"message": "Operation not supported"}`))
	}
}

func startServer(config config.Config) {

	const ApplicationName = "Toyota Test Crypto Server"

	http.HandleFunc("/currency/", handleCurrency)
	log.Printf("Starting service: %s, running on port: %d", ApplicationName, config.PortNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(config.PortNumber)), nil))
}

func main() {

	c := config2.GetConfig()
	log.Printf("Server configuration: %+v \n", c)

	go wsclient.StartWebSocketClient(c)
	startServer(c)
}
