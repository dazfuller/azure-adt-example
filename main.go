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
	client := digitaltwin.NewClient(config, nil)

	from := rec33.Company{}
	builder := digitaltwin.NewBuilder(from, true)
	var err error
	if err = builder.AddJoin(from, rec33.Building{}, "owns", false); err != nil {
		log.Fatal(err)
	}
	if err = builder.WhereId(from, "NHHG"); err != nil {
		log.Fatal(err)
	}

	results, err := digitaltwin.ExecuteBuilder2[rec33.Company, rec33.Building](client, builder)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range results {
		fmt.Printf("%s owns %s (%s)\n", row.Twin1.Name, row.Twin2.Name, row.Twin2.ExternalId)
	}
}
