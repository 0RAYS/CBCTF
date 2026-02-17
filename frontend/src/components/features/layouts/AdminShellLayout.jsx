import { useMemo, useState, useEffect, useCallback } from 'react';
import { useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import AdminSidebar from '../Admin/AdminSidebar';
import AdminTopbar from '../Admin/AdminTopbar';
import AdminGlobalSearch from '../Admin/AdminGlobalSearch';
import Footer from './Footer';

function AdminShellLayout({
  children,
  sections = [],
  logo = '',
  onLogoClick = () => {},
  onNavigate = () => {},
  user = null,
  onAvatarClick = () => {},
  title = '',
  subtitle,
  showSidebar = true,
  showFooter = true,
}) {
  const location = useLocation();
  const { t } = useTranslation();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const activePath = location.pathname;

  // Ctrl+F / Cmd+F opens global search
  const handleKeyDown = useCallback((e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
      e.preventDefault();
      setSearchOpen(true);
    }
  }, []);

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [handleKeyDown]);

  const activeItem = useMemo(() => {
    const items = sections.flatMap((section) => section.items || []);
    const matches = items.filter((item) => activePath.startsWith(item.path));
    return matches.sort((a, b) => b.path.length - a.path.length)[0];
  }, [activePath, sections]);

  const fallbackLabel = useMemo(() => {
    if (activePath.startsWith('/admin/settings')) return t('nav.settings');
    return '';
  }, [activePath, t]);

  const resolvedLogo = logo || t('branding.admin');
  const titleLabel = activeItem?.label || fallbackLabel;
  const resolvedTitle =
    title || (titleLabel ? t('admin.topbar.managementTitle', { label: titleLabel }) : t('admin.title'));
  const resolvedSubtitle = subtitle !== undefined ? subtitle : t('admin.welcome');
  const userName = user?.user?.nickname || user?.user?.username || user?.user?.name || user?.user?.email || 'ADMIN';

  const handleNavigate = (path) => {
    onNavigate(path);
    setSidebarOpen(false);
  };

  const footerConfig = {
    copyright: t('footer.copyright'),
    icp: {
      number: t('footer.icp'),
      link: 'https://beian.miit.gov.cn/',
    },
    links: [
      { label: t('footer.support'), href: '/support', isExternal: false },
      { label: t('footer.contact'), href: '/contact', isExternal: false },
      { label: t('footer.github'), href: 'https://github.com/0RAYS/CBCTF', isExternal: true },
    ],
  };

  return (
    <div className="h-screen w-screen overflow-hidden">
      <div className="fixed top-0 left-0 w-full h-full bg-black" />
      <div className="relative z-10 flex h-full">
        {showSidebar && (
          <AdminSidebar
            sections={sections}
            activePath={activePath}
            onNavigate={handleNavigate}
            logo={resolvedLogo}
            onLogoClick={onLogoClick}
            open={sidebarOpen}
            onClose={() => setSidebarOpen(false)}
          />
        )}
        <div className="flex-1 flex flex-col min-w-0">
          <AdminTopbar
            title={resolvedTitle}
            subtitle={resolvedSubtitle}
            userName={userName}
            pictureSrc={user?.user?.picture}
            onAvatarClick={onAvatarClick}
            onToggleSidebar={() => setSidebarOpen((prev) => !prev)}
            showSidebarToggle={showSidebar}
          />
          <div className={`flex-1 overflow-y-auto px-4 md:px-8 py-6 relative ${showFooter ? 'pb-[90px]' : ''}`}>
            {children}
          </div>
        </div>
      </div>
      {showFooter && (
        <div className="relative z-10">
          <Footer {...footerConfig} />
        </div>
      )}
      <AdminGlobalSearch isOpen={searchOpen} onClose={() => setSearchOpen(false)} />
    </div>
  );
}

export default AdminShellLayout;
