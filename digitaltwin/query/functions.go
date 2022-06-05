package query

import "log"

type Function interface {
	StringFunction | BooleanExpressionFunction
	IsValid() bool
}

type StringFunction int

const (
	Contains StringFunction = iota + 1
	EndsWith
	StartsWith
)

func (sf StringFunction) String() string {
	functions := []string{"CONTAINS", "ENDSWITH", "STARTSWITH"}
	if !sf.IsValid() {
		log.Fatalf("%d is not a valid function type", sf)
	}
	return functions[sf-1]
}

func (sf StringFunction) IsValid() bool {
	switch sf {
	case Contains, EndsWith, StartsWith:
		return true
	}
	return false
}

type BooleanExpressionFunction int

const (
	IsBool BooleanExpressionFunction = iota + 1
	IsDefined
	IsNull
	IsNumber
	IsObject
	IsOfModel
	IsPrimitive
	IsString
)

func (be BooleanExpressionFunction) String() string {
	functions := []string{"IS_BOOL", "IS_DEFINED", "IS_NULL", "IS_NUMBER", "IS_OBJECT", "IS_OF_MODEL", "IS_PRIMITIVE", "IS_STRING"}
	if !be.IsValid() {
		log.Fatalf("%d is not a valid function type", be)
	}
	return functions[be-1]
}

func (be BooleanExpressionFunction) IsValid() bool {
	switch be {
	case IsBool, IsDefined, IsNull, IsNumber, IsObject, IsOfModel, IsPrimitive, IsString:
		return true
	}
	return false
}
