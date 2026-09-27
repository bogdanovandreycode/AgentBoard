# Preparando un lanzamiento para Scoop

## Lo que ya está automatizado

`scripts/package-scoop.ps1 -Version X.Y.Z` ejecuta `npm ci`, compilación frontend, `go test ./...`, compilación de Windows `agentboard.exe` con número de versión, crea un ZIP y calcula el SHA-256 de ese ZIP. Desde **el mismo** archivo, crea `release/agentboard.json` con URL, hash, licencia MIT, CLI-shim, acceso directo, `checkver` y `autoupdate`. El ZIP y el instalador contienen el archivo `LICENSE`. La base de datos se encuentra fuera del directorio de instalación, por lo que `persist` no es necesario en el manifiesto.

`.github/workflows/release.yml` en la etiqueta `vX.Y.Z` ejecuta el mismo paquete en Windows runner, crea el instalador de Inno Setup y adjunta el ZIP, el instalador y el manifiesto a la versión de GitHub. La ejecución manual de un flujo de trabajo solo crea un artefacto para realizar pruebas, sin publicar una versión.

## Orden de publicación para mantenedor

1. Verifique que el código, la documentación, el número de versión y el archivo `LICENSE` estén listos.
2. Ejecute `./scripts/package-scoop.ps1 -Version X.Y.Z` localmente. Verifique la salida de `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` y `agentboard version` después de desembalar. No cambie el ZIP después de que se haya calculado el hash.
3. Cree y envíe la etiqueta `vX.Y.Z`. GitHub Actions publicará una versión con ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` y manifiesto. Verifique los tres archivos en la página de lanzamiento y el ZIP SHA-256 del manifiesto.
4. En una máquina Windows limpia con Scoop, ejecute `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, luego `agentboard version`, `agentboard init` y `agentboard open` en la carpeta de prueba.
5. Para obtener una fuente de actualización permanente, coloque el `agentboard.json` generado en su propio depósito Scoop u ofrézcalo en un depósito público adecuado. Vuelva a consultar `scoop update agentboard` después del próximo lanzamiento. El comando de instalación desde una URL es adecuado para el primer contacto, mientras que el depósito es más conveniente para las actualizaciones.

No sustituya manualmente el valor aleatorio `hash`: Scoop verifica el contenido del ZIP descargado.

Para verificaciones de manifiesto, céntrese en [formato Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [creación de manifiesto](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) y [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
