package executor

import (
	"reflect"
	"testing"
	"time"
)

func TestGetValueFromRecordEmptyCheck(t *testing.T) {
	r1 := Record{ID: "101", Fields: map[string]interface{}{"field1": 10}}
	r2 := Record{ID: "201", Fields: map[string]interface{}{"field1": 10, "field2": "money"}}

	tests := []struct {
		FieldName     string
		Record1       Record
		ExpectedValue interface{}
		ErrorExpected bool
	}{
		{"field1", Record{ID: "r101"}, nil, true},
		{"field1", r1, 10, false},
		{"field2", r2, "money", false},
	}

	for i, test := range tests {
		t.Log("running % test", i+1)
		value, err := GetValueFromRecord(test.FieldName, test.Record1)

		if !test.ErrorExpected && err != nil {
			t.Errorf("GetValueFromRecord(field, record) = not raising error")
		}

		if !test.ErrorExpected && value != test.ExpectedValue {
			t.Errorf("GetValueFromRecord(field, record) = %d; expected = %d", value, test.ExpectedValue)
		}
	}
}

func TestGetValueFromRecordNumberCheck(t *testing.T) {
	value, err := GetValueFromRecord("field1", Record{ID: "101", Fields: map[string]interface{}{"field1": 10}})

	if err != nil {
		t.Errorf("GetValueFromRecord(field, record) = should not raise error")
	}

	if value != 10 {
		t.Errorf("GetValueFromRecord(field, record) = %d; expected = 10", value)
	}
}

func TestParseToType(t *testing.T) {
	tests := []struct {
		FieldValue    interface{}
		FieldType     string
		ExpectedType  reflect.Kind
		ErrorExpected bool
	}{
		{nil, NumberType, reflect.Float64, true},
		{10, NumberType, reflect.Float64, false},
		{float64(10.0), NumberType, reflect.Float64, false},
		{float32(10.0), NumberType, reflect.Float64, false},
		{int64(10), NumberType, reflect.Float64, false},
		{int32(10), NumberType, reflect.Float64, false},
		{int16(10), NumberType, reflect.Float64, false},
		{int8(10), NumberType, reflect.Float64, false},
		{"money", StringType, reflect.String, false},
		{nil, StringType, reflect.String, false},
		{[]int{1, 2, 3}, ArrayType, reflect.Array, true},
		{uint(1), ArrayType, reflect.Array, true},
	}

	for i, test := range tests {
		t.Log("running % test", i+1)
		v, err := ParseToType(test.FieldValue, test.FieldType)

		if !test.ErrorExpected && err != nil {
			t.Errorf("ParseToType(field, type) should not raise error; err=%s", err)
		}

		if !test.ErrorExpected && reflect.TypeOf(v).Kind() != test.ExpectedType {
			t.Errorf("ParseToType(field, type) failed")
		}
	}
}

func TestGetValue(t *testing.T) {
	tests := []struct {
		Field1         interface{}
		Field2         interface{}
		Operator       string
		ExpectedResult bool
		ErrorExpected  bool
	}{
		{10.0, 20.0, "<", true, false},
		{30.0, 20.0, ">", true, false},
		{10.0, 10.0, "<=", true, false},
		{10.0, 10.0, ">=", true, false},
		{10.0, 10.0, "=", true, false},
		{10.0, 10.0, "!=", false, false},
		{"hello world", "hello world", "==", true, false},
	}

	for i, test := range tests {
		t.Log("running % test", i+1)
		v, err := GetValue(test.Field1, test.Field2, test.Operator)

		if !test.ErrorExpected && err != nil {
			t.Errorf("GetValue(field, type) should not raise error; err=%s", err)
		}

		if !test.ErrorExpected && v != test.ExpectedResult {
			t.Errorf("GetValue(f1, f2, op) failed; f1=%q, f2=%q, op=%s", test.Field1, test.Field2, test.Operator)
		}
	}
}

