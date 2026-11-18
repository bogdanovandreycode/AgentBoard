# Pracownicy i połączenie MCP

## Krok 1. Utwórz pracownika

W AgentBoard otwórz **Pracownicy → Dodaj pracownika**. Wybierz profil klienta: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot lub inny klient MCP. Profil wypełnia typowe funkcjonalności (kod, testy, Git); możesz je zmienić. Nazwa jest widoczna na płytce, a `Slug` to krótki identyfikator bez spacji dla komendy MCP. Kliknij **Zapisz**.

Jeden pracownik odpowiada jednej osobowości AI. Twórz różnych pracowników dla różnych klientów lub zespołów. Profil zdolności opisuje specjalizację, ale nie rozszerza uprawnień AI na etapy zadaniowe.

## Krok 2. Skopiuj konfigurację

Otwórz utworzonego pracownika. Blok **Diagnostyka MCP** pokazuje fragment konfiguracji i plik, w którym należy go dodać. Kliknij **Kopiuj konfigurację MCP**. Jeżeli plik już istnieje, dodaj sugerowany serwer do istniejącego obiektu `mcpServers`/`servers`/`mcp` bez usuwania pozostałych serwerów.

Główne polecenie wygląda następująco:

```text
agentboard mcp --worker codex
```

`--worker` — `Slug` utworzonego pracownika. Klient MCP sam uruchamia to polecenie, gdy potrzebuje narzędzi.

### MCP działa bez interfejsu internetowego

`agentboard mcp` to oddzielny proces stdio, który sam otwiera bazę danych AgentBoard. Nie potrzebuje `agentboard open`/`serve`, nie jest to zależne od folderu, z którego uruchamiany jest interfejs sieciowy i od tego, ile instancji interfejsu sieciowego jest uruchomionych. Jednocześnie może działać wielu agentów, każdy w swoim własnym projekcie: każdy klient uruchamia własny proces MCP i wszyscy bezpiecznie korzystają z tej samej bazy danych SQLite.

### Jak wybierany jest projekt

Serwer MCP działa w dowolnym folderze. Przy pierwszym wywołaniu narzędzia jest ono powiązane z projektem:

1. `--project PATH`, jeśli określono (odpowiedni jest dowolny folder w zarejestrowanym projekcie);
2. w przeciwnym wypadku foldery robocze Klienta (root), jeśli Klient je zgłosi;
3. w przeciwnym razie folder, w którym klient uruchomił serwer MCP.

`get_my_board` zwraca wybrany projekt i pracownika do `session` wraz z `binding_source`. Jeśli foldery robocze klienta ulegną zmianie, następne wywołanie narzędzia ponownie wykona powiązanie. Bez `--worker` używany jest jedyny włączony proces roboczy projektu.

Jeśli powiązanie nie jest możliwe, każde narzędzie zwraca błąd, który agent musi przekazać użytkownikowi; dane zadania nie są odczytywane w tym przypadku:

