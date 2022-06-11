package query

import (
	"azure-adt-example/digitaltwin/models/rec33"
	"strings"
	"testing"
)

func TestNewWhereLogical(t *testing.T) {
	subClause, err := NewWhereCondition(rec33.Company{}, "ExternalId", Equals, "Comp1")
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}
	clause, err := NewWhereLogical(And, subClause)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if clause.operator != And || len(clause.conditions) != 1 || clause.conditions[0] != subClause {
		t.Errorf("Invalid clause generated")
	}
}

func TestNewWhereLogical_ErrorConditions(t *testing.T) {
	tests := []struct {
		name       string
		operator   LogicalOperator
		subClauses []IWhere
		expected   string
	}{
		{"NoSubClauses", And, make([]IWhere, 0), "at least one condition must be specified"},
		{"TooManyClausesForNot", Not, make([]IWhere, 2), "NOT only supports a single condition"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewWhereLogical(test.operator, test.subClauses...)
			if err == nil {
				t.Logf("Expected an error but got nil")
				t.FailNow()
			}

			if !strings.Contains(err.Error(), test.expected) {
				t.Errorf("Expected error to contain \"%s\", but got: %v", test.expected, err)
			}
		})
	}
}

func TestWhereLogical_GetSource(t *testing.T) {
	clause := WhereLogical{}
	if clause.GetSource() != nil {
		t.Errorf("Expected nil, but got %v", clause.GetSource())
	}
}

func TestWhereLogical_GenerateClause(t *testing.T) {
	subClause1, _ := NewWhereCondition(rec33.Company{}, "ExternalId", Equals, "Comp1")
	subClause2, _ := NewWhereCondition(rec33.Company{}, "ExternalId", Equals, "Comp2")
	subClause3, _ := NewWhereLogical(Not, subClause2)
	subClause4, _ := NewWhereLogical(Or, subClause1, subClause2)
	tests := []struct {
		name       string
		operator   LogicalOperator
		subClauses []IWhere
		expected   string
	}{
		{"Not", Not, []IWhere{subClause1}, "(NOT company.$dtId = 'Comp1')"},
		{"And", And, []IWhere{subClause1, subClause2}, "(company.$dtId = 'Comp1' AND company.$dtId = 'Comp2')"},
		{"Or", Or, []IWhere{subClause1, subClause2}, "(company.$dtId = 'Comp1' OR company.$dtId = 'Comp2')"},
		{"OrNestedNot", Or, []IWhere{subClause1, subClause3}, "(company.$dtId = 'Comp1' OR (NOT company.$dtId = 'Comp2'))"},
		{"AndNestedOr", And, []IWhere{subClause1, subClause4}, "(company.$dtId = 'Comp1' AND (company.$dtId = 'Comp1' OR company.$dtId = 'Comp2'))"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cond, err := NewWhereLogical(test.operator, test.subClauses...)
			if err != nil {
				t.Logf("Expected nil error, but got %v", err)
				t.FailNow()
			}

			actual := cond.GenerateClause()

			if actual != test.expected {
				t.Errorf("Expected \"%s\", but got \"%s\"", test.expected, actual)
			}
		})
	}
}
