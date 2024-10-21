package deploy

import (
	"fmt"
	"sia-extractor/src/core"
	"sia-extractor/src/utils"
	"slices"
)

func getTipologiasUnicas(asignaturas []core.Asignatura) []string {
	tipologiasUnicas := make([]string, 0)

	for _, asignatura := range asignaturas {

		if slices.Contains(tipologiasUnicas, asignatura.Tipologia) {
			continue
		}

		tipologiasUnicas = append(tipologiasUnicas, asignatura.Tipologia)
	}

	return tipologiasUnicas
}

func CreateDocumentCarrera(carrera string, facultad string, asignaturas []core.Asignatura) DocumentCarrera {
	return DocumentCarrera{
		ID:          carrera,
		Facultad:    facultad,
		Carrera:     carrera,
		Asignaturas: asignaturas,
	}
}

func MergeDataSede() map[string]map[string][]core.Asignatura {
	// Cargar listado de carreras
	var carreras []map[string]string
	utils.LoadJsonFromFile(&carreras, core.Path_Carreras)

	println("Cantidad de grupos: ", totalGrupos)

	dataAsignaturas := make(map[string][]core.Asignatura)

	for i := 0; i < totalGrupos; i++ {
		// Cargar datos de asignaturas de carreras
		path := fmt.Sprintf("%s%d.json", pathToData, i+1)

		// Leer datos de asignaturas
		var data map[string][]core.Asignatura
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
			valueCarrera := carrera["carrera"]

			if len(dataAsignaturas[valueCarrera]) == 0 {
				panic("No se encontraron datos para la carrera: " + valueCarrera)
			}

			dataFacultad[valueCarrera] = dataAsignaturas[valueCarrera]
		}

		merged[facultad] = dataFacultad
	}

	return merged
}
