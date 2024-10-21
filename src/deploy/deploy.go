package deploy

import (
	"context"
	"fmt"
	"sia-extractor/src/core"
	"sia-extractor/src/utils"
	"strconv"
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

	merged := mergeAllData()
	err := utils.SaveJsonToFile(merged, "data.json")
	if err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}

	dbClient.updateListadoCarreras(&merged)
	saveInDatabase(dbClient, &merged)
	dbClient.updateFechaExtraccion(merged["3068 FACULTAD DE MINAS"]["3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA"][0].FechaExtraccion)

}

func mergeAllData() map[string]map[string][]core.Asignatura {
	// Cargar listado de carreras
	var carreras []map[string]string
	utils.LoadJsonFromFile(&carreras, core.Path_Carreras)

	println("Cantidad de grupos: ", totalGrupos)

	var dataAsignaturas = make(map[string][]core.Asignatura)

	for i := 0; i < totalGrupos; i++ {
		// Cargar datos de asignaturas de carreras
		path := pathToData + strconv.Itoa(i+1) + ".json"
		var data map[string][]core.Asignatura

		// Leer datos de asignaturas
		utils.LoadJsonFromFile(&data, path)

		// Agregar asignaturas a consolidado
		for carrera, asignaturas := range data {
			dataAsignaturas[carrera] = asignaturas
		}
	}

	// Agrupar carreras por facultad
	carrerasAgrupadas := utils.GroupBy(carreras, func(carrera map[string]string) string {
		return carrera["facultad"]
	})

	merged := make(map[string]map[string][]core.Asignatura)
	for facultad, carreras := range carrerasAgrupadas {

		dataFacultad := make(map[string][]core.Asignatura)
		for _, carrera := range carreras {
			var valueCarrera string = carrera["carrera"]

			if len(dataAsignaturas[valueCarrera]) == 0 {
				panic("No se encontraron datos para la carrera: " + valueCarrera)
			}

			dataFacultad[valueCarrera] = dataAsignaturas[valueCarrera]
		}

		merged[facultad] = dataFacultad

	}

	return merged
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
