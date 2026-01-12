package crypto

import (
	"context"
	"cryptoserver/errorfmt"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

var (
	ErrListCrypto = errors.New("ListCrypto executed wrong.")
	ErrLimitExceeded = errors.New("Request limit exceeded.")
)

type API struct {
	rootURL string
	key string
	client *http.Client
	cache  *redis.Client
	recordsCount int
}

type CryptoDTO struct {
	Id string 	  `json:"id"`
	Symbol string `json:"symbol"`
	Name string   `json:"name"`
}

type CoinDTO struct {
	Symbol     string 		`json:"symbol"`
	Name       string 		`json:"name"`
	MarketData struct {
		CurrentPrice struct {
			Usd float64 `json:"usd"`
		}`json:"current_price"`
	} `json:"market_data"`
	LastUpdated string  `json:"last_updated"`
}

type CryptoDTOList struct {
	Coins []CryptoDTO `json:"coins"`
}

type OutputCoin struct {
	Symbol string      	 `json:"symbol"`
	Name string        	 `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdated string 	 `json:"last_updated"`
}

type HistoryDTO struct {
	Prices [][]float64 `json:"prices"`
}

type MemObject struct {
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}  

type OutputMem struct {
	Symbol string      	`json:"symbol"`
	History []MemObject `json:"history"`
}

type StatsDTO struct {
	CurrentPrice 			 float64 `json:"current_price"`
	High24h                  float64 `json:"high_24h"`
	Low24h 					 float64 `json:"low_24h"`
	PriceChange24h 			 float64 `json:"price_change_24h"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
}

type Record struct {
	MinPrice           float64 `json:"min_price"`
	MaxPrice           float64 `json:"max_price"`
	AvgPrice           float64 `json:"avg_price"`
	PriceChange        float64 `json:"price_change"`
	PriceChangePercent float64 `json:"price_change_percent"`
	RecordsCount       int     `json:"records_count"`
} 

type OutputStats struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
	Stats        Record  `json:"stats"`
}

func NewAPI() *API {
	api := &API{
		rootURL: "https://api.coingecko.com/api/v3",
		key: "",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		recordsCount: 100,
	}

	_, err := api.cache.Ping(context.Background()).Result()
    if err != nil {
        log.Println("Cannot connect to redis.", err)
        return nil
    }

    log.Println("Successful connection to redis.")
	return api
}

func (api *API) sendCryptoRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-cg-demo-api-key", api.key)
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (api *API) ListCryptos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	url := fmt.Sprintf("%s/coins/list", api.rootURL)
	resp, err := api.sendCryptoRequest(url)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, errorfmt.Jsonize(ErrLimitExceeded), http.StatusBadRequest)
		return
	}

	cryptos := make([]CryptoDTO, 1) 
	if err := json.NewDecoder(resp.Body).Decode(&cryptos); err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	api.cacheCryptoIDSet(cryptos)

	clientJSON, err := json.Marshal(cryptos)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(clientJSON))
}

func (api *API) GetCrypto(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	symbol := chi.URLParam(r, "symbol")
	id, err := api.getID(symbol)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusNotFound)
		return
	}

	url := fmt.Sprintf("%s/coins/%s", api.rootURL, id)
	resp, err := api.sendCryptoRequest(url)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	coin := CoinDTO{}
	if err := json.NewDecoder(resp.Body).Decode(&coin); err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}
	
	if coin.Symbol == "" && coin.Name == "" {
		http.Error(w, errorfmt.Jsonize(ErrLimitExceeded), http.StatusBadRequest)
		return
	}

	formatedCoin := OutputCoin{
		Symbol: coin.Symbol,
		Name: coin.Name,
		CurrentPrice: coin.MarketData.CurrentPrice.Usd,
		LastUpdated: coin.LastUpdated,
	}
	
	clientJSON, err := json.Marshal(formatedCoin)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(clientJSON))
}	

func (api *API) GetHistory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	symbol := chi.URLParam(r, "symbol")
	id, err := api.getID(symbol)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusNotFound)
		return
	}

	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=usd&days=1", api.rootURL, id)
	resp, err := api.sendCryptoRequest(url)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	history := HistoryDTO{}
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	formatedMem := OutputMem{
		Symbol: symbol,
		History: make([]MemObject, len(history.Prices)),
	}

	for i, price := range history.Prices {
		valPrice, ms := price[1], int64(price[0])
		formatedMem.History[i] = MemObject{
			Price: valPrice,
			Timestamp: time.UnixMilli(ms).UTC(),
		}
	}

	clientJSON, err := json.Marshal(formatedMem)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(clientJSON))
}

func (api *API) GetStats(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	symbol := chi.URLParam(r, "symbol")
	id, err := api.getID(symbol)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusNotFound)
		return
	}

	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s&symbols=%s", api.rootURL, id, symbol)
	resp, err := api.sendCryptoRequest(url)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, errorfmt.Jsonize(ErrLimitExceeded), http.StatusBadRequest)
		return
	}

	statsList := make([]StatsDTO, 1)
	if err := json.NewDecoder(resp.Body).Decode(&statsList); err != nil {
		fmt.Println("json err")
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}
	
	stats := statsList[0]
	formatedStats := OutputStats {
		Symbol: symbol,
		CurrentPrice: stats.CurrentPrice,
		Stats: Record {
			MinPrice: stats.Low24h,
			MaxPrice: stats.High24h,
			AvgPrice: (stats.Low24h + stats.High24h) / 2,
			PriceChange: stats.PriceChange24h,
			PriceChangePercent: stats.PriceChangePercentage24h,
			RecordsCount: api.recordsCount,
		},
	}

	clientJSON, err := json.Marshal(formatedStats)
	if err != nil {
		fmt.Println("json err")
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(clientJSON))
}
