# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Previzualizare AgentBoard](../../assets/social-preview.png)

**[Descărcare pentru Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Site de documentare](https://bogdanovandreycode.github.io/AgentBoard/) · [Licență MIT](../../LICENSE)**

AgentBoard este un comitet local de activități în care lucrătorii umani și AI lucrează la sarcini comune, dar au drepturi diferite. Aplicația este lansată cu un fișier `agentboard.exe`, deschide interfața web în browser și oferă lucrătorilor un server MCP separat prin `stdio`. Datele rămân pe computer.

**[Începe de la zero](START_HERE.md)· [Lucrul cu sarcini](TASKS.md)· [Conexiune AI prin MCP](WORKERS_MCP.md)· [Importați JSON](IMPORT.md)· [Setări](SETTINGS.md)· [Rezolvarea problemelor](TROUBLESHOOTING.md)**

## Caracteristici MVP

- Panou local de proiect cu căutare de sarcini, import JSON, coloane și proprietăți personalizate.
- Acces MCP pentru un anumit lucrător cu tranziții limitate de AI și acceptare finală umană a sarcinilor.
- Istoricul general al sarcinilor, instrucțiunile de verificare, artefactele, costurile AI și diagnosticarea conexiunii lucrătorilor.
- Program de instalare Windows, ZIP și manifest Scoop; interfața web este încorporată în fișierul executabil.

AgentBoard este conceput pentru un utilizator local de încredere. Aplicația nu găzduiește proiecte în cloud și nu lansează clienții AI în sine; dacă este necesar, conectați un client compatibil MCP la lucrător.

## În cinci minute

1. Descărcați programul de instalare `agentboard-VERSION-windows-amd64-setup.exe` de la [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Acesta va sugera un folder (în mod implicit `C:\AI\AgentBoard`) și îl va adăuga la `PATH`. Scoop și ZIP sunt, de asemenea, disponibile.
2. Deschideți PowerShell în folderul proiectului, de exemplu `C:\Projects\MyApp`.
3. Executați `agentboard init` (pentru ZIP: calea completă către `agentboard.exe` și `init`).
4. Executați `agentboard open`. `http://127.0.0.1:7337` se va deschide.
5. Adăugați o sarcină utilizând butonul **Sarcina nouă**. Pentru un lucrător AI, deschideți **Lucrători → Adăugați lucrător**, selectați profilul clientului și copiați configurația MCP.

Dacă nu aveți deja un folder de proiect, creați unul în Windows Explorer. Un proiect poate fi orice folder, chiar și fără Git și cod.

## Instalare prin Scoop

În PowerShell cu [Scoop](https://scoop.sh/)] deja instalat după lansare:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Pentru dezvoltatori există [build from source ](INSTALL.md). Fluxul de lucru de lansare creează un program de instalare, ZIP și manifest Scoop cu SHA-256 din același artefact. Actualizarea versiunii instalate prin Scoop: `scoop update agentboard` după adăugarea manifestului în bucket; detalii - [pregătirea lansării](SCOOP_RELEASE.md).

## Cum este structurat consiliul

`Backlog → Features → In progress → Testing → Verification → Complete`

O persoană poate muta sarcini în jurul tablei. AI poate muta numai `Features → In progress → Testing → Verification`; acceptarea finală în `Complete` este efectuată de un om. Coloanele de utilizator sunt destinate oamenilor: sarcina din ele rămâne în starea `Backlog` pentru MCP. Moduri de testare: AI, uman și hibrid. Istoricul, testele, artefactele și costurile AI sunt atașate sarcinii.

## Echipe

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` înregistrează folderul și scrie acolo doar `.agentboard/project.json`. Datele de producție SQLite se află în directorul de configurare a utilizatorului Windows (`%AppData%\AgentBoard\agentboard.db`), în afara proiectului și în afara instalării Scoop. Eliminarea sau actualizarea pachetului nu ar trebui să elimine aceste date. Înainte de a transfera pe alt computer, faceți o copie a bazei de date în timp ce AgentBoard este oprit.

Muncitori - conturi logice; AgentBoard în sine nu rulează Codex, Claude sau orice alt client AI. Clientul începe un proces MCP local pentru un anumit lucrător. MCP este limita permisiunii aplicației, iar pentru izolarea fișierelor, utilizați sandbox-ul clientului AI.

## Pentru dezvoltatori

Stack: Go, SQLite, oficial MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Creați mai întâi frontend-ul, apoi Go: `./scripts/build.ps1`. Fișierele web sunt incluse în binar prin `go:embed`. Arhitectura și API-ul sunt descrise în [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Licență

AgentBoard este un proiect gratuit și open source sub [licență MIT](../../LICENSE). Utilizarea comercială, modificarea, bifurcarea și redistribuirea sunt permise cu condiția menținerii notificării privind drepturile de autor și a licenței.
