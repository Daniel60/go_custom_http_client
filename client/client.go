package client

import "net/http"

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type Client struct {
	httpClient HttpClient
}

func NewClient(hc HttpClient) *Client {
	return &Client{httpClient: hc}
}
