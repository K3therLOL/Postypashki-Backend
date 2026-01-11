package crypto

import (
	"cryptoserver/errorfmt"
	"errors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/go-chi/chi/v5"

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
}

type CryptoDTO struct {
	Id string 	  `json:"id"`
	Symbol string `json:"symbol"`
	Name string   `json:"name"`
}

type CoinDTO struct {
	Symbol string 		`json:"symbol"`
	Name string 		`json:"name"`
	MarketData struct {
		CurrentPrice struct {
			Usd float64 `json:"usd"`
		}`json:"current_price"`
	} `json:"market_data"`
	LastUpdated string  `json:"last_updated"`
}

type OutputCoin struct {
	Symbol string      	 `json:"symbol"`
	Name string        	 `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdated string 	 `json:"last_updated"`
}

type CryptoDTOList struct {
	Coins []CryptoDTO `json:"coins"`
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
	}

	_, err := api.cache.Ping(context.Background()).Result()
    if err != nil {
        log.Println("Cannot connect to redis.", err)
        return nil
    }

    log.Println("Successful connection to redis.")
	return api
}

func (api *API) ListCryptos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	url := fmt.Sprintf("%s/coins/list", api.rootURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	req.Header.Add("x-cg-demo-api-key", api.key)
	resp, err := api.client.Do(req)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, errorfmt.Jsonize(ErrListCrypto), http.StatusBadRequest)
		return
	}

	cryptos := make([]CryptoDTO, 1) 
	if err := json.NewDecoder(resp.Body).Decode(&cryptos); err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	api.cacheCryptoIDSet(cryptos)

	outputJSON, err := json.Marshal(cryptos)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(outputJSON))
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	req.Header.Add("x-cg-demo-api-key", api.key)
	resp, err := api.client.Do(req)
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
	
	outputJSON, err := json.Marshal(formatedCoin)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(outputJSON))
}	
