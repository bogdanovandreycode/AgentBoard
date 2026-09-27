# Instalación y lanzamiento en Windows

## Instalador (recomendado)

Descarga `agentboard-VERSION-windows-amd64-setup.exe` desde la página [Lanzamientos](https://github.com/bogdanovandreycode/AgentBoard/releases). El asistente de instalación le sugerirá una carpeta; el valor predeterminado es `C:\AI\AgentBoard`. Copiará `agentboard.exe` y la documentación, creará un acceso directo y agregará la carpeta seleccionada al sistema `PATH`. Después de la instalación, abra una nueva terminal para que el comando `agentboard` esté disponible.

En la carpeta de su proyecto ejecute:

```powershell
agentboard init
agentboard open
```

La interfaz está integrada en el `agentboard.exe`; No se necesita una instalación separada de Go o Node.js. Los datos se almacenan en `%AppData%\AgentBoard` y se conservan cuando el programa se actualiza o desinstala. La desinstalación a través de "Aplicaciones instaladas" elimina los accesos directos y la entrada de `PATH`.

## Instalación de pala

Si Scoop aún no está instalado, abra PowerShell como usuario habitual y siga las [instrucciones oficiales de Scoop](https://scoop.sh/). Si existen restricciones en su computadora corporativa, comuníquese con su administrador; AgentBoard también se puede iniciar desde un ZIP sin Scoop.

Como alternativa, instale AgentBoard a través de Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

La disponibilidad del manifiesto en GitHub Release se puede verificar en [página Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Después de agregar el manifiesto a Scoop, el depósito se puede instalar por nombre de depósito y actualizar con el comando `scoop update agentboard`.

## ZIP sin pala

Descargue `agentboard-VERSION-windows-amd64.zip` desde Versiones y descomprímalo, por ejemplo, en `C:\Tools\AgentBoard`. En PowerShell en la carpeta del proyecto:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Para un cliente AI, especifique la ruta completa a `agentboard.exe` en su configuración MCP si el programa no está ubicado en `PATH`.

## Construir desde la fuente

Instale la versión Go de `go.mod` y Node.js 22 o posterior. En PowerShell en la raíz del repositorio:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` ensambla la interfaz en `internal/webui/dist`, ejecuta pruebas de Go y ensambla un `agentboard.exe`. El orden es importante: la interfaz está integrada en el binario cuando se construye Go. Si el antiguo `agentboard.exe open` se está ejecutando, deténgalo antes de reconstruirlo (Ctrl+C); de lo contrario, Windows no le permitirá reemplazar el archivo.

## ¿Dónde están los datos?

- `%AppData%\AgentBoard\agentboard.db` - tareas, proyectos, trabajadores y configuraciones. Puede especificar un archivo diferente con el indicador `--db`, pero `init`, `open`/`serve` y `mcp` deben tener la **misma ruta**.
- `<ваш проект>\.agentboard\project.json`: identificador del proyecto. Este archivo no contiene tareas.
- El servidor solo escucha `127.0.0.1:7337` por defecto. Especifique otra dirección `--addr` antes de la ruta del proyecto: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` abre la página del proyecto y reutiliza un servidor que ya se está ejecutando en esa dirección. Si hay varios proyectos abiertos en un navegador, selecciónelos en la lista de la izquierda.
