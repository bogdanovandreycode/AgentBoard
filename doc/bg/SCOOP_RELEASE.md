# Подготовка на издание за Scoop

## Какво вече е автоматизирано

`scripts/package-scoop.ps1 -Version 0.2.0` изпълнява `npm ci`, компилация на интерфейса, `go test ./...`, компилация на Windows `agentboard.exe` с номер на версия, създава ZIP и изчислява SHA-256 на този ZIP. От **същия** файл той създава `release/agentboard.json` с URL, хеш, CLI-shim, пряк път, `checkver` и `autoupdate`. Базата данни се намира извън инсталационната директория, така че `persist` не е необходим в манифеста.

`.github/workflows/release.yml` на етикета `vX.Y.Z` изпълнява същата опаковка в Windows runner, изгражда инсталатора на Inno Setup и прикачва ZIP, инсталатора и манифеста към GitHub Release. Ръчното изпълнение на работен поток създава само артефакт за тестване, без да публикува версия.

## Поръчка за публикуване за поддържащия

1. Проверете дали кодът, документацията и номерът на версията са готови. Дефиниране на лиценза на проекта: манифестът в момента показва `Unknown`, тъй като лицензният файл не е дефиниран в хранилището. Ако изберете лиценз, моля, добавете `LICENSE` и актуализирайте `license` в скрипта преди пускане.
2. Стартирайте `./scripts/package-scoop.ps1 -Version X.Y.Z` локално. Проверете изхода на `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` и `agentboard version` след разопаковането. Не променяйте ZIP кода, след като хешът е изчислен.
3. Създайте и изпратете етикет `vX.Y.Z`. GitHub Actions ще публикува издание с ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` и манифест. Проверете и трите файла на страницата за издание и SHA-256 ZIP от манифеста.
4. На чиста Windows машина със Scoop стартирайте `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, след това `agentboard version`, `agentboard init` и `agentboard open` в тестовата папка.
5. За постоянна емисия за актуализиране, поставете генерирания `agentboard.json` във вашата собствена кофа Scoop или я предложете в подходяща обществена кофа. Проверете отново за `scoop update agentboard` след следващото издание. Командата за инсталиране от URL е подходяща за първо запознаване, докато кофата е по-удобна за актуализации.

Не замествайте ръчно произволната стойност `hash`: Scoop проверява съдържанието на изтегления ZIP файл.

За проверки на манифест се фокусирайте върху [format Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [creating manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) и [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
