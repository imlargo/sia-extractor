package app

import (
	"encoding/json"
	"fmt"
	"os"
	"sia-extractor/src/core"
	"sia-extractor/src/deploy"
	"strconv"
	"time"
)

func App() {

	args := os.Args[1:]
	if len(args) == 0 {
		println("Debe ingresar los argumentos")
		return
	}

	tipo := args[0]

	switch tipo {
	case "paths":
		println("Creando paths")
		core.CreatePathsCarreras()
	case "group":
		println("Agrupando carreras")
		core.GenerarGruposCarreras()
	case "electivas":
		println("Electivas")
		electivas := core.ExtraerElectivas()
		electivasJSON, _ := json.Marshal(electivas)
		os.WriteFile("electivas.json", electivasJSON, 0644)
	case "deploy":
		println("Consolidando datos")
		deploy.DeployData()
	case "test":
		println("Iniciando test")
		grupo, _ := strconv.Atoi(args[1])
		initTime := time.Now()
		println("Grupo asignado: ", grupo)

		var listadoGrupos [][]map[string]string
		bytesGrupos, _ := os.ReadFile(core.Path_Grupos)
		json.Unmarshal(bytesGrupos, &listadoGrupos)
		carrera := listadoGrupos[grupo-1][0]

		core.GetAsignaturasCarrera(core.Codigo{
			Nivel:     core.ValueNivel,
			Sede:      core.ValueSede,
			Facultad:  carrera["facultad"],
			Carrera:   carrera["carrera"],
			Tipologia: core.Tipologia_All,
		})
		fmt.Printf("Tiempo de ejecución final: %v\n", time.Since(initTime))
	case "extract":
		grupo, _ := strconv.Atoi(args[1])
		println("Grupo asignado: ", grupo)
		initTime := time.Now()
		data := core.ExtraerGrupo(grupo)
		if data == nil {
			println("Grupo no encontrado")
			return
		}
		filename := strconv.Itoa(grupo) + ".json"
		finalAsignaturasJSON, _ := json.Marshal(data)
		os.WriteFile(filename, finalAsignaturasJSON, 0644)
		println("")
		println("......................................................")
		fmt.Printf("Tiempo de ejecuciónnnnnnnnnnn final: %v\n", time.Since(initTime))
	default:
		fmt.Println("Comando no reconocido")
	}
}
