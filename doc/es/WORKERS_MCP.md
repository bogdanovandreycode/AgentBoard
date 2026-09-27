# Trabajadores y conexión MCP

## Paso 1. Crear un trabajador

En AgentBoard, abra **Trabajadores → Agregar trabajador**. Seleccione un perfil de cliente: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama vía OpenCode, VS Code Copilot u otro cliente MCP. El perfil completa las características típicas (código, pruebas, Git); puedes cambiarlos. El nombre es visible en el tablero y `Slug` es un identificador corto sin espacios para el comando MCP. Haga clic en **Guardar**.

Un trabajador corresponde a una personalidad de IA. Cree diferentes trabajadores para diferentes clientes o equipos. El perfil de capacidad describe la especialización, pero no extiende los derechos de la IA a las etapas de la tarea.

## Paso 2. Copiar la configuración

Abra el trabajador creado. El bloque **MCP diagnostics** muestra el fragmento de configuración y el archivo donde agregarlo. Haga clic en **Copiar configuración de MCP**. Si el archivo ya existe, agregue el servidor sugerido al objeto `mcpServers`/`servers`/`mcp` existente sin borrar los otros servidores.

El comando principal se ve así:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` debería apuntar a la carpeta que registró mediante `agentboard init`. `--worker` — `Slug` del trabajador creado. El cliente MCP ejecuta este comando por sí mismo cuando necesita herramientas. En el navegador AgentBoard, el servidor web puede ejecutarse por separado.

Después de la verificación, el bloque de diagnóstico muestra la ruta absoluta al `agentboard.exe` en ejecución. Es especialmente útil al instalar desde un ZIP. Al instalar a través de Scoop, puede usar el comando `agentboard` si el cliente ve el mismo `PATH`.

## Paso 3: Agregue un servidor a su cliente

La interfaz del trabajador ya tiene un fragmento listo para usar. A continuación se muestra una explicación de dónde se utiliza:

| Cliente | Dónde insertar | Cómo comprobar del lado del cliente |
| --- | --- | --- |
| Códice | `%USERPROFILE%\.codex\config.toml`, sección `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Código Claude | `.mcp.json` en la carpeta del proyecto | `claude mcp list` |
| Géminis CLI | `%USERPROFILE%\.gemini\settings.json`, objeto `mcpServers` | `/mcp list` en Géminis CLI |
| Cursores | Proyecto `.cursor\mcp.json` | lista de servidores MCP en la configuración del cursor |
| Código abierto | Proyecto `opencode.json`, objeto `mcp` | lista de herramientas MCP en OpenCode |
| Copiloto de código VS | Proyecto `.vscode\mcp.json`, objeto `servers` | comando **MCP: Listar servidores** |

Para Codex y Claude Code, un ejemplo con el proyecto `C:\Projects\MyFirstProject` y el trabajador `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

En JSON, la barra invertida de Windows se duplica; un fragmento listo para usar de la interfaz hace esto automáticamente. Si está utilizando `--db` con una base personalizada, agréguelo a la configuración de MCP `args` y especifique la misma ruta que cuando inició `open`.

**Ollama** proporciona un modelo local, pero no reemplaza al cliente MCP. El perfil **Ollama vía OpenCode** genera una configuración MCP para OpenCode; configure OpenCode por separado en el modelo Ollama. También es adecuado otro cliente compatible con MCP con Ollama.

## Paso 4: Verifica tu conexión

1. Abra la tarjeta de trabajador y haga clic en **Verificar nuevamente**. **Verificación del servidor** debe mostrar la cantidad de herramientas MCP. Se trata de una comprobación y detección del protocolo interno de las herramientas del servidor.
2. Inicie o reinicie el cliente AI después de agregar el archivo de configuración. Pídale que llame a `get_my_board`.
3. **Cliente conectado** aparecerá en AgentBoard y aparecerá una nueva sesión en la lista. Sólo esto confirma la conexión de su cliente. Si hay herramientas, pero el cliente no está conectado, verifique la ruta al programa, el nombre del archivo de configuración y su sintaxis JSON/TOML.

Inicie su sesión de trabajo con `get_my_board`. La IA no recibe tareas de `Backlog` y `Complete`, incluso si conoce su ID. La IA no puede pretender ser humana y no tiene una orden general de "moverse a cualquier parte". Trabajar con el sistema de archivos fuera de AgentBoard depende de las capacidades y la zona de pruebas del cliente seleccionado.

Instrucciones oficiales para el cliente: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
