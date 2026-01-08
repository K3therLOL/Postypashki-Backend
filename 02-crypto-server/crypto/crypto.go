package crypto

import (
	"encoding/json"
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

func (api *API) ListCryptos(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/coins/list", api.rootURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Header.Add("x-cg-demo-api-key", api.key)
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

	fmt.Fprintln(w, cryptos)
	fmt.Fprintln(w, len(cryptos))
}
