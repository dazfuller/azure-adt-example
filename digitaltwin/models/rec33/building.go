package rec33

import "azure-adt-example/digitaltwin/models"

type Building struct {
	models.GenericModel
	Name string `json:"name"`
}

func (Building) Model() string {
	return "dtmi:digitaltwins:rec_3_3:core:Building;1"
}

func (Building) Alias() string {
	return models.GetModelAlias[Building]()
}

func (Building) ValidationClause() string {
	return models.ModelValidationClause[Building]()
}
