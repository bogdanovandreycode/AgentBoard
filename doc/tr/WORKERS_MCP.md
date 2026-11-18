# İşçiler ve MCP bağlantısı

## 1. Adım. Bir çalışan oluşturun

AgentBoard'da **İşçiler → Çalışan ekle**'yi açın. Bir müşteri profili seçin: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, OpenCode aracılığıyla Ollama, VS Code Copilot veya başka bir MCP istemcisi. Profil tipik özellikleri (kod, testler, Git) doldurur; onları değiştirebilirsiniz. Ad kartta görünür ve `Slug`, MCP komutu için boşluk içermeyen kısa bir tanımlayıcıdır. **Kaydet**'i tıklayın.

Bir çalışan bir yapay zeka kişiliğine karşılık gelir. Farklı müşteriler veya ekipler için farklı çalışanlar oluşturun. Yetenek profili uzmanlığı açıklar ancak yapay zekanın haklarını görev aşamalarına genişletmez.

## Adım 2. Yapılandırmayı kopyalayın

Oluşturulan işçiyi açın. **MCP tanılama** bloğu, yapılandırma parçasını ve bunun ekleneceği dosyayı gösterir. **MCP yapılandırmasını kopyala**'yı tıklayın. Dosya zaten mevcutsa, önerilen sunucuyu diğer sunucuları silmeden mevcut `mcpServers`/`servers`/`mcp` nesnesine ekleyin.

Ana komut şuna benzer:

```text
agentboard mcp --worker codex
```

`--worker` — Oluşturulan çalışanın `Slug`'su. MCP istemcisi, araçlara ihtiyaç duyduğunda bu komutu kendisi çalıştırır.

### MCP web arayüzü olmadan çalışır

`agentboard mcp`, AgentBoard veritabanını kendisi açan ayrı bir stdio işlemidir. `agentboard open`/`serve`'ya ihtiyaç duymaz, web arayüzünün başlatıldığı klasöre ve web arayüzünün kaç örneğinin çalıştığına bağlı değildir. Birden fazla aracı, her biri kendi projesinde eşzamanlı olarak çalışabilir: her istemci kendi MCP işlemini çalıştırır ve hepsi aynı SQLite veritabanını güvenli bir şekilde kullanır.

### Bir proje nasıl seçilir

MCP sunucusu herhangi bir klasörde çalışır. Aracı ilk kez çağırdığınızda projeye bağlanır:

1. `--project PATH`, eğer belirtilmişse (kayıtlı proje içindeki herhangi bir klasör uygundur);
2. Aksi takdirde, müşteri bunları rapor ediyorsa, müşterinin çalışma klasörleri (kökler);
3. aksi takdirde istemcinin MCP sunucusunu başlattığı klasör.

`get_my_board`, seçilen projeyi ve çalışanı `session` ile birlikte `binding_source`'ya döndürür. İstemcinin çalışma klasörleri değişirse, araca yapılan bir sonraki çağrı bağlamayı yeniden gerçekleştirir. `--worker` olmadan projenin tek etkin çalışanı kullanılır.

Bağlama mümkün değilse her araç, aracının kullanıcıya iletmesi gereken bir hata döndürür; bu durumda görev verileri okunmaz:

| Kod | Anlamı |
| --- | --- |
| `NO_PROJECTS` | Veritabanında proje yok. Proje klasöründe `agentboard init`'yu yürütün ve bir çalışan oluşturun. Yeniden başlatmaya gerek yok. |
| `PROJECT_REQUIRED` | Çalışma klasörü kayıtlı projelerin dışında. Mesajda projeler listelenir. |
| `PROJECT_NOT_FOUND` | `--project` kayıtlı bir projenin içinde değil. |
| `PROJECT_AMBIGUOUS` | Bir müşterinin çalışma klasörleri birden fazla projeye yayılır. `--project`'yu ekleyin. |
| `NO_WORKERS`, `WORKER_REQUIRED`, `WORKER_NOT_FOUND`, `WORKER_DISABLED` | `--worker` aracılığıyla bir çalışan oluşturun, etkinleştirin veya seçin. |

**Codex** ve **Gemini CLI** yapılandırmaları özeldir ve tüm klasörlerde paylaşılır, dolayısıyla `--project` içermezler: Bu istemcinin çalıştığı her projede aynı bilgi notuna sahip bir çalışan oluşturun. Proje düzeyindeki dosyalar (Claude Code, Cursor, OpenCode, VS Code) projelerinin `--project`'sunu içerir.

### Açık tasarım ve içerik kontrolü

