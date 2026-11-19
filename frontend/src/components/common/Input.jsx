import { useId } from 'react';

/** 通用输入框组件 */
function Input({
  type = 'text',
  value,
  onChange,
  onBlur,
  placeholder,
  disabled = false,
  error,
  icon,
  iconRight,
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

  const sizes = {
    sm: 'h-8 text-sm',
    md: 'h-10',
    lg: 'h-12 text-base',
  };

  const getPadding = () => {
    const base = size === 'sm' ? 'px-3' : 'px-4';
    if (icon && iconRight) return 'pl-10 pr-10';
    if (icon) return 'pl-10 pr-4';
    if (iconRight) return 'pl-4 pr-10';
    return base;
  };

  const inputClasses = `
    min-w-0 bg-neutral-800/60 border rounded-md text-neutral-50 placeholder-neutral-400 font-mono
    focus:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 transition-all duration-200
    ${type === 'number' ? 'tabular-nums' : ''}
    ${sizes[size] || sizes.md}
    ${getPadding()}
    ${fullWidth ? 'w-full' : ''}
    ${error ? 'border-red-400 focus:border-red-400 focus:shadow-error' : 'border-neutral-600/60 focus:border-geek-400 focus:shadow-focus'}
    ${disabled ? 'opacity-50 cursor-not-allowed bg-neutral-800/30' : ''}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  return (
    <div className={`relative min-w-0 ${fullWidth ? 'w-full' : 'inline-block'}`}>
      {label && (
        <label htmlFor={id} className="block text-sm text-neutral-300 mb-1">
          {label}
        </label>
      )}

      <div className="relative">
        {icon && (
          <div className="absolute left-3 top-1/2 transform -translate-y-1/2 text-neutral-400 pointer-events-none">
            {icon}
          </div>
        )}

        <input
          id={id}
          type={type}
          value={value}
          onChange={onChange}
          onBlur={onBlur}
          placeholder={placeholder}
          disabled={disabled}
          aria-invalid={error ? 'true' : undefined}
          aria-describedby={error ? errorId : undefined}
          className={inputClasses}
          {...rest}
        />

        {iconRight && (
          <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-neutral-400">{iconRight}</div>
        )}
      </div>

      {error && (
        <div id={errorId} role="alert" className="mt-2 text-sm text-red-400">
          {error}
        </div>
      )}
    </div>
  );
}

export default Input;
