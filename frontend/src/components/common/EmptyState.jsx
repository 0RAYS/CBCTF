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

  // 默认图标 - 档案/盒子图标
  const defaultIcon = (
    <svg className="w-12 h-12 text-neutral-300/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={1.5}
        d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
      />
    </svg>
  );

  return (
    <div className={`flex flex-col items-center justify-center space-y-3 py-12 ${className}`}>
      {/* 图标 */}
      <div className="flex items-center justify-center">{icon || defaultIcon}</div>

      {/* 主标题 */}
      <span className="font-mono text-neutral-400">{resolvedTitle}</span>

      {/* 描述 */}
      {description && <span className="text-sm text-neutral-500">{description}</span>}

      {/* 操作按钮 */}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}

export default EmptyState;
