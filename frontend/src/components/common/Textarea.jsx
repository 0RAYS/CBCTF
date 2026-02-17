import 'react';

/**
 * 通用文本域组件
 * @param {Object} props
 * @param {string} props.value - 文本域值
 * @param {function} props.onChange - 值变化回调
 * @param {function} props.onBlur - 失焦回调
 * @param {string} props.placeholder - 占位文本
 * @param {boolean} props.disabled - 是否禁用
 * @param {string} props.error - 错误信息
 * @param {number} props.rows - 行数
 * @param {'none'|'vertical'|'horizontal'|'both'} props.resize - 调整大小方式
 * @param {boolean} props.fullWidth - 是否占满宽度
 * @param {string} props.className - 额外的自定义类名
 */
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
  ...rest
}) {
  // 调整大小方式
  const resizeOptions = {
    none: 'resize-none',
    vertical: 'resize-y',
    horizontal: 'resize-x',
    both: 'resize',
  };

  // 基础样式
  const textareaClasses = `
    bg-black/20 border rounded-md px-4 py-2 text-neutral-50 placeholder-neutral-500
    focus:outline-none transition-all duration-200
    ${fullWidth ? 'w-full' : ''}
    ${error ? 'border-red-400 focus:border-red-400 focus:shadow-[0_0_15px_rgba(239,68,68,0.3)]' : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]'}
    ${disabled ? 'opacity-50 cursor-not-allowed bg-black/10' : ''}
    ${resizeOptions[resize] || resizeOptions.none}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  return (
    <div className={`${fullWidth ? 'w-full' : 'inline-block'}`}>
      {/* 文本域 */}
      <textarea
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        placeholder={placeholder}
        disabled={disabled}
        rows={rows}
        className={textareaClasses}
        {...rest}
      />

      {/* 错误信息 */}
      {error && <div className="mt-1 text-sm text-red-400">{error}</div>}
    </div>
  );
}

export default Textarea;
