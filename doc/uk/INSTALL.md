# Встановлення та запуск на Windows

## Установник (рекомендується)

Завантажте `agentboard-VERSION-windows-amd64-setup.exe` зі сторінки [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Майстер установки запропонує папку; за замовчуванням це `C:\AI\AgentBoard`. Він скопіює `agentboard.exe` та документацію, створить ярлик та додасть обрану папку до системного `PATH`. Після встановлення відкрийте новий термінал, щоб команда `agentboard` стала доступною.

У папці вашого проекту виконайте:

```powershell
agentboard init
agentboard open
```

Інтерфейс вбудований у `agentboard.exe`; окрема установка Go чи Node.js не потрібна. Дані зберігаються в `%AppData%\AgentBoard` і зберігаються під час оновлення або видалення програми. Видалення через «Встановлені програми» прибирає ярлики та запис із `PATH`.

## Установка Scoop

Якщо Scoop ще не встановлено, відкрийте PowerShell від свого звичайного користувача та дотримуйтесь [офіційної інструкції Scoop](https://scoop.sh/). У разі обмеження корпоративного комп'ютера зверніться до адміністратора; AgentBoard також можна запустити із ZIP без Scoop.

Альтернативно встановіть AgentBoard через Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Наявність manifest у GitHub Release можна перевірити на сторінці Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Після додавання manifest в Scoop bucket можна встановлювати на ім'я bucket і оновлювати командою `scoop update agentboard`.

## ZIP без Scoop

Завантажте `agentboard-VERSION-windows-amd64.zip` з Releases, розпакуйте, наприклад, `C:\Tools\AgentBoard`. PowerShell в папці проекту:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Для AI-клієнта вкажіть у його MCP конфігурації повний шлях до `agentboard.exe`, якщо програма не знаходиться в `PATH`.

## Складання з вихідних

Встановіть Go версії з `go.mod` та Node.js 22 або новіші. У PowerShell докорінно репозиторія:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` збирає frontend у `internal/webui/dist`, запускає Go-тести та збирає один `agentboard.exe`. Порядок важливий: інтерфейс вбудовується в бінарник при збиранні Go. Якщо старий `agentboard.exe open` запущено, зупиніть його перед перескладанням (Ctrl+C), інакше Windows не дасть замінити файл.

## Де дані

- `%AppData%\AgentBoard\agentboard.db` - завдання, проекти, воркери та налаштування. Можна задати інший файл прапором `--db`, але у `init`, `open`/`serve` та `mcp` має бути **один і той самий шлях**.
- `<ваш проект>\.agentboard\project.json` – ідентифікатор проекту. Цей файл не містить завдань.
- Сервер слухає лише `127.0.0.1:7337` за промовчанням. Іншу адресу вкажіть `--addr` до шляху проекту: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` відкриває сторінку проекту і повторно використовує сервер, що вже працює, на цій адресі. Якщо в одному браузері відкрито декілька проектів, вибирайте їх за допомогою списку зліва.
