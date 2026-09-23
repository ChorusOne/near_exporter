package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (c *nearExporter) getAccountBalance(account string) (float64, error) {
	accountJSON, _ := json.Marshal(account)
	req, err := http.NewRequest("POST", c.rpcAddr, bytes.NewBufferString(fmt.Sprintf(
		`{"jsonrpc":"2.0","id":1,"method":"query","params":{"request_type":"view_account","finality":"final","account_id":%s}}`, accountJSON)))
	if err != nil {
		return 0, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("account query status code %d", resp.StatusCode)
	}

	var response struct {
		Result struct {
			Amount string `json:"amount"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, fmt.Errorf("JSONRPC error: %s", response.Error.Message)
	}
	balance, err := strconv.ParseFloat(response.Result.Amount, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid account balance: %w", err)
	}
	// RPC returns balance in yoctoNEAR
	return balance / 1e24, nil
}
