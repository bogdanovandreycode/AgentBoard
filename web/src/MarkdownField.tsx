import { useEffect, useLayoutEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Message } from "primereact/message";
import { ImageAttachments, ImageViewer } from "./TaskImages";
import { FileAttachments, FileViewer } from "./TaskFiles";
import { fileReferences, imageReferences, isManagedFileURL, isManagedImageURL, nextImageNumber, removeImageReference, uploadImage, type ImageReference } from "./imageReferences";
import { t } from "./i18n";

type Props = {
  value: string;
  onChange: (value: string) => void;
  rows?: number;
  projectID?: string;
  placeholder?: string;
  disabled?: boolean;
  onBusyChange?: (busy: boolean) => void;
  onUploadError?: (message: string) => void;
};

type MarkdownNode = { type: string; url?: string; alt?: string; value?: string; children?: MarkdownNode[] };
// A Markdown image can itself be wrapped in a link. Keep only its label inside
// that link so the two image renderers cannot produce nested interactive buttons.
function linkedImageLabels() {
  return (tree: MarkdownNode) => {
    const visit = (node: MarkdownNode, insideLink = false) => {
      if (insideLink && node.type === "image" && (isManagedImageURL(node.url) || isManagedFileURL(node.url))) {
        node.type = "text";
        node.value = `${isManagedFileURL(node.url) ? "📄" : "🖼"} ${node.alt || t("Image")}`;
      }
      node.children?.forEach((child) => visit(child, insideLink || node.type === "link"));
    };
    visit(tree);
  };
}

export function MarkdownView({ value }: { value: string }) {
  const [selected, setSelected] = useState<ImageReference | null>(null);
  const [selectedFile, setSelectedFile] = useState<ImageReference | null>(null);
  const images = imageReferences(value);
  const files = fileReferences(value);
  return <div className="markdown-view"><ReactMarkdown remarkPlugins={[remarkGfm, linkedImageLabels]} components={{
    a: ({ href, children, node: _node, ...props }) => isManagedImageURL(href)
      ? <button type="button" className="image-link" onClick={() => setSelected(images.find((image) => image.url === href) || { url: href, label: t("Image") })}>{children}</button>
      : isManagedFileURL(href)
        ? <button type="button" className="image-link" onClick={() => setSelectedFile(files.find((file) => file.url === href) || {url: href, label: t("File")})}>{children}</button>
        : <a {...props} href={href}>{children}</a>,
    img: ({ src, alt, node: _node, ...props }) => isManagedImageURL(typeof src === "string" ? src : undefined)
      ? <button type="button" className="image-link" onClick={() => setSelected({ url: src as string, label: alt || t("Image") })}>🖼 {alt || t("Image")}</button>
      : isManagedFileURL(typeof src === "string" ? src : undefined)
        ? <button type="button" className="image-link" onClick={() => setSelectedFile({url: src as string, label: alt || t("File")})}>📄 {alt || t("File")}</button>
        : <img {...props} src={src} alt={alt} loading="lazy" />,
  }}>{value}</ReactMarkdown>
    <ImageAttachments images={images} open={setSelected} />
    <FileAttachments files={files} open={setSelectedFile} />
    <ImageViewer image={selected} close={() => setSelected(null)} />
    <FileViewer file={selectedFile} close={() => setSelectedFile(null)} />
  </div>;
}

