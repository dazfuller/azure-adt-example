package query

import (
	"azure-adt-example/digitaltwin/models"
	"strings"
	"testing"
	"time"
)

type TestModel struct {
	models.GenericModel
	ExampleField           string `json:"example_field"`
	InvalidField           string `json:"example_field_2"`
	NoMapping              int
	ExtraMappingProperties time.Time `json:"ex_mapping,omitempty"`
}

func (tm TestModel) Model() string {
	return "dtmi:digitaltwins:rec_3_3:core:TestModel;1"
}

func (tm TestModel) Alias() string {
	return "testmodel"
}

func TestGetPropertyName_WithValidProperty(t *testing.T) {
	jsonName, err := getPropertyJsonName(TestModel{}, "ExampleField")
	if err != nil {
		t.Logf("Error should be nil, got %s", err)
		t.FailNow()
	}

	expectedName := "example_field"

	if jsonName != expectedName {
		t.Errorf("Expected '%s' but got '%s'", expectedName, jsonName)
	}
}

func TestGetPropertyName_WithValidPropertyWithOmitSet(t *testing.T) {
	jsonName, err := getPropertyJsonName(TestModel{}, "ExtraMappingProperties")
	if err != nil {
		t.Logf("Error should be nil, got %s", err)
		t.FailNow()
	}

	expectedName := "ex_mapping"

	if jsonName != expectedName {
		t.Errorf("Expected '%s' but got '%s'", expectedName, jsonName)
	}
}

func TestGetPropertyName_WithInvalidPropertyName(t *testing.T) {
	_, err := getPropertyJsonName(TestModel{}, "NotValid")
	if err == nil {
		t.Log("Expected to get an error, but got nil")
		t.FailNow()
	}

	expectedErrorString := "field NotValid does not exist"

	if !strings.Contains(err.Error(), expectedErrorString) {
		t.Errorf("Expected error to contain '%s', but got: %v", expectedErrorString, err)
	}
}

func TestGetPropertyName_PropertyWithMissingJsonMapping(t *testing.T) {
	_, err := getPropertyJsonName(TestModel{}, "NoMapping")
	if err == nil {
		t.Log("Expected to get an error, but got nil")
		t.FailNow()
	}

	expectedErrorString := "field query.TestModel.NoMapping does not have a json mapping property"

	if !strings.Contains(err.Error(), expectedErrorString) {
		t.Errorf("Expected error to contain '%s', but got: %v", expectedErrorString, err)
	}
}
