# Scoop için bir sürüm hazırlanıyor

## Zaten otomatik olan şey

`scripts/package-scoop.ps1 -Version X.Y.Z`, `npm ci`'yu, ön uç yapısını, `go test ./...`'yu, Windows derlemesi `agentboard.exe`'yu sürüm numarasıyla çalıştırır, bir ZIP oluşturur ve bu ZIP'in SHA-256'sını hesaplar. **Aynı** dosyadan URL, karma, MIT lisansı, CLI-shim, kısayol, `release/agentboard.json` ve `checkver` içeren `autoupdate`'yu oluşturur. ZIP ve yükleyici `LICENSE` dosyasını içerir. Veritabanı kurulum dizininin dışında olduğundan, bildirimde `persist`'ya gerek yoktur.

`.github/workflows/release.yml` etiketindeki `vX.Y.Z` aynı paketi Windows çalıştırıcısında çalıştırır, Inno Kurulum yükleyicisini oluşturur ve ZIP'i, yükleyiciyi ve bildirimi GitHub Sürümüne ekler. Bir iş akışını manuel olarak çalıştırmak, bir sürüm yayınlamadan yalnızca test amaçlı bir yapıt oluşturur.

## Bakımcı için sipariş yayınlanıyor

1. Kodun, belgelerin, sürüm numarasının ve `LICENSE` dosyasının hazır olduğunu doğrulayın.
2. `./scripts/package-scoop.ps1 -Version X.Y.Z`'yu yerel olarak çalıştırın. Paketi açtıktan sonra `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` ve `agentboard version` çıkışını kontrol edin. Hash hesaplandıktan sonra ZIP'i değiştirmeyin.
3. `vX.Y.Z` etiketini oluşturun ve gönderin. GitHub Actions, ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` ve manifest içeren bir Sürüm yayınlayacak. Sürüm sayfasındaki üç dosyayı ve bildirimden SHA-256 ZIP dosyasını kontrol edin.
4. Scoop'lu temiz bir Windows makinesinde, test klasöründe `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`'yu, ardından `agentboard version`, `agentboard init` ve `agentboard open`'yu çalıştırın.
5. Kalıcı bir güncelleme beslemesi için, oluşturulan `agentboard.json`'yu kendi Scoop kovanıza yerleştirin veya uygun bir genel kovaya sunun. Bir sonraki sürümden sonra `scoop update agentboard`'yu tekrar kontrol edin. Bir URL'den gelen kurulum komutu ilk tanışma için uygundur, kova ise güncellemeler için daha uygundur.

`hash` rastgele değerini manuel olarak değiştirmeyin: Scoop indirilen ZIP içeriğini kontrol eder.

Bildiri kontrolleri için [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [creating manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) ve [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).xml] öğelerine odaklanın.
