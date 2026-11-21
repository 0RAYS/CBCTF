import { useId } from 'react';

/** Shared checkbox box: standalone checkboxes use it, labelled ones use the component below. */
export const CHECKBOX_INPUT_CLASS =
  'w-4 h-4 rounded border-neutral-600/60 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-neutral-800/60';

function Checkbox({ id: suppliedId, checked = false, onChange, label, className = '', disabled = false, ...rest }) {
  const generatedId = useId();
  const id = suppliedId || generatedId;
  return (
    <label htmlFor={id} className={`flex items-center gap-2 text-neutral-300 ${className}`.trim()}>
      <input
        id={id}
        type="checkbox"
        checked={checked}
        onChange={onChange}
        disabled={disabled}
        className={CHECKBOX_INPUT_CLASS}
        {...rest}
      />
      <span>{label}</span>
    </label>
  );
}

export default Checkbox;
