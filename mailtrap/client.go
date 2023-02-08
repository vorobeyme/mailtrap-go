package mailtrap

import (
	"net/http"
	"net/url"
)

const baseURL = "https://api.mailtrap.io"

type Client struct {
	apiKey         string
	defaultBaseURL *url.URL
	httpClient     *http.Client
}

func New(apiKey string) (*Client, error) {
	defaultURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		apiKey:         apiKey,
		defaultBaseURL: defaultURL,
		httpClient:     http.DefaultClient,
	}, nil
}

type Response struct {
	*http.Response
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	return nil, nil
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	return nil, nil
}
