package query

import (
	"azure-adt-example/digitaltwin/models"
	"azure-adt-example/digitaltwin/models/rec33"
	"reflect"
	"strings"
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

func createErrorString(err string) *string {
	return &err
}

func TestBuilder_WhereId(t *testing.T) {
	tests := []struct {
		name     string
		source   models.IModel
		ids      []string
		expected *string
	}{
		{"ValidSingleId", rec33.Company{}, []string{"Comp1"}, nil},
		{"ValidMultipleIds", rec33.Company{}, []string{"Comp1", "Comp2"}, nil},
		{"InvalidSource", rec33.Building{}, []string{"Building1"}, createErrorString("source building is not part of the query")},
		{"NoIds", rec33.Company{}, make([]string, 0), createErrorString("at least one id must be specified")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder(rec33.Company{}, false, false)
			err := builder.WhereId(test.source, test.ids...)

			assertExpectedError(t, err, test.expected)
		})
	}
}

func TestBuilder_WhereClause(t *testing.T) {
	tests := []struct {
		name     string
		source   models.IModel
		property string
		operator Operator
		value    []any
		expected *string
	}{
		{"ValidSingleValue", rec33.Company{}, "Name", NotEquals, []any{"Test"}, nil},
		{"ValidMultipleValues", rec33.Company{}, "Name", NotEquals, []any{"Test1", "Test2"}, nil},
		{"InvalidSource", rec33.Building{}, "Name", Equals, []any{"Test1"}, createErrorString("source building is not part of the query")},
		{"InvalidValues", rec33.Company{}, "Name", GreaterThanOrEqual, []any{nil}, createErrorString("at least one value must be provided")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder(rec33.Company{}, false, false)
			err := builder.WhereClause(test.source, test.property, test.operator, test.value...)

			assertExpectedError(t, err, test.expected)
		})
	}
}

func TestBuilder_WhereStringFunction(t *testing.T) {
	tests := []struct {
		name     string
		source   models.IModel
		property string
		function StringFunction
		value    string
		expected *string
	}{
		{"ValidStringFunction", rec33.Company{}, "Name", Contains, "Test", nil},
		{"InvalidSource", rec33.Building{}, "Name", EndsWith, "est", createErrorString("source building is not part of the query")},
		{"InvalidProperty", rec33.Company{}, "DoesNotExist", StartsWith, "", createErrorString("field DoesNotExist does not exist")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder(rec33.Company{}, false, false)
			err := builder.WhereStringFunction(test.source, test.property, test.function, test.value)

			assertExpectedError(t, err, test.expected)
		})
	}
}

func TestBuilder_WhereBooleanFunction(t *testing.T) {
	tests := []struct {
		name     string
		source   models.IModel
		property string
		function BooleanExpressionFunction
		value    any
		expected *string
	}{
		{"ValidBooleanFunction", rec33.Company{}, "Name", IsString, nil, nil},
		{"InvalidSource", rec33.Building{}, "Name", IsNull, "est", createErrorString("source building is not part of the query")},
		{"InvalidProperty", rec33.Company{}, "DoesNotExist", IsObject, "", createErrorString("field DoesNotExist does not exist")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder(rec33.Company{}, false, false)
			err := builder.WhereBooleanFunction(test.source, test.property, test.function, test.value)

			assertExpectedError(t, err, test.expected)
		})
	}
}

func TestBuilder_WhereLogicalOperator(t *testing.T) {
	subClause1, _ := NewWhereCondition(rec33.Company{}, "ExternalId", Equals, "T1")
	subClause2, _ := NewWhereCondition(rec33.Company{}, "ExternalId", Equals, "T2")
	subClause3, _ := NewWhereCondition(rec33.Building{}, "ExternalId", Equals, "T2")

	tests := []struct {
		name       string
		operator   LogicalOperator
		conditions []IWhere
		expected   *string
	}{
		{"ValidLogical", Or, []IWhere{subClause1, subClause2}, nil},
		{"InvalidSource", And, []IWhere{subClause1, subClause3}, createErrorString("source building is not part of the query")},
		{"InvalidConditions", Not, []IWhere{subClause1, subClause2}, createErrorString("NOT only supports a single condition")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder(rec33.Company{}, false, false)
			err := builder.WhereLogicalOperator(test.operator, test.conditions...)

			assertExpectedError(t, err, test.expected)
		})
	}
}

