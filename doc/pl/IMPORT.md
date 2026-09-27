# Importuj zadania z JSON

Zakładka **Zadania importu** akceptuje plik JSON w kodowaniu UTF-8. Użyj go, jeśli przenosisz zadania z innej usługi. Kliknij **Pobierz szablon JSON**: Szablon zawiera bieżące nazwy właściwości niestandardowych i `Slug` dostępnego pracownika dla tego projektu.

## Krok po kroku

1. Przed importem utwórz niezbędnych pracowników (**Pracownicy**) i właściwości (**Właściwości**).
2. Pobierz szablon, otwórz go w edytorze tekstu, zastąp przykłady własnymi zadaniami i zapisz plik z rozszerzeniem `.json`.
3. Na karcie **Zadania importu** wybierz plik. Interfejs pokaże liczbę zadań.
4. Kliknij **Importuj zadania**. Serwer sprawdzi cały plik: w przypadku wykrycia błędu nie zostanie zapisane ani jedno zadanie z niego. Popraw komunikat o błędzie i spróbuj ponownie.

Minimalny plik:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Przykład z pracownikiem, własnością i zależnością:

```json
{
  "version": 1,
  "tasks": [
    {
      "key": "design",
      "title": "Prepare the design",
      "description": "## Goal\nPrepare the home page mockup.",
      "state": "features",
      "priority": "high",
      "testing_mode": "hybrid",
      "assignee": { "type": "worker", "worker": "codex" },
      "ai_test_instructions": "Check the build.",
      "human_test_instructions": "Review the page in a browser.",
      "properties": { "Department": "Design" }
    },
    {
      "title": "Approve the design",
      "depends_on": ["design"],
      "assignee": { "type": "human" }
    }
  ]
}
```

`version` powinno być `1`, tablica `tasks` - od 1 do 1000 zadań. Wymagany jest `title`. Prawidłowe etapy: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Priorytety: `critical`, `high`, `medium`, `low`; tryby testowe: `ai`, `human`, `hybrid`. Odpowiedzialny: `unassigned`, `human` lub `worker` z istniejącym `worker` (Slug lub ID). `properties` wykorzystuje nazwy lub identyfikatory już utworzonych właściwości. `key` jest unikalny w pliku; `depends_on` odnosi się do takich kluczy. Serwer sam tworzy identyfikatory zadań.

Ponowny import tego samego pliku spowoduje utworzenie nowych problemów, więc sprawdź tablicę przed ponownym kliknięciem. Po zaimportowaniu zadania są rejestrowane jako stworzone przez człowieka; W historii pojawi się wpis systemowy.
