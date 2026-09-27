# TableroAgente

![Vista previa AgentBoard](../../assets/social-preview.png)

**[Descargar para Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Sitio de documentación](https://bogdanovandreycode.github.io/AgentBoard/) · [Licencia MIT](../../LICENSE)**

AgentBoard es un tablero de tareas local donde los trabajadores humanos y de IA trabajan en tareas comunes, pero tienen derechos diferentes. La aplicación se inicia con un archivo `agentboard.exe`, abre la interfaz web en el navegador y proporciona a los trabajadores un servidor MCP independiente a través de `stdio`. Los datos permanecen en su computadora.

**[Comenzar desde cero](START_HERE.md) · [Trabajar con tareas](TASKS.md) · [Conectar AI a través de MCP](WORKERS_MCP.md) · [Importar JSON](IMPORT.md) · [Configuración](SETTINGS.md) · [Resolver problemas](TROUBLESHOOTING.md)**

## Funciones de MVP

- Tablero de proyecto local con búsqueda de tareas, importación JSON, columnas y propiedades personalizadas.
- Acceso MCP para un trabajador específico con transiciones de IA limitadas y aceptación humana final de las tareas.
- Historial general de tareas, instrucciones de verificación, artefactos, costos de IA y diagnóstico de conexión de trabajadores.
- Instalador de Windows, manifiesto ZIP y Scoop; la interfaz web está integrada en el archivo ejecutable.

AgentBoard está diseñado para un usuario local confiable. La aplicación no aloja proyectos en la nube y no lanza clientes de IA por sí misma; si es necesario, conecte un cliente compatible con MCP al trabajador.

## En cinco minutos

1. Descargue el instalador `agentboard-VERSION-windows-amd64-setup.exe` desde [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Sugerirá una carpeta (por defecto `C:\AI\AgentBoard`) y la agregará a `PATH`. Scoop y ZIP también están disponibles.
2. Abra PowerShell en la carpeta de su proyecto, por ejemplo `C:\Projects\MyApp`.
3. Ejecute `agentboard init` (para ZIP: ruta completa a `agentboard.exe` y `init`).
4. Ejecute `agentboard open`. Se abrirá `http://127.0.0.1:7337`.
5. Agregue una tarea usando el botón **Nueva tarea**. Para un trabajador de IA, abra **Trabajadores → Agregar trabajador**, seleccione el perfil del cliente y copie la configuración de MCP.

Si aún no tiene una carpeta de proyecto, cree una en el Explorador de Windows. Un proyecto puede ser cualquier carpeta, incluso sin Git ni código.

## Instalación mediante Scoop

En PowerShell con [Scoop](https://scoop.sh/)] ya instalado después del lanzamiento:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Para los desarrolladores existe [compilación desde la fuente ](INSTALL.md). El flujo de trabajo de lanzamiento crea un instalador, un ZIP y un manifiesto Scoop con SHA-256 a partir del mismo artefacto. Actualización de la versión instalada a través de Scoop: `scoop update agentboard` después de agregar el manifiesto al depósito; detalles - [preparación de lanzamiento](SCOOP_RELEASE.md).

## Cómo está estructurado el tablero

`Backlog → Features → In progress → Testing → Verification → Complete`

Una persona puede mover tareas por el tablero. La IA sólo puede mover `Features → In progress → Testing → Verification`; La aceptación final en `Complete` la realiza un humano. Las columnas de usuario están destinadas a humanos: la tarea que contienen permanece en el estado `Backlog` para MCP. Modos de prueba: IA, Humano e Híbrido. La tarea incluye historia, pruebas, artefactos y costos de IA.

## Equipos

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` registra la carpeta y escribe allí solo `.agentboard/project.json`. Los datos de producción de SQLite se encuentran en el directorio de configuración de usuario de Windows (`%AppData%\AgentBoard\agentboard.db`), fuera del proyecto y fuera de la instalación de Scoop. Eliminar o actualizar el paquete no debería eliminar estos datos. Antes de transferir a otra computadora, haga una copia de la base de datos mientras AgentBoard está detenido.

Trabajadores - cuentas lógicas; AgentBoard en sí no ejecuta Codex, Claude ni ningún otro cliente de IA. El cliente inicia un proceso MCP local para un trabajador específico. MCP es el límite de permisos de la aplicación y, para aislar archivos, utilice el entorno limitado de pruebas del cliente AI.

## Para desarrolladores

Pila: Go, SQLite, SDK oficial de MCP Go, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Primero cree la interfaz y luego vaya: `./scripts/build.ps1`. Los archivos web se incluyen en el binario a través de `go:embed`. La arquitectura y la API se describen en [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Licencia

AgentBoard es un proyecto gratuito y de código abierto bajo [licencia MIT](../../LICENSE). Se permiten el uso comercial, la modificación, la bifurcación y la redistribución siempre que se mantengan el aviso y la licencia de derechos de autor.
