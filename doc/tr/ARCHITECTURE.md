# Mimari ve API

AgentBoard, SQLite çalıştıran yerel bir Go işlemidir. İnsanlar için HTTP API, yapay zeka için MCP adaptörü ve web arayüzü ortak hizmet mantığını paylaşır. SQLite - durum kaynağı; ön uç sunucuyu atlamaz. Yapay zeka araçlarının ayrı bir haklar yüzeyi vardır ve MCP'de `Backlog`/`Complete` durumu mevcut değildir. `System`'nun geçmiş kaydı uygulamanın kendisi tarafından oluşturulur.

## Bileşenler

- `cmd/agentboard` - CLI `init`, `open`, `serve`, `mcp`, `version`.
- `internal/core` - etki alanı modelleri ve hatalar.
- `internal/service` - görev geçişleri, izinler, içe aktarma ve ayarlar.
- `internal/persistence` - SQLite ve geçişler.
- `internal/httpapi` - İnsan HTTP API'si.
- `internal/mcpserver` - Ayrı bir çalışan için MCP araçları.
- `web` - React/TypeScript kullanıcı arayüzü; derleme `internal/webui/dist`'da sona erer ve EXE'ye dahil edilir.

## Temel HTTP rotaları

| Yöntem ve yol | Hedef |
| --- | --- |
| `GET /api/health` | Sunucu kontrolü. |
| `GET /api/projects` | Kayıtlı projeler. |
| `GET /api/projects/{id}/board` | Proje panosu. |
| `GET/PUT /api/projects/{id}/settings` | Sütunların ayarları, sırası ve adları. |
| `POST /api/projects/{id}/tasks` | Bir görev oluşturun. |
| `POST /api/projects/{id}/tasks/import` | JSON sürüm 1'i atomik olarak içe aktarın. |
| `GET/PATCH/DELETE /api/tasks/{id}` | Kart, değiştir, sil. |
| `POST /api/tasks/{id}/move` | Bir kişinin hareket etmesi. |
| `GET/POST /api/projects/{id}/workers` | Bir işçinin listesi ve oluşturulması. |
| `GET /api/workers/{id}/mcp/check` | MCP anlaşmasının ve araçlarının dahili testi. |
| `GET /api/workers/{id}/sessions` | Oturum teşhisi. |
| `GET/POST /api/projects/{id}/properties` | Özel özellikler. |

HTTP yerel güvenilir kullanıcı içindir. Kendi kimlik doğrulaması, ağ kısıtlamaları ve HTTPS'si olmadan bir web bağlantı noktasını İnternet'te yayınlamayın. MCP sunucusu, belirli bir proje ve çalışan için `stdio` kullanılarak başlatılır; `get_my_board` ile başlayın. İzinleri AgentBoard içindeki eylemleri kısıtlar ancak AI istemcisinin dosya sistemi sanal alanının yerini almaz.

Kullanıcı sütunları `tasks.state`'dan ayrı olarak depolanır: çekirdek, görevi `backlog`'da böyle bir sütunda tutar ve `tasks.board_column`, İnsan panosundaki yeri belirler. Bu aynı yapay zeka geçiş modelini korur. Bir sütunun kaldırılması `board_column` görevlerini temizler.
