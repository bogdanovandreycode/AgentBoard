# Comienza aquí

## ¿Qué es AgentBoard?

Imagine un tablero normal con tarjetas de tareas. Creas tareas, asignas a alguien responsable y supervisas el trabajo. Los trabajadores de IA reciben solo las tareas asignadas a través de MCP e informan el progreso. Tú decides cuándo la tarea finalmente estará lista.

**Proyecto**: una carpeta en la computadora y un tablero separado. **Tarea**: una tarjeta con una descripción, persona responsable y etapa. **Trabajador** es el nombre lógico del cliente de IA. **MCP** es la forma en que el cliente de IA se conecta al AgentBoard. **Scoop** es un administrador de instalación de software para Windows.

No se requiere código. Necesita Windows, un navegador, PowerShell y, para que la IA funcione, un cliente de IA instalado con soporte para servidores MCP locales.

## Primer lanzamiento

1. Instale la aplicación según las [instrucciones](INSTALL.md).
2. Cree una carpeta de proyecto en Explorer, por ejemplo `C:\Projects\MyFirstProject`.
3. Abra esta carpeta en el Explorador. Haga clic en la barra de direcciones, escriba `powershell` y presione Entrar.
4. En la ventana que se abre, haga:

```powershell
agentboard init
agentboard open
```

5. Se abrirá un navegador con la dirección `http://127.0.0.1:7337`. Deje abierta la ventana de PowerShell mientras usa la pizarra. Cerrar la ventana detendrá el servidor local, pero las tareas permanecerán.

Si no se encuentra el comando `agentboard`, cierre PowerShell y vuelva a abrirlo después de instalar Scoop. Al instalar desde un ZIP, utilice la ruta completa a `agentboard.exe`.

## Primera tarea

Haga clic en **Nueva tarea**, complete **Título**, si es necesario **Descripción** y luego **Guardar**. Una nueva tarea en `Backlog` está disponible para los humanos. Para que la IA comience a funcionar, asigna un trabajador y transfiere la tarea a `Features`. Descripción paso a paso de los campos - [TASKS.md](TASKS.md).

## Primer trabajador

Abra **Trabajadores → Agregar trabajador**. Seleccione el cliente de IA que está utilizando (por ejemplo Codex o Claude Code), verifique el nombre y el ID corto `Slug`, haga clic en **Guardar**. Abra el trabajador creado: hay una configuración de MCP, un botón de copiar y una verificación del servidor. Copie la configuración al cliente AI según [WORKERS_MCP.md](WORKERS_MCP.md). Una vez conectado, solicite al cliente que llame a `get_my_board`.

## Qué leer a continuación

- [Tareas, pruebas, historias y columnas](TASKS.md)
- [Conectando trabajadores y comprobando MCP](WORKERS_MCP.md)
- [Importar tareas desde JSON](IMPORT.md)
- [Idioma, tema, zona horaria y actualización](SETTINGS.md)
- [Problemas típicos](TROUBLESHOOTING.md)
