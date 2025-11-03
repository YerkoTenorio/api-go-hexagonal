# Arquitectura y Estructura del Proyecto

Este documento resume la arquitectura, convenciones y estructura de carpetas de esta API en Go, pensada para ser reutilizable como plantilla base en futuros proyectos.

## Principios

- Arquitectura Hexagonal (Puertos y Adaptadores): separa dominio/aplicación de infraestructura/presentación.
- Capas con responsabilidades claras y dependencias dirigidas (de afuera hacia adentro): presentación  aplicación  dominio.
- Inversión de dependencias: el dominio define puertos; la infraestructura implementa adaptadores.
- Configuración externa por variables de entorno (.env) y build reproducible.

## Estructura de Carpetas

```
.
 cmd/
    server/           # Entry point del servidor (wire-up, DI, router)
        data/         # Recursos del binario (si aplica)
       - main.go
 modules/
    task/             # Módulo de tareas (ejemplo)
        domain/       # Entidades y puertos (interfaces del repositorio)
        application/  # Casos de uso/servicios (orquestación de lógica)
        infrastructure/# Adaptadores externos (repositorios, DB, etc.)
        presentation/ # HTTP handlers y rutas (Gin)
 shared/
    config/           # Carga de configuración (.env)
    database/         # Conexión SQLite/libSQL (Turso)
 pkg/
    errors/           # Utilidades de error (opcional)
 .github/workflows/    # CI/CD con GitHub Actions
 data/                 # Archivos de datos locales (SQLite)
 .env.example          # Variables de entorno de ejemplo
 README.md             # Guía general del proyecto
 go.mod / go.sum       # Módulo de Go y dependencias
```

## Capas y Responsabilidades

- Dominio (`modules/<feature>/domain`)
  - Entidades (p. ej., `Task`).
  - Puertos (interfaces) para operaciones de repositorio, sin detalles de implementación.
- Aplicación (`modules/<feature>/application`)
  - Servicios (casos de uso) que coordinan reglas de negocio usando puertos del dominio.
  - No deben conocer detalles de base de datos ni HTTP.
- Infraestructura (`modules/<feature>/infrastructure`)
  - Implementa los puertos del dominio (p. ej., repositorio SQLite/libSQL).
  - Aquí viven las dependencias externas (drivers, clientes, etc.).
- Presentación (`modules/<feature>/presentation`)
  - Handlers y rutas HTTP con Gin.
  - Traduce HTTP  llamadas a servicios de aplicación.
- `cmd/server`
  - Ensamblaje (wiring/DI): crea dependencias y las conecta.
  - Inicializa router, middlewares y rutas.
- `shared/config` y `shared/database`
  - Carga de configuración desde `.env`.
  - Factoría de conexiones (SQLite local con `modernc.org/sqlite` o remota con libSQL/Turso).

## Configuración por Entorno

- Archivo `.env` (usa `.env.example` como base):
  - `APP_ENV`, `APP_NAME`, `SERVER_HOST`, `SERVER_PORT`.
  - Local: `DB_PATH`.
  - Remota: `DB_URL` y opcional `DB_AUTH_TOKEN`.
- Estrategia:
  - Si `DB_URL` está definido, se usa libSQL/Turso; `DB_PATH` se ignora.
  - Si no hay `DB_URL`, se usa SQLite local en `data/*.db`.

## Base de Datos

- Local: SQLite (driver `modernc.org/sqlite`).
- Remota: libSQL (Turso). Requiere URL y, si aplica, token.
- `shared/database/sqlite.go` contiene la lógica para abrir la conexión adecuada.

## Endpoints Clave

- Salud:
  - `GET /health`  estado general de la app.
  - `GET /health/db`  ping + versión de la base de datos.
- Tareas (ejemplo de módulo): CRUD y filtrado por estado.

## Testing

- Dominio: tests puros de entidades y reglas.
- Aplicación: tests de servicios con repositorios mock (unitarios).
- Presentación: tests de handlers con servicio mock, validan respuestas HTTP.
- Infraestructura: tests del repositorio con SQLite temporal por prueba (aislados y rápidos).
- Comandos:
  - `go test ./modules/task/infrastructure -v`
  - `go test ./modules/task/presentation/test -v`
  - `go test ./... -race` (CI)

## CI/CD (GitHub Actions)

- Workflow en `.github/workflows/ci.yml`:
  - Dispara en `push` y `pull_request` a `main`.
  - Pasos: checkout, setup Go, cache módulos, `go test ./... -race`, `go build -o server ./cmd/server`, subir artefacto.
- Personalización:
  - Secrets para BD remota (`DB_URL`, `DB_AUTH_TOKEN`).
  - Matriz de versiones de Go/OS si se requiere.
  - Añadir `golangci-lint` y reporte de cobertura.

## Convenciones

- Nombres en kebab-case para repos y directorios top-level.
- Módulos por feature bajo `modules/<feature>`.
- Imports alineados con `go.mod` (`module github.com/<usuario>/<repo>`).
- Handlers enfocan solo capa HTTP; la lógica vive en servicios de aplicación.

## Guía para Crear un Nuevo Módulo (Plantilla)

1. Dominio
   - Crear entidad: `modules/<feature>/domain/<entity>.go`.
   - Definir puertos (interfaces): `modules/<feature>/domain/repository.go`.
2. Aplicación
   - Implementar servicios/casos de uso: `modules/<feature>/application/service.go`.
   - Definir interfaces de servicio si se mockearán en presentación.
3. Infraestructura
   - Implementar repositorio: `modules/<feature>/infrastructure/repository.go`.
   - Usar `shared/database` para acceder a la conexión.
4. Presentación
   - Crear handlers: `modules/<feature>/presentation/handler.go`.
   - Rutas: `modules/<feature>/presentation/routes.go`.
   - Registrar en `cmd/server/main.go`.
5. Wiring en `cmd/server`
   - Construir repositorio (infra) y servicio (aplicación).
   - Inyectarlos en handlers (presentación) y registrar rutas.
6. Tests
   - Dominio/aplicación: unitarios con mocks.
   - Infraestructura: integración con SQLite temporal.
   - Presentación: handlers con servicio mock.

## Comandos Útiles

- Desarrollo: `go run ./cmd/server`
- Tests: `go test ./... -race`
- Build: `go build -o server ./cmd/server`

## Consideraciones de Escalabilidad

- Añade nuevos módulos bajo `modules/` siguiendo el mismo patrón.
- Centraliza configuración en `shared/config` y conexiones externas en `shared/database`.
- Usa interfaces en el dominio para mantener independencia de infraestructura.
- Mantén los handlers delgados; la lógica empresarial vive en servicios.

---

Esta guía sirve como referencia práctica para replicar la arquitectura y estructura en nuevos proyectos. Ajusta nombres de módulos, entidades y adaptadores según tu dominio.
