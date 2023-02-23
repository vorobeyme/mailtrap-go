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
		fmt.Print("Enter method [permissions, resources] or `q` for exit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		// Manage user or token permissions.
		case "permissions":
			var accountAccessID, accountID int
			fmt.Print("Enter an accountAccessID and accountID space-separated values: ")
			fmt.Scanf("%d %d", &accountAccessID, &accountID)
			_, err := client.Permissions.Manage(accountID, accountAccessID, &[]mailtrap.PermissionRequest{})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Permissions have been updated.")

		// Get all resources in account to which the token has admin access.
		case "resources":
			var accountID int
			fmt.Print("Enter an accountID value: ")
			fmt.Scanf("%d %d", &accountID)
			resources, _, err := client.Permissions.ListResources(accountID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n\n", displayResource(resources, 0))

		default:
			fmt.Println("method doesn't exist...")
		}
	}
}

func displayResource(resources []*mailtrap.Resource, tabs int) string {
	var res string
	// var tabStr = strings.Repeat("", tabs)
	for _, r := range resources {
		res = fmt.Sprintf(
			"\nID: %v \nName: %v \nType: %v \nAccess level: %v \nResources: %v",
			r.ID, r.Name, r.Type, r.AccessLevel, displayResource(r.Resource, 2),
		)
	}
	return res
}
