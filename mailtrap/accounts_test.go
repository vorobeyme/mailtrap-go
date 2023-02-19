package mailtrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccountsService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Account{}, "{}")

	u := &Account{
		ID:           1,
		Name:         "account-1",
		AccessLevels: []int{100},
	}
	want := `{
		"id": 1,
		"name": "account-1",	
		"access_levels": [
			100
		]
	}`
	testJSONMarshal(t, u, want)
}

func TestAccountsService_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	expectedAccounts := []*Account{
		{
			ID:           1,
			Name:         "account-1",
			AccessLevels: []int{100},
		},
		{
			ID:           2,
			Name:         "account-2",
			AccessLevels: []int{1000},
		},
	}

	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		res, _ := json.Marshal(expectedAccounts)
		fmt.Fprint(w, string(res))
	})

	accounts, _, err := client.Accounts.List()
	if err != nil {
		t.Errorf("Accounts.List returned error: %v", err)
	}

	if !reflect.DeepEqual(accounts, expectedAccounts) {
		t.Errorf("Accounts.List returned accounts %+v, expected %+v", accounts, expectedAccounts)
	}

	testNewRequestAndDoFail(t, "Accounts.List", client, func() (*Response, error) {
		acc, resp, err := client.Accounts.List()
		if acc != nil {
			t.Errorf("Accounts.List client.BaseURL.Host=%v acc=%#v, want nil", client.defaultBaseURL.Host, acc)
		}
		return resp, err
	})
}
