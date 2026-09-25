# Установка и запуск на Windows

## Установка Scoop

Если Scoop ещё не установлен, откройте PowerShell от своего обычного пользователя и следуйте [официальной инструкции Scoop](https://scoop.sh/). При ограничениях корпоративного компьютера обратитесь к администратору; AgentBoard также можно запустить из ZIP без Scoop.

После публикации первого релиза установите AgentBoard:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Наличие manifest в GitHub Release можно проверить на [странице Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Пока релиза нет, скачивать нечего. После добавления manifest в Scoop bucket можно устанавливать по имени bucket и обновлять командой `scoop update agentboard`.

## ZIP без Scoop

Скачайте `agentboard-VERSION-windows-amd64.zip` из Releases, распакуйте, например, в `C:\Tools\AgentBoard`. В PowerShell в папке проекта:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Для AI-клиента укажите в его MCP-конфигурации полный путь к `agentboard.exe`, если программа не находится в `PATH`.

## Сборка из исходников

Установите Go версии из `go.mod` и Node.js 22 или новее. В PowerShell в корне репозитория:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` собирает frontend в `internal/webui/dist`, запускает Go-тесты и собирает один `agentboard.exe`. Порядок важен: интерфейс встраивается в бинарник при сборке Go. Если старый `agentboard.exe open` запущен, остановите его перед пересборкой (Ctrl+C), иначе Windows не даст заменить файл.

## Где данные

- `%AppData%\AgentBoard\agentboard.db` — задачи, проекты, воркеры и настройки. Можно задать другой файл флагом `--db`, но у `init`, `open`/`serve` и `mcp` должен быть **один и тот же путь**.
- `<ваш проект>\.agentboard\project.json` — идентификатор проекта. Этот файл не содержит задач.
- Сервер слушает только `127.0.0.1:7337` по умолчанию. Другой адрес укажите `--addr` до пути проекта: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` открывает страницу проекта и повторно использует уже работающий сервер на этом адресе. Если в одном браузере открыто несколько проектов, выбирайте их через список слева.
