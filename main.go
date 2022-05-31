package main

import (
	"azure-adt-example/azuread"
	"azure-adt-example/digitaltwin"
	"azure-adt-example/digitaltwin/models/rec33"
	"fmt"
	"log"
)

func main() {
	config := azuread.NewTwinConfiguration()
	auth, err := azuread.GetBearerToken(config)
	if err != nil {
		log.Fatal(err)
	}
	client := digitaltwin.NewClient(config, auth)
	results, err := digitaltwin.GetTwinsOfType[rec33.Company](client)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range results {
		fmt.Printf("%s - %s\n", v.Name, v.TwinModelType())
	}
}
