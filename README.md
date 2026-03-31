# Credit Layer

**[English](#english) | [Español](#español)**

---

## English

### What is Credit Layer?

Credit Layer is an open-source service that acts as a **credits layer** for SaaS products with usage-based pricing. Instead of reimplementing the same ledger logic in every project, you self-host this service and call 3 endpoints from your own backend.

Credit Layer does **not** know anything about payments. You handle your own payment flow (Stripe, Polar, Lemon Squeezy, etc.) and when a payment is confirmed, you call `POST /v1/credit/add`. Credit Layer only manages the ledger.

### The problem it solves

Every developer building a usage-based SaaS has to reinvent the same infrastructure: tracking credits per user, avoiding duplicates, handling race conditions. Credit Layer gives you that infrastructure out of the box.

### How it works

1. You download and self-host Credit Layer on your own infrastructure
2. Create an app to get your API key
3. From your backend, call 3 endpoints:

```
POST /v1/credit/add      → add credits to a user
POST /v1/credit/deduct   → deduct credits from a user
GET  /v1/credit/balance  → get a user's current balance
```

Credit Layer handles the accounting ledger, idempotency, and race conditions internally.

### Requirements

- Go 1.21+
- PostgreSQL 14+

### Quick Start

**1. Clone the repository**

```bash
git clone https://github.com/your-username/credit-layer-back.git
cd credit-layer-back
```

**2. Install dependencies**

```bash
go mod download
```

**3. Configure environment variables**

```bash
cp .example.env .env
```

Edit `.env` with your values:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/credit_layer?sslmode=disable
TEST_DATABASE_URL=postgres://user:password@localhost:5432/credit_layer_test?sslmode=disable
```

> **`sslmode` note:** Use `sslmode=disable` when the app and database are on the same private network. Use `sslmode=require` when connecting to a cloud-hosted database or over the internet.

**4. Run the server**

```bash
go run cmd/main.go
```

Migrations run automatically on startup. The server will be available at `http://localhost:8080`.

### API Reference

All endpoints under `/v1/credit/` require authentication via API key.

**Header:**
```
Authorization: Bearer <your_api_key>
```

---

#### Create an App

```
POST /v1/apps
```

Call this once to register your app and obtain your API key.

**Request body:**
```json
{
  "name": "My App"
}
```

**Response:**
```json
{
  "id": "368f5c53-9bb6-4bb3-952f-a2ff08fac076",
  "name": "My App",
  "api_key": "cl_a1b2c3d4..."
}
```

> Save the `api_key` — it will not be shown again.

---

#### Add Credits

```
POST /v1/credit/add
Authorization: Bearer <api_key>
```

**Request body:**
```json
{
  "user_id": "user_123",
  "amount": 100,
  "description": "Credit purchase",
  "idempotency_key": "stripe_evt_abc123"
}
```

**Response:**
```json
{
  "id": "77312eca-9f51-4083-9076-b40b9b390803",
  "app_id": "368f5c53-9bb6-4bb3-952f-a2ff08fac076",
  "user_id": "user_123",
  "amount": 100,
  "description": "Credit purchase",
  "idempotency_key": "stripe_evt_abc123"
}
```

---

#### Deduct Credits

```
POST /v1/credit/deduct
Authorization: Bearer <api_key>
```

**Request body:**
```json
{
  "user_id": "user_123",
  "amount": 10,
  "description": "Image generation",
  "idempotency_key": "op_xyz789"
}
```

Returns `422 Unprocessable Entity` if the user has insufficient credits.

---

#### Get Balance

```
GET /v1/credit/balance?user_id=user_123
Authorization: Bearer <api_key>
```

**Response:**
```json
{
  "balance": 90
}
```

---

### Field Reference

| Field | Type | Description |
|---|---|---|
| `user_id` | string | The ID of the user in **your** system. Can be any string. |
| `amount` | integer | Number of credits. Always positive — Credit Layer handles the sign internally. |
| `description` | string | Required. Describes the reason for the credit movement. |
| `idempotency_key` | string | Required. A unique key per operation to prevent duplicates. Use the payment event ID or generate a UUID. |

### Docker

**Build the image:**

```bash
docker build -t credit-layer .
```

**Run the container:**

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_URL=postgres://user:password@host:5432/credit_layer?sslmode=disable \
  credit-layer
```

> Never pass a `.env` file into the container. Always inject environment variables via `-e` flags or your orchestrator's secret management (e.g. Kubernetes Secrets, Docker Swarm secrets).

**Security notes:**
- Multi-stage build: the final image contains only the binary and migrations — no Go toolchain, no source code
- Runs as a non-root user (`appuser`)
- Static binary with debug info stripped (`-ldflags="-w -s"`)
- Based on `alpine:3.21` for a minimal attack surface

---

### Running Tests

**Unit tests:**
```bash
go test ./internal/application/...
```

**Integration tests** (requires a running PostgreSQL instance with `credit_layer_test` database):
```bash
go test ./internal/infrastructure/postgres/...
```

**All tests:**
```bash
go test ./...
```

### Architecture

Credit Layer follows a hexagonal architecture:

```
cmd/
  main.go                         entry point, wires everything together
internal/
  domain/                         entities and repository interfaces
  application/                    use cases (business logic)
  infrastructure/postgres/        sqlc-generated DB access
  interfaces/http/                Gin HTTP handlers
db/
  migrations/                     goose migrations
  queries/                        sqlc SQL queries
```

### Tech Stack

| Layer | Technology |
|---|---|
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | PostgreSQL |
| DB Driver | [pgx/v5](https://github.com/jackc/pgx) |
| Query Generation | [sqlc](https://sqlc.dev) |
| Migrations | [goose](https://github.com/pressly/goose) |
| Testing | [testify](https://github.com/stretchr/testify) + [mockery](https://github.com/vektra/mockery) |

---

## Español

### ¿Qué es Credit Layer?

Credit Layer es un servicio open-source que actúa como **capa de créditos** para productos SaaS con precios por uso. En lugar de reimplementar la misma lógica de ledger en cada proyecto, self-hosteás este servicio y llamás 3 endpoints desde tu propio backend.

Credit Layer **no** sabe nada de pagos. Vos manejás tu propio flujo de pagos (Stripe, Polar, Lemon Squeezy, etc.) y cuando se confirma un pago, llamás `POST /v1/credit/add`. Credit Layer solo administra el ledger.

### El problema que resuelve

Todo desarrollador que construye un SaaS con precios por uso tiene que reinventar la misma infraestructura: llevar créditos por usuario, evitar duplicados, manejar race conditions. Credit Layer te da esa infraestructura lista para usar.

### Cómo funciona

1. Descargás y self-hosteás Credit Layer en tu propia infraestructura
2. Creás una app para obtener tu API key
3. Desde tu backend, llamás 3 endpoints:

```
POST /v1/credit/add      → agregar créditos a un usuario
POST /v1/credit/deduct   → descontar créditos de un usuario
GET  /v1/credit/balance  → obtener el balance actual de un usuario
```

Credit Layer maneja internamente el ledger contable, la idempotencia y las race conditions.

### Requisitos

- Go 1.21+
- PostgreSQL 14+

### Inicio Rápido

**1. Clonar el repositorio**

```bash
git clone https://github.com/your-username/credit-layer-back.git
cd credit-layer-back
```

**2. Instalar dependencias**

```bash
go mod download
```

**3. Configurar variables de entorno**

```bash
cp .example.env .env
```

Editá `.env` con tus valores:

```env
PORT=8080
DATABASE_URL=postgres://usuario:password@localhost:5432/credit_layer?sslmode=disable
TEST_DATABASE_URL=postgres://usuario:password@localhost:5432/credit_layer_test?sslmode=disable
```

> **Nota sobre `sslmode`:** Usá `sslmode=disable` cuando la app y la base de datos están en la misma red privada. Usá `sslmode=require` cuando te conectás a una base de datos en la nube o a través de internet.

**4. Correr el servidor**

```bash
go run cmd/main.go
```

Las migraciones corren automáticamente al arrancar. El servidor estará disponible en `http://localhost:8080`.

### Referencia de API

Todos los endpoints bajo `/v1/credit/` requieren autenticación via API key.

**Header:**
```
Authorization: Bearer <tu_api_key>
```

---

#### Crear una App

```
POST /v1/apps
```

Llamá este endpoint una vez para registrar tu app y obtener tu API key.

**Body:**
```json
{
  "name": "Mi App"
}
```

**Respuesta:**
```json
{
  "id": "368f5c53-9bb6-4bb3-952f-a2ff08fac076",
  "name": "Mi App",
  "api_key": "cl_a1b2c3d4..."
}
```

> Guardá el `api_key` — no se mostrará nuevamente.

---

#### Agregar Créditos

```
POST /v1/credit/add
Authorization: Bearer <api_key>
```

**Body:**
```json
{
  "user_id": "user_123",
  "amount": 100,
  "description": "Compra de créditos",
  "idempotency_key": "stripe_evt_abc123"
}
```

**Respuesta:**
```json
{
  "id": "77312eca-9f51-4083-9076-b40b9b390803",
  "app_id": "368f5c53-9bb6-4bb3-952f-a2ff08fac076",
  "user_id": "user_123",
  "amount": 100,
  "description": "Compra de créditos",
  "idempotency_key": "stripe_evt_abc123"
}
```

---

#### Descontar Créditos

```
POST /v1/credit/deduct
Authorization: Bearer <api_key>
```

**Body:**
```json
{
  "user_id": "user_123",
  "amount": 10,
  "description": "Generación de imagen",
  "idempotency_key": "op_xyz789"
}
```

Retorna `422 Unprocessable Entity` si el usuario no tiene créditos suficientes.

---

#### Obtener Balance

```
GET /v1/credit/balance?user_id=user_123
Authorization: Bearer <api_key>
```

**Respuesta:**
```json
{
  "balance": 90
}
```

---

### Referencia de Campos

| Campo | Tipo | Descripción |
|---|---|---|
| `user_id` | string | El ID del usuario en **tu** sistema. Puede ser cualquier string. |
| `amount` | integer | Cantidad de créditos. Siempre positivo — Credit Layer maneja el signo internamente. |
| `description` | string | Requerido. Describe el motivo del movimiento de créditos. |
| `idempotency_key` | string | Requerido. Una clave única por operación para evitar duplicados. Usá el ID del evento de pago o generá un UUID. |

### Docker

**Construir la imagen:**

```bash
docker build -t credit-layer .
```

**Correr el contenedor:**

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_URL=postgres://usuario:password@host:5432/credit_layer?sslmode=disable \
  credit-layer
```

> Nunca copies el archivo `.env` dentro del contenedor. Siempre inyectá las variables de entorno via flags `-e` o el sistema de secretos de tu orquestador (Kubernetes Secrets, Docker Swarm secrets, etc.).

**Medidas de seguridad aplicadas:**
- Build multi-stage: la imagen final solo contiene el binario y las migraciones — sin Go toolchain, sin código fuente
- Corre con usuario no-root (`appuser`)
- Binario estático con debug info eliminada (`-ldflags="-w -s"`)
- Basado en `alpine:3.21` para minimizar la superficie de ataque

---

### Correr Tests

**Tests unitarios:**
```bash
go test ./internal/application/...
```

**Tests de integración** (requiere una instancia de PostgreSQL con la base de datos `credit_layer_test`):
```bash
go test ./internal/infrastructure/postgres/...
```

**Todos los tests:**
```bash
go test ./...
```

### Arquitectura

Credit Layer sigue una arquitectura hexagonal:

```
cmd/
  main.go                         punto de entrada, conecta todo
internal/
  domain/                         entidades e interfaces de repositorio
  application/                    casos de uso (lógica de negocio)
  infrastructure/postgres/        acceso a DB generado por sqlc
  interfaces/http/                handlers HTTP con Gin
db/
  migrations/                     migraciones de goose
  queries/                        queries SQL de sqlc
```

### Stack Tecnológico

| Capa | Tecnología |
|---|---|
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| Base de datos | PostgreSQL |
| Driver DB | [pgx/v5](https://github.com/jackc/pgx) |
| Generación de queries | [sqlc](https://sqlc.dev) |
| Migraciones | [goose](https://github.com/pressly/goose) |
| Testing | [testify](https://github.com/stretchr/testify) + [mockery](https://github.com/vektra/mockery) |

