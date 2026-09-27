# التثبيت والتشغيل على نظام التشغيل Windows

## المثبت (مستحسن)

قم بتنزيل `agentboard-VERSION-windows-amd64-setup.exe` من الصفحة [Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). سيقترح معالج التثبيت مجلدًا؛ الافتراضي هو `C:\AI\AgentBoard`. سيقوم بنسخ `agentboard.exe` والوثائق وإنشاء اختصار وإضافة المجلد المحدد إلى نظام `PATH`. بعد التثبيت، افتح محطة جديدة حتى يصبح الأمر `agentboard` متاحًا.

في مجلد المشروع الخاص بك، قم بتشغيل:

```powershell
agentboard init
agentboard open
```

الواجهة مدمجة في `agentboard.exe`؛ ليست هناك حاجة إلى تثبيت منفصل لـ Go أو Node.js. يتم تخزين البيانات في `%AppData%\AgentBoard` ويتم الاحتفاظ بها عند تحديث البرنامج أو إلغاء تثبيته. يؤدي إلغاء التثبيت عبر "التطبيقات المثبتة" إلى إزالة الاختصارات والإدخال من `PATH`.

## تثبيت المجرفة

إذا لم يكن Scoop مثبتًا بالفعل، فافتح PowerShell كمستخدم عادي واتبع [تعليمات Scoop](https://scoop.sh/) الرسمية. إذا كانت هناك قيود على جهاز الكمبيوتر الخاص بشركتك، فاتصل بالمسؤول الخاص بك؛ يمكن أيضًا إطلاق AgentBoard من ملف ZIP بدون Scoop.

بدلاً من ذلك، قم بتثبيت AgentBoard عبر Scoop:

```powershell
scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json
agentboard version
```

يمكن التحقق من توفر البيان في إصدار GitHub على [page Releases](https://github.com/bogdanovandreycode/AgentBoard/releases). بعد إضافة البيان إلى Scoop، يمكن تثبيت الحاوية حسب اسم الحاوية وتحديثها باستخدام الأمر `scoop update agentboard`.

## ملف مضغوط بدون مغرفة

قم بتنزيل `agentboard-VERSION-windows-amd64.zip` من الإصدارات، وقم بفك ضغطه، على سبيل المثال، في `C:\Tools\AgentBoard`. في PowerShell في مجلد المشروع:

```powershell
& 'C:\Tools\AgentBoard\agentboard.exe' init
& 'C:\Tools\AgentBoard\agentboard.exe' open
```

بالنسبة لعميل AI، حدد المسار الكامل إلى `agentboard.exe` في تكوين MCP الخاص به إذا لم يكن البرنامج موجودًا في `PATH`.

## البناء من المصدر

قم بتثبيت إصدار Go من `go.mod` وNode.js 22 أو الأحدث. في PowerShell في جذر المستودع:

```powershell
cd web
npm ci
cd ..
./scripts/build.ps1
./agentboard.exe version
```

يقوم `scripts/build.ps1` بتجميع الواجهة الأمامية في `internal/webui/dist`، وإجراء اختبارات Go وتجميع `agentboard.exe` واحد. الترتيب مهم: الواجهة مدمجة في الملف الثنائي عند إنشاء Go. إذا كان `agentboard.exe open` القديم قيد التشغيل، فأوقفه قبل إعادة البناء (Ctrl+C)، وإلا فلن يسمح لك Windows باستبدال الملف.

## أين البيانات

- `%AppData%\AgentBoard\agentboard.db` - المهام والمشاريع والعاملين والإعدادات. يمكنك تحديد ملف مختلف باستخدام العلامة `--db`، ولكن يجب أن يكون لكل من `init` و`open`/`serve` و`mcp` **نفس المسار**.
- `<ваш проект>\.agentboard\project.json` — معرف المشروع. هذا الملف لا يحتوي على مهام.
- يستمع الخادم فقط إلى `127.0.0.1:7337` بشكل افتراضي. حدد عنوانًا آخر `--addr` قبل مسار المشروع: `agentboard open --addr 127.0.0.1:7444 C:\Projects\MyProject`.

يفتح `open` صفحة المشروع ويعيد استخدام خادم قيد التشغيل بالفعل على هذا العنوان. إذا كانت هناك عدة مشاريع مفتوحة في متصفح واحد، فقم بتحديدها من خلال القائمة الموجودة على اليسار.
