package core

import (
	"encoding/json"
	"os"
	"sync"
)

func ExtraerGrupo(indexGrupo int) map[string][]Asignatura {

	var listadoGrupos [][]map[string]string
	bytesGrupos, _ := os.ReadFile(Path_Grupos)
	json.Unmarshal(bytesGrupos, &listadoGrupos)
	grupo := listadoGrupos[indexGrupo]

	chanAsignaturas := make(chan []Asignatura, len(grupo))

	var wg sync.WaitGroup
	for _, carrera := range grupo {

		wg.Add(1)

		go func(carrera map[string]string) {
			defer wg.Done()

			codigo := Codigo{
				Nivel:     ValueNivel,
				Sede:      ValueSede,
				Facultad:  carrera["facultad"],
				Carrera:   carrera["carrera"],
				Tipologia: Tipologia_All,
			}

			println("INICIANDOOOO: ", codigo.Carrera)

			var asignaturas []Asignatura = GetAsignaturasCarrera(codigo)

			chanAsignaturas <- asignaturas

		}(carrera)
	}

	go func() {
		wg.Wait()
		close(chanAsignaturas)
	}()

	data := make(map[string][]Asignatura)
	for asignaturas := range chanAsignaturas {
		carrera := asignaturas[0].Carrera
		data[carrera] = asignaturas
	}

	return data
}
