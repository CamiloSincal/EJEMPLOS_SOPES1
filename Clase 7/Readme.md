# Comandos para usar Zot como registro de contenedores

## 1. Iniciar el registro Zot en una VM con DOCKER

Ejecutamos el siguiente comando para iniciar un registro Zot en segundo plano, exponiendo el puerto 5000:

```bash
docker run -d -p 5000:5000 --name zot ghcr.io/project-zot/zot-linux-amd64:latest
```

Esto descargará la imagen de Zot y la ejecutará como un contenedor llamado `zot`.

---

## En la Computadora donde se tenga el desarrollo de la imagen:

### 2. Editamos la configuración de Docker para que pueda subir la imagen a la VM con DOCKER

#### a. Editamos la configuración de Docker:

```bash
sudo nano /etc/docker/daemon.json
```

#### b. Agregamos o modificamos el contenido del archivo (si está vacío, agregamos lo siguiente):

```json
{
  "insecure-registries": ["<IP_VM_DOCKER>:5000"]
}
```

#### c. Reiniciamos Docker para aplicar los cambios:
```bash
sudo systemctl restart docker
```

### 3. Creamos y etiquetamos la imagen para el registro privado

Para subir la imagen a docker, primero es necesario crear la imagen a través del comando:

```bash
docker build -t api-go .
```

Recordando siempre estar en la carpeta del `Dockerfile`

Usaremos de ejemplo la imagen llamada api-go. Cambiamos la etiqueta de la imagen para que apunte al registro privado (reemplazamos `<IP_VM_DOCKER>` por la IP de la VM):

```bash
docker tag api-go:v1 <IP_VM_DOCKER>:5000/api-go:v1
```

### 4. Subimos la imagen al registro Zot

Sube la imagen etiquetada al registro Zot:

```bash
docker push <IP_VM_DOCKER>:5000/api-go:v1
```

### 5. Verificamos las imágenes disponibles en el registro
Consultamos el catálogo de imágenes almacenadas en el registro Zot:

```bash
curl http://<IP_VM_DOCKER>:5000/v2/_catalog
```

*También se puede pegar la URL en el navegador para verificar que funciona*

---

# Configuración de `/etc/docker/daemon.json` en Flatcar

En Flatcar el sistema de archivos raíz es de **solo lectura**, por lo que no podemos editar `/etc/docker/daemon.json` directamente como en una distribución convencional.

---

## Opción 1: Archivo daemon.json en /etc (recomendada para VMs en ejecución)

Esta es la opción más directa cuando la VM ya está corriendo. El directorio `/etc` en Flatcar **sí es escribible**, únicamente `/usr` es de solo lectura.

Creamos el directorio y el archivo de configuración:

```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<EOF
{
  "insecure-registries": ["<IP_VM_DOCKER>:5000"]
}
EOF
```

Reiniciamos Docker para aplicar los cambios:

```bash
sudo systemctl restart docker
```

Verificamos que la configuración se aplicó correctamente:

```bash
docker info | grep -A5 "Insecure Registries"
```

---

## Opción 2: systemd drop-in

Usamos esta opción cuando preferimos pasar la configuración como flags directamente al servicio de Docker mediante un archivo override de systemd.

Creamos el directorio para el drop-in:

```bash
sudo mkdir -p /etc/systemd/system/docker.service.d/
```

Escribimos el archivo de configuración:

```bash
sudo tee /etc/systemd/system/docker.service.d/override.conf <<EOF
[Service]
Environment="DOCKER_OPTS=--insecure-registry=<IP_VM_DOCKER>:5000"
ExecStart=
ExecStart=/usr/bin/dockerd --insecure-registry=<IP_VM_DOCKER>:5000
EOF
```

Recargamos el daemon de systemd y reiniciamos Docker:

```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

---

## Opción 3: Butane / Ignition (aprovisionamiento desde cero)

Utilizamos esta opción cuando configuramos la VM desde cero. Declaramos el archivo `daemon.json` directamente en el manifiesto de Butane, y Ignition lo aplica durante el primer arranque de la máquina.

Configuramos nuestro archivo `butane.yaml` como:

```yaml
variant: flatcar
version: 1.0.0

passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh

storage:
  files:
    - path: /etc/hostname
      contents:
        inline: mi-flatcar-node
      mode: 0644
    - path: /etc/docker/daemon.json
      contents:
        inline: |
          {
            "insecure-registries": [" <IP_VM_DOCKER>:5000"]
          }
      mode: 0644

systemd:
  units:
    - name: hello.service
      enabled: true
      contents: |
        [Unit]
        Description=Hello Flatcar
        After=network-online.target

        [Service]
        Type=oneshot
        ExecStart=/bin/echo "¡Flatcar arrancó correctamente!"
        RemainAfterExit=yes

        [Install]
        WantedBy=multi-user.target
