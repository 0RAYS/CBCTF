import { useTranslation } from 'react-i18next';
import { IconMoodEmpty } from '@tabler/icons-react';

/** EmptyState — zero-data placeholder */
function EmptyState({ icon, title, description, hint, action, className = '' }) {
  const { t } = useTranslation();
  const resolvedTitle = title ?? t('common.noData');

  const defaultIcon = <IconMoodEmpty size={40} stroke={1} className="text-neutral-600" />;

  return (
    <div className={`flex flex-col items-center justify-center py-14 px-6 ${className}`}>
      <div className="mb-4 opacity-60">{icon || defaultIcon}</div>
      <p className="text-sm font-mono text-neutral-200 text-center mb-1">{resolvedTitle}</p>
      {description && (
        <p className="text-sm text-neutral-400 text-center max-w-[360px] leading-relaxed mt-1">{description}</p>
      )}
      {hint && <p className="text-xs text-neutral-400 text-center mt-2 font-mono">{hint}</p>}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}

export default EmptyState;
