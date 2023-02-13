package mailtrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
)

const (
	libVersion = "0.1.0"

	defaultBaseURL   = "https://mailtrap.io/api/"
	sendEmailBaseURL = "https://send.api.mailtrap.io/api/"
	contentType      = "application/json"
)

var (
	userAgent = fmt.Sprintf("mailtrap-go/%s (%s %s) go/%s", libVersion, runtime.GOOS, runtime.GOARCH, runtime.Version())
)

// Client manages communication with the Mailtrap API.
type Client struct {
	// API key used to make authenticated API calls.
	apiKey string

	// Base URL for API requests.
	defaultBaseURL *url.URL

	// Base URL for email sending API requests.
	sendEmailBaseURL *url.URL

	// User agent used when communicating with the API.
	userAgent string

	// HTTP client used to communicate with the API.
	httpClient *http.Client

	// Service for sending emails.
	SendEmail *SendEmailService

	// Services used for communicating with the Mailtrap API.
	Accounts     *AccountsService
	AccountUsers *AccountUsersService
	Permissions  *PermissionsService
	Projects     *ProjectsService
	Inboxes      *InboxesService
	Messages     *MessagesService
	Attachments  *AttachmentsService
}

// New returns a new Mailtrap API client instance.
func New(apiKey string) *Client {
	defaultURL, _ := url.Parse(defaultBaseURL)
	sendEmailURL, _ := url.Parse(sendEmailBaseURL)

	client := &Client{
		apiKey:           apiKey,
		defaultBaseURL:   defaultURL,
		sendEmailBaseURL: sendEmailURL,
		httpClient:       http.DefaultClient,
		userAgent:        userAgent,
	}

	// Create all the public services.
	client.SendEmail = &SendEmailService{client: client}
	client.Accounts = &AccountsService{client: client}
	client.AccountUsers = &AccountUsersService{client: client}
	client.Permissions = &PermissionsService{client: client}
	client.Projects = &ProjectsService{client: client}
	client.Inboxes = &InboxesService{client: client}
	client.Messages = &MessagesService{client: client}
	client.Attachments = &AttachmentsService{client: client}

	return client
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := &Response{Response: resp}
	if err := CheckResponse(resp); err != nil {
		return response, err
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, err
		}
	}

	return response, err
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.retrieveApiURL(path)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return req, nil
}

func (c *Client) retrieveApiURL(path string) (*url.URL, error) {
	u := c.defaultBaseURL
	if path == sendEmailEndpoint {
		u = c.sendEmailBaseURL
	}

	return u.Parse(path)
}

// Response is a Mailtrap response.
// This wraps the standard http.Response returned from Mailtrap.
type Response struct {
	*http.Response
}

// CheckResponse checks the API response for errors and returns them if present.
// A response is considered an error if it has a status code outside the 200-299 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	errResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errResponse)
		if err != nil {
			errResponse.Message = string(data)
		}
	}

	return errResponse
}
