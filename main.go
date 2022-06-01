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

	from := rec33.Company{}
	builder := digitaltwin.NewBuilder(from, true)

	err = builder.AddJoin(from, rec33.Building{}, "owns", true)
	if err != nil {
		log.Fatal(err)
	}

	err = builder.WhereId(from, "NHHG")
	if err != nil {
		log.Fatal(err)
	}

	query, _ := builder.CreateQuery()
	fmt.Printf("Generated query:\n%s\n", query)

	results, err := digitaltwin.GetTwinsOfType[rec33.Building](client)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range results {
		fmt.Printf("%s - %s\n", v.Name, v.TwinModelType())
	}
}
