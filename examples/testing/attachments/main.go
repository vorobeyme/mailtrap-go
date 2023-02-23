// An example of how to use Messages methods.
//
// It's runnable with the following command:
// export MAILTRAP_API_KEY=your_api_key
// go run .
package main

import (
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
		fmt.Print("Enter method [list, get] or `q` for exit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		// Get message attachments by inbox_id and message_id.
		case "list":
			var accountID, inboxID, messageID int
			fmt.Print("Enter an accountID, inboxID and messageID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID)
			attach, _, err := client.Attachments.List(accountID, inboxID, messageID)
			if err != nil {
				log.Fatal(err)
			}
			for _, a := range attach {
				fmt.Printf("\n%s \n\n", displayAttach(a))
			}

		// Get message single attachment by id.
		case "get":
			var accountID, inboxID, messageID, attachmentID int
			fmt.Print("Enter an accountID, inboxID, messageID and attachmentID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID, &messageID, &attachmentID)
			attach, _, err := client.Attachments.Get(accountID, inboxID, messageID, attachmentID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\n%s \n\n", displayAttach(attach))

		default:
			fmt.Println("method doesn't exist...")
		}
	}
}

func displayAttach(a *mailtrap.Attachment) string {
	return fmt.Sprintf(
		"ID: %v \nMessage ID: %v \nFilename: %v \nAttachment type: %v \nContent type: %v \nContent ID: %v"+
			"\nTransfer encoding %v \nAttachment size: %v \nCreated at: %v \nUpdated at: %v \nAttachment human size: %v"+
			"\nDownload path: %v",
		a.ID, a.MessageID, a.Filename, a.AttachmentType, a.ContentType, a.ContentID, a.TransferEncoding, a.AttachmentSize,
		a.CreatedAt, a.UpdatedAt, a.AttachmentHumanSize, a.DownloadPath,
	)
}
