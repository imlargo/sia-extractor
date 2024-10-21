package app

import (
	"fmt"
	"sia-extractor/src/core"
	"sia-extractor/src/deploy"
	"sia-extractor/src/extractor"
	"sia-extractor/src/utils"
	"strconv"
	"time"
)

func handlePaths(args []string) {
	fmt.Println("Creando paths")

	extractor := extractor.NewExtractor()
	extractor.CreatePathsCarreras()
}

func handleGroup(args []string) {
	fmt.Println("Agrupando carreras")
	extractor := extractor.NewExtractor()
	extractor.GenerarGruposCarreras()
}

func handleElectivas(args []string) {
	grupo := GetNumGrupo(args)
	if grupo == -1 {
		return
	}

	initTime := time.Now()
	data := ExtractCarrera(grupo, true)

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
	extractor := extractor.NewExtractor()
	extractor.GetAsignaturasCarrera(core.ConstructCodigo(carrera["facultad"], carrera["carrera"]))
	fmt.Printf("Tiempo de ejecución final: %v\n", time.Since(initTime))
}

func handleExtract(args []string) {
	grupo := GetNumGrupo(args)
	if grupo == -1 {
		return
	}

	if grupo == 40 {
		return
	}

	initTime := time.Now()
	data := ExtractCarrera(grupo, false)

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

func ExtractCarrera(indexGrupo int, electivas bool) map[string]*[]core.Asignatura {
	codigo := core.GetCodigoFromGrupo(indexGrupo)
	if codigo == nil {
		return nil
	}

	var data *[]core.Asignatura
	fmt.Println("Iniciando: ", codigo.Carrera)

	extractor := extractor.NewExtractor()

	if electivas {
		data = extractor.ExtraerElectivas(*codigo)
	} else {
		data = extractor.GetAsignaturasCarrera(*codigo)
	}
	fmt.Println("Finalizado: ", codigo.Carrera)

	if len(*data) != 0 {
		dbClient := deploy.NewDatabaseClient()
		fmt.Println("Connected to MongoDB!")
		defer dbClient.Disconnect()

		dbClient.SaveCarrera(deploy.CreateDocumentCarrera(codigo.Carrera, codigo.Facultad, *data))
	}

	return map[string]*[]core.Asignatura{codigo.Carrera: data}
}
