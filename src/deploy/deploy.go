package deploy

import (
	"context"
	"fmt"
	"sia-extractor/src/core"
	"sia-extractor/src/utils"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

const pathToData string = "artifacts/"

type DocumentCarrera struct {
	ID          string            `json:"_id" bson:"_id"`
	Facultad    string            `json:"facultad" bson:"facultad"`
	Carrera     string            `json:"carrera" bson:"carrera"`
	Asignaturas []core.Asignatura `json:"asignaturas" bson:"asignaturas"`
}

const (
	totalGrupos int = 40
)

func DeployData() {

	dbClient := NewDatabaseClient()
	fmt.Println("Connected to MongoDB!")

	defer dbClient.Disconnect()

	merged := MergeDataSede()
	err := utils.SaveJsonToFile(merged, "data.json")
	if err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}

	dbClient.updateListadoCarreras(&merged)
	saveInDatabase(dbClient, &merged)
	dbClient.updateFechaExtraccion(merged["3068 FACULTAD DE MINAS"]["3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA"][0].FechaExtraccion)

}

func saveInDatabase(dbClient *DatabaseClient, data *map[string]map[string][]core.Asignatura) {

	collFacultades := dbClient.Client.Database("asignaturas").Collection("asignaturas")
	collCarreras := dbClient.Client.Database("asignaturas").Collection("carreras")

	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	for facultad, carreras := range *data {

		wg.Add(1)

		go func(facultad string, carreras map[string][]core.Asignatura) {
			defer wg.Done()

			query := bson.D{{Key: "_id", Value: facultad}}

			_, err := collFacultades.ReplaceOne(context.TODO(), query, carreras)

			if err != nil {
				panic(err)
			}

			fmt.Println("Facultad actualizada: ", facultad)
		}(facultad, carreras)

		for carrera, asignaturas := range carreras {
			wg2.Add(1)

			document := CreateDocumentCarrera(carrera, facultad, asignaturas)

			go func(document DocumentCarrera) {
				defer wg2.Done()

				err := dbClient.SaveCarrera(document, collCarreras)

				if err != nil {
					panic(err)
				}

				fmt.Println("Carrera actualizada: ", document.Carrera)

			}(document)

		}
	}

	wg.Wait()
	wg2.Wait()

	println("Datos actualizados con exito")

}
