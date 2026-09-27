# Instalare și lansare pe Windows

## Instalator (recomandat)

Descărcați `agentboard-VERSION-windows-amd64-setup.exe` de pe pagina [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Expertul de instalare va sugera un folder; implicit este `C:\AI\AgentBoard`. Acesta va copia `agentboard.exe` și documentația, va crea o comandă rapidă și va adăuga folderul selectat în sistem `PATH`. După instalare, deschideți un nou terminal, astfel încât comanda `agentboard` să devină disponibilă.

În folderul de proiect rulați:

```powershell
agentboard init
agentboard open
```

Interfața este încorporată în `agentboard.exe`; Nu este necesară nicio instalare separată a Go sau Node.js. Datele sunt stocate în `%AppData%\AgentBoard` și sunt păstrate atunci când programul este actualizat sau dezinstalat. Dezinstalarea prin „Aplicații instalate” elimină comenzile rapide și intrarea din `PATH`.

## Instalarea Scoop

Dacă Scoop nu este deja instalat, deschideți PowerShell ca utilizator obișnuit și urmați [instrucțiunile oficiale Scoop](https://scoop.sh/). Dacă există restricții pe computerul dvs. corporativ, contactați administratorul; AgentBoard poate fi lansat și dintr-un ZIP fără Scoop.

Alternativ, instalați AgentBoard prin Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Disponibilitatea manifestului în GitHub Release poate fi verificată pe [pagina Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). După adăugarea manifestului în Scoop, compartimentul poate fi instalat după numele compartimentului și actualizat cu comanda `scoop update agentboard`.

## ZIP fără Scoop

Descărcați `agentboard-VERSION-windows-amd64.zip` din versiuni, despachetați, de exemplu, în `C:\Tools\AgentBoard`. În PowerShell în folderul proiectului:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Pentru un client AI, specificați calea completă către `agentboard.exe` în configurația sa MCP dacă programul nu este localizat în `PATH`.

## Construiți din sursă

Instalați versiunea Go de la `go.mod` și Node.js 22 sau o versiune ulterioară. În PowerShell, la rădăcina depozitului:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` asamblează interfața în `internal/webui/dist`, execută teste Go și asamblează un `agentboard.exe`. Ordinea este importantă: interfața este încorporată în binar atunci când este construit Go. Dacă vechiul `agentboard.exe open` rulează, opriți-l înainte de a reconstrui (Ctrl+C), altfel Windows nu vă va permite să înlocuiți fișierul.

## Unde sunt datele

- `%AppData%\AgentBoard\agentboard.db` - sarcini, proiecte, lucrători și setări. Puteți specifica un fișier diferit cu indicatorul `--db`, dar `init`, `open`/`serve` și `mcp` trebuie să aibă **aceeași cale**.
- `<ваш проект>\.agentboard\project.json` — identificatorul proiectului. Acest fișier nu conține sarcini.
- Serverul ascultă numai `127.0.0.1:7337` în mod implicit. Specificați o altă adresă `--addr` înainte de calea proiectului: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` deschide pagina proiectului și reutiliza un server care rulează deja la acea adresă. Dacă mai multe proiecte sunt deschise într-un browser, selectați-le din lista din stânga.
