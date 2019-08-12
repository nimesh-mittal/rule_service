package executor

import (
	"fmt"
	"time"
)

type Ruleset struct {
	ID        string `json:"id" validate:"required"`
	Name      string `json:"name"`
	StartDate time.Time
	EndDate   time.Time
	Enable    bool
	Rules     []Rule
}

type Rule struct {
	ID             string `json:"id" validate:"required"`
	Name           string `json:"name"`
	Priority       int
	Enable         bool
	WhenConditions []WhenCondition
	ThenActions    []ThenAction
}

type WhenCondition struct {
	ID        string `json:"id" validate:"required"`
	Field1    string
	Operator  string
	Field2    interface{}
	Threshold string
}

type ThenAction struct {
	ID       string `json:"id" validate:"required"`
	Field1   string
	Operator string
	Value    interface{}
}

func (r Ruleset) String() string {
	return fmt.Sprintf("%v:%v", r.ID, r.Rules)
}

func (r Rule) String() string {
	return fmt.Sprintf("%v:%v -> %v", r.ID, r.WhenConditions, r.ThenActions)
}

func (r WhenCondition) String() string {
	return fmt.Sprintf("%v:%v %v %v", r.ID, r.Field1, r.Operator, r.Field2)
}

func (r ThenAction) String() string {
	return fmt.Sprintf("%v:%v = %v", r.ID, r.Field1, r.Value)
}

type Record struct {
	ID     string
	Fields map[string]interface{}
}

func (r Record) String() string {
	return fmt.Sprintf("%v:%v", r.ID, r.Fields)
}
