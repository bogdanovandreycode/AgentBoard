import { useState } from "react";
import { Button } from "primereact/button";
import { Dialog } from "primereact/dialog";
import { Message } from "primereact/message";
import type { ImageReference } from "./imageReferences";
import { t } from "./i18n";
import "./TaskImages.css";

export function FileAttachments({ files, open, remove, disabled }: {
  files: ImageReference[];
  open: (file: ImageReference) => void;
  remove?: (file: ImageReference) => void;
  disabled?: boolean;
}) {
  if (!files.length) return null;
  return <ul className="file-attachments" aria-label={t("Attached files")}>
    {files.map((file) => <li key={file.url}>
      <button type="button" className="file-attachment" onClick={() => open(file)} title={file.label}>
        <i className="pi pi-file" aria-hidden="true" /><span>{file.label}</span>
      </button>
      {remove && <button type="button" className="file-remove" disabled={disabled}
        title={t("Remove file reference")} aria-label={`${t("Remove file reference")}: ${file.label}`}
        onClick={() => remove(file)}><i className="pi pi-times" aria-hidden="true" /></button>}
    </li>)}
  </ul>;
}

export function FileViewer({ file, close }: { file: ImageReference | null; close: () => void }) {
  const [revealing, setRevealing] = useState(false);
  const [error, setError] = useState("");
  const reveal = async () => {
    if (!file) return;
    setRevealing(true);
    setError("");
    try {
      const response = await fetch(`${file.url}/reveal`, { method: "POST" });
      if (!response.ok) {
        const result = await response.json().catch(() => ({}));
        throw new Error(result.error || response.statusText);
      }
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : t("Could not open file folder"));
    } finally {
      setRevealing(false);
    }
  };
  return <Dialog visible={!!file} onHide={() => { setError(""); close(); }} header={file?.label || t("File")}
    className="file-viewer" modal draggable={false} dismissableMask blockScroll
    footer={<div className="image-viewer-actions">
      <Button label={t("Download")} icon="pi pi-download" outlined onClick={() => {
        if (!file) return;
        const link = document.createElement("a"); link.href = file.url; link.download = ""; link.click();
      }} />
      <Button label={t("Show in folder")} icon="pi pi-folder-open" outlined loading={revealing} onClick={reveal} />
    </div>}>
    <div className="file-viewer-name"><i className="pi pi-file" aria-hidden="true" /><span>{file?.label}</span></div>
    {error && <Message severity="error" text={error} />}
  </Dialog>;
}
