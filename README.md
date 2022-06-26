# Azure Digital Twin Query API

[![codecov](https://codecov.io/gh/dazfuller/azure-adt-example/branch/main/graph/badge.svg?token=GRK06HOQ61)](https://codecov.io/gh/dazfuller/azure-adt-example)

This is a sample project I created to learn [Go](https://go.dev/). There's an existing SDK for working with
[Azure Digital Twin](https://azure.microsoft.com/services/digital-twins/) but only for the management plain,
so I thought I'd try and write a small application which uses the [Query API](https://docs.microsoft.com/rest/api/digital-twins/dataplane/query/querytwins),
but also handle the authentication flow. I've previously written a Fluent API for the Digital Twin Query API
in [C#](https://dotnet.microsoft.com) and wanted to try something like it for Go using [Generics](https://go.dev/doc/tutorial/generics).
That's what this has become.

The authentication flow needs to read some values from environment variables. You can do this by either
creating the environment variables, or by creating a `.env` file.

```text
TWIN_URL=https://<twin instance>.<region>.digitaltwins.azure.net
TWIN_CLIENT_ID=<application id>
TWIN_CLIENT_SECRET=<client secret>
TWIN_TENANT_ID=<directory id>
```

## Models / Ontology

The builder process uses models which implement the `models.IModel` interface. Each twin in Azure Digital Twin
has a `$dtId`, `$etag`, and `$metadata` value, so there is a `GenericModel` type which can be inherited so that
these items are handled for you. Each model needs to specify its model type so that its type can be validated
if needed. An example model would look as follows.

```go
package myontology

import (
	"azure-adt-example/digitaltwin/models"
)

type MyOntologyType struct {
	models.GenericModel
	Name           string `json:"name"`
}

func (MyOntologyType) Model() string {
	return "dtmi:digitaltwins:rec_3_3:agents:MyOntologyType;1"
}

func (MyOntologyType) Alias() string {
	return models.GetModelAlias[MyOntologyType]()
}

func (MyOntologyType) ValidationClause() string {
	return models.ModelValidationClause[MyOntologyType]()
}
```

The `Alias` and `ValidationClause` can be manually specified, but there are helper methods to generate the
correct values. The `GetModelAlias` returns the name of the type in lowercase.

## Usage

Once the models for the ontology have been specified the twin can be queried as follows.

```go
config := azuread.NewTwinConfiguration()
// An authentication token can be provided, but if it isn't then the client will get it's own
client := digitaltwin.NewClient(config, nil)

from := rec33.Company{}

// Create a new builder using Company as the base twin type
builder := digitaltwin.NewBuilder(from, false)

var err error

// Add a join from Company to Building where the company "owns" the building
if err = builder.AddJoin(from, rec33.Building{}, "owns", false); err != nil {
    log.Fatal(err)
}

// Add a join to Level where it is part of the building
if err = builder.AddJoin(rec33.Building{}, rec33.Level{}, "isPartOf", true); err != nil {
    log.Fatal(err)
}

// Add a where clause for the query
if err = builder.WhereId(from, "<company id>"); err != nil {
    log.Fatal(err)
}

// Execute the generated query and return the company, building, and level objects from the results.
// There are also methods for a single return type and 2 return types
results, err := digitaltwin.ExecuteBuilder3[rec33.Company, rec33.Building, rec33.Level](client, builder)
if err != nil {
    log.Fatal(err)
}

// Output the names of the company, building, and level. Because of Generics the `TwinX` fields are
// typed and so access to all of the properties is available
for _, row := range results {
    fmt.Printf("%s is part of %s owned by %s\n", row.Twin3.Name, row.Twin2.Name, row.Twin1.Name)
}
```

## Issues

This is a side project for teaching myself, but I'm putting it out there in case anyone else finds it
useful.

There's still a number of things I want to implement.

* Extend return types up to 5 types (5 is the maximum number of relationships in a query unless using `MATCH` which is still in preview)
* Clean the interface up a bit to make it more obvious
* Add where clause building which includes all the available functions from Digital Twin
* UNIT TESTING
