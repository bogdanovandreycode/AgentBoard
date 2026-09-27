# Инсталиране и стартиране на Windows

## Инсталатор (препоръчително)

Изтеглете `agentboard-VERSION-windows-amd64-setup.exe` от страницата [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Помощникът за инсталиране ще предложи папка; по подразбиране е `C:\AI\AgentBoard`. Той ще копира `agentboard.exe` и документацията, ще създаде пряк път и ще добави избраната папка към системата `PATH`. След инсталирането отворете нов терминал, така че командата `agentboard` да стане достъпна.

Във вашата папка на проекта стартирайте:

```powershell
agentboard init
agentboard open
```

Интерфейсът е вграден в `agentboard.exe`; Не е необходима отделна инсталация на Go или Node.js. Данните се съхраняват в `%AppData%\AgentBoard` и се запазват, когато програмата се актуализира или деинсталира. Деинсталирането чрез „Инсталирани приложения“ премахва преките пътища и записа от `PATH`.

## Инсталиране на Scoop

Ако Scoop все още не е инсталиран, отворете PowerShell като редовен потребител и следвайте [официалните инструкции на Scoop](https://scoop.sh/). Ако има ограничения за корпоративния ви компютър, свържете се с вашия администратор; AgentBoard може да се стартира и от ZIP без Scoop.

Като алтернатива инсталирайте AgentBoard чрез Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Наличието на манифест в GitHub Release може да се провери на [страница Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). След добавяне на манифест към Scoop, кофата може да бъде инсталирана по име на кофа и актуализирана с командата `scoop update agentboard`.

## ZIP без Scoop

Изтеглете `agentboard-VERSION-windows-amd64.zip` от Releases, разопаковайте, например, в `C:\Tools\AgentBoard`. В PowerShell в папката на проекта:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

За AI клиент, посочете пълния път до `agentboard.exe` в неговата MCP конфигурация, ако програмата не се намира в `PATH`.

## Създаване от източника

Инсталирайте Go версия от `go.mod` и Node.js 22 или по-нова. В PowerShell в основата на хранилището:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` сглобява интерфейса в `internal/webui/dist`, изпълнява Go тестове и сглобява един `agentboard.exe`. Редът е важен: интерфейсът е вграден в двоичния файл, когато Go е изграден. Ако старият `agentboard.exe open` работи, спрете го преди повторно изграждане (Ctrl+C), в противен случай Windows няма да ви позволи да замените файла.

## Къде са данните

- `%AppData%\AgentBoard\agentboard.db` - задачи, проекти, работници и настройки. Можете да посочите различен файл с флага `--db`, но `init`, `open`/`serve` и `mcp` трябва да имат **същия път**.
- `<ваш проект>\.agentboard\project.json` — идентификатор на проекта. Този файл не съдържа задачи.
- Сървърът слуша само `127.0.0.1:7337` по подразбиране. Посочете друг адрес `--addr` преди пътя на проекта: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` отваря страницата на проекта и използва повторно вече работещ сървър на този адрес. Ако няколко проекта са отворени в един браузър, изберете ги от списъка вляво.
