# Începeți de aici

## Ce este AgentBoard

Imaginați-vă o tablă obișnuită cu carduri de sarcini. Creați sarcini, desemnați pe cineva responsabil și supravegheați munca. Lucrătorii AI primesc doar sarcinile care le sunt atribuite prin MCP și raportează progresul. Tu decizi când sarcina este în sfârșit gata.

**Proiect** - un folder pe computer și o placă separată. **Sarcina** - un card cu o descriere, o persoană responsabilă și o etapă. **Lucrător** este numele logic al clientului AI. **MCP** este modul în care clientul AI se conectează la AgentBoard. **Scoop** este un manager de instalare de software pentru Windows.

Nu este necesar niciun cod. Aveți nevoie de Windows, un browser, PowerShell și, pentru ca AI să funcționeze, un client AI instalat cu suport pentru serverele MCP locale.

## Prima lansare

1. Instalați aplicația conform [instrucțiuni ](INSTALL.md).
2. Creați un folder de proiect în Explorer, de exemplu `C:\Projects\MyFirstProject`.
3. Deschideți acest folder în Explorer. Faceți clic în bara de adrese, tastați `powershell` și apăsați Enter.
4. În fereastra care se deschide, faceți:

```powershell
agentboard init
agentboard open
```

5. Se va deschide un browser cu adresa `http://127.0.0.1:7337`. Lăsați fereastra PowerShell deschisă în timp ce utilizați tabla albă. Închiderea ferestrei va opri serverul local, dar sarcinile vor rămâne.

Dacă comanda `agentboard` nu este găsită, închideți PowerShell și redeschideți-l după instalarea Scoop. Când instalați dintr-un fișier ZIP, utilizați calea completă către `agentboard.exe`.

## Prima sarcină

Faceți clic pe **Sarcina nouă**, completați **Titlu**, dacă este necesar **Descriere**, apoi **Salvați**. O nouă sarcină în `Backlog` este disponibilă oamenilor. Pentru ca AI să înceapă să funcționeze, atribuiți un lucrător și transferați sarcina către `Features`. Descrierea pas cu pas a câmpurilor - [TASKS.md](TASKS.md).

## Primul muncitor

Deschideți **Lucrători → Adăugați un lucrător**. Selectați clientul AI pe care îl utilizați (de exemplu Codex sau Claude Code), verificați numele și ID-ul scurt `Slug`, faceți clic pe **Salvați**. Deschideți lucrătorul creat: există o configurație MCP, un buton de copiere și o verificare a serverului. Copiați configurația către clientul AI în conformitate cu [WORKERS_MCP.md](WORKERS_MCP.md). Odată conectat, cereți clientului să sune la `get_my_board`.

## Ce să citești în continuare

- [Sarcini, teste, povești și coloane](TASKS.md)
- [Conectarea lucrătorilor și verificarea MCP](WORKERS_MCP.md)
- [Importați sarcini din JSON](IMPORT.md)
- [Limbă, temă, fus orar și actualizare](SETTINGS.md)
- [Probleme tipice](TROUBLESHOOTING.md)
