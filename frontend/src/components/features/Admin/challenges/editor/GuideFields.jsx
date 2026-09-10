import { IconPlus, IconTrash } from '@tabler/icons-react';
import Button from '../../../../common/Button';
import { inputBaseClass, selectClass, textareaClass } from './editorStyles.js';

export function GuideListHeader({ title, addLabel, onAdd, disabled = false }) {
  return (
    <div className="flex justify-between items-center gap-2 pt-2">
      <span className="text-xs font-mono text-neutral-400">{title}</span>
      <Button
        variant="ghost"
        size="sm"
        align="icon-left"
        icon={<IconPlus size={12} />}
        className="!bg-transparent !text-geek-400 hover:!text-geek-300 !text-xs !h-8"
        disabled={disabled}
        onClick={onAdd}
      >
        {addLabel}
      </Button>
    </div>
  );
}

export function IconButton({ onClick }) {
  return (
    <Button
      variant="ghost"
      size="icon"
      className="!bg-transparent !text-red-400 hover:!text-red-300 !w-8 !h-10"
      onClick={onClick}
    >
      <IconTrash size={14} />
    </Button>
  );
}

export function GuideField({ label, children }) {
  return (
    <label className="block space-y-1">
      <span className="text-[11px] font-mono text-neutral-500">{label}</span>
      {children}
    </label>
  );
}

export function GuideErrors({ errors }) {
  if (!errors?.length) return null;
  return (
    <div className="mt-1 space-y-1">
      {errors.map((error, index) => (
        <div key={index} className="text-[11px] font-mono text-red-300">
          {error}
        </div>
      ))}
    </div>
  );
}

export function GuideInput({ label, value, onChange, placeholder, required, errors, multiline = false }) {
  const Input = multiline ? 'textarea' : 'input';
  return (
    <GuideField label={label}>
      <Input
        className={multiline ? textareaClass : inputBaseClass}
        value={value}
        required={required}
        placeholder={placeholder}
        onChange={(e) => onChange(e.target.value)}
      />
      <GuideErrors errors={errors} />
    </GuideField>
  );
}

export function GuideBoolean({ label, value, onChange }) {
  return (
    <GuideField label={label}>
      <select
        className={selectClass}
        value={value ? 'true' : 'false'}
        onChange={(e) => onChange(e.target.value === 'true')}
      >
        <option value="false">false</option>
        <option value="true">true</option>
      </select>
    </GuideField>
  );
}

export function GuideStringList({ label, placeholder, values = [], onChange, addLabel }) {
  return (
    <div className="space-y-2">
      <GuideListHeader title={label} addLabel={addLabel} onAdd={() => onChange([...values, ''])} />
      {values.map((value, index) => (
        <div key={index} className="grid grid-cols-[1fr_32px] gap-2">
          <input
            className={inputBaseClass}
            value={value}
            placeholder={placeholder}
            onChange={(e) => onChange(values.map((item, i) => (i === index ? e.target.value : item)))}
          />
          <IconButton onClick={() => onChange(values.filter((_, i) => i !== index))} />
        </div>
      ))}
    </div>
  );
}
