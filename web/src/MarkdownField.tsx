import { useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { t } from "./i18n";

type Props = {
  value: string;
  onChange: (value: string) => void;
  rows?: number;
};

export function MarkdownView({ value }: { value: string }) {
  return <div className="markdown-view"><ReactMarkdown remarkPlugins={[remarkGfm]}>{value}</ReactMarkdown></div>;
}

export function MarkdownField({ value, onChange, rows = 5 }: Props) {
  const textarea = useRef<HTMLTextAreaElement>(null);
  const [preview, setPreview] = useState(false);
  const insert = (before: string, after = before, placeholder = t("text")) => {
    const field = textarea.current;
    if (!field) return;
    const start = field.selectionStart;
    const end = field.selectionEnd;
    const selected = value.slice(start, end) || placeholder;
    onChange(value.slice(0, start) + before + selected + after + value.slice(end));
    requestAnimationFrame(() => {
      field.focus();
      field.setSelectionRange(start + before.length, start + before.length + selected.length);
    });
  };
  return <div className="markdown-field">
    <div className="markdown-toolbar" role="toolbar" aria-label={t("Markdown formatting")}>
      <button type="button" title={t("Bold")} aria-label={t("Bold")} onClick={() => insert("**")}>B</button>
      <button type="button" title={t("Italic")} aria-label={t("Italic")} onClick={() => insert("*")}><i>I</i></button>
      <button type="button" title={t("Link")} aria-label={t("Link")} onClick={() => insert("[", "](https://example.com)", t("link text"))}>🔗</button>
      <button type="button" title={t("Heading")} aria-label={t("Heading")} onClick={() => insert("## ", "", t("Heading"))}>H</button>
      <button type="button" title={t("List")} aria-label={t("List")} onClick={() => insert("- ", "", t("List item"))}>☷</button>
      <button type="button" title={t("Code")} aria-label={t("Code")} onClick={() => insert("`")}>{"</>"}</button>
      <button type="button" className="preview-toggle" aria-pressed={preview} onClick={() => setPreview(!preview)}>{preview ? t("Edit") : t("Preview")}</button>
    </div>
    {preview ? <div className="markdown-preview"><MarkdownView value={value || t("Nothing to preview.")} /></div> :
      <textarea ref={textarea} className="p-inputtextarea p-inputtext p-component" rows={rows} value={value} onChange={(event) => onChange(event.target.value)} />}
  </div>;
}
