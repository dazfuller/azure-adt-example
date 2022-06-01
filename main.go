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
	_ = builder.AddJoin(from, rec33.Building{}, "owns", false)

	_ = builder.AddProjection(from)
	_ = builder.AddProjection(rec33.Building{})

	err = builder.WhereId(from, "NHHG")
	if err != nil {
		log.Fatal(err)
	}

	query, _ := builder.CreateQuery()
	fmt.Printf("Generated query:\n%s\n", query)

	_, err = digitaltwin.ExecuteBuilder[rec33.Company, rec33.Building](client, builder)
	if err != nil {
		log.Fatal(err)
	}
	/*
		results, err := digitaltwin.GetTwinsOfType[rec33.Building](client)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range results {
			fmt.Printf("%s - %s\n", v.Name, v.TwinModelType())
		}
	*/
}
