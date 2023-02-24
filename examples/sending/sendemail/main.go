// The command uses package as a CLI tool to show an example of how to send an email.
// It takes an auth token as an env variable and sends the email.
//
// It's runnable with the following command:
//
// export GITHUB_TOKEN=your_api_key
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vorobeyme/mailtrap-go/mailtrap"
)

func main() {
	apiKey := os.Getenv("MAILTRAP_API_KEY")
	if apiKey == "" {
		log.Fatal("No API key present")
	}
	client, err := mailtrap.NewSendingClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	resp, _, err := client.Send(emailRequest())
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}
	fmt.Printf("List of delivered message IDs: %#v \n", resp)
}

func emailRequest() *mailtrap.SendEmailRequest {
	return &mailtrap.SendEmailRequest{
		From: mailtrap.EmailAddress{
			Email: "ches@example.com",
			Name:  "Ches",
		},
		To: []mailtrap.EmailAddress{
			{
				Email: "doe@example.com",
				Name:  "John Doe",
			},
			{
				Email: "smith@example.com",
				Name:  "John Smith",
			},
		},
		Cc: []mailtrap.EmailAddress{
			{
				Email: "email.cc@example.com",
			},
		},
		Bcc: []mailtrap.EmailAddress{
			{
				Email: "email.bcc@example.com",
			},
		},
		Attachments: []mailtrap.EmailAttachment{
			{
				Content:     "PGh0bWw+CiAgICA8aGVhZD4KICAgICAgICA8dGl0bGU+YjY0PC90aXRsZT4KICAgIDwvaGVhZD4KICAgIDxib2R5PgogICAgPHA+SGVsbG8sIHdvcmxkITwvcD4KICAgIDwvYm9keT4KPC9odG1sPg==",
				AttachType:  "text/html",
				Filename:    "index.html",
				Disposition: "attachment",
			},
		},
		CustomVars: map[string]string{
			"user_id":  "1",
			"batch_id": "2",
		},
		Headers: map[string]string{
			"X-Message-Source": "mail.example.com",
		},
		Subject:  "API Client Test",
		Text:     "Hello, world!",
		Category: "API Client",
	}
}
