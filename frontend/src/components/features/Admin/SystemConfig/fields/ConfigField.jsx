import { Input, Select } from '../../../../common';

/**
 * ConfigField - Generic field component for config forms
 * @param {string} label - Field label
 * @param {*} value - Field value
 * @param {Function} onChange - Change handler
 * @param {string} type - Field type: 'text', 'number', 'password', 'select', 'boolean'
 * @param {Array} options - Options for select type (array of {value, label} or strings)
 * @param {string} placeholder - Input placeholder
 */
export function ConfigField({ label, value, onChange, type = 'text', options = [], placeholder = '' }) {
  const handleChange = (event) => {
    const newValue = event.target.value;
    onChange(newValue);
  };

  return (
    <div className="space-y-1">
      <span className="text-xs font-mono text-neutral-400">{label}</span>
      {type === 'select' || type === 'boolean' ? (
        <Select
          size="sm"
          value={type === 'boolean' ? (value ? 'true' : 'false') : value}
          onChange={(event) => {
            const val = event.target.value;
            onChange(type === 'boolean' ? val === 'true' : val);
          }}
          options={options.map((opt) => (typeof opt === 'string' ? { value: opt, label: opt } : opt))}
        />
      ) : (
        <Input size="sm" type={type} value={value ?? ''} onChange={handleChange} placeholder={placeholder} />
      )}
    </div>
  );
}
