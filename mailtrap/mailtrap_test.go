package mailtrap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

// setupTestingClient sets up a test HTTP server for testing API client.
func setupTestingClient() (client *TestingClient, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client, _ = NewTestingClient("api-token")
	url, _ := url.Parse(server.URL)
	client.baseURL = url

	return client, mux, server.Close
}

// setupSendingClient sets up a test HTTP server for sending API client.
func setupSendingClient() (client *SendingClient, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client, _ = NewSendingClient("api-token")
	url, _ := url.Parse(server.URL)
	client.baseURL = url

	return client, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testBadPathParams(t *testing.T, method string, fn func() error) {
	t.Helper()
	if err := fn(); err == nil {
		t.Errorf("%v bad params, err = nil, want error", method)
	}
}

func testNewRequestAndDoFail(t *testing.T, method string, c *client, fn func() (*Response, error)) {
	t.Helper()
	c.baseURL.Host = "!@#$%^&*()_+"
	resp, err := fn()
	if resp != nil {
		t.Errorf("%v client.BaseURL=Host='%v', resp = %#v, want nil", method, c.baseURL.Host, resp)
	}
	if err == nil {
		t.Errorf("%v client.BaseURL=Host='%v', err = nil, want error", method, c.baseURL.Host)
	}
}

// testJSONMarshal tests whether the marshaling produces a JSON
// that corresponds to the want string.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	t.Helper()

	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	if string(w) != string(j) {
		t.Errorf("json.Marshal(%q) \nreturned %s,\nwant %s", v, j, w)
	}
}

func TestNewSendingClient(t *testing.T) {
	apiKey := "api-token"
	expectedBaseURL := sendingAPIURL + apiSuffix

	c, err := NewSendingClient(apiKey)
	if err != nil {
		t.Errorf("Sending client returned error: %v", err)
	}

	if c.apiKey != apiKey {
		t.Errorf("Sending client apiKey is %s, want %s", c.apiKey, apiKey)
	}
	if c.baseURL.String() != expectedBaseURL {
		t.Errorf("Sending client baseURL is %s, want %s", c.baseURL.String(), expectedBaseURL)
	}
	if c.userAgent != userAgent {
		t.Errorf("Sending client userAgent is %s, want %s", c.userAgent, userAgent)
	}
}

func TestNewTestingClient(t *testing.T) {
	apiKey := "api-token"
	expectedBaseURL := testingAPIURL + apiSuffix

	c, err := NewTestingClient(apiKey)
	if err != nil {
		t.Errorf("Testing client returned error: %v", err)
	}

	if c.apiKey != apiKey {
		t.Errorf("Testing client apiKey is %s, want %s", c.apiKey, apiKey)
	}
	if c.baseURL.String() != expectedBaseURL {
		t.Errorf("Testing client baseURL is %s, want %s", c.baseURL.String(), expectedBaseURL)
	}
	if c.userAgent != userAgent {
		t.Errorf("Testing client userAgent is %s, want %s", c.userAgent, userAgent)
	}
}

func TestNewRequest(t *testing.T) {

}

func TestCheckResponse(t *testing.T) {

}

func TestCheckResponseOnUnknownErrorFormat(t *testing.T) {

}
