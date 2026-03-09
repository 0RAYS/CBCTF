/**
 * Tabs - simple tab navigation bar.
 * @param {Object} props
 * @param {Array<{key: string, label: import('react').ReactNode, disabled?: boolean}>} props.items
 * @param {string} props.value
 * @param {(key: string) => void} props.onChange
 * @param {'default'|'compact'} [props.variant='default']
 * @param {string} [props.wrapperClassName]
 * @param {string} [props.containerClassName]
 * @param {string} [props.buttonClassName]
 */
function Tabs({
  items = [],
  value,
  onChange,
  variant = 'default',
  wrapperClassName = 'w-full mx-auto mb-6',
  containerClassName = '',
  buttonClassName = '',
}) {
  const variants = {
    default: {
      container: 'flex border-b border-neutral-700',
      button: 'px-6 py-3',
    },
    compact: {
      container: 'flex flex-wrap border-b border-neutral-700',
      button: 'px-4 py-2',
    },
  };

  const resolvedVariant = variants[variant] || variants.default;

  const content = (
    <div className={`${resolvedVariant.container} ${containerClassName}`.trim()}>
      {items.map((item) => (
        <button
          key={item.key}
          type="button"
          disabled={item.disabled}
          className={`${resolvedVariant.button} text-sm font-medium transition-colors
            focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 ${
            value === item.key ? 'text-geek-400 border-b-2 border-geek-400' : 'text-neutral-400 hover:text-neutral-300'
          } ${item.disabled ? 'opacity-50 cursor-not-allowed' : ''} ${buttonClassName}`.trim()}
          onClick={() => {
            if (!item.disabled) onChange?.(item.key);
          }}
        >
          {item.label}
        </button>
      ))}
    </div>
  );

  if (wrapperClassName === null) return content;

  return <div className={wrapperClassName}>{content}</div>;
}

export default Tabs;
