# Sorun çözme

| Belirti | Ne kontrol edilmeli |
| --- | --- |
| `agentboard` bulunamadı | Scoop'tan sonra PowerShell'i yeniden başlatın. ZIP'ten kurulum yaparken `agentboard.exe`'nun tam yolunu kullanın. |
| Derleme sonrasında eski arayüzü gösteren web sayfası | Çalışan sunucuyu durdurun Ctrl+C. `./scripts/build.ps1`'yu çalıştırın, yeni bir ikili dosya başlatın. Arayüz EXE'de yerleşik olduğundan Vite'ın Go'dan önce derlenmesi gerekir. Sayfayı yenileyin Ctrl+F5. |
| Bağlantı noktası 7337 meşgul | AgentBoard zaten çalışıyor olabilir. `http://127.0.0.1:7337`'yu açın veya eski işlemi sonlandırın. Diğer bağlantı noktası için `--addr` kullanın. |
| Proje bulunamadı | İstediğiniz klasörde `agentboard init` komutunu çalıştırın. Daha sonra ondan `agentboard open` veya `agentboard open C:\путь\к\проекту`. |
| Çalışan görevi görmüyor | Görevin bu belirli çalışana atanması ve `Features`, `In progress`, `Testing` veya `Verification`'da bulunması gerekir. AI, `Backlog`, `Complete` ve özel sütunları görmüyor. |
| MCP sunucu kontrolü mevcut ancak istemci bağlı değil | Sunucu kontrolü harici istemci ayarlarını kontrol etmez. İstemciyi yeniden başlatın, yapılandırma dosyasını, `agentboard.exe`, `--project`, `--worker` ve genel `--db` yolunu kontrol edin. `get_my_board`'yu aramayı isteyin. |
| Çalışan Çevrimdışı | İstemci MCP'yi sonlandırmış veya henüz başlatmamış olabilir. Kalp atışı olmadan geçen 90 saniyenin ardından oturumun bağlantısı kesilmiş sayılır. |
| JSON dosyası içe aktarılmıyor | `version: 1`'yu, gerekli `title`'yu, mevcut `Slug` çalışanlarını ve özellik adlarını kontrol edin. JSON, yorumlara veya sondaki virgüllere izin vermez. |
| `agentboard.exe` yeniden oluşturulamıyor | Windows, çalışan bir EXE'nin yerini alamaz. Sunucuyu Ctrl+C durdurun ve derlemeyi yeniden deneyin. |
| Güncellemeden sonra görevler kayboldu | `--db`'nun başka bir dosyayı işaret etmediğinden ve aynı Windows kullanıcısı olarak oturum açtığınızdan emin olun. Varsayılan taban `%AppData%\AgentBoard`'dur. |

Hata açıklanmıyorsa mesajın tam metnini, `agentboard version` ve Windows sürümlerini ve yeniden denemek için gereken adımları toplayın. Özel proje verilerini veya veritabanı içeriklerini açık bir sayıda yayınlamayın.
