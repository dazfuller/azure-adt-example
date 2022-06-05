package query

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"strconv"
)

type WhereCondition struct {
	source           models.IModel
	property         string
	propertyJsonName string
	operator         Operator
	value            any
}

func NewWhereCondition(source models.IModel, property string, operator Operator, value any) (*WhereCondition, error) {
	jsonPropertyName, err := getPropertyJsonName(source, property)
	if err != nil {
		return nil, err
	}

	if !operator.IsValid() {
		return nil, fmt.Errorf("operator specified is not valid")
	}

	return &WhereCondition{
		source:           source,
		property:         property,
		propertyJsonName: jsonPropertyName,
		operator:         operator,
		value:            value,
	}, nil
}

func (wc *WhereCondition) GenerateClause() string {
	var value string

	switch wc.value.(type) {
	case int, float32, float64:
		value = fmt.Sprintf("%d", wc.value)
	case bool:
		value = strconv.FormatBool(wc.value.(bool))
	default:
		value = fmt.Sprintf("'%s'", wc.value)
	}

	return fmt.Sprintf("%s.%s %s %s", wc.source.Alias(), wc.propertyJsonName, wc.operator, value)
}
