# اتصال العمال وMCP

## الخطوة 1. أنشئ عاملاً

في AgentBoard، افتح **العمال ← إضافة عامل**. حدد ملف تعريف العميل: Codex، أو Claude Code، أو Gemini CLI، أو Cursor، أو OpenCode، أو Ollama عبر OpenCode، أو VS Code Copilot أو عميل MCP آخر. يملأ الملف التعريفي الميزات النموذجية (الكود، والاختبارات، وGit)؛ يمكنك تغييرها. الاسم مرئي على اللوحة، و`Slug` هو معرف قصير بدون مسافات لأمر MCP. انقر **حفظ**.

عامل واحد يتوافق مع شخصية واحدة من الذكاء الاصطناعي. إنشاء عمال مختلفين لعملاء أو فرق مختلفة. يصف ملف تعريف القدرة التخصص، لكنه لا يوسع حقوق الذكاء الاصطناعي لتشمل مراحل المهمة.

## الخطوة 2. انسخ التكوين

افتح العامل الذي تم إنشاؤه. تعرض كتلة **تشخيصات MCP** جزء التكوين والملف الذي تريد إضافته. انقر **نسخ تكوين MCP**. إذا كان الملف موجودًا بالفعل، أضف الخادم المقترح إلى كائن `mcpServers`/`servers`/`mcp` الموجود دون مسح الخوادم الأخرى.

يبدو الأمر الرئيسي كما يلي:

```text
agentboard mcp --project C:\Projects\MyFirstProject --worker codex
```

يجب أن يشير `--project` إلى المجلد الذي قمت بتسجيله عبر `agentboard init`. `--worker` — `Slug` للعامل الذي تم إنشاؤه. يقوم عميل MCP بتشغيل هذا الأمر بنفسه عندما يحتاج إلى أدوات. في متصفح AgentBoard، يمكن تشغيل خادم الويب بشكل منفصل.

بعد التحقق، تعرض الكتلة التشخيصية المسار المطلق لجهاز `agentboard.exe` قيد التشغيل. إنه مفيد بشكل خاص عند التثبيت من ملف ZIP. عند التثبيت عبر Scoop، يمكنك استخدام الأمر `agentboard` إذا رأى العميل نفس `PATH`.

## الخطوة 3: أضف خادمًا إلى عميلك

تحتوي واجهة العامل بالفعل على جزء جاهز. وفيما يلي شرح لمكان استخدامه:

| العميل | مكان الإدراج | كيفية التحقق من جانب العميل |
| --- | --- | --- |
| الدستور الغذائي | `%USERPROFILE%\.codex\config.toml`، قسم `[mcp_servers.agentboard_<slug>]` | `codex mcp list` |
| كلود كود | `.mcp.json` في مجلد المشروع | `claude mcp list` |
| الجوزاء CLI | `%USERPROFILE%\.gemini\settings.json`، الكائن `mcpServers` | `/mcp list` في الجوزاء CLI |
| المؤشر | مشروع `.cursor\mcp.json` | قائمة خوادم MCP في إعدادات المؤشر |
| الكود المفتوح | مشروع `opencode.json`، الكائن `mcp` | قائمة أدوات MCP في OpenCode |
| مساعد الطيار VS Code | مشروع `.vscode\mcp.json`، الكائن `servers` | الأمر **MCP: خوادم القائمة** |

بالنسبة لـ Codex وClaude Code، مثال مع المشروع `C:\Projects\MyFirstProject` والعامل `codex`:

```toml
# %USERPROFILE%\.codex\config.toml
[mcp_servers.agentboard_codex]
command = "agentboard"
args = ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
```

```json
{
  "mcpServers": {
    "agentboard_codex": {
      "type": "stdio",
      "command": "agentboard",
      "args": ["mcp", "--project", "C:\\Projects\\MyFirstProject", "--worker", "codex"]
    }
  }
}
```

في JSON، يتم مضاعفة الخط المائل العكسي لنظام التشغيل Windows؛ يقوم جزء جاهز من الواجهة بذلك تلقائيًا. إذا كنت تستخدم `--db` مع قاعدة مخصصة، فأضفه إلى تكوين `args` MCP وحدد نفس المسار كما هو الحال عند بدء `open`.

**Ollama** يوفر نموذجًا محليًا، لكنه لا يحل محل عميل MCP. يقوم ملف تعريف **Ollama via OpenCode** بإنشاء تكوين MCP لـ OpenCode؛ قم بتكوين OpenCode بشكل منفصل على نموذج Ollama. يعد عميل آخر متوافق مع MCP مع Ollama مناسبًا أيضًا.

## الخطوة 4: التحقق من اتصالك

1. افتح بطاقة العامل وانقر فوق **التحقق مرة أخرى**. **فحص الخادم** يجب أن يُظهر عدد أدوات MCP. هذا هو فحص البروتوكول الداخلي واكتشاف أدوات الخادم.
2. قم بتشغيل أو إعادة تشغيل عميل AI بعد إضافة ملف التكوين. اطلب منه الاتصال بـ `get_my_board`.
3. سيظهر **العميل المتصل** في AgentBoard وستظهر جلسة جديدة في القائمة. هذا فقط يؤكد اتصال عميلك. إذا كانت هناك أدوات، ولكن العميل غير متصل، فتحقق من المسار إلى البرنامج واسم ملف التكوين وبناء جملة JSON/TOML الخاص به.

ابدأ جلسة عملك مع `get_my_board`. لا يتلقى الذكاء الاصطناعي المهام من `Backlog` و`Complete`، حتى لو كان يعرف هويتهم. لا يمكن للذكاء الاصطناعي أن يتظاهر بأنه إنسان وليس لديه أمر عام "بالتحرك إلى أي مكان". يعتمد العمل مع نظام الملفات خارج AgentBoard على إمكانيات العميل المحدد ووضع الحماية الخاص به.

تعليمات العملاء الرسمية: [Codex](https://developers.openai.com/learn/docs-mcp)، [Claude Code](https://docs.anthropic.com/en/docs/claude-code/mcp)، [Gemini CLI](https://geminicli.com/docs/tools/mcp-server/)، [Cursor](https://prod.cursor.com/docs/cli/mcp)، [OpenCode](https://opencode.ai/docs/mcp-servers/)، [VS Code](https://code.visualstudio.com/docs/agent-customization/mcp-servers).
