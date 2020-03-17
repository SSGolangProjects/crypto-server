package wsclient

import (
	"encoding/json"
	"fmt"
	"github.com/basicsbeauty/crypto-server/config"
	"github.com/basicsbeauty/crypto-server/price"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"time"
)

type WebSocketResponse struct {
	Jsonrpc string
	Method  string
	Price   price.Pricing `json:"params"`
}

// StartWebSocketClient - WebSocket setup
func StartWebSocketClient(config config.Config) {
	log.Printf("Connecting to server: %s\n", config.APIRootURL)

	interrupt := make(chan os.Signal, 1)

	// Connection
	conn, _, err := websocket.DefaultDialer.Dial(config.APIRootURL, nil)
	if err != nil {
		log.Fatal("WebSocket: Client: Server connection FAILED", err)
		panic(err)
	}
	defer conn.Close()

	// Subscribe tickers
	for _, ticker := range config.TrackedTickers {

		// Build the message
		message := fmt.Sprintf("{\"method\":\"subscribeTicker\",\"params\":{\"symbol\":\"%s\"},\"id\":786}", ticker)
		log.Printf("WebSocket: Client: Ticker: %s Subscribing with payload: %s\n", ticker, message)

		//jsonPayload, _ := json.Marshal(message)
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Fatalf("WebSocket: Client: Ticker: %s Subscribing FAILED error: %s\n", ticker, err.Error())
		}
		log.Printf("WebSocket: Client: Ticker: %s Subscribing Success\n", ticker)
	}

	// Read ticker updates
	done := make(chan struct{})
	go func() {
		defer close(done)
		var tempResponse WebSocketResponse
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket: Client: Failed to read meassage: Error: %s\n", err.Error())
			}
			//log.Printf("WebSocket: Client: Notification Meassage: %s\n", string(message))

			tempResponse = WebSocketResponse{}
			//log.Printf("WebSocket: Client: Read meassage: Before: %+v\n", tempResponse)
			err = json.Unmarshal(message, &tempResponse)
			if err != nil {
				log.Printf("WebSocket: Client: Read meassage: Unmarshal error: %s, continuing\n", err.Error())
				continue
			}
			price.ProcessUpdate(tempResponse.Price)
		}
	}()

	for {
		select {
		case <-done:
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("WebSocket: Client: Failed to close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
