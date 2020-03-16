package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
)

// Configuration
const ConfigFileName = "config.json"

// Default values
const DefaultPort = 8000
const APIRootUrl = "https://api.hitbtc.com/api/2/public/currency"

type Config struct {
	PortNumber     int      `json:"port"`
	APIRootURL     string   `json:"apiRootUrl"`
	TrackedTickers []string `json:"trackedTickers"`
}

// Func getConfig return default
func getConfig() Config {

	// Setup config with default parameters
	c := Config{
		PortNumber: DefaultPort,
		APIRootURL: APIRootUrl}

	// Check for config file
	fileData, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		log.Printf("Failed to open config file: %s", ConfigFileName)
		log.Print("Using default values for config")
		return c
	}

	// Parse file contents by unmarshalling using JSON
	log.Printf("Parsing config file: %s", ConfigFileName)
	err = json.Unmarshal(fileData, &c)
	if err != nil {
		log.Printf("Failed to parse config file: %s", ConfigFileName)
		log.Print("Using default values for config")
		return c
	}

	return c
}

// WebSocket setup
func setWebSocketClient(config Config) {
	log.Printf("Connecting to server: %s\n", config.APIRootURL)

	// Connection
	c, _, err := websocket.DefaultDialer.Dial(config.APIRootURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Subscribe tickers
	for _, ticker := range config.TrackedTickers {

		// Build the message
		message := fmt.Sprintf("{\"method\":\"subscribeTicker\",\"params\":{\"symbol\":\"%s\"},\"id\":786}", ticker)
		jsonPayload, _ := json.Marshal(message)
		log.Printf("WebSocket: Client: Ticker: %s Subscribing with payload: %s\n", ticker, message)
		log.Printf("WebSocket: Client: Ticker: %s Subscribing with payload: %s\n", ticker, string(jsonPayload))

		//jsonPayload, _ := json.Marshal(message)
		err := c.WriteMessage(websocket.TextMessage, jsonPayload)
		if err != nil {
			log.Fatalf("WebSocket: Client: Ticker: %s Subscribing FAILED error: %s\n", ticker, err.Error())
		}
		log.Printf("WebSocket: Client: Ticker: %s Subscribing Success\n", ticker)
	}

	// Read notifications
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("WebSocket: Client: Failed to read meassage: Error: %s\n", err.Error())
			}
			log.Printf("WebSocket: Client: Read meassage: Error: %s\n", string(message))
		}
	}()

	for {
		select {
		case <-done:
			return
		}
	}
}

func main() {

	c := getConfig()

	bc, _ := json.Marshal(c)
	log.Printf("Server configuration: %s \n", string(bc))

	setWebSocketClient(c)
}
