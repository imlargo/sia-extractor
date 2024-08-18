package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sia-extractor/src/core"
	"strconv"
	"sync"
	"time"
)

func main() {

	// core.CreatePathsCarreras()

	var args []string = os.Args[1:]

	if len(args) == 0 {
		println("Debe ingresar los argumentos")
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
	contentGrupos, _ := os.ReadFile("data/grupos.json")
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

	var finalAsignaturas []core.Asignatura
	for asignaturas := range asignaturasChan {
		finalAsignaturas = append(finalAsignaturas, asignaturas...)
	}

	filename := strconv.Itoa(indexGrupo+1) + ".json"
	finalAsignaturasJSON, _ := json.Marshal(finalAsignaturas)
	os.WriteFile(filename, finalAsignaturasJSON, 0644)
}