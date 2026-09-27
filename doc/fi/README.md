# AgentBoard

[🌐 Languages](../LANGUAGES.md)

![Esikatselu AgentBoard](../../assets/social-preview.png)

**[Lataa Windowsille](https://github.com/bogdanovandreycode/AgentBoard/releases/latest) · [Dokumentaatiosivusto](https://bogdanovandreycode.github.io/AgentBoard/) · [Lisenssi MIT](../../LICENSE)**

AgentBoard on paikallinen tehtävälautakunta, jossa ihmiset ja tekoälytyöntekijät työskentelevät yhteisissä tehtävissä, mutta heillä on erilaiset oikeudet. Sovellus käynnistetään yhdellä tiedostolla `agentboard.exe`, avaa verkkoliittymän selaimessa ja tarjoaa työntekijöille erillisen MCP-palvelimen `stdio`:n kautta. Tiedot jäävät tietokoneellesi.

**[Aloittaa tyhjästä](START_HERE.md)· [Työtehtävien parissa työskenteleminen](TASKS.md)· [AI-yhteys MCP:n kautta](WORKERS_MCP.md)· [Tuo JSON](IMPORT.md)· [Asetukset](SETTINGS.md)· [Ongelmanratkaisu](TROUBLESHOOTING.md)**

## MVP-ominaisuudet

- Paikallinen projektilevy, jossa on tehtävähaku, JSON-tuonti, mukautetut sarakkeet ja ominaisuudet.
- MCP-käyttöoikeus tietylle työntekijälle rajoitetuilla tekoälysiirtymillä ja lopullisella ihmisen hyväksynnällä tehtävissä.
- Yleinen tehtävähistoria, tarkistusohjeet, artefaktit, tekoälykustannukset ja työntekijöiden yhteysdiagnostiikka.
- Windowsin asennusohjelma, ZIP- ja Scoop-luettelo; web-käyttöliittymä on sisäänrakennettu suoritettavaan tiedostoon.

AgentBoard on suunniteltu luotettavalle paikalliselle käyttäjälle. Sovellus ei isännöi projekteja pilvessä eikä käynnistä AI-asiakkaita itse; tarvittaessa yhdistä MCP-yhteensopiva asiakas työntekijään.

## Viidessä minuutissa

1. Lataa `agentboard-VERSION-windows-amd64-setup.exe`-asennusohjelma osoitteesta [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Se ehdottaa kansiota (oletuksena `C:\AI\AgentBoard`) ja lisää sen kansioon `PATH`. Scoop ja ZIP ovat myös saatavilla.
2. Avaa PowerShell projektikansiossasi, esimerkiksi `C:\Projects\MyApp`.
3. Suorita `agentboard init` (ZIP:lle: koko polku `agentboard.exe`:iin ja `init`).
4. Suorita `agentboard open`. `http://127.0.0.1:7337` avautuu.
5. Lisää tehtävä **Uusi tehtävä** -painikkeella. Jos kyseessä on tekoälytyöntekijä, avaa **Työntekijät → Lisää työntekijä**, valitse asiakasprofiili ja kopioi MCP-määritykset.

Jos sinulla ei vielä ole projektikansiota, luo se Windowsin Resurssienhallinnassa. Projekti voi olla mikä tahansa kansio, jopa ilman Gitiä ja koodia.

## Asennus Scoopin kautta

PowerShellissä, jossa [Scoop](https://scoop.sh/)] on jo asennettu julkaisun jälkeen:

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

## Lisenssi

AgentBoard on ilmainen ja avoimen lähdekoodin projekti [lisenssillä MIT](../../LICENSE). Kaupallinen käyttö, muokkaaminen, jakaminen ja uudelleenjakelu on sallittua edellyttäen, että tekijänoikeusilmoitus ja lisenssi säilyvät.
