package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"gustavonovaes.dev/go-sjc-assist-bot/internal/config"
)

var (
	client    *mongo.Client
	appConfig config.Config
)

func init() {
	appConfig = config.New()

	var err error
	log.Println("INFO: Connecting to MongoDB")
	client, err = Connect()
	if err != nil {
		log.Fatalf("ERROR: Fail to connect to MongoDB: %v", err)
	}
	log.Println("INFO: MongoDB connected successfully")
}

func Connect() (*mongo.Client, error) {
	return mongo.Connect(context.Background(), options.Client().ApplyURI(appConfig.MONGODB_URI))
}

func Close() error {
	return client.Disconnect(context.Background())
}

func GetClient() *mongo.Client {
	if client == nil {
		client, _ = Connect()
	}

	return client
}

func GetCollection(name string) *mongo.Collection {
	uri, err := connstring.ParseAndValidate(appConfig.MONGODB_URI)
	if err != nil {
		log.Fatalf("ERROR: Invalid MongoDB URI: %v", err)
	}

	return client.Database(uri.Database).Collection(name)
}

func SaveCollection(name string, data interface{}) error {
	collection := GetCollection(name)
	_, err := collection.InsertOne(context.Background(), data)
	return err
}

func FindCollection(name string, filter interface{}) (*mongo.Cursor, error) {
	collection := GetCollection(name)
	return collection.Find(context.Background(), filter)
}
