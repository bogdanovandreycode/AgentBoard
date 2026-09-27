# Arkkitehtuuri ja API

AgentBoard on paikallinen Go-prosessi, joka käyttää SQLitea. Ihmisten HTTP-sovellusliittymällä, tekoälyn MCP-sovittimella ja verkkoliittymällä on yhteinen palvelulogiikka. SQLite - tilalähde; frontend ei ohita palvelinta. AI-työkaluilla on erillinen oikeuspinta, eikä `Backlog`/`Complete`-tila MCP:ssä ole käytettävissä. `System`:n historiatietue on itse sovelluksen luoma.

## Komponentit

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - verkkotunnuksen mallit ja virheet.
- `internal/service` - tehtävien siirrot, käyttöoikeudet, tuonti ja asetukset.
- `internal/persistence` - SQLite ja siirrot.
- `internal/httpapi` - Ihmisen HTTP-sovellusliittymä.
- `internal/mcpserver` - MCP-työkalut erilliselle työntekijälle.
- `web` - React/TypeScript-käyttöliittymä; kokoonpano päätyy malliin `internal/webui/dist` ja sisältyy EXE-tiedostoon.

## HTTP-perusreitit

| Menetelmä ja polku | Kohde |
| --- | --- |
| `GET /api/health` | Palvelimen tarkistus. |
| `GET /api/projects` | Rekisteröityneet projektit. |
| `GET /api/projects/{id}/board` | Projektin hallitus. |
| `GET/PUT /api/projects/{id}/settings` | Asetukset, järjestys ja sarakkeiden nimet. |
| `POST /api/projects/{id}/tasks` | Luo tehtävä. |
| `POST /api/projects/{id}/tasks/import` | Tuo atomisesti JSON-versio 1. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Kortti, muuta, poista. |
| `POST /api/tasks/{id}/move` | Liikkuminen henkilön toimesta. |
| `GET/POST /api/projects/{id}/workers` | Työntekijän luettelo ja luominen. |
| `GET /api/workers/{id}/mcp/check` | MCP-kättelyn ja työkalujen sisäinen testaus. |
| `GET /api/workers/{id}/sessions` | Istunnon diagnostiikka. |
| `GET/POST /api/projects/{id}/properties` | Mukautetut ominaisuudet. |

HTTP on paikalliselle luotettavalle käyttäjälle. Älä julkaise verkkoporttia Internetiin ilman sen omaa todennusta, verkkorajoituksia ja HTTPS:ää. MCP-palvelin käynnistetään käyttämällä `stdio`:ta tiettyä projektia ja työntekijää varten; aloita `get_my_board`:lla. Sen käyttöoikeudet rajoittavat toimia AgentBoardissa, mutta eivät korvaa AI-asiakkaan tiedostojärjestelmän hiekkalaatikkoa.

Käyttäjäsarakkeet tallennetaan erillään `tasks.state`:sta: ydin pitää tehtävän tällaisessa sarakkeessa `backlog`:ssa, ja `tasks.board_column` määrittää paikan Human boardilla. Tämä säilyttää saman AI-siirtymämallin. Sarakkeen poistaminen tyhjentää `board_column`-tehtävät.
