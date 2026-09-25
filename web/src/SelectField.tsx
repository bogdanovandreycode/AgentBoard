import { Children, isValidElement, type ChangeEvent, type ReactNode, type SelectHTMLAttributes } from "react";
import { Dropdown } from "primereact/dropdown";
import { MultiSelect } from "primereact/multiselect";

type Props = Omit<SelectHTMLAttributes<HTMLSelectElement>, "onChange"> & {
  children: ReactNode;
  onChange?: (event: ChangeEvent<HTMLSelectElement>) => void;
};

function textFrom(children: ReactNode): string {
  return Children.toArray(children).map((child) => {
    if (typeof child === "string" || typeof child === "number") return String(child);
    return isValidElement<{ children?: ReactNode }>(child) ? textFrom(child.props.children) : "";
  }).join("");
}

function optionsFrom(children: ReactNode): { label: string; value: string }[] {
  return Children.toArray(children).flatMap((child) => {
    if (!isValidElement<{ children?: ReactNode; value?: string }>(child)) return [];
    if (child.type === "option") {
      const label = textFrom(child.props.children);
      return [{ label, value: String(child.props.value ?? label) }];
    }
    return optionsFrom(child.props.children);
  });
}

export function SelectField({ children, value, onChange, multiple, className, disabled, ...props }: Props) {
  const options = optionsFrom(children);
  const change = (next: string | string[]) => {
    onChange?.({
      target: {
        value: next,
        selectedOptions: Array.isArray(next) ? next.map((item) => ({ value: item })) : [],
      },
    } as unknown as ChangeEvent<HTMLSelectElement>);
  };
  if (multiple) {
    return <MultiSelect value={value || []} options={options} onChange={(event) => change(event.value || [])}
      className={className} disabled={disabled} display="chip" filter placeholder="Choose dependencies" />;
  }
  return <Dropdown value={value} options={options} onChange={(event) => change(event.value)}
    className={className} disabled={disabled} aria-label={props["aria-label"]} />;
}
