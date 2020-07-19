package query

import (
	"context"
	"encoding/json"
	"os"

	"github.com/mongodb/mongo-tools-common/bsonutil"
	"github.com/pkg/errors"
	"github.com/yolossn/query2metric/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoQuery struct {
	connection string
	client     *mongo.Client
}

func NewMongoConn(connnectionURL string) (CountQuery, error) {
	connnectionString := os.Getenv(connnectionURL)
	if connnectionString == "" {
		return nil, errors.New("connnectionString is empty")
	}
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(connnectionString))
	if err != nil {
		return nil, err
	}
	return &mongoQuery{connnectionURL, mongoClient}, err
}

func (m *mongoQuery) Count(metric config.Metric) (int64, error) {

	ctx := context.Background()
	err := m.client.Ping(ctx, nil)
	if err != nil {
		err = m.client.Connect(ctx)
		if err != nil {
			return 0, err
		}
	}

	query := map[string]interface{}{}
	if metric.Query != "" {
		err = json.Unmarshal([]byte(metric.Query), &query)
		if err != nil {
			return 0, err
		}
	}
	bsonQuery, err := bsonutil.ConvertLegacyExtJSONValueToBSON(query)
	if err != nil {
		return 0, err
	}
	count, err := m.client.Database(metric.Database).Collection(metric.Collection).CountDocuments(ctx, bsonQuery)
	if err != nil {
		return 0, err
	}
	return count, nil
}
