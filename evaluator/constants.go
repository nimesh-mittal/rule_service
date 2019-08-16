package evaluator

import (
	"errors"
)

// general
const (
	EmptyString = ""
	Field1      = "field1"
	Field2      = "field2"
)

// types
const (
	NumberType      = "Number"
	StringType      = "String"
	NumberArrayType = "NumberArray"
	StringArrayType = "StringArray"
)

// errors
var (
	ErrorOperatorNotSupported   = errors.New("operator not supported")
	ErrorFieldNotInRecord       = errors.New("field is missing in the record")
	ErrorFieldTypeNotSupported  = errors.New("field type not supported")
	ErrorUnableToInferFieldType = errors.New("unable to infer field type")
	ErrorUnableToParseField     = errors.New("unable to parse field")
	ErrorFailedToApplyOperator  = errors.New("failed to apply operator on values")
	ErrorNumberCantnotBeNil     = errors.New("number can not be nil")
	ErrorUnknownCause           = errors.New("unknown cause of this error")
	ErrorEmptyMatch             = errors.New("unable to match with any rule")
)

//  rule execution strategy
const (
	MatchFirst           = "MatchFirst"
	MatchHighestPriority = "MatchHighestPriority"
)
