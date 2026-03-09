import { useId } from 'react';

/**
 * 通用输入框组件
 * @param {Object} props
 * @param {'text'|'email'|'password'|'number'|'search'} props.type - 输入框类型
 * @param {string} props.value - 输入框值
 * @param {function} props.onChange - 值变化回调
 * @param {function} props.onBlur - 失焦回调
 * @param {string} props.placeholder - 占位文本
 * @param {boolean} props.disabled - 是否禁用
 * @param {string} props.error - 错误信息
 * @param {React.ReactNode} props.icon - 左侧图标
 * @param {React.ReactNode} props.iconRight - 右侧图标
 * @param {'sm'|'md'|'lg'} props.size - 输入框大小
 * @param {boolean} props.fullWidth - 是否占满宽度
 * @param {string} props.className - 额外的自定义类名
 * @param {string} props.label - 关联标签文本
 * @param {string} props.id - 输入框 id（自动生成可省略）
 */
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

  // 尺寸变体
  const sizes = {
    sm: 'h-8 text-sm',
    md: 'h-10',
    lg: 'h-12 text-base',
  };

  // 计算左右padding
  const getPadding = () => {
    const base = size === 'sm' ? 'px-3' : 'px-4';
    if (icon && iconRight) return 'pl-10 pr-10';
    if (icon) return 'pl-10 pr-4';
    if (iconRight) return 'pl-4 pr-10';
    return base;
  };

  // 基础样式
  const inputClasses = `
    bg-black/20 border rounded-md text-neutral-50 placeholder-neutral-500
    focus:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 transition-all duration-200
    ${sizes[size] || sizes.md}
    ${getPadding()}
    ${fullWidth ? 'w-full' : ''}
    ${error ? 'border-red-400 focus:border-red-400 focus:shadow-[0_0_15px_rgba(239,68,68,0.3)]' : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]'}
    ${disabled ? 'opacity-50 cursor-not-allowed bg-black/10' : ''}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  return (
    <div className={`relative ${fullWidth ? 'w-full' : 'inline-block'}`}>
      {/* 关联标签 */}
      {label && (
        <label htmlFor={id} className="block text-sm text-neutral-300 mb-1">
          {label}
        </label>
      )}

      <div className="relative">
        {/* 左侧图标 */}
        {icon && (
          <div className="absolute left-3 top-1/2 transform -translate-y-1/2 text-neutral-400 pointer-events-none">
            {icon}
          </div>
        )}

        {/* 输入框 */}
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

        {/* 右侧图标 */}
        {iconRight && (
          <div className="absolute right-3 top-1/2 transform -translate-y-1/2 text-neutral-400">{iconRight}</div>
        )}
      </div>

      {/* 错误信息 */}
      {error && (
        <div id={errorId} role="alert" className="mt-1 text-sm text-red-400">
          {error}
        </div>
      )}
    </div>
  );
}

export default Input;
