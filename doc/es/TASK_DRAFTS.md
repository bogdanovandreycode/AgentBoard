# Gestión de proyectos y borradores de tareas.

La página **Todos los proyectos** muestra la carpeta registrada, el número de proyectos publicados
tareas y actividad más reciente. La actividad tiene en cuenta los cambios en las tareas, el historial y
registros de gastos. El botón del lápiz le permite cambiar el nombre del proyecto o seleccionar
carpeta existente en el cuadro de diálogo del sistema. Cambiar la ruta actualiza el registro y el marcador;
Los ID de proyecto, trabajadores, tareas y configuraciones se guardan, los archivos no se mueven.
No se pueden seleccionar carpetas que ya estén registradas en otro proyecto.

## Configuración global

Abra **Configuración global** desde la página o pestaña de selección de proyectos
**Configuración global** en la configuración de cualquier tablero. Estos parámetros son comunes a todos los proyectos.
una base AgentBoard:

- El guardado automático de borradores mediante temporizador está habilitado de forma predeterminada, cada **10 segundos**.
- **Guardar al salir del campo** está habilitado de forma predeterminada y se puede configurar por separado.
- Las comprobaciones de **contexto del proyecto MCP** están habilitadas y deshabilitadas globalmente. su area
  Las acciones y un indicador CLI independiente se describen en la sección [Trabajadores y MCP](WORKERS_MCP.md).

## Borrador de ciclo de vida

El cuadro de diálogo Nueva tarea guarda el trabajo sin terminar como borrador. Línea de estado
distingue entre cambios no guardados, guardado en progreso, borrador guardado y error
conservación. El borrador no aparece en el tablero y no está disponible para los trabajadores de MCP hasta
publicado.

Antes de publicar, use el botón **Crear tarea** para completar **Título, Estado,
Responsable y Descripción**. Antes de esto, **Guardar borrador** guarda un borrador sin crearlo.
tareas. Un formulario vacío se descarta de forma predeterminada cuando se cierra. Si se borra
borrador existente a un estado vacío, se eliminará la próxima vez que se guarde.
Al cerrar un cuadro de diálogo que no está vacío, también se guarda.

El botón **Borradores** en el tablero le permite continuar o eliminar borradores guardados.
Se admiten varios borradores. Después de reiniciar, se abre
borrador. Los borradores se almacenan en SQLite; Con cada cambio, el local se actualiza.
copiar en el navegador, para que las ediciones no guardadas sobrevivan al reinicio o
fallo temporal del servidor. Si el almacenamiento del navegador no está disponible, un cuadro de diálogo advierte
que el borrador debe guardarse en el servidor antes de cerrarlo.

La redacción de borradores utiliza revisiones para evitar que una ventana sobrescriba el trabajo
otro. La respuesta de guardado perdido es consistente si ya hay una en el servidor
datos idénticos. Si el contenido es diferente, se muestra un conflicto y
**Recuperar como un nuevo borrador** guarda el contenido local por separado. Publicación
atómico e idempotente: volver a publicar el mismo borrador después de que se pierda
La respuesta devuelve la misma tarea en lugar de crear una copia.
