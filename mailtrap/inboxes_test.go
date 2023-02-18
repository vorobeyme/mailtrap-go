package mailtrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestInboxesService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Inbox{}, "{}")

	u := inboxMock(1)
	want := `{
		"id": 1,
		"name": "inbox",
		"username": "username",
		"password": "pswd",
		"max_size": 0,
		"status": "active",
		"email_username": "emailusername",
		"email_username_enabled": false,
		"sent_messages_count": 100,
		"forwarded_messages_count": 0,
		"used": false,
		"forward_from_email_address": "frwd@example.com",	  
		"project_id": 2,
		"domain": "localhost",
		"pop3_domain": "localhost",
		"email_domain": "localhost",
		"emails_count": 10,
		"emails_unread_count": 0,
		"last_message_sent_at": null,
		"smtp_ports": [
		  25,
		  2525
		],
		"pop3_ports": [
		  1100
		],
		"max_message_size": 2000,
		"permissions": {
		  "can_read": true,
		  "can_update": false,
		  "can_destroy": false,
		  "can_leave": true
		}
	}`
	testJSONMarshal(t, u, want)
}

func TestInboxesService_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	expectedInboxes := []*Inbox{inboxMock(1), inboxMock(2)}

	mux.HandleFunc("/accounts/1/inboxes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		resp, _ := json.Marshal(expectedInboxes)
		fmt.Fprint(w, string(resp))
	})

	inboxes, _, err := client.Inboxes.List(1)
	if err != nil {
		t.Errorf("Inboxes.List returned error: %v", err)
	}

	if !reflect.DeepEqual(inboxes, expectedInboxes) {
		t.Errorf("Inboxes.List returned %+v, expected %+v", inboxes, expectedInboxes)
	}
}

func TestInboxesService_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	expectedInbox := inboxMock(1)

	mux.HandleFunc("/accounts/1/inboxes/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		resp, _ := json.Marshal(expectedInbox)
		fmt.Fprint(w, string(resp))
	})

	inbox, _, err := client.Inboxes.Get(1, 2)
	if err != nil {
		t.Errorf("Inboxes.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(inbox, expectedInbox) {
		t.Errorf("Inboxes.Get returned %+v, expected %+v", inbox, expectedInbox)
	}
}

func TestInboxesService_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"id":1}`)
	})

	_, err := client.Inboxes.Delete(1, 2)
	if err != nil {
		t.Errorf("Inboxes.Delete returned error: %v", err)
	}
}

func TestInboxesService_CreateInbox(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	var name = "inbox name"

	mux.HandleFunc("/accounts/1/projects/2/inboxes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintf(w, `{"id":1, "name":"%s"}`, name)
	})

	inbox, _, err := client.Inboxes.Create(1, 2, name)
	if err != nil {
		t.Errorf("Inboxes.Create returned error: %v", err)
	}

	expected := &Inbox{ID: 1, Name: name}
	if !reflect.DeepEqual(inbox, expected) {
		t.Errorf("Inboxes.Create returned %+v, expected %+v", inbox, expected)
	}
}

func TestInboxesService_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	var opts = UpdateInboxRequest{
		Name:          "new inbox name",
		EmailUsername: "username",
	}

	mux.HandleFunc("/accounts/1/inboxes/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":3,"name":"new inbox name","email_username":"username"}`)
	})

	inbox, _, err := client.Inboxes.Update(1, 3, &opts)
	if err != nil {
		t.Errorf("Inboxes.Update returned error: %v", err)
	}

	expected := &Inbox{ID: 3, Name: "new inbox name", EmailUsername: "username"}
	if !reflect.DeepEqual(inbox, expected) {
		t.Errorf("Inboxes.Update returned %+v, expected %+v", inbox, expected)
	}
}

func TestInboxesService_Clean(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/clean", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":3,"name":"inbox-clean","email_username":"username"}`)
	})

	inbox, _, err := client.Inboxes.Clean(1, 2)
	if err != nil {
		t.Errorf("Inboxes.Clean returned error: %v", err)
	}

	expected := &Inbox{ID: 3, Name: "inbox-clean", EmailUsername: "username"}
	if !reflect.DeepEqual(inbox, expected) {
		t.Errorf("Inboxes.Clean returned %+v, expected %+v", inbox, expected)
	}
}

func TestInboxesService_MarkAsRead(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/all_read", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":3,"name":"inbox-read","email_username":"username"}`)
	})

	inbox, _, err := client.Inboxes.MarkAsRead(1, 2)
	if err != nil {
		t.Errorf("Inboxes.MarkAsRead returned error: %v", err)
	}

	expected := &Inbox{ID: 3, Name: "inbox-read", EmailUsername: "username"}
	if !reflect.DeepEqual(inbox, expected) {
		t.Errorf("Inboxes.MarkAsRead returned %+v, expected %+v", inbox, expected)
	}
}

func TestInboxesService_ResetCredentials(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/reset_credentials", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":3,"name":"inbox-clean","email_username":"username"}`)
	})

	inbox, _, err := client.Inboxes.ResetCredentials(1, 2)
	if err != nil {
		t.Errorf("Inboxes.ResetCredentials returned error: %v", err)
	}

	expected := &Inbox{ID: 3, Name: "inbox-clean", EmailUsername: "username"}
	if !reflect.DeepEqual(inbox, expected) {
		t.Errorf("Inboxes.ResetCredentials returned %+v, expected %+v", inbox, expected)
	}
}

func TestInboxesService_EnableEmail(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/toggle_email_username", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":3,"name":"inbox-enable","email_username":"username"}`)
	})

	inbox, _, err := client.Inboxes.EnableEmail(1, 2)
	if err != nil {
		t.Errorf("Inboxes.EnableEmail returned error: %v", err)
	}

	expected := &Inbox{ID: 3, Name: "inbox-enable", EmailUsername: "username"}
	if !reflect.DeepEqual(inbox, expected) {
		t.Errorf("Inboxes.EnableEmail returned %+v, expected %+v", inbox, expected)
	}
}

func TestInboxesService_ResetEmail(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/reset_email_username", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":3,"name":"inbox-reset","email_username":"username"}`)
	})

	inbox, _, err := client.Inboxes.ResetEmail(1, 2)
	if err != nil {
		t.Errorf("Inboxes.EnableEmail returned error: %v", err)
	}

	expected := &Inbox{ID: 3, Name: "inbox-reset", EmailUsername: "username"}
	if !reflect.DeepEqual(inbox, expected) {
		t.Errorf("Inboxes.EnableEmail returned %+v, expected %+v", inbox, expected)
	}
}

func inboxMock(ID int) *Inbox {
	return &Inbox{
		ID:                      ID,
		Name:                    "inbox",
		Username:                "username",
		Password:                "pswd",
		MaxSize:                 0,
		Status:                  "active",
		EmailUsername:           "emailusername",
		EmailUsernameEnabled:    false,
		SentMessagesCount:       100,
		ForwardedMessagesCount:  0,
		Used:                    false,
		ForwardFromEmailAddress: "frwd@example.com",
		ProjectID:               2,
		Domain:                  "localhost",
		POP3Domain:              "localhost",
		EmailDomain:             "localhost",
		EmailsCount:             10,
		EmailsUnreadCount:       0,
		LastMessageSentAt:       "",
		SMTPPorts:               []int{25, 2525},
		POP3Ports:               []int{1100},
		MaxMessageSize:          2000,
		Permissions: Permissions{
			CanRead:    true,
			CanUpdate:  false,
			CanDestroy: false,
			CanLeave:   true,
		},
	}
}
