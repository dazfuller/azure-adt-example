package query

import (
	"azure-adt-example/digitaltwin/models"
	"azure-adt-example/digitaltwin/models/rec33"
	"reflect"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	source := rec33.Company{}

	expected := Builder{
		from:          source,
		validateFrom:  true,
		validateExact: true,
		join:          make([]join, 0),
		where:         make([]IWhere, 0),
		project:       make([]models.IModel, 0),
	}

	actual := NewBuilder(source, true, true)

	if reflect.TypeOf(actual.from) != reflect.TypeOf(expected.from) ||
		actual.validateFrom != expected.validateFrom ||
		actual.validateExact != expected.validateExact ||
		len(actual.join) != 0 ||
		len(actual.where) != 0 ||
		len(actual.project) != 0 {
		t.Errorf("Builder %v does not match expected %v", actual, expected)
	}
}

func TestBuilder_AddJoin(t *testing.T) {
	builder := NewBuilder(rec33.Company{}, false, false)

	err := builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", false, false)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if len(builder.join) != 1 {
		t.Logf("Expected single join expression, but received %v", builder.join)
		t.FailNow()
	}

	if reflect.TypeOf(builder.join[0].source) != reflect.TypeOf(rec33.Company{}) {
		t.Errorf("Expected source to be rec33.Company, but got %T", builder.join[0].source)
	}
	if reflect.TypeOf(builder.join[0].target) != reflect.TypeOf(rec33.Building{}) {
		t.Errorf("Expected target to be rec33.Building, but got %T", builder.join[0].target)
	}
	if builder.join[0].relationship != "owns" {
		t.Errorf("Expected relationship to be 'owns', but got '%s'", builder.join[0].relationship)
	}
	if builder.join[0].validateType != false {
		t.Errorf("Expected type validation to be false, but got '%t'", builder.join[0].validateType)
	}
	if builder.join[0].validateExact != false {
		t.Errorf("Expected exact validation to be false, but got '%t'", builder.join[0].validateExact)
	}
}

func TestBuilder_AddJoin_DuplicateJoin(t *testing.T) {
	builder := NewBuilder(rec33.Company{}, false, false)

	err := builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", false, false)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	err = builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", false, false)
	if err == nil {
		t.Log("Expected error but got nil")
		t.FailNow()
	}

	expectedError := "a target of alias 'building' already exists"

	if err.Error() != expectedError {
		t.Errorf("Expected error \"%s\", but got \"%s\"", expectedError, err.Error())
	}
}
