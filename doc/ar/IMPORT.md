# استيراد المهام من JSON

تقبل علامة التبويب **استيراد المهام** ملف JSON بتشفير UTF-8. استخدمه إذا كنت تقوم بنقل المهام من خدمة أخرى. انقر فوق **تنزيل قالب JSON**: يتضمن القالب أسماء الخصائص المخصصة الحالية و`Slug` للعامل المتوفر لهذا المشروع.

## خطوة بخطوة

1. قم بإنشاء العمال اللازمين (**العمال**) والخصائص (**الخصائص**) قبل الاستيراد.
2. قم بتنزيل القالب، وافتحه في محرر نصوص، واستبدل الأمثلة بمهامك الخاصة، واحفظ الملف بامتداد `.json`.
3. في علامة التبويب **استيراد المهام**، حدد ملفًا. سوف تظهر الواجهة عدد المهام.
4. انقر **استيراد المهام**. سيقوم الخادم بفحص الملف بأكمله: إذا تم اكتشاف خطأ، فلن تتم كتابة أي مهمة منه. قم بتصحيح رسالة الخطأ وحاول مرة أخرى.

الحد الأدنى للملف:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

مثال مع العامل والملكية والتبعية:

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

يجب أن يكون `version` `1`، الصفيف `tasks` - من 1 إلى 1000 مهمة. مطلوب `title`. المراحل الصالحة: `backlog`، `features`، `in_progress`، `testing`، `verification`، `complete`. الأولويات: `critical`، `high`، `medium`، `low`؛ أوضاع الاختبار: `ai`، `human`، `hybrid`. المسؤول: `unassigned` أو `human` أو `worker` مع `worker` (Slug أو ID) الحالي. يستخدم `properties` أسماء أو معرفات الخصائص التي تم إنشاؤها بالفعل. `key` فريد داخل الملف؛ يشير `depends_on` إلى مثل هذه المفاتيح. يقوم الخادم بإنشاء معرفات المهام بنفسه.

سيؤدي استيراد نفس الملف مرة أخرى إلى إنشاء مشكلات جديدة، لذا تحقق من اللوحة قبل النقر مرة أخرى. عند استيرادها، يتم تسجيل المهام على أنها من إنشاء الإنسان؛ يظهر إدخال النظام في السجل.
