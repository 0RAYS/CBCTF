import 'react';

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
  ...rest
}) {
  // 尺寸样式类（使用index.css中的预定义类）
  const sizeClasses = {
    sm: 'select-custom-sm',
    md: 'select-custom-md',
    lg: 'select-custom-lg',
  };

  // 基础样式（使用index.css中的.select-custom）
  const selectClasses = `
    select-custom
    ${sizeClasses[size] || sizeClasses.md}
    ${fullWidth ? 'w-full' : ''}
    ${error ? '!border-red-400 focus:!border-red-400 focus:!shadow-[0_0_15px_rgba(239,68,68,0.3)]' : ''}
    ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  return (
    <div className={`${fullWidth ? 'w-full' : 'inline-block'}`}>
      {/* 下拉选择框 */}
      <select value={value} onChange={onChange} disabled={disabled} className={selectClasses} {...rest}>
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
      {error && <div className="mt-1 text-sm text-red-400">{error}</div>}
    </div>
  );
}

export default Select;
