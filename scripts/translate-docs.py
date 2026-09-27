"""Generate localized Markdown from the Russian documentation.

The Russian source lives in doc/ru. Run from the repository root. Translation
uses the same Google Translate endpoint as the existing UI locale generator.
English is edited separately. Committed output is served statically; builds do
not call the translation API.
"""

import argparse
import json
import re
import time
import urllib.parse
import urllib.request
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
DOC = ROOT / "doc"
SOURCE = DOC / "ru"
LOCALE_SOURCE = (ROOT / "web/src/i18n.ts").read_text(encoding="utf-8")
LANGUAGES = dict(re.findall(r'\["([\w-]+)", "([^"]+)"\]', LOCALE_SOURCE.split("export const languages = [", 1)[1].split("] as const", 1)[0]))
LANGUAGES.pop("system")
FILES = sorted(path.name for path in SOURCE.glob("*.md"))
FENCE = re.compile(r"(^[ \t]*```[^\n]*\n[\s\S]*?^[ \t]*```[^\n]*$)", re.MULTILINE)
INLINE_CODE = re.compile(r"`[^`\n]+`")
LINK = re.compile(r"\]\(([^)]+)\)")


def request_translation(value: str, language: str) -> str:
    protected = {}

    def protect(match: re.Match[str]) -> str:
        marker = f"ZXQ{len(protected):04d}ZXQ"
        protected[marker] = match.group(0)
        return marker

    safe_value = re.sub(r"`[^`\n]+`|\]\([^)]+\)", protect, value)
    parameters = urllib.parse.urlencode({"client": "gtx", "sl": "ru", "tl": language, "dt": "t", "q": safe_value})
    request = urllib.request.Request("https://translate.googleapis.com/translate_a/single?" + parameters, headers={"User-Agent": "Mozilla/5.0"})
    for attempt in range(5):
        try:
            with urllib.request.urlopen(request, timeout=35) as response:
                translated = "".join(part[0] for part in json.load(response)[0]).strip()
                if any(marker not in translated for marker in protected):
                    chunks = re.split(r"(`[^`\n]+`|\]\([^)]+\))", value)
                    return "".join(chunk if index % 2 else request_translation(chunk, language) for index, chunk in enumerate(chunks) if chunk)
                for marker, original in protected.items():
                    translated = translated.replace(marker, original)
                return translated
        except Exception:
            if attempt == 4:
                raise
            time.sleep(2 * (attempt + 1))
    raise RuntimeError("Translation failed")


def restore_markup(source: str, translated: str) -> str:
    original_code = INLINE_CODE.findall(source)
    translated_code = INLINE_CODE.findall(translated)
    if len(original_code) != len(translated_code):
        raise ValueError("Inline code count changed during translation")
    codes = iter(original_code)
    translated = INLINE_CODE.sub(lambda _: next(codes), translated)
    original_links = LINK.findall(source)
    translated_links = LINK.findall(translated)
    if len(original_links) != len(translated_links):
        raise ValueError("Link count changed during translation")
    links = iter(original_links)
    translated = LINK.sub(lambda match: "](" + next(links) + ")", translated)
    translated = re.sub(r"(?m)^((?:#\s*){2,6})([^\n]*)$", lambda match: "#" * match.group(1).count("#") + " " + match.group(2).strip(), translated)
    translated = re.sub(r"(?m)^(#{1,6})(?=\S)", r"\1 ", translated)
    translated = re.sub(r"(?m)^(\d+\.)(?=\S)", r"\1 ", translated)
    return translated


def translate_line_safely(line: str, language: str) -> str:
    """Fallback for translators that corrupt Markdown tokens in a full paragraph."""
    tokens = re.split(r"(`[^`\n]+`|\[[^\]\n]+\]\([^)]+\))", line)
    output = []
    for token in tokens:
        if not token:
            continue
        if token.startswith("`"):
            output.append(token)
        elif token.startswith("[") and "](" in token:
            label, destination = token[1:].split("](", 1)
            output.append("[" + request_translation(label, language) + "](" + destination)
        else:
            output.append(request_translation(token, language))
    return "".join(output)


def translate_prose(text: str, language: str) -> str:
    if not text.strip():
        return text
    paragraphs = re.split(r"\n{2,}", text.strip())
    output = []
    index = 0
    while index < len(paragraphs):
        batch = [paragraphs[index]]
        while index + len(batch) < len(paragraphs) and len("\n\n".join(batch + [paragraphs[index + len(batch)]])) < 1600:
            batch.append(paragraphs[index + len(batch)])
        source = "\n\n".join(batch)
        translated = request_translation(source, language)
        if translated.count("\n\n") != len(batch) - 1:
            translated = "\n\n".join(request_translation(paragraph, language) for paragraph in batch)
        try:
            output.append(restore_markup(source, translated))
        except ValueError:
            output.append("\n".join(translate_line_safely(line, language) for line in source.splitlines()))
        index += len(batch)
    return "\n\n".join(output)


def translate_markdown(source: str, language: str) -> str:
    pieces = FENCE.split(source.replace("\r\n", "\n"))
    translated = []
    for piece in pieces:
        if piece.lstrip().startswith("```"):
            translated.append(piece)
        else:
            translated.append(translate_prose(piece, language))
    result = "\n\n".join(part.strip("\n") for part in translated if part.strip()) + "\n"
    return re.sub(r"(?m)^((?:#\s*){2,6})([^\n]*)$", lambda match: "#" * match.group(1).count("#") + " " + match.group(2).strip(), result)


def destination(language: str, filename: str) -> Path:
    return DOC / language / filename


def process_language(language: str, files: list[str]) -> tuple[str, int]:
    if language in ("ru", "en"):
        return language, len(files)
    count = 0
    for filename in files:
        target = destination(language, filename)
        target.parent.mkdir(parents=True, exist_ok=True)
        source = (SOURCE / filename).read_text(encoding="utf-8")
        if filename == "README.md":
            source = source.replace("[🌐 Languages](../LANGUAGES.md)\n\n", "")
        try:
            translated = translate_markdown(source, language)
        except Exception as error:
            raise RuntimeError(f"{language}/{filename}: {error}") from error
        if filename == "README.md":
            translated = translated.replace("# AgentBoard\n", "# AgentBoard\n\n[🌐 Languages](../LANGUAGES.md)\n", 1)
        target.write_text(translated, encoding="utf-8")
        count += 1
    return language, count


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--languages", nargs="*", default=[language for language in LANGUAGES if language not in ("ru", "en")])
    parser.add_argument("--files", nargs="*", default=FILES)
    args = parser.parse_args()
    for language in args.languages:
        if language not in LANGUAGES:
            parser.error(f"Unsupported language: {language}")
    for filename in args.files:
        if filename not in FILES:
            parser.error(f"Unknown document: {filename}")
    with ThreadPoolExecutor(max_workers=3) as pool:
        futures = [pool.submit(process_language, language, args.files) for language in args.languages]
        for future in as_completed(futures):
            language, count = future.result()
            print(f"{language}: {count} documents", flush=True)
