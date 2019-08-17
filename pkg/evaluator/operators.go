package evaluator

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
)

var GreaterThan = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) > v2.(float64) }
var LessThan = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) < v2.(float64) }
var NumericEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) == v2.(float64) }
var NumericNotEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) != v2.(float64) }
var GreaterThanEq = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) >= v2.(float64) }
var LessThanEq = func(v1 interface{}, v2 interface{}) bool { return v1.(float64) <= v2.(float64) }
var StringEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(string) == v2.(string) }
var StringNotEquals = func(v1 interface{}, v2 interface{}) bool { return v1.(string) != v2.(string) }
var NumberInArray = func(v1 interface{}, v2 interface{}) bool { return InNum(v1.(float64), v2.([]float64)) }
var StringInArray = func(v1 interface{}, v2 interface{}) bool { return InString(v1.(string), v2.([]string)) }

func InNum(f1 float64, f2 []float64) bool {
	for _, v := range f2 {
		if v == f1 {
			return true
		}
	}
	return false
}

func InString(f1 string, f2 []string) bool {
	for _, v := range f2 {
		if v == f1 {
			return true
		}
	}
	return false
}

var OperatorType = map[string]string{
	">":          NumberType,
	"<":          NumberType,
	">=":         NumberType,
	"<=":         NumberType,
	"=":          NumberType,
	"==":         StringType,
	"!=":         NumberType,
	"!==":        StringType,
	"in_numbers": NumberArrayType,
	"in_strings": StringArrayType,
	"%":          StringType,
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

	"in_numbers": NumberInArray,
	"in_strings": StringInArray,
}

func ParseToType(field interface{}, fieldType string) (interface{}, error) {

	if fieldType == NumberType {
		if field == nil {
			return nil, ErrorNumberCantnotBeNil
		}

		switch reflect.TypeOf(field).Kind() {
		case reflect.Float64:
			value := field.(float64)
			return value, nil
		case reflect.Float32:
			f := field.(float32)
			return float64(f), nil
		case reflect.Int:
			f := field.(int)
			return float64(f), nil
		case reflect.Int64:
			f, _ := field.(int64)
			return float64(f), nil
		case reflect.Int32:
			f, _ := field.(int32)
			return float64(f), nil
		case reflect.Int16:
			f, _ := field.(int16)
			return float64(f), nil
		case reflect.Int8:
			f, _ := field.(int8)
			return float64(f), nil
		case reflect.String:
			f, err := strconv.ParseFloat(field.(string), 64)
			if err != nil {
				logrus.WithField("field", field).
					WithField("type", "float64").
					Error(err)
				return nil, ErrorUnableToParseField
			}
			return float64(f), nil
		default:
			return nil, ErrorUnableToInferFieldType
		}
	} else if fieldType == StringType {
		if field == nil {
			return "", nil
		}
		return field.(string), nil
	} else if fieldType == NumberArrayType {
		if field == nil {
			return []float64{}, nil
		}
		s := reflect.ValueOf(field)
		var res []float64
		for i := 0; i < s.Len(); i++ {
			res = append(res, s.Index(i).Interface().(float64))
		}
		return res, nil
	} else if fieldType == StringArrayType {
		if field == nil {
			return []string{}, nil
		}
		s := reflect.ValueOf(field)
		var res []string
		for i := 0; i < s.Len(); i++ {
			res = append(res, s.Index(i).Interface().(string))
		}
		return res, nil
	} else {
		logrus.WithField("field", field).
			WithField("fieldType", fieldType).
			Error(ErrorOperatorNotSupported)
		return nil, ErrorOperatorNotSupported
	}
}
