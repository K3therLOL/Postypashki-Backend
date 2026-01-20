package crypto

import (
	"strings"
	"errors"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var ctx = context.Background()
var ErrNoID = errors.New("No id by your symbol.")

func (api *API) cacheCryptoID(crypto CryptoDTO) {
	err := api.cache.SetNX(ctx, crypto.Symbol, crypto.Id, 30 * time.Minute).Err()
	if err != nil {
		log.Println(err.Error())
	}
}

func (api *API) cacheCryptoIDSet(cryptos []CryptoDTO) {
	for _, crypto := range cryptos {
		api.cacheCryptoID(crypto)
	}
}

func (api *API) getID(symbol string) (string, error) {
	id, err := api.cache.Get(ctx, symbol).Result()
	if err == nil {
		fmt.Println("cache boom")
		return id, nil
	}

	// CACHE MISS
	url := fmt.Sprintf("%s/search?query=%s", api.rootURL, symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("x-cg-demo-api-key", api.key)

	resp, err := api.client.Do(req)
	if err != nil {
		return "", err
	}

	cryptos := CryptoDTOList{}
	if err := json.NewDecoder(resp.Body).Decode(&cryptos); err != nil {
		return "", err
	}

	for _, crypto := range cryptos.Coins {
		if strings.EqualFold(crypto.Symbol, symbol) {
			id := crypto.Id	
			return id, nil
		}
	}

	return "", ErrNoID
}
