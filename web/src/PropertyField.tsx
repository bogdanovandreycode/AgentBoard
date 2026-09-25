import { InputText } from "primereact/inputtext";
import { InputTextarea } from "primereact/inputtextarea";
import { SelectField } from "./SelectField";
import { propertyOptions, type PropertyDef } from "./propertyUtils";

export function PropertyField({ definition, value, onChange }: {
  definition: PropertyDef;
  value: string;
  onChange: (value: string) => void;
}) {
  const options = propertyOptions(definition);
  if (definition.Type === "boolean") return <input type="checkbox" checked={value === "true"} onChange={(e) => onChange(String(e.target.checked))} />;
  if (definition.Type === "select") return <SelectField value={value} onChange={(e) => onChange(e.target.value)}>
    <option value="">{definition.Placeholder || "—"}</option>
    {options.map((option) => <option value={option} key={option}>{option}</option>)}
  </SelectField>;
  if (definition.Type === "multi_select") {
    let selected: string[] = [];
    try { const parsed: unknown = JSON.parse(value || "[]"); if (Array.isArray(parsed)) selected = parsed.filter((item): item is string => typeof item === "string"); } catch { /* invalid saved data stays editable */ }
    return <SelectField multiple value={selected} onChange={(e) => onChange(JSON.stringify(Array.from(e.target.selectedOptions, (option) => option.value)))}>
      {options.map((option) => <option value={option} key={option}>{option}</option>)}
    </SelectField>;
  }
  if (definition.Type === "text") return <InputTextarea value={value} placeholder={definition.Placeholder} onChange={(e) => onChange(e.target.value)} rows={2} />;
  const type = definition.Type === "datetime" ? "datetime-local" : definition.Type;
  const shown = definition.Type === "datetime" && value ? new Date(value).toLocaleString("sv-SE").replace(" ", "T").slice(0, 16) : value;
  return <InputText type={type} step={definition.Type === "number" ? "any" : undefined} value={shown}
    placeholder={definition.Placeholder} pattern={definition.Regex || undefined}
    onChange={(e) => onChange(definition.Type === "datetime" && e.target.value ? new Date(e.target.value).toISOString() : e.target.value)} />;
}