func TestInferFieldType(t *testing.T) {
	tests := []struct {
		Operator      string
		ExpectedValue string
		ErrorExpected bool
		Error         error
	}{
		{">", NumberType, false, nil},
		{"<", NumberType, false, nil},
		{"==", StringType, false, nil},
		{"=", NumberType, false, nil},
		{"in", ArrayType, false, nil},
		{"##", "", true, ErrorOperatorNotSupported},
	}

	for i, test := range tests {
		t.Log("running % test", i+1)
		value, err := InferFieldType(test.Operator)

		if !test.ErrorExpected && err != nil {
			t.Errorf("InferFieldType(op) = not raising error")
		}

		if !test.ErrorExpected && value != test.ExpectedValue {
			t.Errorf("InferFieldType(op) = %s; expected = %s", value, test.ExpectedValue)
		}

		if test.ErrorExpected && err != test.Error {
			t.Errorf("got=%v; expected error=%v", err, test.Error)
		}
	}
}

func TestCheckWhenCondition(t *testing.T) {
	w1 := WhenCondition{"1", "field1", ">", "5", ""}
	w2 := WhenCondition{"1", "field1", ">=", "10", ""}
	w3 := WhenCondition{"1", "field1", ">=", "5.0", ""}
	w4 := WhenCondition{"1", "field2", "==", "money", ""}

	r1 := Record{ID: "101", Fields: map[string]interface{}{"field1": 10, "field2": "money"}}

	tests := []struct {
		WhenCondition WhenCondition
		Record        Record
		ExpectedValue bool
		ErrorExpected bool
		Error         error
	}{
		{w1, r1, true, false, nil},
		{w2, r1, true, false, nil},
		{w3, r1, true, false, nil},
		{w4, r1, true, false, nil},
	}

	for i, test := range tests {
		t.Log("running % test", i+1)
		value, err := CheckWhenCondition(test.WhenCondition, test.Record)

		if !test.ErrorExpected && err != nil {
			t.Errorf("CheckWhenCondition(condition, record) = not raising error")
		}

		if !test.ErrorExpected && value != test.ExpectedValue {
			t.Errorf("CheckWhenCondition(condition, record) = %v; expected = %v", value, test.ExpectedValue)
		}

		if test.ErrorExpected && err != test.Error {
			t.Errorf("got=%v; expected error=%v", err, test.Error)
		}
	}
}

func TestCheckRule(t *testing.T) {
	wc1 := []WhenCondition{{"1", "field1", ">", "5", ""}}
	ta1 := []ThenAction{{"1", "field1", "+", 10}}
	w1 := Rule{ID: "1", Name: "rule1", Priority: 0, Enable: true, WhenConditions: wc1, ThenActions: ta1}

	r1 := Record{ID: "101", Fields: map[string]interface{}{"field1": 10, "field2": "money"}}

	tests := []struct {
		Rule          Rule
		Record        Record
		ExpectedValue bool
		ErrorExpected bool
		Error         error
	}{
		{w1, r1, true, false, nil},
	}

	for i, test := range tests {
		t.Log("running % test", i+1)
		value, err := CheckRule(test.Rule, test.Record)

		if !test.ErrorExpected && err != nil {
			t.Errorf("TestCheckRule(rule, record) = not raising error")
		}

		if !test.ErrorExpected && value != test.ExpectedValue {
			t.Errorf("TestCheckRule(rule, record) = %v; expected = %v", value, test.ExpectedValue)
		}

		if test.ErrorExpected && err != test.Error {
			t.Errorf("got=%v; expected error=%v", err, test.Error)
		}
	}
}

