#!/bin/bash

# Este script comprueba el uso del disco y envía una advertencia si el uso supera un umbral especificado (80% en este caso).

# Obtener el uso del disco

# Comandos a utilizar
# $() : Ejecuta un comando y guarda su salida
# df -h / : Muestra el uso de disco de la partición raíz
# | : Pipe, redirige la salida de un comando a la entrada de otro
# grep / : Filtra las líneas que contienen la partición raíz
# awk '{print $5}' : Imprime la quinta columna de un archivo
# sed 's/%//g' : Elimina los caracteres '%' de un archivo

uso_disco=$(df -h / | grep / | awk '{print $5}' | sed 's/%//g')

echo "Uso de disco: $uso_disco%"

# Comprobar si el uso de disco supera el 80%

if [ "$uso_disco" -gt 80 ]; then # Si el uso de disco supera el 80%
    echo "Advertencia: El uso de disco supera el 80%"
else
    echo "Uso de disco normal"
fi