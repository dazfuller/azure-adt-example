package rec33

import (
	"azure-adt-example/digitaltwin/models"
)

type Company struct {
	models.GenericModel
	Logo string `json:"logo"`
	Name string `json:"name"`
}

func (Company) Model() string {
	return "dtmi:digitaltwins:rec_3_3:agents:Company;1"
}

func (Company) Alias() string {
	return models.GetModelAlias[Company]()
}
