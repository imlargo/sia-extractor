package deploy

import (
	"context"
	"fmt"
	"os"
	"sia-extractor/src/core"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName           string = "asignaturas"
	pathCollAll      string = "asignaturas"
	pathCollCarreras string = "carreras"
	pathCollConfig   string = "config"
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

func (db *DatabaseClient) SaveCarrera(document DocumentCarrera) error {

	collCarreras := db.Client.Database(dbName).Collection(pathCollCarreras)
	query := bson.D{{Key: "_id", Value: document.ID}}

	_, err := collCarreras.ReplaceOne(context.TODO(), query, document)

	return err
}

func (db *DatabaseClient) UpdateFechaExtraccion(lastUpdate string) {
	collConfig := db.Client.Database(dbName).Collection(pathCollConfig)
	query := bson.D{{Key: "_id", Value: "metadata"}}

	metadata := map[string]string{
		"lastUpdated": lastUpdate,
	}

	collConfig.ReplaceOne(context.TODO(), query, metadata)
}

func (db *DatabaseClient) UpdateListadoCarreras(data *map[string]map[string][]core.Asignatura) {

	listado := make(map[string]map[string][]string)

	for facultad, carreras := range *data {
		// crear el listado con datos
		listado[facultad] = map[string][]string{}
		for carrera, asignaturas := range carreras {
			tipologiasUnicas := getTipologiasUnicas(asignaturas)
			listado[facultad][carrera] = tipologiasUnicas
		}
	}

	collConfig := db.Client.Database(dbName).Collection(pathCollConfig)
	query := bson.D{{Key: "_id", Value: "listado"}}
	collConfig.ReplaceOne(context.TODO(), query, listado)
}

func (db *DatabaseClient) SaveDataCarreras(data *map[string]map[string][]core.Asignatura) {

	var wg sync.WaitGroup

	for facultad, carreras := range *data {
		for carrera, asignaturas := range carreras {
			wg.Add(1)

			go func(carrera string, facultad string, asignaturas []core.Asignatura) {
				defer wg.Done()

				document := CreateDocumentCarrera(carrera, facultad, asignaturas)
				if err := db.SaveCarrera(document); err != nil {
					println("Error al guardar carrera: ", err)
				}

				fmt.Println("Carrera actualizada: ", document.Carrera)

			}(carrera, facultad, asignaturas)
		}
	}

	wg.Wait()

	println("Datos actualizados con exito")
}

func (db *DatabaseClient) SaveDataSede(data *map[string]map[string][]core.Asignatura) {

	collFacultades := db.Client.Database(dbName).Collection(pathCollAll)

	for facultad, carreras := range *data {

		query := bson.D{{Key: "_id", Value: facultad}}
		_, err := collFacultades.ReplaceOne(context.TODO(), query, carreras)
		if err != nil {
			panic(err)
		}
		fmt.Println("Facultad actualizada: ", facultad)
	}

}
