# Ustawienia

Otwórz **Ustawienia** w bocznym menu wybranego projektu. Zmiany zapisywane są przyciskiem **Zapisz ustawienia** i zapisywane dla projektu w bazie AgentBoard.

## Język i wygląd

Pole **Język** otwiera listę rozwijaną wyszukiwania. Domyślnie wybrana jest opcja **Podążaj za systemem**, która uwzględnia język przeglądarki. Dostępne języki ze zrzutów ekranu: arabski, portugalski (Brazylia), chiński uproszczony, czeski, duński, holenderski, angielski, fiński, francuski, niemiecki, włoski, japoński, koreański, norweski Bokmål, polski, rosyjski, hiszpański, szwedzki, turecki, ukraiński, wietnamski; dodatkowo białoruski, rumuński i bułgarski. **Śledź system** przyjmuje język przeglądarki. Główne podpisy są tłumaczone ręcznie, pozostałe linie interfejsu użytkownika mają wstępne tłumaczenie automatyczne. Przed publikacją zaleca się sprawdzenie tłumaczeń przez native speakerów; nazwy klientów, polecenia, pola JSON i dane użytkownika pozostają bez tłumaczenia.

**Schemat kolorów**: Ciemny (oryginalny motyw), Jasny, Czarny, Ubuntu i Windows. **Strefa czasowa** steruje wyświetlaniem dat; dane są nadal przechowywane w formacie UTC. **Systemowa strefa czasowa** wykorzystuje ustawienia komputera.

## Kolumny

Cztery etapy `Features`, `In progress`, `Testing`, `Verification` są ustalone w tej kolejności: nie można zmienić ich nazwy ani usunąć. Pozostałe kolumny można uporządkować za pomocą dostępnych strzałek. Wpisz nazwę i kliknij **Dodaj kolumnę**, aby utworzyć kolumnę niestandardową. Jest przeznaczony do zadań ludzkich: AI widzi zadanie takie jak `Backlog` i nie otrzymuje go przez MCP. Usunięcie kolumny podczas zapisywania powoduje przeniesienie jej zadań do zwykłego `Backlog`.

## Sieć i MCP

**Odświeżanie tablicy** ustawia częstotliwość odświeżania planszy i kart w sekundach (1–60). **Odświeżenie pracowników** aktualizuje status pracowników (2–120 sekund). To jest odpytywanie interfejsu internetowego, a nie częstotliwość wyzwalania AI. Pracownicy nie rozpoczynają pracy automatycznie.

Adres serwera WWW jest ustawiany podczas uruchamiania CLI, na przykład `agentboard open --addr 127.0.0.1:7444`. Aby zmienić adres, należy zrestartować serwer. Wartość domyślna to `127.0.0.1:7337`. MCP działa poprzez oddzielne polecenie lokalne `agentboard mcp --project ... --worker ...` i jest niezależne od portu internetowego. W przypadku bazy niestandardowej we wszystkich poleceniach należy podać ten sam `--db`. Skopiuj konfigurację każdego klienta z karty roboczej.
