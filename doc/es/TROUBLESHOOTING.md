# Resolución de problemas

| Síntoma | Qué comprobar |
| --- | --- |
| `agentboard` no encontrado | Reinicie PowerShell después de Scoop. Al instalar desde un ZIP, utilice la ruta completa a `agentboard.exe`. |
| Página web que muestra la interfaz antigua después de la compilación | Detenga el servidor en ejecución Ctrl+C. Ejecute `./scripts/build.ps1`, inicie un nuevo binario. Vite debe compilarse antes que Go porque la interfaz está integrada en el EXE. Actualizar la página Ctrl+F5. |
| Puerto 7337 ocupado | Es posible que AgentBoard ya se esté ejecutando. Abra `http://127.0.0.1:7337` o finalice el proceso anterior. Para otros puertos utilice `--addr`. |
| Proyecto no encontrado | En la carpeta deseada, ejecute `agentboard init`. Luego `agentboard open` o `agentboard open C:\путь\к\проекту`. |
| El trabajador no ve la tarea | La tarea debe estar asignada a este trabajador en particular y estar ubicada en `Features`, `In progress`, `Testing` o `Verification`. AI no ve `Backlog`, `Complete` ni columnas personalizadas. |
| Existe la verificación del servidor MCP, pero el cliente no está conectado | La verificación del servidor no verifica la configuración del cliente externo. Reinicie el cliente, verifique su archivo de configuración, ruta a `agentboard.exe`, `--project`, `--worker` y `--db` general. Pida llamar a `get_my_board`. |
| Trabajador desconectado | Es posible que el cliente haya finalizado o que aún no haya iniciado MCP. Después de 90 segundos sin latidos, la sesión se considera desconectada. |
| El archivo JSON no se importa | Verifique `version: 1`, `title` requerido, trabajadores `Slug` existentes y nombres de propiedades. JSON no permite comentarios ni comas finales. |
| No se puede reconstruir `agentboard.exe` | Windows no puede reemplazar un EXE en ejecución. Detenga el servidor Ctrl+C e intente la compilación nuevamente. |
| Las tareas desaparecieron después de la actualización | Compruebe que `--db` no apunte a otro archivo y que haya iniciado sesión como el mismo usuario de Windows. La base predeterminada está en `%AppData%\AgentBoard`. |

Si el error no se describe, recopile el texto exacto del mensaje, las versiones `agentboard version` y Windows, y los pasos para volver a intentarlo. No publique datos de proyectos privados o contenidos de bases de datos en una edición abierta.
