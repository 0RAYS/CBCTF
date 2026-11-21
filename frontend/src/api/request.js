import axios from 'axios';
import { toast } from '@/utils/toast';
import i18n from '@/i18n';
import { API_CONFIG } from '@/api/config.js';

let nprogressPromise = null;

function loadNProgress() {
  if (!nprogressPromise) {
    nprogressPromise = import('@/utils/nprogress');
  }
  return nprogressPromise;
}

function finishRequestLoading() {
  loadNProgress().then(({ finishLoading }) => finishLoading());
}

const request = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

let requestCount = 0;

const updateGlobalLoading = (count) => {
  requestCount = Math.max(0, count);
  // 使用动态导入避免循环依赖
  import('@/store').then(({ store }) => {
    import('@/store/app').then(({ setGlobalLoading }) => {
      store.dispatch(setGlobalLoading(requestCount > 0));
    });
  });
};

request.interceptors.request.use(
  async (config) => {
    if (!config.noLoading) {
      loadNProgress().then(({ startLoading }) => startLoading());
      updateGlobalLoading(requestCount + 1);
    }
    config.headers['Accept-Language'] = i18n.language;
    return config;
  },
  (error) => {
    finishRequestLoading();
    updateGlobalLoading(requestCount - 1);
    return Promise.reject(error);
  }
);

request.interceptors.response.use(
  async (response) => {
    const { config } = response;
    if (!config.noLoading) {
      finishRequestLoading();
      updateGlobalLoading(requestCount - 1);
    }
    if (response.data instanceof Blob && response.headers?.['file'] === 'true') {
      return response;
    }

    let { code, msg, data } = response.data;

    if (response.data instanceof Blob) {
      let parsed;
      try {
        const text = await response.data.text();
        parsed = JSON.parse(text);
      } catch {
        return Promise.reject(new Error(response.data));
      }
      code = parsed.code;
      msg = parsed.msg;
      data = parsed.data;
    }

    if (code !== 200 && !config.noToast) {
      toast.warning({
        description: msg,
      });
    }
    if (code === 401) {
      import('@/store').then(({ store }) => {
        import('@/store/user').then(({ logout }) => {
          store.dispatch(logout());
        });
      });
    }

    return { code, msg, data };
  },
  async (error) => {
    if (!error.config?.noLoading) {
      finishRequestLoading();
      updateGlobalLoading(requestCount - 1);
    }

    let errorMessage;
    if (error.response) {
      const { status, data } = error.response;
      // responseType: 'blob' 时错误体也是 Blob, 需异步解析
      let resolved = data;
      if (data instanceof Blob) {
        try {
          resolved = JSON.parse(await data.text());
        } catch {
          resolved = null;
        }
      }
      switch (status) {
        case 401:
          errorMessage = i18n.t('errors.unauthorized');
          // 使用动态导入避免循环依赖
          import('@/store').then(({ store }) => {
            import('@/store/user').then(({ logout }) => {
              store.dispatch(logout());
            });
          });
          break;
        case 403:
          errorMessage = i18n.t('errors.forbidden');
          break;
        case 404:
          errorMessage = i18n.t('errors.notFound');
          break;
        case 500:
          errorMessage = i18n.t('errors.serverError');
          break;
        case 502:
          errorMessage = i18n.t('errors.badGateway');
          break;
        case 503:
          errorMessage = i18n.t('errors.serviceUnavailable');
          break;
        case 504:
          errorMessage = i18n.t('errors.gatewayTimeout');
          break;
        default:
          errorMessage = resolved?.msg || i18n.t('errors.requestFailed');
      }
    } else if (error.request) {
      if (error.code === 'ECONNABORTED') {
        errorMessage = i18n.t('errors.requestTimeout');
      } else {
        errorMessage = i18n.t('errors.networkError');
      }
    } else {
      errorMessage = error.message;
    }

    // Forms with inline errors own their feedback.
    if (!error.config?.noToast) {
      toast.warning({
        description: errorMessage,
      });
    }
    return Promise.reject(error);
  }
);

export default request;
