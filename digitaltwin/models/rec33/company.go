package rec33

import (
	"azure-adt-example/digitaltwin/models"
)

type Company struct {
	models.GenericModel
	Logo           string `json:"logo"`
	Name           string `json:"name"`
	PrimaryCssFile string `json:"primaryCssFile"`
	ThemeCssFile   string `json:"themeCssFile"`
	Motto          string `json:"motto"`
}

func (Company) Model() string {
	return "dtmi:digitaltwins:rec_3_3:agents:Company;1"
}
