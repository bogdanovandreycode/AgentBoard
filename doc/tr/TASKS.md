# Görevlerle çalışma

## Görev ekle

**Pano** sekmesinde **Yeni görev**'i tıklayın. Başlığı ve açıklamayı girin. Açıklama ve test talimatları Markdown'ı destekler: metni vurgulayın ve kalın, italik, bağlantı ve liste için biçimlendirme çubuğunu kullanın. **Kaydet**'i tıklayın.

Alanlar:

| Alan | ne işe yarar |
| --- | --- |
| Başlık | Görevin kısa adı. |
| Açıklama | Ne yapılması gerekiyor ve işin hazır olduğu nasıl anlaşılacak? |
| Eyalet | Çalışma aşaması; yeni bir görev genellikle `Backlog`'da başlar. |
| Öncelik | `critical`, `high`, `medium` veya `low`. |
| Sorumlu | Bir kişi, belirli bir çalışan veya randevusu olmayan. |
| Test modu | `AI` - AI'yi kontrol eder; `Human` - insan kontrolleri; `Hybrid` - her ikisi de. |
| Bağımlılıklar | Daha önce tamamlanması gereken görevler. |
| Yapay Zeka/İnsan testi talimatları | Uygun incelemeciye yönelik talimatlar. |
| Özel özellikler | **Özellikler** sekmesinde oluşturulan ek alanlar. |

Bir AI görevi için önce bir çalışan oluşturun, Sorumlu'da onu seçin, ardından kartı `Features`'ya aktarın. Yapay zeka, kendisine atanan görevleri yalnızca `Features`'dan `Verification`'ya kadar dört aşamada görüyor.

## Görevi taşı

Kartı sütunların arasına sürükleyin. Yatay kaydırmayla tüm sütunlardan oluşan görünüm ve geniş sütunlar arasında geçiş yapabilirsiniz. `Backlog` ve `Complete` insan kontrollüdür. Yapay zeka, `Features → In progress → Testing → Verification` görevini yalnızca özel MCP araçlarıyla ilerletebilir. Tamamlanmamış manuel teste sahip bir görev, yapay zeka doğrulamasını geçmemelidir.

## Görev kartı

Açıklamayı, sahibini, test talimatlarını ve geçmiş, testler, yapılar ve yapay zeka maliyetlerine ilişkin sekmeleri görmek için bir kartı tıklayın. **Görevi düzenle** içeriği değiştirir. Kişinin yorumu genel hikayeye eklenir. Üst filtrelerde arama; başlığa, açıklamaya, kimliğe ve çalışana göre arama yapar; ayrı bir düğme geniş bir aramayı açar.

## Sütunlar ve özellikler

**Ayarlar** sekmesinde, bir kişinin kullanabileceği sütunların sırasını değiştirebilir ve kendinizinkini ekleyebilirsiniz. Dört AI aşaması sabittir ve aynı sırada ilerler. Kullanıcı sütunu, bir kişi tarafından ertelenen görevlerin yeridir: Yapay zeka için böyle bir görevin durumu `Backlog`'dur. Bir sütun silindiğinde görevleri normal `Backlog`'ya geri döner.

**Özellikler** sekmesinde metin, sayı, bayrak, tarih, seçim ve URL gibi alanlar ekleyebilirsiniz. `Human only` alanı yapay zekadan gizler; `Agent read` okumaya izin verir, `Agent read/write` ayrıca desteklenen araçlar aracılığıyla yazmaya da olanak tanır. Bu, yapay zekanın görev adımlarına ilişkin haklarını değiştirmez.
