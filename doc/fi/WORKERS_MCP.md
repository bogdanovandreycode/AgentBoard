# Työntekijät ja MCP-yhteys

## Vaihe 1. Luo työntekijä

Avaa AgentBoardissa **Työntekijät → Lisää työntekijä**. Valitse asiakasprofiili: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama kautta OpenCode, VS Code Copilot tai muu MCP-asiakas. Profiili täyttää tyypilliset ominaisuudet (koodi, testit, Git); voit muuttaa niitä. Nimi näkyy taululla, ja `Slug` on lyhyt tunniste ilman välilyöntejä MCP-komennolle. Napsauta **Tallenna**.

Yksi työntekijä vastaa yhtä AI-persoonallisuutta. Luo erilaisia ​​työntekijöitä eri asiakkaille tai ryhmille. Kykyprofiili kuvaa erikoistumista, mutta ei ulota tekoälyn oikeuksia tehtävävaiheisiin.

## Vaihe 2. Kopioi kokoonpano

Avaa luotu työntekijä. **MCP-diagnostiikka** -lohko näyttää konfiguraatiofragmentin ja tiedoston, johon se lisätään. Napsauta **Kopioi MCP-asetukset**. Jos tiedosto on jo olemassa, lisää ehdotettu palvelin olemassa olevaan `mcpServers`/`servers`/`mcp`-objektiin poistamatta muita palvelimia.

Pääkomento näyttää tältä:

```text
agentboard mcp --worker codex
```

`--worker` — luodun työntekijän `Slug`. MCP-asiakas suorittaa tämän komennon itse, kun se tarvitsee työkaluja.

### MCP toimii ilman verkkokäyttöliittymää

`agentboard mcp` on erillinen stdio-prosessi, joka itse avaa AgentBoard-tietokannan. Se ei vaadi `agentboard open`/`serve`:ta, se ei riipu kansiosta, josta verkkokäyttöliittymä käynnistetään, ja kuinka monta web-käyttöliittymän esiintymää on käynnissä. Useita agentteja voi ajaa samanaikaisesti, kukin omassa projektissaan: jokainen asiakas käyttää omaa MCP-prosessiaan, ja ne kaikki käyttävät turvallisesti samaa SQLite-tietokantaa.

### Miten projekti valitaan

MCP-palvelin toimii missä tahansa kansiossa. Kun kutsut työkalua ensimmäistä kertaa, se on sidottu projektiin:

1. `--project PATH`, jos määritetty (mikä tahansa rekisteröidyn projektin sisällä oleva kansio sopii);
2. muuten asiakkaan työkansiot (juuret), jos asiakas ilmoittaa ne;
3. muuten kansio, jossa asiakas käynnisti MCP-palvelimen.

`get_my_board` palauttaa valitun projektin ja työntekijän `session`:lle yhdessä `binding_source`:n kanssa. Jos asiakkaan työkansiot muuttuvat, seuraava työkalun kutsu suorittaa sidoksen uudelleen. Ilman `--worker`:ta käytetään projektin ainoaa käytössä olevaa työntekijää.

Jos sitominen ei ole mahdollista, jokainen työkalu palauttaa virheen, joka agentin on välitettävä käyttäjälle; tehtävätietoja ei lueta tässä tapauksessa:

| Koodi | Merkitys |
| --- | --- |
| `NO_PROJECTS` | Tietokannassa ei ole projekteja. Suorita `agentboard init` projektikansiossa ja luo työntekijä. Uudelleenkäynnistystä ei tarvita. |
| `PROJECT_REQUIRED` | Työkansio on rekisteröityjen projektien ulkopuolella. Viesti listaa projekteja. |
| `PROJECT_NOT_FOUND` | `--project` ei ole rekisteröidyn projektin sisällä. |
| `PROJECT_AMBIGUOUS` | Asiakkaan työkansiot kattavat useita projekteja. Lisää `--project`. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | Luo, ota käyttöön tai valitse työntekijä `--worker`:n kautta. |

**Codexin** ja **Gemini CLI:n** määritykset ovat mukautettuja ja jaettu kaikille kansioille, joten niissä ei ole `--project`: Luo työntekijä, jolla on sama slug jokaiseen projektiin, jossa tämä asiakas on käynnissä. Projektitason tiedostot (Claude Code, Cursor, OpenCode, VS Code) sisältävät projektin `--project`.

### Selkeä suunnittelu ja kontekstin tarkistus

Jos `--project` on määritetty ja asiakkaan käynnistyskansio tai työkansio on eri projektissa, työkalut pyytävät vahvistusta MCP-eliitioinnin kautta. Asiakkaat, joilla ei ole lomaketta, saavat `PROJECT_CONTEXT_MISMATCH`:n ennen tehtävien käyttöä. Korjaa kokoonpano, poista `--project` automaattista tunnistusta varten tai vahvista nimenomaisesti, että aiot työskennellä toisen projektin kanssa `--allow-project-mismatch`-lipulla. Työkansioiden muuttaminen nollaa aiemmat vahvistukset. Automaattisesti määritellyt projektit eivät vaadi vahvistusta.

