# Nastavení

Otevřete **Nastavení** v postranní nabídce vybraného projektu. Změny se ukládají tlačítkem **Uložit nastavení** a ukládají se pro projekt do databáze AgentBoard.

## Jazyk a vzhled

Pole **Jazyk** otevře rozevírací seznam vyhledávání. Ve výchozím nastavení je vybrána možnost **Sledovat systém**, která přebírá jazyk prohlížeče. Dostupné jazyky ze snímků obrazovky: arabština, portugalština (Brazílie), zjednodušená čínština, čeština, dánština, holandština, angličtina, finština, francouzština, němčina, italština, japonština, korejština, norština bokmål, polština, ruština, španělština, švédština, turečtina, ukrajinština, vietnamština; navíc běloruština, rumunština a bulharština. **Sledovat systém** přebírá jazyk prohlížeče. Hlavní podpisy jsou přeloženy ručně, zbývající řádky uživatelského rozhraní mají předběžný automatický překlad. Před zveřejněním je vhodné provést korektury překladů rodilými mluvčími; jména klientů, příkazy, pole JSON a uživatelská data zůstávají bez překladu.

**Barevné schéma**: Tmavá (původní motiv), Světlá, Černá, Ubuntu a Windows. **Časové pásmo** řídí zobrazení dat; data jsou nadále uložena v UTC. **Časové pásmo systému** používá nastavení počítače.

## Sloupce

Čtyři stupně `Features`, `In progress`, `Testing`, `Verification` jsou pevně dané v tomto pořadí: nelze je přejmenovat ani odstranit. Další sloupce lze přeskupit pomocí dostupných šipek. Zadejte název a kliknutím na **Přidat sloupec** vytvořte vlastní sloupec. Je navržen pro lidské úkoly: AI vidí úkol jako `Backlog` a nepřijímá ho prostřednictvím MCP. Odstraněním sloupce při ukládání se jeho úkoly přenesou do běžného `Backlog`.

## Web a MCP

**Obnovení desky** nastaví obnovovací frekvenci desky a karet v sekundách (1–60). **Obnovení pracovníka** aktualizuje stav pracovníků (2–120 sekund). Toto je dotazování webového rozhraní, nikoli frekvence spouštění AI. Pracovníci se nespouštějí automaticky.

Adresa webového serveru se nastavuje při spouštění CLI, například `agentboard open --addr 127.0.0.1:7444`. Pro změnu adresy je nutné restartovat server. Výchozí nastavení je `127.0.0.1:7337`. MCP funguje prostřednictvím samostatného místního příkazu `agentboard mcp --project ... --worker ...` a je nezávislý na webovém portu. Pro nestandardní základ zadejte ve všech příkazech stejný `--db`. Zkopírujte konfiguraci každého klienta z karty pracovníka.
