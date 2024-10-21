package app

import (
	"encoding/json"
	"fmt"
	"os"
	"sia-extractor/src/core"
	"sia-extractor/src/deploy"
	"sia-extractor/src/utils"
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
		electivas := core.ExtraerElectivas(core.ConstructCodigo("3068 FACULTAD DE MINAS", "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA"))
		err := utils.SaveJsonToFile(electivas, "electivas.json")
		if err != nil {
			fmt.Println("Error al guardar archivo: ", err)
		}

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

		core.GetAsignaturasCarrera(core.ConstructCodigo(carrera["facultad"], carrera["carrera"]))
		fmt.Printf("Tiempo de ejecución final: %v\n", time.Since(initTime))
	case "extract":
		grupo, _ := strconv.Atoi(args[1])
		initTime := time.Now()

		data := ExtractCarrera(grupo)

		if data == nil {
			println("Grupo no encontrado")
			return
		}

		filename := strconv.Itoa(grupo) + ".json"
		err := utils.SaveJsonToFile(data, filename)
		if err != nil {
			fmt.Println("Error al guardar archivo: ", err)
		}
		println("")
		println("......................................................")
		fmt.Printf("Tiempo de ejecución final: %v\n", time.Since(initTime))
	default:
		fmt.Println("Comando no reconocido")
	}
}

func ExtractCarrera(indexGrupo int) map[string]*[]core.Asignatura {

	codigo := core.GetCodigoFromGrupo(indexGrupo)

	if codigo == nil {
		return nil
	}

	var data *[]core.Asignatura

	println("Iniciando: ", codigo.Carrera)
	if codigo.Facultad == core.ValuesElectiva.FacultadPor {
		data = core.ExtraerElectivas(*codigo)
	} else {
		data = core.GetAsignaturasCarrera(*codigo)
	}
	println("Finalizado: ", codigo.Carrera)

	return map[string]*[]core.Asignatura{codigo.Carrera: data}
}
