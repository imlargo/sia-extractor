package main

import (
	"encoding/json"
	"os"
	"sia-extractor/core"
)

func main() {

	codigo := core.Codigo{
		Nivel:     core.ValueNivel,
		Sede:      core.ValueSede,
		Facultad:  "3068 FACULTAD DE MINAS",
		Carrera:   "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA",
		Tipologia: "TODAS MENOS  LIBRE ELECCIÓN",
	}

	asignaturas := core.GetAsignaturasCarrera(codigo)

	dataAsignaturasJSON, _ := json.Marshal(asignaturas)
	os.WriteFile("asignaturas.json", dataAsignaturasJSON, 0644)

	// core.CreatePathsCarreras()
}
