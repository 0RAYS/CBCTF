import { useId } from 'react';

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
    bg-black/20 border rounded-md px-4 py-2 text-neutral-50 placeholder-neutral-500
    focus:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 transition-all duration-200
    ${fullWidth ? 'w-full' : ''}
    ${error ? 'border-red-400 focus:border-red-400 focus:shadow-error' : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-focus'}
    ${disabled ? 'opacity-50 cursor-not-allowed bg-black/10' : ''}
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