**Yleiset asetukset → MCP-projektikonteksti** hallitsee näitä tarkastuksia kaikille saman tietokannan projekteille ja työntekijöille. Oletuksena käytössä. Jos se ei ole käytössä, MCP käyttää määritettyä projektia ilman vahvistuksia tai yhteensopivuusvirheitä. Jo avoimet istunnot käyttävät asetusta seuraavan kerran, kun työkalua kutsutaan; tehtävien näkyvyyttä, siirtymiä ja ihmistestaussääntöjä noudatetaan edelleen.

Tarkistuksen jälkeen diagnoosilohko näyttää absoluuttisen polun käynnissä olevaan `agentboard.exe`:hen. Se on erityisen hyödyllinen ZIP-tiedostosta asennettaessa. Kun asennat Scoopin kautta, voit käyttää `agentboard`-komentoa, jos asiakas näkee saman `PATH`:n.

## Vaihe 3: Lisää palvelin asiakkaallesi

Työntekijän käyttöliittymässä on jo valmis fragmentti. Alla on selitys sen käytöstä:

| Asiakas | Mihin lisätä | Kuinka tarkistaa asiakkaan puolella |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, osa `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` projektikansiossa | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, esine `mcpServers` | `/mcp list` Gemini CLI:ssä |
| Kursori | `.cursor\mcp.json`-projekti | luettelo MCP-palvelimista Kohdistimen asetuksissa |
| OpenCode | `opencode.json` projekti, kohde `mcp` | luettelo MCP-työkaluista OpenCode |
| VS Code Copilot | `.vscode\mcp.json` projekti, kohde `servers` | komento **MCP: List Servers** |

Esimerkkejä Codexista (käyttäjän kokoonpano, projekti tunnistettu automaattisesti) ja Claude Codesta (projektitiedosto) projektin `C:\Projects\MyFirstProject` ja työntekijän `codex` kanssa:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

JSONissa Windowsin kenoviiva kaksinkertaistuu; käyttöliittymän valmis fragmentti tekee tämän automaattisesti. Luodut kokoonpanot sisältävät verkkopalvelimen todellisen `--db`:n, joten HTTP ja MCP käyttävät samaa perustaa, vaikka asiakasympäristö olisikin erilainen.

**Ollama** tarjoaa paikallisen mallin, mutta ei korvaa MCP-asiakasta. **Ollama kautta OpenCode** -profiili luo MCP-määrityksen OpenCodea varten; määritä OpenCode erikseen Ollama-mallissa. Toinen Ollaman MCP-yhteensopiva asiakas sopii myös.

## Vaihe 4: Tarkista yhteys

1. Avaa työntekijäkortti ja napsauta **Tarkista uudelleen**. **Palvelimen tarkistus** suorittaa todellisen stdio-suoritettavan tiedoston tämän projektin, työntekijän ja tukiaseman kanssa, tarkistaa kättelyn, työkaluluettelon ja sen, että projekti ja työntekijä sijaitsevat.
2. Käynnistä tai käynnistä AI-asiakas uudelleen määritystiedoston lisäämisen jälkeen. Pyydä häntä soittamaan `get_my_board`.
3. **Asiakas yhdistetty** näkyy AgentBoardissa ja uusi istunto tulee näkyviin luetteloon. Vain tämä vahvistaa asiakkaasi yhteyden. Jos työkaluja on, mutta asiakasta ei ole yhdistetty, tarkista ohjelman polku, määritystiedoston nimi ja sen JSON/TOML-syntaksi.

Aloita työistunto `get_my_board`:lla. Tekoäly ei vastaanota tehtäviä `Backlog`:lta ja `Complete`:lta, vaikka se tietäisi heidän tunnuksensa. AI ei voi teeskennellä olevansa ihminen, eikä sillä ole yleistä "siirry minnekään" -komentoa. Työskentely tiedostojärjestelmän kanssa AgentBoardin ulkopuolella riippuu valitun asiakkaan ominaisuuksista ja hiekkalaatikosta.

Jokainen työkalu ilmoittaa `title` ja kaikki neljä MCP-käyttäytymisvinkkiä (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) eksplisiittisiksi loogisiksi arvoiksi. Lukutyökalut ovat vain luku -muotoisia ja tehokkaita; tallennustyökalut vain lisäävät tietueita tai siirtävät tehtäviä eteenpäin, joten mikään ei häiritse tai ylitä paikallista AgentBoard-pohjaa.

Viralliset asiakasohjeet: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp)], [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
