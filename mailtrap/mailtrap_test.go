package mailtrap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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
	c, _ := NewTestingClient("")

	inURL, outURL := "/accounts/1/projects", testingAPIURL+"api/accounts/1/projects"
	inBody := &PermissionRequest{
		ResourceID:   1,
		ResourceType: "account",
		AccessLevel:  "100",
	}
	outBody := `{"resource_id":1,"resource_type":"account","access_level":"100"}`

	req, _ := c.NewRequest(http.MethodPost, inURL, inBody)

	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	body, _ := ioutil.ReadAll(req.Body)
	if strings.TrimSpace(string(body)) != outBody {
		t.Errorf("NewRequest(%v) Body = %v, expected %v", inBody, strings.TrimSpace(string(body)), outBody)
	}

	userAgent := req.Header.Get("User-Agent")
	if c.userAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.userAgent)
	}
}

func TestDo(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	type account struct {
		ID string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"ID":"1234567890"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(account)
	_, _ = client.Do(req, body)

	want := &account{"1234567890"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpBadRequest(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	resp, err := client.Do(req, nil)

	if err == nil {
		t.Fatal("Expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

func TestDo_redirectLoop(t *testing.T) {
	client, mux, teardown := setupSendingClient()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "", http.StatusFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestDo_decodeUndefinedType(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "From: <jd@example.com> To: <info@example.com> Subject: Hello, world!")
	})

	req, _ := client.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "image/jpeg")

	body := new(int)
	if _, err := client.Do(req, body); err == nil {
		t.Error("Expected error to be returned.")
	}
}

func TestCheckResponse(t *testing.T) {
	t.Skip()
}
