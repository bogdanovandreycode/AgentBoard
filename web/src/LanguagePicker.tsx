import { useState } from "react";
import { languages, setLanguage, t } from "./i18n";

export function LanguagePicker({ value, onChange }: { value: string; onChange: (value: string) => void }) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const label = languages.find(([code]) => code === value)?.[1] || value;
  return <div className="language-picker" onBlur={(event) => { if (!event.currentTarget.contains(event.relatedTarget)) setOpen(false); }}>
    <button type="button" className="language-trigger" aria-expanded={open} onClick={() => setOpen(!open)}>{label}<span>⌄</span></button>
    <div className={`language-menu ${open ? "open" : ""}`} aria-hidden={!open}><input tabIndex={open ? 0 : -1} placeholder={t("Search")} value={search} onChange={(e) => setSearch(e.target.value)} />
      <div className="language-options">{languages.filter(([, name]) => name.toLocaleLowerCase().includes(search.toLocaleLowerCase())).map(([code, name]) =>
        <button type="button" tabIndex={open ? 0 : -1} key={code} className={code === value ? "selected" : ""} onClick={() => { setLanguage(code); onChange(code); setOpen(false); setSearch(""); }}>{name}{code === value && <span>✓</span>}</button>
      )}</div></div>
  </div>;
}
