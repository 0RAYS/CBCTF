import { useId } from 'react';

/**
 * 通用下拉选择组件
 * @param {Object} props
 * @param {string|number} props.value - 当前选中值
 * @param {function} props.onChange - 值变化回调
 * @param {Array<{value, label, disabled}>} props.options - 选项列表
 * @param {string} props.placeholder - 占位文本
 * @param {boolean} props.disabled - 是否禁用
 * @param {string} props.error - 错误信息
 * @param {'sm'|'md'|'lg'} props.size - 选择框大小
 * @param {boolean} props.fullWidth - 是否占满宽度
 * @param {string} props.className - 额外的自定义类名
 * @param {string} props.label - 关联标签文本
 * @param {string} props.id - 选择框 id（自动生成可省略）
 */
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

  // 尺寸样式类（使用index.css中的预定义类）
  const sizeClasses = {
    sm: 'select-custom-sm',
    md: 'select-custom-md',
    lg: 'select-custom-lg',
  };

  // 基础样式（使用index.css中的.select-custom）
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
      {/* 关联标签 */}
      {label && (
        <label htmlFor={id} className="block text-sm text-neutral-300 mb-1">
          {label}
        </label>
      )}

      {/* 下拉选择框 */}
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
        {/* 占位选项 */}
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}

        {/* 选项列表 */}
        {options.map((option, index) => (
          <option key={index} value={option.value} disabled={option.disabled}>
            {option.label}
          </option>
        ))}
      </select>

      {/* 错误信息 */}
      {error && (
        <div id={errorId} role="alert" className="mt-1 text-sm text-red-400">
          {error}
        </div>
      )}
    </div>
  );
}

export default Select;
