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
}) {
  const { t } = useTranslation();
  const resolvedLogo = logo || t('branding.admin');

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
        className={`fixed md:static z-40 top-0 left-0 h-full w-[240px] border-r border-neutral-300/30 bg-black/60 backdrop-blur-[2px] transition-transform duration-200 flex flex-col ${
          open ? 'translate-x-0' : '-translate-x-full'
        } md:translate-x-0`}
      >
        <div className="p-4 border-b border-neutral-300/30">
          <div className="relative">
            <Button variant="outline" size="md" className="w-full justify-center" onClick={onLogoClick}>
              {resolvedLogo}
            </Button>
            <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
            <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
          </div>
        </div>

        <div className="flex-1 overflow-y-auto px-3 py-4 space-y-6">
          {sections.map((section) => (
            <div key={section.id}>
              {section.title && (
                <div className="text-xs font-mono text-neutral-500 uppercase tracking-[0.2em] mb-3">
                  {section.title}
                </div>
              )}
              <div className="space-y-2">
                {section.items.map((item) => {
                  const isActive = activePath.startsWith(item.path);
                  return (
                    <motion.button
                      key={item.id}
                      type="button"
                      className={`w-full text-left px-4 py-3 border rounded-md font-mono text-sm transition-all duration-150 ${
                        isActive
                          ? 'border-geek-400/50 text-geek-300 bg-geek-400/10 shadow-[0_0_12px_rgba(89,126,247,0.15)]'
                          : 'border-neutral-300/20 text-neutral-300 hover:text-neutral-100 hover:border-neutral-300/60 hover:bg-white/5'
                      }`}
                      onClick={() => onNavigate(item.path)}
                      whileHover={{ x: 2 }}
                      whileTap={{ scale: 0.98 }}
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
