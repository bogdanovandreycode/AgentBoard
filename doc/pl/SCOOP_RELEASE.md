# Przygotowanie wydania dla Scoop

## Co jest już zautomatyzowane

`scripts/package-scoop.ps1 -Version X.Y.Z` wykonuje `npm ci`, kompilację frontendową, `go test ./...`, kompilację Windows `agentboard.exe` z numerem wersji, tworzy plik ZIP i oblicza SHA-256 tego ZIP. Z **tego samego** pliku tworzy `release/agentboard.json` z adresem URL, skrótem, licencją MIT, podkładką CLI, skrótem, `checkver` i `autoupdate`. ZIP i instalator zawierają plik `LICENSE`. Baza danych znajduje się poza katalogiem instalacyjnym, więc `persist` nie jest potrzebny w manifeście.

`.github/workflows/release.yml` na tagu `vX.Y.Z` uruchamia ten sam pakiet w Windows Runner, tworzy instalator Inno Setup i dołącza plik ZIP, instalator i manifest do wydania GitHub. Ręczne uruchomienie przepływu pracy tworzy jedynie artefakt do testowania, bez publikowania wersji.

## Polecenie wysłania dla konserwatora

1. Sprawdź, czy kod, dokumentacja, numer wersji i plik `LICENSE` są gotowe.
2. Uruchom `./scripts/package-scoop.ps1 -Version X.Y.Z` lokalnie. Sprawdź dane wyjściowe `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` i `agentboard version` po rozpakowaniu. Nie zmieniaj kodu ZIP po obliczeniu skrótu.
3. Utwórz i prześlij tag `vX.Y.Z`. GitHub Actions opublikuje wydanie w formacie ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` i manifestem. Sprawdź wszystkie trzy pliki na stronie Wydanie i plik ZIP SHA-256 z manifestu.
4. Na czystym komputerze z systemem Windows i aplikacją Scoop uruchom `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, następnie `agentboard version`, `agentboard init` i `agentboard open` w folderze testowym.
5. Aby uzyskać stałą aktualizację, umieść wygenerowany plik `agentboard.json` we własnym wiadrze Scoop lub udostępnij go w odpowiednim wiadrze publicznym. Sprawdź ponownie `scoop update agentboard` po wydaniu następnej wersji. Polecenie instalacji z adresu URL jest odpowiednie dla pierwszego znajomego, a wiadro jest wygodniejsze w przypadku aktualizacji.

Nie zastępuj ręcznie losowej wartości `hash`: Scoop sprawdza zawartość pobranego pliku ZIP.

W przypadku sprawdzania manifestu skup się na [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [tworzenie manifestu](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) i [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
