// Initial/manual loads remain with the caller, along with its in-flight and scope guards.
// Automatic refresh pauses in background tabs and while offline, then catches up once.
export function startVisiblePolling(callback, delay) {
  if (!(delay > 0)) return () => {};
  let timer = null;
  let disposed = false;
  const canPoll = () =>
    (typeof document === 'undefined' || document.visibilityState !== 'hidden') &&
    (typeof navigator === 'undefined' || navigator.onLine !== false);
  const stop = () => {
    clearInterval(timer);
    timer = null;
  };
  const start = () => {
    timer = setInterval(() => {
      if (canPoll()) callback();
    }, delay);
  };
  const update = () => {
    if (disposed) return;
    if (!canPoll()) {
      stop();
    } else if (timer === null) {
      start();
      callback();
    }
  };
  if (canPoll()) start();
  if (typeof document !== 'undefined') document.addEventListener('visibilitychange', update);
  if (typeof window !== 'undefined') {
    window.addEventListener('online', update);
    window.addEventListener('offline', update);
  }
  return () => {
    disposed = true;
    stop();
    if (typeof document !== 'undefined') document.removeEventListener('visibilitychange', update);
    if (typeof window !== 'undefined') {
      window.removeEventListener('online', update);
      window.removeEventListener('offline', update);
    }
  };
}
