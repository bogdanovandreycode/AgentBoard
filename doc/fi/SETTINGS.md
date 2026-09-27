# Asetukset

Avaa **Asetukset** valitun projektin sivuvalikosta. Muutokset tallennetaan **Tallenna asetukset** -painikkeella ja tallennetaan projektia varten AgentBoard-tietokantaan.

## Kieli ja ulkonäkö

**Kieli**-kenttä avaa avattavan hakuluettelon. Oletuksena on valittuna **Seuraa järjestelmää** -vaihtoehto, joka ottaa selaimen kielen. Kuvakaappauksista saatavilla olevat kielet: arabia, portugali (Brasilia), yksinkertaistettu kiina, tšekki, tanska, hollanti, englanti, suomi, ranska, saksa, italia, japani, korea, norja bokmål, puola, venäjä, espanja, ruotsi, turkki, ukraina, vietnam; lisäksi valkovenäläinen, romania ja bulgaria. **Seuraa järjestelmää** käyttää selaimen kieltä. Tärkeimmät allekirjoitukset käännetään manuaalisesti, muilla käyttöliittymäriveillä on alustava automaattinen käännös. Ennen julkistamista on suositeltavaa oikolukea äidinkielenään puhuvien käännökset; asiakasnimet, komennot, JSON-kentät ja käyttäjätiedot jäävät ilman käännöstä.

**Värimalli**: Tumma (alkuperäinen teema), vaalea, musta, Ubuntu ja Windows. **Aikavyöhyke** ohjaa päivämäärien näyttöä; tiedot säilyvät edelleen UTC:ssä. **Järjestelmän aikavyöhyke** käyttää tietokoneen asetuksia.

## Sarakkeet

Neljä vaihetta `Features`, `In progress`, `Testing`, `Verification` on kiinteät tässä järjestyksessä: niitä ei voi nimetä uudelleen tai poistaa. Muut sarakkeet voidaan järjestää uudelleen käyttämällä käytettävissä olevia nuolia. Kirjoita nimi ja napsauta **Lisää sarake** luodaksesi mukautetun sarakkeen. Se on suunniteltu ihmistehtäviin: tekoäly näkee tehtävän, kuten `Backlog`, eikä vastaanota sitä MCP:n kautta. Sarakkeen poistaminen tallennuksen aikana siirtää sen tehtävät tavalliseen `Backlog`:hen.

## Web ja MCP

**Palvelun virkistys** määrittää laudan ja korttien virkistystaajuuden sekunneissa (1–60). **Työntekijän päivitys** päivittää työntekijöiden tilan (2–120 sekuntia). Tämä on verkkokäyttöliittymän kyselyä, ei AI-laukaisutaajuutta. Työntekijät eivät käynnisty automaattisesti.

Verkkopalvelimen osoite asetetaan CLI:tä käynnistettäessä, esimerkiksi `agentboard open --addr 127.0.0.1:7444`. Jos haluat muuttaa osoitetta, palvelin on käynnistettävä uudelleen. Oletus on `127.0.0.1:7337`. MCP toimii erillisen paikallisen komennon `agentboard mcp --project ... --worker ...` kautta ja on riippumaton verkkoportista. Jos kyseessä on ei-standardi kanta, määritä sama `--db` kaikissa komennoissa. Kopioi kunkin asiakkaan asetukset työntekijäkortilta.
