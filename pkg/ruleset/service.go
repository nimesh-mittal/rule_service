package ruleset

import (
	"github.com/sirupsen/logrus"
	"rule_service/commons"
	"rule_service/config"
	"rule_service/models"
	"rule_service/pkg/evaluator"
)

type RulesetServiceContext struct {
	repo *RulesetRepoContext
}

func NewRulesetServiceContext(flowContext *models.FlowContext) *RulesetServiceContext {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).
		Info("setting up ruleset service")
	repo, _ := NewRulesetRepoContext(config.GetInstance().Database.MongoURL)

	ctx := RulesetServiceContext{repo: repo}

	return &ctx
}

// Evaluate record for ruleset id
func (ctx *RulesetServiceContext) Evaluate(flowContext *models.FlowContext,
	rulesetID string, record *map[string]interface{}, strategy string) (*evaluator.Rule, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("evaluating ruleset")

	ruleset, err := ctx.repo.GetRuleset(flowContext, rulesetID)

	if err != nil {
		return nil, err
	}

	rule, err := evaluator.CheckRuleset(ruleset, &evaluator.Record{Fields: *record}, strategy)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (ctx *RulesetServiceContext) ListRuleset(flowContext *models.FlowContext, limit int, offset int) (*[]evaluator.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("listing rulesets")

	rulesets, err := ctx.repo.ListRuleset(flowContext, limit, offset)

	if err != nil {
		return nil, err
	}

	return rulesets, nil
}

func (ctx *RulesetServiceContext) GetRuleset(flowContext *models.FlowContext, id string) (*evaluator.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("get ruleset by id")

	c, err := ctx.repo.GetRuleset(flowContext, id)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (ctx *RulesetServiceContext) CreateRuleset(flowContext *models.FlowContext, ruleset *evaluator.Ruleset) (*evaluator.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("creating ruleset")

	c, err := ctx.repo.CreateRuleset(flowContext, ruleset)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (ctx *RulesetServiceContext) UpdateRuleset(flowContext *models.FlowContext, id string, entity *evaluator.Ruleset) (string, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("updating ruleset")

	id, err := ctx.repo.UpdateRuleset(flowContext, id, entity)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (ctx *RulesetServiceContext) DeleteRuleset(flowContext *models.FlowContext, id string) (*evaluator.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("deleting ruleset")

	c, err := ctx.repo.DeleteRuleset(flowContext, id)

	if err != nil {
		return nil, err
	}
	return c, nil
}
