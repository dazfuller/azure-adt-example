package query

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"reflect"
	"strconv"
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

func typeToString(value any) string {
	if value == nil {
		return "''"
	}
	switch value.(type) {
	case int:
		return fmt.Sprintf("%d", value)
	case float32:
		return fmt.Sprintf("%f", value)
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', 8, 64)
	case bool:
		return strconv.FormatBool(value.(bool))
	default:
		return fmt.Sprintf("'%s'", value)
	}
}
