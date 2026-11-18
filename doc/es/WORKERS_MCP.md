# Trabajadores y conexión MCP

## Paso 1. Crear un trabajador

En AgentBoard, abra **Trabajadores → Agregar trabajador**. Seleccione un perfil de cliente: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama vía OpenCode, VS Code Copilot u otro cliente MCP. El perfil completa las características típicas (código, pruebas, Git); puedes cambiarlos. El nombre es visible en el tablero y `Slug` es un identificador corto sin espacios para el comando MCP. Haga clic en **Guardar**.

Un trabajador corresponde a una personalidad de IA. Cree diferentes trabajadores para diferentes clientes o equipos. El perfil de capacidad describe la especialización, pero no extiende los derechos de la IA a las etapas de la tarea.

## Paso 2. Copiar la configuración

Abra el trabajador creado. El bloque **MCP diagnostics** muestra el fragmento de configuración y el archivo donde agregarlo. Haga clic en **Copiar configuración de MCP**. Si el archivo ya existe, agregue el servidor sugerido al objeto `mcpServers`/`servers`/`mcp` existente sin borrar los otros servidores.

El comando principal se ve así:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` del trabajador creado. El cliente MCP ejecuta este comando por sí mismo cuando necesita herramientas.

### MCP funciona sin una interfaz web

`agentboard mcp` es un proceso estándar independiente que abre la base de datos AgentBoard. No necesita `agentboard open`/`serve`, no depende de la carpeta desde la que se inicia la interfaz web ni de cuántas instancias de la interfaz web se están ejecutando. Se pueden ejecutar varios agentes simultáneamente, cada uno en su propio proyecto: cada cliente ejecuta su propio proceso MCP y todos utilizan de forma segura la misma base de datos SQLite.

### Cómo se selecciona un proyecto

El servidor MCP se ejecuta en cualquier carpeta. La primera vez que llama a la herramienta, está vinculada al proyecto:

1. `--project PATH`, si se especifica (cualquier carpeta dentro del proyecto registrado es adecuada);
2. en caso contrario, las carpetas de trabajo del cliente (raíces), si el cliente las informa;
3. de lo contrario, la carpeta en la que el cliente inició el servidor MCP.

`get_my_board` devuelve el proyecto y el trabajador seleccionados a `session` junto con `binding_source`. Si las carpetas de trabajo del cliente cambian, la siguiente llamada a la herramienta realiza el enlace nuevamente. Sin `--worker`, se utiliza el único trabajador habilitado del proyecto.

Si la vinculación no es posible, cada herramienta devuelve un error que el agente debe transmitir al usuario; Los datos de la tarea no se leen en este caso:

| Código | Significado |
| --- | --- |
| `NO_PROJECTS` | No hay proyectos en la base de datos. Ejecute `agentboard init` en la carpeta del proyecto y cree un trabajador. No es necesario reiniciar. |
| `PROJECT_REQUIRED` | La carpeta de trabajo está fuera de los proyectos registrados. El mensaje enumera los proyectos. |
| `PROJECT_NOT_FOUND` | `--project` no está dentro de un proyecto registrado. |
| `PROJECT_AMBIGUOUS` | Las carpetas de trabajo de un cliente abarcan varios proyectos. Agregue `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Cree, habilite o seleccione un trabajador a través de `--worker`. |

Las configuraciones para **Codex** y **Gemini CLI** son personalizadas y compartidas en todas las carpetas, por lo que no tienen `--project`: cree un trabajador con el mismo slug en cada proyecto donde se ejecuta este cliente. Los archivos a nivel de proyecto (Claude Code, Cursor, OpenCode, VS Code) contienen `--project` de su proyecto.

### Diseño explícito y verificación de contexto.

Si se especifica `--project` y la carpeta de inicio o la carpeta de trabajo del cliente está en un proyecto diferente, las herramientas solicitan confirmación mediante la obtención de MCP. Los clientes sin obtención de formulario reciben `PROJECT_CONTEXT_MISMATCH` antes de acceder a las tareas. Corrija la configuración, elimine `--project` para la detección automática o confirme explícitamente que tiene la intención de trabajar con otro proyecto con el indicador `--allow-project-mismatch`. Al cambiar las carpetas de trabajo se restablecen las confirmaciones anteriores. Los proyectos definidos automáticamente no requieren confirmación.

