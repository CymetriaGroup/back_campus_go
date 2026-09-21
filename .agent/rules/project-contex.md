# Contexto y reglas arquitectónicas del proyecto

## Enfoque general

Este backend es un monolito modular en Go. Cada capacidad de negocio se organiza como un módulo con arquitectura hexagonal propia. No se deben volver a crear capas globales que mezclen entidades, servicios, puertos o repositorios de módulos diferentes.

Estructura base de un módulo:

```text
internal/modules/<module>/
├── domain/
├── application/
├── infrastructure/
└── delivery/
```

Los módulos representan capacidades o bounded contexts del negocio, no tablas ni operaciones CRUD individuales. Ejemplos: `users`, `auth`, `tenants`, `campuses`, `enrollment`, `scheduling` y `billing`.

## Domain

- Contiene entidades, value objects, reglas, eventos y errores propios del módulo.
- Debe ser Go puro y no depender de Gin, GORM, PostgreSQL, JSON ni detalles de transporte.
- No usar tags `gorm`, `json`, `binding`, `form` o `uri` en entidades de dominio.
- El dominio nunca importa `application`, `infrastructure`, `delivery` ni otro adaptador.

## Application

- Contiene casos de uso, comandos, consultas, resultados y puertos requeridos por el módulo.
- Los puertos se definen del lado consumidor, no junto a su implementación técnica.
- Los casos de uso dependen de abstracciones y del dominio, nunca de GORM, Gin o repositorios concretos.
- Los DTO de aplicación no tienen versionado HTTP ni tags de Gin.
- Se puede separar en `commands/`, `queries/`, `ports/` o archivos individuales cuando el tamaño real del módulo lo justifique.

## Infrastructure

- Contiene adaptadores de salida: persistencia, proveedores, correo, cache, seguridad y clientes externos.
- Las implementaciones PostgreSQL/GORM viven en `infrastructure/persistence/postgres/`.
- Los modelos GORM viven en `infrastructure/persistence/postgres/models/` y nunca sustituyen a las entidades de dominio.
- Los mappers entre modelos GORM y entidades viven junto al adaptador, normalmente en `infrastructure/persistence/postgres/mappers/`.
- Ningún otro módulo debe importar modelos GORM, mappers o repositorios concretos de este módulo.
- Las migraciones SQL versionadas en `migrations/` son la fuente de verdad del esquema. No usar `AutoMigrate` como mecanismo de despliegue.

## Delivery

- Contiene adaptadores de entrada, por ejemplo HTTP, gRPC, CLI o consumers.
- Los contratos HTTP versionados viven en `delivery/http/v1/`, `delivery/http/v2/`, etc.
- Requests, responses, controllers/handlers y rutas HTTP pertenecen a Delivery, no a Application.
- Los tipos con tags `json`, `binding`, `form` o `uri` deben mantenerse en Delivery.
- Un controller traduce `HTTP request -> input de aplicación` y `resultado de aplicación -> HTTP response`.
- Se puede usar `controller.go` en módulos pequeños y dividir en `controllers/`, `requests/` y `responses/` cuando el tamaño lo justifique.

## Dependencias entre módulos

- Un módulo no puede importar `infrastructure` ni `delivery` de otro módulo.
- Un módulo puede consumir contratos públicos mínimos del `domain` o `application` de otro módulo cuando la dependencia sea necesaria y explícita.
- Preferir puertos de aplicación o eventos de integración para reducir acoplamiento entre módulos.
- No acceder directamente a tablas pertenecientes a otro módulo.
- No compartir DTO HTTP, modelos GORM ni repositorios concretos entre módulos.
- Evitar dependencias circulares. Si aparecen, revisar los límites del dominio o introducir un contrato/evento explícito.

## Platform y Shared

- `internal/platform/` contiene infraestructura transversal y composición: configuración, servidor, router raíz, conexión de base de datos, logging, observabilidad y middleware global.
- Platform puede ensamblar módulos, pero no contiene reglas de negocio.
- `internal/shared/` solo debe existir para conceptos realmente transversales y estables.
- No mover entidades de negocio a `shared` solo porque dos módulos las utilicen; primero definir qué módulo es propietario del concepto.
- Mantener el shared kernel mínimo para evitar convertirlo en un nuevo núcleo acoplado.

## Persistencia con GORM

- PostgreSQL es la base de datos objetivo y GORM es el ORM.
- Entidad de dominio y modelo de persistencia son tipos diferentes.
- Los repositorios PostgreSQL convierten entre ambos mediante mappers explícitos.
- Los puertos de repositorio reciben y retornan tipos de dominio, nunca `*gorm.DB` ni modelos de persistencia.
- La conexión GORM se configura una sola vez en Platform y se inyecta en los repositorios de cada módulo.
- Las transacciones deben exponerse mediante abstracciones para que Application no dependa directamente de GORM.

## Convenciones para cambios futuros

- Antes de agregar una funcionalidad, identificar el módulo propietario.
- No crear paquetes globales `domain`, `services`, `ports`, `models` o `repositories` que mezclen módulos.
- Mantener imports apuntando hacia adentro: Delivery e Infrastructure dependen de Application/Domain; Domain no depende de capas externas.
- Favorecer cohesión por módulo sobre organización horizontal global.
- No crear un módulo por endpoint, caso CRUD o tabla; modularizar por capacidad del negocio.
- Evitar abstracciones vacías y carpetas sin responsabilidad concreta, pero preservar desde el inicio los límites modulares conocidos del sistema Campus.
