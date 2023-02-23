// An example of how to use Messages methods.
//
// It's runnable with the following command:
// export MAILTRAP_API_KEY=your_api_key
// go run .
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vorobeyme/mailtrap-go/mailtrap"
)

var client *mailtrap.TestingClient

func main() {
	apiKey := os.Getenv("MAILTRAP_API_KEY")
	if apiKey == "" {
		log.Fatal("No API key present")
	}
	client, _ = mailtrap.NewTestingClient(apiKey)

	for {
		var method string
		fmt.Print("Enter method [list, get, update, delete, forward, report, raw, text, html, source, eml] or `q` for exit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		// Get all messages in inboxes.
		case "list":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)

			messages, err := list(accountID, inboxID)
			if err != nil {
				log.Fatalf("list: %v\n", err)
			}
			if len(messages) == 0 {
				fmt.Println("Inbox is empty.")
				continue
			}
			for _, m := range messages {
				displayMessage(m)
			}

		// Get email message by ID.
		case "get":
			var accountID, inboxID, messageID int
			fmt.Print("Enter an accountID, inboxID and messageID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID)
			m, err := get(accountID, inboxID, messageID)
			if err != nil {
				log.Fatalf("get: %v\n", err)
			}
			displayMessage(m)

		// Update message attributes (only the is_read attribute).
		case "update":
			var (
				accountID, inboxID, messageID int
				isRead                        bool
			)
			fmt.Print("Enter an accountID, inboxID and messageID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID)
			fmt.Print("Enter an isRead bool value: ")
			fmt.Scanf("%t", &isRead)
			m, err := update(accountID, inboxID, messageID, isRead)
			if err != nil {
				log.Fatalf("update: %v\n", err)
			}
			displayMessage(m)

		// Delete message from inbox.
		case "delete":
			var accountID, inboxID, messageID int
			fmt.Print("Enter an accountID, inboxID and messageID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID)
			err := delete(accountID, inboxID, messageID)
			if err != nil {
				log.Fatalf("delete: %v\n", err)
			}
			fmt.Println("message deleted.")

		// Forward message to an email address.
		// The email address must be confirmed by the recipient in advance.
		case "forward":
			var (
				accountID, inboxID, messageID int
				email                         string
			)
			fmt.Print("Enter an accountID, inboxID and messageID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID)
			fmt.Print("Enter an forward email value: ")
			fmt.Scanf("%s", &email)

			err := forward(accountID, inboxID, messageID, email)
			if err != nil {
				log.Fatalf("forward: %v\n", err)
			}
			fmt.Println("message forwarded.")

		// Get a brief spam report by message ID.
		case "report":
			var accountID, inboxID, messageID int
			fmt.Print("Enter an accountID, inboxID and messageID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID)
			r, err := spamReport(accountID, inboxID, messageID)
			if err != nil {
				log.Fatalf("report: %v\n", err)
			}

			fmt.Printf(
				"\tCode: %v \n\tMessage: %v \n\tVersion: %v \n\tScore: %v \n\tSpam: %v \n\tThreshold: %v \n\tDetails: %v",
				r.Report.ResponseCode, r.Report.ResponseMessage, r.Report.ResponseVersion, r.Report.Score,
				r.Report.Spam, r.Report.Threshold, r.Report.Details,
			)

		// Get message body
		case "raw", "text", "html", "source", "eml":
			var accountID, inboxID, messageID int
			fmt.Print("Enter an accountID, inboxID and messageID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID)
			m, err := viewAs(accountID, inboxID, messageID, method)
			if err != nil {
				log.Fatalf("email body: %v\n", err)
			}
			fmt.Printf("%q", m)

		default:
			fmt.Println("method doesn't exist...")
		}
	}
}

func list(accID, inboxID int) ([]*mailtrap.Message, error) {
	messages, _, err := client.Messages.List(accID, inboxID)
	if err != nil {
		log.Fatal(err)
	}
	return messages, err
}

func get(accID, inboxID, messageID int) (*mailtrap.Message, error) {
	message, _, err := client.Messages.Get(accID, inboxID, messageID)
	if err != nil {
		log.Fatal(err)
	}
	return message, err
}

func update(accID, inboxID, messageID int, isRead bool) (*mailtrap.Message, error) {
	message, _, err := client.Messages.Update(accID, inboxID, messageID, &mailtrap.UpdateMessageRequest{IsRead: isRead})
	if err != nil {
		log.Fatal(err)
	}
	return message, err
}

func delete(accID, inboxID, messageID int) error {
	_, err := client.Messages.Delete(accID, inboxID, messageID)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func forward(accID, inboxID, messageID int, email string) error {
	_, err := client.Messages.Forward(accID, inboxID, messageID, email)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func spamReport(accID, inboxID, messageID int) (*mailtrap.SpamReport, error) {
	report, _, err := client.Messages.SpamReport(accID, inboxID, messageID)
	if err != nil {
		log.Fatal(err)
	}
	return report, nil
}

func viewAs(accID, inboxID, messageID int, asType string) (string, error) {
	var (
		body string
		err  error
	)
	switch asType {
	case "raw":
		body, _, err = client.Messages.AsRaw(accID, inboxID, messageID)
	case "text":
		body, _, err = client.Messages.AsText(accID, inboxID, messageID)
	case "html":
		body, _, err = client.Messages.AsHTML(accID, inboxID, messageID)
	case "source":
		body, _, err = client.Messages.AsHTMLSource(accID, inboxID, messageID)
	case "eml":
		body, _, err = client.Messages.AsEML(accID, inboxID, messageID)
	default:
		return "", errors.New("undefined body type")
	}
	return body, err
}

func displayMessage(m *mailtrap.Message) {
	fmt.Printf(
		"\tID: %v \n\tInbox ID: %v \n\tSubject: %v \n\tSent at: %v \n\tFrom email: %v \n\tFrom name: %v \n\tTo email: %v"+
			"\n\tTo name: %v \n\tEmail size: %v \n\tIs read: %v \n\tCreated at: %v \n\tUpdated at: %v"+
			"\n\tHTML body size: %v \n\tText body size: %v \n\tHuman size: %v \n\tHTML path: %v \n\tTxt path: %v"+
			"\n\tRaw path: %v \n\tDownload path: %v \n\tHTML source path: %v \n\tBlacklists report info: %v \n\tSMTP information: %v",
		m.ID, m.InboxID, m.Subject, m.SentAt, m.FromEmail, m.FromName, m.ToEmail, m.ToName, m.EmailSize, m.IsRead, m.CreatedAt,
		m.UpdatedAt, m.HTMLBodySize, m.TextBodySize, m.HumanSize, m.HTMLPath, m.TxtPath, m.RawPath, m.DownloadPath, m.HTMLSourcePath,
		m.BlacklistsReportInfo, m.SMTPInfo,
	)
}
