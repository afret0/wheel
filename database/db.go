package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/afret0/wheel/log"
)

type MongoDB struct {
	logger *logrus.Logger
	client *mongo.Client
	db     *mongo.Database
}

// Deprecated: GetMongoDB 方法已废弃，请使用 NewMongoDB 直接创建实例。
func GetMongoDB(opt *options.ClientOptions, database string) *MongoDB {
	return NewMongoDB(opt, database)
}

func NewMongoDB(opt *options.ClientOptions, database string) *MongoDB {
	//if m != nil {
	//	return m
	//}

	m := new(MongoDB)
	m.logger = log.GetLogger()
	m.client = m.newClient(context.Background(), opt)
	m.db = m.client.Database(database)

	return m
}

func (m *MongoDB) newClient(ctx context.Context, opt *options.ClientOptions) *mongo.Client {
	m.logger.Infoln("MongoDB is starting to connect...")

	//uri := config.GetConfig().GetString("mongo")
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(20))
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		//m.logger.Panicf("MongoDB connection failed, Err: %s, uri: %s", err.Error(), opt.GetURI())
		panic(fmt.Sprintf("MongoDB connection failed, Err: %s, uri: %s", err.Error(), opt.GetURI()))
	} else {
		m.logger.Infof("MongoDB connection success...")
	}
	return client
}

func (m *MongoDB) Ping(ctx context.Context) {
	err := m.client.Ping(ctx, readpref.SecondaryPreferred())
	if err != nil {
		//m.logger.Panicf("mongoDB ping Err: %s", err.Error())
		panic(err)
	} else {
		m.logger.Info("mongoDB ping success...")
	}
}

func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.db
}

func (m *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	return m.db.Collection(collectionName)
}

func (m *MongoDB) GetClient() *mongo.Client {
	return m.client
}

// func (m *MongoDB) GetDatabase(name ...string) *mongo.Database {
// 	if m.db == nil {
// 		m.db = m.client.Database("guoguo")
// 		if len(name) > 0 {
// 			m.db = m.client.Database(name[0])
// 		}
// 	}
// 	return m.db
// }

// func (m *MongoDB) GetCollection(col string) *mongo.Collection {
// 	if m.db == nil {
// 		m.db = m.GetDatabase()
// 	}
// 	return m.db.Collection(col)
// }

func (m *MongoDB) Disconnect() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		//m.logger.Panicf("mongoDB disconnect Err: %s", err.Error())
		panic(fmt.Sprintf("mongoDB disconnect Err: %s", err.Error()))
	} else {
		m.logger.Info("mongoDB disconnect success...")
	}
}

func (m *MongoDB) CheckCollectionExist(CollectionName string) {
	collectionL, err := m.GetDatabase().ListCollectionNames(context.Background(), bson.M{"name": CollectionName})
	if err != nil {
		log.GetLogger().Panic(err)
	}
	if len(collectionL) > 0 {
		log.GetLogger().Infof("collection %s exists", CollectionName)
	} else {
		log.GetLogger().Infof("collection %s not exists, 开始创建", CollectionName)
		err = m.GetDatabase().CreateCollection(context.Background(), CollectionName)
		if err != nil {
			log.GetLogger().Panic(err)
		}
		log.GetLogger().Infof("collection %s 创建成功", CollectionName)
	}
}
