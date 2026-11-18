# Proje yönetimi ve görev taslakları

**Tüm projeler** sayfasında kayıtlı klasör, yayınlanan proje sayısı gösterilir.
görevler ve en son etkinlik. Etkinlik; görevler, geçmiş ve özelliklerdeki değişiklikleri hesaba katar.
gider kayıtları. Kalem düğmesi projeyi yeniden adlandırmanıza veya seçmenize olanak tanır.
sistem iletişim kutusundaki mevcut klasör. Yolu değiştirmek kaydı ve işaretçiyi günceller;
Proje kimlikleri, çalışanlar, görevler ve ayarlar kaydedilir, dosyalar taşınmaz.
Zaten başka bir projeye kayıtlı olan klasörler seçilemez.

## Genel ayarlar

Proje seçim sayfasından veya sekmesinden **Genel ayarları** açın
**Genel ayarlar** herhangi bir panonun ayarlarında. Bu parametreler tüm projeler için ortaktır.
bir AgentBoard tabanı:

- Taslakların zamanlayıcıya göre otomatik olarak kaydedilmesi varsayılan olarak her **10 saniyede bir** etkinleştirilir.
- **Alandan ayrılırken kaydetme** varsayılan olarak etkindir ve ayrı olarak yapılandırılabilir.
- **MCP proje içeriği** kontrolleri genel olarak etkinleştirilir ve devre dışı bırakılır. Onların alanı
  eylemler ve ayrı bir CLI bayrağı [İşçiler ve MCP](WORKERS_MCP.md).

## Taslak yaşam döngüsü

Yeni Görev iletişim kutusu, tamamlanmamış çalışmayı taslak olarak kaydeder. Durum satırı
kaydedilmemiş değişiklikler, devam eden kaydetme, kayıtlı taslak ve hata arasında ayrım yapar
koruma. Taslak kurulda görünmez ve MCP çalışanlarının kullanımına sunulmaz.
yayınlandı.

Yayımlamadan önce **Başlık, Durum, kısmını doldurmak için **Görev oluştur** düğmesini kullanın.
Sorumlu ve Açıklama**. Bundan önce **Taslağı kaydet**, taslağı oluşturmadan kaydeder
görevler. Boş bir form kapatıldığında varsayılan olarak atılır. Temizlenirse
Mevcut taslağı boş duruma getirirseniz, bir sonraki kaydedişinizde silinir.
Boş olmayan bir iletişim kutusunu kapatmak da onu kaydeder.

Panodaki **Taslaklar** düğmesi kayıtlı taslaklara devam etmenize veya silmenize olanak tanır.
Çoklu taslaklar desteklenir. Yeniden başlatmanın ardından açık
taslak. Taslaklar SQLite'ta saklanır; Her değişiklikte yerel olan güncellenir
tarayıcıya kopyalayın, böylece kaydedilmemiş düzenlemeler yeniden başlatmalarda veya
geçici sunucu hatası. Tarayıcı depolama alanı mevcut değilse bir iletişim kutusu uyarır
taslağın kapatılmadan önce sunucuya kaydedilmesi gerektiğini.

Taslak yazarken, bir pencerenin işin üzerine yazılmasını önlemek için revizyonlar kullanılır
diğeri. Sunucuda zaten bir tane varsa, kayıp kaydetme yanıtı tutarlıdır
özdeş veriler. İçerik farklıysa bir çakışma gösterilir ve
**Yeni taslak olarak kurtar** yerel içeriği ayrı olarak kaydeder. Yayın
atomik ve değişmez: aynı taslağı kaybolduktan sonra yeniden yayınlayın
yanıt, bir kopya oluşturmak yerine aynı görevi döndürür.
