package mailtrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAttachmentsService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Attachment{}, "{}")

	u := attachment(1)
	want := `{
		"id": 1,
		"message_id": 2,
		"filename": "test.csv",
		"attachment_type": "inline",
		"content_type": "plain/text",
		"content_id": null,
		"transfer_encoding": null,
		"attachment_size": 0,
		"created_at": "2023-02-13T21:05:55.687Z",
		"updated_at": "2023-02-13T21:05:55.687Z",
		"attachment_human_size": "0 Bytes",
		"download_path": "/api/accounts/1/inboxes/2/messages/3/attachments/4/download"
	}`
	testJSONMarshal(t, u, want)
}

func TestAttachmentsService_List(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	expectedAttachments := []*Attachment{attachment(1), attachment(2)}

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/attachments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		res, _ := json.Marshal(expectedAttachments)
		fmt.Fprint(w, string(res))
	})

	attachments, _, err := client.Attachments.List(1, 2, 3)
	if err != nil {
		t.Errorf("Attachments.List returned error: %v", err)
	}

	if !reflect.DeepEqual(attachments, expectedAttachments) {
		t.Errorf("Attachments.List returned %+v, expected %+v", attachments, expectedAttachments)
	}

	_, _, err = client.Attachments.List(-1, -2, -3)
	if err == nil {
		t.Error("Attachments.List bad params err = nil, want error")
	}

	client.baseURL.Host = "!@#$%^&*()_+"
	attach, resp, err := client.Attachments.List(1, 2, 3)

	if attach != nil {
		t.Errorf("Attachments.List client.BaseURL.Host='invalid' attach = %#v, want nil", attach)
	}
	if resp != nil {
		t.Errorf("Attachments.List client.BaseURL=Host='invalid' resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("Attachments.List client.BaseURL=Host='invalid' err = nil, want error")
	}
}

func TestAttachmentsService_Get(t *testing.T) {
	client, mux, teardown := setupTestingClient()
	defer teardown()

	expectedAttachment := attachment(1)

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/attachments/4", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		res, _ := json.Marshal(expectedAttachment)
		fmt.Fprint(w, string(res))
	})

	attachments, _, err := client.Attachments.Get(1, 2, 3, 4)
	if err != nil {
		t.Errorf("Attachments.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(attachments, expectedAttachment) {
		t.Errorf("Attachments.Get returned %+v, expected %+v", attachments, expectedAttachment)
	}

	_, _, err = client.Attachments.Get(-1, -2, -3, -4)
	if err == nil {
		t.Error("Attachments.Get bad params err = nil, want error")
	}

	client.baseURL.Host = "!@#$%^&*()_+"
	attach, resp, err := client.Attachments.Get(1, 2, 3, 4)

	if attach != nil {
		t.Errorf("Attachments.Get client.BaseURL.Host='invalid' attach = %#v, want nil", attach)
	}
	if resp != nil {
		t.Errorf("Attachments.List client.BaseURL=Host='invalid' resp = %#v, want nil", resp)
	}
	if err == nil {
		t.Error("Attachments.Get client.BaseURL=Host='invalid' err = nil, want error")
	}
}

func TestAttachmentsService_Get_notFound(t *testing.T) {
	t.Skip()
}

func attachment(ID int) *Attachment {
	return &Attachment{
		ID:                  ID,
		MessageID:           2,
		Filename:            "test.csv",
		AttachmentType:      "inline",
		ContentType:         "plain/text",
		ContentID:           "",
		TransferEncoding:    "",
		AttachmentSize:      0,
		CreatedAt:           "2023-02-13T21:05:55.687Z",
		UpdatedAt:           "2023-02-13T21:05:55.687Z",
		AttachmentHumanSize: "0 Bytes",
		DownloadPath:        "/api/accounts/1/inboxes/2/messages/3/attachments/4/download",
	}
}
