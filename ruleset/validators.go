package ruleset

import (
	uuid2 "github.com/google/uuid"
	"rule_service/executor"
)

func ToRuleset(rulesetDTO RulesetDTO) executor.Ruleset {
	var entity executor.Ruleset
	entity.ID = uuid2.New().String()
	entity.Name = rulesetDTO.Name
	entity.StartDate = rulesetDTO.StartDate
	entity.EndDate = rulesetDTO.EndDate
	entity.Enable = rulesetDTO.Enable
	entity.Rules = rulesetDTO.Rules
	for i := 0; i < len(entity.Rules); i++ {
		entity.Rules[i].ID = uuid2.New().String()
		for w := 0; w < len(entity.Rules[i].WhenConditions); w++ {
			entity.Rules[i].WhenConditions[w].ID = uuid2.New().String()
		}

		for t := 0; t < len(entity.Rules[i].ThenActions); t++ {
			entity.Rules[i].ThenActions[t].ID = uuid2.New().String()
		}
	}
	return entity
}
