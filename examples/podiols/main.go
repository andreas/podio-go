// Command podiols lists the content of your podio account.
package main

import (
	"fmt"

	"github.com/andreas/podio-go"
)

func main() {
	authToken, err := podio.AuthWithUserCredentials("my-client-id", "my-client-secret", "my-username", "my-password")

	if err != nil {
		fmt.Println("Auth failed: ", err)
		return
	}

	client := podio.NewClient(authToken)
	orgs, err := client.GetOrganizations()

	if err != nil {
		fmt.Println("Failed to get orgs: ", err)
		return
	}

	for _, org := range orgs {
		fmt.Println("Org: ", org.Name)

		spaces, err := client.GetSpaces(org.Id)
		if err != nil {
			fmt.Println("Failed to get spaces: ", err)
			continue
		}

		for _, space := range spaces {
			fmt.Println("Space: ", space.Name)

			apps, err := client.GetApps(space.Id)
			if err != nil {
				fmt.Println("Failed to get apps: ", err)
				continue
			}

			for _, app := range apps {
				fmt.Println("App: ", app.Name)

				items, err := client.GetItems(app.Id)
				if err != nil {
					fmt.Println("Failed to get items: ", err)
					continue
				}

				for _, item := range items.Items {
					fmt.Printf("Item: %v\n", item)
				}
			}
		}
	}
}
