#!/bin/bash

# Array de comandos a ejecutar
commands=(
    "go run . test 1 -rod=show,devtools"
    "go run . test 2 -rod=show,devtools"
    "go run . test 3 -rod=show,devtools"
)

# Ejecutar cada comando en paralelo
for i in "${!commands[@]}"; do
  ${commands[$i]} >> "output_$i.txt" &
done

# Esperar a que todos los comandos terminen
wait