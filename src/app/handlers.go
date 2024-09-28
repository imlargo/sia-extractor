package app

import (
	"encoding/json"
	"os"
	"sia-extractor/src/core"
	"strconv"
	"sync"
)

func testExtraccion(indexGrupo int) {
	var listadoGrupos [][]map[string]string
	contentGrupos, _ := os.ReadFile(core.Path_Grupos)
	json.Unmarshal(contentGrupos, &listadoGrupos)
	println("Accediendo a grupo: ", indexGrupo+1)
	var grupoAsignado []map[string]string = listadoGrupos[indexGrupo+1]

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

			core.GetAsignaturasCarrera(codigo)

		}(carrera)
	}

	wg.Wait()
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
