import axios from 'axios';
import { toast } from '../utils/toast';
import i18n from '../i18n';
import { API_CONFIG } from './config.js';
import FingerprintJS from '@fingerprintjs/fingerprintjs';
import { startLoading, finishLoading } from '../utils/nprogress';

const NONCE_KEY = 'LXM_NONCE';
const FINGERPRINT_KEY = 'LXM';

const fpPromise = FingerprintJS.load();

// SHA-256 哈希, 带非安全上下文回退
async function sha256Hex(input) {
  if (crypto.subtle) {
    const data = new TextEncoder().encode(input);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    return Array.from(new Uint8Array(hashBuffer))
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('');
  }
  // 非 HTTPS 环境回退：直接返回原始拼接字符串, 后端会做 double-MD5
  return input;
}

// 获取或创建浏览器实例唯一标识
function getOrCreateNonce() {
  let nonce = localStorage.getItem(NONCE_KEY);
  if (!nonce) {
    nonce =
      typeof crypto.randomUUID === 'function'
        ? crypto.randomUUID()
        : ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) =>
            (c ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))).toString(16)
          );
    localStorage.setItem(NONCE_KEY, nonce);
  }
  return nonce;
}

// 混合指纹：visitorId（硬件特征）+ nonce（实例标识）→ SHA-256
async function generateFingerprint() {
  try {
    const fp = await fpPromise;
    const result = await fp.get();
    const nonce = getOrCreateNonce();
    return await sha256Hex(result.visitorId + ':' + nonce);
  } catch {
    // FingerprintJS 失败时仅用 nonce, 仍保证唯一性
    const nonce = getOrCreateNonce();
    return await sha256Hex('fp-fallback:' + nonce);
  }
}

let magicNum = localStorage.getItem(FINGERPRINT_KEY) || '';
let fingerprintPromise = null;

async function ensureFingerprint() {
  if (magicNum) {
    return magicNum;
  }

  if (!fingerprintPromise) {
    fingerprintPromise = generateFingerprint()
      .then((fingerprint) => {
        magicNum = fingerprint;
        localStorage.setItem(FINGERPRINT_KEY, fingerprint);
        return fingerprint;
      })
      .finally(() => {
        fingerprintPromise = null;
      });
  }

  return fingerprintPromise;
}

ensureFingerprint();

// 清除设备指纹缓存（登出时调用）
export function clearFingerprint() {
  magicNum = '';
  fingerprintPromise = null;
  localStorage.removeItem(FINGERPRINT_KEY);
  localStorage.removeItem(NONCE_KEY);
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
  requestCount = count;
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
      startLoading();
      updateGlobalLoading(requestCount + 1);
    }
    config.headers['X-M'] = await ensureFingerprint();
    config.headers['Accept-Language'] = i18n.language;
    return config;
  },
  (error) => {
    // 请求错误时减少计数
    finishLoading();
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
      finishLoading();
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
          clearFingerprint();
          store.dispatch(logout());
        });
      });
    }

    return { code, msg, data };
  },
  async (error) => {
    // 减少请求计数
    finishLoading();
    updateGlobalLoading(requestCount - 1);

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
              clearFingerprint();
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
