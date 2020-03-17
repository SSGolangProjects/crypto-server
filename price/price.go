package price

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/basicsbeauty/crypto-server/config"
	"log"
	"sync"
)

type Pricing struct {
	Id          string `json:"id"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	Symbol      string `json:"symbol"`
	FeeCurrency string `json:"feeCurrency"`
	FullName    string `json:"fullName"`
}

var temp Pricing
var readMutex *sync.Mutex
var updateMutex *sync.Mutex
var priceMap map[string]Pricing

// Pricing, initialize mutex and pricing map
func init() {
	log.Printf("Pricing: Initalizing mutex\n")
	readMutex = &sync.Mutex{}
	updateMutex = &sync.Mutex{}

	// Initialize ticker:price map
	priceMap = make(map[string]Pricing)
	for _, ticker := range config.DefaultTickers {
		switch ticker {
		case "BTCUSD":
			priceMap[ticker] = Pricing{Id: "BTC", FullName: "Bitcoin", FeeCurrency: "USD"}
		case "ETHBTC":
			priceMap[ticker] = Pricing{Id: "ETH", FullName: "Ethereum", FeeCurrency: "BTC"}
		default:
			log.Printf("Pricing: Init: Unknow ticker: %s, continuing", ticker)
			continue
		}
	}
}

// updatePriceData - Update pricing data of a ticker
func (p Pricing) updatePriceData(src Pricing) Pricing {
	p.Ask = src.Ask
	p.Bid = src.Bid
	p.Last = src.Last
	p.Open = src.Open
	p.Low = src.Low
	p.High = src.High
	return p
}

func (p *Pricing) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Id          string `json:"id"`
			Ask         string `json:"ask"`
			Bid         string `json:"bid"`
			Last        string `json:"last"`
			Open        string `json:"open"`
			Low         string `json:"low"`
			High        string `json:"high"`
			FeeCurrency string `json:"feeCurrency"`
			FullName    string `json:"fullName"`
		}{
			Id:          p.Id,
			Ask:         p.Ask,
			Bid:         p.Bid,
			Last:        p.Last,
			Open:        p.Open,
			Low:         p.Low,
			High:        p.High,
			FeeCurrency: p.FeeCurrency,
			FullName:    p.FullName,
		})
}

// IsSymbolSupported - Check if a symbol is supported
func IsSymbolSupported(symbol string) bool {
	if _, ok := priceMap[symbol]; ok {
		return true
	} else {
		return false
	}
}

// ProcessUpdate - Update pricing for a symbol
func ProcessUpdate(p Pricing) {

	if !IsSymbolSupported(p.Symbol) {
		return
	}

	updateMutex.Lock()
	priceMap[p.Symbol] = priceMap[p.Symbol].updatePriceData(p)
	//temp := &priceMap[p.Symbol].updatePriceData(p)
	//temp.updatePriceData(p)
	//priceMap[p.Symbol] = temp
	updateMutex.Unlock()
}

// GetPricingBySymbol - Get pricing for a symbol
func GetPricingBySymbol(symbol string) (Pricing, error) {
	readMutex.Lock()
	defer readMutex.Unlock()
	if details, ok := priceMap[symbol]; ok {
		return details, nil
	} else {
		return Pricing{}, errors.New(fmt.Sprintf("Pricing: Not found for symbol: %s", symbol))
	}
}

// GetAllPricing - Get pricing for all symbols
func GetAllPricing() []Pricing {
	var rvalue []Pricing
	for _, detail := range priceMap {
		rvalue = append(rvalue, detail)
	}
	return rvalue
}