export function MarkdownField({ value, onChange, rows = 5, projectID, placeholder, disabled, onBusyChange, onUploadError }: Props) {
  const textarea = useRef<HTMLTextAreaElement>(null);
  const fileInput = useRef<HTMLInputElement>(null);
  const documentInput = useRef<HTMLInputElement>(null);
  const current = useRef({ value, onChange, onBusyChange, onUploadError });
  useLayoutEffect(() => { current.current = { value, onChange, onBusyChange, onUploadError }; }, [value, onChange, onBusyChange, onUploadError]);
  const uploadController = useRef<AbortController | null>(null);
  const selection = useRef({ start: value.length, end: value.length });
  const [preview, setPreview] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [dragging, setDragging] = useState(false);
  const [selected, setSelected] = useState<ImageReference | null>(null);
  const [selectedFile, setSelectedFile] = useState<ImageReference | null>(null);
  const blocked = disabled || busy;
  useEffect(() => () => { uploadController.current?.abort(); current.current.onBusyChange?.(false); }, []);
  const rememberSelection = () => {
    const field = textarea.current;
    if (field) selection.current = { start: field.selectionStart, end: field.selectionEnd };
  };
  const attach = async (files: File[], kind: "images" | "files" | "auto" = "images") => {
    if (!projectID || disabled || uploadController.current || !files.length) return;
    const controller = new AbortController();
    uploadController.current = controller;
    setBusy(true);
    setError("");
    current.current.onBusyChange?.(true);
    const { start, end } = selection.current;
    let inserted = "";
    let number = nextImageNumber(current.current.value);
    try {
      for (const file of files) {
        try {
          const raster = ["image/png", "image/jpeg", "image/webp", "image/gif", "image/bmp"].includes(file.type);
          const image = kind === "images" || (kind === "auto" && raster);
          const url = await uploadImage(projectID, file, controller.signal, image ? "images" : "files");
          if (controller.signal.aborted) return;
          const name = file.name.replace(/[[\]\\]/g, "\\$&").replace(/[\r\n]/g, " ");
          const reference = image ? `[🖼 ${t("Screenshot")} ${number++}](${url})` : `[📄 ${name}](${url})`;
          const original = current.current.value;
          const prefix = inserted ? " " : (start > 0 && !/\s/.test(original[start - 1]) ? " " : "");
          const suffix = !inserted && end < original.length && !/\s/.test(original[end]) ? " " : "";
          const addition = prefix + reference;
          const offset = start + inserted.length;
          const next = original.slice(0, offset) + addition + (inserted ? "" : suffix) + original.slice(inserted ? offset : end);
          inserted += addition;
          current.current.value = next;
          current.current.onChange(next);
        } catch (cause) {
          if (controller.signal.aborted) return;
          const message = `${file.name}: ${cause instanceof Error ? cause.message : t("Attachment upload failed")}`;
          setError(message);
          current.current.onUploadError?.(message);
        }
      }
    } finally {
      uploadController.current = null;
      if (!controller.signal.aborted) {
        setBusy(false);
        current.current.onBusyChange?.(false);
        requestAnimationFrame(() => {
          textarea.current?.focus();
          textarea.current?.setSelectionRange(start + inserted.length, start + inserted.length);
          rememberSelection();
        });
      }
    }
  };
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
  return <div className={`markdown-field${dragging ? " dragging" : ""}`} onDragOver={(event) => {
    if (projectID && !blocked && event.dataTransfer.types.includes("Files")) { event.preventDefault(); setDragging(true); }
  }} onDragLeave={(event) => { if (!event.currentTarget.contains(event.relatedTarget as Node | null)) setDragging(false); }} onDrop={(event) => {
    setDragging(false);
    if (!projectID || !event.dataTransfer.files.length) return;
    event.preventDefault();
    void attach(Array.from(event.dataTransfer.files), "auto");
  }}>
    <div className="markdown-toolbar" role="toolbar" aria-label={t("Markdown formatting")}>
      <button type="button" disabled={blocked} title={t("Bold")} aria-label={t("Bold")} onClick={() => insert("**")}>B</button>
      <button type="button" disabled={blocked} title={t("Italic")} aria-label={t("Italic")} onClick={() => insert("*")}><i>I</i></button>
      <button type="button" disabled={blocked} title={t("Link")} aria-label={t("Link")} onClick={() => insert("[", "](https://example.com)", t("link text"))}>🔗</button>
      <button type="button" disabled={blocked} title={t("Heading")} aria-label={t("Heading")} onClick={() => insert("## ", "", t("Heading"))}>H</button>
      <button type="button" disabled={blocked} title={t("List")} aria-label={t("List")} onClick={() => insert("- ", "", t("List item"))}>☷</button>
      <button type="button" disabled={blocked} title={t("Code")} aria-label={t("Code")} onClick={() => insert("`")}>{"</>"}</button>
      {projectID && <>
        <button type="button" className="markdown-attachment-button" disabled={blocked} title={t("Attach images")} aria-label={t("Attach images")} onClick={() => fileInput.current?.click()}><i className="pi pi-image" aria-hidden="true" /> <span>{t("Images")}</span></button>
        <button type="button" className="markdown-attachment-button" disabled={blocked} title={t("Attach files")} aria-label={t("Attach files")} onClick={() => documentInput.current?.click()}><i className="pi pi-file" aria-hidden="true" /> <span>{t("Files")}</span></button>
      </>}
      <button type="button" disabled={blocked} className="preview-toggle" aria-pressed={preview} onClick={() => setPreview(!preview)}>{preview ? t("Edit") : t("Preview")}</button>
    </div>
    {projectID && <input ref={fileInput} type="file" accept="image/png,image/jpeg,image/webp,image/gif,image/bmp" multiple hidden onChange={(event) => {
      const files = Array.from(event.target.files || []);
      event.target.value = "";
      void attach(files);
    }} />}
    {projectID && <input ref={documentInput} type="file" multiple hidden onChange={(event) => {
      const files = Array.from(event.target.files || []); event.target.value = ""; void attach(files, "files");
    }} />}
    {preview ? <div className="markdown-preview"><MarkdownView value={value || t("Nothing to preview.")} /></div> :
      <textarea ref={textarea} className="p-inputtextarea p-inputtext p-component" rows={rows} value={value} disabled={blocked}
        placeholder={placeholder} onSelect={rememberSelection} onBlur={rememberSelection} onChange={(event) => onChange(event.target.value)}
        onPaste={(event) => {
          const files = Array.from(event.clipboardData.items).filter((item) => item.kind === "file").map((item) => item.getAsFile()).filter((file): file is File => !!file);
          if (projectID && files.length) { event.preventDefault(); rememberSelection(); void attach(files, "auto"); }
        }} />}
    {!preview && <ImageAttachments images={imageReferences(value)} open={setSelected} disabled={blocked}
      remove={(image) => onChange(removeImageReference(current.current.value, image.url))} />}
    {!preview && <FileAttachments files={fileReferences(value)} open={setSelectedFile} disabled={blocked}
      remove={(file) => onChange(removeImageReference(current.current.value, file.url))} />}
    {busy && <div className="markdown-upload-status" role="status"><i className="pi pi-spin pi-spinner" aria-hidden="true" />{t("Uploading attachments…")}</div>}
    {error && <Message className="markdown-upload-error" severity="error" text={error} />}
    {projectID && <div className="markdown-drop-hint">{t("Drop files here or paste a screenshot")}</div>}
    <ImageViewer image={selected} close={() => setSelected(null)} />
    <FileViewer file={selectedFile} close={() => setSelectedFile(null)} />
  </div>;
}
