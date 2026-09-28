import { useLayoutEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { useTranslation } from 'react-i18next';

/** Shared, viewport-bound menu. Actions: { key, label, icon, onClick, disabled, danger, hidden }. */
export default function ContextMenu({ position, actions, onClose }) {
  const { t } = useTranslation();
  const menuRef = useRef(null);
  const visibleActions = actions.filter((action) => action && !action.hidden);

  useLayoutEffect(() => {
    const menu = menuRef.current;
    const bounds = menu.getBoundingClientRect();
    menu.style.left = `${Math.max(8, Math.min(position.x, window.innerWidth - bounds.width - 8))}px`;
    menu.style.top = `${Math.max(8, Math.min(position.y, window.innerHeight - bounds.height - 8))}px`;
    (menu.querySelector('button:not(:disabled)') || menu).focus({
      preventScroll: true,
    });

    const dismiss = (event) => {
      if (!menu.contains(event.target)) onClose(false);
    };
    const close = () => onClose(false);
    const keydown = (event) => {
      if (event.key === 'Escape' || event.key === 'Tab') {
        event.preventDefault();
        event.stopImmediatePropagation();
        onClose();
        return;
      }
      if (!['ArrowDown', 'ArrowUp', 'Home', 'End'].includes(event.key)) return;
      event.preventDefault();
      event.stopImmediatePropagation();
      const items = [...menu.querySelectorAll('button:not(:disabled)')];
      const current = items.indexOf(document.activeElement);
      const next =
        event.key === 'Home'
          ? 0
          : event.key === 'End'
            ? items.length - 1
            : (current + (event.key === 'ArrowDown' ? 1 : -1) + items.length) % items.length;
      items[next]?.focus();
    };
    document.addEventListener('pointerdown', dismiss, true);
    document.addEventListener('contextmenu', dismiss, true);
    document.addEventListener('scroll', dismiss, true);
    document.addEventListener('keydown', keydown, true);
    window.addEventListener('resize', close);
    window.addEventListener('blur', close);
    return () => {
      document.removeEventListener('pointerdown', dismiss, true);
      document.removeEventListener('contextmenu', dismiss, true);
      document.removeEventListener('scroll', dismiss, true);
      document.removeEventListener('keydown', keydown, true);
      window.removeEventListener('resize', close);
      window.removeEventListener('blur', close);
    };
  }, [position, onClose]);

  return createPortal(
    <div
      ref={menuRef}
      role="menu"
      aria-label={t('common.contextMenu.label')}
      tabIndex={-1}
      className="fixed z-[10000] min-w-44 w-max max-w-[calc(100vw-1rem)] max-h-[calc(100dvh-1rem)] overflow-y-auto rounded-lg border border-neutral-600 bg-neutral-900 p-1.5 shadow-2xl outline-none"
      style={{ left: position.x, top: position.y }}
      onClick={(event) => event.stopPropagation()}
      onContextMenu={(event) => {
        event.preventDefault();
        event.stopPropagation();
      }}
    >
      {visibleActions.length ? (
        visibleActions.map((action) => (
          <button
            key={action.key}
            type="button"
            role="menuitem"
            tabIndex={-1}
            disabled={action.disabled}
            className={`flex w-full items-center gap-3 rounded px-3 py-2 text-left text-sm font-mono whitespace-normal break-words outline-none transition-colors disabled:opacity-40 disabled:cursor-not-allowed ${action.danger ? 'text-red-400 hover:bg-red-400/10 focus:bg-red-400/10' : 'text-neutral-200 hover:bg-geek-400/15 focus:bg-geek-400/15 focus:text-geek-300'}`}
            onClick={() => {
              onClose();
              action.onClick();
            }}
          >
            {action.icon && (
              <span className="shrink-0" aria-hidden="true">
                {action.icon}
              </span>
            )}
            <span>{action.label}</span>
          </button>
        ))
      ) : (
        <p className="px-3 py-2 text-sm text-neutral-500">{t('common.contextMenu.empty')}</p>
      )}
    </div>,
    document.body
  );
}
