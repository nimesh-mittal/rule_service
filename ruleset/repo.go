package ruleset

import (
	"context"
	"errors"
	"rule_service/commons"
	"rule_service/executor"
	"rule_service/models"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/sirupsen/logrus"
)

type RulesetRepoContext struct {
	DB *mongo.Client
}

func NewRulesetRepoContext(dialect string, dbURL string) (*RulesetRepoContext, error) {

	// setup mongo client
	ctx1, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx1, dbURL)

	if err != nil {
		logrus.WithField("URL", dbURL).Fatal("failed to connect database")
		return nil, err
	}

	// ping mongo to test connection
	ctx2, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx2, nil)

	if err != nil {
		logrus.Info("failed to ping database", err)
	}

	defer logrus.WithField("URL", dbURL).Info("mongo database setup completed")
	return &RulesetRepoContext{DB: client}, nil
}

func (ctx *RulesetRepoContext) SafeClose() {
	err := ctx.DB.Disconnect(context.TODO())

	if err != nil {
		logrus.Fatal("failed to close database connection", err)
	}
}

func (ctx *RulesetRepoContext) ListRuleset(flowContext *models.FlowContext, limit int, offset int) (*[]executor.Ruleset, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	collection := ctx.DB.Database(commons.RULESET_DB).Collection(commons.RULESET_COLLECTION)

	// TODO: take filter as parameter
	filter := bson.D{}
	cursor, err := collection.Find(context.TODO(), filter, findOptions)

	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, err
	}

	var rulesets []executor.Ruleset
	for cursor.Next(context.TODO()) {
		var rs executor.Ruleset
		err := cursor.Decode(&rs)
		if err != nil {
			return nil, err
		}
		rulesets = append(rulesets, rs)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &rulesets, nil
}

// TODO: make it efficient, current implementation is very hacky
func (ctx *RulesetRepoContext) GetRuleset(flowContext *models.FlowContext, rsid string) (*executor.Ruleset, error) {
	findOptions := options.Find()
	findOptions.SetLimit(1)
	findOptions.SetSkip(0)

	collection := ctx.DB.Database(commons.RULESET_DB).Collection(commons.RULESET_COLLECTION)
	filter := bson.D{{Key: "id", Value: rsid}}
	cursor, err := collection.Find(context.TODO(), filter, findOptions)

	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, err
	}

	var rulesets []executor.Ruleset
	for cursor.Next(context.TODO()) {
		var rs executor.Ruleset
		err := cursor.Decode(&rs)
		if err != nil {
			return nil, err
		}
		rulesets = append(rulesets, rs)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(rulesets) <= 0 {
		return nil, errors.New("record not found")
	}

	return &rulesets[0], nil
}

func (ctx *RulesetRepoContext) CreateRuleset(flowContext *models.FlowContext, ruleset *executor.Ruleset) (*executor.Ruleset, error) {
	collection := ctx.DB.Database(commons.RULESET_DB).Collection(commons.RULESET_COLLECTION)
	ctx1, _ := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := collection.InsertOne(ctx1, ruleset)

	if err != nil {
		return nil, err
	} else {
		logrus.WithField(commons.TrackingID, flowContext).Info("created record with id ", resp.InsertedID)
	}

	return ruleset, nil
}

func GetBSON(entity *executor.Ruleset) *bson.D {
	// TODO: use reflect
	name := bson.E{Key: "name", Value: entity.Name}
	startDate := bson.E{Key: "start_date", Value: entity.StartDate}
	endDate := bson.E{Key: "end_date", Value: entity.EndDate}
	enable := bson.E{Key: "enable", Value: entity.Enable}
	rules := bson.E{Key: "rules", Value: entity.Rules}

	return &bson.D{name, startDate, endDate, enable, rules}
}

func (ctx *RulesetRepoContext) UpdateRuleset(flowContext *models.FlowContext, rsid string, entity *executor.Ruleset) (string, error) {
	collection := ctx.DB.Database(commons.RULESET_DB).Collection(commons.RULESET_COLLECTION)

	filter := bson.D{{Key: "id", Value: rsid}}
	new := bson.D{{Key: "$set", Value: GetBSON(entity)}}
	res, err := collection.UpdateOne(context.TODO(), filter, new)

	if err != nil {
		return commons.EMPTY, err
	}

	logrus.WithField(commons.TrackingID, flowContext).
		Info("updated MatchedCount, ModifiedCount, UpsertedCount, UpsertedID ",
			res.MatchedCount, res.ModifiedCount, res.UpsertedCount, res.UpsertedID)

	return rsid, nil
}

func (ctx *RulesetRepoContext) DeleteRuleset(flowContext *models.FlowContext, rsid string) (*executor.Ruleset, error) {
	collection := ctx.DB.Database(commons.RULESET_DB).Collection(commons.RULESET_COLLECTION)

	if rsid == commons.EMPTY {
		logrus.Error("ruleset id is empty")
		return nil, errors.New("ruleset id is empty")
	}

	ruleset, err := ctx.GetRuleset(flowContext, rsid)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	filter := bson.D{{Key: "id", Value: rsid}}
	ctx1, _ := context.WithTimeout(context.Background(), 5*time.Second)
	deleteResult, err := collection.DeleteOne(ctx1, filter)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	logrus.WithField(commons.TrackingID, flowContext).Info("number of ruleset deleted is ", deleteResult.DeletedCount)
	return ruleset, nil
}
