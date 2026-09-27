# Arbejde med opgaver

## Tilføj opgave

På fanen **Tavle** skal du klikke på **Ny opgave**. Udfyld titel og beskrivelse. Beskrivelsen og testinstruktionerne understøtter Markdown: Fremhæv tekst og brug formateringslinjen til fed, kursiv, link og liste. Klik på **Gem**.

Felter:

| Felt | Hvad gør |
| --- | --- |
| Titel | Kort navn på opgaven. |
| Beskrivelse | Hvad skal der gøres, og hvordan man forstår, at arbejdet er klar. |
| Stat | Arbejdsfase; en ny opgave starter normalt ved `Backlog`. |
| Prioritet | `critical`, `high`, `medium` eller `low`. |
| Ansvarlig | En person, en bestemt arbejder eller uden en aftale. |
| Testtilstand | `AI` - kontrollerer AI; `Human` - menneskelige kontroller; `Hybrid` - begge. |
| Afhængigheder | Opgaver, der bør løses tidligere. |
| AI/Human testinstruktioner | Instruktioner til den relevante anmelder. |
| Brugerdefinerede egenskaber | Yderligere felter oprettet på fanen **Egenskaber**. |

For en AI-opgave skal du først oprette en arbejder, vælge den i Ansvarlig og derefter overføre kortet til `Features`. AI ser kun opgaver, der er tildelt den i fire trin fra `Features` til `Verification`.

## Flyt opgave

Træk kortet mellem kolonnerne. Du kan skifte mellem en visning med alle kolonner og brede kolonner med vandret rulning. `Backlog` og `Complete` er menneskekontrollerede. AI kan kun fremme opgave `Features → In progress → Testing → Verification` gennem specielle MCP-værktøjer. En opgave med ufuldendt manuel test bør ikke bestå AI-verifikation.

## Opgavekort

Klik på et kort for at se beskrivelse, ejer, testinstruktioner og faner for historik, test, artefakter og AI-omkostninger. **Rediger opgave** ændrer indholdet. Personens kommentar føjes til den overordnede historie. Søg i de øverste filtre søgninger efter titel, beskrivelse, ID og arbejder; en separat knap åbner en stor søgning.

## Kolonner og egenskaber

På fanen **Indstillinger** kan du ændre rækkefølgen af kolonner, der er tilgængelige for en person og tilføje dine egne. De fire AI-stadier er faste og fortsætter i samme rækkefølge. Brugerkolonnen er et sted for opgaver, der er udskudt af en person: for AI har en sådan opgave status `Backlog`. Når en kolonne slettes, vender dens opgaver tilbage til normal `Backlog`.

På fanen **Egenskaber** kan du tilføje felter som tekst, tal, flag, dato, valg og URL. `Human only` skjuler feltet for AI; `Agent read` tillader læsning, `Agent read/write` tillader også skrivning gennem understøttede værktøjer. Dette ændrer ikke AI'ens rettigheder til opgavetrinene.
