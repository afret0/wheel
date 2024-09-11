package database

import (
	"context"
	"github.com/afret0/wheel/log"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type Repository struct {
	collection *mongo.Collection
	logger     *logrus.Logger
}

func GetRepository(db *MongoDB, collection string) *Repository {

	if collection == "" {
		panic("collection name is empty")
	}

	r := new(Repository)
	r.collection = db.GetCollection(collection)
	r.logger = log.GetLogger()

	return r
}

func (r *Repository) FindOne(ctx context.Context, entity interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	one := r.collection.FindOne(ctx, filter, opts...)
	err := one.Decode(entity)
	return err
}

func (r *Repository) Find(ctx context.Context, entityList interface{}, filter interface{}, opts ...*options.FindOptions) error {
	resultV := reflect.ValueOf(entityList)
	if resultV.Kind() != reflect.Ptr || resultV.Elem().Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}

	sliceV := resultV.Elem()
	elemT := sliceV.Type().Elem()

	cur, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer func() { _ = cur.Close(ctx) }()

	for cur.Next(ctx) {
		elemP := reflect.New(elemT)
		if err = cur.Decode(elemP.Interface()); err != nil {
			return err
		}
		sliceV = reflect.Append(sliceV, elemP.Elem())
	}
	if err = cur.Err(); err != nil {
		return err
	}
	resultV.Elem().Set(sliceV)
	return nil
}

func (r *Repository) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := r.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := r.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) InsertOne(ctx context.Context, entity interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := r.collection.InsertOne(ctx, entity, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) FindOneAndUpdate(ctx context.Context, entity interface{}, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	one := r.collection.FindOneAndUpdate(ctx, filter, update, opts...)
	err := one.Decode(entity)
	return err
}

func (r *Repository) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := r.collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := r.collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
