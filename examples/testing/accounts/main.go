// An example of how to use Account methods.
//
// It's runnable with the following command:
// export MAILTRAP_API_KEY=your_api_key
// go run .
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
	client := mailtrap.New(apiKey)

	accounts, _, err := client.Accounts.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range accounts {
		fmt.Println("ID:", v.ID, "\nName:", v.Name, "\nAccess levels:", v.AccessLevels)
	}
}
