# Ayarlar

Seçilen projenin yan menüsünde **Ayarlar**'ı açın. Değişiklikler **Ayarları kaydet** düğmesiyle kaydedilir ve proje için AgentBoard veritabanında saklanır.

## Dil ve görünüm

**Dil** alanı bir arama açılır listesini açar. Varsayılan olarak, tarayıcı dilini alan **Sistemi takip et** seçeneği seçilidir. Ekran görüntülerinden kullanılabilen diller: Arapça, Portekizce (Brezilya), Basitleştirilmiş Çince, Çekçe, Danca, Felemenkçe, İngilizce, Fince, Fransızca, Almanca, İtalyanca, Japonca, Korece, Norveççe Bokmål, Lehçe, Rusça, İspanyolca, İsveççe, Türkçe, Ukraynaca, Vietnamca; ayrıca Belarusça, Romence ve Bulgarca. **Takip sistemi** tarayıcı dilini alır. Ana imzalar manuel olarak çevrilir, geri kalan kullanıcı arayüzü satırlarının ön otomatik çevirisi vardır. Halka açıklanmadan önce, çevirilerin anadili İngilizce olan kişiler tarafından yeniden okunması tavsiye edilir; istemci adları, komutlar, JSON alanları ve kullanıcı verileri tercüme edilmeden kalır.

**Renk şeması**: Koyu (orijinal tema), Açık, Siyah, Ubuntu ve Windows. **Saat dilimi** tarihlerin görüntülenmesini kontrol eder; veriler UTC'de saklanmaya devam eder. **Sistem saat dilimi** bilgisayar ayarlarını kullanır.

## Sütunlar

Dört aşama `Features`, `In progress`, `Testing`, `Verification` bu sırayla sabitlenmiştir: yeniden adlandırılamaz veya silinemez. Diğer sütunlar mevcut oklar kullanılarak yeniden düzenlenebilir. Bir ad girin ve özel bir sütun oluşturmak için **Sütun ekle**'yi tıklayın. İnsan görevleri için tasarlanmıştır: Yapay zeka, `Backlog` gibi bir görevi görür ve onu MCP aracılığıyla almaz. Kaydederken bir sütunun silinmesi, görevlerini normal `Backlog`'ya aktarır.

## Web ve MCP

**Pano yenileme**, panonun ve kartların yenileme hızını saniye cinsinden (1-60) ayarlar. **Çalışan yenileme** çalışanların durumunu günceller (2-120 saniye). Bu, AI tetikleme frekansı değil, web arayüzü yoklamasıdır. İşçiler otomatik olarak başlamazlar.

Web sunucusu adresi CLI başlatılırken ayarlanır, örneğin `agentboard open --addr 127.0.0.1:7444`. Adresi değiştirmek için sunucunun yeniden başlatılması gerekir. Varsayılan `127.0.0.1:7337`'dur. MCP ayrı bir yerel komut `agentboard mcp --project ... --worker ...` aracılığıyla çalışır ve web bağlantı noktasından bağımsızdır. Standart olmayan bir taban için tüm komutlarda aynı `--db`'yu belirtin. Her istemcinin yapılandırmasını çalışan kartından kopyalayın.
