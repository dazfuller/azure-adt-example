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

	results, err := digitaltwin.ExecuteBuilder[rec33.Company, rec33.Building](client, builder)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range results {
		fmt.Printf("%s owns %s\n", row[0].(rec33.Company).Name, row[1].(rec33.Building).Name)
	}
}
