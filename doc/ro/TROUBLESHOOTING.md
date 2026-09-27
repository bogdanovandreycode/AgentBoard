# Rezolvarea problemelor

| Simptom | Ce trebuie verificat |
| --- | --- |
| `agentboard` nu a fost găsit | Reporniți PowerShell după Scoop. Când instalați dintr-un fișier ZIP, utilizați calea completă către `agentboard.exe`. |
| Pagina web care arată interfața veche după compilare | Opriți serverul care rulează Ctrl+C. Executați `./scripts/build.ps1`, lansați un nou binar. Vite trebuie să fie construit înainte de Go, deoarece interfața este încorporată în EXE. Actualizează pagina Ctrl+F5. |
| Portul 7337 ocupat | AgentBoard poate rula deja. Deschideți `http://127.0.0.1:7337` sau încheiați vechiul proces. Pentru alt port utilizați `--addr`. |
| Proiectul nu a fost găsit | În folderul dorit, executați `agentboard init`. Apoi `agentboard open` de la acesta sau `agentboard open C:\путь\к\проекту`. |
| Lucrătorul nu vede sarcina | Sarcina trebuie să fie atribuită acestui lucrător și să fie localizată în `Features`, `In progress`, `Testing` sau `Verification`. AI nu vede `Backlog`, `Complete` și coloanele personalizate. |
| Verificarea serverului MCP există, dar clientul nu este conectat | Verificarea serverului nu verifică setările clientului extern. Reporniți clientul, verificați fișierul de configurare al acestuia, calea către `agentboard.exe`, `--project`, `--worker` și general `--db`. Cereți să sunați la `get_my_board`. |
| Lucrător Offline | Este posibil ca clientul să fi terminat sau să nu fi început încă MCP. După 90 de secunde fără bătăi ale inimii, sesiunea este considerată deconectată. |
| Fișierul JSON nu se importă | Verificați `version: 1`, `title` necesar, lucrătorii `Slug` existenți și numele proprietăților. JSON nu permite comentarii sau virgule finale. |
| Nu se poate reconstrui `agentboard.exe` | Windows nu poate înlocui un EXE care rulează. Opriți serverul Ctrl+C și încercați din nou construirea. |
| Sarcinile au dispărut după actualizare | Verificați dacă `--db` nu indică un alt fișier și că sunteți conectat ca același utilizator Windows. Baza implicită este în `%AppData%\AgentBoard`. |

Dacă eroarea nu este descrisă, colectați textul exact al mesajului, versiunile `agentboard version` și Windows și pașii pentru a reîncerca. Nu publicați date private de proiect sau conținutul bazei de date într-o problemă deschisă.
