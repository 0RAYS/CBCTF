import { useState, useCallback, useRef, useEffect } from 'react';
import { toast } from '../utils/toast';
import { useTranslation } from 'react-i18next';

/**
 * 通用 API 调用 Hook
 * 提供统一的加载状态、错误处理和成功消息显示
 *
 * @param {Function} apiFn - API 函数
 * @param {Object} options - 配置选项
 * @param {string} options.successMessage - 成功提示消息
 * @param {string|boolean} options.errorMessage - 错误提示消息，false 则不显示
 * @param {Function} options.onSuccess - 成功回调
 * @param {Function} options.onError - 错误回调
 * @param {boolean} options.autoReset - 是否自动重置数据，默认 false
 * @returns {Object} { data, loading, error, execute, reset }
 *
 * @example
 * const { data, loading, execute } = useApiCall(
 *   getContestVictims,
 *   {
 *     successMessage: '获取成功',
 *     errorMessage: '获取失败'
 *   }
 * );
 *
 * // 调用 API
 * const handleFetch = async () => {
 *   const result = await execute(contestId, params);
 *   if (result.success) {
 *     console.log(result.data);
 *   }
 * };
 */
export function useApiCall(apiFn, options = {}) {
  const { successMessage, errorMessage, onSuccess, onError, autoReset = false } = options;

  const { t } = useTranslation();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // 使用 ref 追踪组件挂载状态
  const mountedRef = useRef(true);

  // 执行 API 调用
  const execute = useCallback(
    async (...args) => {
      setLoading(true);
      setError(null);

      // 如果配置了自动重置，清空之前的数据
      if (autoReset) {
        setData(null);
      }

      try {
        const response = await apiFn(...args);

        // 检查组件是否还在挂载
        if (!mountedRef.current) {
          return { success: false, error: new Error('Component unmounted') };
        }

        // 处理响应
        if (response.code === 200 || response.code === 201) {
          setData(response.data);

          // 显示成功消息
          if (successMessage) {
            toast.success({ description: successMessage });
          }

          // 调用成功回调
          if (onSuccess) {
            onSuccess(response.data);
          }

          return { success: true, data: response.data };
        } else {
          throw new Error(response.msg || response.message || 'Request failed');
        }
      } catch (err) {
        // 检查组件是否还在挂载
        if (!mountedRef.current) {
          return { success: false, error: err };
        }

        setError(err);

        // 显示错误消息（如果配置了）
        if (errorMessage !== false) {
          const message = errorMessage || err.message || t('toast.common.operationFailed');
          toast.danger({ description: message });
        }

        // 调用错误回调
        if (onError) {
          onError(err);
        }

        return { success: false, error: err };
      } finally {
        if (mountedRef.current) {
          setLoading(false);
        }
      }
    },
    [apiFn, successMessage, errorMessage, autoReset, onSuccess, onError, t]
  );

  // 重置状态
  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);

  // 清理函数
  useEffect(() => {
    mountedRef.current = true;

    return () => {
      mountedRef.current = false;
    };
  }, []);

  return {
    data,
    loading,
    error,
    execute,
    reset,
  };
}

/**
 * 带自动执行的 API 调用 Hook
 * 组件挂载时自动执行一次 API 调用
 *
 * @param {Function} apiFn - API 函数
 * @param {Array} deps - 依赖数组，变化时重新执行
 * @param {Object} options - 配置选项（同 useApiCall）
 * @returns {Object} { data, loading, error, refetch, reset }
 *
 * @example
 * const { data, loading, refetch } = useAutoApiCall(
 *   () => getContestVictims(contestId),
 *   [contestId],
 *   { successMessage: '加载成功' }
 * );
 */
export function useAutoApiCall(apiFn, deps = [], options = {}) {
  const { data, loading, error, execute, reset } = useApiCall(apiFn, options);

  // 自动执行
  useEffect(() => {
    execute();
  }, deps);

  return {
    data,
    loading,
    error,
    refetch: execute,
    reset,
  };
}
