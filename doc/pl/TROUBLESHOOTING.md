# Rozwiązywanie problemów

| Objaw | Co sprawdzić |
| --- | --- |
| Nie znaleziono `agentboard` | Uruchom ponownie PowerShell po Scoop. Podczas instalacji z pliku ZIP użyj pełnej ścieżki do `agentboard.exe`. |
| Strona internetowa pokazująca stary interfejs po kompilacji | Zatrzymaj działający serwer Ctrl+C. Wykonaj `./scripts/build.ps1`, uruchom nowy plik binarny. Vite musi zostać skompilowany przed Go, ponieważ interfejs jest wbudowany w plik EXE. Odśwież stronę Ctrl+F5. |
| Port 7337 zajęty | AgentBoard może już działać. Otwórz `http://127.0.0.1:7337` lub zakończ stary proces. W przypadku innego portu użyj `--addr`. |
| Nie znaleziono projektu | W żądanym folderze wykonaj `agentboard init`. Następnie `agentboard open` z niego lub `agentboard open C:\путь\к\проекту`. |
| Pracownik nie widzi zadania | Zadanie musi być przypisane do tego konkretnego pracownika i znajdować się w `Features`, `In progress`, `Testing` lub `Verification`. AI nie widzi `Backlog`, `Complete` i kolumn niestandardowych. |
| `NO_PROJECTS` | W bazie MCP nie ma żadnych projektów. Wykonaj `agentboard init` w folderze projektu i utwórz proces roboczy; agent może powtórzyć połączenie bez ponownego uruchamiania. Sprawdź, czy MCP i interfejs sieciowy korzystają z tego samego `--db`. |
| `PROJECT_REQUIRED` / `PROJECT_AMBIGUOUS` | Bez `--project` ani folder roboczy klienta, ani folder startowy MCP nie znajdują się dokładnie w tym samym zarejestrowanym projekcie. Uruchom agenta w folderze projektu, zarejestruj go poprzez `agentboard init` lub dodaj `--project`. Nie ma potrzeby uruchamiania interfejsu internetowego. |
| `PROJECT_NOT_FOUND` / `WORKER_NOT_FOUND` / `WORKER_REQUIRED` / `NO_WORKERS` | Sprawdź `--project`, ślimak roboczy w tym projekcie i `--db`. Wiadomość zawiera listę zarejestrowanych projektów lub włączonych pracowników. Niestandardowa konfiguracja MCP bez `--project` wymaga pracownika z tym samym ślimakiem w każdym projekcie. |
| `PROJECT_CONTEXT_MISMATCH` | Folder startowy MCP lub folder roboczy klienta nie pasuje do określonego `--project`. Potwierdź pytanie, jeśli jest to zamierzone, w przeciwnym razie usuń `--project` (automatyczne wykrywanie) lub wskaż właściwy. |
| Kontrola serwera MCP istnieje, ale klient nie jest podłączony | Kontrola serwera nie sprawdza ustawień klienta zewnętrznego. Uruchom ponownie klienta, sprawdź jego plik konfiguracyjny, ścieżkę do `agentboard.exe`, `--project`, `--worker` i ogólny `--db`. Poproś o połączenie z `get_my_board`. |
| Pracownik offline | Klient mógł zakończyć lub jeszcze nie rozpocząć MCP. Po 90 sekundach bez pulsu sesję uznaje się za rozłączoną. |
| Plik JSON nie importuje | Sprawdź `version: 1`, wymagane `title`, istniejących pracowników `Slug` i nazwy właściwości. JSON nie pozwala na komentarze ani końcowe przecinki. |
| Nie można odbudować `agentboard.exe` | System Windows nie może zastąpić działającego pliku EXE. Zatrzymaj serwer Ctrl+C i ponów próbę kompilacji. |
| Zadania zniknęły po aktualizacji | Sprawdź, czy `--db` nie wskazuje innego pliku i czy jesteś zalogowany jako ten sam użytkownik systemu Windows. Domyślna baza znajduje się w `%AppData%\AgentBoard`. |

Jeśli błąd nie jest opisany, zbierz dokładny tekst komunikatu, wersję `agentboard version` i Windows oraz kroki, aby spróbować ponownie. Nie publikuj prywatnych danych projektu ani zawartości bazy danych w otwartym numerze.
