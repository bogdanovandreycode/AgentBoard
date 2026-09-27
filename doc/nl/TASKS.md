# Werken met taken

## Taak toevoegen

Klik op het tabblad **Bestuur** op **Nieuwe taak**. Vul de titel en beschrijving in. De beschrijving en testinstructies ondersteunen Markdown: markeer tekst en gebruik de opmaakbalk voor vet, cursief, link en lijst. Klik op **Opslaan**.

Velden:

| Veld | Wat doet |
| --- | --- |
| Titel | Korte naam van de taak. |
| Beschrijving | Wat er moet gebeuren en hoe u kunt begrijpen dat het werk klaar is. |
| Staat | Werkfase; een nieuwe taak begint meestal bij `Backlog`. |
| Prioriteit | `critical`, `high`, `medium` of `low`. |
| Verantwoordelijk | Een persoon, een specifieke medewerker of zonder afspraak. |
| Testmodus | `AI` - controleert AI; `Human` - menselijke controles; `Hybrid` - beide. |
| Afhankelijkheden | Taken die eerder voltooid moeten zijn. |
| AI/Menselijke testinstructies | Instructies voor de juiste reviewer. |
| Aangepaste eigenschappen | Extra velden gemaakt op het tabblad **Eigenschappen**. |

Voor een AI-taak maakt u eerst een werknemer aan, selecteert u deze in Verantwoordelijk en draagt ​​u vervolgens de kaart over naar `Features`. AI ziet alleen taken die eraan zijn toegewezen in vier fasen, van `Features` tot `Verification`.

## Verplaats taak

Sleep de kaart tussen de kolommen. U kunt schakelen tussen een weergave met alle kolommen en brede kolommen met horizontaal scrollen. `Backlog` en `Complete` worden door mensen bestuurd. AI kan taak `Features → In progress → Testing → Verification` alleen vooruit helpen via speciale MCP-tools. Een taak met onvoltooide handmatige tests mag de AI-verificatie niet doorstaan.

## Taakkaart

Klik op een kaart om de beschrijving, eigenaar, testinstructies en tabbladen voor geschiedenis, tests, artefacten en AI-kosten te bekijken. **Taak bewerken** wijzigt de inhoud. De opmerking van de persoon wordt toegevoegd aan het algemene verhaal. Zoeken in de topfilters zoekt op titel, omschrijving, ID en werknemer; een aparte knop opent een grote zoekopdracht.

## Kolommen en eigenschappen

Op het tabblad **Instellingen** kunt u de volgorde van de kolommen die beschikbaar zijn voor een persoon wijzigen en uw eigen kolommen toevoegen. De vier AI-fasen liggen vast en verlopen in dezelfde volgorde. De gebruikerskolom is een plek voor taken die door een persoon zijn uitgesteld: voor AI heeft zo’n taak de status `Backlog`. Wanneer een kolom wordt verwijderd, keren de taken ervan terug naar de normale `Backlog`.

Op het tabblad **Eigenschappen** kunt u velden toevoegen zoals tekst, getal, vlag, datum, selectie en URL. `Human only` verbergt het veld voor AI; `Agent read` maakt lezen mogelijk, `Agent read/write` maakt ook schrijven via ondersteunde tools mogelijk. Dit verandert niets aan de rechten van de AI op de taakstappen.
