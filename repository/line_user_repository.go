package repository

import (
	"context"

	"github.com/samber/lo"
	"github.com/xuhaojun/cinnox-homework/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LineUserRepository struct {
	c *mongo.Collection
}

func NewLineUserRepository(db *mongo.Database) *LineUserRepository {
	return &LineUserRepository{c: db.Collection(CollNameLineUser)}
}

func (repo *LineUserRepository) UpsertMany(lusers []*dto.LineUser) (*mongo.BulkWriteResult, error) {
	ops := lo.Map(lusers, func(item *dto.LineUser, index int) mongo.WriteModel {
		op := mongo.NewUpdateOneModel()
		op.SetUpsert(true)
		op.SetFilter(bson.M{"id": item.Id})
		op.SetUpdate(bson.M{"$set": item})
		return op
	})
	return repo.c.BulkWrite(context.Background(), ops)
}

func (repo *LineUserRepository) FindAll() ([]*dto.LineUser, error) {
	cur, err := repo.c.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var lineUsers []*dto.LineUser
	err = cur.All(context.Background(), &lineUsers)
	if err != nil {
		return nil, err
	}
	return lineUsers, nil
}

func (repo *LineUserRepository) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := repo.c.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return err
	}
	return nil
}
