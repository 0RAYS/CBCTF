import { createContext, useCallback, useContext, useEffect, useState } from 'react';
import ToastContainer from './ToastContainer';
import { registerToastHandler, unregisterToastHandler } from '../../utils/toast';

const ToastContext = createContext(null);

let toastIdCounter = 0;

export const ToastProvider = ({ children, position = 'top-right', maxToasts = 5 }) => {
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
        return [{ id, title, description, color, timeout, hasCloseButton }, ...prevToasts].slice(0, maxToasts);
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

  // 清除所有toast
  const clearToasts = useCallback(() => {
    setToasts([]);
  }, []);

  // Register handler for the singleton toast module
  useEffect(() => {
    registerToastHandler(addToast);
    return () => unregisterToastHandler();
  }, [addToast]);

  // 提供给外部应用使用的上下文值
  const value = {
    addToast,
    removeToast,
    clearToasts,
  };

  return (
    <ToastContext.Provider value={value}>
      {children}
      <ToastContainer position={position} toasts={toasts} removeToast={removeToast} />
    </ToastContext.Provider>
  );
};

ToastProvider.displayName = 'ToastProvider';

// 创建Hook以便在组件中使用
// eslint-disable-next-line react-refresh/only-export-components
export const useToast = () => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

export default ToastProvider;
