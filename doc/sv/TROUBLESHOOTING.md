# Problemlösning

| Symptom | Vad ska man kontrollera |
| --- | --- |
| `agentboard` hittades inte | Starta om PowerShell efter Scoop. När du installerar från en ZIP, använd den fullständiga sökvägen till `agentboard.exe`. |
| Webbsida som visar gammalt gränssnitt efter build | Stoppa den körande servern Ctrl+C. Kör `./scripts/build.ps1`, starta en ny binär. Vite måste byggas före Go eftersom gränssnittet är inbyggt i EXE. Uppdatera sidan Ctrl+F5. |
| Hamn 7337 upptagen | AgentBoard kanske redan körs. Öppna `http://127.0.0.1:7337` eller avsluta den gamla processen. För andra portar använd `--addr`. |
| Projektet hittades inte | I önskad mapp, kör `agentboard init`. Sedan `agentboard open` från den eller `agentboard open C:\путь\к\проекту`. |
| Arbetaren ser inte uppgiften | Uppgiften måste tilldelas just den här arbetaren och finnas i `Features`, `In progress`, `Testing` eller `Verification`. AI ser inte `Backlog`, `Complete` och anpassade kolumner. |
| MCP-serverkontroll finns, men klienten är inte ansluten | Serverkontrollen kontrollerar inte de externa klientinställningarna. Starta om klienten, kontrollera dess konfigurationsfil, sökväg till `agentboard.exe`, `--project`, `--worker` och allmänna `--db`. Be att få ringa `get_my_board`. |
| Arbetare offline | Klienten kan ha avslutat eller kanske inte startat MCP ännu. Efter 90 sekunder utan hjärtslag anses sessionen vara frånkopplad. |
| JSON-fil importeras inte | Kontrollera `version: 1`, obligatoriska `title`, befintliga `Slug`-arbetare och egendomsnamn. JSON tillåter inte kommentarer eller avslutande kommatecken. |
| Det går inte att bygga om `agentboard.exe` | Windows kan inte ersätta ett körande EXE. Stoppa servern Ctrl+C och försök bygga igen. |
| Uppgifter försvann efter uppdatering | Kontrollera att `--db` inte pekar på en annan fil och att du är inloggad som samma Windows-användare. Standardbasen är i `%AppData%\AgentBoard`. |

Om felet inte beskrivs, samla in den exakta texten i meddelandet, `agentboard version`- och Windows-versionerna och steg för att försöka igen. Publicera inte privata projektdata eller databasinnehåll i ett öppet nummer.
