package mailtrap

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAccountUsersService_Marshal(t *testing.T) {
	testJSONMarshal(t, &AccountUser{}, "{}")

	u := &AccountUser{
		ID:            1,
		SpecifierType: "user",
		Resources: []AccountUserResources{
			{
				ResourceType: "account",
				ResourceID:   2,
				AccessLevel:  1000,
			},
		},
		Specifier: AccountUserSpecifier{
			ID:    3,
			Email: "jd@example.com",
			Name:  "John",
		},
		Permissions: Permissions{
			CanRead:    true,
			CanUpdate:  false,
			CanDestroy: false,
			CanLeave:   false,
		},
	}
	want := `{
		"id": 1,
		"specifier_type": "user",
		"resources": [
		  {
			"resource_type": "account",
			"resource_id": 2,
			"access_level": 1000
		  }
		],
		"specifier": {
		  "id": 3,
		  "email": "jd@example.com",
		  "name": "John"
		},
		"permissions": {
		  "can_read": true,
		  "can_update": false,
		  "can_destroy": false,
		  "can_leave": false
		}
	}`
	testJSONMarshal(t, u, want)
}

func TestAccountUsersService_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/account_accesses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1, "specifier":{"id":1}},{"id":2, "specifier":{"id":2}}]`)
	})

	accountUsers, _, err := client.AccountUsers.List(1, nil)
	if err != nil {
		t.Errorf("Accounts.ListAccounts returned error: %v", err)
	}

	want := []*AccountUser{
		{ID: 1, Specifier: AccountUserSpecifier{ID: 1}},
		{ID: 2, Specifier: AccountUserSpecifier{ID: 2}},
	}
	if !reflect.DeepEqual(want, accountUsers) {
		t.Errorf("AccountUsers.List returned %+v, want %+v", accountUsers, want)
	}
}

func TestAccountUsersService_List_withQueryParams(t *testing.T) {
	t.Skip()
}

func TestAccountUsersService_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/account_accesses/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.AccountUsers.Delete(1, 2)
	if err != nil {
		t.Errorf("AccountUsers.Delete returned error: %v", err)
	}
}
