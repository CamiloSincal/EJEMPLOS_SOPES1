#!/bin/bash
#Bucle while
contador=1
while [ $contador -le 5 ]; do # Mientras contador sea menor o igual a 5
    echo "Iteracion $contador"
    contador=$((contador + 1)) # Incrementar contador
done