#!/bin/bash

# Archivo de entrada
archivo_entrada="sistema.log"

# Archivo de salida
archivo_salida="reporte.txt"

# Procesar el archivo de entrada y guardar los resultados en el archivo de salida

# Explicación de comandos a utilizar
# grep <patron> <archivo> : Busca un patron en un archivo
# | : Pipe, redirige la salida de un comando a la entrada de otro
# cut -d <delimitador> -f <campo> : Extrae un campo de un archivo delimitado
# sort : Ordena las líneas de un archivo
# uniq -c : Cuenta las líneas repetidas de un archivo
# > : Redirige la salida de un comando a un archivo

grep "ERROR" $archivo_entrada | cut -d':' -f 2 - | sort | uniq -c > $archivo_salida

echo "Reporte generado en $archivo_salida"