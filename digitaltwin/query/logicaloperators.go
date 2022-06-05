package query

import "log"

type LogicalOperator int

const (
	And LogicalOperator = iota + 1
	Or
	Not
)

func (lo LogicalOperator) String() string {
	operators := []string{"AND", "OR", "NOT"}
	if !lo.IsValid() {
		log.Fatalf("%d is not a valid logical operator type", lo)
	}
	return operators[lo-1]
}

func (lo LogicalOperator) IsValid() bool {
	switch lo {
	case And, Or, Not:
		return true
	}
	return false
}
