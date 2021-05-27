package crud

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Resource struct {
	DB     *mongo.Database
	Client *mongo.Client
}

// 根据环境变量快速初始化数据库链接
func SimpleNewDB() (*Resource, error) {
	return NewDB(os.Getenv("MONGO_URL"), os.Getenv("MONGO_DBNAME"))
}

// 初始化数据库链接
func NewDB(uri string, databaseName string) (*Resource, error) {

	if uri == "" {
		return nil, errors.New("mongo uri is empty")
	}
	if databaseName == "" {
		return nil, errors.New("mongo databaseName is empty")
	}

	// Replace the uri string with your MongoDB deployment's connection string.
	// uri := os.Getenv("MONGO_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctxping, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctxping, readpref.Primary())
	if err != nil {
		return nil, err
	}
	// fmt.Println(client)
	// databaseName :=os.Getenv("MONGO_DBNAME")
	db := client.Database(databaseName)

	return &Resource{
		DB:     db,
		Client: client,
	}, nil
}
