# Začněte zde

## Co je AgentBoard

Představte si běžnou desku s kartami úkolů. Vy vytváříte úkoly, přidělujete někoho zodpovědného a dohlížíte na práci. Pracovníci AI dostávají pouze své přidělené úkoly prostřednictvím MCP a hlásí pokrok. Vy rozhodnete, kdy bude úkol konečně připraven.

**Projekt** - složka v počítači a samostatná deska. **Úkol** - karta s popisem, odpovědnou osobou a etapou. **Worker** je logické jméno klienta AI. **MCP** je způsob, jakým se klient AI připojuje k AgentBoard. **Scoop** je správce instalace softwaru pro Windows.

Není vyžadován žádný kód. Potřebujete Windows, prohlížeč, PowerShell a aby umělá inteligence fungovala, nainstalovaný klient AI s podporou místních serverů MCP.

## První spuštění

1. Nainstalujte aplikaci podle [pokynů](INSTALL.md).
2. Vytvořte složku projektu v Průzkumníku, například `C:\Projects\MyFirstProject`.
3. Otevřete tuto složku v Průzkumníkovi. Klikněte do adresního řádku, napište `powershell` a stiskněte Enter.
4. V okně, které se otevře, proveďte:

```powershell
agentboard init
agentboard open
```

5. Otevře se prohlížeč s adresou `http://127.0.0.1:7337`. Při používání tabule nechte okno PowerShellu otevřené. Zavřením okna se zastaví místní server, ale úlohy zůstanou.

Pokud příkaz `agentboard` není nalezen, zavřete PowerShell a znovu jej otevřete po instalaci Scoop. Při instalaci ze ZIP použijte úplnou cestu k `agentboard.exe`.

## První úkol

Klikněte na **Nový úkol**, vyplňte **Název**, v případě potřeby **Popis** a poté **Uložit**. Nová úloha v `Backlog` je dostupná lidem. Aby AI mohla začít pracovat, přiřaďte pracovníka a přeneste úkol do `Features`. Podrobný popis polí - [TASKS.md](TASKS.md).

## První pracovník

Otevřete **Workers → Add worker**. Vyberte klienta AI, kterého používáte (například Codex nebo Claude Code), zkontrolujte jméno a krátké ID `Slug`, klikněte na **Uložit**. Otevřete vytvořeného pracovníka: je zde konfigurace MCP, tlačítko kopírování a kontrola serveru. Zkopírujte konfiguraci do AI klienta podle [WORKERS_MCP.md](WORKERS_MCP.md). Po připojení požádejte klienta, aby zavolal na `get_my_board`.

## Co číst dále

- [Úkoly, testy, příběhy a sloupky](TASKS.md)
- [Připojování pracovníků a kontrola MCP](WORKERS_MCP.md)
- [Importovat úlohy z JSON](IMPORT.md)
- [Jazyk, téma, časové pásmo a aktualizace](SETTINGS.md)
- [Typické problémy](TROUBLESHOOTING.md)
