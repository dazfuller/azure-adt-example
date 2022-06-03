package rec33

import "azure-adt-example/digitaltwin/models"

type Level struct {
	models.GenericModel
	Name            string `json:"name"`
	Number          int32  `json:"levelNumber"`
	PersonCapacity  int32  `json:"personCapacity"`
	PersonOccupancy int32  `json:"personOccupancy"`
}

func (Level) Model() string {
	return "dtmi:digitaltwins:rec_3_3:core:Level;1"
}

func (Level) Alias() string {
	return models.GetModelAlias[Level]()
}

func (Level) ValidationClause() string {
	return models.ModelValidationClause[Level]()
}
