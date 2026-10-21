import { useState } from "react";
import { Button } from "primereact/button";
import { Dialog } from "primereact/dialog";
import { Message } from "primereact/message";
import type { ImageReference } from "./imageReferences";
import { t } from "./i18n";
import "./TaskImages.css";

export function ImageAttachments({ images, open, remove, disabled }: {
  images: ImageReference[];
  open: (image: ImageReference) => void;
  remove?: (image: ImageReference) => void;
  disabled?: boolean;
}) {
  if (!images.length) return null;
  return <ul className="image-attachments" aria-label={t("Attached images")}>
    {images.map((image) => <li key={image.url}>
      <button type="button" className="image-attachment" onClick={() => open(image)} title={image.label}>
        <img src={image.url} alt="" loading="lazy" />
        <span>{image.label}</span>
      </button>
      {remove && <button type="button" className="image-remove" disabled={disabled}
        title={t("Remove image reference")} aria-label={`${t("Remove image reference")}: ${image.label}`}
        onClick={() => remove(image)}><i className="pi pi-times" aria-hidden="true" /></button>}
    </li>)}
  </ul>;
}

export function ImageViewer({ image, close }: { image: ImageReference | null; close: () => void }) {
  const [revealing, setRevealing] = useState(false);
  const [error, setError] = useState("");
  const [loadError, setLoadError] = useState(false);
  const reveal = async () => {
    if (!image) return;
    setRevealing(true);
    setError("");
    try {
      const response = await fetch(`${image.url}/reveal`, { method: "POST" });
      if (!response.ok) {
        const result = await response.json().catch(() => ({}));
        throw new Error(result.error || response.statusText);
      }
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : t("Could not open image folder"));
    } finally {
      setRevealing(false);
    }
  };
  const onClose = () => { setError(""); setLoadError(false); close(); };
  return <Dialog visible={!!image} onHide={onClose} header={image?.label || t("Image")}
    className="image-viewer" modal draggable={false} dismissableMask blockScroll
    footer={<div className="image-viewer-actions">
      <Button label={t("Download")} icon="pi pi-download" outlined onClick={() => {
        if (!image) return;
        const link = document.createElement("a");
        link.href = `${image.url}?download=1`;
        link.download = "";
        link.click();
      }} />
      <Button label={t("Show in folder")} icon="pi pi-folder-open" outlined loading={revealing} onClick={reveal} />
    </div>}>
    {image && <img className="image-viewer-content" src={image.url} alt={image.label}
      onError={() => setLoadError(true)} onLoad={() => setLoadError(false)} />}
    {loadError && <Message severity="error" text={t("Image could not be loaded")} />}
    {error && <Message severity="error" text={error} />}
  </Dialog>;
}
