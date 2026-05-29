import { useTranslation } from 'react-i18next';
import { setLanguage } from '../../i18n';

/* Inline flag SVGs — no external deps, no network calls */

function FlagCN({ className }) {
  return (
    <svg viewBox="0 0 20 14" className={className} aria-hidden="true" xmlns="http://www.w3.org/2000/svg">
      <rect width="20" height="14" fill="#DE2910" />
      {/* large star */}
      <polygon
        points="4,2 4.588,3.809 6.472,3.809 4.944,4.932 5.532,6.742 4,5.618 2.468,6.742 3.056,4.932 1.528,3.809 3.412,3.809"
        fill="#FFDE00"
      />
      {/* small stars */}
      <polygon
        points="8,1 8.294,1.904 9.236,1.904 8.471,2.466 8.764,3.371 8,2.809 7.236,3.371 7.529,2.466 6.764,1.904 7.706,1.904"
        fill="#FFDE00"
      />
      <polygon
        points="10,3 10.294,3.904 11.236,3.904 10.471,4.466 10.764,5.371 10,4.809 9.236,5.371 9.529,4.466 8.764,3.904 9.706,3.904"
        fill="#FFDE00"
      />
      <polygon
        points="10,6 10.294,6.904 11.236,6.904 10.471,7.466 10.764,8.371 10,7.809 9.236,8.371 9.529,7.466 8.764,6.904 9.706,6.904"
        fill="#FFDE00"
      />
      <polygon
        points="8,8 8.294,8.904 9.236,8.904 8.471,9.466 8.764,10.371 8,9.809 7.236,10.371 7.529,9.466 6.764,8.904 7.706,8.904"
        fill="#FFDE00"
      />
    </svg>
  );
}

function FlagUS({ className }) {
  return (
    <svg viewBox="0 0 20 14" className={className} aria-hidden="true" xmlns="http://www.w3.org/2000/svg">
      {/* stripes */}
      <rect width="20" height="14" fill="#B22234" />
      <rect y="1.077" width="20" height="1.077" fill="#fff" />
      <rect y="3.231" width="20" height="1.077" fill="#fff" />
      <rect y="5.385" width="20" height="1.077" fill="#fff" />
      <rect y="7.538" width="20" height="1.077" fill="#fff" />
      <rect y="9.692" width="20" height="1.077" fill="#fff" />
      <rect y="11.846" width="20" height="1.077" fill="#fff" />
      {/* canton */}
      <rect width="8" height="7.538" fill="#3C3B6E" />
      {/* stars — 5 rows of 6 and 4 rows of 5, simplified as dots */}
      {[0, 1, 2, 3, 4, 5].map((i) => (
        <circle key={`a${i}`} cx={0.667 + i * 1.333} cy="0.769" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4].map((i) => (
        <circle key={`b${i}`} cx={1.333 + i * 1.333} cy="1.538" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4, 5].map((i) => (
        <circle key={`c${i}`} cx={0.667 + i * 1.333} cy="2.308" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4].map((i) => (
        <circle key={`d${i}`} cx={1.333 + i * 1.333} cy="3.077" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4, 5].map((i) => (
        <circle key={`e${i}`} cx={0.667 + i * 1.333} cy="3.846" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4].map((i) => (
        <circle key={`f${i}`} cx={1.333 + i * 1.333} cy="4.615" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4, 5].map((i) => (
        <circle key={`g${i}`} cx={0.667 + i * 1.333} cy="5.385" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4].map((i) => (
        <circle key={`h${i}`} cx={1.333 + i * 1.333} cy="6.154" r="0.28" fill="#fff" />
      ))}
      {[0, 1, 2, 3, 4, 5].map((i) => (
        <circle key={`i${i}`} cx={0.667 + i * 1.333} cy="6.923" r="0.28" fill="#fff" />
      ))}
    </svg>
  );
}

const LANG_CONFIG = {
  'zh-CN': {
    Flag: FlagCN,
    label: '中文',
    code: '中',
    next: 'en',
    nextLabel: 'EN',
  },
  en: {
    Flag: FlagUS,
    label: 'English',
    code: 'EN',
    next: 'zh-CN',
    nextLabel: '中',
  },
};

function LanguageSwitcher({ size = 'sm', className = '' }) {
  const { i18n, t } = useTranslation();

  const currentLang = i18n.language in LANG_CONFIG ? i18n.language : 'en';
  const { Flag, label, code, next, nextLabel } = LANG_CONFIG[currentLang];

  const handleToggle = () => {
    setLanguage(next);
  };

  const sizeMap = {
    sm: { btn: 'h-8 px-2 gap-1.5 text-xs', flag: 'w-[18px] h-[12.6px] rounded-[1px]' },
    md: { btn: 'h-9 px-2.5 gap-1.5 text-xs', flag: 'w-5 h-[14px] rounded-[1px]' },
    lg: { btn: 'h-10 px-3 gap-2 text-sm', flag: 'w-[22px] h-[15.4px] rounded-[1px]' },
  };
  const { btn, flag } = sizeMap[size] || sizeMap.sm;

  return (
    <button
      type="button"
      onClick={handleToggle}
      title={`${t('common.switchTo')} ${nextLabel}`}
      aria-label={`${t('common.currentLanguage')}: ${label}. ${t('common.switchTo')} ${nextLabel}`}
      className={[
        'group inline-flex items-center justify-center rounded-md',
        'border border-neutral-300/30 bg-black/40',
        'font-mono font-medium tracking-wide text-neutral-300',
        'transition-all duration-150 ease-out',
        'hover:border-geek-400/50 hover:bg-geek-900/30 hover:text-neutral-100',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/60',
        btn,
        className,
      ].join(' ')}
    >
      <span className={`shrink-0 overflow-hidden shadow-sm ${flag}`}>
        <Flag className="h-full w-full" />
      </span>
      <span className="select-none leading-none">{code}</span>
    </button>
  );
}

export default LanguageSwitcher;
