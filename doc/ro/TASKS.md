# Lucrul cu sarcini

## Adăugați sarcină

În fila **Tablou**, dați clic pe **Sarcina nouă**. Completați titlul și descrierea. Descrierea și instrucțiunile de testare acceptă Markdown: evidențiați textul și utilizați bara de formatare pentru aldine, cursive, link și listă. Faceți clic pe **Salvați**.

Câmpuri:

| Câmp | Ce înseamnă |
| --- | --- |
| Titlu | Numele scurt al sarcinii. |
| Descriere | Ce trebuie făcut și cum să înțelegeți că lucrarea este gata. |
| Stat | Etapa de lucru; o sarcină nouă începe de obicei la `Backlog`. |
| Prioritate | `critical`, `high`, `medium` sau `low`. |
| Responsabil | O persoană, un anumit lucrător sau fără programare. |
| Modul de testare | `AI` - verifică AI; `Human` - verificări umane; `Hybrid` - ambele. |
| Dependente | Sarcini care ar trebui finalizate mai devreme. |
| Instrucțiuni de testare AI/Human | Instrucțiuni pentru examinatorul corespunzător. |
| Proprietăți personalizate | Câmpuri suplimentare create în fila **Proprietăți**. |

Pentru o sarcină AI, creați mai întâi un lucrător, selectați-l în Responsabil, apoi transferați cardul pe `Features`. AI vede doar sarcinile care îi sunt atribuite în patru etape de la `Features` la `Verification`.

## Mutați sarcina

Trageți cardul între coloane. Puteți comuta între o vizualizare cu toate coloanele și coloanele largi cu defilare orizontală. `Backlog` și `Complete` sunt controlate de om. AI poate avansa sarcina `Features → In progress → Testing → Verification` numai prin instrumente speciale MCP. O sarcină cu testare manuală nefinalizată nu trebuie să treacă de verificarea AI.

## Card de activitate

Faceți clic pe un card pentru a vedea descrierea, proprietarul, instrucțiunile de testare și file pentru istorie, teste, artefacte și costuri AI. **Editați sarcina** modifică conținutul. Comentariul persoanei este adăugat la povestea generală. Caută în filtrele de top căutări după titlu, descriere, ID și lucrător; un buton separat deschide o căutare mare.

## Coloane și proprietăți

În fila **Setări** puteți modifica ordinea coloanelor disponibile pentru o persoană și puteți să le adăugați pe ale dvs. Cele patru etape AI sunt fixe și procedează în aceeași ordine. Coloana utilizator este un loc pentru sarcini amânate de o persoană: pentru AI, o astfel de sarcină are statutul `Backlog`. Când o coloană este ștearsă, sarcinile acesteia revin la normal `Backlog`.

În fila **Proprietăți** puteți adăuga câmpuri precum text, număr, steag, dată, selecție și URL. `Human only` ascunde câmpul de AI; `Agent read` permite citirea, `Agent read/write` permite și scrierea prin instrumente acceptate. Acest lucru nu schimbă drepturile AI la pașii sarcinii.
