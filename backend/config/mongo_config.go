package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
init loads environment variables from a .env file
*/
func init() {
	_ = godotenv.Load()
}

type DatabaseConfig struct {
	Client  *mongo.Client
	DbName  string
}

/*
mongoClient is the MongoDB client instance
*/
var mongoClient *DatabaseConfig

/*
ConnectToMongo establishes a connection to the MongoDB database grabbing the mongo data from enviromental variables
*/
func ConnectToMongo() error {
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	dbName := os.Getenv("MONGO_DB")

	if user == "" || password == "" || host == "" || dbName == "" {
		return fmt.Errorf("faltan variables de entorno: MONGO_USER, MONGO_PASSWORD, MONGO_HOST o MONGO_DB")
	}

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", user, password, host, dbName)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)

	mongoClient = &DatabaseConfig{
		Client: client,
		DbName: dbName,
	}

	return err
}

/*
 Get MongoClient is a helper function that helps to get the MongoDB client instance	
*/
func GetMongoClient() *DatabaseConfig {
	if mongoClient == nil {
		if err := ConnectToMongo(); err != nil {
			panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
		}
	}
	return mongoClient

}

/*
CloseMongoConnection closes the MongoDB connection
*/
func CloseMongoConnection() error {
	if mongoClient != nil && mongoClient.Client != nil {
		return mongoClient.Client.Disconnect(context.Background())
	}	
	return nil
}
