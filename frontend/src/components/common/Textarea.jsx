import { useId } from 'react';
import { INPUT_BASE_CLASS, INPUT_BORDER_CLASS } from '@/components/common/Input';

/** 通用文本域组件 */
function Textarea({
  value,
  onChange,
  onBlur,
  placeholder,
  disabled = false,
  error,
  rows = 4,
  resize = 'none',
  fullWidth = true,
  className = '',
  label,
  id: externalId,
  ...rest
}) {
  const generatedId = useId();
  const id = externalId || generatedId;
  const errorId = `${id}-error`;

  const resizeOptions = {
    none: 'resize-none',
    vertical: 'resize-y',
    horizontal: 'resize-x',
    both: 'resize',
  };

  const textareaClasses = `
    ${INPUT_BASE_CLASS}
    px-4 py-2 placeholder-neutral-500
    ${fullWidth ? 'w-full' : ''}
    ${error ? 'border-red-400 focus:border-red-400 focus:shadow-error' : INPUT_BORDER_CLASS}
    ${disabled ? 'opacity-50 cursor-not-allowed bg-neutral-800/30' : ''}
    ${resizeOptions[resize] || resizeOptions.none}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  return (
    <div className={`${fullWidth ? 'w-full' : 'inline-block'}`}>
      {label && (
        <label htmlFor={id} className="block text-sm text-neutral-300 mb-1">
          {label}
        </label>
      )}

      <textarea
        id={id}
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        placeholder={placeholder}
        disabled={disabled}
        rows={rows}
        aria-invalid={error ? 'true' : undefined}
        aria-describedby={error ? errorId : undefined}
        className={textareaClasses}
        {...rest}
      />

      {error && (
        <div id={errorId} role="alert" className="mt-2 text-sm text-red-400">
          {error}
        </div>
      )}
    </div>
  );
}

export default Textarea;
