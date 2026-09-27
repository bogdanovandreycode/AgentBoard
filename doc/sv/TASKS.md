# Arbeta med uppgifter

## Lägg till uppgift

På fliken **Tavla** klickar du på **Ny uppgift**. Fyll i rubrik och beskrivning. Beskrivningen och testinstruktionerna stöder Markdown: markera text och använd formateringsfältet för fetstil, kursiv, länk och lista. Klicka på **Spara**.

Fält:

| Fält | Vad gör |
| --- | --- |
| Titel | Kort namn på uppgiften. |
| Beskrivning | Vad behöver göras och hur man förstår att arbetet är klart. |
| Stat | Arbetsstadium; en ny uppgift börjar vanligtvis vid `Backlog`. |
| Prioritet | `critical`, `high`, `medium` eller `low`. |
| Ansvarig | En person, en specifik arbetare eller utan tidsbokning. |
| Testläge | `AI` - kontrollerar AI; `Human` - mänskliga kontroller; `Hybrid` - båda. |
| Beroenden | Uppgifter som bör slutföras tidigare. |
| AI/Human testinstruktioner | Instruktioner till lämplig granskare. |
| Anpassade egenskaper | Ytterligare fält skapas på fliken **Egenskaper**. |

För en AI-uppgift, skapa först en arbetare, välj den i Ansvarig och överför sedan kortet till `Features`. AI ser bara uppgifter som tilldelats den i fyra steg från `Features` till `Verification`.

## Flytta uppgift

Dra kortet mellan kolumnerna. Du kan växla mellan en vy med alla kolumner och breda kolumner med horisontell rullning. `Backlog` och `Complete` är människokontrollerade. AI kan bara föra fram uppgiften `Features → In progress → Testing → Verification` genom speciella MCP-verktyg. En uppgift med ofullbordad manuell testning bör inte klara AI-verifiering.

## Uppgiftskort

Klicka på ett kort för att se beskrivning, ägare, testinstruktioner och flikar för historik, tester, artefakter och AI-kostnader. **Redigera uppgift** ändrar innehållet. Personens kommentar läggs till den övergripande berättelsen. Sök i de översta filtren sökningar efter titel, beskrivning, ID och arbetare; en separat knapp öppnar en stor sökning.

## Kolumner och egenskaper

På fliken **Inställningar** kan du ändra ordningen på kolumner som är tillgängliga för en person och lägga till dina egna. De fyra AI-stegen är fasta och fortsätter i samma ordning. Användarkolumnen är en plats för uppgifter som skjutits upp av en person: för AI har en sådan uppgift statusen `Backlog`. När en kolumn raderas återgår dess uppgifter till normala `Backlog`.

På fliken **Egenskaper** kan du lägga till fält som text, nummer, flagga, datum, urval och URL. `Human only` döljer fältet från AI; `Agent read` tillåter läsning, `Agent read/write` tillåter också skrivning genom verktyg som stöds. Detta ändrar inte AI:s rättigheter till uppgiftsstegen.
