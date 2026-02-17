import { useState, useEffect, useCallback, useRef } from 'react';

/**
 * 防抖搜索 Hook
 * @param {Function} searchFn - 搜索函数，接收 query 参数
 * @param {Object} options - 配置选项
 * @param {number} options.delay - 防抖延迟时间（毫秒），默认 300ms
 * @param {number} options.minLength - 最小搜索长度，默认 1
 * @param {boolean} options.immediate - 是否立即执行首次搜索，默认 false
 * @returns {Object} { query, setQuery, results, loading, error, reset }
 *
 * @example
 * const { query, setQuery, results, loading } = useDebounceSearch(
 *   async (q) => {
 *     const response = await searchAPI(q);
 *     return response.data;
 *   },
 *   { delay: 300, minLength: 2 }
 * );
 */
export function useDebounceSearch(searchFn, options = {}) {
  const { delay = 300, minLength = 1, immediate = false } = options;

  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // 使用 ref 追踪最新的搜索请求
  const abortControllerRef = useRef(null);
  const mountedRef = useRef(true);
  // 用 ref 持有 searchFn，避免因 inline 函数引用变化导致 executeSearch 重建
  const searchFnRef = useRef(searchFn);
  useEffect(() => {
    searchFnRef.current = searchFn;
  });

  // 搜索函数
  const executeSearch = useCallback(async (searchQuery) => {
    // 取消之前的请求
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    // 创建新的 AbortController
    abortControllerRef.current = new AbortController();

    setLoading(true);
    setError(null);

    try {
      const data = await searchFnRef.current(searchQuery, {
        signal: abortControllerRef.current.signal,
      });

      // 只有组件还在挂载时才更新状态
      if (mountedRef.current) {
        setResults(data || []);
        setLoading(false);
      }
    } catch (err) {
      // 忽略取消的请求
      if (err.name === 'AbortError') {
        return;
      }

      if (mountedRef.current) {
        setError(err);
        setResults([]);
        setLoading(false);
      }
    }
  }, []);

  // 防抖效果
  useEffect(() => {
    // 如果查询为空或长度不足，清空结果
    if (!query || query.length < minLength) {
      setResults([]);
      setLoading(false);
      return;
    }

    // 立即执行搜索（如果配置了 immediate）
    if (immediate && query.length === minLength) {
      executeSearch(query);
      return;
    }

    // 设置防抖定时器
    const timer = setTimeout(() => {
      executeSearch(query);
    }, delay);

    // 清理函数
    return () => {
      clearTimeout(timer);
    };
  }, [query, delay, minLength, immediate, executeSearch]);

  // 组件卸载时清理
  useEffect(() => {
    mountedRef.current = true;

    return () => {
      mountedRef.current = false;
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  // 重置函数
  const reset = useCallback(() => {
    setQuery('');
    setResults([]);
    setError(null);
    setLoading(false);
  }, []);

  return {
    query,
    setQuery,
    results,
    loading,
    error,
    reset,
  };
}

/**
 * 简单的防抖 Hook（不包含搜索逻辑）
 * @param {*} value - 要防抖的值
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {*} 防抖后的值
 *
 * @example
 * const debouncedSearchTerm = useDebounce(searchTerm, 500);
 *
 * useEffect(() => {
 *   if (debouncedSearchTerm) {
 *     // 执行搜索
 *   }
 * }, [debouncedSearchTerm]);
 */
export function useDebounce(value, delay = 300) {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(timer);
    };
  }, [value, delay]);

  return debouncedValue;
}
