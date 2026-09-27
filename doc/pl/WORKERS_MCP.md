# Pracownicy i połączenie MCP

## Krok 1. Utwórz pracownika

W AgentBoard otwórz **Pracownicy → Dodaj pracownika**. Wybierz profil klienta: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama via OpenCode, VS Code Copilot lub inny klient MCP. Profil wypełnia typowe funkcjonalności (kod, testy, Git); możesz je zmienić. Nazwa jest widoczna na płytce, a `Slug` to krótki identyfikator bez spacji dla komendy MCP. Kliknij **Zapisz**.

Jeden pracownik odpowiada jednej osobowości AI. Twórz różnych pracowników dla różnych klientów lub zespołów. Profil zdolności opisuje specjalizację, ale nie rozszerza uprawnień AI na etapy zadaniowe.

## Krok 2. Skopiuj konfigurację

Otwórz utworzonego pracownika. Blok **Diagnostyka MCP** pokazuje fragment konfiguracji i plik, w którym należy go dodać. Kliknij **Kopiuj konfigurację MCP**. Jeżeli plik już istnieje, dodaj sugerowany serwer do istniejącego obiektu `mcpServers`/`servers`/`mcp` bez usuwania pozostałych serwerów.

Główne polecenie wygląda następująco:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project` powinien wskazywać folder zarejestrowany za pomocą `agentboard init`. `--worker` — `Slug` utworzonego pracownika. Klient MCP sam uruchamia to polecenie, gdy potrzebuje narzędzi. W przeglądarce AgentBoard serwer WWW może działać oddzielnie.

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

Dla Codexu i Claude Code przykład z projektem `C:\Projects\MyFirstProject` i pracownikiem `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

W JSON ukośnik odwrotny systemu Windows jest podwojony; gotowy fragment z interfejsu robi to automatycznie. Jeśli używasz `--db` z niestandardową podstawą, dodaj ją do konfiguracji MCP `args` i określ tę samą ścieżkę, co przy uruchamianiu `open`.

**Ollama** zapewnia model lokalny, ale nie zastępuje klienta MCP. Profil **Ollama via OpenCode** generuje konfigurację MCP dla OpenCode; osobno skonfiguruj OpenCode w modelu Ollama. Odpowiedni jest również inny klient kompatybilny z MCP z Ollamą.

## Krok 4: Sprawdź swoje połączenie

1. Otwórz kartę pracownika i kliknij **Sprawdź ponownie**. **Kontrola serwera** powinna pokazać liczbę narzędzi MCP. Jest to wewnętrzna kontrola protokołu i wykrywanie narzędzi serwerowych.
2. Uruchom lub zrestartuj klienta AI po dodaniu pliku konfiguracyjnego. Poproś go, aby zadzwonił do `get_my_board`.
3. Na AgentBoard pojawi się komunikat **Klient podłączony**, a na liście pojawi się nowa sesja. Tylko to potwierdza połączenie Twojego klienta. Jeśli narzędzia są dostępne, ale klient nie jest podłączony, sprawdź ścieżkę do programu, nazwę pliku konfiguracyjnego i jego składnię JSON/TOML.

Rozpocznij sesję roboczą z `get_my_board`. AI nie otrzymuje zadań od `Backlog` i `Complete`, nawet jeśli zna ich identyfikator. Sztuczna inteligencja nie może udawać człowieka i nie ma ogólnego polecenia „przesuń się gdziekolwiek”. Praca z systemem plików poza AgentBoard zależy od możliwości i piaskownicy wybranego klienta.

Oficjalne instrukcje klienta: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
