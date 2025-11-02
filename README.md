# Task Manager API (Arquitectura Hexagonal)

API de tareas implementada con Go, Gin y SQLite/libSQL (Turso), siguiendo principios de Arquitectura Hexagonal (puertos y adaptadores).

## Características

- Gestión de tareas (CRUD y filtrado por estado).
- Capa de aplicación y dominio separadas de infraestructura y presentación.
- Conexión local (`modernc.org/sqlite`) o remota (`libSQL` de Turso).
- Endpoints de salud:
  - `GET /health` — estado general de la app.
  - `GET /health/db` — estado de la base de datos (ping + versión).

## Estructura

- `modules/task/domain/` — entidades y puertos de dominio.
- `modules/task/application/` — casos de uso (servicios).
- `modules/task/infrastructure/` — repositorios (adaptadores externos).
- `modules/task/presentation/` — handlers y rutas HTTP.
- `shared/config/` — configuración (`.env`).
- `shared/database/` — conexión SQLite/libSQL.
- `cmd/server/` — wire-up del servidor y DI.

## Configuración

Variables de entorno relevantes (usar `.env`):

```
APP_ENV=development        # development | production
APP_NAME=Task Manager API
SERVER_HOST=localhost
SERVER_PORT=8080

# Base de datos local
DB_PATH=./data/tasks.db

# Base de datos remota (Turso/libSQL)
# DB_URL=libsql://<host>:<port>?insecure=true
# DB_AUTH_TOKEN=<tu_token>
```

- Local: establece `DB_PATH` y deja vacío `DB_URL`.
- Remota: define `DB_URL` y opcionalmente `DB_AUTH_TOKEN`. `DB_PATH` se ignora.

Sugerencia: copia el archivo de ejemplo y ajusta tus valores:

```
cp .env.example .env
```

## Ejecutar

- Desarrollo:
  ```bash
  go run .\cmd\server
  ```
- Producción (logs más silenciosos):
  - Ajusta `APP_ENV=production` en `.env`.
  ```bash
  go run .\cmd\server
  ```

El servidor se inicia en `http://<SERVER_HOST>:<SERVER_PORT>/`.

## Endpoints

- Salud:
  - `GET /health`
  - `GET /health/db`
- Tareas (según handlers ya implementados):
  - `GET /tasks/:id`
  - `GET /tasks`
  - `GET /tasks/status?completed=<true|false>`
  - `POST /tasks`
  - `PUT /tasks/:id`
  - `DELETE /tasks/:id`

## Tests

- Infraestructura (SQLiteTaskRepository):
  ```bash
  go test .\modules\task\infrastructure -v
  ```
- Presentación (handlers):
  ```bash
  go test .\modules\task\presentation\test -v
  ```

Los tests de infraestructura usan SQLite local temporal por prueba (aislado y rápido). Los de presentación mockean el servicio.

## Notas

- Si `8080` está ocupado, usa `SERVER_PORT=8081`.
- Para exposición externa, configura `SERVER_HOST=0.0.0.0`.
- `/health/db` responde `503` si la BD no está disponible.

## Próximos pasos (opcional)

- CI/CD con `go test`, `golangci-lint`.
- `.env.example` ya incluido para facilitar la configuración.
- Migraciones/Seeds si agregas más entidades.

## CI/CD (GitHub Actions)

- Ya se agregó un workflow en `.github/workflows/ci.yml` que:
  - Ejecuta `go test ./... -race` en cada `push` y `pull_request` a `main/master`.
  - Compila el binario del servidor y sube el artefacto.
- Cómo activarlo:
  - Sube este repo a GitHub (origen remoto).
  - Asegúrate de que la rama se llama `main` o `master`.
  - GitHub Actions se activará automáticamente en el siguiente `push`.
- Personalización:
  - Si agregas tests que dependan de una BD remota, usa `Secrets` del repositorio para `DB_URL` y `DB_AUTH_TOKEN` y exporta esas variables en el job antes de correr los tests.
  - Puedes añadir matrices de versiones de Go u OS si lo necesitas.