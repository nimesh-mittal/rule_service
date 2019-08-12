package ruleset

import (
	"github.com/sirupsen/logrus"
	"rule_service/commons"
	"rule_service/config"
	"rule_service/executor"
	"rule_service/models"
)

type RulesetServiceContext struct {
	repo *RulesetRepoContext
}

func NewRulesetServiceContext(flowContext *models.FlowContext) *RulesetServiceContext {
	repo, _ := NewRulesetRepoContext(config.GetInstance().Database.Dialect,
		config.GetInstance().Database.MongoURL)

	ctx := RulesetServiceContext{repo: repo}

	return &ctx
}

func (ctx *RulesetServiceContext) Evaluate(flowContext *models.FlowContext,
	rulesetID string, record *map[string]interface{}, strategy string) (*executor.Rule, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("evaluating ruleset")

	ruleset, err := ctx.repo.GetRuleset(flowContext, rulesetID)

	if err != nil {
		return nil, err
	}

	rule, err := executor.CheckRuleset(ruleset, &executor.Record{Fields: *record}, strategy)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (ctx *RulesetServiceContext) ListRuleset(flowContext *models.FlowContext, limit int, offset int) (*[]executor.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("listing rulesets")

	rulesets, err := ctx.repo.ListRuleset(flowContext, limit, offset)

	if err != nil {
		return nil, err
	}

	return rulesets, nil
}

func (ctx *RulesetServiceContext) GetRuleset(flowContext *models.FlowContext, rsid string) (*executor.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("get geofence by id")

	c, err := ctx.repo.GetRuleset(flowContext, rsid)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (ctx *RulesetServiceContext) CreateRuleset(flowContext *models.FlowContext, ruleset *executor.Ruleset) (*executor.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("creating ruleset")

	c, err := ctx.repo.CreateRuleset(flowContext, ruleset)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (ctx *RulesetServiceContext) UpdateRuleset(flowContext *models.FlowContext, rsid string, entity *executor.Ruleset) (string, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("updating geofence")

	rsid, err := ctx.repo.UpdateRuleset(flowContext, rsid, entity)

	if err != nil {
		return "", err
	}

	return rsid, nil
}

func (ctx *RulesetServiceContext) DeleteRuleset(flowContext *models.FlowContext, rsid string) (*executor.Ruleset, error) {
	logrus.WithField(commons.TrackingID, flowContext.TrackingID).Info("deleting geofence")

	c, err := ctx.repo.DeleteRuleset(flowContext, rsid)

	if err != nil {
		return nil, err
	}
	return c, nil
}