`--project` belirtilirse ve istemcinin başlangıç klasörü veya çalışma klasörü farklı bir projedeyse araçlar, MCP oluşturma yoluyla onay ister. Form oluşturmayan istemciler, görevlere erişmeden önce `PROJECT_CONTEXT_MISMATCH` alırlar. Yapılandırmayı düzeltin, otomatik algılama için `--project`'yu kaldırın veya `--allow-project-mismatch` bayrağıyla başka bir projeyle çalışmayı planladığınızı açıkça onaylayın. Çalışma klasörlerini değiştirmek önceki onayları sıfırlar. Otomatik olarak tanımlanan projeler onay gerektirmez.

**Genel ayarlar → MCP proje içeriği**, aynı veritabanındaki tüm projeler ve çalışanlar için bu kontrolleri kontrol eder. Varsayılan olarak etkindir. Devre dışı bırakılırsa MCP, belirtilen projeyi onaylar veya uyumsuzluk hataları olmadan kullanır. Zaten açık olan oturumlar, aracın bir sonraki çağrılışında bu ayarı uygular; görev görünürlüğü, geçişler ve insan testi kurallarına uyulmaya devam edilir.

Kontrolden sonra tanılama bloğu çalışan `agentboard.exe`'ya giden mutlak yolu gösterir. Özellikle ZIP'ten yükleme yaparken kullanışlıdır. Scoop aracılığıyla kurulum yaparken istemci aynı `agentboard`'yu görürse `PATH` komutunu kullanabilirsiniz.

## 3. Adım: İstemcinize bir sunucu ekleyin

Çalışan arayüzünde zaten hazır bir parça var. Aşağıda nerede kullanıldığına dair bir açıklama bulunmaktadır:

| Müşteri | Nereye eklenmeli | İstemci tarafında nasıl kontrol edilir |
| --- | --- | --- |
| Kodeks | `%USERPROFILE%\.codex\config.toml`, bölüm `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| Claude Kodu | `.mcp.json` proje klasöründe | `claude mcp list` |
| İkizler CLI | `%USERPROFILE%\.gemini\settings.json`, nesne `mcpServers` | Gemini CLI'de `/mcp list` |
| İmleç | `.cursor\mcp.json` projesi | İmleç ayarlarında MCP sunucularının listesi |
| Açık Kod | `opencode.json` projesi, nesne `mcp` | OpenCode'daki MCP araçlarının listesi |
| VS Kodu Yardımcı Pilotu | `.vscode\mcp.json` projesi, nesne `servers` | komutu **MCP: Sunucuları Listele** |

`C:\Projects\MyFirstProject` projesi ve çalışan `codex` ile Codex (kullanıcı yapılandırması, proje otomatik olarak algılanır) ve Claude Kodu (proje dosyası) örnekleri:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

JSON'da Windows ters eğik çizgi iki katına çıkar; arayüzdeki hazır bir parça bunu otomatik olarak yapar. Oluşturulan konfigürasyonlar web sunucusunun gerçek `--db`'sunu içerir, dolayısıyla istemci ortamı farklı olsa bile HTTP ve MCP aynı tabanı kullanır.

**Ollama** yerel bir model sağlar ancak MCP istemcisinin yerini almaz. **Ollama via OpenCode** profili, OpenCode için bir MCP yapılandırması oluşturur; OpenCode'u Ollama modelinde ayrı olarak yapılandırın. Ollama ile MCP uyumlu başka bir istemci de uygundur.

## 4. Adım: Bağlantınızı kontrol edin

1. Çalışan kartını açın ve **Tekrar kontrol et** seçeneğine tıklayın. **Sunucu kontrolü** bu proje, çalışan ve taban ile gerçek bir çalıştırılabilir stdio çalıştırır, el sıkışmayı, araç listesini ve proje ile çalışanın konumunu kontrol eder.
2. Yapılandırma dosyasını ekledikten sonra AI istemcisini başlatın veya yeniden başlatın. Ondan `get_my_board`'yu aramasını isteyin.
3. **İstemci bağlandı** AgentBoard'da görünecek ve listede yeni bir oturum görünecektir. Yalnızca bu, müşterinizin bağlantısını doğrular. Araçlar varsa ancak istemci bağlı değilse programın yolunu, yapılandırma dosyasının adını ve JSON/TOML sözdizimini kontrol edin.

Çalışma oturumunuza `get_my_board` ile başlayın. AI, kimliklerini bilse bile `Backlog` ve `Complete`'dan görev almaz. Yapay zeka insanmış gibi davranamaz ve genel bir "herhangi bir yere git" komutuna sahip değildir. AgentBoard dışında dosya sistemiyle çalışmak, seçilen istemcinin yeteneklerine ve sanal alanına bağlıdır.

Her araç, `title` ve dört MCP davranış ipucunun tamamını (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`) açık boole değerleri olarak bildirir. Okuma araçları salt okunurdur ve bağımsızdır; kayıt araçları yalnızca kayıt ekler veya görevleri ileriye taşır; dolayısıyla hiçbiri kesintiye uğramaz veya yerel AgentBoard tabanının ötesine geçmez.

Resmi müşteri talimatları: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
