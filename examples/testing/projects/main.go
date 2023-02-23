// An example of how to use Project methods.
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
		fmt.Print("Enter method [list, get, create, update, delete] or `q` for exit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		// List projects and their inboxes to which the API token has access.
		case "list":
			var accountID int
			fmt.Print("Enter an accountID value: ")
			fmt.Scanf("%d", &accountID)
			projects, err := list(accountID)
			if err != nil {
				log.Fatalf("list: %v\n", err)
			}
			for _, p := range projects {
				displayProject(p)
			}

		// Get the project and its inboxes.
		case "get":
			var accountID, projectID int
			fmt.Print("Enter an accountID and projectID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &projectID)
			p, err := get(accountID, projectID)
			if err != nil {
				log.Fatalf("get: %v\n", err)
			}
			displayProject(p)

		// Create project.
		case "create":
			var accountID int
			fmt.Print("Enter an accountID value: ")
			fmt.Scanf("%d", &accountID)
			p, err := create(accountID, "Project name")
			if err != nil {
				log.Fatalf("create: %v\n", err)
			}
			displayProject(p)

		// Update project name.
		case "update":
			var (
				accountID, projectID int
				name                 string
			)
			fmt.Print("Enter an accountID and projectID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &projectID)
			fmt.Print("Enter a project name: ")
			fmt.Scanf("%s", &name)
			p, err := update(accountID, projectID, name)
			if err != nil {
				log.Fatalf("update: %v\n", err)
			}
			displayProject(p)

		// Delete project and its inboxes.
		case "delete":
			var accountID, projectID int
			fmt.Print("Enter an accountID and projectID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &projectID)
			err := delete(accountID, projectID)
			if err != nil {
				log.Fatalf("delete: %v\n", err)
			}
			fmt.Println("project deleted.")
		default:
			fmt.Println("method doesn't exist...")
		}
	}
}

func list(accID int) ([]*mailtrap.Project, error) {
	projects, _, err := client.Projects.List(accID)
	if err != nil {
		log.Fatal(err)
	}
	return projects, err
}

func get(accID, projectID int) (*mailtrap.Project, error) {
	project, _, err := client.Projects.Get(accID, projectID)
	if err != nil {
		log.Fatal(err)
	}
	return project, err
}

func create(accID int, name string) (*mailtrap.Project, error) {
	project, _, err := client.Projects.Create(accID, name)
	if err != nil {
		log.Fatal(err)
	}
	return project, err
}

func update(accID, projectID int, name string) (*mailtrap.Project, error) {
	project, _, err := client.Projects.Update(accID, projectID, name)
	if err != nil {
		log.Fatal(err)
	}
	return project, err
}

func delete(accID, projectID int) error {
	_, err := client.Projects.Delete(accID, projectID)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func displayProject(p *mailtrap.Project) {
	inbx := prepareInboxes(p.Inboxes)
	permission := preparePermissions(p.Permissions)
	fmt.Printf(
		"\nID: %v \nName: %v \nShare links: \n\tAdmin: %v \n\tViewer: %v \nInboxes: \n%v \nPermissions: \n%v\n\n",
		p.ID, p.Name, p.ShareLinks.Admin, p.ShareLinks.Viewer, strings.Join(inbx, "\n"), permission,
	)
}

func prepareInboxes(inboxes []mailtrap.Inbox) []string {
	var res []string
	for _, i := range inboxes {
		permissions := preparePermissions(i.Permissions)
		inbx := fmt.Sprintf(
			"\tID: %v \n\tName: %v \n\tUsername: %v \n\tPassword: %v \n\tMaxs size: %v \n\tStatus: %v \n\tEmail username: %v"+
				"\n\tEmail username enabled: %v \n\tSend msg count: %v \n\tForwarded msg count: %v \n\tProject ID: %v \n\tDomain: %v"+
				"\n\tPOP3 domain: %v \n\tEmail domain: %v \n\tEmails count: %v \n\tEmails unread count: %v \n\tLast message sent at: %v"+
				"\n\tSMTP ports: %v \n\tPOP3 ports: %v \n\tMax msg size: %v \n\tPermissions: %v",
			i.ID, i.Name, i.Username, i.Password, i.MaxSize, i.Status, i.EmailUsername, i.EmailUsernameEnabled, i.SentMessagesCount,
			i.ForwardedMessagesCount, i.ProjectID, i.Domain, i.POP3Domain, i.EmailDomain, i.EmailsCount, i.EmailsUnreadCount, i.LastMessageSentAt,
			i.SMTPPorts, i.POP3Ports, i.MaxMessageSize, permissions,
		)
		res = append(res, inbx)
	}
	return res
}

func preparePermissions(p mailtrap.Permissions) string {
	return fmt.Sprintf(
		"\tCan read: %v \n\tCan update: %v \n\tCan destroy: %v \n\tCan leave: %v",
		p.CanRead, p.CanUpdate, p.CanDestroy, p.CanLeave,
	)
}
