package app

import (
	"fmt"
	"strconv"
)

func GetNumGrupo(args []string) int {
	if len(args) < 2 {
		fmt.Println("Debe ingresar el número de grupo")
		return -1
	}

	grupoStr := args[1]

	grupo, err := strconv.Atoi(grupoStr)
	if err != nil {
		fmt.Println("Error al convertir grupo: ", err)
		return -1
	}

	return grupo
}
