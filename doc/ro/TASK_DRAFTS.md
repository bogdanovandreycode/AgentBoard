# Management de proiect și proiecte de sarcini

Pagina **Toate proiectele** arată folderul înregistrat, numărul de publicate
sarcini și ultima activitate. Activitatea ia în considerare schimbările în sarcini, istoric și
evidențele cheltuielilor. Butonul creion vă permite să redenumiți proiectul sau să selectați
folderul existent în dialogul de sistem. Schimbarea căii actualizează înregistrarea și marcatorul;
ID-urile proiectului, lucrătorii, sarcinile și setările sunt salvate, fișierele nu sunt mutate.
Folderele care sunt deja înregistrate la alt proiect nu pot fi selectate.

## Setări globale

Deschideți **Setări globale** din pagina sau fila de selecție a proiectelor
**Setări globale** în setările oricărei plăci. Acești parametri sunt comuni tuturor proiectelor
o bază AgentBoard:

- Salvarea automată a schițelor prin cronometru este activată în mod implicit, la fiecare **10 secunde**.
- **Salvarea la părăsirea câmpului** este activată implicit și poate fi configurată separat.
- Verificările **contextul proiectului MCP** sunt activate și dezactivate la nivel global. Zona lor
  acțiunile și un indicator CLI separat sunt descrise în secțiunea [Lucrători și MCP](WORKERS_MCP.md).

## Schiță de ciclu de viață

Fereastra de dialog Nouă sarcină salvează munca neterminată ca schiță. Linia de stare
face distincție între modificările nesalvate, salvarea în curs, schița salvată și eroarea
conservare. Proiectul nu apare pe tablă și nu este disponibil pentru lucrătorii MCP până când
publicat.

Înainte de a publica, utilizați butonul **Creați sarcină** pentru a completa **Titlu, stare,
Responsabil și Descriere**. Înainte de aceasta, **Salvare schiță** salvează o schiță fără a crea
sarcini. Un formular gol este eliminat în mod implicit când este închis. Dacă este șters
schița existentă într-o stare goală, este ștearsă data viitoare când este salvată.
Închiderea unui dialog care nu este gol îl salvează și.

Butonul **Schițe** de pe tablă vă permite să continuați sau să ștergeți schițele salvate.
Mai multe schițe acceptate. După o repornire, se deschide
draft. Schițele sunt stocate în SQLite; Cu fiecare modificare, cea locală este actualizată
copiați în browser, astfel încât editările nesalvate să supraviețuiască repornirilor sau
defecțiune temporară a serverului. Dacă stocarea în browser nu este disponibilă, un dialog avertizează
că schița trebuie salvată pe server înainte de închidere.

Scrierea schițelor folosește revizuiri pentru a împiedica o fereastră să suprascrie lucrările
alta. Răspunsul de salvare pierdut este consistent dacă există deja unul pe server
date identice. Dacă conținutul este diferit, este afișat un conflict și
**Recuperare ca o schiță nouă** salvează separat conținutul local. Publicare
atomic și idempotent: republicați același proiect după ce este pierdut
răspunsul returnează aceeași sarcină în loc să creeze o copie.
