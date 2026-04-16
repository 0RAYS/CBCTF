import { useCallback, useEffect, useState } from 'react';
import ToastContainer from './ToastContainer';
import { registerToastHandler, unregisterToastHandler } from '../../utils/toast';

let toastIdCounter = 0;

const ToastProvider = ({ children, position = 'top-right', maxToasts = 5 }) => {
  const [toasts, setToasts] = useState([]);

  // 移除toast
  const removeToast = useCallback((id) => {
    setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id));
  }, []);

  // 添加toast
  const addToast = useCallback(
    ({ title, description, color = 'default', timeout = 3000, hasCloseButton = true }) => {
      const id = `toast-${++toastIdCounter}`;
      setToasts((prevToasts) => {
        for (const toast of prevToasts) {
          if (toast.description === description) {
            return prevToasts;
          }
        }
        // FIFO: 若已满, 移除最旧的（末尾）, 再添加新通知
        const existing = prevToasts.length >= maxToasts ? prevToasts.slice(0, maxToasts - 1) : prevToasts;
        return [{ id, title, description, color, timeout, hasCloseButton }, ...existing];
      });

      // 设置自动关闭
      if (timeout > 0) {
        setTimeout(() => {
          removeToast(id);
        }, timeout);
      }

      return id;
    },
    [maxToasts, removeToast]
  );

  // Register handler for the singleton toast module
  useEffect(() => {
    registerToastHandler(addToast);
    return () => unregisterToastHandler();
  }, [addToast]);

  return (
    <>
      {children}
      <ToastContainer position={position} toasts={toasts} removeToast={removeToast} />
    </>
  );
};

ToastProvider.displayName = 'ToastProvider';

export default ToastProvider;
