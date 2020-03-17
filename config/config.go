package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Configuration file
const ConfigFileName = "config.json"

// Default values
const DefaultPort = 8000
const APIRootUrl = "wss://api.hitbtc.com/api/2/ws"

var DefaultTickers = []string{"BTCUSD", "ETHBTC"}

const GetAllTickerValue = "ALL"

type Config struct {
	PortNumber     int      `json:"port"`
	APIRootURL     string   `json:"apiRootUrl"`
	TrackedTickers []string `json:"trackedTickers"`
}

// Func getConfig return default
func GetConfig() Config {

	// Setup config with default parameters
	c := Config{
		PortNumber:     DefaultPort,
		APIRootURL:     APIRootUrl,
		TrackedTickers: DefaultTickers,
	}

	// Check for config file
	fileData, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		log.Printf("Failed to open config file: %s\n", ConfigFileName)
		log.Print("Using default values for config\n")
		return c
	}

	// Parse file contents by unmarshalling using JSON
	log.Printf("Parsing config file: %s\n", ConfigFileName)
	err = json.Unmarshal(fileData, &c)
	if err != nil {
		log.Printf("Failed to parse config file: %s\n", ConfigFileName)
		log.Print("Using default values for config\n")
		return c
	}

	return c
}
