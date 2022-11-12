package repository

import (
	"context"

	"github.com/xuhaojun/cinnox-homework/dto"
	"github.com/xuhaojun/cinnox-homework/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LineEventRepository struct {
	c *mongo.Collection
}

func NewLineEventRepository(db *mongo.Database) *LineEventRepository {
	return &LineEventRepository{c: db.Collection(CollNameLineEvent)}
}

func (repo *LineEventRepository) InsertMany(events []*dto.LineEvent) (*mongo.InsertManyResult, error) {
	return repo.c.InsertMany(context.Background(), util.ToSliceOfAny(events))
}

func (repo *LineEventRepository) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uniqueId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := repo.c.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return err
	}
	return nil
}
