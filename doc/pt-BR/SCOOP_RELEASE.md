# Preparando um lançamento para Scoop

## O que já está automatizado

`scripts/package-scoop.ps1 -Version X.Y.Z` executa `npm ci`, compilação de frontend, `go test ./...`, compilação do Windows `agentboard.exe` com número de versão, cria um ZIP e calcula o SHA-256 desse ZIP. A partir do **mesmo** arquivo ele cria `release/agentboard.json` com URL, hash, licença MIT, CLI-shim, atalho, `checkver` e `autoupdate`. O ZIP e o instalador contêm o arquivo `LICENSE`. O banco de dados reside fora do diretório de instalação, portanto `persist` não é necessário no manifesto.

`.github/workflows/release.yml` na tag `vX.Y.Z` executa o mesmo pacote no executor do Windows, cria o instalador do Inno Setup e anexa o ZIP, o instalador e o manifesto ao lançamento do GitHub. A execução manual de um fluxo de trabalho cria apenas um artefato para teste, sem publicar uma versão.

## Postagem de ordem para mantenedor

1. Verifique se o código, a documentação, o número da versão e o arquivo `LICENSE` estão prontos.
2. Execute `./scripts/package-scoop.ps1 -Version X.Y.Z` localmente. Verifique a saída `release/agentboard-X.Y.Z-windows-amd64.zip`, `release/agentboard.json` e `agentboard version` após desembalar. Não altere o ZIP após o cálculo do hash.
3. Crie e envie a tag `vX.Y.Z`. GitHub Actions publicará um lançamento com ZIP, `agentboard-X.Y.Z-windows-amd64-setup.exe` e manifesto. Verifique todos os três arquivos na página Release e o ZIP SHA-256 do manifesto.
4. Em uma máquina Windows limpa com Scoop, execute `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`, depois `agentboard version`, `agentboard init` e `agentboard open` na pasta de teste.
5. Para um feed de atualização permanente, coloque o `agentboard.json` gerado em seu próprio bucket Scoop ou ofereça-o em um bucket público adequado. Volte para `scoop update agentboard` após o próximo lançamento. O comando de instalação de uma URL é adequado para o primeiro contato, enquanto o bucket é mais conveniente para atualizações.

Não substitua manualmente o valor aleatório `hash`: Scoop verifica o conteúdo do ZIP baixado.

Para verificações de manifesto, concentre-se em [formato Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests), [criando manifest](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) e [autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
