package repository

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/xuhaojun/cinnox-homework/dto"
	"github.com/xuhaojun/cinnox-homework/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CinnoxMessageRepository struct {
	c *mongo.Collection
}

func NewCinnoxMessageRepository(db *mongo.Database) *CinnoxMessageRepository {
	return &CinnoxMessageRepository{c: db.Collection(CollNameCinnoxMessage)}
}

func (repo *CinnoxMessageRepository) SaveLineMessages(msgs []*dto.CinnoxMessage[linebot.Message]) (*mongo.InsertManyResult, error) {
	return repo.c.InsertMany(context.Background(), util.ToSliceOfAny(msgs))
}

func (repo *CinnoxMessageRepository) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "sourceType", Value: 1}, {Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := repo.c.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return err
	}
	return nil
}
