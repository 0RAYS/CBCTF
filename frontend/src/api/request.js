import axios from 'axios';
import { toast } from '../utils/toast';
import i18n from '../i18n';
import { API_CONFIG } from './config.js';

let nprogressPromise = null;

function loadNProgress() {
  if (!nprogressPromise) {
    nprogressPromise = import('../utils/nprogress');
  }
  return nprogressPromise;
}

function startRequestLoading() {
  loadNProgress().then(({ startLoading }) => startLoading());
}

function finishRequestLoading() {
  loadNProgress().then(({ finishLoading }) => finishLoading());
}

// 创建 Axios 实例
const request = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 正在进行的请求数量
let requestCount = 0;

// 更新全局loading状态
const updateGlobalLoading = (count) => {
  requestCount = Math.max(0, count);
  // 使用动态导入避免循环依赖
  import('../store').then(({ store }) => {
    import('../store/app').then(({ setGlobalLoading }) => {
      store.dispatch(setGlobalLoading(requestCount > 0));
    });
  });
};

// 请求拦截器
request.interceptors.request.use(
  async (config) => {
    // 如果没有设置 noLoading 标识, 则执行全局 loading 逻辑
    if (!config.noLoading) {
      startRequestLoading();
      updateGlobalLoading(requestCount + 1);
    }
    config.headers['Accept-Language'] = i18n.language;
    return config;
  },
  (error) => {
    // 请求错误时减少计数
    finishRequestLoading();
    updateGlobalLoading(requestCount - 1);
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  async (response) => {
    // 如果没有设置 noLoading 标识, 则执行全局 loading 逻辑
    const { config } = response;
    if (!config.noLoading) {
      finishRequestLoading();
      updateGlobalLoading(requestCount - 1);
    }
    // 处理文件下载
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

    // 处理业务错误
    if (code !== 200) {
      toast.warning({
        description: msg,
      });
    }
    if (code === 401) {
      import('../store').then(({ store }) => {
        import('../store/user').then(({ logout }) => {
          store.dispatch(logout());
        });
      });
    }

    return { code, msg, data };
  },
  async (error) => {
    // 减少请求计数
    if (!error.config?.noLoading) {
      finishRequestLoading();
      updateGlobalLoading(requestCount - 1);
    }

    // 处理错误响应
    let errorMessage;
    if (error.response) {
      const { status, data } = error.response;
      // responseType: 'blob' 时错误体也是 Blob, 需异步解析
      const resolveData = async () => {
        if (data instanceof Blob) {
          try {
            const text = await data.text();
            return JSON.parse(text);
          } catch {
            return null;
          }
        }
        return data;
      };

      const resolved = await resolveData();
      switch (status) {
        case 401:
          errorMessage = i18n.t('errors.unauthorized');
          // 使用动态导入避免循环依赖
          import('../store').then(({ store }) => {
            import('../store/user').then(({ logout }) => {
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

    // 显示错误信息
    toast.warning({
      description: errorMessage,
    });
    return Promise.reject(error);
  }
);

export default request;
