# Arbeide med oppgaver

## Legg til oppgave

På **Tavle**-fanen klikker du på **Ny oppgave**. Fyll inn tittel og beskrivelse. Beskrivelsen og testinstruksjonene støtter Markdown: uthev tekst og bruk formateringslinjen for fet, kursiv, lenke og liste. Klikk på **Lagre**.

Felter:

| Felt | Hva gjør |
| --- | --- |
| Tittel | Kort navn på oppgaven. |
| Beskrivelse | Hva må gjøres og hvordan forstå at arbeidet er klart. |
| Stat | Arbeidsfasen; en ny oppgave starter vanligvis ved `Backlog`. |
| Prioritet | `critical`, `high`, `medium` eller `low`. |
| Ansvarlig | En person, en bestemt arbeider eller uten avtale. |
| Testmodus | `AI` - sjekker AI; `Human` - menneskelige kontroller; `Hybrid` - begge. |
| Avhengigheter | Oppgaver som bør gjennomføres tidligere. |
| AI/Human testinstruksjoner | Instruksjoner til riktig anmelder. |
| Egenskaper | Ytterligere felt opprettet i fanen **Egenskaper**. |

For en AI-oppgave må du først opprette en arbeider, velge den i Ansvarlig, og deretter overføre kortet til `Features`. AI ser bare oppgaver som er tildelt den i fire trinn fra `Features` til `Verification`.

## Flytt oppgave

Dra kortet mellom kolonnene. Du kan bytte mellom en visning med alle kolonner og brede kolonner med horisontal rulling. `Backlog` og `Complete` er menneskekontrollert. AI kan bare fremme oppgave `Features → In progress → Testing → Verification` gjennom spesielle MCP-verktøy. En oppgave med ufullført manuell testing skal ikke bestå AI-verifisering.

## Oppgavekort

Klikk på et kort for å se beskrivelsen, eieren, testinstruksjonene og fanene for historikk, tester, artefakter og AI-kostnader. **Rediger oppgave** endrer innholdet. Personens kommentar legges til den generelle historien. Søk i toppfiltrene søk etter tittel, beskrivelse, ID og arbeider; en egen knapp åpner et stort søk.

## Kolonner og egenskaper

I fanen **Innstillinger** kan du endre rekkefølgen på kolonner som er tilgjengelige for en person og legge til dine egne. De fire AI-stadiene er faste og fortsetter i samme rekkefølge. Brukerkolonnen er et sted for oppgaver utsatt av en person: for AI har en slik oppgave statusen `Backlog`. Når en kolonne slettes, går oppgavene tilbake til normal `Backlog`.

I fanen **Egenskaper** kan du legge til felt som tekst, tall, flagg, dato, utvalg og URL. `Human only` skjuler feltet for AI; `Agent read` tillater lesing, `Agent read/write` tillater også skriving gjennom støttede verktøy. Dette endrer ikke AIs rettigheter til oppgavetrinnene.
