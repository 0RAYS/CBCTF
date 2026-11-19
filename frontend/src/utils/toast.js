let handler = null;

// Queue for toasts dispatched before ToastProvider mounts
const pendingQueue = [];

// Called by ToastProvider on mount
export const registerToastHandler = (addToast) => {
  handler = addToast;
  while (pendingQueue.length > 0) {
    handler(pendingQueue.shift());
  }
};

// Called by ToastProvider on unmount
export const unregisterToastHandler = () => {
  handler = null;
};

const dispatch = (params) => {
  if (handler) {
    return handler(params);
  }
  pendingQueue.push(params);
  return null;
};

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

const toast = (params) => dispatch({ color: 'default', ...params });

Object.entries(toastMethods).forEach(([method, defaults]) => {
  toast[method] = (params = {}) =>
    dispatch({
      ...defaults,
      ...params,
      title: params.title || defaults.title,
    });
});

export { toast };
