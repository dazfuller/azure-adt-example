package query

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"log"
	"reflect"
)

type WhereFunction[F Function] struct {
	source           models.IModel
	property         string
	propertyJsonName string
	function         F
	value            any
}

func NewWhereFunction[F Function](source models.IModel, property string, function F, value any) (*WhereFunction[F], error) {
	jsonPropertyName, err := getPropertyJsonName(source, property)
	if err != nil {
		return nil, err
	}

	if !function.IsValid() {
		return nil, fmt.Errorf("function %v specified is not valid", function)
	}

	return &WhereFunction[F]{
		source:           source,
		property:         property,
		propertyJsonName: jsonPropertyName,
		function:         function,
		value:            value,
	}, nil
}

func (wf *WhereFunction[F]) GenerateClause() string {
	var expression string
	v1 := reflect.ValueOf(wf.function)
	v := v1.Type().Name()
	sfName := reflect.TypeOf(*new(StringFunction)).Name()
	beName := reflect.TypeOf(*new(BooleanExpressionFunction)).Name()
	switch v {
	case sfName:
		expression = fmt.Sprintf("%s(%s.%s, '%s')", wf.function, wf.source.Alias(), wf.propertyJsonName, wf.value)
	case beName:
		switch x := BooleanExpressionFunction(v1.Int()); x {
		case IsOfModel:
			boolValue, ok := wf.value.(bool)
			if !ok {
				boolValue = false
			}
			var exactMatch string
			if boolValue {
				exactMatch = ", exact"
			} else {
				exactMatch = ""
			}
			expression = fmt.Sprintf("%s(%s, '%s'%s)", wf.function, wf.source.Alias(), wf.source.Model(), exactMatch)
		default:
			expression = fmt.Sprintf("%s(%s.%s)", wf.function, wf.source.Alias(), wf.propertyJsonName)
		}
	}
	log.Print(v)

	return expression
}

func ModelValidationClause(source models.IModel) *WhereFunction[BooleanExpressionFunction] {
	wf, err := NewWhereFunction(source, "ExternalId", IsOfModel, false)
	if err != nil {
		log.Fatal(err)
	}
	return wf
}
