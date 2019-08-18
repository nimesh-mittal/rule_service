package ruleset

import (
	"context"
	"errors"
	"rule_service/commons"
	"rule_service/models"
	"rule_service/pkg/evaluator"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/sirupsen/logrus"
)

// RepoContext holds repo state
type RepoContext struct {
	DB *mongo.Client
}

// NewRepoContext initialises Ruleset repo
func NewRepoContext(dbURL string) (*RepoContext, error) {

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
	return &RepoContext{DB: client}, nil
}

func (ctx *RepoContext) SafeClose() {
	err := ctx.DB.Disconnect(context.TODO())

	if err != nil {
		logrus.Fatal("failed to close database connection", err)
	}
}

func (ctx *RepoContext) ListRuleset(flowContext *models.FlowContext, limit int, offset int) (*[]evaluator.Ruleset, error) {
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

	var rulesets []evaluator.Ruleset
	for cursor.Next(context.TODO()) {
		var rs evaluator.Ruleset
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

// GetRuleset get ruleset by ID
func (ctx *RepoContext) GetRuleset(flowContext *models.FlowContext, rsid string) (*evaluator.Ruleset, error) {

	collection := ctx.DB.Database(commons.RULESET_DB).Collection(commons.RULESET_COLLECTION)
	filter := bson.D{{Key: "id", Value: rsid}}

	var rs evaluator.Ruleset
	err := collection.FindOne(context.TODO(), filter).Decode(&rs)

	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func (ctx *RepoContext) CreateRuleset(flowContext *models.FlowContext, ruleset *evaluator.Ruleset) (*evaluator.Ruleset, error) {
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

func GetBSON(entity *evaluator.Ruleset) *bson.D {
	// TODO: use reflect
	name := bson.E{Key: "name", Value: entity.Name}
	startDate := bson.E{Key: "start_date", Value: entity.StartDate}
	endDate := bson.E{Key: "end_date", Value: entity.EndDate}
	enable := bson.E{Key: "enable", Value: entity.Enable}
	rules := bson.E{Key: "rules", Value: entity.Rules}

	return &bson.D{name, startDate, endDate, enable, rules}
}

func (ctx *RepoContext) UpdateRuleset(flowContext *models.FlowContext, rsid string, entity *evaluator.Ruleset) (string, error) {
	collection := ctx.DB.Database(commons.RULESET_DB).Collection(commons.RULESET_COLLECTION)

	filter := bson.D{{Key: "id", Value: rsid}}
	newRecord := bson.D{{Key: "$set", Value: GetBSON(entity)}}
	res, err := collection.UpdateOne(context.TODO(), filter, newRecord)

	if err != nil {
		return commons.EMPTY, err
	}

	logrus.WithField(commons.TrackingID, flowContext).
		Info("updated MatchedCount, ModifiedCount, UpsertedCount, UpsertedID ",
			res.MatchedCount, res.ModifiedCount, res.UpsertedCount, res.UpsertedID)

	return rsid, nil
}

func (ctx *RepoContext) DeleteRuleset(flowContext *models.FlowContext, rsid string) (*evaluator.Ruleset, error) {
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
