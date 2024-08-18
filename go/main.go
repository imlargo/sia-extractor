package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sia-extractor/core"
	"sync"
	"time"
)

const (
	cantidad_por_grupo int = 4
)

func main() {

	// core.CreatePathsCarreras()

	fmt.Println(len(os.Args), os.Args)

	initTime := time.Now()

	extraerTodo()

	println("")
	println("......................................................")
	fmt.Printf("Tiempo de ejecuciónnnnnnnnnnn final: %v\n", time.Since(initTime))

	// group()

}

func group() {
	var listadoCarreras []map[string]string

	contentCarreras, _ := os.ReadFile(core.Path_Carreras)
	json.Unmarshal(contentCarreras, &listadoCarreras)

	stacks := int(math.Ceil(float64(len(listadoCarreras)) / float64(cantidad_por_grupo)))

	println("Cantidad de stacks: ", (stacks))

	var grupos [][]map[string]string
	for i := 0; i < stacks; i++ {
		var grupo []map[string]string

		for j := 0; j < cantidad_por_grupo; j++ {
			if (i*cantidad_por_grupo)+j < len(listadoCarreras) {
				grupo = append(grupo, listadoCarreras[(i*cantidad_por_grupo)+j])
			}
		}

		grupos = append(grupos, grupo)
	}

	dataGruposJSON, _ := json.Marshal(grupos)
	os.WriteFile("grupos.json", dataGruposJSON, 0644)

}

func extraerTodo() {
	var listadoCarreras []map[string]string

	contentCarreras, _ := os.ReadFile(core.Path_Carreras)
	json.Unmarshal(contentCarreras, &listadoCarreras)

	var wg sync.WaitGroup
	for _, carrera := range listadoCarreras[0:4] {

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
