package executor

var GreaterThan = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) > v2.(float64) }
var LessThan = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) < v2.(float64) }
var NumericEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) == v2.(float64) }
var NumericNotEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) != v2.(float64) }
var GreaterThanEq = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) >= v2.(float64) }
var LessThanEq = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) <= v2.(float64) }
var StringEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(string) == v2.(string) }
var StringNotEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(string) != v2.(string) }

var OperatorType = map[string]string{
	">":      NumberType,
	"<":      NumberType,
	">=":     NumberType,
	"<=":     NumberType,
	"=":      NumberType,
	"==":     StringType,
	"!=":     NumberType,
	"!==":    StringType,
	"in":     ArrayType,
	"not_in": ArrayType,
	"%":      StringType,
}

var OperatorFunction = map[string]func(interface{}, interface{}) bool{
	">":  GreaterThan,
	"<":  LessThan,
	">=": GreaterThanEq,
	"<=": LessThanEq,

	"=":  NumericEquals,
	"==": StringEquals,

	"!=":  NumericNotEquals,
	"!==": StringNotEquals,

	"in":     NumericEquals,
	"not_in": NumericEquals,
}
