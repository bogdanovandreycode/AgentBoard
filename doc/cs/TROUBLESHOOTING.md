# Řešení problémů

| Symptom | Co zkontrolovat |
| --- | --- |
| `agentboard` nenalezeno | Po Scoop restartujte PowerShell. Při instalaci ze ZIP použijte úplnou cestu k `agentboard.exe`. |
| Webová stránka zobrazující staré rozhraní po sestavení | Zastavte běžící server Ctrl+C. Spusťte `./scripts/build.ps1` a spusťte nový binární soubor. Vite se musí sestavit před Go, protože rozhraní je zabudováno do EXE. Obnovte stránku Ctrl+F5. |
| Port 7337 obsazený | AgentBoard již možná běží. Otevřete `http://127.0.0.1:7337` nebo ukončete starý proces. Pro jiný port použijte `--addr`. |
| Projekt nenalezen | V požadované složce spusťte `agentboard init`. Z něj pak `agentboard open` nebo `agentboard open C:\путь\к\проекту`. |
| Pracovník nevidí úkol | Úloha musí být přiřazena tomuto konkrétnímu pracovníkovi a musí být umístěna v `Features`, `In progress`, `Testing` nebo `Verification`. AI nevidí `Backlog`, `Complete` a vlastní sloupce. |
| Kontrola serveru MCP existuje, ale klient není připojen | Kontrola serveru nekontroluje nastavení externího klienta. Restartujte klienta, zkontrolujte jeho konfigurační soubor, cestu k `agentboard.exe`, `--project`, `--worker` a obecné `--db`. Požádejte o zavolání na `get_my_board`. |
| Pracovník offline | Klient možná ukončil nebo ještě nespustil MCP. Po 90 sekundách bez srdečního tepu je relace považována za odpojenou. |
| Soubor JSON se neimportuje | Zkontrolujte `version: 1`, požadované `title`, stávající pracovníky `Slug` a názvy vlastností. JSON nepovoluje komentáře ani koncové čárky. |
| Nelze přestavět `agentboard.exe` | Systém Windows nemůže nahradit běžící EXE. Zastavte server Ctrl+C a zkuste sestavení znovu. |
| Úkoly zmizely po aktualizaci | Zkontrolujte, zda `--db` neodkazuje na jiný soubor a zda jste přihlášeni jako stejný uživatel Windows. Výchozí základna je v `%AppData%\AgentBoard`. |

Pokud chyba není popsána, shromážděte přesný text zprávy, verze `agentboard version` a Windows a postup opakujte. Nezveřejňujte data soukromého projektu nebo obsah databáze v otevřeném problému.
