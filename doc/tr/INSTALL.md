# Windows'ta kurulum ve başlatma

## Yükleyici (önerilir)

`agentboard-VERSION-windows-amd64-setup.exe`'yu [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases).jpg] sayfasından indirin. Kurulum sihirbazı bir klasör önerecektir; varsayılan `C:\AI\AgentBoard`'dur. `agentboard.exe` ve belgeleri kopyalayacak, bir kısayol oluşturacak ve seçilen klasörü `PATH` sistemine ekleyecektir. Kurulumdan sonra `agentboard` komutunun kullanılabilir olması için yeni bir terminal açın.

Proje klasörünüzde şunu çalıştırın:

```powershell
agentboard init
agentboard open
```

Arayüz `agentboard.exe`'nun içine yerleştirilmiştir; Go veya Node.js'nin ayrı kurulumuna gerek yoktur. Veriler `%AppData%\AgentBoard`'da saklanır ve program güncellendiğinde veya kaldırıldığında korunur. "Yüklü Uygulamalar" yoluyla kaldırma, `PATH`'daki kısayolları ve girişi kaldırır.

## Scoop'un Kurulumu

Scoop henüz kurulu değilse normal kullanıcınız olarak PowerShell'i açın ve [resmi Scoop](https://scoop.sh/) talimatlarını izleyin. Kurumsal bilgisayarınızda kısıtlamalar varsa yöneticinizle iletişime geçin; AgentBoard, Scoop olmadan bir ZIP'ten de başlatılabilir.

Alternatif olarak AgentBoard'u Scoop aracılığıyla yükleyin:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

GitHub Sürümündeki manifestin kullanılabilirliği [sayfa Releases](https://github.com/bogdanovandreycode/AgentBoard/releases) sayfasından kontrol edilebilir. Scoop'a manifest eklendikten sonra klasör, paket adına göre kurulabilir ve `scoop update agentboard` komutuyla güncellenebilir.

## Kepçesiz ZIP

Sürümlerden `agentboard-VERSION-windows-amd64.zip`'yu indirin, paketi örneğin `C:\Tools\AgentBoard`'ya açın. Proje klasöründeki PowerShell'de:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

Bir AI istemcisi için, program `agentboard.exe`'da bulunmuyorsa MCP yapılandırmasında `PATH`'nun tam yolunu belirtin.

## Kaynaktan derle

Go sürümünü `go.mod` ve Node.js 22 veya üzeri sürümlerden yükleyin. Deponun kökündeki PowerShell'de:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

`scripts/build.ps1`, ön ucu `internal/webui/dist`'ya birleştirir, Go testlerini çalıştırır ve bir `agentboard.exe`'yu birleştirir. Sıra önemlidir: Go oluşturulduğunda arayüz ikili dosyaya yerleşiktir. Eski `agentboard.exe open` çalışıyorsa, yeniden oluşturmadan önce onu durdurun (Ctrl+C), aksi halde Windows dosyayı değiştirmenize izin vermez.

## Veriler nerede

- `%AppData%\AgentBoard\agentboard.db` - görevler, projeler, çalışanlar ve ayarlar. `--db` bayrağıyla farklı bir dosya belirtebilirsiniz ancak `init`, `open`/`serve` ve `mcp` **aynı yola** sahip olmalıdır.
- `<ваш проект>\.agentboard\project.json` — proje tanımlayıcı. Bu dosya görev içermiyor.
- Sunucu varsayılan olarak yalnızca `127.0.0.1:7337`'yu dinler. Proje yolundan önce başka bir `--addr` adresi belirtin: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

`open` proje sayfasını açar ve bu adreste halihazırda çalışmakta olan bir sunucuyu yeniden kullanır. Bir tarayıcıda birden fazla proje açıksa bunları soldaki listeden seçin.
