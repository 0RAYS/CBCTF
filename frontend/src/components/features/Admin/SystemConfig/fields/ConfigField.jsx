import { Input, Select } from '@/components/common';

/** Generic field component for config forms. */
export function ConfigField({
  label,
  value,
  onChange,
  type = 'text',
  options = [],
  placeholder = '',
  disabled = false,
  description,
  min,
  step,
}) {
  const handleChange = (event) => {
    const newValue = event.target.value;
    onChange?.(newValue);
  };

  return (
    <div className="space-y-1 [&_label]:text-xs [&_label]:font-mono [&_label]:text-neutral-400">
      {type === 'select' || type === 'boolean' ? (
        <Select
          label={label}
          size="sm"
          value={type === 'boolean' ? (value ? 'true' : 'false') : value}
          disabled={disabled}
          onChange={(event) => {
            const val = event.target.value;
            onChange?.(type === 'boolean' ? val === 'true' : val);
          }}
          options={options.map((opt) => (typeof opt === 'string' ? { value: opt, label: opt } : opt))}
        />
      ) : (
        <Input
          label={label}
          size="sm"
          type={type}
          value={value ?? ''}
          onChange={handleChange}
          placeholder={placeholder}
          disabled={disabled}
          min={min}
          step={step}
        />
      )}
      {description && <p className="text-xs leading-relaxed text-neutral-500">{description}</p>}
    </div>
  );
}
