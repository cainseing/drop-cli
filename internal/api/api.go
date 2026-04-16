package api

import (
	"fmt"

	"github.com/cainseing/drop-cli/internal/config"
	"github.com/go-resty/resty/v2"
)

type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

type DropRequest struct {
	Blob      string `json:"blob"`
	TTL       int    `json:"ttl"`
	Reads     int    `json:"reads"`
	Signature string `json:"signature,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

type GetDropResponse struct {
	Blob           string `json:"blob"`
	RemainingReads int    `json:"remaining_reads"`
	Signature      string `json:"signature,omitempty"`
	Sender         string `json:"sender,omitempty"`
	Provider       string `json:"provider,omitempty"`
}

func newClient() *resty.Client {
	client := resty.New()
	client.SetBaseURL(config.ApiURL)
	client.SetHeader("X-Drop-Client", "drop-cli-v1")
	client.SetHeader("User-Agent", "DropCLI/v1.0")
	return client
}

func postBlob(blob string, ttl int, reads int, sig string, sender string, provider string) (string, error) {
	client := newClient()

	body := DropRequest{
		Blob:      blob,
		TTL:       ttl * 60,
		Reads:     reads,
		Signature: sig,
		Sender:    sender,
		Provider:  provider,
	}

	var result struct {
		Id string `json:"id"`
	}

	errorResponse := ErrorResponse{}

	resp, err := client.R().
		SetBody(body).
		SetResult(&result).
		SetError(&errorResponse).
		Post("/blob")

	if err != nil {
		return "", fmt.Errorf("request to API failed, please try again")
	}

	if resp.IsError() {
		return "", fmt.Errorf("request to API failed: %s", errorResponse.Message)
	}

	return result.Id, nil
}

func getBlob(id string) (*GetDropResponse, error) {
	client := newClient()

	result := GetDropResponse{}

	resp, err := client.R().
		SetResult(&result).
		Get("/blob/" + id)

	print(result.Provider)

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("drop was not found")
	}

	if err != nil || resp.IsError() {
		return nil, fmt.Errorf("request to API failed, please try again")
	}

	return &result, nil
}

func purgeBlob(id string) (bool, error) {
	client := newClient()

	resp, err := client.R().
		Delete("/blob/" + id)

	if resp.StatusCode() == 404 {
		return false, fmt.Errorf("drop was not found")
	}

	if err != nil || resp.IsError() {
		return false, fmt.Errorf("request to API failed, please try again")
	}

	return true, nil
}
