package mailtrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestMessagesService_Marshal(t *testing.T) {
	testJSONMarshal(t, &Message{}, "{}")

	u := messageMock(1)
	want := `{		
		"id": 1,
		"inbox_id": 2,
		"subject": "Test email",
		"sent_at": "2023-02-14T19:29:59.295Z",
		"from_email": "john@example.com",
		"from_name": "John",
		"to_email": "mary@xample.com",
		"to_name": "Mary",
		"email_size": 30,
		"is_read": false,
		"created_at": "2023-02-14T19:29:59.295Z",
		"updated_at": "2023-02-14T19:29:59.295Z",
		"html_body_size": 200,
		"text_body_size": 100,
		"human_size": "300 Bytes",
		"html_path": "/api/accounts/1/inboxes/2/messages/3/body.html",
		"txt_path": "/api/accounts/1/inboxes/2/messages/3/body.txt",
		"raw_path": "/api/accounts/1/inboxes/2/messages/3/body.raw",
		"download_path": "/api/accounts/1/inboxes/2/messages/3/body.eml",
		"html_source_path": "/api/accounts/1/inboxes/2/messages/3/body.htmlsource",
		"blacklists_report_info": false,
		"smtp_information": {
			"ok": true,
			"data": {
				"mail_from_addr": "john@xample.com",
				"client_ip": "127.0.0.1"
			}
		}
	}`
	testJSONMarshal(t, u, want)
}

func TestMessagesService_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	expectedMessages := []*Message{messageMock(1), messageMock(2)}

	mux.HandleFunc("/accounts/1/inboxes/2/messages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		resp, _ := json.Marshal(expectedMessages)
		fmt.Fprint(w, string(resp))
	})

	messages, _, err := client.Messages.List(1, 2)
	if err != nil {
		t.Errorf("Messages.List returned error: %v", err)
	}

	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Errorf("Messages.List returned %+v, expected %+v", messages, expectedMessages)
	}
}

func TestMessagesService_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	expectedMessage := messageMock(1)

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		resp, _ := json.Marshal(expectedMessage)
		fmt.Fprint(w, string(resp))
	})

	message, _, err := client.Messages.Get(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(message, expectedMessage) {
		t.Errorf("Messages.Get returned %+v, expected %+v", message, expectedMessage)
	}
}

func TestMessagesService_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	updReq := &UpdateMessageRequest{
		IsRead: false,
	}

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"ID":3,"is_read":false}`)
	})

	message, _, err := client.Messages.Update(1, 2, 3, updReq)
	if err != nil {
		t.Errorf("Messages.Update returned error: %v", err)
	}

	expected := &Message{ID: 3, IsRead: false}
	if !reflect.DeepEqual(message, expected) {
		t.Errorf("Messages.Update returned %+v, expected %+v", message, expected)
	}
}

func TestMessagesService_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Messages.Delete(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.Delete returned error: %v", err)
	}
}

func TestMessagesService_Forward(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/forward", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"message":"Your email message has been successfully forwarded"}`)
	})

	_, err := client.Messages.Forward(1, 2, 3, "email@example.com")
	if err != nil {
		t.Errorf("Messages.Forward returned error: %v", err)
	}
}

func TestMessagesService_Forward_invalidEmail(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	_, err := client.Messages.Forward(1, 2, 3, "emailexample.com")

	const errMessage = "forward 'email' is invalid"
	if err.Error() != errMessage {
		t.Errorf("Messages.Forward error is %v, want %s", err, errMessage)
	}
}

func TestMessagesService_SpamReport(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/spam_report", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"report":{"ResponseCode":1,"Spam":false}}`)
	})

	report, _, err := client.Messages.SpamReport(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.SpamReport returned error: %v", err)
	}

	expected := new(SpamReport)
	expected.Report.ResponseCode = 1
	expected.Report.Spam = false

	if !reflect.DeepEqual(report, expected) {
		t.Errorf("Messages.SpamReport returned %+v, expected %+v", report, expected)
	}
}

func TestMessagesService_AsRaw(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	rawBody := `
	From: Ches Sparrow <ches@example.com>
	To: John Doe <jd@example.com>
	Subject: You are awesome!
	Content-Type: multipart/alternative; boundary="boundary-string"

	--boundary-string
	Content-Type: text/plain; charset="utf-8"
	Content-Transfer-Encoding: quoted-printable
	Content-Disposition: inline

	Congrats for sending test email with Mailtrap!

	Inspect it using the tabs above and learn how this email can be improved.
	Now send your email using our fake SMTP server and integration of your choice!

	Good luck! Hope it works.
	`

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/body.raw", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "text/plain")
		fmt.Fprint(w, rawBody)
	})

	rawResp, _, err := client.Messages.AsRaw(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.AsRaw returned error: %v", err)
	}
	if !reflect.DeepEqual(rawResp, rawBody) {
		t.Errorf("Messages.AsRaw returned %+v, expected %+v", rawResp, rawBody)
	}
}

func TestMessagesService_AsText(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	textBody := `
	Congrats for sending test email with Mailtrap!

	Inspect it using the tabs above and learn how this email can be improved.
	Now send your email using our fake SMTP server and integration of your choice!

	Good luck! Hope it works.
	`

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/body.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "text/plain")
		fmt.Fprint(w, textBody)
	})

	textResp, _, err := client.Messages.AsText(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.AsText returned error: %v", err)
	}
	if !reflect.DeepEqual(textResp, textBody) {
		t.Errorf("Messages.AsText returned %+v, expected %+v", textResp, textBody)
	}
}

