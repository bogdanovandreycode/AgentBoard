# Valmistellaan julkaisua Scoopille

## Mikä on jo automatisoitu

`scripts/package-scoop.ps1 -Version 0.2.0` suorittaa `npm ci`, frontend build, `go test ./...`, Windows build `agentboard.exe` versionumerolla, luo ZIP-tiedoston ja laskee tämän ZIP:n SHA-256:n. **samasta** tiedostosta se luo `release/agentboard.json`:n URL-osoitteella, hashilla, CLI-välilevyllä, pikakuvakkeella, `checkver` ja `autoupdate`. Tietokanta sijaitsee asennushakemiston ulkopuolella, joten `persist`:ta ei tarvita luettelossa.

`.github/workflows/release.yml` tagissa `vX.Y.Z` käyttää samaa pakkausta Windows Runnerissa, rakentaa Inno Setup -asennusohjelman ja liittää ZIP-tiedoston, asennusohjelman ja luettelon GitHub-julkaisuun. Työnkulun manuaalinen suorittaminen luo vain artefaktin testausta varten julkaisematta julkaisua.

## Kirjausjärjestys ylläpitäjälle

1. Tarkista, että koodi, dokumentaatio ja versionumero ovat valmiit. Määritä projektin lisenssi: manifesti osoittaa tällä hetkellä `Unknown`, koska lisenssitiedostoa ei ole määritetty arkistossa. Jos valitset lisenssin, lisää `LICENSE` ja päivitä `license` skriptiin ennen julkaisua.
2. Suorita `./scripts/package-scoop.ps1 -Version X.Y.Z` paikallisesti. Tarkista `release/agentboard-X.Y.Z-windows-amd64.zip`-, `release/agentboard.json`- ja `agentboard version`-ulostulot pakkauksesta purkamisen jälkeen. Älä muuta ZIP-koodia sen jälkeen, kun tiiviste on laskettu.
3. Luo ja lähetä tunniste `vX.Y.Z`. GitHub Actions julkaisee julkaisun ZIP-, `agentboard-X.Y.Z-windows-amd64-setup.exe`- ja manifestin kanssa. Tarkista kaikki kolme tiedostoa Julkaisusivulta ja SHA-256 ZIP luettelosta.
4. Suorita puhtaalla Windows-koneella Scoopilla `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, sitten `agentboard version`, `agentboard init` ja `agentboard open` testikansiossa.
5. Pysyvää päivityssyötettä varten aseta luotu `agentboard.json` omaan Scoop-ämpäriisi tai tarjoa se sopivaan julkiseen ämpäriin. Tarkista `scoop update agentboard` seuraavan julkaisun jälkeen. Asennuskomento URL-osoitteesta sopii ensituttamiseen, kun taas ämpäri on kätevämpi päivityksiin.

Älä korvaa manuaalisesti satunnaista arvoa `hash`: Scoop tarkistaa ladatun ZIP-tiedoston sisällön.

Luettelotarkistuksia varten keskity [muotoon Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [luontoluettelo](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) ja [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
