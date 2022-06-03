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
	builder := digitaltwin.NewBuilder(from, false)
	var err error
	if err = builder.AddJoin(from, rec33.Building{}, "owns", false); err != nil {
		log.Fatal(err)
	}
	if err = builder.AddJoin(rec33.Building{}, rec33.Level{}, "isPartOf", true); err != nil {
		log.Fatal(err)
	}
	if err = builder.WhereId(from, "NHHG"); err != nil {
		log.Fatal(err)
	}

	results, err := digitaltwin.ExecuteBuilder3[rec33.Company, rec33.Building, rec33.Level](client, builder)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range results {
		//fmt.Printf("%s owns %s (%s)\n", row.Twin1.Name, row.Twin2.Name, row.Twin2.ExternalId)
		fmt.Printf("%s is part of %s owned by %s\n", row.Twin3.Name, row.Twin2.Name, row.Twin1.Name)
	}
}
