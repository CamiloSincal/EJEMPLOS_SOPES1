# Rust
## Guía de instalación
### Windows

#### 1. Descargar rustup

Vamos a [https://rustup.rs](https://rustup.rs) y descargamos el instalador `rustup-init.exe`.

#### 2. Ejecutar el instalador

Abrimos `rustup-init.exe` y seguimos las instrucciones. Por defecto se instala:

- `rustc` — el compilador
- `cargo` — el gestor de paquetes
- `rustup` — el gestor de versiones

Cuando nos pregunte por el tipo de instalación, elegimos la opción `1) Proceed with standard installation`.

#### 3. Verificar la instalación

Abrimos una nueva terminal (CMD o PowerShell) y ejecutamos:

```
rustc --version
cargo --version
```

> **Nota:** Es necesario tener instalado [Microsoft C++ Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/) para compilar algunos proyectos, el instalador avisará si hace falta.

---

### Linux

#### 1. Instalar con rustup

Abrimos una terminal y ejecutamos:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

Cuando pregunte por el tipo de instalación, elegimos la opción `1) Proceed with standard installation`.

#### 2. Cargar el entorno

Después de la instalación, cargamos las variables de entorno en la sesión actual:

```bash
source "$HOME/.cargo/env"
```

Para que se cargue automáticamente en futuras sesiones, usamos:

```bash
echo 'source "$HOME/.cargo/env"' >> ~/.bashrc
```

Luego recargamos el archivo para aplicar los cambios en la sesión actual:

```bash
source ~/.bashrc
```

#### 3. Verificar la instalación

```bash
rustc --version
cargo --version
```

---

### Comandos básicos de cargo

| Comando | Descripción |
|---|---|
| `cargo new mi_proyecto` | Crea un nuevo proyecto |
| `cargo build` | Compila el proyecto |
| `cargo run` | Compila y ejecuta |
| `cargo test` | Ejecuta las pruebas |
| `cargo update` | Actualiza las dependencias |

---

### Actualizar o desinstalar Rust

```bash
# Actualizar
rustup update

# Desinstalar
rustup self uninstall
```

---

## Conceptos básicos de Rust

Rust es un lenguaje de programación de sistemas enfocado en tres pilares: **seguridad de memoria**, **rendimiento** y **concurrencia**. A diferencia de otros lenguajes, logra esto sin necesidad de un recolector de basura.

### Crear, compilar y ejecutar un proyecto

Rust usa `cargo` como herramienta principal para gestionar proyectos. Es la forma recomendada frente a usar `rustc` directamente.

#### 1. Creamos el proyecto

```bash
cargo new mi_proyecto
```

Esto genera la siguiente estructura:

```
mi_proyecto/
├── Cargo.toml       # metadatos y dependencias del src/
    └── main.rs      # punto de entrada del programa
```

El archivo `main.rs` ya viene con un "Hola, mundo!" por defecto:

```rust
fn main() {
    println!("Hello, world!");
}
```

#### 2. Entramos al directorio del proyecto

```bash
cd mi_proyecto
```

Todos los comandos de `cargo` deben ejecutarse desde la raíz del proyecto (donde está el `Cargo.toml`).

#### 3. Compilamos el proyecto

```bash
cargo build
```

Esto genera el ejecutable en `target/debug/mi_proyecto`. Es una compilación de desarrollo, sin optimizaciones.

Si queremos una compilación optimizada para producción:

```bash
cargo build --release
```

El ejecutable se genera en `target/release/mi_proyecto`.

#### 4. Ejecutamos el proyecto

```bash
cargo run
```

Este comando compila y ejecuta en un solo paso. Es el más usado durante el desarrollo. Si el código no cambió desde la última compilación, `cargo` omite la compilación y ejecuta directamente.

#### 5. Verificamos sin compilar

Si solo queremos saber si el código tiene errores sin generar el ejecutable:

```bash
cargo check
```

Es más rápido que `cargo build` y es útil mientras escribimos código.

---

**Resumen del flujo típico:**

| Comando | Cuándo usarlo |
|---|---|
| `cargo new nombre` | Al iniciar un proyecto nuevo |
| `cargo check` | Para verificar errores rápidamente |
| `cargo build` | Para compilar en modo desarrollo |
| `cargo run` | Para compilar y ejecutar en desarrollo |
| `cargo build --release` | Para compilar la versión final optimizada |

### Variables y mutabilidad

Por defecto, las variables en Rust son **inmutables** *(una vez que se les asigna un valor a una variable, no se les puede cambiar ese valor ni reasignar otro en su lugar)*. Para hacerlas mutables usamos `mut`:

```rust
let x = 5;        // inmutable
let mut y = 10;   // mutable
y = 20;           // válido
```

### Tipos de datos básicos

```rust
let entero: i32 = 42;
let flotante: f64 = 3.14;
let booleano: bool = true;
let caracter: char = 'R';
let texto: &str = "Hola, Rust";
```

### Funciones

```rust
fn sumar(a: i32, b: i32) -> i32 {
    a + b  // sin punto y coma = valor de retorno
}

fn main() {
    let resultado = sumar(3, 7);
    println!("Resultado: {}", resultado);
}
```

### Ownership (propiedad)

Es el concepto más importante de Rust. Cada valor tiene un único dueño, y cuando ese dueño sale del alcance, el valor se libera automáticamente de memoria. Esto elimina errores comunes como punteros nulos o memory leaks.

```rust
let s1 = String::from("hola");
let s2 = s1;  // s1 ya no es válido, s2 es el nuevo dueño

// println!("{}", s1);  // error: s1 fue movido
println!("{}", s2);     // válido
```

### Control de flujo

```rust
// Condicional
let numero = 7;
if numero > 5 {
    println!("mayor que 5");
} else {
    println!("menor o igual a 5");
}

// Bucle
for i in 0..5 {
    println!("{}", i);
}
```

### Manejo de errores con `Result`

Rust no usa excepciones. En su lugar, las funciones que pueden fallar retornan un tipo `Result<T, E>`:

```rust
use std::fs;

match fs::read_to_string("archivo.txt") {
    Ok(contenido) => println!("{}", contenido),
    Err(error)    => println!("Error: {}", error),
}
```

---

# Ejemplos prácticos de Rust

---

## Ejemplo 1: Calculadora de notas

Combina: variables, tipos de datos, funciones, control de flujo y mutabilidad.

```rust
// Una función que recibe un vector de notas (f64) y devuelve el promedio
// El símbolo & significa que "prestamos" el vector sin tomar su ownership
fn calcular_promedio(notas: &Vec<f64>) -> f64 {
    let mut suma = 0.0; // variable mutable porque la vamos a modificar en el bucle

    for nota in notas {
        suma += nota; // acumulamos cada nota en suma
    }

    suma / notas.len() as f64 // retorno implícito: sin punto y coma al final
}

// Una función que recibe el promedio y devuelve una calificación como texto
fn calificacion(promedio: f64) -> &'static str {
    // &'static str es un texto de duración indefinida (vive todo el programa)
    if promedio >= 90.0 {
        "Excelente"
    } else if promedio >= 70.0 {
        "Aprobado"
    } else {
        "Reprobado"
    }
}

fn main() {
    // Vec es un arreglo dinámico, equivalente a una lista en otros lenguajes
    let notas: Vec<f64> = vec![85.0, 92.0, 78.0, 60.0, 95.0];

    let promedio = calcular_promedio(&notas); // pasamos &notas para no ceder el ownership

    println!("Notas: {:?}", notas);           // {:?} imprime el vector completo
    println!("Promedio: {:.2}", promedio);    // {:.2} limita a 2 decimales
    println!("Calificación: {}", calificacion(promedio));
}
```

**Salida esperada:**
```
Notas: [85.0, 92.0, 78.0, 60.0, 95.0]
Promedio: 82.00
Calificación: Aprobado
```

---

## Ejemplo 2: Registro de productos con manejo de errores

Combina: ownership, funciones, control de flujo, tipos de datos y manejo de errores con `Result`.

```rust
// Definimos una estructura para representar un producto
// #[derive(Debug)] permite imprimir la estructura con {:?}
#[derive(Debug)]
struct Producto {
    nombre: String,
    precio: f64,
    stock: u32, // u32 = entero sin signo (no puede ser negativo)
}

// Función que intenta realizar una venta
// Devuelve Result: Ok si tuvo éxito, Err con un mensaje si falló
fn vender(producto: &mut Producto, cantidad: u32) -> Result<f64, String> {
    // Verificamos si hay suficiente stock
    if cantidad > producto.stock {
        // Err detiene la operación y devuelve el mensaje de error
        return Err(format!(
            "Stock insuficiente. Disponible: {}, solicitado: {}",
            producto.stock, cantidad
        ));
    }

    producto.stock -= cantidad; // reducimos el stock
    let total = producto.precio * cantidad as f64; // calculamos el total

    Ok(total) // Ok indica que todo salió bien y devuelve el valor
}

fn main() {
    // mut porque vamos a modificar el stock del producto
    let mut producto = Producto {
        nombre: String::from("Teclado mecánico"),
        precio: 350.0,
        stock: 5,
    };

    println!("Producto: {:?}\n", producto);

    // Intentamos dos ventas: una válida y una que excede el stock
    let ventas = vec![2, 10]; // queremos vender 2 unidades, luego 10

    for cantidad in ventas {
        // match maneja los dos posibles resultados de vender()
        match vender(&mut producto, cantidad) {
            Ok(total) => println!(
                "Venta exitosa: {} unidades. Total: Q{:.2}. Stock restante: {}",
                cantidad, total, producto.stock
            ),
            Err(mensaje) => println!("Error en la venta: {}", mensaje),
        }
    }
}
```

**Salida esperada:**
```
Producto: Producto { nombre: "Teclado mecánico", precio: 350.0, stock: 5 }

Venta exitosa: 2 unidades. Total: Q700.00. Stock restante: 3
Error en la venta: Stock insuficiente. Disponible: 3, solicitado: 10
```

---

## Ejemplo 3: Procesador de nombres de archivos

Combina: ownership, funciones, control de flujo, tipos de datos, mutabilidad y manejo de errores con `Result`.

```rust
// Función que recibe el nombre de un archivo y extrae su extensión
// Option<&str> significa que puede devolver un texto (&str) o nada (None)
fn obtener_extension(nombre: &str) -> Option<&str> {
    // find busca el último punto en el nombre del archivo
    // rfind busca desde el final, útil para casos como "archivo.backup.txt"
    match nombre.rfind('.') {
        Some(indice) => Some(&nombre[indice + 1..]), // devolvemos todo después del punto
        None => None, // el archivo no tiene extensión
    }
}

// Función que clasifica el archivo según su extensión
// Devuelve Result: Ok con la categoría, Err si la extensión no es reconocida
fn clasificar_archivo(nombre: &str) -> Result<String, String> {
    // Obtenemos la extensión; si no hay, devolvemos error de inmediato
    let extension = match obtener_extension(nombre) {
        Some(ext) => ext,
        None => return Err(format!("'{}' no tiene extensión", nombre)),
    };

    // Clasificamos según la extensión
    // match con múltiples patrones usando |
    let categoria = match extension {
        "rs"                    => "Código fuente Rust",
        "txt" | "md"            => "Documento de texto",
        "png" | "jpg" | "jpeg"  => "Imagen",
        "mp3" | "wav"           => "Audio",
        "zip" | "tar" | "gz"    => "Archivo comprimido",
        otra => return Err(format!("Extensión '.{}' no reconocida", otra)),
    };

    Ok(String::from(categoria))
}

fn main() {
    // Lista de archivos a procesar
    let archivos = vec![
        "main.rs",
        "notas.txt",
        "foto_viaje.jpg",
        "backup.zip",
        "documento",      // sin extensión
        "video.avi",      // extensión no reconocida
        "README.md",
    ];

    println!("Procesando {} archivos:\n", archivos.len());

    for archivo in &archivos { // & para no ceder el ownership del vector
        // Construimos la salida según el resultado de clasificar_archivo
        let resultado = match clasificar_archivo(archivo) {
            Ok(categoria) => format!("✓ {:<25} → {}", archivo, categoria),
            Err(error)    => format!("✗ {:<25} → {}", archivo, error),
        };

        println!("{}", resultado);
    }
}
```

**Salida esperada:**
```
Procesando 7 archivos:

✓ main.rs                   → Código fuente Rust
✓ notas.txt                 → Documento de texto
✓ foto_viaje.jpg            → Imagen
✓ backup.zip                → Archivo comprimido
✗ documento                 → 'documento' no tiene extensión
✗ video.avi                 → Extensión '.avi' no reconocida
✓ README.md                 → Documento de texto
```

## Seguir aprendiendo

Para profundizar en Rust, la documentación oficial es el mejor punto de partida:

- [The Rust Programming Language (El libro oficial)](https://doc.rust-lang.org/book/) — guía completa desde cero, gratuita y en línea.