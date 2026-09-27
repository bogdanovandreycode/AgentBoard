# AgentBoard

[🌐 Languages](../LANGUAGES.md)

AgentBoard este un comitet local de activități în care lucrătorii umani și AI lucrează la sarcini comune, dar au drepturi diferite. Aplicația începe cu un singur fișier`agentboard.exe`, deschide interfața web în browser și oferă lucrătorilor un server MCP separat prin`stdio`. Datele rămân pe computer.

**[Începe de la zero](START_HERE.md)· [Lucrul cu sarcini](TASKS.md)· [Conexiune AI prin MCP](WORKERS_MCP.md)· [Importați JSON](IMPORT.md)· [Setări](SETTINGS.md)· [Rezolvarea problemelor](TROUBLESHOOTING.md)**

## În cinci minute

1. Descărcați programul de instalare`agentboard-VERSION-windows-amd64-setup.exe`din [Lansări](https://github.com/bogdanovandreycode/AgentBoard/releases). Acesta va sugera un folder (implicit`C:\AI\AgentBoard`) și îl va adăuga la`PATH`. Scoop și ZIP sunt, de asemenea, disponibile.
2. Deschideți PowerShell în folderul de proiect, cum ar fi`C:\Projects\MyApp`.
3. Executați`agentboard init`(pentru ZIP: calea completă către`agentboard.exe`Şi`init`).
4. Executați`agentboard open`. Se va deschide`http://127.0.0.1:7337`.
5. Adăugați o sarcină utilizând butonul **Sarcina nouă**. Pentru un lucrător AI, deschideți **Lucrători → Adăugați lucrător**, selectați profilul clientului și copiați configurația MCP.

Dacă nu aveți deja un folder de proiect, creați unul în Windows Explorer. Un proiect poate fi orice folder, chiar și fără Git și cod.

## Instalare prin Scoop

În PowerShell cu [Scoop] deja instalat](https://scoop.sh/)după eliberare:

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
