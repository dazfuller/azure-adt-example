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
	In
	NotIn
)

func (o Operator) String() string {
	operators := []string{"=", "!=", "<", ">", "<=", ">=", "IN", "NIN"}
	if !o.IsValid() {
		log.Fatalf("%d is not a valid operator type", o)
	}
	return operators[o-1]
}

func (o Operator) ToName() string {
	operators := []string{"Equals", "NotEquals", "LessThan", "GreaterThan", "LessThanOrEqual", "GreaterThanOrEqual", "In", "NotIn"}
	if !o.IsValid() {
		log.Fatalf("%d is not a valid operator type", o)
	}
	return operators[o-1]
}

func (o Operator) IsValid() bool {
	switch o {
	case Equals, NotEquals, LessThan, GreaterThan, LessThanOrEqual, GreaterThanOrEqual, In, NotIn:
		return true
	}
	return false
}
