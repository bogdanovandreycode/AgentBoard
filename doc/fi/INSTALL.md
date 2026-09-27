# Asennus ja käynnistäminen Windowsissa

## Asennusohjelma (suositus)

Lataa `agentboard-VERSION-windows-amd64-setup.exe` sivulta [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Ohjattu asennustoiminto ehdottaa kansiota; oletus on `C:\AI\AgentBoard`. Se kopioi `agentboard.exe`:n ja dokumentaation, luo pikakuvakkeen ja lisää valitun kansion `PATH`-järjestelmään. Asennuksen jälkeen avaa uusi pääte, jotta komento `agentboard` tulee saataville.

Suorita projektikansiossasi:

```powershell
agentboard init
agentboard open
```

Liitäntä on sisäänrakennettu `agentboard.exe`; Erillistä Go- tai Node.js-asennusta ei tarvita. Tiedot tallennetaan `%AppData%\AgentBoard`:hen ja säilytetään, kun ohjelma päivitetään tai sen asennus poistetaan. Asennuksen poistaminen "Asennettujen sovellusten" kautta poistaa pikakuvakkeet ja merkinnän tiedostosta `PATH`.

## Scoopin asentaminen

Jos Scooppia ei ole vielä asennettu, avaa PowerShell tavallisena käyttäjänä ja noudata [virallisia Scoop](https://scoop.sh/) ohjeita. Jos yrityksesi tietokoneellasi on rajoituksia, ota yhteyttä järjestelmänvalvojaan. AgentBoard voidaan käynnistää myös ZIP-kortista ilman Scoopia.

Vaihtoehtoisesti voit asentaa AgentBoardin Scoopin kautta:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

GitHub-julkaisun luettelon saatavuus voidaan tarkistaa osoitteesta [sivulla Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). Kun luettelo on lisätty Scoopiin, ämpäri voidaan asentaa kauhan nimellä ja päivittää `scoop update agentboard`-komennolla.

## ZIP ilman Scooppia

Lataa `agentboard-VERSION-windows-amd64.zip` julkaisuista, pura esimerkiksi tiedostoon `C:\Tools\AgentBoard`. PowerShellissä projektikansiossa:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Jos kyseessä on AI-asiakas, määritä `agentboard.exe`:n koko polku sen MCP-kokoonpanossa, jos ohjelma ei ole `PATH`:ssa.

## Rakenna lähteestä

Asenna Go-versio `go.mod`:sta ja Node.js 22:sta tai uudemmasta. PowerShellissä arkiston juuressa:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1` kokoaa etuosan `internal/webui/dist`:ksi, suorittaa Go-testejä ja kokoaa yhden `agentboard.exe`:n. Järjestys on tärkeä: käyttöliittymä on sisäänrakennettu binaariin, kun Go rakennetaan. Jos vanha `agentboard.exe open` on käynnissä, pysäytä se ennen uudelleen rakentamista (Ctrl+C), muuten Windows ei salli sinun vaihtaa tiedostoa.

## Missä tiedot ovat

- `%AppData%\AgentBoard\agentboard.db` - tehtävät, projektit, työntekijät ja asetukset. Voit määrittää eri tiedoston `--db`-lipulla, mutta `init`:lla, `open`/`serve`:lla ja `mcp`:lla on oltava **sama polku**.
- `<ваш проект>\.agentboard\project.json` — projektin tunniste. Tämä tiedosto ei sisällä tehtäviä.
- Palvelin kuuntelee vain `127.0.0.1:7337` oletusarvoisesti. Määritä toinen osoite `--addr` ennen projektipolkua: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` avaa projektisivun ja käyttää uudelleen jo käynnissä olevaa palvelinta kyseisessä osoitteessa. Jos yhdessä selaimessa on auki useita projekteja, valitse ne vasemmalla olevasta luettelosta.
