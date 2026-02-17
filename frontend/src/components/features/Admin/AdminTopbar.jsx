import LanguageSwitcher from '../../common/LanguageSwitcher';
import Avatar from '../../common/Avatar';
import { useTranslation } from 'react-i18next';

function AdminTopbar({
  title = '',
  subtitle = '',
  userName = '',
  pictureSrc = '',
  onAvatarClick = () => {},
  onToggleSidebar = () => {},
  showSidebarToggle = true,
}) {
  const { t } = useTranslation();
  const resolvedUserName = userName || t('admin.topbar.defaultUser');

  return (
    <div className="sticky top-0 z-20 flex items-center justify-between h-[70px] px-4 md:px-8 bg-black/40 backdrop-blur-[2px] border-b border-neutral-300/30">
      <div className="flex items-center gap-3">
        {showSidebarToggle && (
          <button
            type="button"
            onClick={onToggleSidebar}
            className="md:hidden w-9 h-9 border border-neutral-300/30 rounded-md flex items-center justify-center text-neutral-300 hover:text-neutral-100 hover:border-neutral-100 transition-colors"
            aria-label={t('admin.topbar.toggleSidebar')}
          >
            <svg viewBox="0 0 24 24" className="w-5 h-5" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 6H20" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
              <path d="M4 12H20" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
              <path d="M4 18H20" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
            </svg>
          </button>
        )}
        <div>
          <div className="text-lg font-mono text-neutral-100 tracking-wider">{title}</div>
          {subtitle && <div className="text-xs font-mono text-neutral-500 mt-1">{subtitle}</div>}
        </div>
      </div>

      <div className="flex items-center gap-4">
        <LanguageSwitcher size="sm" />
        <div className="flex items-center gap-3 cursor-pointer" onClick={onAvatarClick} role="button" tabIndex={0}>
          <div className="text-right">
            <div className="text-sm font-mono text-neutral-200">{resolvedUserName}</div>
            <div className="text-xs text-neutral-500">{t('admin.topbar.role')}</div>
          </div>
          <div className="relative group">
            <div className="w-[42px] h-[42px] border border-neutral-300 rounded-md overflow-hidden transition-all duration-200 hover:border-neutral-100 hover:shadow-[0_0_15px_rgba(179,179,179,0.3)]">
              <Avatar src={pictureSrc} name={userName} size={42} shape="rounded" />
            </div>
            <div className="absolute -top-[2px] -right-[2px] w-[8px] h-[8px] border-t border-r border-neutral-300 group-hover:border-neutral-100" />
            <div className="absolute -bottom-[2px] -left-[2px] w-[8px] h-[8px] border-b border-l border-neutral-300 group-hover:border-neutral-100" />
          </div>
        </div>
      </div>
    </div>
  );
}

export default AdminTopbar;
