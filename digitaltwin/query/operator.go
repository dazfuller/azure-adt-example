package query

import "log"

type Operator int

const (
	Equals Operator = iota + 1
	NotEquals
	LessThan
	GreaterThan
	LessThanOrEqual
	GreaterThanOrEqual
)

func (o Operator) String() string {
	operators := []string{"=", "!=", "<", ">", "<=", ">="}
	if !o.IsValid() {
		log.Fatalf("%d is not a valid operator type", o)
	}
	return operators[o-1]
}

func (o Operator) IsValid() bool {
	switch o {
	case Equals, NotEquals, LessThan, GreaterThan, LessThanOrEqual, GreaterThanOrEqual:
		return true
	}
	return false
}
