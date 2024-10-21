package deploy

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoDbClient() *mongo.Client {
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		panic("No se ha definido la variable de entorno MONGO_URI")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}
