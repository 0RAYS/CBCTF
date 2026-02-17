import 'react';

/**
 * 状态标签组件
 * @param {Object} props
 * @param {'success'|'warning'|'error'|'info'|'default'} props.type - 标签类型
 * @param {string} props.text - 标签文本
 * @param {string} props.className - 额外的自定义类名
 */
function StatusTag({ type = 'default', text, className = '' }) {
  // 标签类型样式
  const types = {
    success: 'bg-green-400/20 text-green-400',
    warning: 'bg-yellow-400/20 text-yellow-400',
    error: 'bg-red-400/20 text-red-400',
    info: 'bg-geek-400/20 text-geek-400',
    default: 'bg-neutral-400/20 text-neutral-400',
  };

  const tagClasses = `
    px-3 py-1 rounded-full text-xs font-mono inline-block
    ${types[type] || types.default}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  return <span className={tagClasses}>{text}</span>;
}

export default StatusTag;
