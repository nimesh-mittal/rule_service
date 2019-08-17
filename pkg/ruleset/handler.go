package ruleset

import (
	"encoding/json"
	"github.com/go-chi/chi"
	uuid2 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"opensource/rule_service/commons"
	"opensource/rule_service/models"
	"opensource/rule_service/pkg/evaluator"
	"strconv"
)

type RulesetContext struct {
	service *RulesetServiceContext
}

func NewRulesetContext() *RulesetContext {
	flowContext := &models.FlowContext{TrackingID: uuid2.New().String()}
	service := NewRulesetServiceContext(flowContext)
	return &RulesetContext{service: service}
}

func (bc *RulesetContext) SafeClose() {

}

func (ctx *RulesetContext) NewRulesetRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", ctx.ListRuleset)
	r.Get("/{RulesetID}", ctx.GetRuleset)
	r.Post("/_evaluate", ctx.EvaluateRuleset)
	r.Post("/", ctx.CreateRuleset)
	r.Put("/{RulesetID}", ctx.UpdateRuleset)
	r.Delete("/{RulesetID}", ctx.DeleteRuleset)

	return r
}

func (ctx *RulesetContext) EvaluateRuleset(w http.ResponseWriter, r *http.Request) {
	flowContext := &models.FlowContext{TrackingID: uuid2.New().String()}

	strategy := r.URL.Query().Get("strategy")
	if strategy == "" {
		logrus.Info("strategy is not set so setting it to MatchFirst")
		strategy = evaluator.MatchFirst
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var evaluateRequestDTO EvaluateRequestDTO
	err := decoder.Decode(&evaluateRequestDTO)

	if err != nil {
		res := commons.MakeResp(nil, commons.EMPTY, err)
		w.Write(res)
		return
	}

	rule, err := ctx.service.Evaluate(flowContext, evaluateRequestDTO.RulesetID, &evaluateRequestDTO.Record, strategy)
	logrus.Info(rule)
	record := evaluator.ApplyRule(rule, &evaluator.Record{Fields: evaluateRequestDTO.Record})

	evaluateResposeDTO := EvaluateResposeDTO{Record: record.Fields, MatchingRule: rule}
	response := commons.MakeResp(evaluateResposeDTO, commons.EMPTY, err)

	w.Write(response)
	return
}

func (ctx *RulesetContext) GetRuleset(w http.ResponseWriter, r *http.Request) {
	flowContext := &models.FlowContext{TrackingID: uuid2.New().String()}

	rsid := chi.URLParam(r, "RulesetID")
	logrus.Info("fetching Ruleset for id ", rsid)

	res, err := ctx.service.GetRuleset(flowContext, rsid)

	if err == nil {
		res := commons.MakeResp(res, commons.EMPTY, nil)
		w.Write(res)
		return
	} else {
		res := commons.MakeResp(nil, commons.EMPTY, err)
		w.Write(res)
		return
	}
}

func (ctx *RulesetContext) ListRuleset(w http.ResponseWriter, r *http.Request) {
	flowContext := &models.FlowContext{TrackingID: uuid2.New().String()}

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)

	offsetStr := r.URL.Query().Get("offset")
	offset, _ := strconv.Atoi(offsetStr)
	logrus.Info("fetching Ruleset for limit,offset = ", limit, offset)

	res, err := ctx.service.ListRuleset(flowContext, limit, offset)

	if err == nil {
		res := commons.MakeResp(res, commons.EMPTY, nil)
		w.Write(res)
		return
	} else {
		res := commons.MakeResp(nil, commons.EMPTY, err)
		w.Write(res)
		return
	}
}

func (ctx *RulesetContext) CreateRuleset(w http.ResponseWriter, r *http.Request) {
	flowContext := &models.FlowContext{TrackingID: uuid2.New().String()}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var rulesetDTO RulesetDTO
	err := decoder.Decode(&rulesetDTO)

	if err == nil {
		entity := ToRuleset(rulesetDTO)

		ctx.service.CreateRuleset(flowContext, &entity)
		res := commons.MakeResp("ruleset created successfully", commons.EMPTY, nil)
		w.Write(res)
		return
	} else {
		res := commons.MakeResp(nil, commons.EMPTY, err)
		w.Write(res)
		return
	}
}

func (ctx *RulesetContext) UpdateRuleset(w http.ResponseWriter, r *http.Request) {
	flowContext := &models.FlowContext{TrackingID: uuid2.New().String()}

	rsid := chi.URLParam(r, "RulesetID")
	logrus.Info("updating Ruleset for id ", rsid)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var rulesetDTO RulesetDTO
	err := decoder.Decode(&rulesetDTO)

	if err == nil {
		entity := ToRuleset(rulesetDTO)

		ctx.service.UpdateRuleset(flowContext, rsid, &entity)
		res := commons.MakeResp("ruleset updated successfully", commons.EMPTY, nil)
		w.Write(res)
		return
	} else {
		res := commons.MakeResp(nil, commons.EMPTY, err)
		w.Write(res)
		return
	}
}

func (ctx *RulesetContext) DeleteRuleset(w http.ResponseWriter, r *http.Request) {
	flowContext := &models.FlowContext{TrackingID: uuid2.New().String()}

	rsid := chi.URLParam(r, "RulesetID")
	logrus.Info("deleting Ruleset for id ", rsid)

	ctx.service.DeleteRuleset(flowContext, rsid)
}
