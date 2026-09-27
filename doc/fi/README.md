# AgentBoard

[🌐 Languages](../LANGUAGES.md)

AgentBoard on paikallinen tehtävälautakunta, jossa ihmiset ja tekoälytyöntekijät työskentelevät yhteisissä tehtävissä, mutta heillä on erilaiset oikeudet. Sovellus alkaa yhdellä tiedostolla`agentboard.exe`, avaa verkkokäyttöliittymän selaimessa ja tarjoaa työntekijöille erillisen MCP-palvelimen kautta`stdio`. Tiedot jäävät tietokoneellesi.

**[Aloittaa tyhjästä](START_HERE.md)· [Työtehtävien parissa työskenteleminen](TASKS.md)· [AI-yhteys MCP:n kautta](WORKERS_MCP.md)· [Tuo JSON](IMPORT.md)· [Asetukset](SETTINGS.md)· [Ongelmanratkaisu](TROUBLESHOOTING.md)**

## Viidessä minuutissa

1. Lataa asennusohjelma`agentboard-VERSION-windows-amd64-setup.exe`alkaen [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Se ehdottaa kansiota (oletus`C:\AI\AgentBoard`) ja lisää sen`PATH`. Scoop ja ZIP ovat myös saatavilla.
2. Avaa PowerShell projektikansiossasi kuten`C:\Projects\MyApp`.
3. Suorita`agentboard init`(ZIP: koko polku osoitteeseen`agentboard.exe`Ja`init`).
4. Suorita`agentboard open`. Avautuu`http://127.0.0.1:7337`.
5. Lisää tehtävä **Uusi tehtävä** -painikkeella. Jos kyseessä on tekoälytyöntekijä, avaa **Työntekijät → Lisää työntekijä**, valitse asiakasprofiili ja kopioi MCP-määritykset.

Jos sinulla ei vielä ole projektikansiota, luo se Windowsin Resurssienhallinnassa. Projekti voi olla mikä tahansa kansio, jopa ilman Gitiä ja koodia.

## Asennus Scoopin kautta

PowerShellissä, jossa [Scoop] on jo asennettu](https://scoop.sh/)julkaisun jälkeen:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Kehittäjille on [build lähteestä ](INSTALL.md). Julkaisutyönkulku luo asennusohjelman, ZIP- ja Scoop-luettelon SHA-256:lla samasta artefaktista. Scoopin kautta asennetun version päivittäminen: `scoop update agentboard` luettelon lisäämisen jälkeen ämpäriin; tiedot - [julkaisun valmistelu](SCOOP_RELEASE.md).

## Kuinka hallitus on rakennettu

`Backlog → Features → In progress → Testing → Verification → Complete`

Henkilö voi siirtää tehtäviä laudalla. AI voi siirtää vain `Features → In progress → Testing → Verification`; Lopullisen hyväksynnän `Complete`:hen suorittaa ihminen. Käyttäjäsarakkeet on tarkoitettu ihmisille: niissä oleva tehtävä pysyy MCP:n tilassa `Backlog`. Testitilat: AI, Human ja Hybrid. Historia, testit, esineet ja tekoälykustannukset liitetään tehtävään.

## Joukkueet

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` rekisteröi kansion ja kirjoittaa siihen vain `.agentboard/project.json`. SQLiten tuotantotiedot sijaitsevat Windowsin käyttäjän asetushakemistossa (`%AppData%\AgentBoard\agentboard.db`), projektin ulkopuolella ja Scoop-asennuksen ulkopuolella. Paketin poistaminen tai päivittäminen ei saa poistaa näitä tietoja. Ennen kuin siirrät toiseen tietokoneeseen, tee tietokannasta kopio AgentBoardin ollessa pysäytettynä.

Työntekijät - loogiset tilit; AgentBoard itse ei käytä Codexia, Claudea tai muita tekoälyasiakkaita. Asiakas aloittaa paikallisen MCP-prosessin tietylle työntekijälle. MCP on sovelluksen käyttöoikeusraja, ja tiedostojen eristämiseen käytä AI-asiakkaan hiekkalaatikkoa.

## Kehittäjille

Pino: Go, SQLite, virallinen MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Rakenna ensin käyttöliittymä ja sitten Go: `./scripts/build.ps1`. Web-tiedostot sisältyvät binaariin `go:embed`:n kautta. Arkkitehtuuri ja API on kuvattu julkaisussa [doc/ARCHITECTURE.md](ARCHITECTURE.md).
