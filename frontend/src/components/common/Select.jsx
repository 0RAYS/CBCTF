import { useId } from 'react';

/** 通用下拉选择组件 */
function Select({
  value,
  onChange,
  options = [],
  placeholder,
  disabled = false,
  error,
  size = 'md',
  fullWidth = true,
  className = '',
  label,
  id: externalId,
  ...rest
}) {
  const generatedId = useId();
  const id = externalId || generatedId;
  const errorId = `${id}-error`;

  // 预定义类定义在 index.css 中
  const sizeClasses = {
    sm: 'select-custom-sm',
    md: 'select-custom-md',
    lg: 'select-custom-lg',
  };

  const selectClasses = `
    select-custom focus-visible:ring-2 focus-visible:ring-geek-400/70
    ${sizeClasses[size] || sizeClasses.md}
    ${fullWidth ? 'w-full' : ''}
    ${error ? '!border-red-400 focus:!border-red-400 focus:!shadow-error' : ''}
    ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
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

      <select
        id={id}
        value={value}
        onChange={onChange}
        disabled={disabled}
        aria-invalid={error ? 'true' : undefined}
        aria-describedby={error ? errorId : undefined}
        className={selectClasses}
        {...rest}
      >
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}

        {options.map((option, index) => (
          <option key={index} value={option.value} disabled={option.disabled}>
            {option.label}
          </option>
        ))}
      </select>

      {error && (
        <div id={errorId} role="alert" className="mt-2 text-sm text-red-400">
          {error}
        </div>
      )}
    </div>
  );
}

export default Select;
