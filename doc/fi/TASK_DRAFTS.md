# Projektinhallinta ja tehtäväluonnokset

**Kaikki projektit** -sivulla näkyy rekisteröity kansio ja julkaistujen määrä
tehtäviä ja viimeisintä toimintaa. Toiminnassa otetaan huomioon muutokset tehtävissä, historiassa ja
kulukirjaukset. Kynäpainikkeella voit nimetä projektin uudelleen tai valita
olemassa oleva kansio järjestelmävalintaikkunassa. Polun muuttaminen päivittää rekisteröinnin ja merkin;
Projektitunnukset, työntekijät, tehtävät ja asetukset tallennetaan, tiedostoja ei siirretä.
Toiseen projektiin jo rekisteröityjä kansioita ei voi valita.

## Yleiset asetukset

Avaa **Yleiset asetukset** projektinvalintasivulta tai -välilehdeltä
**Yleiset asetukset** minkä tahansa levyn asetuksissa. Nämä parametrit ovat yhteisiä kaikille projekteille
yksi AgentBoard-kanta:

- Luonnosten automaattinen tallennus ajastimella on oletusarvoisesti käytössä **10 sekunnin välein**.
- **Tallennus kentältä poistuttaessa** on oletuksena käytössä ja se voidaan määrittää erikseen.
- **MCP-projektikonteksti** -tarkistukset ovat käytössä ja pois käytöstä maailmanlaajuisesti. Heidän alueensa
  toiminnot ja erillinen CLI-lippu on kuvattu osiossa [Workers and MCP](WORKERS_MCP.md).

## Luonnoksen elinkaari

Uusi tehtävä -valintaikkuna tallentaa keskeneräiset työt luonnoksena. Tilarivi
erottaa tallentamattomat muutokset, tallennuksen käynnissä, tallennetun luonnoksen ja virheen välillä
säilyttäminen. Luonnos ei näy taululla, eikä se ole MCP:n työntekijöiden saatavilla ennen kuin
julkaistu.

Täytä ennen julkaisua **Luo tehtävä** -painiketta **Otsikko, Osavaltio,
Vastuullinen ja kuvaus**. Ennen tätä **Tallenna luonnos** tallentaa luonnoksen luomatta
tehtäviä. Tyhjä lomake hylätään oletuksena, kun se suljetaan. Jos tyhjennetään
olemassa oleva luonnos tyhjäksi, se poistetaan, kun se seuraavan kerran tallennetaan.
Ei-tyhjän valintaikkunan sulkeminen myös tallentaa sen.

Taulukon **Luonnokset**-painikkeella voit jatkaa tai poistaa tallennettuja luonnoksia.
Useita luonnoksia tuetaan. Uudelleenkäynnistyksen jälkeen avataan
luonnos. Luonnokset tallennetaan SQLiteen; Jokaisen muutoksen yhteydessä paikallinen päivitetään
kopioi selaimeen, jotta tallentamattomat muokkaukset säilyvät uudelleen käynnistyksissä tai
väliaikainen palvelinvika. Jos selaimen tallennustilaa ei ole saatavilla, valintaikkuna varoittaa
että luonnos on tallennettava palvelimelle ennen sulkemista.

Luonnosten kirjoittaminen käyttää versioita estääkseen yhtä ikkunaa korvaamasta työtä
toinen. Kadonnut tallennusvastaus on johdonmukainen, jos sellainen on jo palvelimella
identtiset tiedot. Jos sisältö on erilainen, näytetään ristiriita ja
**Palauta uutena luonnoksena** tallentaa paikallisen sisällön erikseen. Julkaisu
atomic and idempotent: julkaise sama luonnos uudelleen, kun se on kadonnut
vastaus palauttaa saman tehtävän kopion luomisen sijaan.
