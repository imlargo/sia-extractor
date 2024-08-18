package main

import "sia-extractor/core"

func main() {

	/*
		codigo := core.Codigo{
			Nivel:     core.ValueNivel,
			Sede:      core.ValueSede,
			Facultad:  "3068 FACULTAD DE MINAS",
			Carrera:   "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA",
			Tipologia: "TODAS MENOS  LIBRE ELECCIÓN",
		}

		core.GetAsignaturasCarrera(codigo)
	*/

	core.CreatePathsCarreras()
}
