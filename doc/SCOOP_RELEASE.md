# Подготовка релиза для Scoop

## Что уже автоматизировано

`scripts/package-scoop.ps1 -Version 0.1.0` выполняет `npm ci`, сборку frontend, `go test ./...`, сборку Windows `agentboard.exe` с номером версии, создаёт ZIP и рассчитывает SHA-256 этого ZIP. Из **этого же** файла он создаёт `release/agentboard.json` с URL, хешем, CLI-shim, ярлыком, `checkver` и `autoupdate`. База живёт вне каталога установки, поэтому `persist` в manifest не нужен.

`.github/workflows/release.yml` на теге `vX.Y.Z` запускает ту же упаковку в Windows runner и прикладывает ZIP и manifest к GitHub Release. Ручной запуск workflow создаёт только artifact для проверки, без публикации релиза.

## Порядок публикации для сопровождающего

1. Проверьте, что код, документация и номер версии готовы. Определите лицензию проекта: сейчас manifest указывает `Unknown`, потому что файл лицензии в репозитории не задан. Если вы выбираете лицензию, добавьте `LICENSE` и обновите `license` в скрипте перед релизом.
2. Локально выполните `./scripts/package-scoop.ps1 -Version X.Y.Z`. Проверьте `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` и вывод `agentboard version` после распаковки. Не меняйте ZIP после расчёта хеша.
3. Создайте и отправьте тег `vX.Y.Z`. GitHub Actions опубликует Release с тем же именем ZIP и manifest. Проверьте, что оба файла доступны на странице Release и что скачанный ZIP имеет SHA-256 из manifest.
4. На чистой Windows-машине с Scoop выполните `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, затем `agentboard version`, `agentboard init` и `agentboard open` в тестовой папке.
5. Для постоянного канала обновлений поместите сгенерированный `agentboard.json` в собственный Scoop bucket или предложите его в подходящий публичный bucket. Проверяйте `scoop update agentboard` после следующего релиза. Команда установки из URL подходит для первого знакомства, а bucket удобнее для обновлений.

Сейчас на GitHub ещё нет релиза. Поэтому установка из URL пока недоступна; сам manifest до публикации не добавлен в репозиторий с вымышленным хешем. Не подставляйте вручную случайное значение `hash`: Scoop сверяет содержимое скачанного ZIP.

Для проверок manifest ориентируйтесь на [формат Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [создание manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) и [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
