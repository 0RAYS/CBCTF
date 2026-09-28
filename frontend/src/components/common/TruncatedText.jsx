/** Visual truncation only: keep the full value in the DOM and native tooltip. */
export default function TruncatedText({ children, className = '', maxWidth }) {
  const text = children === null || children === undefined ? '-' : String(children);
  return (
    <span className={`block min-w-0 max-w-full truncate ${className}`} style={{ maxWidth }} title={text}>
      {text}
    </span>
  );
}
