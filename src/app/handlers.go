package app

import (
	"fmt"
	"sia-extractor/src/core"
	"sia-extractor/src/deploy"
	"sia-extractor/src/utils"
	"strconv"
	"time"
)

func handlePaths(args []string) {
	fmt.Println("Creando paths")

	extractor := core.NewExtractor()
	extractor.CreatePathsCarreras()
}

func handleGroup(args []string) {
	fmt.Println("Agrupando carreras")
	extractor := core.NewExtractor()
	extractor.GenerarGruposCarreras()
}

func handleElectivas(args []string) {
	fmt.Println("Electivas")
	extractor := core.NewExtractor()
	electivas := extractor.ExtraerElectivas(core.ConstructCodigo("3068 FACULTAD DE MINAS", "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA"))
	if err := utils.SaveJsonToFile(electivas, "electivas.json"); err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}
}

func handleDeploy(args []string) {
	fmt.Println("Consolidando datos")
	deploy.DeployData()
}

func handleTest(args []string) {

	grupo := GetNumGrupo(args)
	if grupo == -1 {
		return
	}

	initTime := time.Now()
	fmt.Println("Grupo asignado: ", grupo)

	var listadoGrupos [][]map[string]string
	if err := utils.LoadJsonFromFile(&listadoGrupos, core.Path_Grupos); err != nil {
		fmt.Println("Error al cargar grupos: ", err)
		return
	}

	carrera := listadoGrupos[grupo-1][0]
	extractor := core.NewExtractor()
	extractor.GetAsignaturasCarrera(core.ConstructCodigo(carrera["facultad"], carrera["carrera"]))
	fmt.Printf("Tiempo de ejecución final: %v\n", time.Since(initTime))
}

func handleExtract(args []string) {
	grupo := GetNumGrupo(args)
	if grupo == -1 {
		return
	}

	initTime := time.Now()
	data := ExtractCarrera(grupo)

	if data == nil {
		fmt.Println("Grupo no encontrado")
		return
	}

	filename := strconv.Itoa(grupo) + ".json"
	if err := utils.SaveJsonToFile(data, filename); err != nil {
		fmt.Println("Error al guardar archivo: ", err)
	}
	fmt.Println("......................................................")
	fmt.Printf("Tiempo de ejecución final: %v\n", time.Since(initTime))
}

func ExtractCarrera(indexGrupo int) map[string]*[]core.Asignatura {
	codigo := core.GetCodigoFromGrupo(indexGrupo)
	if codigo == nil {
		return nil
	}

	var data *[]core.Asignatura
	fmt.Println("Iniciando: ", codigo.Carrera)

	extractor := core.NewExtractor()

	if codigo.Facultad == core.ValuesElectiva.FacultadPor {
		data = extractor.ExtraerElectivas(*codigo)
	} else {
		data = extractor.GetAsignaturasCarrera(*codigo)
	}
	fmt.Println("Finalizado: ", codigo.Carrera)

	return map[string]*[]core.Asignatura{codigo.Carrera: data}
}
