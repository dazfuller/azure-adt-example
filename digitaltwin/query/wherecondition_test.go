package query

import (
	"azure-adt-example/digitaltwin/models/rec33"
	"reflect"
	"strings"
	"testing"
)

func TestNewWhereCondition_SimpleCreation(t *testing.T) {
	cond, err := NewWhereCondition(rec33.Company{}, "Name", Equals, "Test")
	if err != nil {
		t.Logf("Expected nil for error, but got %v", err)
		t.FailNow()
	}

	if reflect.TypeOf(cond.source) != reflect.TypeOf(rec33.Company{}) {
		t.Errorf("Expected source to be of type rec33.Company, but got %T", cond.source)
	} else if cond.property != "Name" {
		t.Errorf("Expected property to be 'Name', but got %s", cond.property)
	} else if cond.propertyJsonName != "name" {
		t.Errorf("Expected JSON property name to be 'name', but got %s", cond.propertyJsonName)
	} else if cond.operator != Equals {
		t.Errorf("Expected operator to be '==', but got %s", cond.operator)
	} else if len(cond.value) != 1 || cond.value[0] != "Test" {
		t.Errorf("Expected value to have length of 1 and value 'Test' but got length %d and value %s", len(cond.value), cond.value[0])
	}
}

func TestNewWhereCondition_InvalidPropertyName(t *testing.T) {
	_, err := NewWhereCondition(rec33.Company{}, "Invalid", Equals, "Test")
	if err == nil {
		t.Log("Expected an error but got nil")
		t.FailNow()
	}

	expectedErrorString := "field Invalid does not exist"

	if !strings.Contains(err.Error(), expectedErrorString) {
		t.Errorf("Expected error to contain '%s', but got: %v", expectedErrorString, err)
	}
}

func TestNewWhereCondition_WithInvalidOperator(t *testing.T) {
	_, err := NewWhereCondition(rec33.Company{}, "Name", 99, "Test")
	if err == nil {
		t.Log("Expected an error but got nil")
		t.FailNow()
	}

	expectedErrorString := "operator specified is not valid"

	if err.Error() != expectedErrorString {
		t.Errorf("Expected error to be '%s', but got: %v", expectedErrorString, err)
	}
}

func TestNewWhereCondition_NoValues(t *testing.T) {
	_, err := NewWhereCondition(rec33.Company{}, "Name", Equals, make([]any, 0)...)
	if err == nil {
		t.Log("Expected an error but got nil")
		t.FailNow()
	}

	expectedErrorString := "at least one value must be provided"

	if err.Error() != expectedErrorString {
		t.Errorf("Expected error to be '%s', but got: %v", expectedErrorString, err)
	}
}

func TestNewWhereCondition_NilValue(t *testing.T) {
	_, err := NewWhereCondition(rec33.Company{}, "Name", Equals, nil)
	if err == nil {
		t.Log("Expected an error but got nil")
		t.FailNow()
	}

	expectedErrorString := "at least one value must be provided"

	if err.Error() != expectedErrorString {
		t.Errorf("Expected error to be '%s', but got: %v", expectedErrorString, err)
	}
}

func TestNewWhereCondition_InvalidLength(t *testing.T) {
	tests := []struct {
		operator Operator
	}{
		{In},
		{NotIn},
	}

	for _, test := range tests {
		_, err := NewWhereCondition(rec33.Company{}, "Name", test.operator, make([]any, 101)...)
		if err == nil {
			t.Errorf("Expected an error for %s but got nil", test.operator)
		} else if err.Error() != "IN and NIN operators do not support more than 100 values as part of the query" {
			t.Errorf("Invalid error recieved, got %s", err)
		}
	}
}

func TestWhereCondition_GetSource(t *testing.T) {
	cond, err := NewWhereCondition(rec33.Company{}, "Name", Equals, "Test")
	if err != nil {
		t.Logf("Expected nil for error, but got %v", err)
		t.FailNow()
	}

	expected := reflect.TypeOf(rec33.Company{})

	if reflect.TypeOf(cond.GetSource()) != expected {
		t.Errorf("Expected type of rec33.Company, but got %T", cond.GetSource())
	}
}

func TestWhereCondition_typeToString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"Int", 73, "73"},
		{"Float32", float32(123.32), "123.320000"},
		{"Float64", 123.32, "123.32000000"},
		{"Boolean", false, "false"},
		{"String", "testing", "'testing'"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := typeToString(test.input)
			if actual != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, actual)
			}
		})
	}
}

func TestWhereCondition_GenerateClause(t *testing.T) {
	tests := []struct {
		field    string
		operator Operator
		value    []any
		expected string
	}{
		{"ExternalId", Equals, []any{"Test"}, "company.$dtId = 'Test'"},
		{"ExternalId", NotEquals, []any{"Test"}, "company.$dtId != 'Test'"},
		{"Name", LessThan, []any{12}, "company.name < 12"},
		{"Name", GreaterThan, []any{99}, "company.name > 99"},
		{"Name", LessThanOrEqual, []any{12}, "company.name <= 12"},
		{"Name", GreaterThanOrEqual, []any{99}, "company.name >= 99"},
		{"ExternalId", In, []any{"Comp1", "Comp2"}, "company.$dtId IN ['Comp1', 'Comp2']"},
		{"ExternalId", NotIn, []any{"Comp1", "Comp2"}, "company.$dtId NIN ['Comp1', 'Comp2']"},
	}

	for _, test := range tests {
		t.Run(test.operator.ToName(), func(t *testing.T) {
			cond, _ := NewWhereCondition(rec33.Company{}, test.field, test.operator, test.value...)

			actualClause := cond.GenerateClause()

			if test.expected != actualClause {
				t.Errorf("Expected \"%s\", but got \"%s\"", test.expected, actualClause)
			}
		})
	}
}
