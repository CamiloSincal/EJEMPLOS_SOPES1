use std::io::{Cursor, Read};
use tiny_http::{Server, Response, Header, Method, StatusCode};
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Clone)]
struct Mensaje {
    usuario: String,
    pais: String,
    mensaje: String,
}

#[derive(Serialize)]
struct Respuesta {
    usuario: String,
    pais: String,
    mensaje: String,
}

const GO_SERVICE_URL: &str = "http://go-api-service:8081/messages";

fn main() {
    let server = Server::http("0.0.0.0:8080").unwrap();
    println!("Servidor Rust en http://0.0.0.0:8080");
    println!("Forwarding a Go API: {}", GO_SERVICE_URL);

    for request in server.incoming_requests() {
        if request.url() == "/" || request.url() == "/health" {
            let body = b"ok";
            let _ = request.respond(Response::new(
                StatusCode(200), vec![], Cursor::new(body.to_vec()), Some(body.len()), None,
            ));
            continue;
        }

        if request.method() != &Method::Post || request.url() != "/messages" {
            eprintln!("[WARN] Ruta no encontrada: {} {}", request.method(), request.url());
            let body = b"404 - Ruta no encontrada";
            let _ = request.respond(Response::new(
                StatusCode(404), vec![], Cursor::new(body.to_vec()), Some(body.len()), None,
            ));
            continue;
        }

        handle(request);
    }
}

fn handle(mut req: tiny_http::Request) {
    let remote = req.remote_addr().map(|a| a.to_string()).unwrap_or_else(|| "desconocido".to_string());
    println!("[INFO] POST /messages desde {}", remote);

    // 1. Leer el body entrante
    let mut body = String::new();
    req.as_reader().read_to_string(&mut body).unwrap();
    println!("[INFO] Body recibido: {}", body);

    // 2. Validar que es JSON válido antes de forwarding
    let datos: Mensaje = match serde_json::from_str(&body) {
        Ok(d) => d,
        Err(e) => {
            eprintln!("[ERROR] JSON inválido: {}", e);
            let msg = format!("{{\"error\": \"JSON inválido: {}\"}}", e);
            let bytes = msg.into_bytes();
            let header = Header::from_bytes(b"Content-Type", b"application/json").unwrap();
            let _ = req.respond(Response::new(
                StatusCode(400), vec![header], Cursor::new(bytes.clone()), Some(bytes.len()), None,
            ));
            return;
        }
    };

    println!("[INFO] usuario={} | pais={} | mensaje={}", datos.usuario, datos.pais, datos.mensaje);

    // 3. Hacer POST al servicio de Go
    println!("[INFO] Forwarding a {}", GO_SERVICE_URL);
    match forward_to_go(&body) {
        Ok((status_code, response_body)) => {
            println!("[INFO] Respuesta de Go API - status={} body={}", status_code, response_body);
            let bytes = response_body.into_bytes();
            let header = Header::from_bytes(b"Content-Type", b"application/json").unwrap();
            let _ = req.respond(Response::new(
                StatusCode(status_code),
                vec![header],
                Cursor::new(bytes.clone()),
                Some(bytes.len()),
                None,
            ));
        }
        
        Err(e) => {
            eprintln!("[ERROR] Error al contactar Go API: {}", e);
            let msg = format!("{{\"error\": \"Error al contactar Go API: {}\"}}", e);
            let bytes = msg.into_bytes();
            let header = Header::from_bytes(b"Content-Type", b"application/json").unwrap();
            let _ = req.respond(Response::new(
                StatusCode(502),
                vec![header],
                Cursor::new(bytes.clone()),
                Some(bytes.len()),
                None,
            ));
        }
    }
}

fn forward_to_go(body: &str) -> Result<(u16, String), Box<dyn std::error::Error>> {
    let response = ureq::post(GO_SERVICE_URL)
        .set("Content-Type", "application/json")
        .send_string(body)?;

    let status = response.status();
    let response_body = response.into_string()?;

    Ok((status, response_body))
}