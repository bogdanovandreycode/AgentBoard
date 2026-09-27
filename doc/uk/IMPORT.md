# Імпорт завдань з JSON

Вкладка **Import tasks** приймає файл JSON у кодуванні UTF-8. Використовуйте її, якщо ви переносите завдання з іншого сервісу. Натисніть **Download JSON template**: шаблон включає поточні імена власних властивостей і `Slug` доступного воркера цього проекту.

## Покроково

1. Створіть необхідні воркери (**Workers**) та властивості (**Properties**) до імпорту.
2. Завантажте шаблон, відкрийте його у текстовому редакторі, замініть приклади своїми завданнями та збережіть файл із розширенням `.json`.
3. На вкладці **Import tasks** виберіть файл. Інтерфейс покаже кількість завдань.
4. Натисніть **Import tasks**. Сервер перевірить весь файл: якщо виявиться помилка, жодне завдання не буде записано. Виправте повідомлення про помилку та повторіть.

Мінімальний файл:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Приклад з воркером, властивістю та залежністю:

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

`version` повинен бути `1`, масив `tasks` - від 1 до 1000 завдань. `title` є обов'язковим. Допустимі етапи: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Пріоритети: `critical`, `high`, `medium`, `low`; режими перевірки: `ai`, `human`, `hybrid`. Відповідальний: `unassigned`, `human` або `worker` з існуючим `worker` (Slug або ID). `properties` використовує імена або ID вже створених властивостей. `key` унікальний усередині файлу; `depends_on` посилається на такі ключі. ID задач сервер створює сам.

Повторний імпорт файлу створить нові завдання, тому перед повторним натисканням перевірте дошку. При імпорті завдання записуються як створені людиною; історія з'являється системна запис.
