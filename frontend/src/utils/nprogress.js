let nprogressPromise = null;

const loadNProgress = async () => {
  if (!nprogressPromise) {
    nprogressPromise = Promise.all([
      import('nprogress'),
      import('nprogress/nprogress.css'),
      import('./nprogress.css'),
    ]).then(([{ default: NProgress }]) => {
      NProgress.configure({
        showSpinner: false,
        minimum: 0.1,
        trickleSpeed: 200,
        easing: 'ease',
        speed: 500,
      });
      return NProgress;
    });
  }

  return nprogressPromise;
};

// 计数器
let requestCount = 0;

// 开始加载
export const startLoading = () => {
  requestCount++;
  if (requestCount === 1) {
    loadNProgress().then((NProgress) => {
      if (requestCount > 0) {
        NProgress.start();
      }
    });
  }
};

// 结束加载
export const finishLoading = () => {
  requestCount = Math.max(0, requestCount - 1);
  if (requestCount === 0) {
    loadNProgress().then((NProgress) => {
      if (requestCount === 0) {
        NProgress.done();
      }
    });
  }
};

// 强制结束加载
// export const forceFinishLoading = () => {
//   requestCount = 0;
//   NProgress.done(true);
// };
