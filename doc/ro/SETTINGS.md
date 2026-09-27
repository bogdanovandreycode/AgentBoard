# Setari

Deschideți **Setări** în meniul lateral al proiectului selectat. Modificările sunt salvate cu butonul **Salvare setări** și stocate pentru proiect în baza de date AgentBoard.

## Limbă și aspect

Câmpul **Limbă** deschide o listă derulantă de căutare. În mod implicit, este selectată opțiunea **Urmăriți sistemul**, care ia limba browserului. Limbi disponibile din capturi de ecran: arabă, portugheză (Brazilia), chineză simplificată, cehă, daneză, olandeză, engleză, finlandeză, franceză, germană, italiană, japoneză, coreeană, norvegiană bokmål, poloneză, rusă, spaniolă, suedeză, turcă, ucraineană, vietnameză; în plus, belarusă, română și bulgară. **Urmărește sistemul** folosește limba browserului. Semnăturile principale sunt traduse manual, liniile UI rămase au o traducere automată preliminară. Înainte de lansarea publică, este recomandabil să corectați traducerile de către vorbitori nativi; numele clientului, comenzile, câmpurile JSON și datele utilizatorului rămân fără traducere.

**Schema de culori**: Întunecat (temă originală), Lumină, Negru, Ubuntu și Windows. **Fusul orar** controlează afișarea datelor; datele continuă să fie stocate în UTC. **Fusul orar al sistemului** folosește setările computerului.

## Coloane

Cele patru etape `Features`, `In progress`, `Testing`, `Verification` sunt fixate în această ordine: nu pot fi redenumite sau șterse. Alte coloane pot fi rearanjate folosind săgețile disponibile. Introduceți un nume și faceți clic pe **Adăugați coloană** pentru a crea o coloană personalizată. Este conceput pentru sarcini umane: AI vede o sarcină precum `Backlog` și nu o primește prin MCP. Ștergerea unei coloane în timpul salvării transferă sarcinile acesteia la `Backlog` obișnuit.

## Web și MCP

**Reîmprospătarea tablei** setează rata de reîmprospătare a tablei și a cărților în secunde (1–60). **Reîmprospătarea lucrătorilor** actualizează starea lucrătorilor (2–120 de secunde). Acesta este sondajul interfeței web, nu frecvența de declanșare a AI. Muncitorii nu pornesc automat.

Adresa serverului web este setată la pornirea CLI, de exemplu `agentboard open --addr 127.0.0.1:7444`. Pentru a schimba adresa, serverul trebuie repornit. Valoarea implicită este `127.0.0.1:7337`. MCP operează printr-o comandă locală separată `agentboard mcp --project ... --worker ...` și este independent de portul web. Pentru o bază non-standard, specificați același `--db` în toate comenzile. Copiați configurația fiecărui client de pe cardul lucrătorului.
