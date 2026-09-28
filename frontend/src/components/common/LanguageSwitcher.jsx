import { IconLanguage } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { setLanguage } from '../../i18n';

const LANGUAGES = {
  'zh-CN': { label: '简体中文', code: '中文' },
  en: { label: 'English', code: 'EN' },
};

const SIZES = {
  sm: { button: 'h-9 min-w-20 px-2.5 gap-1.5 text-xs', icon: 16, label: 'w-6' },
  md: { button: 'h-10 min-w-22 px-3 gap-2 text-xs', icon: 16, label: 'w-6' },
  lg: { button: 'h-11 min-w-24 px-4 gap-2 text-sm', icon: 18, label: 'w-7' },
};

function LanguageSwitcher({ size = 'sm', className = '' }) {
  const { i18n, t } = useTranslation();
  const currentLang = i18n.resolvedLanguage in LANGUAGES ? i18n.resolvedLanguage : 'en';
  const nextLang = currentLang === 'en' ? 'zh-CN' : 'en';
  const current = LANGUAGES[currentLang];
  const switchLabel = `${t('common.switchTo')} ${LANGUAGES[nextLang].label}`;
  const styles = SIZES[size] || SIZES.sm;

  return (
    <button
      type="button"
      onClick={() => setLanguage(nextLang)}
      title={switchLabel}
      aria-label={`${t('common.currentLanguage')}: ${current.label}. ${switchLabel}`}
      className={[
        'group inline-flex shrink-0 items-center justify-center rounded-md',
        'border border-neutral-600/60 bg-neutral-800/40',
        'font-mono font-medium text-neutral-300 whitespace-nowrap',
        'transition-colors duration-150',
        'hover:border-geek-400/60 hover:bg-geek-400/10 hover:text-neutral-100',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 focus-visible:ring-offset-2 focus-visible:ring-offset-neutral-900',
        styles.button,
        className,
      ].join(' ')}
    >
      <IconLanguage
        size={styles.icon}
        className="shrink-0 text-neutral-400 transition-colors group-hover:text-geek-400"
        aria-hidden="true"
      />
      <span className="h-3.5 w-px shrink-0 bg-neutral-600/60" aria-hidden="true" />
      <span lang={currentLang} className={`select-none text-center leading-none ${styles.label}`}>
        {current.code}
      </span>
    </button>
  );
}

export default LanguageSwitcher;
