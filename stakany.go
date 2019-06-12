package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

type ReadMsg struct {
	LastUpdateID int         `json:"lastUpdateId"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`

	Bests Bests `json:"-"`
}

type Bests struct {
	Bid struct {
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
	} `json:"bid"`

	Ask struct {
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
	} `json:"ask"`
}

func (msg *ReadMsg) FindBests() {

	for i, bid := range msg.Bids {
		val, _ := strconv.ParseFloat(bid[0], 64)
		if i == 0 || msg.Bests.Bid.Price < val {
			msg.Bests.Bid.Price = val
			msg.Bests.Bid.Amount, _ = strconv.ParseFloat(bid[1], 64)
		}
	}

	for i, ask := range msg.Asks {
		val, _ := strconv.ParseFloat(ask[0], 64)
		if i == 0 || msg.Bests.Ask.Price > val {
			msg.Bests.Ask.Price = val
			msg.Bests.Ask.Amount, _ = strconv.ParseFloat(ask[1], 64)
		}
	}
}

func main() {
	u := url.URL{
		Scheme: "wss",
		Host:   "stream.binance.com:9443",
		Path:   "ws/btcusdt@depth10",
	}

	log.Println("Connecting to " + u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic("Cannot connect to binance: " + err.Error())
	}
	log.Println("Connected to " + u.String())
	defer c.Close()

	for {
		msg := ReadMsg{}

		err := c.ReadJSON(&msg)
		if err != nil {
			panic("Message error: " + err.Error())
		}
		msg.FindBests()
		bests, _ := json.Marshal(msg.Bests)
		fmt.Println(string(bests))
	}
}
