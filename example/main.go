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
	}
}
