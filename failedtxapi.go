package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	fbfailedtx "github.com/metachris/flashbots-failed-tx"
)

func GetFailedFlashbotsTransactions() (response []fbfailedtx.BlockWithFailedTx, err error) {
	url := fmt.Sprintf("%s/failedTx", Cfg.FlashbotsTxServerUrl)
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}
