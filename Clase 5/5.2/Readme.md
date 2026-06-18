# Flatcar Container Linux
Flatcar es una distribución de linux minimalista y especializada para la gestión de contenedores. Este forma parte de los **sistemas operativos inmutables**, lo que implica que algunos o todos los sistemas de archivos pertenecientes al SO son de solo lectura, por lo cuál no tiene gestor de paquetes y todo se ejecuta como contenedores.

---

## 1. Conceptos clave antes de empezar

| Concepto | Descripción |
|---|---|
| **Ignition** | Sistema de configuración que se ejecuta **una sola vez** al primer arranque. Configura usuarios, claves SSH, unidades systemd, etc. |
| **Butane** | Formato YAML legible para humanos que se compila a JSON de Ignition. |
| **Canal** | `stable`, `beta`, `alpha` — elige según el caso de uso. |
| **Imagen OEM** | Flatcar distribuye imágenes específicas por plataforma: `qemu`, `vmware`, `azure`, `aws`, etc. |

> **Flujo de trabajo básico:**
> `Escribir Butane (.yaml)` → `compilar a Ignition (.json)` → `pasar al arranque de la VM`

---

## 2. Opción A — Linux: QEMU/KVM

### 2.1 Requisitos

```bash
# Debian/Ubuntu
sudo apt install qemu-system-x86 qemu-utils wget

# Arch
sudo pacman -S qemu-full wget

# Fedora/RHEL
sudo dnf install qemu-kvm wget
```

Verifica que KVM esté disponible:

```bash
ls /dev/kvm   # debe existir
```

### 2.2 Descargar la imagen de QEMU

```bash
# Deben elegir el canal (stable recomendado para producción)
CHANNEL=stable
VERSION=$(curl -s "https://www.flatcar.org/releases-json/releases-${CHANNEL}.json" \
  | python3 -c "import sys,json; d=json.load(sys.stdin); print(sorted(d.keys())[-1])")

echo "Última versión $CHANNEL: $VERSION"

# Descargar imagen comprimida
wget "https://${CHANNEL}.release.flatcar-linux.net/amd64-usr/${VERSION}/flatcar_production_qemu_image.img.bz2"

# Descomprimir
bunzip2 flatcar_production_qemu_image.img.bz2
```

### 2.3 Crear disco de datos (opcional pero siempre es recomendado)

```bash
qemu-img create -f qcow2 data.qcow2 20G
```

### 2.4 Preparar Ignition y arrancar

```bash
qemu-system-x86_64 \
  -enable-kvm \
  -m 2048 \
  -smp 2 \
  -cpu host \
  -drive "file=flatcar_production_qemu_image.img,format=raw,if=virtio" \
  -drive "file=data.qcow2,format=qcow2,if=virtio" \
  -fw_cfg "name=opt/org.flatcar-linux/config,file=ignition.json" \
  -net nic,model=virtio \
  -net user,hostfwd=tcp::2222-:22 \
  -nographic
```

> `hostfwd=tcp::2222-:22` redirige el puerto 2222 del host al SSH de la VM.

### 2.5 Conectar por SSH

```bash
ssh -p 2222 core@localhost
```

---

## 3. Opción B — Windows: VMware Workstation/Player

### 3.1 Requisitos

- VMware Workstation Pro o Player (versión 16+)
- Descompresor para .bz2

### 3.2 Descargar la imagen de VMware

