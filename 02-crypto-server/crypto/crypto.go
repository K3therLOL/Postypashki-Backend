package crypto

import (
	"encoding/json"
	"sync"
	"fmt"
	"net/http"
	"time"
)

type API struct {
	rootURL string
	key string
	client *http.Client
}

func NewAPI() *API {
	return &API{
		rootURL: "https://api.coingecko.com/api/v3",
		key: "",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type CryptoDTO struct {
	Id string 	  `json:"id"`
	Symbol string `json:"symbol"`
	Name string   `json:"name"`
}

type Coin struct {
	Symbol string 		`json:"symbol"`
	Name string 		`json:"name"`
	CurrentPrice string `json:"current_price"`
	LastUpdated string  `json:"last_updated"`
}

type OutputCrypto struct {
	Cryptos []Coin `json:"cryptos"`
}

func (api *API) ListCryptos(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/coins/list", api.rootURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//req.Header.Add("x-cg-demo-api-key", api.key)
	resp, err := api.client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "ListCrypto executed wrong.", http.StatusBadRequest)
		return
	}

	cryptos := make([]CryptoDTO, 1)
	if err := json.NewDecoder(resp.Body).Decode(&cryptos); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ch := make(chan Coin)
	wg := sync.WaitGroup{}
	for _, crypto := range cryptos {
		wg.Go(func() {
			url := fmt.Sprintf("%s/coins/%s", api.rootURL, crypto.Id)
			resp, err := http.Get(url)
			if err != nil {
				return
			}

			coin := Coin{}
			if err := json.NewDecoder(resp.Body).Decode(&coin); err != nil {
				return
			}

			ch <- coin
		})
	}

	wg.Wait()
	go func() { 
		close(ch) 
	}()

	coins := OutputCrypto{}
	coins.Cryptos = make([]Coin, 1)
	for coin := range ch {
		coins.Cryptos = append(coins.Cryptos, coin)
	}

	outputJSON, err := json.Marshal(coins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, cryptos)
	fmt.Fprintln(w, len(cryptos))
	fmt.Fprintln(w, string(outputJSON))
}
