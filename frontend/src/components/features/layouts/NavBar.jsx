import { useEffect, useRef, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { motion, AnimatePresence } from 'motion/react';
import { Button, LanguageSwitcher, Avatar } from '../../common';
import { useTranslation } from 'react-i18next';
import { EASE_T2, EASE_T3 } from '../../../config/motion';

function NavBar({
  tabs = [],
  activeTab,
  onTabChange,
  logo = 'DEEP DIVE',
  pictureSrc = '',
  userName = '',
  className = '',
  onLogoClick = () => {},
  onPictureClick = () => {},
}) {
  const { t } = useTranslation();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const menuButtonRef = useRef(null);
  const location = useLocation();

  useEffect(() => {
    setMobileMenuOpen(false);
  }, [location, activeTab]);

  useEffect(() => {
    if (!mobileMenuOpen) return;
    const handleKeyDown = (event) => {
      if (event.key === 'Escape') {
        event.preventDefault();
        setMobileMenuOpen(false);
        menuButtonRef.current?.focus();
      }
    };
    const desktop = window.matchMedia('(min-width: 1280px)');
    const closeOnDesktop = () => {
      if (desktop.matches) setMobileMenuOpen(false);
    };
    closeOnDesktop();
    document.addEventListener('keydown', handleKeyDown);
    desktop.addEventListener('change', closeOnDesktop);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      desktop.removeEventListener('change', closeOnDesktop);
    };
  }, [mobileMenuOpen]);

  const handleTabChange = (id) => {
    onTabChange(id);
    setMobileMenuOpen(false);
  };

  return (
    <>
      <div
        className={`fixed top-0 left-0 w-full h-[80px] flex items-center justify-between gap-3 px-3 sm:px-6 lg:px-8 bg-neutral-900/70 backdrop-blur-[4px] border-b border-neutral-700/60 z-40 ${className}`}
      >
        {/* 左侧区域包装 */}
        <div className="flex min-w-0 flex-1 items-center">
          {/* Logo区域 */}
          <div className="relative min-w-0 max-w-[240px] xl:max-w-[200px]">
            <Button
              variant="outline"
              size="lg"
              className="w-full min-w-0 h-[50px] px-3! sm:px-4! font-mono text-sm! sm:text-lg! tracking-wider [&>span]:min-w-0"
              title={logo}
              onClick={() => {
                setMobileMenuOpen(false);
                onLogoClick();
              }}
            >
              <span className="block truncate">{logo}</span>
            </Button>
            {/* Logo装饰角 */}
            <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
            <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
          </div>

          {/* Keep the six contest tabs collapsed on tablets. */}
          <nav aria-label={t('common.mainNavigation')} className="hidden xl:flex shrink-0 ml-6 gap-2">
            {tabs.map((tab) => (
              <Button
                key={tab.id}
                variant={activeTab === tab.id ? 'primary' : 'outline'}
                size="sm"
                className="min-w-0 px-3! whitespace-nowrap"
                aria-current={activeTab === tab.id ? 'page' : undefined}
                onClick={() => onTabChange(tab.id)}
              >
                {tab.label}
              </Button>
            ))}
          </nav>
        </div>

        {/* 右侧语言切换与头像区域 */}
        <div className="flex shrink-0 items-center gap-2 sm:gap-3">
          <LanguageSwitcher size="sm" className="min-h-10" />
          <button
            type="button"
            aria-label={t('common.openUserMenu')}
            className="relative group focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 rounded-md"
            onClick={() => {
              setMobileMenuOpen(false);
              onPictureClick();
            }}
          >
            <div className="w-[45px] h-[45px] border border-neutral-300 rounded-md overflow-hidden transition-colors duration-200 hover:border-neutral-100 cursor-pointer">
              <Avatar src={pictureSrc} name={userName} size={45} shape="rounded" />
            </div>
            {/* 装饰角 */}
            <div className="absolute -top-[2px] -right-[2px] w-[8px] h-[8px] border-t border-r border-neutral-300 group-hover:border-neutral-100 rounded-tr-none"></div>
            <div className="absolute -bottom-[2px] -left-[2px] w-[8px] h-[8px] border-b border-l border-neutral-300 group-hover:border-neutral-100 rounded-bl-none"></div>
          </button>

          {/* Collapsed navigation trigger. */}
          {tabs.length > 0 && (
            <button
              type="button"
              ref={menuButtonRef}
              aria-label={t('common.toggleMenu')}
              aria-expanded={mobileMenuOpen}
              aria-controls="mobile-nav-menu"
              className="xl:hidden shrink-0 w-10 h-10 border border-neutral-300/40 rounded-md flex items-center justify-center text-neutral-300 hover:text-neutral-100 hover:border-neutral-100 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 relative overflow-hidden"
              onClick={() => setMobileMenuOpen((prev) => !prev)}
            >
              <AnimatePresence mode="wait" initial={false}>
                {mobileMenuOpen ? (
                  <motion.svg
                    key="close"
                    viewBox="0 0 24 24"
                    className="w-5 h-5 absolute"
                    fill="none"
                    initial={{ opacity: 0, rotate: -45 }}
                    animate={{ opacity: 1, rotate: 0 }}
                    exit={{ opacity: 0, rotate: 45 }}
                    transition={{ duration: 0.18, ease: EASE_T3 }}
                  >
                    <path d="M6 6L18 18M6 18L18 6" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
                  </motion.svg>
                ) : (
                  <motion.svg
                    key="open"
                    viewBox="0 0 24 24"
                    className="w-5 h-5 absolute"
                    fill="none"
                    initial={{ opacity: 0, rotate: 45 }}
                    animate={{ opacity: 1, rotate: 0 }}
                    exit={{ opacity: 0, rotate: -45 }}
                    transition={{ duration: 0.18, ease: EASE_T3 }}
                  >
                    <path d="M4 6H20M4 12H20M4 18H20" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
                  </motion.svg>
                )}
              </AnimatePresence>
            </button>
          )}
        </div>
      </div>

      {/* 移动端导航菜单 */}
      <AnimatePresence>
        {mobileMenuOpen && tabs.length > 0 && (
          <>
            <motion.div
              className="fixed inset-0 bg-neutral-900/70 z-30 xl:hidden"
              onClick={() => setMobileMenuOpen(false)}
              aria-hidden="true"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.18, ease: EASE_T3 }}
            />
            <motion.nav
              id="mobile-nav-menu"
              aria-label={t('common.mainNavigation')}
              className="fixed top-[80px] left-0 right-0 max-h-[calc(100dvh-80px)] overflow-y-auto z-40 bg-neutral-800/95 border-b border-neutral-600/50 xl:hidden"
              initial={{ opacity: 0, y: -8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
              transition={{ duration: 0.22, ease: EASE_T2 }}
            >
              <div className="flex flex-col p-4 gap-2">
                {tabs.map((tab) => (
                  <Button
                    key={tab.id}
                    variant={activeTab === tab.id ? 'primary' : 'outline'}
                    size="sm"
                    fullWidth
                    aria-current={activeTab === tab.id ? 'page' : undefined}
                    onClick={() => handleTabChange(tab.id)}
                  >
                    {tab.label}
                  </Button>
                ))}
              </div>
            </motion.nav>
          </>
        )}
      </AnimatePresence>
    </>
  );
}

export default NavBar;
