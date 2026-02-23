package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

type DropRequest struct {
	Blob  string `json:"blob"`
	TTL   int    `json:"ttl"`
	Reads int    `json:"reads"`
}

type GetDropResponse struct {
	Blob           string `json:"blob"`
	RemainingReads int    `json:"remaining_reads"`
}

func newClient() *resty.Client {
	client := resty.New()
	client.SetHeader("X-Drop-Client", "drop-cli-v1")
	client.SetHeader("User-Agent", "DropCLI/v1.0")
	return client
}

func postBlob(blob string, ttl int, reads int) (string, error) {
	client := newClient()

	body := DropRequest{
		Blob:  blob,
		TTL:   ttl * 60,
		Reads: reads,
	}

	var result struct {
		Id string `json:"id"`
	}

	errorResponse := ErrorResponse{}

	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		SetError(&errorResponse).
		Post(viper.GetString("api_url") + "/blob")

	if err != nil {
		return "", fmt.Errorf("Request to API failed, please try again")
	}

	if resp.IsError() {
		return "", fmt.Errorf("Request to API failed: %s", errorResponse.Message)
	}

	return result.Id, nil
}

func getBlob(id string) (*GetDropResponse, error) {
	client := newClient()

	result := GetDropResponse{}

	resp, err := client.R().
		SetResult(&result).
		Get(viper.GetString("api_url") + "/blob/" + id)

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("This drop was not found")
	}

	if err != nil || resp.IsError() {
		return nil, fmt.Errorf("Request to API failed, please try again.")
	}

	return &result, nil
}

func purgeBlob(id string) (bool, error) {
	client := newClient()

	resp, err := client.R().
		Delete(viper.GetString("api_url") + "/blob/" + id)

	if resp.StatusCode() == 404 {
		return false, fmt.Errorf("This drop was not found")
	}

	if err != nil || resp.IsError() {
		return false, fmt.Errorf("Request to API failed, please try again")
	}

	return true, nil
}
