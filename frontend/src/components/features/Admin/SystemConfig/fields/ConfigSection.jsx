/**
 * ConfigSection - Wrapper component for configuration sections
 * @param {string} title - Section title
 * @param {Component} icon - Icon component from @tabler/icons-react
 * @param {ReactNode} children - Section content
 */
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
