# إعداد إصدار للمغرفة

## ما هو آلي بالفعل

ينفذ `scripts/package-scoop.ps1 -Version X.Y.Z` `npm ci`، وبناء الواجهة الأمامية، `go test ./...`، ويقوم Windows ببناء `agentboard.exe` برقم الإصدار، وينشئ ملف ZIP ويحسب SHA-256 لهذا الرمز البريدي. من **نفس** الملف، يقوم بإنشاء `release/agentboard.json` مع URL، والتجزئة، وترخيص MIT، وCLI-shim، والاختصار، و`checkver`، و`autoupdate`. يحتوي ملف ZIP والمثبت على الملف `LICENSE`. توجد قاعدة البيانات خارج دليل التثبيت، لذا ليست هناك حاجة إلى `persist` في البيان.

يقوم `.github/workflows/release.yml` الموجود على العلامة `vX.Y.Z` بتشغيل نفس الحزمة في Windows runner، ويقوم بإنشاء برنامج تثبيت Inno Setup وإرفاق ملف ZIP والمثبت والبيان بإصدار GitHub. يؤدي تشغيل سير العمل يدويًا إلى إنشاء عنصر للاختبار فقط، دون نشر إصدار.

## ترحيل أمر المشرف

1. تأكد من أن الكود والوثائق ورقم الإصدار والملف `LICENSE` جاهزون.
2. قم بتشغيل `./scripts/package-scoop.ps1 -Version X.Y.Z` محليًا. تحقق من مخرجات `release/agentboard-X.Y.Z-windows-amd64.zip` و`release/agentboard.json` و`agentboard version` بعد التفريغ. لا تقم بتغيير الرمز البريدي بعد حساب التجزئة.
3. قم بإنشاء العلامة `vX.Y.Z` وإرسالها. ستنشر GitHub Actions إصدارًا يتضمن ZIP و`agentboard-X.Y.Z-windows-amd64-setup.exe` والبيان. تحقق من الملفات الثلاثة جميعها في صفحة الإصدار وSHA-256 ZIP من البيان.
4. على جهاز Windows نظيف مزود بـ Scoop، قم بتشغيل `scoop install https://github.com/bogdanovandreycode/AgentBoard/releases/latest/download/agentboard.json`، ثم `agentboard version`، و`agentboard init`، و`agentboard open` في مجلد الاختبار.
5. للحصول على موجز تحديث دائم، ضع `agentboard.json` الذي تم إنشاؤه في حاوية Scoop الخاصة بك أو قم بتقديمه في حاوية عامة مناسبة. تحقق مرة أخرى من `scoop update agentboard` بعد الإصدار التالي. يعد أمر التثبيت من عنوان URL مناسبًا للتعارف الأول، بينما يكون الجرافة أكثر ملاءمة للتحديثات.

لا تقم باستبدال القيمة العشوائية `hash` يدويًا: يتحقق Scoop من محتويات ملف ZIP الذي تم تنزيله.

بالنسبة لعمليات التحقق من البيان، ركز على [تنسيق Scoop](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests)، و[إنشاء البيان](https://github.com/ScoopInstaller/Scoop/wiki/Creating-an-app-manifest) و[autoupdate](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifest-Autoupdate).
