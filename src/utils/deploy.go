package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sia-extractor/src/core"
	"slices"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentCarrera struct {
	ID          string            `json:"_id" bson:"_id"`
	Facultad    string            `json:"facultad" bson:"facultad"`
	Carrera     string            `json:"carrera" bson:"carrera"`
	Asignaturas []core.Asignatura `json:"asignaturas" bson:"asignaturas"`
}

const (
	totalGrupos int = 39
)

func DeployData() {

	var carreras []map[string]string
	bytes, _ := os.ReadFile(core.Path_Carreras)
	json.Unmarshal(bytes, &carreras)

	println("Cantidad de grupos: ", totalGrupos)

	var dataAsignaturas = make(map[string][]core.Asignatura)

	for i := 0; i < totalGrupos; i++ {
		var path string = "./" + strconv.Itoa(i+1) + ".json"
		var data map[string][]core.Asignatura

		// unmarshall json
		bytes, _ := os.ReadFile(path)
		json.Unmarshal(bytes, &data)

		for carrera, asignaturas := range data {
			dataAsignaturas[carrera] = asignaturas
		}
	}

	carrerasAgrupadas := groupBy(carreras, func(carrera map[string]string) string {
		return carrera["facultad"]
	})

	var merged = make(map[string]map[string][]core.Asignatura)
	for facultad, carreras := range carrerasAgrupadas {

		var dataFacultad = make(map[string][]core.Asignatura)
		for _, carrera := range carreras {
			var valueCarrera string = carrera["carrera"]
			dataFacultad[valueCarrera] = dataAsignaturas[valueCarrera]
		}

		merged[facultad] = dataFacultad

	}

	dataMerged, _ := json.Marshal(merged)
	os.WriteFile("data.json", dataMerged, 0644)

	saveInDatabase(&merged)

}

func groupBy[T any](array []map[string]T, function func(map[string]T) string) map[string][]map[string]T {

	result := make(map[string][]map[string]T)

	for _, item := range array {
		key := function(item)
		result[key] = append(result[key], item)
	}

	return result
}

func saveInDatabase(data *map[string]map[string][]core.Asignatura) {
	var uri string = os.Getenv("MONGO_URI")

	if uri == "" {
		println("No se ha definido la variable de entorno MONGO_URI")
		return
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collFacultades := client.Database("asignaturas").Collection("asignaturas")
	collCarreras := client.Database("asignaturas").Collection("carreras")
	collConfig := client.Database("asignaturas").Collection("config")

	fmt.Println("Connected to MongoDB!")

	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	listado := make(map[string]map[string][]string)

	for facultad, carreras := range *data {

		// crear el listado con datos

		listado[facultad] = map[string][]string{}
		for carrera, asignaturas := range carreras {

			tipologias := make([]string, 0)

			for _, asignatura := range asignaturas {

				if slices.Contains(tipologias, asignatura.Tipologia) {
					continue
				}

				tipologias = append(tipologias, asignatura.Tipologia)
			}

			listado[facultad][carrera] = tipologias
		}

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

			document := DocumentCarrera{
				ID:          carrera,
				Facultad:    facultad,
				Carrera:     carrera,
				Asignaturas: asignaturas,
			}

			go func(document DocumentCarrera) {
				defer wg2.Done()

				query := bson.D{{Key: "_id", Value: document.ID}}

				_, err := collCarreras.ReplaceOne(context.TODO(), query, document)

				if err != nil {
					panic(err)
				}

				fmt.Println("Carrera actualizada: ", document.Carrera)

			}(document)

		}
	}

	query := bson.D{{Key: "_id", Value: "listado"}}
	collConfig.ReplaceOne(context.TODO(), query, listado)

	wg.Wait()
	wg2.Wait()

	println("Datos actualizados con exito")

}
