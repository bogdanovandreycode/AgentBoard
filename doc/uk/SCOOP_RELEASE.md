# Підготовка релізу для Scoop

## Що вже автоматизовано

`scripts/package-scoop.ps1 -Version 0.2.0` виконує `npm ci`, складання frontend, `go test ./...`, складання Windows `agentboard.exe` з номером версії, створює ZIP і розраховує SHA-256 цього ZIP. З цього ж файлу він створює `release/agentboard.json` з URL, хешем, CLI-shim, ярликом, `checkver` і `autoupdate`. База живе поза каталогом установки, тому `persist` в manifest не потрібний.

`.github/workflows/release.yml` на тезі `vX.Y.Z` запускає ту ж упаковку в Windows runner, збирає інсталятор Inno Setup і прикладає ZIP, інсталятор і manifest до GitHub Release. Ручний запуск workflow створює тільки artifact для перевірки без публікації релізу.

## Порядок публікації для супроводжуючого

1. Перевірте, чи код, документація та номер версії готові. Визначте ліцензію проекту: зараз manifest вказує на `Unknown`, тому що файл ліцензії в репозиторії не заданий. Якщо ви вибираєте ліцензію, додайте `LICENSE` та оновіть `license` у скрипті перед релізом.
2. Виконайте локально `./scripts/package-scoop.ps1 -Version X.Y.Z`. Перевірте `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` та виведення `agentboard version` після розпакування. Не змінюйте ZIP після розрахунку хешу.
3. Створіть та надішліть тег `vX.Y.Z`. GitHub Actions опублікує Release з ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` та manifest. Перевірте всі три файли на сторінці Release та SHA-256 ZIP з manifest.
4. На чистій машині Windows з Scoop виконайте `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, потім `agentboard version`, `agentboard init` і `agentboard open` у тестовій папці.
5. Для постійного каналу оновлень помістіть згенерований `agentboard.json` у власний Scoop bucket або запропонуйте його у відповідний публічний bucket. Перевіряйте `scoop update agentboard` після наступного релізу. Команда установки з URL підходить для першого знайомства, а bucket зручніший для оновлень.

Не підставляйте вручну випадкове значення `hash`: Scoop звіряє вміст завантаженого ZIP.

Для перевірок manifest орієнтуйтеся на [формат Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [створення manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) та [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate)].
