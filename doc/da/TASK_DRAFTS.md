# Projektledelse og opgaveudkast

Siden **Alle projekter** viser den registrerede mappe, antallet af publicerede
opgaver og seneste aktivitet. Aktivitet tager højde for ændringer i opgaver, historik og
udgiftsoptegnelser. Blyantsknappen giver dig mulighed for at omdøbe projektet eller vælge
eksisterende mappe i systemdialogen. Ændring af stien opdaterer registreringen og markøren;
Projekt-id'er, arbejdere, opgaver og indstillinger gemmes, filer flyttes ikke.
Mapper, der allerede er registreret til et andet projekt, kan ikke vælges.

## Globale indstillinger

Åbn **Globale indstillinger** fra siden eller fanen til projektvalg
**Globale indstillinger** i indstillingerne på ethvert board. Disse parametre er fælles for alle projekter
én AgentBoard-base:

- Automatisk lagring af kladder med timer er aktiveret som standard hvert **10. sekund**.
- **Gem når du forlader feltet** er aktiveret som standard og kan konfigureres separat.
- **MCP-projektkontekst**-tjek er aktiveret og deaktiveret globalt. Deres område
  handlinger og et separat CLI-flag er beskrevet i afsnittet [Workers and MCP](WORKERS_MCP.md).

## Udkast til livscyklus

Dialogboksen Ny opgave gemmer ufærdigt arbejde som en kladde. Statuslinje
skelner mellem ikke-gemte ændringer, lagring i gang, gemt kladde og fejl
bevarelse. Udkastet vises ikke på tavlen og er ikke tilgængeligt for MCP-medarbejdere indtil kl
offentliggjort.

Før du udgiver, skal du bruge knappen **Opret opgave** til at udfylde **Titel, tilstand,
Ansvarlig og Beskrivelse**. Før dette gemmer **Gem kladde** en kladde uden at oprette
opgaver. En tom formular kasseres som standard, når den lukkes. Hvis ryddet
eksisterende kladde til en tom tilstand, slettes den næste gang den gemmes.
Lukning af en ikke-tom dialog gemmer den også.

Knappen **Kladder** på tavlen giver dig mulighed for at fortsætte eller slette gemte kladder.
Flere udkast understøttes. Efter en genstart åbnes
udkast. Kladder gemmes i SQLite; Med hver ændring opdateres den lokale
kopier i browseren, så ikke-gemte redigeringer overlever genstart eller
midlertidig serverfejl. Hvis browserlagring ikke er tilgængelig, advarer en dialogboks
at kladden skal gemmes på serveren inden lukning.

Skrivning af kladder bruger revisioner til at forhindre et vindue i at overskrive arbejde
en anden. Mistet lagringssvar er konsistent, hvis der allerede er en på serveren
identiske data. Hvis indholdet er anderledes, vises en konflikt og
**Gendan som en ny kladde** gemmer lokalt indhold separat. Udgivelse
atomisk og idempotent: Genudgiv det samme udkast, når det er tabt
svar returnerer den samme opgave i stedet for at oprette en kopi.