```

Para obtener la clave pública en windows:

Abrimos PowerShell y ejecutamos:

```powershell
ssh-keygen -t ed25519 -C "flatcar" -f $env:USERPROFILE\.ssh\id_flatcar
```

Presionamos Enter dos veces para omitir passphrase y obtenemos nuestra clave pública.
Vemos la clave generada:
```powershell
cat $env:USERPROFILE\.ssh\id_flatcar.pub
```

Si el archivo .bu no funciona será necesario agregar manualmente la clave ssh a la VM:
```powershell
echo "ssh-ed25519 AAAAC3Nza...tu-clave-completa" | update-ssh-keys -a ignition
```

Para conectarnos:
```powershell
ssh -i $env:USERPROFILE\.ssh\id_flatcar core@ip_vm
```

---

## Resumen

| Opción | Cuándo usarla |
|---|---|
| `daemon.json` en `/etc` | VM ya en ejecución, cambio rápido |
| systemd drop-in | Preferimos configurar a nivel de servicio |
| Butane / Ignition | Aprovisionamiento de la VM desde cero |


---

## En las maquinas virtuales con unicamente containerd y ctr

### 6. Descargar la imagen desde el registro Zot

Descargamos la imagen desde el registro privado para comprobar que está disponible:

```bash
sudo ctr images pull --plain-http <IP_VM_DOCKER>:5000/api-go:v1
```

Listamos las imágenes

```bash
sudo ctr images ls
```

```bash
# si el comando anterior no funciona probamos con este
sudo ctr images list
```

Listo ya podemos conectarnos con el Registro de Contenedores privados de ZOT en la maquina virtual, ahora pongamoslo a prueba con containerd.

---

# Containerd

**containerd** es un runtime de contenedores de nivel industrial que maneja el ciclo de vida completo de los contenedores en un sistema host. Es el componente principal que Docker utiliza internamente, pero también puede usarse de forma independiente.

Características principales de containerd:

- **Runtime estándar:** Implementa las especificaciones OCI Runtime y Image
- **Gestión de imágenes:** Descarga, almacena y gestiona imágenes de contenedores
- **Ciclo de vida:** Crea, ejecuta, detiene y elimina contenedores
- **Snapshots:** Maneja sistemas de archivos en capas para los contenedores
- **Networking:** Proporciona capacidades básicas de red para contenedores

`ctr` es la herramienta de línea de comandos que viene con containerd, similar a como docker es el cliente para Docker Engine.

---

## Comandos básicos

### 1. Para verificar si containerd está corriendo

```bash
sudo systemctl status containerd
```

```bash
ctr --version
```

### 2. Para listar las imágenes disponibles

```bash
sudo ctr images ls
```

```bash
# si el comando anterior no funciona probamos con
sudo ctr images list
```

### 3. Para descargar una imagen desde un registry (como Docker Hub)

```bash
sudo ctr images pull docker.io/library/hello-world:latest
```

Otra opción es con un registro privado como Zot en una maquina virtual

```bash
sudo ctr images pull --plain-http <IP_VM1_DOCKER>:5000/api-go:v1
```

### 4. Para levantar el contenedor a partir de la imagen

```bash
sudo ctr run -t --rm docker.io/library/hello-world:latest my-hello
```

De este comando es necesario saber que:
- `-d` : ejecuta el contenedor en segundo plano (detached).
- `--tty` : asigna una terminal al contenedor.
- `--plain-http` : utiliza la red del host, permitiendo que el contenedor acceda directamente a los puertos y servicios del host.
- `<IP_VM_DOCKER>:5000/api-go:v1` : imagen a ejecutar, obtenida desde el registro privado.
- `my-api-go` : nombre local para el contenedor.

### 5. Para listar los contenedores activos

```bash
sudo ctr containers ls
```

### 6. Para crear y ejecutar un contenedor en pasos separados

```bash
sudo ctr images pull docker.io/library/alpine:latest
sudo ctr containers create docker.io/library/alpine:latest my-alpine
sudo ctr tasks start -d my-alpine
```

### 7. Para subir una imagen a un registry (como Zot o Docker Hub)

```bash
sudo ctr images tag docker.io/library/hello-world:latest localhost:5000/hello-world:latest
sudo ctr images push <IP_VM1_DOCKER>:5000/hello-world:latest
```

### 8. Para eliminar una imagen

```bash
sudo ctr images rm docker.io/library/hello-world:latest
```

### 9. Para eliminar un contenedor

Primero detenemos la tarea del contenedor con:

```bash
sudo ctr tasks list         # para obtener el nombre de la tarea
sudo ctr task kill <nombre-de-tarea>               # para detener la tarea
sudo ctr task kill --signal SIGKILL <nombre-de-tarea>  # para obligar a detener la tarea
```

Con la tarea detenida ya podemos eliminar el contenedor:

```bash
sudo ctr containers delete my-hello
```

### 10. Para acceder a un shell en un contenedor

Si el contenedor está corriendo y tiene `/bin/sh`:

```bash
sudo ctr tasks exec -t --exec-id myexecid my-alpine /bin/sh
```

### 11. Para eliminar todo

```bash
sudo ctr tasks list -q | xargs -I {} sudo ctr task kill {}
sudo ctr containers list -q | xargs -I {} sudo ctr containers delete {}
sudo ctr images list -q | xargs -I {} sudo ctr images remove {}
```
