# Probleemoplossing

| Symptoom | Wat te controleren |
| --- | --- |
| `agentboard` niet gevonden | Start PowerShell opnieuw na Scoop. Wanneer u installeert vanuit een ZIP, gebruik dan het volledige pad naar `agentboard.exe`. |
| Webpagina met oude interface na build | Stop de actieve server met Ctrl+C. Voer `./scripts/build.ps1` uit en start een nieuw binair bestand. Vite moet vóór Go bouwen omdat de interface in de EXE is ingebouwd. Vernieuw de pagina Ctrl+F5. |
| Poort 7337 bezet | AgentBoard is mogelijk al actief. Open `http://127.0.0.1:7337` of beëindig het oude proces. Gebruik voor andere poorten `--addr`. |
| Project niet gevonden | Voer in de gewenste map `agentboard init` uit. Vervolgens `agentboard open` ervan of `agentboard open C:\путь\к\проекту`. |
| De werknemer ziet de taak niet | De taak moet aan deze specifieke werknemer worden toegewezen en zich bevinden in `Features`, `In progress`, `Testing` of `Verification`. AI ziet `Backlog`, `Complete` en aangepaste kolommen niet. |
| `NO_PROJECTS` | Er zijn geen projecten in de MCP-database. Voer `agentboard init` uit in de projectmap en maak een werknemer aan; de agent kan het gesprek herhalen zonder opnieuw op te starten. Controleer of MCP en webinterface dezelfde `--db` gebruiken. |
| `PROJECT_REQUIRED` / `PROJECT_AMBIGUOUS` | Zonder `--project` bevinden noch de werkmap van de klant, noch de MCP-opstartmap zich in precies hetzelfde geregistreerde project. Voer de agent uit in de projectmap, registreer deze via `agentboard init` of voeg `--project` toe. Het is niet nodig om de webinterface te starten. |
| `PROJECT_NOT_FOUND` / `WORKER_NOT_FOUND` / `WORKER_REQUIRED` / `NO_WORKERS` | Bekijk `--project`, de werkslak in dit project en `--db`. Het bericht vermeldt geregistreerde projecten of ingeschakelde werknemers. Voor een aangepaste MCP-configuratie zonder `--project` is in elk project een medewerker nodig met dezelfde slug. |
| `PROJECT_CONTEXT_MISMATCH` | De MCP-opstartmap of clientwerkmap komt niet overeen met de opgegeven `--project`. Bevestig de vraag als dit opzettelijk is, verwijder anders `--project` (automatische detectie) of geef de juiste aan. |
| MCP-servercontrole bestaat, maar de client is niet verbonden | Bij de servercontrole worden de externe clientinstellingen niet gecontroleerd. Start de client opnieuw op, controleer het configuratiebestand, het pad naar `agentboard.exe`, `--project`, `--worker` en algemeen `--db`. Vraag om `get_my_board` te bellen. |
| Werknemer offline | Het kan zijn dat de client MCP heeft beëindigd of nog niet heeft gestart. Na 90 seconden zonder hartslag wordt de sessie als verbroken beschouwd. |
| JSON-bestand importeert niet | Controleer `version: 1`, vereiste `title`, bestaande `Slug`-werknemers en eigendomsnamen. JSON staat geen opmerkingen of komma's toe. |
| Kan `agentboard.exe` | niet opnieuw opbouwen Windows kan een actieve EXE niet vervangen. Stop de server met Ctrl+C en probeer de build opnieuw. |
| Taken verdwenen na update | Controleer of `--db` niet naar een ander bestand verwijst en of u bent ingelogd als dezelfde Windows-gebruiker. De standaardbasis bevindt zich in `%AppData%\AgentBoard`. |

Als de fout niet wordt beschreven, verzamel dan de exacte tekst van het bericht, de `agentboard version`- en Windows-versies en de stappen om het opnieuw te proberen. Publiceer geen privéprojectgegevens of database-inhoud in een open probleem.
