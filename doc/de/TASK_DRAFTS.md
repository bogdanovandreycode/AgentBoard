# Projektmanagement und Aufgabenentwürfe

Auf der Seite **Alle Projekte** werden der registrierte Ordner und die Anzahl der veröffentlichten Projekte angezeigt
Aufgaben und letzte Aktivität. Die Aktivität berücksichtigt Änderungen in Aufgaben, Verlauf usw
Spesenabrechnungen. Mit der Stifttaste können Sie das Projekt umbenennen oder auswählen
vorhandenen Ordner im Systemdialog. Durch Ändern des Pfads werden die Registrierung und die Markierung aktualisiert.
Projekt-IDs, Mitarbeiter, Aufgaben und Einstellungen werden gespeichert, Dateien werden nicht verschoben.
Ordner, die bereits bei einem anderen Projekt registriert sind, können nicht ausgewählt werden.

## Globale Einstellungen

Öffnen Sie **Globale Einstellungen** auf der Projektauswahlseite oder dem Tab
**Globale Einstellungen** in den Einstellungen eines beliebigen Boards. Diese Parameter gelten für alle Projekte
eine AgentBoard-Basis:

- Das automatische Speichern von Entwürfen per Timer ist standardmäßig alle **10 Sekunden** aktiviert.
- **Speichern beim Verlassen des Feldes** ist standardmäßig aktiviert und kann separat konfiguriert werden.
- **MCP-Projektkontextprüfungen** werden global aktiviert und deaktiviert. Ihr Gebiet
  Aktionen und ein separates CLI-Flag werden im Abschnitt [Worker und MCP](WORKERS_MCP.md).

## Lebenszyklusentwurf

Das Dialogfeld „Neue Aufgabe“ speichert nicht abgeschlossene Arbeiten als Entwurf. Statuszeile
unterscheidet zwischen nicht gespeicherten Änderungen, laufendem Speichern, gespeichertem Entwurf und Fehler
Erhaltung. Der Entwurf erscheint nicht auf der Tafel und steht den MCP-Mitarbeitern erst dann zur Verfügung
veröffentlicht.

Verwenden Sie vor der Veröffentlichung die Schaltfläche **Aufgabe erstellen**, um **Titel, Status,
Verantwortlich und Beschreibung**. Zuvor speichert **Entwurf speichern** einen Entwurf, ohne ihn zu erstellen
Aufgaben. Ein leeres Formular wird beim Schließen standardmäßig verworfen. Wenn gelöscht
Wenn Sie den vorhandenen Entwurf in einen leeren Zustand versetzen, wird er beim nächsten Speichern gelöscht.
Durch das Schließen eines nicht leeren Dialogs wird dieser ebenfalls gespeichert.

Mit der Schaltfläche **Entwürfe** auf dem Board können Sie gespeicherte Entwürfe fortsetzen oder löschen.
Mehrere Entwürfe werden unterstützt. Nach einem Neustart öffnet sich das
Entwurf. Entwürfe werden in SQLite gespeichert; Bei jeder Änderung wird die lokale aktualisiert
Kopieren Sie es im Browser, damit nicht gespeicherte Änderungen einen Neustart überstehen oder
vorübergehender Serverausfall. Wenn kein Browserspeicher verfügbar ist, wird ein Dialogfeld mit einer Warnung angezeigt
dass der Entwurf vor dem Schließen auf dem Server gespeichert werden muss.

Beim Schreiben von Entwürfen werden Revisionen verwendet, um zu verhindern, dass ein Fenster die Arbeit überschreibt
ein anderer. Eine verlorene Speicherantwort ist konsistent, wenn bereits eine auf dem Server vorhanden ist
identische Daten. Wenn der Inhalt unterschiedlich ist, wird ein Konflikt angezeigt und
**Als neuen Entwurf wiederherstellen** speichert lokale Inhalte separat. Veröffentlichung
atomar und idempotent: Veröffentlichen Sie denselben Entwurf erneut, nachdem er verloren gegangen ist
Die Antwort gibt dieselbe Aufgabe zurück, anstatt eine Kopie zu erstellen.
