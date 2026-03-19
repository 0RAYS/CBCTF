function FormField({ label, className = '', children }) {
  return (
    <div className={className}>
      {label && <label className="block text-sm font-medium text-neutral-400 mb-1">{label}</label>}
      {children}
    </div>
  );
}

export default FormField;