func TestBuilder_AddProjection(t *testing.T) {
	tests := []struct {
		name     string
		source   models.IModel
		expected *string
	}{
		{"ValidProjection", rec33.Company{}, nil},
		{"InvalidProjection", rec33.Building{}, createErrorString("source building is not part of the query")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder(rec33.Company{}, false, false)
			err := builder.AddProjection(test.source)

			assertExpectedError(t, err, test.expected)
		})
	}
}

func TestBuilder_AddProjection_AlreadyExists(t *testing.T) {
	builder := NewBuilder(rec33.Company{}, false, false)
	_ = builder.AddProjection(rec33.Company{})
	_ = builder.AddProjection(rec33.Company{})

	if len(builder.project) > 1 {
		t.Errorf("Expected only a single projection, but %d exist", len(builder.project))
	}
}

func TestBuilder_AddProjection_FromJoin(t *testing.T) {
	builder := NewBuilder(rec33.Company{}, false, false)
	_ = builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", false, false)
	_ = builder.AddProjection(rec33.Company{})
	err := builder.AddProjection(rec33.Building{})

	if err != nil {
		t.Errorf("Adding building as projection from join clause should work, but got error %v", err)
	}
}

func TestBuilder_CreateQuery(t *testing.T) {
	builder := NewBuilder(rec33.Company{}, false, false)
	_ = builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", false, false)
	_ = builder.WhereId(rec33.Company{}, "Comp1")

	officeClause, _ := NewWhereFunction(rec33.Building{}, "ExternalId", StartsWith, "Office")
	warehouseClause, _ := NewWhereFunction(rec33.Building{}, "ExternalId", StartsWith, "Warehouse")

	_ = builder.WhereLogicalOperator(Or, officeClause, warehouseClause)
	_ = builder.AddProjection(rec33.Company{})
	_ = builder.AddProjection(rec33.Building{})

	expected := "SELECT company, building FROM digitaltwins company JOIN building RELATED company.owns WHERE company.$dtId = 'Comp1' AND (STARTSWITH(building.$dtId, 'Office') OR STARTSWITH(building.$dtId, 'Warehouse'))"

	actual, err := builder.CreateQuery()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	} else if *actual != expected {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, *actual)
	}
}

func TestBuilder_CreateQuery_DefaultProjection(t *testing.T) {
	builder := NewBuilder(rec33.Company{}, false, false)
	_ = builder.WhereId(rec33.Company{}, "Comp1")

	expected := "SELECT company FROM digitaltwins company WHERE company.$dtId = 'Comp1'"

	actual, err := builder.CreateQuery()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	} else if *actual != expected {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, *actual)
	}
}

func TestBuilder_CreateQuery_ModelValidate(t *testing.T) {
	builder := NewBuilder(rec33.Company{}, true, false)
	_ = builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", true, true)
	_ = builder.WhereId(rec33.Company{}, "Comp1")
	_ = builder.AddProjection(rec33.Company{})
	_ = builder.AddProjection(rec33.Building{})

	expected := "SELECT company, building FROM digitaltwins company JOIN building RELATED company.owns WHERE company.$dtId = 'Comp1' AND IS_OF_MODEL(company, 'dtmi:digitaltwins:rec_3_3:agents:Company;1') AND IS_OF_MODEL(building, 'dtmi:digitaltwins:rec_3_3:core:Building;1', exact)"

	actual, err := builder.CreateQuery()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	} else if *actual != expected {
		t.Errorf("Expected:\n%s\nActual:\n%s", expected, *actual)
	}
}

func assertExpectedError(t *testing.T, actual error, expected *string) {
	if expected == nil && actual != nil {
		t.Errorf("Expected nil error, but got %v", actual)
	} else if expected != nil {
		if actual == nil {
			t.Errorf("Expected error containing '%s', but got nil", *expected)
		} else if !strings.Contains(actual.Error(), *expected) {
			t.Errorf("Expected error containing '%s', but got '%s'", *expected, actual)
		}
	}
}
