# Zarządzanie projektami i projekty zadań

Strona **Wszystkie projekty** pokazuje zarejestrowany folder i liczbę opublikowanych
zadania i ostatnia aktywność. Aktywność uwzględnia zmiany w zadaniach, historii i
ewidencja wydatków. Przycisk ołówka pozwala zmienić nazwę projektu lub wybrać
istniejący folder w oknie systemowym. Zmiana ścieżki aktualizuje rejestrację i znacznik;
Identyfikatory projektów, pracownicy, zadania i ustawienia są zapisywane, pliki nie są przenoszone.
Nie można wybrać folderów, które są już zarejestrowane w innym projekcie.

## Ustawienia globalne

Otwórz **Ustawienia globalne** na stronie lub karcie wyboru projektu
**Ustawienia globalne** w ustawieniach dowolnej tablicy. Parametry te są wspólne dla wszystkich projektów
jedna baza AgentBoard:

- Automatyczne zapisywanie wersji roboczych według timera jest domyślnie włączone co **10 sekund**.
- **Zapis przy wyjściu z pola** jest domyślnie włączony i można go osobno skonfigurować.
- Kontrole **kontekstu projektu MCP** są włączone i wyłączone globalnie. Ich obszar
  akcje i osobną flagę CLI opisano w sekcji [Workers i MCP](WORKERS_MCP.md).

## Projekt cyklu życia

Okno dialogowe Nowe zadanie zapisuje niedokończoną pracę jako wersję roboczą. Linia stanu
rozróżnia niezapisane zmiany, zapisywanie w toku, zapisaną wersję roboczą i błąd
konserwacja. Projekt nie pojawia się na tablicy i nie jest dostępny dla pracowników MCP do czasu
opublikowane.

Przed publikacją użyj przycisku **Utwórz zadanie**, aby wypełnić **Tytuł, Stan,
Odpowiedzialny i opis**. Wcześniej **Zapisz wersję roboczą** zapisuje wersję roboczą bez tworzenia
zadania. Pusty formularz jest domyślnie odrzucany po zamknięciu. Jeśli wyczyszczone
istniejącą wersję roboczą do stanu pustego, zostanie ona usunięta przy następnym zapisaniu.
Zamknięcie niepustego okna dialogowego również go zapisuje.

Przycisk **Wersje robocze** na tablicy umożliwia kontynuację lub usunięcie zapisanych wersji roboczych.
Obsługiwanych jest wiele wersji roboczych. Po ponownym uruchomieniu otwarte
projekt. Wersje robocze są przechowywane w SQLite; Przy każdej zmianie aktualizacja lokalna jest aktualizowana
skopiuj w przeglądarce, aby niezapisane zmiany przetrwały ponowne uruchomienie lub
tymczasowa awaria serwera. Jeśli pamięć przeglądarki nie jest dostępna, pojawi się ostrzeżenie
że wersja robocza musi zostać zapisana na serwerze przed zamknięciem.

Pisanie wersji roboczych wykorzystuje poprawki, aby zapobiec nadpisaniu pracy przez jedno okno
inny. Utracona odpowiedź na zapis jest spójna, jeśli taka już istnieje na serwerze
identyczne dane. Jeśli treść jest inna, zostanie wyświetlony konflikt i
**Odzyskaj jako nową wersję roboczą** oddzielnie zapisuje zawartość lokalną. Publikacja
atomowy i idempotentny: ponownie opublikuj tę samą wersję roboczą po jej utracie
odpowiedź zwraca to samo zadanie, zamiast tworzyć kopię.