func TestCheckRuleset(t *testing.T) {

	whenConditions1 := []WhenCondition{
		{"1", "field1", ">", "5", EmptyString},
		{"2", "field1", "<", "15", ""},
	}

	whenConditions2 := []WhenCondition{
		{"2", "field2", "==", "money", EmptyString},
	}

	whenConditions3 := []WhenCondition{
		{"3", "field2", "=", "money", EmptyString},
	}

	whenConditions4 := []WhenCondition{
		{"4", "field2", "#", "money", EmptyString},
	}

	whenConditions5 := []WhenCondition{
		{"5", "field2", "%", "money", EmptyString},
	}

	then1 := []ThenAction{{"1", "field1", "=", 10}}
	rule1 := Rule{ID: "r103", Name: "rule1", Enable: true, WhenConditions: whenConditions1, ThenActions: then1}
	rule2 := Rule{ID: "r104", Name: "rule2", Enable: true, WhenConditions: whenConditions2, ThenActions: then1}
	rule3 := Rule{ID: "r105", Name: "rule3", Enable: true, WhenConditions: whenConditions3, ThenActions: then1}
	rule4 := Rule{ID: "r105", Name: "rule4", Enable: true, WhenConditions: whenConditions4, ThenActions: then1}
	rule5 := Rule{ID: "r105", Name: "rule4", Enable: true, WhenConditions: whenConditions5, ThenActions: then1}

	ruleset1 := Ruleset{ID: "1", Name: "ruleset1", StartDate: time.Now().Add(-1 * time.Hour),
		EndDate: time.Now().Add(1 * time.Hour), Enable: true, Rules: []Rule{rule1}}
	ruleset2 := Ruleset{ID: "2", Name: "ruleset2", StartDate: time.Now().Add(-1 * time.Hour),
		EndDate: time.Now().Add(1 * time.Hour), Enable: true, Rules: []Rule{rule1, rule2}}
	ruleset3 := Ruleset{ID: "3", Name: "ruleset3", StartDate: time.Now(), EndDate: time.Now(), Enable: false,
		Rules: []Rule{rule1}}
	ruleset4 := Ruleset{ID: "4", Name: "ruleset4", StartDate: time.Now().Add(1 * time.Hour), EndDate: time.Now(),
		Enable: true, Rules: []Rule{rule1}}
	ruleset5 := Ruleset{ID: "5", Name: "ruleset5", StartDate: time.Now().Add(-1 * time.Hour),
		EndDate: time.Now().Add(1 * time.Hour), Enable: true, Rules: []Rule{rule3}}
	ruleset6 := Ruleset{ID: "6", Name: "ruleset6", StartDate: time.Now().Add(-1 * time.Hour),
		EndDate: time.Now().Add(1 * time.Hour), Enable: true, Rules: []Rule{rule4}}
	ruleset7 := Ruleset{ID: "7", Name: "ruleset7", StartDate: time.Now().Add(-1 * time.Hour),
		EndDate: time.Now().Add(1 * time.Hour), Enable: true, Rules: []Rule{rule5}}

	record1 := Record{ID: "101", Fields: map[string]interface{}{Field1: 10, Field2: "money"}}
	record2 := Record{ID: "201", Fields: map[string]interface{}{Field2: "money"}}
	record3 := Record{ID: "301", Fields: map[string]interface{}{Field1: 30, Field2: "no money"}}
	record4 := Record{ID: "401", Fields: map[string]interface{}{Field1: "a", Field2: "no money"}}

	tests := []struct {
		Ruleset       Ruleset
		Record        Record
		ExpectedValue string
		ErrorExpected bool
		Error         error
	}{
		{ruleset1, record1, rule1.ID, false, nil},
		{ruleset2, record1, rule1.ID, false, nil},
		{ruleset2, record2, rule2.ID, false, nil},
		{ruleset1, record2, EmptyString, true, ErrorEmptyMatch},
		{ruleset2, record3, EmptyString, true, ErrorEmptyMatch},
		{ruleset1, record3, EmptyString, true, ErrorEmptyMatch},
		{ruleset3, record3, EmptyString, true, ErrorEmptyMatch},
		{ruleset4, record3, EmptyString, true, ErrorEmptyMatch},
		{ruleset5, record3, EmptyString, true, ErrorUnableToParseField},
		{ruleset6, record3, EmptyString, true, ErrorOperatorNotSupported},
		{ruleset1, record4, EmptyString, true, ErrorUnableToParseField},
		{ruleset7, record2, EmptyString, true, ErrorOperatorNotSupported},
	}

	for i, test := range tests {
		t.Log("running % test", i+1)
		value, err := CheckRuleset(&test.Ruleset, &test.Record, MatchFirst)

		if !test.ErrorExpected && err != nil {
			t.Errorf("CheckRuleset(ruleset, record) = error is not expected; err=%v", err)
		}

		if !test.ErrorExpected && value.ID != test.ExpectedValue {
			t.Errorf("CheckRuleset(ruleset, record) = %v; expected = %v", value, test.ExpectedValue)
		}

		if test.ErrorExpected && err != test.Error {
			t.Errorf("got=%v; expected error=%v", err, test.Error)
		}
	}
}

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseToType(10, NumberType)
	}
}
