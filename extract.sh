#!/bin/bash

# Obtener el número de grupo y la cantidad de elementos como argumentos
grupo=$1
cantidad=$2

# Calcular el rango basado en el número de grupo y la cantidad de elementos
start=$(( (grupo - 1) * cantidad + 1 ))
end=$(( grupo * cantidad ))

# Ejecutar el comando parallel con el rango calculado
parallel -j $cantidad --ungroup "go run . extract {1}" ::: $(seq $start $end)