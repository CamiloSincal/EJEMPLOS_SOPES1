#!/bin/bash

# Array con los 3 comandos
comandos=(
  "docker run -d roldyoran/go-client"
  "docker run -d alpine sh -c 'while true; do echo 2^20 | bc > /dev/null; sleep 2; done'"
  "docker run -d alpine sleep 240"
)

# Seleccionar uno aleatorio
indice=$((RANDOM % 3))

echo "Ejecutando contenedor $((indice + 1)) de 3..."
eval "${comandos[$indice]}"