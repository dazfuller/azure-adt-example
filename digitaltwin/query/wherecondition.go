package query

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"strconv"
	"strings"
)

type WhereCondition struct {
	source           models.IModel
	property         string
	propertyJsonName string
	operator         Operator
	value            []any
}

func NewWhereCondition(source models.IModel, property string, operator Operator, value ...any) (*WhereCondition, error) {
	jsonPropertyName, err := getPropertyJsonName(source, property)
	if err != nil {
		return nil, err
	}

	if !operator.IsValid() {
		return nil, fmt.Errorf("operator specified is not valid")
	}

	if len(value) == 0 || (len(value) == 1 && value[0] == nil) {
		return nil, fmt.Errorf("at least one value must be provided")
	}

	if (operator == In || operator == NotIn) && len(value) > 100 {
		return nil, fmt.Errorf("IN and NIN operators do not support more than 100 values as part of the query")
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

	switch wc.operator {
	case In, NotIn:
		valueCollection := make([]string, len(wc.value))
		for i, v := range wc.value {
			valueCollection[i] = typeToString(v)
		}
		value = fmt.Sprintf("[%s]", strings.Join(valueCollection, ", "))
	default:
		value = typeToString(wc.value[0])
	}

	return fmt.Sprintf("%s.%s %s %s", wc.source.Alias(), wc.propertyJsonName, wc.operator, value)
}

func (wc *WhereCondition) GetSource() models.IModel {
	return wc.source
}

func typeToString(value any) string {
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
