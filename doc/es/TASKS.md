# Trabajar con tareas

## Agregar tarea

En la pestaña **Tablero**, haz clic en **Nueva tarea**. Complete el título y la descripción. La descripción y las instrucciones de prueba admiten Markdown: resalte el texto y use la barra de formato para negrita, cursiva, enlace y lista. Haga clic en **Guardar**.

Campos:

| Campo | ¿Qué significa |
| --- | --- |
| Título | Nombre corto de la tarea. |
| Descripción | Qué hay que hacer y cómo entender que el trabajo está listo. |
| Estado | Etapa de trabajo; una nueva tarea normalmente comienza en `Backlog`. |
| Prioridad | `critical`, `high`, `medium` o `low`. |
| Responsable | Una persona, un trabajador concreto o sin cita previa. |
| Modo de prueba | `AI` - comprueba la IA; `Human` - controles humanos; `Hybrid` - ambos. |
| Dependencias | Tareas que deberían completarse antes. |
| Instrucciones de prueba de IA/humanos | Instrucciones al revisor correspondiente. |
| Propiedades personalizadas | Campos adicionales creados en la pestaña **Propiedades**. |

Para una tarea de IA, primero cree un trabajador, selecciónelo en Responsable y luego transfiera la tarjeta a `Features`. La IA solo ve las tareas asignadas en cuatro etapas, desde `Features` hasta `Verification`.

## Mover tarea

Arrastra la tarjeta entre las columnas. Puede cambiar entre una vista de todas las columnas y columnas anchas con desplazamiento horizontal. `Backlog` y `Complete` están controlados por humanos. La IA solo puede avanzar en la tarea `Features → In progress → Testing → Verification` a través de herramientas MCP especiales. Una tarea con pruebas manuales incompletas no debería pasar la verificación de IA.

## Tarjeta de tarea

Haga clic en una tarjeta para ver la descripción, el propietario, las instrucciones de prueba y las pestañas de historial, pruebas, artefactos y costos de IA. **Editar tarea** cambia el contenido. El comentario de la persona se agrega a la historia general. Buscar en los filtros superiores busca por título, descripción, ID y trabajador; un botón separado abre una búsqueda grande.

## Columnas y propiedades

En la pestaña **Configuración** puedes cambiar el orden de las columnas disponibles para una persona y agregar las tuyas propias. Las cuatro etapas de la IA son fijas y avanzan en el mismo orden. La columna de usuario es un lugar para las tareas pospuestas por una persona: para la IA, dicha tarea tiene el estado `Backlog`. Cuando se elimina una columna, sus tareas vuelven a la normalidad `Backlog`.

En la pestaña **Propiedades** puedes agregar campos como texto, número, bandera, fecha, selección y URL. `Human only` oculta el campo a la IA; `Agent read` permite leer, `Agent read/write` también permite escribir a través de herramientas compatibles. Esto no cambia los derechos de la IA sobre los pasos de la tarea.
