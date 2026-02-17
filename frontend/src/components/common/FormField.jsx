/**
 * 表单字段包装组件
 * @param {Object} props
 * @param {string} props.label - 字段标签
 * @param {boolean} props.required - 是否必填
 * @param {React.ReactNode} props.children - 表单控件（Input、Select、Textarea等）
 * @param {string} props.className - 额外的自定义类名
 * @param {'default'|'subtle'|'mono'} [props.variant='default'] - 标签风格
 */

const labelStyles = {
  default: 'text-neutral-300 text-sm font-medium mb-2',
  subtle: 'text-neutral-400 text-sm mb-1',
  mono: 'text-neutral-400 text-sm font-mono mb-2',
};

function FormField({ label, required = false, children, className = '', variant = 'default' }) {
  return (
    <div className={className}>
      {label && (
        <label className={`block ${labelStyles[variant] || labelStyles.default}`}>
          {label}
          {required && <span className="text-red-400 ml-1">*</span>}
        </label>
      )}
      {children}
    </div>
  );
}

export default FormField;
