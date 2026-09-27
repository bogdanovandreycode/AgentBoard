# Zacznij tutaj

## Co to jest AgentBoard

Wyobraź sobie zwykłą tablicę z kartami zadań. Tworzysz zadania, wyznaczasz osobę odpowiedzialną i nadzorujesz pracę. Pracownicy AI otrzymują tylko przydzielone im zadania za pośrednictwem MCP i raportują postęp. Ty decydujesz, kiedy zadanie będzie ostatecznie gotowe.

**Projekt** - folder na komputerze i osobna tablica. **Zadanie** - karta z opisem, osobą odpowiedzialną i etapem. **Pracownik** to logiczna nazwa klienta AI. **MCP** to sposób, w jaki klient AI łączy się z AgentBoard. **Scoop** to menedżer instalacji oprogramowania dla systemu Windows.

Nie wymaga kodu. Potrzebujesz systemu Windows, przeglądarki, PowerShell i do działania AI zainstalowanego klienta AI z obsługą lokalnych serwerów MCP.

## Pierwsze uruchomienie

1. Zainstaluj aplikację zgodnie z [instrukcjami](INSTALL.md).
2. Utwórz folder projektu w Eksploratorze, na przykład `C:\Projects\MyFirstProject`.
3. Otwórz ten folder w Eksploratorze. Kliknij pasek adresu, wpisz `powershell` i naciśnij Enter.
4. W oknie, które zostanie otwarte wykonaj:

```powershell
agentboard init
agentboard open
```

5. Otworzy się przeglądarka z adresem `http://127.0.0.1:7337`. Podczas korzystania z tablicy pozostaw okno programu PowerShell otwarte. Zamknięcie okna zatrzyma serwer lokalny, ale zadania pozostaną.

Jeśli polecenie `agentboard` nie zostanie znalezione, zamknij PowerShell i otwórz go ponownie po zainstalowaniu Scoop. Podczas instalacji z pliku ZIP użyj pełnej ścieżki do `agentboard.exe`.

## Pierwsze zadanie

Kliknij **Nowe zadanie**, w razie potrzeby uzupełnij **Tytuł** **Opis**, następnie **Zapisz**. Nowe zadanie w `Backlog` jest dostępne dla ludzi. Aby sztuczna inteligencja zaczęła działać, przydziel pracownika i przekaż zadanie `Features`. Opis pól krok po kroku - [TASKS.md](TASKS.md).

## Pierwszy pracownik

Otwórz **Pracownicy → Dodaj pracownika**. Wybierz klienta AI, którego używasz (np. Codex lub Claude Code), sprawdź nazwę i krótki identyfikator `Slug`, kliknij **Zapisz**. Otwórz utworzonego pracownika: jest tam konfiguracja MCP, przycisk kopiowania i sprawdzenie serwera. Skopiuj konfigurację do klienta AI zgodnie z [WORKERS_MCP.md](WORKERS_MCP.md). Po nawiązaniu połączenia poproś klienta o wywołanie `get_my_board`.

## Co dalej czytać

- [Zadania, testy, historie i kolumny](TASKS.md)
- [Podłączanie pracowników i sprawdzanie MCP](WORKERS_MCP.md)
- [Importuj zadania z JSON](IMPORT.md)
- [Język, motyw, strefa czasowa i aktualizacja](SETTINGS.md)
- [Typowe problemy](TROUBLESHOOTING.md)
