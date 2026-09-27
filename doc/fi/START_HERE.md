# Aloita tästä

## Mikä on AgentBoard

Kuvittele tavallinen taulu, jossa on tehtäväkortit. Luot tehtäviä, annat vastuuhenkilön ja valvot työtä. Tekoälytyöntekijät saavat vain heille määrätyt tehtävät MCP:n kautta ja raportoivat edistymisestä. Sinä päätät, milloin tehtävä on vihdoin valmis.

**Projekti** - kansio tietokoneella ja erillinen kortti. **Tehtävä** - kortti, jossa on kuvaus, vastuuhenkilö ja vaihe. **Worker** on AI-asiakkaan looginen nimi. **MCP** on tapa, jolla AI-asiakas muodostaa yhteyden AgentBoardiin. **Scoop** on ohjelmiston asennushallinta Windowsille.

Koodia ei vaadita. Tarvitset Windowsin, selaimen, PowerShellin ja, jotta tekoäly toimii, asennetun tekoälyasiakkaan, joka tukee paikallisia MCP-palvelimia.

## Ensimmäinen käynnistys

1. Asenna sovellus [ohjeet](INSTALL.md).
2. Luo projektikansio Explorerissa, esimerkiksi `C:\Projects\MyFirstProject`.
3. Avaa tämä kansio Explorerissa. Napsauta osoitepalkkia, kirjoita `powershell` ja paina Enter.
4. Tee avautuvassa ikkunassa:

```powershell
agentboard init
agentboard open
```

5. Avautuu selain, jonka osoite on `http://127.0.0.1:7337`. Jätä PowerShell-ikkuna auki, kun käytät taulua. Ikkunan sulkeminen pysäyttää paikallisen palvelimen, mutta tehtävät säilyvät.

Jos komentoa `agentboard` ei löydy, sulje PowerShell ja avaa se uudelleen Scoopin asentamisen jälkeen. Kun asennat ZIP-tiedostosta, käytä koko polkua `agentboard.exe`.

## Ensimmäinen tehtävä

Napsauta **Uusi tehtävä**, täytä **Otsikko**, tarvittaessa **Kuvaus** ja sitten **Tallenna**. Uusi tehtävä `Backlog`:ssa on ihmisten käytettävissä. Jotta tekoäly alkaa toimia, määritä työntekijä ja siirrä tehtävä `Features`:lle. Vaiheittainen kuvaus kentistä - [TASKS.md](TASKS.md).

## Ensimmäinen työntekijä

Avaa **Työntekijät → Lisää työntekijä**. Valitse käyttämäsi AI-asiakas (esimerkiksi Codex tai Claude Code), tarkista nimi ja lyhyt tunnus `Slug`, napsauta **Tallenna**. Avaa luotu työntekijä: siellä on MCP-kokoonpano, kopiointipainike ja palvelimen tarkistus. Kopioi konfiguraatio AI-asiakasohjelmaan [WORKERS_MCP.md](WORKERS_MCP.md):n mukaan. Kun yhteys on muodostettu, pyydä asiakasta soittamaan `get_my_board`.

## Mitä lukea seuraavaksi

- [Tehtävät, testit, tarinat ja kolumnit](TASKS.md)
- [Työntekijöiden yhdistäminen ja MCP](WORKERS_MCP.md):n tarkistaminen
- [Tuo tehtäviä JSON](IMPORT.md):sta
- [Kieli, teema, aikavyöhyke ja päivitys](SETTINGS.md)
- [Tyypilliset ongelmat](TROUBLESHOOTING.md)
