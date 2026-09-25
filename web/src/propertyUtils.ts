export type PropertyDef = {
  ID: string;
  Name: string;
  Type: string;
  Options: string;
  Visibility: string;
  Placeholder: string;
  Regex: string;
  DefaultValue: string;
};

export function propertyOptions(definition: PropertyDef): string[] {
  try {
    const value: unknown = JSON.parse(definition.Options || "[]");
    return Array.isArray(value) ? value.filter((item): item is string => typeof item === "string") : [];
  } catch { return []; }
}

export function propertyDisplay(definition: Pick<PropertyDef, "Type">, value: string): string {
  if (definition.Type === "boolean") return value === "true" ? "✓" : "✗";
  if (definition.Type === "multi_select") {
    try { return (JSON.parse(value) as string[]).join(", "); } catch { return value; }
  }
  return value;
}
