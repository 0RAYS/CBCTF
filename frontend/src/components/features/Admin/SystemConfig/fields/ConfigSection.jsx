/** Wrapper component for configuration sections. */
export function ConfigSection({ title, icon: Icon, children }) {
  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2">
        <Icon size={18} className="text-geek-400" />
        <h2 className="text-sm font-mono text-neutral-100">{title}</h2>
      </div>
      <div className="space-y-2">{children}</div>
    </div>
  );
}
