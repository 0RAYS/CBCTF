import { useTranslation } from 'react-i18next';
import { setLanguage } from '../../i18n';

function LanguageSwitcher({ size = 'sm', className = '' }) {
  const { i18n, t } = useTranslation();

  const handleToggle = () => {
    const nextLanguage = i18n.language === 'zh-CN' ? 'en' : 'zh-CN';
    setLanguage(nextLanguage);
  };

  const sizeClasses = {
    sm: 'w-8 h-8',
    md: 'w-9 h-9',
    lg: 'w-10 h-10',
  };

  return (
    <button
      type="button"
      onClick={handleToggle}
      aria-label={t('common.language')}
      className={`inline-flex items-center justify-center rounded-md border border-neutral-300/40 bg-black/40 text-neutral-300 transition-colors hover:border-geek-400/60 hover:text-neutral-100 hover:shadow-[0_0_10px_rgba(89,126,247,0.2)] ${sizeClasses[size] || sizeClasses.sm} ${className}`}
    >
      <svg viewBox="0 0 24 24" fill="none" className="h-4 w-4" aria-hidden="true">
        <path
          d="M12 3a9 9 0 1 0 0 18m0-18c2.5 2.2 4 5.5 4 9s-1.5 6.8-4 9m0-18C9.5 5.2 8 8.5 8 12s1.5 6.8 4 9m-8-9h16"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
    </button>
  );
}

export default LanguageSwitcher;
