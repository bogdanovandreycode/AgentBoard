# Valmistellaan julkaisua Scoopille

## Mikä on jo automatisoitu

`scripts/package-scoop.ps1 -Version X.Y.Z` suorittaa `npm ci`, frontend build, `go test ./...`, Windows build `agentboard.exe` versionumerolla, luo ZIP-tiedoston ja laskee tämän ZIP:n SHA-256:n. **samasta** tiedostosta se luo `release/agentboard.json`:n URL-osoitteella, hashilla, MIT-lisenssillä, CLI-välilevyllä, pikakuvakkeella, `checkver` ja `autoupdate`. ZIP ja asennusohjelma sisältävät tiedoston `LICENSE`. Tietokanta sijaitsee asennushakemiston ulkopuolella, joten `persist`:ta ei tarvita luettelossa.

`.github/workflows/release.yml` tagissa `vX.Y.Z` suorittaa saman paketin Windows Runnerissa, rakentaa Inno Setup -asennusohjelman ja liittää ZIP-tiedoston, asennusohjelman ja luettelon GitHub-julkaisuun. Työnkulun manuaalinen suorittaminen luo vain artefaktin testausta varten julkaisematta julkaisua.

## Kirjausjärjestys ylläpitäjälle

1. Varmista, että koodi, dokumentaatio, versionumero ja tiedosto `LICENSE` ovat valmiit.
2. Run `./scripts/package-scoop.ps1 -Version X.Y.Z` locally. Tarkista `release/agentboard-X.Y.Z-windows-amd64.zip`-, `release/agentboard.json`- ja `agentboard version`-ulostulot pakkauksesta purkamisen jälkeen. Älä muuta ZIP-koodia sen jälkeen, kun tiiviste on laskettu.
3. Create and submit tag `vX.Y.Z`. GitHub Actions julkaisee julkaisun ZIP-, `agentboard-X.Y.Z-windows-amd64-setup.exe`- ja manifestin kanssa. Tarkista kaikki kolme tiedostoa Julkaisusivulta ja SHA-256 ZIP luettelosta.
4. Suorita puhtaalla Windows-koneella Scoopilla `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, sitten `agentboard version`, `agentboard init` ja `agentboard open` testikansiossa.
5. Pysyvää päivityssyötettä varten aseta luotu `agentboard.json` omaan Scoop-ämpäriisi tai tarjoa se sopivaan julkiseen ämpäriin. Tarkista `scoop update agentboard` seuraavan julkaisun jälkeen. Asennuskomento URL-osoitteesta sopii ensituttamiseen, kun taas ämpäri on kätevämpi päivityksiin.

Älä korvaa manuaalisesti satunnaista arvoa `hash`: Scoop tarkistaa ladatun ZIP-tiedoston sisällön.

Luettelon tarkistuksia varten keskity [muotoon Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [luontoluettelo](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) ja [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