Se descarga desde el navegador (https://stable.release.flatcar-linux.net/amd64-usr/3975.2.2/flatcar_production_vmware_image.vmdk.bz2) y luego se descomprime para obtener `flatcar_production_vmware_image.vmdk`.

### 3.3 Crear la VM en VMware

1. Abrimos VMware → **Create a New Virtual Machine** → *Custom (advanced)*.
2. Seleccionamos **"I will install the operating system later"**.
3. Guest OS: **Linux** → *Other Linux 5.x or later kernel 64-bit*.
4. En el paso de disco duro, elegimos **"Use an existing virtual disk"** y lo apuntamos al `.vmdk` descargado.
5. Ajusta la RAM a mínimo **1024 MB** (recomendado 2048 MB) y CPUs a **2**.
6. **No arrancamos la máquina aún.**


### 3.4 Instalar Butane

```bash
# Linux
wget https://github.com/coreos/butane/releases/latest/download/butane-x86_64-unknown-linux-gnu
chmod +x butane-x86_64-unknown-linux-gnu
sudo mv butane-x86_64-unknown-linux-gnu /usr/local/bin/butane

# Windows — descargamos el .exe desde:
# https://github.com/coreos/butane/releases
```

### 3.5 Archivo Butane mínimo (`config.bu`)

```yaml
variant: flatcar
version: 1.0.0

passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-ed25519 AAAAC3Nza... tu_clave_publica_ssh

storage:
  files:
    - path: /etc/hostname
      contents:
        inline: mi-flatcar-node
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

### 3.6 Compilar a Ignition

Con el butane.exe en el mismo directorio que el config.bu, en el PowerShell ingresamos lo siguiente para comprobar su funcionamiento:

```shell
.\butane.exe --version
```

Si nos muestra la versión de butante entonces todo funciona correctamente y procedemos a ejecutar el siguiente comando:
```bash
.\butane.exe --pretty --strict config.bu > ignition.json
```

### 3.7 Agregamos el archivo Ignition como unidad CD-ROM

Flatcar para VMware lee Ignition desde una imagen ISO especial llamada `config-drive`.


En el archivo `.vmx` de la VM (lo abrimos con el bloc de notas), agregamos al final:

```
guestinfo.ignition.config.data.encoding = "base64"
guestinfo.ignition.config.data = "<CONTENIDO_BASE64_DE_ignition.json>"
```

Para generar el base64 en PowerShell:

```powershell
$json = Get-Content "ignition.json" -Raw
$bytes = [System.Text.Encoding]::UTF8.GetBytes($json)
[Convert]::ToBase64String($bytes)
```

Pegamos el resultado en el `.vmx` y en VMware, agregamos la ISO como unidad CD-ROM en la configuración de la VM.

### 3.8 Arrancar la VM

1. Iniciamos la VM en VMware.
2. Flatcar arrancará, aplicará Ignition y reiniciará.
3. Conectamos a la IP que muestra la consola (o configuramos port forwarding en la red NAT de VMware).

```powershell
# Desde Windows con OpenSSH instalado
ssh core@<IP_DE_LA_VM>
```

---

## 5. Primeros pasos dentro de Flatcar

Una vez dentro de la VM con `ssh core@...`:

### 5.1 Información del sistema

```bash
# Versión del OS
cat /etc/os-release

# Kernel
uname -r

# Estado de actualizaciones automáticas
update_engine_client -status
```

### 5.2 Correr un contenedor

```bash
# Docker viene incluido
docker run --rm hello-world

# Verificar estado de Docker
systemctl status docker
```

### 5.3 Desplegar un contenedor como servicio systemd

```bash
# Crear un unit file para nginx
sudo tee /etc/systemd/system/nginx.service > /dev/null <<'EOF'
[Unit]
Description=NGINX en contenedor
After=docker.service
Requires=docker.service

[Service]
Restart=always
ExecStartPre=-/usr/bin/docker stop nginx
ExecStartPre=-/usr/bin/docker rm nginx
ExecStart=/usr/bin/docker run --name nginx -p 80:80 nginx:alpine
ExecStop=/usr/bin/docker stop nginx

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now nginx
```

### 5.4 Sistema de archivos — qué es escribible

| Ruta | Acceso | Nota |
|---|---|---|
| `/` | Solo lectura | Raíz del OS, inmutable |
| `/etc` | Escritura parcial | Solo lo configurado vía Ignition o overlayfs |
| `/var` | Escritura | Datos persistentes, logs |
| `/home` | Escritura | Directorios de usuario |
| `/opt` | Escritura | Binarios adicionales |
| `/tmp` | Escritura | Temporal, se borra al reiniciar |

> Para instalar binarios extra usamos `/opt/bin/` y lo agregamos al PATH.

---


## Recursos adicionales

- [Documentación oficial de Flatcar](https://www.flatcar.org/docs/latest/)
- [Repositorio de Butane (GitHub)](https://github.com/coreos/butane)
- [Especificación de Ignition](https://coreos.github.io/ignition/configuration-v3_4/)
- [Releases de Flatcar](https://www.flatcar.org/releases)