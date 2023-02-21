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

var client *mailtrap.Client

func main() {
	apiKey := os.Getenv("MAILTRAP_API_KEY")
	if apiKey == "" {
		log.Fatal("No API key present")
	}
	client = mailtrap.New(apiKey)

	for {
		var method string
		fmt.Print("Enter method [list, get, create, update, delete, clean, mark-read, reset-credentials, reset-email, enable-email] or `q` for exit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		// Get a list of inboxes.
		case "list":
			var accountID int
			fmt.Print("Enter an accountID value: ")
			fmt.Scanf("%d %d", &accountID)
			inboxes, _, err := client.Inboxes.List(accountID)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range inboxes {
				fmt.Printf("\n%s \n\n", displayInbox(i))
			}

		// Get inbox attributes by inbox id.
		case "get":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			inbox, _, err := client.Inboxes.Get(accountID, inboxID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\n%s \n\n", displayInbox(inbox))

		// Create an inbox in a project.
		case "create":
			var (
				accountID, inboxID int
				name               string
			)
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			fmt.Print("Enter a name value: ")
			fmt.Scanf("%s", &name)
			inbox, _, err := client.Inboxes.Create(accountID, inboxID, name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nCreated inbox:\n%s \n\n", displayInbox(inbox))

		//
		case "update":
			var (
				accountID, inboxID int
				name, email        string
			)
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			fmt.Print("Enter a name value: ")
			fmt.Scanf("%s", &name)
			fmt.Print("Enter an email username value: ")
			fmt.Scanf("%s", &email)
			inbox, _, err := client.Inboxes.Update(accountID, inboxID, &mailtrap.UpdateInboxRequest{
				Name:          name,
				EmailUsername: email,
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nUpdated inbox:\n%s \n\n", displayInbox(inbox))

		// Delete an inbox with all its emails.
		case "delete":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			_, err := client.Inboxes.Delete(accountID, inboxID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inbox deleted.")

		// Delete all messages (emails) from inbox.
		case "clean":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			inbox, _, err := client.Inboxes.Clean(accountID, inboxID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nInbox cleaned:\n%s \n\n", displayInbox(inbox))

		// Mark all messages in the inbox as read.
		case "mark-read":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			inbox, _, err := client.Inboxes.MarkAsRead(accountID, inboxID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nMark all messages as read:\n%s \n\n", displayInbox(inbox))

		// Reset SMTP credentials of the inbox.
		case "reset-credentials":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			inbox, _, err := client.Inboxes.ResetCredentials(accountID, inboxID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nReset credentials:\n%s \n\n", displayInbox(inbox))

		// Reset username of email address per inbox.
		case "reset-email":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			inbox, _, err := client.Inboxes.ResetEmail(accountID, inboxID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nReset username of email address:\n%s \n\n", displayInbox(inbox))

		// Turn the email address of the inbox on/off.
		case "enable-email":
			var accountID, inboxID int
			fmt.Print("Enter an accountID and inboxID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &inboxID)
			inbox, _, err := client.Inboxes.EnableEmail(accountID, inboxID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nToggle email:\n%s \n\n", displayInbox(inbox))

		default:
			fmt.Println("method doesn't exist...")
		}
	}
}

func displayInbox(i *mailtrap.Inbox) string {
	permissions := fmt.Sprintf(
		"\n\tCan read: %v \n\tCan update: %v \n\tCan destroy: %v \n\tCan leave: %v",
		i.Permissions.CanRead, i.Permissions.CanUpdate, i.Permissions.CanDestroy, i.Permissions.CanLeave,
	)

	return fmt.Sprintf(
		"ID: %v \nName: %v \nUsername: %v \nPassword: %v \nMaxs size: %v \nStatus: %v \nEmail username: %v"+
			"\nEmail username enabled: %v \nSend msg count: %v \nForwarded msg count: %v \nProject ID: %v \nDomain: %v"+
			"\nPOP3 domain: %v \nEmail domain: %v \nEmails count: %v \nEmails unread count: %v \nLast message sent at: %v"+
			"\nSMTP ports: %v \nPOP3 ports: %v \nMax msg size: %v \nPermissions: %v",
		i.ID, i.Name, i.Username, i.Password, i.MaxSize, i.Status, i.EmailUsername, i.EmailUsernameEnabled, i.SentMessagesCount,
		i.ForwardedMessagesCount, i.ProjectID, i.Domain, i.POP3Domain, i.EmailDomain, i.EmailsCount, i.EmailsUnreadCount, i.LastMessageSentAt,
		i.SMTPPorts, i.POP3Ports, i.MaxMessageSize, permissions,
	)
}
