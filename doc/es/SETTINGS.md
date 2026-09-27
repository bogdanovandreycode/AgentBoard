# Configuración

Abra **Configuración** en el menú lateral del proyecto seleccionado. Los cambios se guardan con el botón **Guardar configuración** y se almacenan para el proyecto en la base de datos de AgentBoard.

## Idioma y apariencia

El campo **Idioma** abre una lista desplegable de búsqueda. De forma predeterminada, está seleccionada la opción **Seguir sistema**, que toma el idioma del navegador. Idiomas disponibles en las capturas de pantalla: árabe, portugués (Brasil), chino simplificado, checo, danés, holandés, inglés, finlandés, francés, alemán, italiano, japonés, coreano, noruego bokmål, polaco, ruso, español, sueco, turco, ucraniano, vietnamita; además bielorruso, rumano y búlgaro. **El sistema de seguimiento** toma el idioma del navegador. Las firmas principales se traducen manualmente, las líneas restantes de la interfaz de usuario tienen una traducción automática preliminar. Antes de su publicación pública, es recomendable revisar las traducciones realizadas por hablantes nativos; Los nombres de clientes, comandos, campos JSON y datos de usuario permanecen sin traducción.

**Esquema de colores**: Oscuro (tema original), Claro, Negro, Ubuntu y Windows. **Zona horaria** controla la visualización de fechas; los datos continúan almacenándose en UTC. **Zona horaria del sistema** utiliza la configuración de la computadora.

## Columnas

Las cuatro etapas `Features`, `In progress`, `Testing`, `Verification` están fijadas en este orden: no se pueden cambiar de nombre ni eliminar. Otras columnas se pueden reorganizar usando las flechas disponibles. Ingrese un nombre y haga clic en **Agregar columna** para crear una columna personalizada. Está diseñado para tareas humanas: la IA ve una tarea como `Backlog` y no la recibe a través de MCP. Eliminar una columna mientras se guarda transfiere sus tareas al `Backlog` normal.

## Web y MCP

**Actualización del tablero** establece la frecuencia de actualización del tablero y las tarjetas en segundos (1–60). **Actualización de trabajadores** actualiza el estado de los trabajadores (2 a 120 segundos). Se trata de un sondeo de interfaz web, no de una frecuencia de activación de IA. Los trabajadores no comienzan automáticamente.

La dirección del servidor web se establece al iniciar la CLI, por ejemplo `agentboard open --addr 127.0.0.1:7444`. Para cambiar la dirección, se debe reiniciar el servidor. El valor predeterminado es `127.0.0.1:7337`. MCP opera a través de un comando local independiente `agentboard mcp --project ... --worker ...` y es independiente del puerto web. Para una base no estándar, especifique el mismo `--db` en todos los comandos. Copie la configuración de cada cliente de la tarjeta de trabajador.