| Kod | Znaczenie |
| --- | --- |
| `NO_PROJECTS` | Brak projektów w bazie danych. Wykonaj `agentboard init` w folderze projektu i utwórz proces roboczy. Nie ma potrzeby ponownego uruchamiania. |
| `PROJECT_REQUIRED` | Folder roboczy znajduje się poza zarejestrowanymi projektami. Wiadomość zawiera listę projektów. |
| `PROJECT_NOT_FOUND` | `--project` nie znajduje się w zarejestrowanym projekcie. |
| `PROJECT_AMBIGUOUS` | Foldery robocze klienta obejmują wiele projektów. Dodaj `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Utwórz, włącz lub wybierz pracownika za pomocą `--worker`. |

Konfiguracje **Codex** i **Gemini CLI** są niestandardowe i współdzielone we wszystkich folderach, więc nie mają `--project`: Utwórz proces roboczy z tym samym ślimakiem w każdym projekcie, w którym działa ten klient. Pliki na poziomie projektu (Claude Code, Cursor, OpenCode, VS Code) zawierają `--project` swojego projektu.

### Jawny projekt i sprawdzanie kontekstu

Jeśli określono `--project`, a folder startowy lub folder roboczy klienta znajduje się w innym projekcie, narzędzia proszą o potwierdzenie poprzez wywołanie MCP. Klienci bez wywoływania formularza otrzymują `PROJECT_CONTEXT_MISMATCH` przed uzyskaniem dostępu do zadań. Popraw konfigurację, usuń `--project` w celu automatycznego wykrywania lub wyraźnie potwierdź, że zamierzasz pracować z innym projektem z flagą `--allow-project-mismatch`. Zmiana folderów roboczych resetuje poprzednie potwierdzenia. Projekty definiowane automatycznie nie wymagają potwierdzenia.

**Ustawienia globalne → Kontekst projektu MCP** kontroluje te kontrole dla wszystkich projektów i pracowników tej samej bazy danych. Domyślnie włączone. Jeśli ta opcja jest wyłączona, MCP używa określonego projektu bez potwierdzeń i błędów niezgodności. Już otwarte sesje zastosują ustawienie przy następnym wywołaniu narzędzia; w dalszym ciągu przestrzegane są zasady dotyczące widoczności zadań, przejść i testów na ludziach.

Po sprawdzeniu blok diagnostyczny pokazuje ścieżkę bezwzględną do działającego `agentboard.exe`. Jest to szczególnie przydatne podczas instalacji z pliku ZIP. Podczas instalacji za pośrednictwem Scoop możesz użyć komendy `agentboard`, jeśli klient zobaczy ten sam `PATH`.

## Krok 3: Dodaj serwer do swojego klienta

Interfejs workera posiada już gotowy fragment. Poniżej znajduje się wyjaśnienie, gdzie jest ono używane:

| Klient | Gdzie wstawić | Jak sprawdzić po stronie klienta |
| --- | --- | --- |
| Kodeks | `%USERPROFILE%\.codex\config.toml`, przekrój `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Kod Claude'a | `.mcp.json` w folderze projektu | `claude mcp list` |
| Bliźnięta CLI | `%USERPROFILE%\.gemini\settings.json`, obiekt `mcpServers` | `/mcp list` w Gemini CLI |
| Kursor | Projekt `.cursor\mcp.json` | lista serwerów MCP w Ustawieniach kursora |
| Otwarty kod | Projekt `opencode.json`, obiekt `mcp` | lista narzędzi MCP w OpenCode |
| Drugi pilot kodu VS | Projekt `.vscode\mcp.json`, obiekt `servers` | polecenie **MCP: Lista serwerów** |

Przykłady dla Codexu (konfiguracja użytkownika, projekt wykryty automatycznie) i Claude Code (plik projektu) z projektem `C:\Projects\MyFirstProject` i workerem `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

W JSON ukośnik odwrotny systemu Windows jest podwojony; gotowy fragment z interfejsu robi to automatycznie. Wygenerowane konfiguracje zawierają rzeczywisty kod `--db` serwera WWW, więc protokoły HTTP i MCP korzystają z tej samej bazy, nawet jeśli środowisko klienta jest inne.

**Ollama** zapewnia model lokalny, ale nie zastępuje klienta MCP. Profil **Ollama via OpenCode** generuje konfigurację MCP dla OpenCode; osobno skonfiguruj OpenCode w modelu Ollama. Odpowiedni jest również inny klient kompatybilny z MCP z Ollamą.

## Krok 4: Sprawdź swoje połączenie

1. Otwórz kartę pracownika i kliknij **Sprawdź ponownie**. **Sprawdzanie serwera** uruchamia prawdziwy plik wykonywalny stdio z tym projektem, procesem roboczym i bazą, sprawdza uzgadnianie, listę narzędzi oraz lokalizację projektu i procesu roboczego.
2. Uruchom lub zrestartuj klienta AI po dodaniu pliku konfiguracyjnego. Poproś go, aby zadzwonił do `get_my_board`.
3. Na AgentBoard pojawi się komunikat **Klient podłączony**, a na liście pojawi się nowa sesja. Tylko to potwierdza połączenie Twojego klienta. Jeśli narzędzia są dostępne, ale klient nie jest podłączony, sprawdź ścieżkę do programu, nazwę pliku konfiguracyjnego i jego składnię JSON/TOML.

Rozpocznij sesję roboczą z `get_my_board`. AI nie otrzymuje zadań od `Backlog` i `Complete`, nawet jeśli zna ich identyfikator. Sztuczna inteligencja nie może udawać człowieka i nie ma ogólnego polecenia „przesuń się gdziekolwiek”. Praca z systemem plików poza AgentBoard zależy od możliwości i piaskownicy wybranego klienta.

Każde narzędzie deklaruje `title` i wszystkie cztery wskazówki dotyczące zachowania MCP (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) jako jawne wartości logiczne. Narzędzia do czytania są tylko do odczytu i idempotentne; narzędzia do nagrywania jedynie dodają rekordy lub przesuwają zadania do przodu, więc żadne z nich nie zakłóca pracy ani nie wykracza poza lokalną bazę AgentBoard.

Oficjalne instrukcje klienta: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
