#!/bin/bash
# Un script para crea un archivo comprimido que contiene los archivos de un directorio espeficio

# origen (este directorio)
origen="."
# destino un archivo comprimido afuera de este
destino="../"

# checamos si el destino existe, sino lo creamos

if [ ! -d "$destino" ]; then # -d es un flag para verificar si un directorio existe
    mkdir "$destino"
fi

# checamos si tenemos permisos para escribir en el destino

if [ ! -w "$destino" ]; then # -w es un flag para verificar si tenemos permisos de escritura
    echo "No tienes permisos para escribir en $destino"
    exit 1
fi

# creamos el backup
# tar -> es un comando para crear archivos comprimidos
# -czvf -> c: crea un archivo, z: comprime, v: verbose, f: nombre del archivo
# "$destino/backup_$(date +%Y%m%d).tar.gz" -> nombre del archivo
# "$origen" -> directorio a comprimir
tar -czvf "$destino/backup_$(date +%Y%m%d).tar.gz" "$origen"

echo "Backup creado en $destino/backup_$(date +%Y%m%d).tar.gz"