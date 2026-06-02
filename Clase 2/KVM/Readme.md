# GUIA DE INSTALACIÓN DE KVM

Hipervisor KVM: Guía de principiante
(Basado en el artículo de Tytus Kurek)

## ¿Qué es KVM?

Proporciona a cada máquina virtual todos los servicios típicos del sistema físico, incluyendo BIOS virtual
(sistema básico de entrada/salida) y hardware virtual, como procesador, memoria, almacenamiento, tarjetas
de red, etc. Cada máquina virtual simula completamente una máquina física.

KVM se conecta directamente al código del kernel y le permite funcionar como hipervisor.

## Proceso de instalación

### 1. Instalar paquetes requeridos

```bash
sudo apt -y install bridge-utils cpu-checker libvirt-clients libvirt-daemon qemu qemu-kvm
```

### 2. Verificar capacidades de virtualización

```bash
kvm-ok
```

### 3. Ejecutar VM (Con comando)

```bash
sudo virt-install --name ubuntu-guest --os-variant ubuntu20.04 \
  --vcpus 2 --ram 2048 \
  --location http://ftp.ubuntu.com/ubuntu/dists/focal/main/installer-amd64/ \
  --network bridge=virbr0,model=virtio \
  --graphics none \
  --extra-args='console=ttyS0,115200n8 serial'
```

## KVM con virt-manager (Para interfaz gráfica) Ubuntu

### 1. Descargar virt-manager

```bash
apt-get install virt-manager
```

### 2. Ejecutar virt

```bash
sudo virt-manager
```

## Recomendaciones para la virtualización inicial de Ubuntu LTS

- 4k + 2 CPUs por ejemplo (Asignar según capacidades de cada máquina)
- Para el disco duro virtual se recomienda 20GB (Server) o de 30GB - 50GB para el desktop

## Recomendaciones para configuración (server)

- Actualizar siempre el instalador
- Usar distribución de teclado en español latinoamericano
- Usar la base de Ubuntu normal
- El proxy y mirror address no tocarlo y darle a continuar
- Dejar particiones por defecto
- No habilitar Ubuntu Pro
- No instalar SSH
- Instalar docker
- Al hacer upgrade reiniciar todos los servicios