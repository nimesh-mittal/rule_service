package ruleset

import (
	"rule_service/evaluator"
	"time"
)

type RulesetDTO struct {
	Name      string `json:"name"`
	StartDate time.Time
	EndDate   time.Time
	Enable    bool
	Rules     []evaluator.Rule
}

type EvaluateResposeDTO struct {
	Record       map[string]interface{}
	MatchingRule *evaluator.Rule
}

type EvaluateRequestDTO struct {
	Record    map[string]interface{}
	RulesetID string
}
