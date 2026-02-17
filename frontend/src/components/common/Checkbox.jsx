/**
 * Styled checkbox with label
 * @param {Object} props
 * @param {string} props.id - Checkbox ID (also used for label's htmlFor)
 * @param {string} props.label - Label text
 * @param {boolean} props.checked - Checked state
 * @param {Function} props.onChange - Change handler (receives event)
 * @param {string} [props.className] - Additional wrapper className
 */
function Checkbox({ id, label, checked, onChange, className = '' }) {
  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <input
        type="checkbox"
        id={id}
        checked={checked}
        onChange={onChange}
        className="w-4 h-4 rounded border-neutral-300/30 text-geek-400
                  focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
      />
      <label htmlFor={id} className="text-sm font-mono text-neutral-400">
        {label}
      </label>
    </div>
  );
}

export default Checkbox;
