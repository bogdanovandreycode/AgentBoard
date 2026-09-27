# Problemløsning

| Symptom | Hvad skal du tjekke |
| --- | --- |
| `agentboard` ikke fundet | Genstart PowerShell efter Scoop. Når du installerer fra en ZIP, skal du bruge den fulde sti til `agentboard.exe`. |
| Webside, der viser gammel grænseflade efter build | Stop den kørende server Ctrl+C. Udfør `./scripts/build.ps1`, start en ny binær. Vite skal bygge før Go, fordi grænsefladen er indbygget i EXE. Opdater siden Ctrl+F5. |
| Havn 7337 optaget | AgentBoard kører muligvis allerede. Åbn `http://127.0.0.1:7337` eller afslut den gamle proces. Brug `--addr` til anden port. |
| Projekt ikke fundet | Udfør `agentboard init` i den ønskede mappe. Så `agentboard open` fra det eller `agentboard open C:\путь\к\проекту`. |
| Arbejderen kan ikke se opgaven | Opgaven skal tildeles denne særlige arbejder og være placeret i `Features`, `In progress`, `Testing` eller `Verification`. AI kan ikke se `Backlog`, `Complete` og brugerdefinerede kolonner. |
| MCP-servertjek findes, men klienten er ikke forbundet | Serverkontrollen kontrollerer ikke de eksterne klientindstillinger. Genstart klienten, tjek dens konfigurationsfil, sti til `agentboard.exe`, `--project`, `--worker` og generelt `--db`. Bed om at ringe til `get_my_board`. |
| Arbejder offline | Klienten kan have afsluttet eller måske ikke startet MCP endnu. Efter 90 sekunder uden hjerteslag, betragtes sessionen som afbrudt. |
| JSON-fil importeres ikke | Tjek `version: 1`, påkrævet `title`, eksisterende `Slug` arbejdere og ejendomsnavne. JSON tillader ikke kommentarer eller efterfølgende kommaer. |
| Kan ikke genopbygge `agentboard.exe` | Windows kan ikke erstatte en kørende EXE. Stop serveren Ctrl+C og prøv at bygge igen. |
| Opgaver forsvandt efter opdatering | Tjek, at `--db` ikke peger på en anden fil, og at du er logget ind som den samme Windows-bruger. Standardbasen er i `%AppData%\AgentBoard`. |

Hvis fejlen ikke er beskrevet, skal du indsamle den nøjagtige tekst i meddelelsen, `agentboard version`- og Windows-versionerne og trinene for at prøve igen. Udgiv ikke private projektdata eller databaseindhold i en åben udgave.
