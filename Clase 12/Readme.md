# Primer despliegue en k8s
En esta sección se explicará como realizar un primer despliegue de un sistema distribuido simple, el cual se conformará de:
- Deployment de API REST con RUST
- Deployment de API REST con Go + Client gRPC de Go
- Deployment de Server gRPC de Go
- Deployment de RabbitMQ
- Deployment Consumer de RabbitMQ con Go

## Arquitectura de la clase

![arqui-clase12](./imgs/diagram.PNG)

> Para este ejemplo/ejercico se considerará que el Deployment de RabbitMQ ya se encuentra funcionando según las instrucciones de la Clase 11. Además, no se profundizará en los Scripts, ya que el objetivo principal es comprender kubernetes.

--- 

## Estructura de archivos y partes importantes

La estructura general para esta clase es:

```
Clase 12/
├── Scripts/                # Código de cada API y del consumer
|    ├──  rest_rust  
|    ├──  rest_client_go  
|    ├──  server_go 
|    ├──  consumer_go           
├── Manifiestos/           # Archivos YML de Kubernetes
|    ├──  consumer-deployment.yml 
|    ├──  gateway.yml
|    ├──  go_rest_client_grpc.yml
|    ├──  go_server_grpc.yml
|    ├──  rest_rust.yml
```
### API REST con RUST
Se tomará el código de la clase 10, con la diferencia de que ahora la api de rest en rust reenvia el código hacia una api rest de go que está en otro deployment. De este nuevo código lo necesario a tomar en cuenta es el "enlace" de conexión a go:
```go
// URL del servicio Go en Kubernetes (nombre-servicio.namespace.svc.cluster.local:puerto)
// Si están en el mismo namespace, solo es: http://nombre-servicio:8081/messages
const GO_SERVICE_URL: &str = "http://go-api-service:8081/messages";
```
### API REST con GO
Para la API REST de la conexión hacia el Cliente de gRPC se configura con:
```go
const grpcAddress = "localhost:50051"
```

Al estar tanto la api rest y el cliente grpc en el mismo pod de kubernetes, su comunicación es más simple y se puede utilizar `localhost:50051` como url de comunicación.

### Client gRPC de Go
Para el cliente de gRPC también se usará el nombre del servicio del servidor de gRPC.
```go
// Dirección del servidor gRPC externo (Deployment separado en K8s)
// Formato: nombre-service:puerto
const grpcServerAddress = "grpc-server-service:50051"
```

### Server gRPC de Go
Para este punto, el server gRPC obtendrá los datos del Client y servirá como el producer de rabbitMQ (como se vió en clases anteriores), por ende la variables más imporantes a tener en cuenta son:
```go
// Leer variables de entorno
rabbitURL := getEnv("RABBIT_URL", "amqp://guest:guest@rabbitmq-cluster.rabbitmq-system.svc.cluster.local:5672/")
queueName := getEnv("QUEUE_NAME", "mensajes")
```

Al igual que en el ejemplo de RabbitMQ de la clase 11, se usan variables de entorno con variables por defecto, así que es preferible no modificar estas líneas. Además, se debe tomar en cuenta que el texto de `@rabbitmq-cluster.rabbitmq-system.svc.cluster.local:5672/` tiene el formato de `nombre-servicio.namespace.svc.cluster.local:puerto`, así que si no se están usando namespaces es necesario modificar esta línea con la dirección del servicio correcto.

### Consumer de RabbitMQ con Go
El consumer de RabbitMQ a implementar es el mismo utilizado en la Clase 11, y la única consideración es la misma respecto al server gRPC y su conexión con RabbitMQ:

```go
rabbitURL := envOrDefault("RABBITMQ_URL", "amqp://guest:guest@rabbitmq-cluster.rabbitmq-system.svc.cluster.local:5672/")
queueName := envOrDefault("RABBITMQ_QUEUE", "mensajes")
```

`rabbitURL` debe tener el formato correcto según si está o no en un namespace y `queueName` debe ser igual tanto en el server gRPC como en el consumer.

> Tomar en cuenta que para todos los archivos de Go (a excepción del consumer) deben tener los archivos de proto para poder establecer comunicación a través de gRPC

## Despliegue
Para el despliegue en kubernetes primero hacemos build de las imagenes y las cargamos a nuestro servidor con ZOT.

```bash
docker build -t <DIRECCION-ZOT>/rust-api:latest . 
docker push <DIRECCION-ZOT>/rust-api:latest

docker build -t <DIRECCION-ZOT>/go-rest-api:latest . 
docker push <DIRECCION-ZOT>/go-rest-api:latest

docker build -t <DIRECCION-ZOT>/go-grpc-client:latest . 
docker push <DIRECCION-ZOT>/go-grpc-client:latest

docker build -t <DIRECCION-ZOT>/go-grpc-server:latest .
docker push <DIRECCION-ZOT>/go-grpc-server:latest

docker build -t <DIRECCION-ZOT>/rabbitmq-consumer:latest .
docker push <DIRECCION-ZOT>/rabbitmq-consumer:latest
```
Antes de hacer el apply, hay que tener en cuenta algunas cosas relacionados a los manifiestos de kubernetes de manera general:
- Para esta clase se usan archivos YAML multidocumento, es decir se usan `---` para separar los `.yml` en diferentes secciones para desplegar diferentes objetos en un solo archivo.
- Es necesario cambiar `image` en todos los manifiestos con la IP del servidor ZOT.
- Es necesario revisar para cada `.yml` el nombre del objeto en la `metadata`, estas deben conincidir con las utilizadas en los Scripts para que la comunicación se ejecute correctamente.
- Es necesario cambiar la variable de entorno `RABBITMQ_URL` con el usuario, contraseña y url correctos del clúster.
- Es necesario cambiar la variable de entorno `RABBITMQ_QUEUE` para que conicidan entre el gRPC Server(Producer) y el Consumer.
  
De manera indivudual, el archivo `go_rest_client_grpc.yml` es el único con algo nuevo:
```yml
spec:
      containers:
        - name: rest-api
          image:<DIRECCION-DE-ZOT>/go-rest-api:latest
          ports:
            - containerPort: 8081

        - name: grpc-client
          image:<DIRECCION-DE-ZOT>/go-grpc-client:latest
          ports:
            - containerPort: 50051
```
Si nos damos cuenta, en un solo deployment se están definiendo 2 contenedores. Esto habilita la posiblidad (como se mencionó anteriormente) de comunicar ambos contenedores únicamente con `localhost`.

Con esto claro lo primero será desplegar toda nuestra aplicación o como mínimo al API REST de rust para evitar problemas con el Gateway API:

```bash
kubectl apply -f rest_rust.yml
kubectl apply -f go_rest_client.yml
kubectl apply -f go_server_grpc.yml
kubectl apply -f consumer-deployment.yml
```

Finalmente desplegamos nuestro Gateway:
```bash
kubectl apply -f gateway.yml
```

Y para enviar solicitudes obtenemos el ADRESS del Gateway con:
```bash
kubectl get gateway rust-api-gateway
```
