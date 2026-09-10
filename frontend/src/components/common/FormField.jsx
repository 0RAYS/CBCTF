function FormField({ label, htmlFor, className = '', children }) {
  return (
    <div className={className}>
      {label && (
        <label htmlFor={htmlFor} className="block text-sm font-medium text-neutral-400 mb-1">
          {label}
        </label>
      )}
      {children}
    </div>
  );
}

export default FormField;