**Configuración global → Contexto del proyecto MCP** controla estas comprobaciones para todos los proyectos y trabajadores de la misma base de datos. Habilitado de forma predeterminada. Si está deshabilitado, MCP utiliza el proyecto especificado sin confirmaciones ni errores de discrepancia. Las sesiones ya abiertas aplican la configuración la próxima vez que se llama a la herramienta; Se siguen siguiendo las reglas de visibilidad de tareas, transiciones y pruebas humanas.

Después de la verificación, el bloque de diagnóstico muestra la ruta absoluta al `agentboard.exe` en ejecución. Es especialmente útil al instalar desde un ZIP. Al instalar a través de Scoop, puede usar el comando `agentboard` si el cliente ve el mismo `PATH`.

## Paso 3: Agregue un servidor a su cliente

La interfaz del trabajador ya tiene un fragmento listo para usar. A continuación se muestra una explicación de dónde se utiliza:

| Cliente | Dónde insertar | Cómo comprobar del lado del cliente |
| --- | --- | --- |
| Códice | `%USERPROFILE%\.codex\config.toml`, sección `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Código Claude | `.mcp.json` en la carpeta del proyecto | `claude mcp list` |
| Géminis CLI | `%USERPROFILE%\.gemini\settings.json`, objeto `mcpServers` | `/mcp list` en Gemini CLI |
| Cursores | Proyecto `.cursor\mcp.json` | lista de servidores MCP en la configuración del cursor |
| Código abierto | Proyecto `opencode.json`, objeto `mcp` | lista de herramientas MCP en OpenCode |
| Copiloto de código VS | Proyecto `.vscode\mcp.json`, objeto `servers` | comando **MCP: Listar servidores** |

Ejemplos de Codex (configuración de usuario, proyecto detectado automáticamente) y Claude Code (archivo de proyecto) con el proyecto `C:\Projects\MyFirstProject` y el trabajador `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

En JSON, la barra invertida de Windows se duplica; un fragmento listo para usar de la interfaz hace esto automáticamente. Las configuraciones generadas contienen el `--db` real del servidor web, por lo que HTTP y MCP utilizan la misma base, incluso si el entorno del cliente es diferente.

**Ollama** proporciona un modelo local, pero no reemplaza al cliente MCP. El perfil **Ollama vía OpenCode** genera una configuración MCP para OpenCode; configure OpenCode por separado en el modelo Ollama. También es adecuado otro cliente compatible con MCP con Ollama.

## Paso 4: Verifica tu conexión

1. Abra la tarjeta de trabajador y haga clic en **Verificar nuevamente**. **Verificación del servidor** ejecuta un ejecutable stdio real con este proyecto, trabajador y base, verifica el protocolo de enlace, la lista de herramientas y que el proyecto y el trabajador estén ubicados.
2. Inicie o reinicie el cliente AI después de agregar el archivo de configuración. Pídale que llame a `get_my_board`.
3. **Cliente conectado** aparecerá en AgentBoard y aparecerá una nueva sesión en la lista. Sólo esto confirma la conexión de su cliente. Si hay herramientas, pero el cliente no está conectado, verifique la ruta al programa, el nombre del archivo de configuración y su sintaxis JSON/TOML.

Inicie su sesión de trabajo con `get_my_board`. La IA no recibe tareas de `Backlog` y `Complete`, incluso si conoce su ID. La IA no puede pretender ser humana y no tiene una orden general de "moverse a cualquier parte". Trabajar con el sistema de archivos fuera de AgentBoard depende de las capacidades y la zona de pruebas del cliente seleccionado.

Cada herramienta declara `title` y las cuatro sugerencias de comportamiento de MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) como valores booleanos explícitos. Las herramientas de lectura son de sólo lectura e idempotentes; Las herramientas de grabación solo agregan registros o avanzan tareas, por lo que ninguna es disruptiva ni va más allá de la base local de AgentBoard.

Instrucciones oficiales para el cliente: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
