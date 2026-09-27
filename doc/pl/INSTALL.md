# Instalacja i uruchomienie w systemie Windows

## Instalator (zalecane)

Pobierz `agentboard-VERSION-windows-amd64-setup.exe` ze strony [Wersje ](https://github.com/bogdanovandreycode/AgentBoard/releases). Kreator instalacji zasugeruje folder; wartość domyślna to `C:\AI\AgentBoard`. Skopiuje `agentboard.exe` i dokumentację, utworzy skrót i doda wybrany folder do systemu `PATH`. Po instalacji otwórz nowy terminal, aby udostępnić polecenie `agentboard`.

W folderze projektu uruchom:

```powershell
agentboard init
agentboard open
```

Interfejs jest wbudowany w `agentboard.exe`; Nie jest wymagana osobna instalacja Go ani Node.js. Dane są przechowywane w `%AppData%\AgentBoard` i są zachowywane podczas aktualizacji lub dezinstalacji programu. Odinstalowanie poprzez „Zainstalowane aplikacje” usuwa skróty i wpis z `PATH`.

## Instalowanie Scoopa

Jeśli Scoop nie jest jeszcze zainstalowany, otwórz PowerShell jako zwykły użytkownik i postępuj zgodnie z [oficjalnymi instrukcjami Scoop](https://scoop.sh/). Jeśli na komputerze firmowym obowiązują ograniczenia, skontaktuj się z administratorem; AgentBoard można także uruchomić z pliku ZIP bez Scoop.

Alternatywnie zainstaluj AgentBoard poprzez Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Dostępność manifestu w wydaniu GitHub można sprawdzić na stronie [strona Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Po dodaniu manifestu do Scoop, wiadro można zainstalować według nazwy wiadra i zaktualizować za pomocą komendy `scoop update agentboard`.

## ZIP bez miarki

Pobierz `agentboard-VERSION-windows-amd64.zip` z Releases, rozpakuj np. do `C:\Tools\AgentBoard`. W PowerShell w folderze projektu:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

W przypadku klienta AI określ pełną ścieżkę do `agentboard.exe` w jego konfiguracji MCP, jeśli program nie znajduje się w `PATH`.

## Kompiluj ze źródła

Zainstaluj wersję Go z `go.mod` i Node.js 22 lub nowszą. W PowerShell w katalogu głównym repozytorium:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` montuje frontend w `internal/webui/dist`, przeprowadza testy Go i montuje jeden `agentboard.exe`. Kolejność jest ważna: interfejs jest wbudowany w plik binarny, gdy budowany jest Go. Jeśli stary `agentboard.exe open` jest uruchomiony, zatrzymaj go przed odbudowaniem (Ctrl+C), w przeciwnym razie system Windows nie pozwoli na zastąpienie pliku.

## Gdzie są dane

- `%AppData%\AgentBoard\agentboard.db` - zadania, projekty, pracownicy i ustawienia. Możesz określić inny plik z flagą `--db`, ale `init`, `open`/`serve` i `mcp` muszą mieć **tę samą ścieżkę**.
- `<ваш проект>\.agentboard\project.json` — identyfikator projektu. Ten plik nie zawiera zadań.
- Serwer domyślnie nasłuchuje tylko `127.0.0.1:7337`. Przed ścieżką projektu podaj inny adres `--addr`: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` otwiera stronę projektu i ponownie wykorzystuje już działający serwer pod tym adresem. Jeśli w jednej przeglądarce otwartych jest kilka projektów, wybierz je z listy po lewej stronie.
