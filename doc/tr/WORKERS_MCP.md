# İşçiler ve MCP bağlantısı

## 1. Adım. Bir çalışan oluşturun

AgentBoard'da **İşçiler → Çalışan ekle**'yi açın. Bir müşteri profili seçin: Codex, Claude Code, Gemini CLI, Cursor, OpenCode, OpenCode aracılığıyla Ollama, VS Code Copilot veya başka bir MCP istemcisi. Profil tipik özellikleri (kod, testler, Git) doldurur; onları değiştirebilirsiniz. Ad kartta görünür ve `Slug`, MCP komutu için boşluk içermeyen kısa bir tanımlayıcıdır. **Kaydet**'i tıklayın.

Bir çalışan bir yapay zeka kişiliğine karşılık gelir. Farklı müşteriler veya ekipler için farklı çalışanlar oluşturun. Yetenek profili uzmanlığı açıklar ancak yapay zekanın haklarını görev aşamalarına genişletmez.

## Adım 2. Yapılandırmayı kopyalayın

Oluşturulan işçiyi açın. **MCP tanılama** bloğu, yapılandırma parçasını ve bunun ekleneceği dosyayı gösterir. **MCP yapılandırmasını kopyala**'yı tıklayın. Dosya zaten mevcutsa, önerilen sunucuyu diğer sunucuları silmeden mevcut `mcpServers`/`servers`/`mcp` nesnesine ekleyin.

Ana komut şuna benzer:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

`--project`, `agentboard init` aracılığıyla kaydettiğiniz klasörü işaret etmelidir. `--worker` — Oluşturulan çalışanın `Slug`'su. MCP istemcisi, araçlara ihtiyaç duyduğunda bu komutu kendisi çalıştırır. AgentBoard tarayıcısında web sunucusu ayrı olarak çalışabilir.

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

Codex ve Claude Code için `C:\Projects\MyFirstProject` projesi ve çalışan `codex` ile bir örnek:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
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

JSON'da Windows ters eğik çizgi iki katına çıkar; arayüzdeki hazır bir parça bunu otomatik olarak yapar. `--db`'yu özel bir tabanla kullanıyorsanız, bunu `args` MCP yapılandırmasına ekleyin ve `open`'yu başlatırken kullandığınız yolun aynısını belirtin.

**Ollama** yerel bir model sağlar ancak MCP istemcisinin yerini almaz. **Ollama via OpenCode** profili, OpenCode için bir MCP yapılandırması oluşturur; OpenCode'u Ollama modelinde ayrı olarak yapılandırın. Ollama ile MCP uyumlu başka bir istemci de uygundur.

## 4. Adım: Bağlantınızı kontrol edin

1. Çalışan kartını açın ve **Tekrar kontrol et** seçeneğine tıklayın. **Sunucu kontrolü** MCP araçlarının sayısını göstermelidir. Bu, sunucu araçlarının dahili protokol kontrolü ve tespitidir.
2. Yapılandırma dosyasını ekledikten sonra AI istemcisini başlatın veya yeniden başlatın. Ondan `get_my_board`'yu aramasını isteyin.
3. **İstemci bağlandı** AgentBoard'da görünecek ve listede yeni bir oturum görünecektir. Yalnızca bu, müşterinizin bağlantısını doğrular. Araçlar varsa ancak istemci bağlı değilse programın yolunu, yapılandırma dosyasının adını ve JSON/TOML sözdizimini kontrol edin.

Çalışma oturumunuza `get_my_board` ile başlayın. AI, kimliklerini bilse bile `Backlog` ve `Complete`'dan görev almaz. Yapay zeka insanmış gibi davranamaz ve genel bir "herhangi bir yere git" komutuna sahip değildir. AgentBoard dışında dosya sistemiyle çalışmak, seçilen istemcinin yeteneklerine ve sanal alanına bağlıdır.

Resmi müşteri talimatları: [Codex](https://developers.openai.com/learn/docs-mcp), [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp), [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/), [Cursor](https://prod.cursor.com/docs/cli/mcp), [OpenCode](https://opencode.ai/docs/mcp-servers/), [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
