# Työntekijät ja MCP-yhteys

## Vaihe 1. Luo työntekijä

Avaa AgentBoardissa **Työntekijät → Lisää työntekijä**. Valitse asiakasprofiili: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, Ollama kautta OpenCode, VS Code Copilot tai muu MCP-asiakas. Profiili täyttää tyypilliset ominaisuudet (koodi, testit, Git); voit muuttaa niitä. Nimi näkyy taululla, ja `Slug` on lyhyt tunniste ilman välilyöntejä MCP-komennolle. Napsauta **Tallenna**.

Yksi työntekijä vastaa yhtä AI-persoonallisuutta. Luo erilaisia ​​työntekijöitä eri asiakkaille tai ryhmille. Kykyprofiili kuvaa erikoistumista, mutta ei ulota tekoälyn oikeuksia tehtävävaiheisiin.

## Vaihe 2. Kopioi kokoonpano

Avaa luotu työntekijä. **MCP-diagnostiikka** -lohko näyttää konfiguraatiofragmentin ja tiedoston, johon se lisätään. Napsauta **Kopioi MCP-asetukset**. Jos tiedosto on jo olemassa, lisää ehdotettu palvelin olemassa olevaan `mcpServers`/`servers`/`mcp`-objektiin poistamatta muita palvelimia.

Pääkomento näyttää tältä:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project`:n pitäisi osoittaa kansioon, jonka rekisteröit `agentboard init`:n kautta. `--worker` — luodun työntekijän `Slug`. MCP-asiakas suorittaa tämän komennon itse, kun se tarvitsee työkaluja. AgentBoard-selaimessa verkkopalvelin voi toimia erikseen.

Tarkistuksen jälkeen diagnoosilohko näyttää absoluuttisen polun käynnissä olevaan `agentboard.exe`:hen. Se on erityisen hyödyllinen ZIP-tiedostosta asennettaessa. Kun asennat Scoopin kautta, voit käyttää `agentboard`-komentoa, jos asiakas näkee saman `PATH`:n.

## Vaihe 3: Lisää palvelin asiakkaallesi

Työntekijän käyttöliittymässä on jo valmis fragmentti. Alla on selitys sen käytöstä:

| Asiakas | Mihin lisätä | Kuinka tarkistaa asiakkaan puolella |
| --- | --- | --- |
| Codex | `%USERPROFILE%\.codex\config.toml`, osa `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Code | `.mcp.json` projektikansiossa | `claude mcp list` |
| Gemini CLI | `%USERPROFILE%\.gemini\settings.json`, esine `mcpServers` | `/mcp list` Gemini CLI:ssä |
| Kursori | `.cursor\mcp.json` projekti | luettelo MCP-palvelimista Kohdistimen asetuksissa |
| OpenCode | `opencode.json` projekti, kohde `mcp` | luettelo MCP-työkaluista OpenCode |
| VS Code Copilot | `.vscode\mcp.json`-projekti, kohde `servers` | komento **MCP: List Servers** |

Codexille ja Claude Codelle esimerkki projektista `C:\Projects\MyFirstProject` ja työntekijästä `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

JSONissa Windowsin kenoviiva kaksinkertaistuu; käyttöliittymän valmis fragmentti tekee tämän automaattisesti. Jos käytät `--db`:ta mukautetulla alustalla, lisää se `args` MCP-kokoonpanoon ja määritä sama polku kuin käynnistettäessä `open`.

**Ollama** tarjoaa paikallisen mallin, mutta ei korvaa MCP-asiakasta. **Ollama kautta OpenCode** -profiili luo MCP-määrityksen OpenCodea varten; määritä OpenCode erikseen Ollama-mallissa. Myös toinen MCP-yhteensopiva asiakas Ollaman kanssa sopii.

## Vaihe 4: Tarkista yhteys

1. Avaa työntekijäkortti ja napsauta **Tarkista uudelleen**. **Palvelintarkistuksen** pitäisi näyttää MCP-työkalujen lukumäärä. Tämä on sisäinen protokollatarkistus ja palvelintyökalujen havaitseminen.
2. Käynnistä tai käynnistä AI-asiakas uudelleen määritystiedoston lisäämisen jälkeen. Pyydä häntä soittamaan `get_my_board`.
3. **Asiakas yhdistetty** näkyy AgentBoardissa ja uusi istunto tulee näkyviin luetteloon. Vain tämä vahvistaa asiakkaasi yhteyden. Jos työkaluja on, mutta asiakasta ei ole yhdistetty, tarkista ohjelman polku, määritystiedoston nimi ja sen JSON/TOML-syntaksi.

Aloita työistunto `get_my_board`:lla. Tekoäly ei vastaanota tehtäviä `Backlog`:lta ja `Complete`:lta, vaikka se tietäisi heidän tunnuksensa. Tekoäly ei voi teeskennellä olevansa ihminen, eikä sillä ole yleistä "siirry minnekään" -komentoa. Työskentely tiedostojärjestelmän kanssa AgentBoardin ulkopuolella riippuu valitun asiakkaan ominaisuuksista ja hiekkalaatikosta.

Viralliset asiakasohjeet: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp)], [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
