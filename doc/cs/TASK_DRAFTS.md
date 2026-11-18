# Řízení projektů a návrhy úkolů

Stránka **Všechny projekty** zobrazuje registrovanou složku a počet publikovaných
úkoly a poslední aktivita. Aktivita zohledňuje změny v úkolech, historii a
záznamy o výdajích. Tlačítko tužky umožňuje přejmenovat projekt nebo vybrat
existující složku v dialogovém okně systému. Změna cesty aktualizuje registraci a značku;
ID projektů, pracovníci, úkoly a nastavení se uloží, soubory se nepřesouvají.
Složky, které jsou již zaregistrovány do jiného projektu, nelze vybrat.

## Globální nastavení

Na stránce výběru projektu nebo na kartě otevřete **Globální nastavení**
**Globální nastavení** v nastavení libovolné desky. Tyto parametry jsou společné pro všechny projekty
jedna základna AgentBoard:

- Automatické ukládání konceptů pomocí časovače je ve výchozím nastavení povoleno každých **10 sekund**.
- **Ukládání při opuštění pole** je ve výchozím nastavení povoleno a lze jej konfigurovat samostatně.
- **Kontext projektu MCP** je povolen a zakázán globálně. Jejich oblast
  akce a samostatný příznak CLI jsou popsány v části [Workers and MCP](WORKERS_MCP.md).

## Návrh životního cyklu

Dialog Nový úkol uloží nedokončenou práci jako koncept. Stavový řádek
rozlišuje mezi neuloženými změnami, probíhajícím ukládáním, uloženým konceptem a chybou
konzervace. Koncept se neobjeví na nástěnce a není k dispozici pracovníkům MCP do
zveřejněno.

Před publikováním pomocí tlačítka **Vytvořit úkol** vyplňte **Název, Stát,
Odpovědnost a popis**. Předtím **Uložit koncept** uloží koncept bez vytvoření
úkoly. Prázdný formulář je při zavření ve výchozím nastavení zahozen. Pokud je vymazáno
existující koncept do prázdného stavu, bude při příštím uložení odstraněn.
Zavřením neprázdného dialogu se také uloží.

Tlačítko **Drafts** na nástěnce umožňuje pokračovat nebo smazat uložené koncepty.
Podporováno více návrhů. Po restartu se otevře
návrh. Koncepty jsou uloženy v SQLite; Při každé změně se aktualizuje ta místní
zkopírujte v prohlížeči, takže neuložené úpravy přežijí restarty resp
dočasné selhání serveru. Pokud není úložiště prohlížeče k dispozici, zobrazí se varování
že koncept musí být před uzavřením uložen na serveru.

Psaní konceptů používá revize, aby se zabránilo přepsání práce v jednom okně
další. Odezva na ztracené uložení je konzistentní, pokud již nějaká na serveru existuje
identické údaje. Pokud je obsah odlišný, zobrazí se konflikt a
**Obnovit jako nový koncept** ukládá místní obsah samostatně. Publikace
atomický a idempotentní: znovu publikujte stejný koncept poté, co je ztracen
odpověď vrátí stejný úkol, nikoli vytvoření kopie.
