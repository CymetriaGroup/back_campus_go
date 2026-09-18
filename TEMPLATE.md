# Estado de la implementación del template

Este repositorio incluye una referencia ejecutable de arquitectura hexagonal. Arranca sin dependencias externas mediante adaptadores en memoria; PostgreSQL y Redis aparecen en Compose como puntos de extensión, pero la API no se conecta a ellos todavía.

## Implementado

- Dominio sin etiquetas HTTP/GORM y errores explícitos.
- Puertos de usuarios, autenticación, tokens, hashing, cache y correo.
- CRUD de usuarios, búsqueda, filtros, orden y paginación (máximo 100).
- Login, refresh con rotación, logout, access tokens HMAC con expiración y RBAC.
- bcrypt y cache concurrente en memoria con TTL.
- DTOs, validación Gin y respuestas consistentes.
- Request ID, logs JSON, recovery, CORS, rate limit, límite de body y security headers.
- `/health`, `/ready`, timeouts y graceful shutdown.
- Configuración por entorno con validación fail-fast en producción.
- Migración SQL, OpenAPI, Docker multi-stage, Compose, Makefile y GitHub Actions.
- Pruebas unitarias, HTTP y detector de carreras.

## Puntos de extensión deliberados

- Implementar `ports.UserRepository` con PostgreSQL/GORM o `database/sql`.
- Implementar `ports.Cache` con Redis.
- Sustituir `LogMailer` por SMTP/SES/SendGrid.
- Sustituir el token HMAC compacto por una librería JWT auditada si se requiere interoperabilidad JWT.
- Conectar `/ready` a checks reales de PostgreSQL y Redis.
- Añadir OpenTelemetry y métricas Prometheus según la plataforma de despliegue.

## Inicio rápido

```bash
cp .env.example .env
make run
```

Credenciales de desarrollo por defecto: `admin@example.com` / `admin1234`. No deben usarse en producción.

```bash
curl -s http://localhost:8080/health
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"admin1234"}'
```

Use el `access_token` retornado:

```bash
curl -s http://localhost:8080/api/v1/users \
  -H 'Authorization: Bearer ACCESS_TOKEN'
```

Consulte `docs/swagger/openapi.yaml` para el contrato y `migrations/` para el esquema inicial.
