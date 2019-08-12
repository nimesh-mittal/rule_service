package executor

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"time"
)

func ApplyRule(rule *Rule, record *Record) *Record {
	if rule == nil {
		return record
	}

	for _, ta := range rule.ThenActions {
		//Todo: handle error
		value, _ := strconv.ParseFloat(ta.Value.(string), 64)
		if ta.Operator == "=" {
			record.Fields[ta.Field1] = value
		} else if ta.Operator == "+" {
			record.Fields[ta.Field1] = record.Fields[ta.Field1].(float64) + value
		} else if ta.Operator == "-" {
			record.Fields[ta.Field1] = record.Fields[ta.Field1].(float64) - value
		} else if ta.Operator == "*" {
			record.Fields[ta.Field1] = record.Fields[ta.Field1].(float64) * value
		}
	}
	return record
}

func CheckRuleset(rs *Ruleset, record *Record, strategy string) (Rule, error) {

	if CheckRulesetIneligible(*rs) {
		return Rule{}, ErrorEmptyMatch
	}

	var matchedRules []Rule
	for _, rule := range rs.Rules {
		ruleMatched, err := CheckRule(rule, *record)
		if err != nil {
			if err == ErrorFieldNotInRecord {
				logrus.
					WithField("rule", rule).
					WithField("err", err).
					Warn("skip processing rule")
				continue
			} else {
				logrus.Error("error in rule processing", rule)
				return Rule{}, err
			}
		}

		if ruleMatched {
			matchedRules = append(matchedRules, rule)
		}

		if strategy == MatchFirst && len(matchedRules) > 0 {
			return matchedRules[0], nil
		}
	}

	if strategy == MatchHighestPriority && len(matchedRules) > 0 {
		priority := 0
		rule := matchedRules[0]
		for _, r := range matchedRules {
			if r.Priority > priority {
				priority = r.Priority
				rule = r
			}
		}
		return rule, nil
	}

	return Rule{}, ErrorEmptyMatch
}

func CheckRulesetIneligible(rs Ruleset) bool {
	if rs.Enable && time.Now().After(rs.StartDate) && time.Now().Before(rs.EndDate) {
		return false
	}

	logrus.
		WithField("Enable", rs.Enable).
		WithField("StartDate", rs.StartDate).
		WithField("EndDate", rs.EndDate).
		Warn("rule is not eligible")
	return true
}

func CheckRule(rule Rule, record Record) (bool, error) {
	response := true

	for _, when := range rule.WhenConditions {
		res, err := CheckWhenCondition(when, record)

		if err != nil {
			logrus.WithField("when", when).Error("error in when condition processing")
			return false, err
		}

		response = response && res
	}

	return response, nil
}

func CheckWhenCondition(whenCondition WhenCondition, record Record) (bool, error) {
	op := whenCondition.Operator
	fieldName := whenCondition.Field1
	fieldValue := whenCondition.Field2

	value1, err := GetValueFromRecord(fieldName, record)
	if err != nil {
		return false, err
	}

	fieldType, err := InferFieldType(op)
	if err != nil {
		return false, err
	}

	value2, err := ParseToType(fieldValue, fieldType)
	if err != nil {
		return false, err
	}

	value1, err = ParseToType(value1, fieldType)
	if err != nil {
		return false, err
	}

	value, err := GetValue(value1, value2, op)
	if err != nil {
		logrus.Error(ErrorFailedToApplyOperator, value1, value2, op)
		return false, err
	}

	return value, nil
}

func InferFieldType(op string) (string, error) {
	val, ok := OperatorType[op]
	if !ok {
		return "", ErrorOperatorNotSupported
	}
	return val, nil
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
				logrus.
					WithField("field", field).
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
	} else {
		logrus.
			WithField("field", field).
			WithField("fieldType", fieldType).
			Error(ErrorOperatorNotSupported)
		return nil, ErrorOperatorNotSupported
	}

	return nil, ErrorUnknownCause
}

func GetValue(f1 interface{}, f2 interface{}, op string) (bool, error) {

	function, ok := OperatorFunction[op]
	if !ok {
		logrus.
			WithField("f1", f1).
			WithField("f2", f2).
			WithField("op", op).
			Error(ErrorOperatorNotSupported)
		return false, ErrorOperatorNotSupported
	}

	return function(f1, f2), nil
}

func GetValueFromRecord(fieldName string, record Record) (interface{}, error) {
	value, ok := record.Fields[fieldName]

	if !ok {
		logrus.
			WithField("fieldName", fieldName).
			WithField("record", record).
			Error(ErrorFieldNotInRecord)
		return nil, ErrorFieldNotInRecord
	}

	return value, nil
}
