# Innstillinger

Åpne **Innstillinger** i sidemenyen til det valgte prosjektet. Endringer lagres med knappen **Lagre innstillinger** og lagres for prosjektet i AgentBoard-databasen.

## Språk og utseende

Feltet **Språk** åpner en rullegardinliste for søk. Som standard er alternativet **Følg system** valgt, som tar nettleserspråket. Tilgjengelige språk fra skjermbilder: arabisk, portugisisk (Brasil), forenklet kinesisk, tsjekkisk, dansk, nederlandsk, engelsk, finsk, fransk, tysk, italiensk, japansk, koreansk, norsk bokmål, polsk, russisk, spansk, svensk, tyrkisk, ukrainsk, vietnamesisk; i tillegg hviterussisk, rumensk og bulgarsk. **Følg systemet** tar nettleserspråket. Hovedsignaturene oversettes manuelt, de resterende UI-linjene har en foreløpig automatisk oversettelse. Før offentlig utgivelse er det tilrådelig å korrekturlese oversettelser fra morsmål; klientnavn, kommandoer, JSON-felt og brukerdata forblir uten oversettelse.

**Fargeskjema**: Mørk (originalt tema), lys, svart, Ubuntu og Windows. **Tidssone** kontrollerer visningen av datoer; dataene fortsetter å bli lagret i UTC. **Systemets tidssone** bruker datamaskininnstillinger.

## Kolonner

De fire stadiene `Features`, `In progress`, `Testing`, `Verification` er fikset i denne rekkefølgen: de kan ikke gis nytt navn eller slettes. Andre kolonner kan omorganiseres ved å bruke de tilgjengelige pilene. Skriv inn et navn og klikk på **Legg til kolonne** for å opprette en egendefinert kolonne. Den er designet for menneskelige oppgaver: AI ser en oppgave som `Backlog` og mottar den ikke gjennom MCP. Hvis du sletter en kolonne mens du lagrer, overføres oppgavene til den vanlige `Backlog`.

## Web og MCP

**Brettoppdatering** angir oppdateringsfrekvensen for brettet og kortene i sekunder (1–60). **Arbeideroppdatering** oppdaterer statusen til arbeidere (2–120 sekunder). Dette er nettgrensesnittspørring, ikke AI-utløserfrekvens. Arbeidere starter ikke automatisk.

Webserveradressen settes når du starter CLI, for eksempel `agentboard open --addr 127.0.0.1:7444`. For å endre adressen må serveren startes på nytt. Standard er `127.0.0.1:7337`. MCP opererer gjennom en egen lokal kommando `agentboard mcp --project ... --worker ...` og er uavhengig av nettporten. For en ikke-standard base, spesifiser den samme `--db` i alle kommandoer. Kopier konfigurasjonen til hver klient fra arbeiderkortet.
