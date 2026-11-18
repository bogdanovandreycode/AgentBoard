# Prosjektledelse og oppgaveutkast

Siden **Alle prosjekter** viser den registrerte mappen, antall publiserte
oppgaver og siste aktivitet. Aktivitet tar hensyn til endringer i oppgaver, historikk og
utgiftsposter. Blyantknappen lar deg gi nytt navn til prosjektet eller velge
eksisterende mappe i systemdialogen. Endring av banen oppdaterer registreringen og markøren;
Prosjekt-IDer, arbeidere, oppgaver og innstillinger lagres, filer flyttes ikke.
Mapper som allerede er registrert i et annet prosjekt kan ikke velges.

## Globale innstillinger

Åpne **Globale innstillinger** fra prosjektvalgsiden eller -fanen
**Globale innstillinger** i innstillingene til ethvert brett. Disse parameterne er felles for alle prosjekter
én AgentBoard-base:

- Automatisk lagring av utkast med timer er aktivert som standard, hvert **10. sekund**.
- **Lagre når du forlater feltet** er aktivert som standard og kan konfigureres separat.
- **MCP-prosjektkontekst**-sjekker er aktivert og deaktivert globalt. Området deres
  handlinger og et eget CLI-flagg er beskrevet i avsnittet [Workers and MCP](WORKERS_MCP.md).

## Utkast til livssyklus

Dialogboksen Ny oppgave lagrer uferdig arbeid som et utkast. Statuslinje
skiller mellom ulagrede endringer, lagring pågår, lagret utkast og feil
bevaring. Utkastet vises ikke på tavlen og er ikke tilgjengelig for MCP-arbeidere før
publisert.

Før du publiserer, bruk knappen **Opprett oppgave** for å fylle ut **Tittel, tilstand,
Ansvarlig og beskrivelse**. Før dette lagrer **Lagre utkast** et utkast uten å opprette
oppgaver. Et tomt skjema forkastes som standard når det lukkes. Hvis ryddet
eksisterende utkast til tom tilstand, slettes det neste gang det lagres.
Hvis du lukker en ikke-tom dialogboks, lagres den også.

**Utkast**-knappen på tavlen lar deg fortsette eller slette lagrede utkast.
Flere utkast støttes. Etter en omstart åpnes
utkast. Utkast lagres i SQLite; Med hver endring oppdateres den lokale
kopier i nettleseren, så ulagrede redigeringer overlever omstart eller
midlertidig serverfeil. Hvis nettleserlagring ikke er tilgjengelig, varsler en dialogboks
at utkastet må lagres på serveren før stenging.

Å skrive utkast bruker revisjoner for å forhindre at ett vindu overskriver arbeid
en annen. Tapt lagringssvar er konsekvent hvis det allerede er en på serveren
identiske data. Hvis innholdet er annerledes, vises en konflikt og
**Gjenopprett som et nytt utkast** lagrer lokalt innhold separat. Publikasjon
atomisk og idempotent: publiser det samme utkastet på nytt etter at det er tapt
respons returnerer den samme oppgaven i stedet for å lage en kopi.
