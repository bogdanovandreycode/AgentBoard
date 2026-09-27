# Problemløsning

| Symptom | Hva du bør sjekke |
| --- | --- |
| `agentboard` ikke funnet | Start PowerShell på nytt etter Scoop. Når du installerer fra en ZIP, bruk hele banen til `agentboard.exe`. |
| Nettside som viser gammelt grensesnitt etter bygg | Stopp den kjørende serveren Ctrl+C. Kjør `./scripts/build.ps1`, start en ny binær. Vite må bygges før Go fordi grensesnittet er innebygd i EXE. Oppdater siden Ctrl+F5. |
| Port 7337 opptatt | AgentBoard kjører kanskje allerede. Åpne `http://127.0.0.1:7337` eller avslutt den gamle prosessen. For annen port bruk `--addr`. |
| Prosjekt ikke funnet | I ønsket mappe, kjør `agentboard init`. Deretter `agentboard open` fra den eller `agentboard open C:\путь\к\проекту`. |
| Arbeideren ser ikke oppgaven | Oppgaven må tildeles denne bestemte arbeideren og være plassert i `Features`, `In progress`, `Testing` eller `Verification`. AI ser ikke `Backlog`, `Complete` og egendefinerte kolonner. |
| MCP-serversjekk finnes, men klienten er ikke tilkoblet | Serverkontrollen sjekker ikke de eksterne klientinnstillingene. Start klienten på nytt, sjekk konfigurasjonsfilen, banen til `agentboard.exe`, `--project`, `--worker` og generell `--db`. Be om å ringe `get_my_board`. |
| Arbeider frakoblet | Klienten kan ha avsluttet eller kanskje ikke startet MCP ennå. Etter 90 sekunder uten hjerteslag, anses økten som frakoblet. |
| JSON-fil importeres ikke | Sjekk `version: 1`, obligatorisk `title`, eksisterende `Slug` arbeidere og eiendomsnavn. JSON tillater ikke kommentarer eller etterfølgende kommaer. |
| Kan ikke gjenoppbygge `agentboard.exe` | Windows kan ikke erstatte en kjørende EXE. Stopp serveren Ctrl+C og prøv å bygge på nytt. |
| Oppgaver forsvant etter oppdatering | Sjekk at `--db` ikke peker til en annen fil og at du er logget på som samme Windows-bruker. Standardbasen er i `%AppData%\AgentBoard`. |

Hvis feilen ikke er beskrevet, samler du inn den nøyaktige teksten i meldingen, `agentboard version` og Windows-versjonene, og trinnene for å prøve på nytt. Ikke publiser private prosjektdata eller databaseinnhold i en åpen utgave.
