package mailtrap

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

// setup sets up a test HTTP server.
func setup() (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client = New("api-token")

	url, _ := url.Parse(server.URL)
	client.defaultBaseURL = url
	client.sendEmailBaseURL = url

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

func testURL(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.RequestURI; got != want {
		t.Errorf("Request url: %+v, want %s", got, want)
	}
}

func testParams(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.URL.RawQuery; got != want {
		t.Errorf("Request query: %s, want %s", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		t.Fatalf("Failed to Read Body: %v", err)
	}

	if got := buffer.String(); got != want {
		t.Errorf("Request body: %s, want %s", got, want)
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

func TestNewClient(t *testing.T) {
	apiKey := "api-token"
	expectedBaseURL := defaultBaseURL
	expectedSendURL := sendEmailBaseURL

	c := New(apiKey)

	if c.apiKey != apiKey {
		t.Errorf("Client apiKey is %s, want %s", c.apiKey, apiKey)
	}
	if c.defaultBaseURL.String() != expectedBaseURL {
		t.Errorf("Client defaultBaseURL is %s, want %s", c.defaultBaseURL.String(), expectedBaseURL)
	}
	if c.sendEmailBaseURL.String() != expectedSendURL {
		t.Errorf("Client sendEmailBaseURL is %s, want %s", c.sendEmailBaseURL.String(), expectedSendURL)
	}
	if c.userAgent != userAgent {
		t.Errorf("Client userAgent is %s, want %s", c.userAgent, userAgent)
	}
}

func TestNewRequest(t *testing.T) {

}

func TestCheckResponse(t *testing.T) {

}

func TestCheckResponseOnUnknownErrorFormat(t *testing.T) {

}
