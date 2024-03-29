package mailtrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjectsService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Project{}, "{}")

	u := projectMock(1)
	want := `{
		"id": 1,
		"name": "project-1",
		"share_links": {
			"admin": "https://localhost/projects/1/share/foo",
			"viewer": "https://localhost/projects/1/share/bar"
		},
		"inboxes": [
			{
				"id": 2,
				"name": "inbox-1",
				"username": "username",
				"password": "password",
				"max_size": 3,
				"status": "active",
				"email_username": "email-username",
				"email_username_enabled": true,
				"sent_messages_count": 4,
				"forwarded_messages_count": 5,
				"used": true,
				"forward_from_email_address": "forward@example.com",
				"project_id": 1,
				"domain": "localhost",
				"pop3_domain": "pop3-domain",
				"email_domain": "email-domain",
				"emails_count": 6,
				"emails_unread_count": 7,
				"last_message_sent_at": null,
				"smtp_ports": [
				  25,
				  2525
				],
				"pop3_ports": [
				  1100,
				  9950
				],
				"max_message_size": 100,
				"permissions": {
				  "can_read": true,
				  "can_update": true,
				  "can_destroy": true,
				  "can_leave": false
				}
			}
		],
		"permissions": {
			"can_read": true,
			"can_update": true,
			"can_destroy": true,
			"can_leave": true
		}
	}`
	testJSONMarshal(t, u, want)
}

func TestProjectsService_List(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	expectedProjects := []*Project{
		{
			ID:   1,
			Name: "project-1",
		},
		{
			ID:   2,
			Name: "project-2",
		},
	}

	mux.HandleFunc("/accounts/1/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		resp, _ := json.Marshal(expectedProjects)
		fmt.Fprint(w, string(resp))
	})

	projects, _, err := client.Projects.List(1)
	if err != nil {
		t.Errorf("Projects.List returned error: %v", err)
	}

	if !reflect.DeepEqual(projects, expectedProjects) {
		t.Errorf("Projects.List returned %+v, expected %+v", projects, expectedProjects)
	}

	testBadPathParams(t, "Projects.List", func() error {
		_, _, err = client.Projects.List(-1)
		return err
	})

	testNewRequestAndDoFail(t, "Projects.List", &client.client, func() (*Response, error) {
		project, resp, err := client.Projects.List(1)
		if project != nil {
			t.Errorf("Projects.List client.BaseURL.Host=%v project=%#v, want nil", client.baseURL.Host, project)
		}
		return resp, err
	})
}

func TestProjectsService_Get(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	mux.HandleFunc("/accounts/1/projects/20", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":20, "name":"project", "inboxes":[]}`)
	})

	project, _, err := client.Projects.Get(1, 20)
	if err != nil {
		t.Errorf("Projects.Get returned error: %v", err)
	}

	expected := &Project{ID: 20, Name: "project", Inboxes: []Inbox{}}
	if !reflect.DeepEqual(project, expected) {
		t.Errorf("Projects.Get returned %+v, expected %+v", project, expected)
	}

	testBadPathParams(t, "Projects.Get", func() error {
		_, _, err = client.Projects.Get(1, -20)
		return err
	})

	testNewRequestAndDoFail(t, "Projects.Get", &client.client, func() (*Response, error) {
		project, resp, err := client.Projects.Get(1, 2)
		if project != nil {
			t.Errorf("Projects.Get client.BaseURL.Host=%v project=%#v, want nil", client.baseURL.Host, project)
		}
		return resp, err
	})
}

func TestProjectsService_Create(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	var name = "Project name"

	mux.HandleFunc("/accounts/1/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintf(w, `{"id":1, "name":"%s"}`, name)
	})

	project, _, err := client.Projects.Create(1, name)
	if err != nil {
		t.Errorf("Projects.Create returned error: %v", err)
	}

	expected := &Project{ID: 1, Name: name}
	if !reflect.DeepEqual(project, expected) {
		t.Errorf("Projects.Create returned %+v, expected %+v", project, expected)
	}

	testBadPathParams(t, "Projects.Create", func() error {
		_, _, err = client.Projects.Create(-1, "")
		return err
	})

	testNewRequestAndDoFail(t, "Projects.Create", &client.client, func() (*Response, error) {
		project, resp, err := client.Projects.Create(1, "")
		if project != nil {
			t.Errorf("Projects.Create client.BaseURL.Host=%v project=%#v, want nil", client.baseURL.Host, project)
		}
		return resp, err
	})
}

func TestProjectsService_Update(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	var name = "New project name"

	mux.HandleFunc("/accounts/1/projects/21", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"name":"New project name"}`)
	})

	project, _, err := client.Projects.Update(1, 21, name)
	if err != nil {
		t.Errorf("Projects.Update returned error: %v", err)
	}

	expected := &Project{Name: name}
	if !reflect.DeepEqual(project, expected) {
		t.Errorf("Projects.Update returned %+v, expected %+v", project, expected)
	}

	testBadPathParams(t, "Projects.Update", func() error {
		_, _, err = client.Projects.Update(1, -20, "")
		return err
	})

	testNewRequestAndDoFail(t, "Projects.Update", &client.client, func() (*Response, error) {
		project, resp, err := client.Projects.Update(1, 2, "")
		if project != nil {
			t.Errorf("Projects.Update client.BaseURL.Host=%v project=%#v, want nil", client.baseURL.Host, project)
		}
		return resp, err
	})
}

func TestProjectsService_Delete(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	mux.HandleFunc("/accounts/1/projects/20", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"id":20}`)
	})

	resp, err := client.Projects.Delete(1, 20)
	if err != nil {
		t.Errorf("Projects.Delete returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Projects.Delete wrong status code: %d. Expected %d", resp.StatusCode, http.StatusOK)
	}

	testNewRequestAndDoFail(t, "Projects.Delete", &client.client, func() (*Response, error) {
		return client.Projects.Delete(1, 2)
	})
}

func projectMock(ID int) *Project {
	return &Project{
		ID:   ID,
		Name: "project-1",
		ShareLinks: struct {
			Admin  string `json:"admin"`
			Viewer string `json:"viewer"`
		}{
			Admin:  "https://localhost/projects/1/share/foo",
			Viewer: "https://localhost/projects/1/share/bar",
		},
		Inboxes: []Inbox{
			{
				ID:                      2,
				Name:                    "inbox-1",
				Username:                "username",
				Password:                "password",
				MaxSize:                 3,
				Status:                  "active",
				EmailUsername:           "email-username",
				EmailUsernameEnabled:    true,
				SentMessagesCount:       4,
				ForwardedMessagesCount:  5,
				Used:                    true,
				ForwardFromEmailAddress: "forward@example.com",
				ProjectID:               1,
				Domain:                  "localhost",
				POP3Domain:              "pop3-domain",
				EmailDomain:             "email-domain",
				EmailsCount:             6,
				EmailsUnreadCount:       7,
				LastMessageSentAt:       "",
				SMTPPorts:               []int{25, 2525},
				POP3Ports:               []int{1100, 9950},
				MaxMessageSize:          100,
				Permissions: Permissions{
					CanRead:    true,
					CanUpdate:  true,
					CanDestroy: true,
					CanLeave:   false,
				},
			},
		},
		Permissions: Permissions{
			CanRead:    true,
			CanUpdate:  true,
			CanDestroy: true,
			CanLeave:   true,
		},
	}
}
