package deploy

import (
	"context"
	"fmt"
	"sia-extractor/src/core"
	"sia-extractor/src/utils"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	client := getMongoDbClient()
	fmt.Println("Connected to MongoDB!")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	merged := mergeAllData()
	err := utils.SaveJsonToFile(merged, "data.json")
	if err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}

	updateListadoCarreras(client, &merged)
	saveInDatabase(client, &merged)
	updateFechaExtraccion(client, merged["3068 FACULTAD DE MINAS"]["3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA"][0].FechaExtraccion)

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
	carrerasAgrupadas := groupBy(carreras, func(carrera map[string]string) string {
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

func updateFechaExtraccion(client *mongo.Client, lastUpdate string) {
	collConfig := client.Database("asignaturas").Collection("config")
	query := bson.D{{Key: "_id", Value: "metadata"}}

	metadata := map[string]string{
		"lastUpdated": lastUpdate,
	}

	collConfig.ReplaceOne(context.TODO(), query, metadata)
}

func updateListadoCarreras(client *mongo.Client, data *map[string]map[string][]core.Asignatura) {

	listado := make(map[string]map[string][]string)

	for facultad, carreras := range *data {
		// crear el listado con datos
		listado[facultad] = map[string][]string{}
		for carrera, asignaturas := range carreras {
			tipologiasUnicas := getTipologiasUnicas(asignaturas)
			listado[facultad][carrera] = tipologiasUnicas
		}
	}

	collConfig := client.Database("asignaturas").Collection("config")
	query := bson.D{{Key: "_id", Value: "listado"}}
	collConfig.ReplaceOne(context.TODO(), query, listado)
}

func saveInDatabase(client *mongo.Client, data *map[string]map[string][]core.Asignatura) {

	collFacultades := client.Database("asignaturas").Collection("asignaturas")
	collCarreras := client.Database("asignaturas").Collection("carreras")

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

	wg.Wait()
	wg2.Wait()

	println("Datos actualizados con exito")

}
