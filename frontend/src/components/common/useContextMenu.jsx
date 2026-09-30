import { useCallback, useEffect, useId, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ContextMenu from './ContextMenu';

/** Shared row bindings for List and bespoke tables. Refreshing data invalidates the open target. */
export default function useContextMenu(data, getActions, disabled = false) {
  const { t } = useTranslation();
  const hintId = useId();
  const [target, setTarget] = useState(null);
  const trigger = useRef(null);
  const close = useCallback((restoreFocus = true) => {
    if (restoreFocus && trigger.current?.isConnected) trigger.current.focus({ preventScroll: true });
    setTarget(null);
  }, []);
  useEffect(() => {
    setTarget(null);
  }, [data, disabled]);

  const open = (element, item, x, y) => {
    trigger.current = element;
    setTarget({ item, data, position: { x, y } });
  };
  const getRowProps = (item) =>
    !getActions || disabled
      ? {}
      : {
          tabIndex: 0,
          'aria-haspopup': 'menu',
          'aria-describedby': hintId,
          onContextMenu: (event) => {
            event.preventDefault();
            event.stopPropagation();
            const rect = event.currentTarget.getBoundingClientRect();
            open(event.currentTarget, item, event.clientX || rect.left + 24, event.clientY || rect.top + 24);
          },
          onKeyDown: (event) => {
            if (event.key !== 'ContextMenu' && !(event.shiftKey && event.key === 'F10')) return;
            event.preventDefault();
            event.stopPropagation();
            const rect = event.currentTarget.getBoundingClientRect();
            open(event.currentTarget, item, Math.max(8, rect.left + 24), Math.max(8, rect.top + 24));
          },
        };

  return {
    getRowProps,
    hint: getActions && (
      <p id={hintId} className="px-4 py-2 text-xs text-neutral-500">
        {t('common.contextMenu.hint')}
      </p>
    ),
    menu: target && target.data === data && !disabled && getActions && (
      <ContextMenu position={target.position} actions={getActions(target.item)} onClose={close} />
    ),
  };
}
