# Tuo tehtäviä JSONista

**Tuontitehtävät** -välilehti hyväksyy JSON-tiedoston UTF-8-koodauksella. Käytä sitä, jos siirrät tehtäviä toisesta palvelusta. Napsauta **Lataa JSON-malli**: Malli sisältää tämän projektin käytettävissä olevan työntekijän nykyiset mukautetut ominaisuusnimet ja `Slug`.

## Askel askeleelta

1. Luo tarvittavat työntekijät (**Workers**) ja ominaisuudet (**Properties**) ennen tuontia.
2. Lataa malli, avaa se tekstieditorissa, korvaa esimerkit omilla tehtävilläsi ja tallenna tiedosto `.json`-tunnisteella.
3. Valitse tiedosto **Tuo tehtävät** -välilehdeltä. Käyttöliittymä näyttää tehtävien määrän.
4. Napsauta **Tuo tehtäviä**. Palvelin tarkistaa koko tiedoston: jos virhe havaitaan, siitä ei kirjoiteta yhtään tehtävää. Korjaa virheilmoitus ja yritä uudelleen.

Minimi tiedosto:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Esimerkki työntekijästä, omaisuudesta ja huollosta:

```json
{
  "version": 1,
  "tasks": [
    {
      "key": "design",
      "title": "Prepare the design",
      "description": "## Goal\nPrepare the home page mockup.",
      "state": "features",
      "priority": "high",
      "testing_mode": "hybrid",
      "assignee": { "type": "worker", "worker": "codex" },
      "ai_test_instructions": "Check the build.",
      "human_test_instructions": "Review the page in a browser.",
      "properties": { "Department": "Design" }
    },
    {
      "title": "Approve the design",
      "depends_on": ["design"],
      "assignee": { "type": "human" }
    }
  ]
}
```

`version`:n tulisi olla `1`, taulukon `tasks` - 1 - 1000 tehtävää. `title` vaaditaan. Kelvolliset vaiheet: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Prioriteetit: `critical`, `high`, `medium`, `low`; testitilat: `ai`, `human`, `hybrid`. Vastuuhenkilö: `unassigned`, `human` tai `worker` olemassa olevan `worker`:n kanssa (Slug tai ID). `properties` käyttää jo luotujen ominaisuuksien nimiä tai tunnuksia. `key` on ainutlaatuinen tiedostossa; `depends_on` viittaa tällaisiin avaimiin. Palvelin luo itse tehtävätunnukset.

Saman tiedoston tuominen uudelleen luo uusia ongelmia, joten tarkista taulu ennen kuin napsautat uudelleen. Kun tehtävät tuodaan, ne kirjataan ihmisen luomina; Järjestelmämerkintä näkyy historiassa.
