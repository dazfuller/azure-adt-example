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
		t.Logf("Expected source to be of type rec33.Company, but got %T", cond.source)
		t.Fail()
	} else if cond.property != "Name" {
		t.Logf("Expected property to be 'Name', but got %s", cond.property)
		t.Fail()
	} else if cond.propertyJsonName != "name" {
		t.Logf("Expected JSON property name to be 'name', but got %s", cond.propertyJsonName)
		t.Fail()
	} else if cond.operator != Equals {
		t.Logf("Expected operator to be '==', but got %s", cond.operator)
		t.Fail()
	} else if len(cond.value) != 1 || cond.value[0] != "Test" {
		t.Logf("Expected value to have length of 1 and value 'Test' but got length %d and value %s", len(cond.value), cond.value[0])
		t.Fail()
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
		t.Logf("Expected error to contain '%s', but got: %v", expectedErrorString, err)
		t.Fail()
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
		t.Logf("Expected error to be '%s', but got: %v", expectedErrorString, err)
		t.Fail()
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
		t.Logf("Expected error to be '%s', but got: %v", expectedErrorString, err)
		t.Fail()
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
		t.Logf("Expected error to be '%s', but got: %v", expectedErrorString, err)
		t.Fail()
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
		t.Logf("Expected type of rec33.Company, but got %T", cond.GetSource())
		t.Fail()
	}
}

func TestWhereCondition_typeToString_Int(t *testing.T) {
	result := typeToString(73)

	if result != "73" {
		t.Logf("Expected 73, but got %s", result)
		t.Fail()
	}
}

func TestWhereCondition_typeToString_Float32(t *testing.T) {
	result := typeToString(float32(123.32))

	if result != "123.320000" {
		t.Logf("Expected 123.320000, but got %s", result)
		t.Fail()
	}
}

func TestWhereCondition_typeToString_Float64(t *testing.T) {
	result := typeToString(123.32)

	if result != "123.32000000" {
		t.Logf("Expected 123.32000000, but got %s", result)
		t.Fail()
	}
}

func TestWhereCondition_typeToString_Bool(t *testing.T) {
	result := typeToString(false)

	if result != "false" {
		t.Logf("Expected false, but got %s", result)
		t.Fail()
	}
}
