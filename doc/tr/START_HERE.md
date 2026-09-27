# Buradan başlayın

## AgentBoard Nedir?

Görev kartlarının olduğu normal bir tahta düşünün. Görevler yaratırsınız, sorumlu birini atarsınız ve işi denetlersiniz. Yapay zeka çalışanları, MCP aracılığıyla yalnızca kendilerine atanan görevleri alır ve ilerlemeyi bildirir. Görevin nihayet ne zaman hazır olacağına siz karar verirsiniz.

**Proje** - bilgisayardaki bir klasör ve ayrı bir pano. **Görev** - Açıklamayı, sorumlu kişiyi ve aşamayı içeren bir kart. **İşçi**, AI istemcisinin mantıksal adıdır. **MCP**, AI istemcisinin AgentBoard'a bağlanma şeklidir. **Scoop** Windows için bir yazılım yükleme yöneticisidir.

Kod gerekmez. Yapay zekanın çalışması için Windows'a, bir tarayıcıya, PowerShell'e ve yerel MCP sunucularını destekleyen yüklü bir yapay zeka istemcisine ihtiyacınız var.

## İlk lansman

1. Uygulamayı [talimatlar](INSTALL.md)'ya göre yükleyin.
2. Explorer'da `C:\Projects\MyFirstProject` gibi bir proje klasörü oluşturun.
3. Bu klasörü Explorer'da açın. Adres çubuğuna tıklayın, `powershell` yazın ve Enter tuşuna basın.
4. Açılan pencerede şunları yapın:

```powershell
agentboard init
agentboard open
```

5. `http://127.0.0.1:7337` adresine sahip bir tarayıcı açılacaktır. Beyaz tahtayı kullanırken PowerShell penceresini açık bırakın. Pencerenin kapatılması yerel sunucuyu durduracaktır ancak görevler kalacaktır.

`agentboard` komutu bulunamazsa PowerShell'i kapatın ve Scoop'u yükledikten sonra yeniden açın. ZIP'ten kurulum yaparken `agentboard.exe`'ya giden tam yolu kullanın.

## İlk görev

**Yeni görev**'i tıklayın, **Başlık**'ı, gerekiyorsa **Açıklama**'yı girin ve ardından **Kaydet**'i girin. `Backlog`'da insanlara yeni bir görev sunuluyor. Yapay zekanın çalışmaya başlaması için bir çalışan atayın ve görevi `Features`'ya aktarın. Alanların adım adım açıklaması - [TASKS.md](TASKS.md).

## İlk işçi

**İşçiler → İşçi ekle**'yi açın. Kullandığınız AI istemcisini seçin (örneğin Codex veya Claude Code), `Slug` adını ve kısa kimliğini kontrol edin ve **Kaydet**'i tıklayın. Oluşturulan işçiyi açın: bir MCP yapılandırması, bir kopyalama düğmesi ve bir sunucu kontrolü vardır. Yapılandırmayı [WORKERS_MCP.md](WORKERS_MCP.md)'ya göre AI istemcisine kopyalayın. Bağlandıktan sonra istemciden `get_my_board`'yu aramasını isteyin.

## Bundan sonra ne okunmalı?

- [Görevler, testler, hikayeler ve sütunlar](TASKS.md)
- [İşçilerin bağlanması ve MCP](WORKERS_MCP.md)'nun kontrol edilmesi
- [JSON](IMPORT.md)'dan görevleri içe aktar
- [Dil, tema, saat dilimi ve güncelleme](SETTINGS.md)
- [Tipik sorunlar](TROUBLESHOOTING.md)
