package mailtrap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestPermissionsService_Marshal(t *testing.T) {
	testJSONMarshal(t, Resource{}, "{}")

	u := &Resource{
		ID:          1,
		Name:        "permission-1",
		Type:        "project",
		AccessLevel: 100,
		Resource:    []*Resource{},
	}
	want := `{
		"id": 1,
		"name": "permission-1",
		"type": "project",
		"access_level": 100,
		"resources": []
	}`

	testJSONMarshal(t, u, want)
}

func TestPermissionsService_GetResources(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	expectedResources := []*Resource{
		{ID: 1, Name: "foo", Type: "account", AccessLevel: 1, Resource: []*Resource{}},
		{ID: 2, Name: "bar", Type: "project", AccessLevel: 100, Resource: []*Resource{}},
	}

	mux.HandleFunc("/accounts/1/permissions/resources", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		res, _ := json.Marshal(expectedResources)
		fmt.Fprint(w, string(res))
	})

	resources, _, err := client.Permissions.ListResources(1)
	if err != nil {
		t.Errorf("Permissions.ListResources returned error: %v", err)
	}

	if !reflect.DeepEqual(resources, expectedResources) {
		t.Errorf("Permissions.ListResources returned %+v, expected %+v", resources, expectedResources)
	}

	_, _, err = client.Permissions.ListResources(-1)
	if err == nil {
		t.Error("Permissions.ListResources bad params err = nil, want error")
	}

	client.baseURL.Host = "!@#$%^&*()_+"
	resources, resp, err := client.Permissions.ListResources(1)

	if resources != nil {
		t.Errorf("Permissions.ListResources client.BaseURL.Host='invalid' resources = %#v, want nil", resources)
	}
	if resp != nil {
		t.Errorf("Permissions.ListResources client.BaseURL=Host='invalid' resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("Permissions.ListResources client.BaseURL=Host='invalid' err = nil, want error")
	}
}

func TestPermissionsService_Manage(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	opt := &[]PermissionRequest{
		{ResourceID: 1, ResourceType: "inbox"},
		{ResourceID: 2, ResourceType: "project"},
	}

	mux.HandleFunc("/accounts/1/account_accesses/2/permissions/bulk", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		reqBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("Permissions.Manage mock did not work.")
		}

		req := strings.TrimSuffix(string(reqBytes), "\n")
		expectedReq := `{"permissions":[{"resource_id":1,"resource_type":"inbox"},{"resource_id":2,"resource_type":"project"}]}`

		if req != expectedReq {
			t.Errorf("Permissions.Manage expected req != req:\n expected %+v\n got %+v\n", expectedReq, req)
		}

		fmt.Fprint(w, `{"message":"Permissions have been updated!"}`)
	})

	_, err := client.Permissions.Manage(1, 2, opt)
	if err != nil {
		t.Errorf("Permissions.Manage returned error: %v", err)
	}

	_, err = client.Permissions.Manage(-1, -2, &[]PermissionRequest{{}})
	if err == nil {
		t.Error("Permissions.Manage bad params err = nil, want error")
	}

	client.baseURL.Host = "!@#$%^&*()_+"
	resp, err := client.Permissions.Manage(1, 2, nil)

	if resp != nil {
		t.Errorf("Permissions.Manage client.BaseURL=Host='invalid' resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("Permissions.Manage client.BaseURL=Host='invalid' err = nil, want error")
	}
}
