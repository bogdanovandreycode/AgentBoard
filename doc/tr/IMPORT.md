# Görevleri JSON'dan içe aktarın

**Görevleri içe aktar** sekmesi, UTF-8 kodlamasında bir JSON dosyasını kabul eder. Başka bir hizmetten görevleri aktarıyorsanız bunu kullanın. **JSON şablonunu indir** seçeneğini tıklayın: Şablon, bu proje için mevcut çalışanın geçerli özel özellik adlarını ve `Slug`'sunu içerir.

## Adım adım

1. İçe aktarmadan önce gerekli çalışanları (**İşçiler**) ve özellikleri (**Özellikler**) oluşturun.
2. Şablonu indirin, bir metin düzenleyicide açın, örnekleri kendi görevlerinizle değiştirin ve dosyayı `.json` uzantısıyla kaydedin.
3. **Görevleri içe aktar** sekmesinde bir dosya seçin. Arayüz görev sayısını gösterecektir.
4. **Görevleri içe aktar**'ı tıklayın. Sunucu tüm dosyayı kontrol edecektir: bir hata tespit edilirse, ondan tek bir görev yazılmayacaktır. Hata mesajını düzeltip tekrar deneyin.

Asgari dosya:

```json
{
  "version": 1,
  "tasks": [
    { "title": "Plan the project" },
    { "title": "Review the result", "state": "features", "priority": "high" }
  ]
}
```

Çalışan, mülk ve bağımlılıkla ilgili örnek:

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

`version`, `1`, `tasks` dizisi olmalıdır - 1 ila 1000 görev arasında. `title` gereklidir. Geçerli aşamalar: `backlog`, `features`, `in_progress`, `testing`, `verification`, `complete`. Öncelikler: `critical`, `high`, `medium`, `low`; test modları: `ai`, `human`, `hybrid`. Sorumlu: Mevcut `unassigned` (Slug veya ID) ile `human`, `worker` veya `worker`. `properties`, önceden oluşturulmuş özelliklerin adlarını veya kimliklerini kullanır. `key` dosya içinde benzersizdir; `depends_on` bu tür anahtarları ifade eder. Sunucu, görev kimliklerini kendisi oluşturur.

Aynı dosyayı tekrar içe aktarmak yeni sorunlara yol açacağından tekrar tıklamadan önce panoyu kontrol edin. İçe aktarıldığında görevler insan yapımı olarak kaydedilir; Geçmişte bir sistem girişi görünür.
