# Indstillinger

Åbn **Indstillinger** i sidemenuen for det valgte projekt. Ændringer gemmes med knappen **Gem indstillinger** og gemmes for projektet i AgentBoard-databasen.

## Sprog og udseende

Feltet **Sprog** åbner en søgerulleliste. Som standard er indstillingen **Følg system** valgt, som tager browsersproget. Tilgængelige sprog fra skærmbilleder: arabisk, portugisisk (Brasilien), kinesisk forenklet, tjekkisk, dansk, hollandsk, engelsk, finsk, fransk, tysk, italiensk, japansk, koreansk, norsk bokmål, polsk, russisk, spansk, svensk, tyrkisk, ukrainsk, vietnamesisk; desuden hviderussisk, rumænsk og bulgarsk. **Følg system** tager browsersproget. Hovedsignaturerne oversættes manuelt, de resterende UI-linjer har en foreløbig automatisk oversættelse. Før offentlig udgivelse er det tilrådeligt at korrekturlæse oversættelser fra modersmålstalende; klientnavne, kommandoer, JSON-felter og brugerdata forbliver uden oversættelse.

**Farveskema**: Mørk (originalt tema), Lys, Sort, Ubuntu og Windows. **Tidszone** styrer visningen af ​​datoer; dataene bliver fortsat lagret i UTC. **Systemets tidszone** bruger computerindstillinger.

## Kolonner

De fire trin `Features`, `In progress`, `Testing`, `Verification` er faste i denne rækkefølge: de kan ikke omdøbes eller slettes. Andre kolonner kan omarrangeres ved hjælp af de tilgængelige pile. Indtast et navn, og klik på **Tilføj kolonne** for at oprette en brugerdefineret kolonne. Den er designet til menneskelige opgaver: AI ser en opgave som `Backlog` og modtager den ikke via MCP. Sletning af en kolonne, mens du gemmer, overfører dens opgaver til den almindelige `Backlog`.

## Web og MCP

**Brætopdatering** indstiller opdateringshastigheden for brættet og kortene i sekunder (1–60). **Medarbejderopdatering** opdaterer medarbejdernes status (2-120 sekunder). Dette er webgrænsefladeafstemning, ikke AI-triggerfrekvens. Arbejdere starter ikke automatisk.

Webserveradressen indstilles ved start af CLI, for eksempel `agentboard open --addr 127.0.0.1:7444`. For at ændre adressen skal serveren genstartes. Standard er `127.0.0.1:7337`. MCP fungerer gennem en separat lokal kommando `agentboard mcp --project ... --worker ...` og er uafhængig af webporten. For en ikke-standardbase skal du angive den samme `--db` i alle kommandoer. Kopier konfigurationen af ​​hver klient fra arbejderkortet.
