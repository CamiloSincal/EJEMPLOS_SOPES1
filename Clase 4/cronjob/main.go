package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	scriptPath := "./script_cron.sh"

	hacerEjecutable(scriptPath)
	agregarCronJob(scriptPath)
	verificarCronJob(scriptPath)

	log.Println("Cronjob agregado exitosamente")
}

func hacerEjecutable(scriptPath string) {
	// os.Chmod cambia los permisos del archivo. 0755 es una notación octal:
	// el dueño puede leer/escribir/ejecutar (7), y otros solo leer/ejecutar (5)
	err := os.Chmod(scriptPath, 0755)
	if err != nil {
		log.Fatal("Error al hacer el script ejecutable:", err)
	}
	fmt.Printf("Script %s ahora es ejecutable\n", scriptPath)
}

func agregarCronJob(rutaScript string) {
	// "* * * * *" significa "ejecutar cada minuto"
	// El formato es: minuto hora día-del-mes mes día-de-la-semana
	expresionCron := "* * * * *"

	// ">> archivo.log 2>&1" redirige tanto la salida normal como los errores al archivo de log
	comandoCron := fmt.Sprintf("%s %s >> %s.log 2>&1", expresionCron, rutaScript, rutaScript)

	// Este comando hace tres cosas encadenadas con pipes (|):
	// 1. "crontab -l" lista los cronjobs existentes (2>/dev/null suprime el error si no hay ninguno)
	// 2. "echo" agrega la nueva línea al final
	// 3. "crontab -" Pasa el bloque de texto como entrada estándar para establecer la nueva configuración del usuario.
	cmd := exec.Command("bash", "-c",
		fmt.Sprintf("(crontab -l 2>/dev/null; echo \"%s\") | crontab -", comandoCron))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error agregando cronjob: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Cronjob agregado: %s", comandoCron)
}

func verificarCronJob(rutaScript string) {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("No se pudieron listar cronjobs (puede estar vacío): %v", err)
	} else {
		log.Printf("=== Cronjobs Actuales ===\n%s=== Fin de Cronjobs ===", string(output))
	}
}