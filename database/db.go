package database

import (
	"context"

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

//var m *MongoDB

func GetMongoDB(opt *options.ClientOptions) *MongoDB {
	//if m != nil {
	//	return m
	//}

	m := new(MongoDB)
	m.logger = log.GetLogger()
	m.client = m.newClient(context.Background(), opt)

	return m
}

func (m *MongoDB) newClient(ctx context.Context, opt *options.ClientOptions) *mongo.Client {
	m.logger.Infoln("MongoDB is starting to connect...")

	//uri := config.GetConfig().GetString("mongo")
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(20))
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		m.logger.Panicf("MongoDB connection failed, Err: %s, uri: %s", err.Error(), opt.GetURI())
	} else {
		m.logger.Infof("MongoDB connection succeed...")
	}
	return client
}

func (m *MongoDB) Ping(ctx context.Context) {
	err := m.client.Ping(ctx, readpref.SecondaryPreferred())
	if err != nil {
		m.logger.Panicf("mongoDB ping Err: %s", err.Error())
	} else {
		m.logger.Info("mongoDB ping succeed...")
	}
}

func (m *MongoDB) GetDatabase(name ...string) *mongo.Database {
	if m.db == nil {
		m.db = m.client.Database("guoguo")
		if len(name) > 0 {
			m.db = m.client.Database(name[0])
		}
	}
	return m.db
}

func (m *MongoDB) GetCollection(col string) *mongo.Collection {
	if m.db == nil {
		m.db = m.GetDatabase()
	}
	return m.db.Collection(col)
}

func (m *MongoDB) Disconnect() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		m.logger.Panicf("mongoDB disconnect Err: %s", err.Error())
	} else {
		m.logger.Info("mongoDB disconnect succeed...")
	}
}
