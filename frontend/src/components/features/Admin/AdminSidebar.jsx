import { useMemo } from 'react';
import { motion } from 'motion/react';
import { Button } from '../../common';
import { useTranslation } from 'react-i18next';

function AdminSidebar({
  sections = [],
  activePath = '',
  onNavigate = () => {},
  logo = '',
  onLogoClick = () => {},
  open = false,
  onClose = () => {},
  footerOffset = 0,
}) {
  const { t } = useTranslation();
  const resolvedLogo = logo || t('branding.admin');

  const activeItemId = useMemo(() => {
    const allItems = sections.flatMap((s) => s.items || []);
    const matches = allItems.filter((item) => activePath.startsWith(item.path));
    const winner = matches.sort((a, b) => b.path.length - a.path.length)[0];
    return winner?.id ?? null;
  }, [activePath, sections]);

  return (
    <>
      {open && (
        <div
          className="fixed inset-0 bg-black/60 z-30 md:hidden"
          onClick={onClose}
          role="button"
          tabIndex={0}
          aria-label={t('admin.sidebar.close')}
        />
      )}
      <aside
        aria-label={t('admin.sidebar.navigation')}
        className={`fixed md:static z-40 top-0 left-0 w-[240px] border-r border-neutral-600/40 bg-neutral-900/80 backdrop-blur-[4px] transition-transform duration-200 flex flex-col ${
          open ? 'translate-x-0' : '-translate-x-full'
        } md:translate-x-0`}
        style={{
          height: footerOffset > 0 ? `calc(100vh - ${footerOffset}px)` : '100vh',
        }}
      >
        <div className="p-4 border-b border-neutral-600/40">
          <div className="relative">
            <Button variant="outline" size="md" className="w-full justify-center" onClick={onLogoClick}>
              {resolvedLogo}
            </Button>
            <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
            <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
          </div>
        </div>

        <div className="flex-1 overflow-y-auto px-3 py-5 space-y-5">
          {sections.map((section) => (
            <div key={section.id}>
              {section.title && (
                <div className="text-[10px] font-mono text-neutral-500 uppercase tracking-[0.22em] mb-2 px-1">
                  {section.title}
                </div>
              )}
              <div className="space-y-1">
                {section.items.map((item) => {
                  const isActive = item.id === activeItemId;
                  return (
                    <motion.button
                      key={item.id}
                      type="button"
                      className={`w-full text-left px-3 py-2.5 border rounded font-mono text-sm transition-all duration-150 ${
                        isActive
                          ? 'border-geek-400/40 text-geek-300 bg-geek-400/10 shadow-glow-primary'
                          : 'border-transparent text-neutral-400 hover:text-neutral-100 hover:border-neutral-600/60 hover:bg-neutral-700/30'
                      }`}
                      onClick={() => onNavigate(item.path)}
                      whileHover={{ x: 2 }}
                      transition={{ duration: 0.12 }}
                    >
                      {item.label}
                    </motion.button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </aside>
    </>
  );
}

export default AdminSidebar;
