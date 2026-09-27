# Installation och start på Windows

## Installationsprogram (rekommenderas)

Ladda ner `agentboard-VERSION-windows-amd64-setup.exe` från sidan [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Installationsguiden kommer att föreslå en mapp; standard är `C:\AI\AgentBoard`. Den kopierar `agentboard.exe` och dokumentation, skapar en genväg och lägger till den valda mappen till systemet `PATH`. Efter installationen, öppna en ny terminal så att kommandot `agentboard` blir tillgängligt.

Kör i din projektmapp:

```powershell
agentboard init
agentboard open
```

Gränssnittet är inbyggt i `agentboard.exe`; Ingen separat installation av Go eller Node.js behövs. Data lagras i `%AppData%\AgentBoard` och bevaras när programmet uppdateras eller avinstalleras. Avinstallation via "Installerade applikationer" tar bort genvägar och posten från `PATH`.

## Installerar Scoop

Om Scoop inte redan är installerat, öppna PowerShell som din vanliga användare och följ [officiella Scoop](https://scoop.sh/) instruktioner. Om det finns begränsningar på din företagsdator, kontakta din administratör; AgentBoard kan också startas från en ZIP utan Scoop.

Alternativt installera AgentBoard via Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Tillgängligheten för manifest i GitHub Release kan kontrolleras på [sida Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Efter att ha lagt till manifest i Scoop, kan hinken installeras med hinknamn och uppdateras med kommandot `scoop update agentboard`.

## ZIP utan Scoop

Ladda ner `agentboard-VERSION-windows-amd64.zip` från Releases, packa upp till exempel i `C:\Tools\AgentBoard`. I PowerShell i projektmappen:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

För en AI-klient, ange den fullständiga sökvägen till `agentboard.exe` i dess MCP-konfiguration om programmet inte finns i `PATH`.

## Bygg från källan

Installera Go-versionen från `go.mod` och Node.js 22 eller senare. I PowerShell i roten av förvaret:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` monterar frontend till `internal/webui/dist`, kör Go-tester och monterar en `agentboard.exe`. Ordningen är viktig: gränssnittet är inbyggt i binären när Go byggs. Om den gamla `agentboard.exe open` körs, stoppa den innan du bygger om (Ctrl+C), annars tillåter inte Windows dig att ersätta filen.

## Var finns data

- `%AppData%\AgentBoard\agentboard.db` - uppgifter, projekt, arbetare och inställningar. Du kan ange en annan fil med `--db`-flaggan, men `init`, `open`/`serve` och `mcp` måste ha **samma sökväg**.
- `<ваш проект>\.agentboard\project.json` — projektidentifierare. Den här filen innehåller inga uppgifter.
- Servern lyssnar bara på `127.0.0.1:7337` som standard. Ange en annan adress `--addr` före projektsökvägen: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` öppnar projektsidan och återanvänder en server som redan körs på den adressen. Om flera projekt är öppna i en webbläsare, välj dem i listan till vänster.
