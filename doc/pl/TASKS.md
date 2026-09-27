# Praca z zadaniami

## Dodaj zadanie

W zakładce **Tablica** kliknij **Nowe zadanie**. Wypełnij tytuł i opis. Opis i instrukcje testowania obsługują Markdown: zaznacz tekst i użyj paska formatowania dla pogrubienia, kursywy, łącza i listy. Kliknij **Zapisz**.

Pola:

| Pole | Co robi |
| --- | --- |
| Tytuł | Krótka nazwa zadania. |
| Opis | Co należy zrobić i jak zrozumieć, że praca jest gotowa. |
| stan | Etap pracy; nowe zadanie zwykle zaczyna się od `Backlog`. |
| Priorytet | `critical`, `high`, `medium` lub `low`. |
| Odpowiedzialny | Osoba, konkretny pracownik lub bez umówionego spotkania. |
| Tryb testowy | `AI` - sprawdza AI; `Human` - kontrole ludzkie; `Hybrid` - oba. |
| Zależności | Zadania, które należy wykonać wcześniej. |
| Instrukcje dotyczące testu na sztuczną inteligencję/człowieka | Instrukcje dla odpowiedniego recenzenta. |
| Właściwości niestandardowe | Dodatkowe pola utworzone w zakładce **Właściwości**. |

W przypadku zadania AI najpierw utwórz pracownika, wybierz go w obszarze Odpowiedzialny, a następnie przenieś kartę do `Features`. AI widzi przydzielone jej zadania tylko w czterech etapach od `Features` do `Verification`.

## Przenieś zadanie

Przeciągnij kartę pomiędzy kolumnami. Możesz przełączać się między widokiem obejmującym wszystkie kolumny i szerokimi kolumnami z przewijaniem w poziomie. `Backlog` i `Complete` są kontrolowane przez człowieka. Sztuczna inteligencja może przyspieszyć zadanie `Features → In progress → Testing → Verification` jedynie za pomocą specjalnych narzędzi MCP. Zadanie z niezakończonymi testami ręcznymi nie powinno przejść weryfikacji AI.

## Karta zadań

Kliknij kartę, aby wyświetlić opis, właściciela, instrukcje testów oraz zakładki z historią, testami, artefaktami i kosztami sztucznej inteligencji. **Edytuj zadanie** zmienia treść. Komentarz danej osoby jest dodawany do ogólnej historii. Szukaj w najlepszych filtrach, wyszukuje według tytułu, opisu, identyfikatora i pracownika; osobny przycisk otwiera duże wyszukiwanie.

## Kolumny i właściwości

W zakładce **Ustawienia** możesz zmienić kolejność kolumn dostępnych dla danej osoby oraz dodać własną. Cztery etapy AI są stałe i przebiegają w tej samej kolejności. Kolumna użytkownika to miejsce na zadania odłożone przez osobę: dla AI takie zadanie ma status `Backlog`. Po usunięciu kolumny jej zadania wracają do normalnego stanu `Backlog`.

W zakładce **Właściwości** możesz dodać pola takie jak tekst, liczba, flaga, data, wybór i adres URL. `Human only` ukrywa pole przed AI; `Agent read` umożliwia odczyt, `Agent read/write` umożliwia także zapis za pomocą obsługiwanych narzędzi. Nie zmienia to uprawnień AI do poszczególnych etapów zadania.
