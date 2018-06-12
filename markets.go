package main

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"sort"
)

const MARKET_API = "https://api.coinmarketcap.com/v2/"

type Market struct{}

type MarketListingResponse struct {
	Data []struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Slug   string `json:"website_slug"`
	}
}

type MarketTickerResponse struct {
	Data map[string]MarketTickerObject
}

type MarketTickerObject struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	Rank              int     `json:"rank"`
	CirculatingSupply float64 `json:"circulating_supply"`
	TotalSupply       float64 `json:"total_supply"`
	MaxSupply         float64 `json:"max_supply"`
	Quotes            map[string]MarketTickerObjectQuotes
}

type MarketTickerObjectQuotes struct {
	Price            float64 `json:"price"`
	Volume24h        float64 `json:"volume_24h"`
	MarketCap        float64 `json:"market_cap"`
	PercentChange1h  float64 `json:"percent_change_1h"`
	PercentChange24h float64 `json:"percent_change_24h"`
	PercentChange7d  float64 `json:"percent_change_7d"`
}

type MarketGlobalResponse struct {
	Data struct {
		ActiveCryptocurrencies       int     `json:"active_cryptocurrencies"`
		ActiveMarkets                int     `json:"active_markets"`
		BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap"`
		Quotes struct {
			USD struct {
				TotalMarketCap float64 `json:"total_market_cap"`
				TotalVolume24H float64 `json:"total_volume_24h"`
			} `json:"USD"`
		} `json:"quotes"`
		LastUpdated int `json:"last_updated"`
	} `json:"data"`
}

func (m Market) call(url string, res interface{}) error {
	var urlBuffer bytes.Buffer
	urlBuffer.WriteString(MARKET_API)
	urlBuffer.WriteString(url)
	url = urlBuffer.String()

	r, err := http.Get(url)

	if err == nil {
		defer r.Body.Close()
	}

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	return nil
}

func (m Market) listing() MarketListingResponse {
	res := MarketListingResponse{}
	err := m.call("listings", &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func (m Market) ticker() map[string]MarketTickerObject {
	res := MarketTickerResponse{}
	err := m.call("ticker", &res)
	if err != nil {
		log.Fatal(err)
	}

	datas := make([]string, 0, len(res.Data))
	for data := range res.Data {
		datas = append(datas, data)
	}
	sort.Strings(datas)
	r := make(map[string]MarketTickerObject)
	for _, data := range datas {
		r[data] = res.Data[data]
	}

	return r
}

func (m Market) global() MarketGlobalResponse {
	res := MarketGlobalResponse{}
	err := m.call("global", &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func FloatToInt(f float64) int64 {
	c, _ := strconv.ParseInt(fmt.Sprintf("%.f", f), 0, 0)
	return c
}