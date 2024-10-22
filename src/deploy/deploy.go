package deploy

import (
	"fmt"
	"sia-extractor/src/core"
	"sia-extractor/src/utils"
)

const pathToData string = "artifacts/"

type DocumentCarrera struct {
	ID          string            `json:"_id" bson:"_id"`
	Facultad    string            `json:"facultad" bson:"facultad"`
	Carrera     string            `json:"carrera" bson:"carrera"`
	Asignaturas []core.Asignatura `json:"asignaturas" bson:"asignaturas"`
}

const (
	totalGrupos int = 40
)

func DeployData() {

	dbClient := NewDatabaseClient()
	fmt.Println("Connected to MongoDB!")

	defer dbClient.Disconnect()

	merged := MergeDataSede()
	if err := utils.SaveJsonToFile(merged, "data.json"); err != nil {
		println("Error al guardar archivo JSON: ", err)
	}

	dbClient.UpdateListadoCarreras(&merged)
	dbClient.SaveDataCarreras(&merged)
	// dbClient.SaveDataSede(&merged)
	dbClient.UpdateFechaExtraccion(merged["3068 FACULTAD DE MINAS"]["3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA"][0].FechaExtraccion)

}
