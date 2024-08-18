package utils

import (
	"encoding/json"
	"os"
	"sia-extractor/src/core"
	"strconv"
)

const (
	totalGrupos int = 13
)

func MergeData() {

	var carreras []map[string]string
	bytes, _ := os.ReadFile(core.Path_Carreras)
	json.Unmarshal(bytes, &carreras)

	println("Cantidad de grupos: ", totalGrupos)

	var dataAsignaturas = make(map[string][]core.Asignatura)

	for i := 0; i < totalGrupos; i++ {
		var path string = "artifacts/" + strconv.Itoa(i+1) + ".json"
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

}

func groupBy[T any](array []map[string]T, function func(map[string]T) string) map[string][]map[string]T {

	result := make(map[string][]map[string]T)

	for _, item := range array {
		key := function(item)
		result[key] = append(result[key], item)
	}

	return result
}
