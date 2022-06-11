package query

import (
	"azure-adt-example/digitaltwin/models/rec33"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewWhereFunction_BooleanExpression_SimpleGeneration(t *testing.T) {
	f, err := NewWhereFunction(rec33.Company{}, "ExternalId", IsDefined, nil)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if reflect.TypeOf(f.source) != reflect.TypeOf(rec33.Company{}) {
		t.Errorf("Expected source to be of type rec33.Company, but got %T", f.source)
	} else if f.property != "ExternalId" {
		t.Errorf("Expected property to be 'ExternalId', but got %s", f.property)
	} else if f.propertyJsonName != "$dtId" {
		t.Errorf("Expected JSON property name to be '$dtId', but got %s", f.propertyJsonName)
	} else if f.function != IsDefined {
		t.Errorf("Expected operator to be 'IS_DEFINED', but got %s", f.function)
	} else if f.value != nil {
		t.Errorf("Expected value to be nil, but got %v", f.value)
	}
}

func TestNewWhereFunction_ErrorConditions(t *testing.T) {
	tests := []struct {
		name              string
		property          string
		function          BooleanExpressionFunction
		expectedSubstring string
	}{
		{"InvalidPropertyName", "Invalid", IsNumber, "field Invalid does not exist"},
		{"InvalidFunction", "ExternalId", 17, "function specified is not valid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewWhereFunction[BooleanExpressionFunction](rec33.Company{}, test.property, test.function, nil)
			if err == nil {
				t.Error("Expected an error but got nil")
			} else if !strings.Contains(err.Error(), test.expectedSubstring) {
				t.Errorf("Error did not contain substring '%s': %v", test.expectedSubstring, err)
			}
		})
	}
}

func TestWhereFunction_GetSource(t *testing.T) {
	f, err := NewWhereFunction(rec33.Company{}, "ExternalId", IsDefined, nil)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	expected := reflect.TypeOf(rec33.Company{})
	actual := reflect.TypeOf(f.GetSource())

	if actual != expected {
		t.Errorf("Expected type of rec33.Company but got %T", f.GetSource())
	}
}

func TestModelValidationClause(t *testing.T) {
	model := rec33.Company{}
	expected := WhereFunction[BooleanExpressionFunction]{
		source:           model,
		property:         "ExternalId",
		propertyJsonName: "$dtId",
		function:         IsOfModel,
		value:            false,
	}

	f, err := ModelValidationClause(model, false)
	if err != nil {
		t.Logf("Expected error to be nil, but got %v", err)
		t.FailNow()
	}

	if reflect.TypeOf(f.source) != reflect.TypeOf(expected.source) ||
		f.property != expected.property ||
		f.propertyJsonName != expected.propertyJsonName ||
		f.function != expected.function ||
		f.value != expected.value {
		t.Errorf("Generated model %v did not match expected model %v", f, expected)
	}
}

func TestWhereFunction_BooleanExpression_GenerateClause(t *testing.T) {
	company := rec33.Company{}
	tests := []struct {
		name     string
		function BooleanExpressionFunction
		property string
		value    any
		expected string
	}{
		{"IsOfModelExact", IsOfModel, "ExternalId", true, fmt.Sprintf("IS_OF_MODEL(company, '%s', exact)", company.Model())},
		{"IsOfModelNonExact", IsOfModel, "ExternalId", false, fmt.Sprintf("IS_OF_MODEL(company, '%s')", company.Model())},
		{"IsOfModelNonBooleanValue", IsOfModel, "ExternalId", 99, fmt.Sprintf("IS_OF_MODEL(company, '%s')", company.Model())},
		{"IsBoolean", IsBool, "Name", nil, "IS_BOOL(company.name)"},
		{"IsDefined", IsDefined, "Name", nil, "IS_DEFINED(company.name)"},
		{"IsNull", IsNull, "Name", nil, "IS_NULL(company.name)"},
		{"IsNumber", IsNumber, "Name", nil, "IS_NUMBER(company.name)"},
		{"IsObject", IsObject, "Name", nil, "IS_OBJECT(company.name)"},
		{"IsPrimitive", IsPrimitive, "Name", nil, "IS_PRIMITIVE(company.name)"},
		{"IsString", IsString, "Name", nil, "IS_STRING(company.name)"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := NewWhereFunction(company, test.property, test.function, test.value)
			if err != nil {
				t.Logf("Expect nil error, but got %v", err)
				t.FailNow()
			}

			actual := f.GenerateClause()
			if actual != test.expected {
				t.Errorf("Expected \"%s\", but got \"%s\"", test.expected, actual)
			}
		})
	}
}

func TestWhereFunction_StringFunction_GenerateClause(t *testing.T) {
	company := rec33.Company{}
	tests := []struct {
		name     string
		function StringFunction
		property string
		value    any
		expected string
	}{
		{"Contains", Contains, "Name", "something", "CONTAINS(company.name, 'something')"},
		{"ContainsNilValue", Contains, "Name", nil, "CONTAINS(company.name, '')"},
		{"ContainsNonStringValue", Contains, "Name", 123, "CONTAINS(company.name, '123')"},
		{"EndsWith", EndsWith, "Name", "end", "ENDSWITH(company.name, 'end')"},
		{"EndsWithNilValue", EndsWith, "Name", nil, "ENDSWITH(company.name, '')"},
		{"StartsWith", StartsWith, "Name", "start", "STARTSWITH(company.name, 'start')"},
		{"StartsWithNilValue", StartsWith, "Name", nil, "STARTSWITH(company.name, '')"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := NewWhereFunction(company, test.property, test.function, test.value)
			if err != nil {
				t.Logf("Expect nil error, but got %v", err)
				t.FailNow()
			}

			actual := f.GenerateClause()
			if actual != test.expected {
				t.Errorf("Expected \"%s\", but got \"%s\"", test.expected, actual)
			}
		})
	}
}
