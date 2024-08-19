package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sia-extractor/src/core"
	"sia-extractor/src/utils"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// core.CreatePathsCarreras()

	var args []string = os.Args[1:]

	if len(args) == 0 {
		println("Debe ingresar los argumentos")
		return
	}

	if args[0] == "paths" {
		println("Creando paths")
		core.CreatePathsCarreras()
		return
	}

	if args[0] == "mongo" {
		saveInDatabase()
		return
	}

	if args[0] == "group" {
		println("Agrupando carreras")
		core.GenerarGruposCarreras()
		return
	}

	if args[0] == "merge" {
		println("Consolidando datos")
		utils.MergeData()
		return
	}

	grupoAsignado, _ := strconv.Atoi(args[0])

	println("Grupo asignado: ", grupoAsignado)

	initTime := time.Now()

	extraerTodo(grupoAsignado - 1)
	println("")
	println("......................................................")
	fmt.Printf("Tiempo de ejecuciónnnnnnnnnnn final: %v\n", time.Since(initTime))

	// core.GenerarGruposCarreras()

}

func extraerTodo(indexGrupo int) {

	var listadoGrupos [][]map[string]string
	contentGrupos, _ := os.ReadFile(core.Path_Grupos)
	json.Unmarshal(contentGrupos, &listadoGrupos)
	var grupoAsignado []map[string]string = listadoGrupos[indexGrupo]

	asignaturasChan := make(chan []core.Asignatura, len(grupoAsignado))

	var wg sync.WaitGroup
	for _, carrera := range grupoAsignado {

		wg.Add(1)

		go func(carrera map[string]string) {
			defer wg.Done()

			codigo := core.Codigo{
				Nivel:     core.ValueNivel,
				Sede:      core.ValueSede,
				Facultad:  carrera["facultad"],
				Carrera:   carrera["carrera"],
				Tipologia: core.Tipologia_All,
			}

			println("INICIANDOOOO: ", codigo.Carrera)

			var asignaturas []core.Asignatura = core.GetAsignaturasCarrera(codigo)

			asignaturasChan <- asignaturas

		}(carrera)
	}

	go func() {
		wg.Wait()
		close(asignaturasChan)
	}()

	var finalAsignaturas map[string][]core.Asignatura = make(map[string][]core.Asignatura)
	for asignaturas := range asignaturasChan {
		var carrera string = asignaturas[0].Carrera
		finalAsignaturas[carrera] = asignaturas
	}

	filename := strconv.Itoa(indexGrupo+1) + ".json"
	finalAsignaturasJSON, _ := json.Marshal(finalAsignaturas)
	os.WriteFile(filename, finalAsignaturasJSON, 0644)
}

func saveInDatabase() {
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

	coll := client.Database("asignaturas").Collection("asignaturas")

	fmt.Println("Connected to MongoDB!", coll.Name())
}
