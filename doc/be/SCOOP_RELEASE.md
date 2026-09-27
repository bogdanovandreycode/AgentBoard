# Падрыхтоўка рэлізу для Scoop

## Што ўжо аўтаматызавана

`scripts/package-scoop.ps1 -Version 0.2.0` выконвае `npm ci`, зборку frontend, `go test ./...`, зборку Windows `agentboard.exe` з нумарам версіі, стварае ZIP і разлічвае SHA-256 гэтага ZIP. З **гэтага ж** файла ён стварае `release/agentboard.json` з URL, хешам, CLI-shim, цэтлікам, `checkver` і `autoupdate`. База жыве па-за каталогам усталёўкі, таму `persist` у manifest не патрэбен.

`.github/workflows/release.yml` на тэгу `vX.Y.Z` запускае тое ж пакаванне ў Windows runner, збірае ўсталёўнік Inno Setup і прыкладвае ZIP, усталёўшчык і manifest да GitHub Release. Ручны запуск workflow стварае толькі artifact для праверкі, без публікацыі рэлізу.

## Парадак публікацыі для суправаджаючага

1. Праверце, што код, дакументацыя і нумар версіі гатовы. Вызначыце ліцэнзію праекту: цяпер manifest паказвае `Unknown`, таму што файл ліцэнзіі ў рэпазітары не зададзены. Калі вы выбіраеце ліцэнзію, дадайце `LICENSE` і абновіце `license` у скрыпце перад рэлізам.
2. Лакальна выканайце `./scripts/package-scoop.ps1 -Version X.Y.Z`. Праверце `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` і выснову `agentboard version` пасля распакавання. Не мяняйце ZIP пасля разліку хеша.
3. Стварыце і дашліце тэг `vX.Y.Z`. GitHub Actions апублікуе Release з ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` і manifest. Праверце ўсе тры файлы на старонцы Release і SHA-256 ZIP з manifest.
4. На чыстай Windows-машыне з Scoop выканайце `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, затым `agentboard version`, `agentboard init` і `agentboard open` у тэставай тэчцы.
5. Для сталага канала абнаўленняў змесціце згенераваны `agentboard.json` ва ўласны Scoop bucket або прапануйце яго ў прыдатны публічны bucket. Правярайце `scoop update agentboard` пасля наступнага рэлізу. Каманда ўстаноўкі з URL падыходзіць для першага знаёмства, а bucket зручней для абнаўленняў.

Не падстаўляйце ўручную выпадковае значэнне `hash`: Scoop звярае змесціва запампаванага ZIP.

Для праверак manifest арыентуйцеся на [фармат Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [стварэнне manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) і [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
