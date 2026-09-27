# Імпарт задач з JSON

Укладка **Import tasks** прымае файл JSON у кадоўцы UTF-8. Выкарыстоўвайце яе, калі пераносіце задачы з іншага сэрвісу. Націсніце **Download JSON template**: шаблон уключае бягучыя імёны карыстацкіх уласцівасцяў і `Slug` даступнага воркера гэтага праекта.

## Пакрокава

1. Стварыце неабходныя воркеры (**Workers**) і ўласцівасці (**Properties**) да імпарту.
2. Запампуйце шаблон, адкрыйце яго ў тэкставым рэдактары, заменіце прыклады сваімі задачамі і захавайце файл з пашырэннем `.json`.
3. На ўкладцы **Import tasks** абярыце файл. Інтэрфейс пакажа колькасць задач.
4. Націсніце **Import tasks**. Сервер праверыць увесь файл: калі выявіцца памылка, ні адна задача з яго не будзе запісана. Выпраўце паведамленне пра памылку і паспрабуйце.

Мінімальны файл:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Прыклад з воркерам, уласцівасцю і залежнасцю:

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

`version` павінен быць `1`, масіў `tasks` - ад 1 да 1000 задач. `title` абавязковы. Дапушчальныя этапы: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Прыярытэты: `critical`, `high`, `medium`, `low`; рэжымы праверкі: `ai`, `human`, `hybrid`. Адказны: `unassigned`, `human` ці `worker` з існуючым `worker` (Slug альбо ID). `properties` выкарыстоўвае імёны ці ID ужо створаных уласцівасцяў. `key` унікальны ўнутры файла; `depends_on` спасылаецца на такія ключы. ID задач сервер стварае сам.

Паўторны імпарт таго ж файла створыць новыя задачы, таму перад паўторным націскам праверце дошку. Пры імпарце задачы запісваюцца як створаныя чалавекам; у гісторыі з'яўляецца сістэмны запіс.
