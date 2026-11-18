# Ongelmanratkaisu

| Oire | Mitä tarkistaa |
| --- | --- |
| `agentboard` ei löytynyt | Käynnistä PowerShell uudelleen Scoopin jälkeen. Kun asennat ZIP-tiedostosta, käytä koko polkua `agentboard.exe`. |
| Web-sivu, joka näyttää vanhan käyttöliittymän rakentamisen jälkeen | Pysäytä käynnissä oleva palvelin Ctrl+C. Suorita `./scripts/build.ps1`, käynnistä uusi binaari. Viten on rakennettava ennen Goa, koska käyttöliittymä on sisäänrakennettu EXE:hen. Päivitä sivu Ctrl+F5. |
| Portti 7337 varattu | AgentBoard saattaa olla jo käynnissä. Avaa `http://127.0.0.1:7337` tai lopeta vanha prosessi. Käytä muita portteja varten `--addr`. |
| Projektia ei löydy | Suorita halutussa kansiossa `agentboard init`. Sitten `agentboard open` siitä tai `agentboard open C:\путь\к\проекту`. |
| Työntekijä ei näe tehtävää | Tehtävä on osoitettava tälle tietylle työntekijälle, ja sen on sijaittava paikassa `Features`, `In progress`, `Testing` tai `Verification`. Tekoäly ei näe `Backlog`-, `Complete`- ja mukautettuja sarakkeita. |
| `NO_PROJECTS` | MCP-tietokannassa ei ole projekteja. Suorita `agentboard init` projektikansiossa ja luo työntekijä; agentti voi toistaa puhelun käynnistämättä uudelleen. Tarkista, että MCP ja verkkokäyttöliittymä käyttävät samaa `--db`:ta. |
| `PROJECT_REQUIRED` / `PROJECT_AMBIGUOUS` | Ilman `--project`:ta asiakkaan työkansio tai MCP-käynnistyskansio eivät ole täsmälleen samassa rekisteröidyssä projektissa. Suorita agentti projektikansiossa, rekisteröi se `agentboard init`:n kautta tai lisää `--project`. Verkkokäyttöliittymää ei tarvitse käynnistää. |
| `PROJECT_NOT_FOUND` / `WORKER_NOT_FOUND` / `WORKER_REQUIRED` / `NO_WORKERS` | Katso `--project`, tämän projektin työntekijäetana ja `--db`. Viestissä luetellaan rekisteröidyt projektit tai käytössä olevat työntekijät. Mukautettu MCP-kokoonpano ilman `--project`:ta tarvitsee työntekijän, jolla on sama slug jokaiseen projektiin. |
| `PROJECT_CONTEXT_MISMATCH` | MCP-käynnistyskansio tai asiakkaan työkansio ei vastaa määritettyä `--project`:ta. Vahvista kysymys, jos tämä on tahallista, muussa tapauksessa poista `--project` (automaattinen tunnistus) tai ilmoita oikea. |
| MCP-palvelimen tarkistus on olemassa, mutta asiakas ei ole yhteydessä | Palvelimen tarkistus ei tarkista ulkoisen asiakkaan asetuksia. Käynnistä asiakas uudelleen, tarkista sen määritystiedosto, polku kohteisiin `agentboard.exe`, `--project`, `--worker` ja yleinen `--db`. Pyydä soittamaan `get_my_board`. |
| Työntekijä offline-tilassa | Asiakas on saattanut lopettaa tai ei ole vielä käynnistänyt MCP:tä. Kun 90 sekuntia ilman sykettä, istunto katsotaan katkenneeksi. |
| JSON-tiedostoa ei tuoda | Tarkista `version: 1`, vaaditut `title`, nykyiset `Slug` työntekijät ja kiinteistöjen nimet. JSON ei salli kommentteja tai pilkkuja. |
| Ei voida rakentaa uudelleen `agentboard.exe` | Windows ei voi korvata käynnissä olevaa EXE:tä. Pysäytä palvelin Ctrl+C ja yritä rakentaa uudelleen. |
| Tehtävät katosivat päivityksen jälkeen | Tarkista, että `--db` ei osoita toiseen tiedostoon ja että olet kirjautunut sisään samana Windows-käyttäjänä. Oletustukiasema sijaitsee `%AppData%\AgentBoard`:ssa. |

Jos virhettä ei ole kuvattu, kerää viestin tarkka teksti, `agentboard version`- ja Windows-versiot sekä vaiheet yrittääksesi uudelleen. Älä julkaise yksityisiä projektitietoja tai tietokannan sisältöä avoimessa numerossa.
