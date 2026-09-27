# Práce s úkoly

## Přidat úkol

Na kartě **Nástěnka** klikněte na **Nový úkol**. Vyplňte název a popis. Popis a pokyny k testování podporují Markdown: zvýrazněte text a použijte formátovací lištu pro tučné písmo, kurzívu, odkaz a seznam. Klikněte na **Uložit**.

Pole:

| Pole | Co dělá |
| --- | --- |
| Název | Krátký název úkolu. |
| Popis | Co je třeba udělat a jak pochopit, že práce je připravena. |
| stát | Pracovní fáze; nový úkol obvykle začíná na `Backlog`. |
| Priorita | `critical`, `high`, `medium` nebo `low`. |
| Zodpovědný | Člověk, konkrétní pracovník nebo bez objednání. |
| Testovací režim | `AI` - kontroluje AI; `Human` - lidské kontroly; `Hybrid` - oba. |
| Závislosti | Úkoly, které by měly být dokončeny dříve. |
| Pokyny pro AI/Human test | Pokyny pro příslušného recenzenta. |
| Vlastní vlastnosti | Další pole vytvořená na kartě **Vlastnosti**. |

Pro úlohu AI nejprve vytvořte pracovníka, vyberte jej v Zodpovědné a poté přeneste kartu do `Features`. AI vidí pouze úkoly, které jsou jí přiřazeny ve čtyřech fázích od `Features` do `Verification`.

## Přesunout úkol

Přetáhněte kartu mezi sloupce. Vodorovným posouváním můžete přepínat mezi zobrazením všech sloupců a širokými sloupci. `Backlog` a `Complete` jsou ovládány lidmi. Umělá inteligence může pokročit v úkolu `Features → In progress → Testing → Verification` pouze prostřednictvím speciálních nástrojů MCP. Úloha s nedokončeným ručním testováním by neměla projít ověřením AI.

## Karta úkolu

Kliknutím na kartu zobrazíte popis, vlastníka, pokyny k testování a karty pro historii, testy, artefakty a náklady na umělou inteligenci. **Upravit úkol** změní obsah. Komentář osoby se přidá k celkovému příběhu. Hledat v horních filtrech vyhledává podle názvu, popisu, ID a pracovníka; samostatné tlačítko otevírá velké vyhledávání.

## Sloupce a vlastnosti

Na záložce **Nastavení** můžete změnit pořadí sloupců dostupných osobě a přidat vlastní. Čtyři stupně umělé inteligence jsou pevně dané a postupují ve stejném pořadí. Uživatelský sloupec je místem pro úkoly odložené osobou: pro AI má takový úkol stav `Backlog`. Po odstranění sloupce se jeho úlohy vrátí do normálního stavu `Backlog`.

Na kartě **Vlastnosti** můžete přidat pole jako text, číslo, příznak, datum, výběr a URL. `Human only` skryje pole před AI; `Agent read` umožňuje čtení, `Agent read/write` umožňuje i zápis prostřednictvím podporovaných nástrojů. To nemění práva AI na kroky úkolu.
