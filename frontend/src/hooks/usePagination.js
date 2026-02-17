import { useState, useEffect, useCallback, useRef } from 'react';
import { toast } from '../utils/toast';

/**
 * 分页数据管理 Hook
 * 提供完整的分页功能，包括数据获取、页码管理、加载状态等
 *
 * @param {Function} fetchFn - 数据获取函数，接收 { limit, offset } 参数
 * @param {Object} options - 配置选项
 * @param {number} options.pageSize - 每页数据量，默认 20
 * @param {number} options.initialPage - 初始页码，默认 1
 * @param {boolean} options.autoFetch - 是否自动获取数据，默认 true
 * @param {Function} options.onSuccess - 成功回调
 * @param {Function} options.onError - 错误回调
 * @param {string} options.errorMessage - 错误提示消息
 * @returns {Object} { data, loading, error, currentPage, totalPages, totalCount, setPage, nextPage, prevPage, refetch, reset }
 *
 * @example
 * const {
 *   data,
 *   loading,
 *   currentPage,
 *   totalPages,
 *   setPage
 * } = usePagination(
 *   async ({ limit, offset }) => {
 *     const response = await getContainers({ limit, offset });
 *     return {
 *       data: response.data.list,
 *       total: response.data.total,
 *     };
 *   },
 *   { pageSize: 20 }
 * );
 */
export function usePagination(fetchFn, options = {}) {
  const { pageSize = 20, initialPage = 1, autoFetch = true, onSuccess, onError, errorMessage } = options;

  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(initialPage);
  const [totalCount, setTotalCount] = useState(0);

  const mountedRef = useRef(true);
  const abortControllerRef = useRef(null);

  // 计算总页数
  const totalPages = Math.ceil(totalCount / pageSize);

  // 获取数据
  const fetch = useCallback(
    async (page = currentPage) => {
      // 取消之前的请求
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }

      abortControllerRef.current = new AbortController();

      setLoading(true);
      setError(null);

      try {
        const offset = (page - 1) * pageSize;
        const result = await fetchFn({
          limit: pageSize,
          offset,
          signal: abortControllerRef.current.signal,
        });

        if (!mountedRef.current) return;

        // 处理返回数据
        // 支持两种返回格式：
        // 1. { data: [...], total: 100 }
        // 2. { data: [...], count: 100 }
        const items = result.data || result;
        const total = result.total ?? result.count ?? items.length;

        setData(items);
        setTotalCount(total);
        setLoading(false);

        if (onSuccess) {
          onSuccess({ data: items, total });
        }
      } catch (err) {
        // 忽略取消的请求
        if (err.name === 'AbortError') {
          return;
        }

        if (!mountedRef.current) return;

        setError(err);
        setLoading(false);

        if (errorMessage) {
          toast.danger({ description: errorMessage });
        }

        if (onError) {
          onError(err);
        }
      }
    },
    [currentPage, pageSize, fetchFn, onSuccess, onError, errorMessage]
  );

  // 设置页码
  const setPage = useCallback(
    (page) => {
      if (page < 1 || page > totalPages) return;
      setCurrentPage(page);
    },
    [totalPages]
  );

  // 下一页
  const nextPage = useCallback(() => {
    if (currentPage < totalPages) {
      setCurrentPage((prev) => prev + 1);
    }
  }, [currentPage, totalPages]);

  // 上一页
  const prevPage = useCallback(() => {
    if (currentPage > 1) {
      setCurrentPage((prev) => prev - 1);
    }
  }, [currentPage]);

  // 跳转到第一页
  const firstPage = useCallback(() => {
    setCurrentPage(1);
  }, []);

  // 跳转到最后一页
  const lastPage = useCallback(() => {
    if (totalPages > 0) {
      setCurrentPage(totalPages);
    }
  }, [totalPages]);

  // 重置
  const reset = useCallback(() => {
    setData([]);
    setCurrentPage(initialPage);
    setTotalCount(0);
    setError(null);
    setLoading(false);
  }, [initialPage]);

  // 当页码变化时获取数据
  useEffect(() => {
    if (autoFetch) {
      fetch(currentPage);
    }
  }, [currentPage]);

  // 清理函数
  useEffect(() => {
    mountedRef.current = true;

    return () => {
      mountedRef.current = false;
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  return {
    data,
    loading,
    error,
    currentPage,
    totalPages,
    totalCount,
    pageSize,
    setPage,
    nextPage,
    prevPage,
    firstPage,
    lastPage,
    refetch: fetch,
    reset,
    // 辅助信息
    hasNextPage: currentPage < totalPages,
    hasPrevPage: currentPage > 1,
    isEmpty: data.length === 0 && !loading,
  };
}

/**
 * 简化的分页 Hook（仅管理页码状态）
 * @param {number} totalPages - 总页数
 * @param {number} initialPage - 初始页码
 * @returns {Object} { currentPage, setPage, nextPage, prevPage, firstPage, lastPage }
 *
 * @example
 * const { currentPage, setPage, nextPage, prevPage } = useSimplePagination(10);
 */
export function useSimplePagination(totalPages, initialPage = 1) {
  const [currentPage, setCurrentPage] = useState(initialPage);

  const setPage = useCallback(
    (page) => {
      if (page < 1 || page > totalPages) return;
      setCurrentPage(page);
    },
    [totalPages]
  );

  const nextPage = useCallback(() => {
    if (currentPage < totalPages) {
      setCurrentPage((prev) => prev + 1);
    }
  }, [currentPage, totalPages]);

  const prevPage = useCallback(() => {
    if (currentPage > 1) {
      setCurrentPage((prev) => prev - 1);
    }
  }, [currentPage]);

  const firstPage = useCallback(() => {
    setCurrentPage(1);
  }, []);

  const lastPage = useCallback(() => {
    if (totalPages > 0) {
      setCurrentPage(totalPages);
    }
  }, [totalPages]);

  return {
    currentPage,
    setPage,
    nextPage,
    prevPage,
    firstPage,
    lastPage,
    hasNextPage: currentPage < totalPages,
    hasPrevPage: currentPage > 1,
  };
}
