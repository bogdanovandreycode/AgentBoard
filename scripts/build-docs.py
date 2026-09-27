"""Build the multilingual Markdown documentation as a static GitHub Pages site."""

import html
import os
import re
import shutil
from pathlib import Path

import markdown


ROOT = Path(__file__).resolve().parents[1]
DOC = ROOT / "doc"
SITE = ROOT / "site"
LOCALE_SOURCE = (ROOT / "web/src/i18n.ts").read_text(encoding="utf-8")
LANGUAGES = dict(re.findall(r'\["([\w-]+)", "([^"]+)"\]', LOCALE_SOURCE.split("export const languages = [", 1)[1].split("] as const", 1)[0]))
LANGUAGES.pop("system")
DOCUMENTS = ["README.md", *sorted(path.name for path in DOC.glob("*.md") if path.name != "LANGUAGES.md")]
LINK = re.compile(r"\]\(([^)]+)\)")


def source_path(language: str, name: str) -> Path:
    if language == "en":
        return ROOT / name if name == "README.md" else DOC / name
    return DOC / language / name


def output_path(language: str, name: str) -> Path:
    folder = SITE if language == "en" else SITE / language
    if name == "README.md":
        return folder / "index.html"
    if language == "en":
        folder /= "doc"
    return folder / (name.removesuffix(".md") + ".html")


def relative_url(current: Path, destination: Path) -> str:
    return Path(os.path.relpath(destination, current.parent)).as_posix()


def markdown_html(text: str, source: Path, output: Path) -> str:
    targets = {source_path(language, name).resolve(): output_path(language, name) for language in LANGUAGES for name in DOCUMENTS}
    targets[(DOC / "LANGUAGES.md").resolve()] = SITE / "doc" / "LANGUAGES.html"
    targets[(ROOT / "LICENSE").resolve()] = SITE / "LICENSE"
    targets[(ROOT / "assets/social-preview.png").resolve()] = SITE / "assets/social-preview.png"

    def replace_link(match: re.Match[str]) -> str:
        destination, separator, fragment = match.group(1).partition("#")
        if "://" in destination or destination.startswith("mailto:"):
            return match.group(0)
        target = targets.get((source.parent / destination).resolve())
        return "](" + (relative_url(output, target) if target else destination) + (separator + fragment if separator else "") + ")"

    text = LINK.sub(replace_link, text)
    return markdown.markdown(text, extensions=["fenced_code", "tables", "toc", "sane_lists"])


def build_page(language: str, name: str) -> None:
    source = source_path(language, name)
    output = output_path(language, name)
    output.parent.mkdir(parents=True, exist_ok=True)
    content = source.read_text(encoding="utf-8")
    title_match = re.search(r"(?m)^# (.+)$", content)
    title = title_match.group(1) if title_match else name.removesuffix(".md")
    links = []
    for document in DOCUMENTS:
        target = source_path(language, document)
        heading = re.search(r"(?m)^# (.+)$", target.read_text(encoding="utf-8"))
        label = heading.group(1) if heading else document.removesuffix(".md")
        selected = ' aria-current="page"' if name == document else ""
        links.append(f'<a href="{html.escape(relative_url(output, output_path(language, document)))}"{selected}>{html.escape(label)}</a>')
    options = []
    for code, label in LANGUAGES.items():
        selected = " selected" if code == language else ""
        destination = output_path(code, name if name != "LANGUAGES.md" else "README.md")
        options.append(f'<option value="{html.escape(relative_url(output, destination))}"{selected}>{html.escape(label)}</option>')
    icon = relative_url(output, SITE / "assets" / "agentboard-icon.png")
    styles = relative_url(output, SITE / "assets" / "site.css")
    page = f"""<!doctype html>
<html lang="{html.escape(language)}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="color-scheme" content="dark light">
  <title>{html.escape(title)} · AgentBoard</title>
  <link rel="icon" type="image/png" href="{html.escape(icon)}">
  <link rel="stylesheet" href="{html.escape(styles)}">
</head>
<body>
  <header class="site-header"><a class="brand" href="{html.escape(relative_url(output, output_path(language, 'README.md')))}"><img src="{html.escape(icon)}" alt="">AgentBoard</a><label class="language">Language <select aria-label="Documentation language" onchange="location.href=this.value">{''.join(options)}</select></label></header>
  <div class="layout"><nav class="sidebar" aria-label="Documentation">{''.join(links)}<a href="{html.escape(relative_url(output, SITE / 'doc' / 'LANGUAGES.html'))}">Languages</a></nav><main class="content">{markdown_html(content, source, output)}</main></div>
</body>
</html>
"""
    output.write_text(page, encoding="utf-8")


def main() -> None:
    SITE.mkdir(exist_ok=True)
    assets = SITE / "assets"
    assets.mkdir(exist_ok=True)
    shutil.copyfile(ROOT / "web/public/agentboard-icon.png", assets / "agentboard-icon.png")
    shutil.copyfile(ROOT / "assets/social-preview.png", assets / "social-preview.png")
    shutil.copyfile(ROOT / "LICENSE", SITE / "LICENSE")
    shutil.copyfile(ROOT / "doc/site.css", assets / "site.css")
    for language in LANGUAGES:
        for name in DOCUMENTS:
            if not source_path(language, name).is_file():
                raise FileNotFoundError(source_path(language, name))
            build_page(language, name)
    build_page("en", "LANGUAGES.md")
    print(f"Built {len(LANGUAGES) * len(DOCUMENTS) + 1} pages in {SITE}")


if __name__ == "__main__":
    main()
