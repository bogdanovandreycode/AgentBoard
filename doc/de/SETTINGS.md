# Einstellungen

Öffnen Sie **Einstellungen** im Seitenmenü des ausgewählten Projekts. Änderungen werden mit der Schaltfläche **Einstellungen speichern** gespeichert und für das Projekt in der AgentBoard-Datenbank gespeichert.

## Sprache und Aussehen

Das Feld **Sprache** öffnet eine Such-Dropdown-Liste. Standardmäßig ist die Option **System folgen** ausgewählt, die die Browsersprache übernimmt. Verfügbare Sprachen aus Screenshots: Arabisch, Portugiesisch (Brasilien), Chinesisch (vereinfacht), Tschechisch, Dänisch, Niederländisch, Englisch, Finnisch, Französisch, Deutsch, Italienisch, Japanisch, Koreanisch, Norwegisch (Bokmål), Polnisch, Russisch, Spanisch, Schwedisch, Türkisch, Ukrainisch, Vietnamesisch; zusätzlich Weißrussisch, Rumänisch und Bulgarisch. **System folgen** übernimmt die Browsersprache. Die Hauptsignaturen werden manuell übersetzt, die restlichen UI-Zeilen verfügen über eine vorläufige automatische Übersetzung. Vor der Veröffentlichung empfiehlt es sich, die Übersetzungen von Muttersprachlern Korrektur lesen zu lassen; Clientnamen, Befehle, JSON-Felder und Benutzerdaten bleiben ohne Übersetzung.

**Farbschema**: Dunkel (Originalthema), Hell, Schwarz, Ubuntu und Windows. **Zeitzone** steuert die Anzeige von Datumsangaben; Die Daten werden weiterhin in UTC gespeichert. **Systemzeitzone** verwendet Computereinstellungen.

## Spalten

Die vier Stufen `Features`, `In progress`, `Testing`, `Verification` sind in dieser Reihenfolge festgelegt: Sie können weder umbenannt noch gelöscht werden. Andere Spalten können mithilfe der verfügbaren Pfeile neu angeordnet werden. Geben Sie einen Namen ein und klicken Sie auf **Spalte hinzufügen**, um eine benutzerdefinierte Spalte zu erstellen. Es ist für menschliche Aufgaben konzipiert: KI sieht eine Aufgabe wie `Backlog` und empfängt sie nicht über MCP. Durch das Löschen einer Spalte beim Speichern werden deren Aufgaben an den regulären `Backlog` übertragen.

## Web und MCP

**Board-Aktualisierung** legt die Aktualisierungsrate des Boards und der Karten in Sekunden fest (1–60). **Worker-Aktualisierung** aktualisiert den Status der Worker (2–120 Sekunden). Hierbei handelt es sich um Webinterface-Abfragen, nicht um die AI-Triggerfrequenz. Worker starten nicht automatisch.

Die Webserveradresse wird beim Starten der CLI festgelegt, zum Beispiel `agentboard open --addr 127.0.0.1:7444`. Um die Adresse zu ändern, muss der Server neu gestartet werden. Der Standardwert ist `127.0.0.1:7337`. MCP arbeitet über einen separaten lokalen Befehl `agentboard mcp --project ... --worker ...` und ist unabhängig vom Web-Port. Geben Sie für eine nicht standardmäßige Basis in allen Befehlen denselben `--db` an. Kopieren Sie die Konfiguration jedes Clients von der Worker-Karte.
