// Toast handler registered by ToastProvider
let handler = null;

// Queue for toasts dispatched before ToastProvider mounts
const pendingQueue = [];

// Called by ToastProvider on mount
export const registerToastHandler = (addToast) => {
  handler = addToast;
  // Flush any toasts queued before provider was ready
  while (pendingQueue.length > 0) {
    handler(pendingQueue.shift());
  }
};

// Called by ToastProvider on unmount
export const unregisterToastHandler = () => {
  handler = null;
};

// Core dispatch – buffers if provider not yet mounted
const dispatch = (params) => {
  if (handler) {
    return handler(params);
  }
  pendingQueue.push(params);
  return null;
};

// Shortcut method defaults
const toastMethods = {
  default: { color: 'default', title: 'Info', timeout: 3000 },
  primary: { color: 'primary', title: 'Success', timeout: 3000 },
  secondary: { color: 'secondary', title: 'Success', timeout: 3000 },
  success: { color: 'success', title: 'Success', timeout: 3000 },
  warning: { color: 'warning', title: 'Warning', timeout: 10000 },
  danger: { color: 'danger', title: 'Error', timeout: 10000 },
  info: { color: 'primary', title: 'Info', timeout: 3000 },
  notice: { color: 'success', title: 'Notice', timeout: 0 },
};

// Main toast function
const toast = (params) => dispatch({ color: 'default', ...params });

// Add shortcut methods (.success(), .danger(), etc.)
Object.entries(toastMethods).forEach(([method, defaults]) => {
  toast[method] = (params = {}) =>
    dispatch({
      ...defaults,
      ...params,
      title: params.title || defaults.title,
    });
});

export { toast };