func TestMessagesService_AsHTML(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	htmlBody := `
	<!doctype html>
	<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	</head>
	<body style="font-family: sans-serif;">
		<div style="display: block; margin: auto; max-width: 600px;" class="main">
		<h1 style="font-size: 18px; font-weight: bold; margin-top: 20px">
			Congrats for sending test email with Mailtrap!
		</h1>
		<p>Inspect it using the tabs you see above and learn how this email can be improved.</p>
		<img alt="Alt" src="https://assets.examples.com/integration-examples/welcome.png">
		<p>Now send your email using our fake SMTP server and integration of your choice!</p>
		<p>Good luck! Hope it works.</p>
		</div>
		<style> .main { background-color: white; } a:hover
		{ border-left-width: 1em; min-height: 2em; }
		</style>
	</body>
	</html>
	`

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/body.html", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "text/html")
		fmt.Fprint(w, htmlBody)
	})

	htmlResp, _, err := client.Messages.AsHTML(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.AsHTML returned error: %v", err)
	}
	if !reflect.DeepEqual(htmlResp, htmlBody) {
		t.Errorf("Messages.AsHTML returned %+v, expected %+v", htmlResp, htmlBody)
	}
}

func TestMessagesService_AsHTMLSource(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	htmlSrcBody := `
	<!doctype html>
	<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
		</head>
		<body style="font-family: sans-serif;">
		<div style="display: block; margin: auto; max-width: 600px;" class="main">
			<h1 style="font-size: 18px; font-weight: bold; margin-top: 20px">
				Congrats for sending test email with Mailtrap!
			</h1>
			<p>Inspect it using the tabs you see above and learn how this email can be improved.</p>
			<img alt="Alt" src="https://assets.examples.com/integration-examples/welcome.png">
			<p>Now send your email using our fake SMTP server and integration of your choice!</p>
			<p>Good luck! Hope it works.</p>
		</div>
		<style> .main { background-color: white; } a:hover
			{ border-left-width: 1em; min-height: 2em; }
		</style>
		</body>
	</html>
	`

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/body.htmlsource", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "text/html")
		fmt.Fprint(w, htmlSrcBody)
	})

	htmlSrcResp, _, err := client.Messages.AsHTMLSource(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.AsHTMLSource returned error: %v", err)
	}
	if !reflect.DeepEqual(htmlSrcResp, htmlSrcBody) {
		t.Errorf("Messages.AsHTMLSource returned %+v, expected %+v", htmlSrcResp, htmlSrcBody)
	}
}

func TestMessagesService_AsEml(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	emlBody := `
	From: Ches Sparrow <ches@example.com>
	To: John Doe <jd@example.com>
	Subject: You are awesome!
	Content-Type: multipart/alternative; boundary="boundary-string"

	--boundary-string
	Content-Type: text/plain; charset="utf-8"
	Content-Transfer-Encoding: quoted-printable
	Content-Disposition: inline

	Congrats for sending test email with Mailtrap!

	Inspect it using the tabs above and learn how this email can be improved.
	Now send your email using our fake SMTP server and integration of your choice!

	Good luck! Hope it works.

	--boundary-string
	Content-Type: text/html; charset="utf-8"
	Content-Transfer-Encoding: quoted-printable
	Content-Disposition: inline

	<!doctype html>
	<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	</head>
	<body>
		<p>Now send your email using our fake SMTP server and integration of your choice!</p>
	</body>
	</html>
	--boundary-string--
	`

	mux.HandleFunc("/accounts/1/inboxes/2/messages/3/body.eml", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "message/rfc822")
		fmt.Fprint(w, emlBody)
	})

	emlResp, _, err := client.Messages.AsEML(1, 2, 3)
	if err != nil {
		t.Errorf("Messages.AsEML returned error: %v", err)
	}
	if !reflect.DeepEqual(emlResp, emlBody) {
		t.Errorf("Messages.AsEML returned %+v, expected %+v", emlResp, emlBody)
	}
}

func messageMock(ID int) *Message {
	var smtp = new(MessageSMTPInfo)
	smtp.Ok = true
	smtp.Data.MailFromAddr = "john@xample.com"
	smtp.Data.ClientIP = "127.0.0.1"

	datetime, _ := time.Parse(time.RFC3339, "2023-02-14T19:29:59.295Z")

	return &Message{
		ID:                   ID,
		InboxID:              2,
		Subject:              "Test email",
		SentAt:               datetime,
		FromEmail:            "john@example.com",
		FromName:             "John",
		ToEmail:              "mary@xample.com",
		ToName:               "Mary",
		EmailSize:            30,
		IsRead:               false,
		CreatedAt:            datetime,
		UpdatedAt:            datetime,
		HTMLBodySize:         200,
		TextBodySize:         100,
		HumanSize:            "300 Bytes",
		HTMLPath:             "/api/accounts/1/inboxes/2/messages/3/body.html",
		TxtPath:              "/api/accounts/1/inboxes/2/messages/3/body.txt",
		RawPath:              "/api/accounts/1/inboxes/2/messages/3/body.raw",
		DownloadPath:         "/api/accounts/1/inboxes/2/messages/3/body.eml",
		HTMLSourcePath:       "/api/accounts/1/inboxes/2/messages/3/body.htmlsource",
		BlacklistsReportInfo: false,
		SMTPInfo:             smtp,
	}
}
