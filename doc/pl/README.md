# Tabela Agentów

![Podgląd AgentBoard](../../assets/social-preview.png)

**[Pobierz dla Windows](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Witryna dokumentacji](https://bogdanovandreycode.github.io/AgentBoard/) · [Licencja MIT](../../LICENSE)**

AgentBoard to lokalna tablica zadań, na której ludzie i pracownicy AI pracują nad wspólnymi zadaniami, ale mają różne uprawnienia. Aplikacja uruchamia się z jednym plikiem `agentboard.exe`, otwiera interfejs sieciowy w przeglądarce i zapewnia pracownikom oddzielny serwer MCP poprzez `stdio`. Dane pozostają na Twoim komputerze.

**[Zacznij od zera](START_HERE.md) · [Praca z zadaniami](TASKS.md) · [Podłączanie AI przez MCP](WORKERS_MCP.md) · [Importuj JSON](IMPORT.md) · [Ustawienia](SETTINGS.md) · [Rozwiązywanie problemów](TROUBLESHOOTING.md)**

## Funkcje MVP

- Lokalna tablica projektowa z wyszukiwaniem zadań, importem JSON, niestandardowymi kolumnami i właściwościami.
- Dostęp MCP dla konkretnego pracownika z ograniczonymi przejściami AI i ostateczną akceptacją zadań przez człowieka.
- Ogólna historia zadań, instrukcje sprawdzania, artefakty, koszty AI i diagnostyka połączeń pracowników.
- Instalator Windows, manifest ZIP i Scoop; interfejs sieciowy jest wbudowany w plik wykonywalny.

AgentBoard jest przeznaczony dla zaufanego użytkownika lokalnego. Aplikacja nie hostuje projektów w chmurze i nie uruchamia sama klientów AI; w razie potrzeby podłącz klienta kompatybilnego z MCP do pracownika.

## Za pięć minut

1. Pobierz instalator `agentboard-VERSION-windows-amd64-setup.exe` z [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Zaproponuje folder (domyślnie `C:\AI\AgentBoard`) i doda go do `PATH`. Dostępne są również Scoop i ZIP.
2. Otwórz PowerShell w folderze projektu, na przykład `C:\Projects\MyApp`.
3. Wykonaj `agentboard init` (dla ZIP: pełna ścieżka do `agentboard.exe` i `init`).
4. Wykonaj `agentboard open`. Otworzy się `http://127.0.0.1:7337`.
5. Dodaj zadanie za pomocą przycisku **Nowe zadanie**. W przypadku pracownika AI otwórz **Pracownicy → Dodaj pracownika**, wybierz profil klienta i skopiuj konfigurację MCP.

Jeśli nie masz jeszcze folderu projektu, utwórz go w Eksploratorze Windows. Projektem może być dowolny folder, nawet bez Gita i kodu.

## Instalacja poprzez Scoop

W PowerShell z [Scoop](https://scoop.sh/)] już zainstalowanym po wydaniu:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Dla programistów dostępna jest [kompilacja ze źródła ](INSTALL.md). Przepływ pracy związany z wydaniem tworzy manifest instalatora, ZIP i Scoop z SHA-256 z tego samego artefaktu. Aktualizacja wersji zainstalowanej poprzez Scoop: `scoop update agentboard` po dodaniu manifestu do segmentu; szczegóły - [przygotowanie wydania](SCOOP_RELEASE.md).

## Struktura zarządu

`Backlog → Features → In progress → Testing → Verification → Complete`

Osoba może przenosić zadania po tablicy. AI może poruszać się tylko `Features → In progress → Testing → Verification`; ostateczna akceptacja do `Complete` dokonywana jest przez człowieka. Kolumny użytkownika są przeznaczone dla ludzi: zadanie w nich pozostaje w stanie `Backlog` dla MCP. Tryby testowe: AI, ludzki i hybrydowy. Historia, testy, artefakty i koszty AI są dołączone do zadania.

## Zespoły

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` rejestruje folder i zapisuje w nim tylko `.agentboard/project.json`. Dane produkcyjne SQLite znajdują się w katalogu konfiguracyjnym użytkownika Windows (`%AppData%\AgentBoard\agentboard.db`), poza projektem i poza instalacją Scoop. Usunięcie lub aktualizacja pakietu nie powinna usuwać tych danych. Przed przeniesieniem na inny komputer wykonaj kopię bazy danych, gdy AgentBoard jest zatrzymany.

Pracownicy - konta logiczne; Sam AgentBoard nie obsługuje Codexu, Claude ani żadnego innego klienta AI. Klient uruchamia lokalny proces MCP dla konkretnego pracownika. MCP to granica uprawnień aplikacji, a do izolacji plików użyj piaskownicy klienta AI.

## Dla programistów

Stos: Go, SQLite, oficjalny zestaw SDK MCP Go, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Najpierw zbuduj frontend, a następnie Go: `./scripts/build.ps1`. Pliki internetowe są zawarte w pliku binarnym poprzez `go:embed`. Architekturę i API opisano w [doc/ARCHITECTURE.md](ARCHITECTURE.md).

## Licencja

AgentBoard to darmowy projekt typu open source objęty [licencją MIT](../../LICENSE). Komercyjne wykorzystanie, modyfikacja, rozwidlanie i redystrybucja są dozwolone pod warunkiem zachowania informacji o prawach autorskich i licencji.
