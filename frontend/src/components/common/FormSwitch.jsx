function FormSwitch({ id, checked = false, onChange, label, className = '' }) {
  return (
    <label htmlFor={id} className={`flex items-center gap-2 text-neutral-300 ${className}`.trim()}>
      <input
        id={id}
        type="checkbox"
        checked={checked}
        onChange={onChange}
        className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
      />
      <span>{label}</span>
    </label>
  );
}

export default FormSwitch;
