# Nebula Challenge - TLS Security Scanner API

API RESTful desarrollada en **Go + Gin** que permite escanear la seguridad TLS/SSL de dominios utilizando la API de **SSL Labs** (v2), filtrar los resultados relevantes y almacenarlos en **MongoDB**.

Proporciona endpoints para iniciar escaneos asíncronos, consultar estado en tiempo real y obtener reportes filtrados con calificación, protocolos, cifrado, HSTS, certificado y veredicto de seguridad.

## Características principales

- Escaneo asíncrono de dominios con polling (no bloquea la respuesta)
- Filtrado inteligente del reporte SSL Labs (más de 2000 líneas → ~10-20 campos útiles)
- Veredicto claro de seguridad (Excelente / Buena / Aceptable / Deficiente / Muy mala)
- Almacenamiento en MongoDB de reportes filtrados
- Soporte para agregaciones avanzadas vía endpoint `/aggregate`
- Manejo robusto de errores y validaciones
- Concurrencia segura con mutex

## Tecnologías utilizadas

- **Go** (1.25.6)
- **Gin** - Framework web ligero y rápido
- **gjson** - Parsing ultrarrápido y eficiente de JSON grande
- **MongoDB** + **mongo-driver**
- **strings.Builder** para generación eficiente de summaries
- **UUID** para IDs de escaneo

## Ejecutar la aplicación y las pruebas

### Requisitos previos

- Go 1.25.6 o superior
- MongoDB (local o Atlas – recomendado Atlas gratuito para pruebas)
- Node.js 18+ (solo para ejecutar los tests)
- Git

### 1. Clonar el repositorio

```bash
git clone https://github.com/Garffieldd/Nebula-Challenge.git
cd Nebula-Challenge
```


### 2. Configuración de variables de entorno (MongoDB)

Para utilizar los endpoints que interactúan con la base de datos, es necesario crear un archivo .env dentro del directorio backend/ con la siguiente estructura de ejemplo:

MONGO_USER=nebula_user
MONGO_PASSWORD=supersecret123
MONGO_HOST=cluster0.xxxxx.mongodb.net
MONGO_DB=nebula_tls

### 3. Ejecutar el backend

```bash
cd Nebula-Challengue/backend
go mod tidy        # Descarga dependencias (solo la primera vez)
go run .
```

La API quedará disponible por defecto en: http://localhost:8080

### 4. Ejecutar la prueba ( si se desea )

```bash
cd Nebula-Challengue/backend
npm i
cd test
node test.js #Es necesario que el backend esté corriendo antes de ejecutar las pruebas.
```



## Endpoints disponibles

### Endpoints de información de dominios (CRUD + agregación)

| Método | Endpoint                           | Descripción                                                                                 | Body / Params                              |
|--------|------------------------------------|---------------------------------------------------------------------------------------------|--------------------------------------------|
| GET    | `/domains-info`                    | Obtiene todos los registros de dominios escaneados                                          | -                                          |
| GET    | `/domains-info/:id`                | Obtiene un registro específico por su ID                                                    | `:id` (ObjectID de MongoDB)                |
| POST   | `/create-domain-info`              | Crea un nuevo registro de información de dominio (manual o para pruebas)                    | JSON con estructura FilteredTLSReport      |
| POST   | `/domains-info/aggregate`          | Ejecuta una agregación personalizada en MongoDB (pipeline flexible)                         | Array de etapas MongoDB Aggregation        |
| DELETE | `/domains-info/:id`                | Elimina un registro por su ID                                                               | `:id` (ObjectID de MongoDB)                |

### Endpoints de escaneo TLS con SSL Labs

| Método | Endpoint                        | Descripción                                                                                 | Body / Params                              |
|--------|---------------------------------|---------------------------------------------------------------------------------------------|--------------------------------------------|
| POST   | `/start-scan`                   | Inicia un nuevo escaneo TLS asíncrono para un dominio                                       | `{ "domain": "www.ejemplo.com" }`          |
| GET    | `/scan-status/:scanRequestID`   | Consulta el estado del escaneo y, cuando esté completo, devuelve el reporte filtrado       | `:scanRequestID` (UUID devuelto por /start-scan) |

**Ejemplo de respuesta completa en `/scan-status/:id` cuando termina:**
```json
{
  "status": "complete",
  "report": {
    "host": "www.ejemplo.com",
    "webProtocol": "https",
    "endpoints": [ ... ],
    "summary": "Análisis TLS para www.ejemplo.com - Calificación general: A ...",
    "timestamp": "2026-01-19T11:20:00Z"
  }
}



