package query

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"reflect"
	"strings"
)

type IWhere interface {
	GenerateClause() string
	GetSource() models.IModel
}

func getPropertyJsonName(source models.IModel, property string) (string, error) {
	rm := reflect.ValueOf(source)
	field, ok := rm.Type().FieldByName(property)

	if !ok {
		return "", fmt.Errorf("field %s does not exist on model %T", property, source)
	}

	jsonName, ok := field.Tag.Lookup("json")
	if !ok {
		return "", fmt.Errorf("field %T.%s does not have a json mapping property", source, property)
	}

	jsonPropertyName := strings.Split(jsonName, ",")[0]
	return jsonPropertyName, nil
}
