# Importar tareas desde JSON

La pestaña **Importar tareas** acepta un archivo JSON en codificación UTF-8. Úselo si está transfiriendo tareas desde otro servicio. Haga clic en **Descargar plantilla JSON**: la plantilla incluye los nombres de propiedades personalizadas actuales y `Slug` del trabajador disponible para este proyecto.

## Paso a paso

1. Cree los trabajadores (**Trabajadores**) y las propiedades (**Propiedades**) necesarios antes de importar.
2. Descargue la plantilla, ábrala en un editor de texto, reemplace los ejemplos con sus propias tareas y guarde el archivo con la extensión `.json`.
3. En la pestaña **Importar tareas**, seleccione un archivo. La interfaz mostrará la cantidad de tareas.
4. Haga clic en **Importar tareas**. El servidor comprobará todo el archivo: si se detecta un error, no se escribirá ni una sola tarea. Corrija el mensaje de error y vuelva a intentarlo.

Archivo mínimo:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Ejemplo con trabajador, propiedad y dependencia:

```json
{
  "version": 1,
  "tasks": [
    {
      "key": "design",
      "title": "Prepare the design",
      "description": "## Goal\nPrepare the home page mockup.",
      "state": "features",
      "priority": "high",
      "testing_mode": "hybrid",
      "assignee": { "type": "worker", "worker": "codex" },
      "ai_test_instructions": "Check the build.",
      "human_test_instructions": "Review the page in a browser.",
      "properties": { "Department": "Design" }
    },
    {
      "title": "Approve the design",
      "depends_on": ["design"],
      "assignee": { "type": "human" }
    }
  ]
}
```

`version` debe ser `1`, matriz `tasks`: de 1 a 1000 tareas. Se requiere `title`. Etapas válidas: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioridades: `critical`, `high`, `medium`, `low`; Modos de prueba: `ai`, `human`, `hybrid`. Responsable: `unassigned`, `human` o `worker` con `worker` existente (Slug o ID). `properties` utiliza los nombres o ID de propiedades ya creadas. `key` es único dentro del archivo; `depends_on` se refiere a dichas claves. El servidor crea los ID de las tareas por sí mismo.

Importar el mismo archivo nuevamente creará nuevos problemas, así que revise el tablero antes de hacer clic nuevamente. Cuando se importan, las tareas se registran como creadas por humanos; Aparece una entrada del sistema en el historial.
