import { useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

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
  const insert = (before: string, after = before, placeholder = "text") => {
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
    <div className="markdown-toolbar" role="toolbar" aria-label="Markdown formatting">
      <button type="button" title="Bold" aria-label="Bold" onClick={() => insert("**")}>B</button>
      <button type="button" title="Italic" aria-label="Italic" onClick={() => insert("*")}><i>I</i></button>
      <button type="button" title="Link" aria-label="Link" onClick={() => insert("[", "](https://example.com)", "link text")}>🔗</button>
      <button type="button" title="Heading" aria-label="Heading" onClick={() => insert("## ", "", "Heading")}>H</button>
      <button type="button" title="List" aria-label="List" onClick={() => insert("- ", "", "List item")}>☷</button>
      <button type="button" title="Code" aria-label="Code" onClick={() => insert("`")}>{"</>"}</button>
      <button type="button" className="preview-toggle" aria-pressed={preview} onClick={() => setPreview(!preview)}>{preview ? "Edit" : "Preview"}</button>
    </div>
    {preview ? <div className="markdown-preview"><MarkdownView value={value || "Nothing to preview."} /></div> :
      <textarea ref={textarea} className="p-inputtextarea p-inputtext p-component" rows={rows} value={value} onChange={(event) => onChange(event.target.value)} />}
  </div>;
}
