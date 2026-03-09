import 'react';
import { useTranslation } from 'react-i18next';

/**
 * 空状态展示组件
 * @param {Object} props
 * @param {React.ReactNode} props.icon - 自定义图标（默认使用档案图标）
 * @param {string} props.title - 主标题
 * @param {string} props.description - 副标题/描述
 * @param {React.ReactNode} props.action - 操作按钮或其他交互元素
 * @param {string} props.className - 额外的自定义类名
 */
function EmptyState({ icon, title, description, action, className = '' }) {
  const { t } = useTranslation();
  const resolvedTitle = title ?? t('common.noData');

  const defaultIcon = (
    <svg className="w-14 h-14 text-neutral-300/40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={1}
        d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
      />
    </svg>
  );

  return (
    <div className={`flex flex-col items-center justify-center space-y-3 py-16 px-4 ${className}`}>
      <div className="flex items-center justify-center mb-1">{icon || defaultIcon}</div>
      <span className="text-lg font-mono text-neutral-300">{resolvedTitle}</span>
      {description && (
        <span className="text-sm text-neutral-400 text-center max-w-[360px] leading-relaxed">{description}</span>
      )}
      {action && <div className="mt-2">{action}</div>}
    </div>
  );
}

export default EmptyState;
