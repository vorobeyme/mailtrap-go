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
		fmt.Print("Enter method [list, remove] or `q` for exit: ")
		fmt.Scanf("%s", &method)
		if strings.ToLower(method) == "q" {
			return
		}

		switch method {
		// Get list of all account users.
		case "list":
			var accountID int
			fmt.Print("Enter an accountID value: ")
			fmt.Scanf("%d %d", &accountID)
			users, _, err := client.AccountUsers.List(accountID, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, u := range users {
				fmt.Printf(
					"ID: %v \nSpecifier type: %v \nSpecifier ID: %v \nSpecifier name: %v \nSpecifier email: %v"+
						"\nCan read %v \nCan update: %v \nCan destroy: %v \nCan leave: %v \nResources: %+v",
					u.ID, u.SpecifierType, u.Specifier.ID, u.Specifier.Name, u.Specifier.Email, u.Permissions.CanRead,
					u.Permissions.CanUpdate, u.Permissions.CanDestroy, u.Permissions.CanLeave, u.Resources,
				)
			}

		// Remove user by their ID.
		case "remove":
			var accountID, accountAccessID int
			fmt.Print("Enter an accountID and accountAccessID space-separated values: ")
			fmt.Scanf("%d %d", &accountID, &accountAccessID)
			_, err := client.AccountUsers.Delete(accountID, accountAccessID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("User removed.")

		default:
			fmt.Println("method doesn't exist...")
		}
	}
}
