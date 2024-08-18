package main

import (
	"encoding/json"
	"os"
	"sia-extractor/core"
	"sync"
	"time"
)

func main() {

	// core.CreatePathsCarreras()

	initTime := time.Now()

	extraerTodo()

	println("")
	println("......................................................")
	println("Tiempo total de ejecución: ", time.Since(initTime))

}

func extraerTodo() {
	var listadoCarreras []map[string]string

	contentCarreras, _ := os.ReadFile("carreras.json")
	json.Unmarshal(contentCarreras, &listadoCarreras)

	var wg sync.WaitGroup
	for _, carrera := range listadoCarreras[0:3] {

		wg.Add(1)

		go func(carrera map[string]string) {
			defer wg.Done()

			println("--------------------- INICIANDOOOO ", carrera["carrera"], "---------------------")

			codigo := core.Codigo{
				Nivel:     core.ValueNivel,
				Sede:      core.ValueSede,
				Facultad:  carrera["facultad"],
				Carrera:   carrera["carrera"],
				Tipologia: core.Tipologia_All,
			}

			var asignaturas []core.Asignatura = core.GetAsignaturasCarrera(codigo)
			dataAsignaturasJSON, _ := json.Marshal(asignaturas)
			var filename string = codigo.Carrera + ".json"
			os.WriteFile(filename, dataAsignaturasJSON, 0644)

		}(carrera)
	}

	wg.Wait()

}
