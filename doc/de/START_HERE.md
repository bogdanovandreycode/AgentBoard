# Beginnen Sie hier

## Was ist AgentBoard

Stellen Sie sich eine normale Tafel mit Aufgabenkarten vor. Sie erstellen Aufgaben, weisen einen Verantwortlichen zu und überwachen die Arbeit. KI-Mitarbeiter erhalten über MCP nur die ihnen zugewiesenen Aufgaben und melden den Fortschritt. Sie entscheiden, wann die Aufgabe endlich fertig ist.

**Projekt** – ein Ordner auf dem Computer und ein separates Board. **Aufgabe** – eine Karte mit Beschreibung, verantwortlicher Person und Phase. **Worker** ist der logische Name des AI-Clients. **MCP** ist die Art und Weise, wie der KI-Client eine Verbindung zum AgentBoard herstellt. **Scoop** ist ein Softwareinstallationsmanager für Windows.

Kein Code erforderlich. Sie benötigen Windows, einen Browser, PowerShell und, damit AI funktioniert, einen installierten AI-Client mit Unterstützung für lokale MCP-Server.

## Erster Start

1. Installieren Sie die Anwendung gemäß der [Anweisungen](INSTALL.md).
2. Erstellen Sie im Explorer einen Projektordner, zum Beispiel `C:\Projects\MyFirstProject`.
3. Öffnen Sie diesen Ordner im Explorer. Klicken Sie in die Adressleiste, geben Sie `powershell` ein und drücken Sie die Eingabetaste.
4. Führen Sie im folgenden Fenster Folgendes aus:

```powershell
agentboard init
agentboard open
```

5. Es öffnet sich ein Browser mit der Adresse `http://127.0.0.1:7337`. Lassen Sie das PowerShell-Fenster geöffnet, während Sie das Whiteboard verwenden. Durch das Schließen des Fensters wird der lokale Server gestoppt, die Aufgaben bleiben jedoch bestehen.

Wenn der Befehl `agentboard` nicht gefunden wird, schließen Sie PowerShell und öffnen Sie es nach der Installation von Scoop erneut. Verwenden Sie bei der Installation von einer ZIP-Datei den vollständigen Pfad zu `agentboard.exe`.

## Erste Aufgabe

Klicken Sie auf **Neue Aufgabe**, geben Sie **Titel** und ggf. **Beschreibung** ein und klicken Sie dann auf **Speichern**. Eine neue Aufgabe in `Backlog` steht Menschen zur Verfügung. Damit die KI mit der Arbeit beginnen kann, weisen Sie einen Arbeiter zu und übertragen Sie die Aufgabe an `Features`. Schritt-für-Schritt-Beschreibung der Felder - [TASKS.md](TASKS.md).

## Erster Arbeiter

Öffnen Sie **Arbeiter → Arbeiter hinzufügen**. Wählen Sie den AI-Client aus, den Sie verwenden (z. B. Codex oder Claude Code), überprüfen Sie den Namen und die Kurz-ID `Slug` und klicken Sie auf **Speichern**. Öffnen Sie den erstellten Worker: Es gibt eine MCP-Konfiguration, eine Schaltfläche zum Kopieren und eine Serverprüfung. Kopieren Sie die Konfiguration gemäß [WORKERS_MCP.md](WORKERS_MCP.md). Sobald die Verbindung hergestellt ist, bitten Sie den Client, `get_my_board` anzurufen.

## Was Sie als nächstes lesen sollten

- [Aufgaben, Tests, Geschichten und Kolumnen](TASKS.md)
- [Arbeiter verbinden und MCP](WORKERS_MCP.md) überprüfen
- [Aufgaben aus JSON](IMPORT.md) importieren
- [Sprache, Thema, Zeitzone und Update](SETTINGS.md)
- [Typische Probleme](TROUBLESHOOTING.md)
