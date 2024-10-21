package deploy

import (
	"context"
	"os"
	"sia-extractor/src/core"
	"slices"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoDbClient() *mongo.Client {
	var uri string = os.Getenv("MONGO_URI")

	if uri == "" {
		panic("No se ha definido la variable de entorno MONGO_URI")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func getTipologiasUnicas(asignaturas []core.Asignatura) []string {
	tipologiasUnicas := make([]string, 0)

	for _, asignatura := range asignaturas {

		if slices.Contains(tipologiasUnicas, asignatura.Tipologia) {
			continue
		}

		tipologiasUnicas = append(tipologiasUnicas, asignatura.Tipologia)
	}

	return tipologiasUnicas
}
