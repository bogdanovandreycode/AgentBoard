# Mit Aufgaben arbeiten

## Aufgabe hinzufügen

Klicken Sie auf der Registerkarte **Board** auf **Neue Aufgabe**. Geben Sie den Titel und die Beschreibung ein. Die Beschreibung und Testanweisungen unterstützen Markdown: Markieren Sie Text und verwenden Sie die Formatierungsleiste für Fett, Kursiv, Link und Liste. Klicken Sie auf **Speichern**.

Felder:

| Feld | Was bedeutet |
| --- | --- |
| Titel | Kurzname der Aufgabe. |
| Beschreibung | Was muss getan werden und wie erkennt man, dass die Arbeit fertig ist? |
| Staat | Arbeitsphase; Eine neue Aufgabe beginnt normalerweise bei `Backlog`. |
| Priorität | `critical`, `high`, `medium` oder `low`. |
| Verantwortlich | Eine Person, ein bestimmter Mitarbeiter oder ohne Termin. |
| Testmodus | `AI` – prüft die KI; `Human` – menschliche Kontrollen; `Hybrid` - beides. |
| Abhängigkeiten | Aufgaben, die früher erledigt werden sollten. |
| KI/Mensch-Testanweisungen | Anweisungen an den entsprechenden Gutachter. |
| Benutzerdefinierte Eigenschaften | Zusätzliche Felder, die auf der Registerkarte **Eigenschaften** erstellt wurden. |

Erstellen Sie für eine KI-Aufgabe zunächst einen Arbeiter, wählen Sie ihn unter „Verantwortlich“ aus und übertragen Sie dann die Karte an `Features`. Die KI sieht die ihr zugewiesenen Aufgaben nur in vier Stufen von `Features` bis `Verification`.

## Aufgabe verschieben

Ziehen Sie die Karte zwischen den Spalten. Sie können zwischen einer Vollspaltenansicht und breiten Spalten mit horizontalem Scrollen wechseln. `Backlog` und `Complete` werden vom Menschen kontrolliert. KI kann die Aufgabe `Features → In progress → Testing → Verification` nur durch spezielle MCP-Tools vorantreiben. Eine Aufgabe mit unvollständigen manuellen Tests sollte die KI-Verifizierung nicht bestehen.

## Aufgabenkarte

Klicken Sie auf eine Karte, um Beschreibung, Besitzer, Testanweisungen und Registerkarten für Verlauf, Tests, Artefakte und KI-Kosten anzuzeigen. **Aufgabe bearbeiten** ändert den Inhalt. Der Kommentar der Person wird zur Gesamtgeschichte hinzugefügt. Die Suche in den oberen Filtern erfolgt nach Titel, Beschreibung, ID und Arbeiter; Eine separate Schaltfläche öffnet eine große Suche.

## Spalten und Eigenschaften

Auf der Registerkarte **Einstellungen** können Sie die Reihenfolge der für eine Person verfügbaren Spalten ändern und eigene hinzufügen. Die vier KI-Stufen sind festgelegt und laufen in der gleichen Reihenfolge ab. Die Benutzerspalte ist ein Ort für Aufgaben, die von einer Person verschoben wurden: Für AI hat eine solche Aufgabe den Status `Backlog`. Wenn eine Spalte gelöscht wird, kehren ihre Aufgaben zum Normalzustand zurück `Backlog`.

Auf der Registerkarte **Eigenschaften** können Sie Felder wie Text, Nummer, Flagge, Datum, Auswahl und URL hinzufügen. `Human only` verbirgt das Feld vor der KI; `Agent read` ermöglicht das Lesen, `Agent read/write` ermöglicht auch das Schreiben über unterstützte Tools. Die Rechte der KI an den Aufgabenschritten werden dadurch nicht geändert.
