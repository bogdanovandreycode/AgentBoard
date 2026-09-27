# AjanBoard

[🌐 Languages](../LANGUAGES.md)

AgentBoard, insan ve yapay zeka çalışanlarının ortak görevler üzerinde çalıştığı ancak farklı haklara sahip olduğu yerel bir görev panosudur. Uygulama tek bir `agentboard.exe` dosyasıyla başlatılır, tarayıcıda web arayüzünü açar ve çalışanlara `stdio` aracılığıyla ayrı bir MCP sunucusu sağlar. Veriler bilgisayarınızda kalır.

**[Sıfırdan başlayın](START_HERE.md) · [Görevlerle çalışma](TASKS.md) · [MCP](WORKERS_MCP.md) aracılığıyla yapay zekaya bağlanma · [JSON](IMPORT.md)'yu içe aktarma · [Ayarlar](SETTINGS.md) · [Sorunları çözme](TROUBLESHOOTING.md)**

## Beş dakika içinde

1. `agentboard-VERSION-windows-amd64-setup.exe` yükleyicisini [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases).jpg adresinden indirin. Bir klasör önerecektir (varsayılan olarak `C:\AI\AgentBoard`) ve onu `PATH`'ya ekleyecektir. Scoop ve ZIP de mevcuttur.
2. Proje klasörünüzde PowerShell'i açın; örneğin `C:\Projects\MyApp`.
3. `agentboard init`'yu çalıştırın (ZIP için: `agentboard.exe` ve `init`'ya giden tam yol).
4. `agentboard open`'yu çalıştırın. `http://127.0.0.1:7337` açılacaktır.
5. **Yeni görev** düğmesini kullanarak bir görev ekleyin. Bir AI çalışanı için **Çalışanlar → Çalışan ekle**'yi açın, müşteri profilini seçin ve MCP yapılandırmasını kopyalayın.

Henüz bir proje klasörünüz yoksa Windows Gezgini'nde bir tane oluşturun. Bir proje, Git ve kod olmasa bile herhangi bir klasör olabilir.

## Scoop aracılığıyla kurulum

Sürümden sonra [Scoop](https://scoop.sh/)] zaten kurulu olan PowerShell'de:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

Geliştiriciler için [kaynak ](INSTALL.md)'dan derlenmiştir. Sürüm iş akışı, aynı yapıdan SHA-256 içeren bir yükleyici, ZIP ve Scoop bildirimi oluşturur. Scoop aracılığıyla yüklenen sürümün güncellenmesi: `scoop update agentboard`, manifestoyu pakete ekledikten sonra; ayrıntılar - [sürüm hazırlığı](SCOOP_RELEASE.md).

## Yönetim kurulu nasıl yapılandırılmıştır

`Backlog → Features → In progress → Testing → Verification → Complete`

Bir kişi görevleri tahtanın etrafında taşıyabilir. Yapay zeka yalnızca `Features → In progress → Testing → Verification`'yu hareket ettirebilir; `Complete`'ya nihai kabul bir insan tarafından gerçekleştirilir. Kullanıcı sütunları insanlara yöneliktir: İçlerindeki görev, MCP için `Backlog` durumunda kalır. Test modları: Yapay Zeka, İnsan ve Hibrit. Geçmiş, testler, yapılar ve yapay zeka maliyetleri göreve eklenir.

## Takımlar

```text
agentboard init [--db PATH] [project-path]
agentboard open [--addr 127.0.0.1:7337] [--db PATH] [project-path]
agentboard serve [--addr 127.0.0.1:7337] [--db PATH]
agentboard mcp --project PROJECT_PATH --worker WORKER_SLUG [--db PATH]
agentboard version
```

`init` klasörü kaydeder ve oraya yalnızca `.agentboard/project.json` yazar. SQLite üretim verileri, Windows kullanıcı yapılandırma dizininde (`%AppData%\AgentBoard\agentboard.db`), projenin dışında ve Scoop kurulumunun dışında bulunur. Paketin kaldırılması veya güncellenmesi bu verileri kaldırmamalıdır. Başka bir bilgisayara aktarmadan önce AgentBoard durdurulduğunda veritabanının bir kopyasını alın.

İşçiler - mantıksal hesaplar; AgentBoard'un kendisi Codex, Claude veya başka herhangi bir AI istemcisini çalıştırmaz. İstemci, belirli bir çalışan için yerel bir MCP işlemi başlatır. MCP, uygulama izin sınırıdır ve dosya izolasyonu için AI istemci sanal alanını kullanın.

## Geliştiriciler için

Yığın: Go, SQLite, resmi MCP Go SDK, React, TypeScript, Vite, PrimeReact, TanStack Query, dnd-kit. Önce ön uç oluşturun, ardından Go: `./scripts/build.ps1`. Web dosyaları `go:embed` aracılığıyla ikili dosyaya dahil edilir. Mimari ve API, [doc/ARCHITECTURE.md](ARCHITECTURE.md)'da açıklanmıştır.
