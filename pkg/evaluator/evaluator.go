package evaluator

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// ApplyRule returns modified record
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

// CheckRuleset returns matched rule based on the provided strategy
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

// checks for the eligibility of the rule
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

// CheckRule returns true if the rule resolves for record values
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

// CheckWhenCondition return true if when condition resolve for the record
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

	LHSType := fieldType
	if fieldType == NumberArrayType {
		LHSType = NumberType
	} else if fieldType == StringArrayType {
		LHSType = StringType
	}

	value1, err = ParseToType(value1, LHSType)
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

// InferFieldType infers field data type for the operator
func InferFieldType(op string) (string, error) {
	val, ok := OperatorType[op]
	if !ok {
		return "", ErrorOperatorNotSupported
	}
	return val, nil
}

// GetValue returns value after applying op on f1 and f2 field
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

// GetValueFromRecord returns value from the record for field fieldName
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
