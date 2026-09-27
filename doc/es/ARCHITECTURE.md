# Arquitectura y API

AgentBoard es un proceso Go local que ejecuta SQLite. La API HTTP para humanos, el adaptador MCP para IA y la interfaz web comparten una lógica de servicio común. SQLite - fuente del estado; La interfaz no pasa por alto el servidor. Las herramientas de IA tienen una superficie de derechos separada y el estado `Backlog`/`Complete` en MCP no está disponible. El registro histórico de `System` lo crea la propia aplicación.

## Componentes

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - modelos de dominio y errores.
- `internal/service` - transiciones de tareas, permisos, importación y configuración.
- `internal/persistence` - SQLite y migraciones.
- `internal/httpapi` - API HTTP humana.
- `internal/mcpserver` - Herramientas MCP para un trabajador independiente.
- `web` - Interfaz de usuario de React/TypeScript; el ensamblaje termina en `internal/webui/dist` y se incluye en el EXE.

## Rutas HTTP básicas

| Método y camino | Destino |
| --- | --- |
| `GET /api/health` | Comprobación del servidor. |
| `GET /api/projects` | Proyectos registrados. |
| `GET /api/projects/{id}/board` | Tablero de proyecto. |
| `GET/PUT /api/projects/{id}/settings` | Configuración, orden y nombres de columnas. |
| `POST /api/projects/{id}/tasks` | Crea una tarea. |
| `POST /api/projects/{id}/tasks/import` | Importar atómicamente JSON versión 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Tarjeta, cambiar, borrar. |
| `POST /api/tasks/{id}/move` | Moviéndose por una persona. |
| `GET/POST /api/projects/{id}/workers` | Listado y creación de un trabajador. |
| `GET /api/workers/{id}/mcp/check` | Pruebas internas del protocolo de enlace y las herramientas de MCP. |
| `GET /api/workers/{id}/sessions` | Diagnóstico de sesión. |
| `GET/POST /api/projects/{id}/properties` | Propiedades personalizadas. |

HTTP es para el usuario de confianza local. No publique un puerto web en Internet sin su propia autenticación, restricciones de red y HTTPS. El servidor MCP se inicia utilizando `stdio` para un proyecto y trabajador específicos; comience con `get_my_board`. Sus permisos restringen las acciones dentro de AgentBoard, pero no reemplazan el entorno limitado del sistema de archivos del cliente AI.

Las columnas de usuario se almacenan por separado de `tasks.state`: el núcleo mantiene la tarea en dicha columna en `backlog`, y `tasks.board_column` determina el lugar en el tablero humano. Esto mantiene el mismo modelo de transición de IA. Al eliminar una columna se borran las tareas `board_column`.
