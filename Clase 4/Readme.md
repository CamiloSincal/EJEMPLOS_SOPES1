# Introducción a Golang

### 1. Instalación de golang en linux

#### Opción 1: Descarga de la página oficial directamente

```
https://go.dev/dl/
```

#### Opción 2: Descarga de la página oficial con curl

```bash
curl -OL https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
```

##### Verificar el archivo descargado

```bash
sha256sum <FILENAME_OF_GO_TARBALL>
```

##### Seguir procedimiento indicado por la página de golang:

```
https://go.dev/doc/install
```

# Introducción a golang

A continuación se presenta un ejemplo básico en lenguaje go:

```go
package main

import "fmt"

func main() {
    fmt.Println("¡Hola, mundo!")
}
```

Ejecución de script con Golag
```bash
go run hola_mundo.go
```

Si se quiere, es posible compilar y ejecutar el archivo de go con:
```bash
go build hola_mundo.go

./holamundo
```

# Daemon en Go
## EJ1 - Daemon simple en Go
### 1. Código del daemon y compilación
Este daemon escribe un log cada 5 segundos, así se podrá ver si se corre en segundo plano
```go
package main

import (
	"log"
	"os"
	"time"
)

func main() {
	// Abrir archivo de log
	f,err := os.OpenFile("/var/log/mydaemon.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Daemon iniciado")

	for {
		log.Println("Daemon ejecutando tarea...")
		time.Sleep(5 * time.Second) // Simula una tarea que se ejecuta cada 5 segundos
	}
}
```

Para compilar el código se usa:
```bash
go build -o /usr/local/bin/mydaemon main.go
```

### 2. Creación del archivo systemd

Un archivo de systemd es necesario porque actúa como el puente entre la aplicación y el sistema operativo. Se debe crear `/etc/systemd/system/mydaemon.service` con este contenido:

```ini
[Unit]
Description=Mi daemon en Go
After=network.target

[Service]
ExecStart=/usr/local/bin/mydaemon
Restart=always

[Install]
WantedBy=multi-user.target
```

### Inicio automático

La sección `[Install]` con `WantedBy=multi-user.target` le dice a systemd que inicie el daemon automáticamente cuando el sistema arranque. Sin esto, se tendría que iniciar el programa manualmente cada vez que se reinicie el servidor.

### Supervisión y recuperación

La línea `Restart=always` permite que si el daemon se cae por cualquier razón (un error, falta de memoria, etc.), systemd lo reiniciará automáticamente. Esto es especialmente importante para servicios que deben estar disponibles 24/7.

### Dependencias

`After=network.target` asegura que el daemon no intente iniciarse antes de que la red esté disponible. Esto evita errores si la aplicación necesita conectarse a bases de datos, APIs, o cualquier recurso de red.

### Integración con el sistema

Con systemd es posible:

- Ver el estado: `systemctl status mydaemon`
- Ver logs: `journalctl -u mydaemon`
- Iniciar/detener: `systemctl start/stop mydaemon`
- Habilitar al inicio: `systemctl enable mydaemon`

## 3. Activar y arrancar el daemon

```bash
sudo systemctl daemon-reload       # recargar systemd
sudo systemctl enable --now mydaemon
```

## 4. Comandos de verificación

Estado del servicio:

```bash
systemctl status mydaemon
```

Logs en journald:

```bash
journalctl -u mydaemon -f
```

Logs en el archivo que configuramos (`/var/log/mydaemon.log`):

```bash
tail -f /var/log/mydaemon.log
```

La salida a observar debería ser:

```
Aug 21 16:10:12 myhost mydaemon[1234]: Daemon iniciado
Aug 21 16:10:17 myhost mydaemon[1234]: Daemon sigue vivo...
Aug 21 16:10:22 myhost mydaemon[1234]: Daemon sigue vivo...
```

## EJ2 - Cronjob

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Abrir del script existente
	scriptPath := "./script_cron.sh"

	// 1. Hacer el script ejecutable
	hacerEjecutable(scriptPath)

	// 2. Agregar el script al cron
	agregarCronJob(scriptPath)

	// 3. Verificar que el cronjob se agregó correctamente
	verificarCronJob(scriptPath)

	log.Println("Cronjob agregado exitosamente")
}

// Función que cambia el script para que sea ejecutable
func hacerEjecutable(scriptPath string) {
	err := os.Chmod(scriptPath, 0755)
	if err != nil {
		log.Fatal("Error al hacer el script ejecutable:", err)
	}
	fmt.Printf("Script %s ahora es ejecutable\n", scriptPath)
}

func agregarCronJob(rutaScript string) {
    expresionCron := "* * * * *"
    comandoCron := fmt.Sprintf("%s %s >> %s.log 2>&1", expresionCron, rutaScript, rutaScript)

    cmd := exec.Command("bash", "-c",
        fmt.Sprintf("(crontab -l 2>/dev/null; echo \"%s\") | crontab -", comandoCron))

    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("Error agregando cronjob: %v\nOutput: %s", err, string(output))
    }

    log.Printf("Cronjob agregado: %s", comandoCron)
}

// verificarCronJobs lista todos los cronjobs configurados
func verificarCronJobs() {
    cmd := exec.Command("crontab", "-l")
    output, err := cmd.CombinedOutput()

    if err != nil {
        log.Printf("No se pudieron listar cronjobs (puede estar vacío): %v", err)
    } else {
        log.Printf("=== Cronjobs Actuales ===\n%s=== Fin de Cronjobs ===", string(output))
    }
}
```

Para ejecutar y analizar el cronjob se usan los siguientes comandos:
```bash
# Eliminar el cronjob
crontab -l | grep -v "script_cron.sh" | crontab -

# Ejecutar el código
go run main.go

# Verificar que el cronjob fue agregado
crontab -l

# Verificar el log generado
cat script_cron.sh.log

# Verificar en tiempo real
tail -f script_cron.sh.log
```