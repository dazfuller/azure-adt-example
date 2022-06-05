package query

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"strings"
)

type WhereLogical struct {
	operator   LogicalOperator
	conditions []IWhere
}

func NewWhereLogical(operator LogicalOperator, conditions ...IWhere) (*WhereLogical, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("at least one condition must be specified")
	} else if operator == Not && len(conditions) > 1 {
		return nil, fmt.Errorf("NOT only supports a single condition")
	}

	return &WhereLogical{operator: operator, conditions: conditions}, nil
}

func (wl *WhereLogical) GenerateClause() string {
	conditions := make([]string, len(wl.conditions))

	for i, c := range wl.conditions {
		conditions[i] = c.GenerateClause()
	}

	if wl.operator == Not {
		return fmt.Sprintf("(%s %s)", wl.operator, conditions[0])
	}

	return fmt.Sprintf("(%s)", strings.Join(conditions, fmt.Sprintf(" %s ", wl.operator)))
}

func (wl *WhereLogical) GetSource() models.IModel {
	return nil
}
