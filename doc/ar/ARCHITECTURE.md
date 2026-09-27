# الهندسة المعمارية وواجهة برمجة التطبيقات

AgentBoard هي عملية Go محلية تقوم بتشغيل SQLite. تشترك واجهة برمجة تطبيقات HTTP للبشر ومحول MCP للذكاء الاصطناعي وواجهة الويب في منطق الخدمة المشترك. SQLite - مصدر الحالة؛ الواجهة الأمامية لا تتجاوز الخادم. تحتوي أدوات الذكاء الاصطناعي على سطح حقوق منفصل ولا تتوفر حالة `Backlog`/`Complete` في MCP. يتم إنشاء سجل المحفوظات من `System` بواسطة التطبيق نفسه.

## المكونات

- `cmd/agentboard` - CLI `init`، `open`، `serve`، `mcp`، `version`.
- `internal/core` - نماذج المجال والأخطاء.
- `internal/service` - انتقالات المهام والأذونات والاستيراد والإعدادات.
- `internal/persistence` - SQLite والترحيلات.
- `internal/httpapi` - واجهة برمجة تطبيقات HTTP البشرية.
- `internal/mcpserver` - أدوات MCP لعامل منفصل.
- `web` - واجهة المستخدم React/TypeScript؛ ينتهي التجميع في `internal/webui/dist` ويتم تضمينه في EXE.

## مسارات HTTP الأساسية

| الطريقة والمسار | الوجهة |
| --- | --- |
| `GET /api/health` | فحص الخادم. |
| `GET /api/projects` | المشاريع المسجلة. |
| `GET /api/projects/{id}/board` | مجلس المشروع. |
| `GET/PUT /api/projects/{id}/settings` | الإعدادات والترتيب وأسماء الأعمدة. |
| `POST /api/projects/{id}/tasks` | قم بإنشاء مهمة. |
| `POST /api/projects/{id}/tasks/import` | استيراد إصدار JSON 1 ذريًا. |
| `GET/PATCH/DELETE /api/tasks/{id}` | البطاقة، التغيير، الحذف. |
| `POST /api/tasks/{id}/move` | الانتقال من قبل شخص. |
| `GET/POST /api/projects/{id}/workers` | قائمة وإنشاء عامل. |
| `GET /api/workers/{id}/mcp/check` | الاختبار الداخلي لمصافحة وأدوات MCP. |
| `GET /api/workers/{id}/sessions` | تشخيص الجلسة. |
| `GET/POST /api/projects/{id}/properties` | خصائص مخصصة. |

HTTP مخصص للمستخدم المحلي الموثوق به. لا تنشر منفذ ويب على الإنترنت دون المصادقة الخاصة به وقيود الشبكة وHTTPS. يتم تشغيل خادم MCP باستخدام `stdio` لمشروع وعامل محدد؛ ابدأ بـ `get_my_board`. تقيد أذوناتها الإجراءات داخل AgentBoard، ولكنها لا تحل محل صندوق حماية نظام ملفات عميل الذكاء الاصطناعي.

يتم تخزين أعمدة المستخدم بشكل منفصل عن `tasks.state`: يحتفظ المركز بالمهمة في مثل هذا العمود في `backlog`، ويحدد `tasks.board_column` المكان على اللوحة البشرية. يحافظ هذا على نفس نموذج انتقال الذكاء الاصطناعي. تؤدي إزالة عمود إلى مسح مهام `board_column`.
