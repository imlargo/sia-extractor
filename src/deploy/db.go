package deploy

import (
	"context"
	"os"
	"sia-extractor/src/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseClient struct {
	Client *mongo.Client
}

func NewDatabaseClient() *DatabaseClient {
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		panic("No se ha definido la variable de entorno MONGO_URI")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &DatabaseClient{
		Client: client,
	}

}

func (db *DatabaseClient) Disconnect() {
	if err := db.Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (db *DatabaseClient) SaveCarrera(document DocumentCarrera, collCarreras *mongo.Collection) error {
	query := bson.D{{Key: "_id", Value: document.ID}}

	_, err := collCarreras.ReplaceOne(context.TODO(), query, document)

	return err
}

func (db *DatabaseClient) updateFechaExtraccion(lastUpdate string) {
	collConfig := db.Client.Database("asignaturas").Collection("config")
	query := bson.D{{Key: "_id", Value: "metadata"}}

	metadata := map[string]string{
		"lastUpdated": lastUpdate,
	}

	collConfig.ReplaceOne(context.TODO(), query, metadata)
}

func (db *DatabaseClient) updateListadoCarreras(data *map[string]map[string][]core.Asignatura) {

	listado := make(map[string]map[string][]string)

	for facultad, carreras := range *data {
		// crear el listado con datos
		listado[facultad] = map[string][]string{}
		for carrera, asignaturas := range carreras {
			tipologiasUnicas := getTipologiasUnicas(asignaturas)
			listado[facultad][carrera] = tipologiasUnicas
		}
	}

	collConfig := db.Client.Database("asignaturas").Collection("config")
	query := bson.D{{Key: "_id", Value: "listado"}}
	collConfig.ReplaceOne(context.TODO(), query, listado)
}
