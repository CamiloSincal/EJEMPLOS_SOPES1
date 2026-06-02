# Bash scripts

Es un archivo de texto plano que contiene una secuencia de comandos para el intérprete Bourne-Again SHell (Bash) y se usa para automatizar tareas repetitivas, gestionar archivos, ejecutar programas y realizar otras operaciones del sistema operativo al agrupar múltiples comandos en un solo archivo ejecutable.

## Creación y ejecución manual de scripts de bash

### 1. Creación del script

Creamos un archivo con extensión `*.sh`, por ejemplo:

```bash
#!/bin/bash
echo "Hola Mundo"
```

### 2. Cambio de los permisos del archivo

Es necesario cambiar los permisos para hacer el archivo ejecutable:

```bash
chmod +x nombre_script.sh
```

### 3. Verificación de permiso

```bash
ls -l nombre_script.sh
```

### 4. Ejecutar archivo

Con los permisos cambiados el script se ejecuta con:

```bash
./nombre_script.sh
```

Alternativamente, se puede ejecutar invocando directamente el intérprete bash para no hacer el cambio de permisos:

```bash
bash nombre_script.sh
```

