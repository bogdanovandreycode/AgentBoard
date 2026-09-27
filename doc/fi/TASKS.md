# Työskentely tehtävien kanssa

## Lisää tehtävä

Napsauta **Palvelu**-välilehdellä **Uusi tehtävä**. Täytä otsikko ja kuvaus. Kuvaus ja testausohjeet tukevat Markdownia: korosta tekstiä ja käytä muotoilupalkkia lihavointiin, kursivointiin, linkkiin ja luetteloihin. Napsauta **Tallenna**.

Kentät:

| Kenttä | Mitä tekee |
| --- | --- |
| Otsikko | Tehtävän lyhyt nimi. |
| Kuvaus | Mitä pitää tehdä ja miten ymmärtää, että työ on valmis. |
| valtio | Työ vaiheessa; uusi tehtävä alkaa yleensä `Backlog`:sta. |
| Prioriteetti | `critical`, `high`, `medium` tai `low`. |
| Vastuullinen | Henkilö, tietty työntekijä tai ilman ajanvarausta. |
| Testaustila | `AI` - tarkistaa AI; `Human` - ihmistarkastukset; `Hybrid` - molemmat. |
| Riippuvuudet | Tehtävät, jotka tulisi suorittaa aikaisemmin. |
| Tekoäly/ihmistestiohjeet | Ohjeet asianmukaiselle arvioijalle. |
| Mukautetut ominaisuudet | **Ominaisuudet**-välilehdelle luodut lisäkentät. |

AI-tehtävää varten luo ensin työntekijä, valitse se Vastuussa ja siirrä sitten kortti `Features`:hen. Tekoäly näkee sille osoitetut tehtävät vain neljässä vaiheessa `Features` - `Verification`.

## Siirrä tehtävä

Vedä kortti sarakkeiden väliin. Voit vaihtaa kaikki sarakkeet sisältävän näkymän ja leveiden sarakkeiden välillä vaakasuuntaisella vierityksellä. `Backlog` ja `Complete` ovat ihmisen ohjaamia. Tekoäly voi edistää tehtävää `Features → In progress → Testing → Verification` vain erityisten MCP-työkalujen avulla. Tehtävän, jonka manuaalinen testaus on kesken, ei pitäisi läpäistä tekoälyn varmennusta.

## Tehtäväkortti

Napsauta korttia nähdäksesi kuvauksen, omistajan, testiohjeet ja välilehdet historiasta, testeistä, esineistä ja tekoälykuluista. **Muokkaa tehtävää** muuttaa sisältöä. Henkilön kommentti lisätään koko tarinaan. Haku ylhäältä suodattaa haut otsikon, kuvauksen, tunnuksen ja työntekijän mukaan. erillinen painike avaa laajan haun.

## Sarakkeet ja ominaisuudet

**Asetukset**-välilehdellä voit muuttaa henkilön käytettävissä olevien sarakkeiden järjestystä ja lisätä omasi. Neljä AI-vaihetta ovat kiinteät ja etenevät samassa järjestyksessä. Käyttäjäsarake on paikka henkilön siirtämille tehtäville: tekoälyssä tällaisen tehtävän tila on `Backlog`. Kun sarake poistetaan, sen tehtävät palautuvat normaaliksi `Backlog`.

**Ominaisuudet**-välilehdelle voit lisätä kenttiä, kuten tekstiä, numeroa, lippua, päivämäärää, valintaa ja URL-osoitetta. `Human only` piilottaa kentän tekoälyltä; `Agent read` mahdollistaa lukemisen, `Agent read/write` mahdollistaa myös kirjoittamisen tuettujen työkalujen kautta. Tämä ei muuta tekoälyn oikeuksia tehtävän vaiheisiin.
