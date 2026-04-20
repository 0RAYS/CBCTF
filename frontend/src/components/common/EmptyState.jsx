import { useTranslation } from 'react-i18next';

/**
 * EmptyState — zero-data placeholder
 * @param {ReactNode} icon         Custom icon (defaults to inbox slot SVG)
 * @param {string}    title        Primary message
 * @param {string}    description  Secondary explanation
 * @param {string}    hint         Short next-action hint (e.g. "Click + Add to create one")
 * @param {ReactNode} action       CTA button or link
 * @param {string}    className    Extra Tailwind classes
 */
function EmptyState({ icon, title, description, hint, action, className = '' }) {
  const { t } = useTranslation();
  const resolvedTitle = title ?? t('common.noData');

  const defaultIcon = (
    <svg className="w-10 h-10 text-neutral-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1}>
      <rect x="3" y="3" width="18" height="18" rx="2" />
      <path strokeLinecap="round" d="M3 9h18M9 21V9" />
    </svg>
  );

  return (
    <div className={`flex flex-col items-center justify-center py-14 px-6 ${className}`}>
      <div className="mb-4 opacity-60">{icon || defaultIcon}</div>
      <p className="text-sm font-mono text-neutral-300 mb-1">{resolvedTitle}</p>
      {description && (
        <p className="text-xs text-neutral-500 text-center max-w-[320px] leading-relaxed mt-1">{description}</p>
      )}
      {hint && (
        <p className="text-xs text-neutral-600 text-center mt-2 font-mono">{hint}</p>
      )}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}

export default EmptyState;
