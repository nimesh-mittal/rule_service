package ruleset

import (
	"rule_service/executor"
	"time"
)

type RulesetDTO struct {
	Name      string `json:"name"`
	StartDate time.Time
	EndDate   time.Time
	Enable    bool
	Rules     []executor.Rule
}

type EvaluateResposeDTO struct {
	Record       map[string]interface{}
	MatchingRule *executor.Rule
}

type EvaluateRequestDTO struct {
	Record    map[string]interface{}
	RulesetID string
}
