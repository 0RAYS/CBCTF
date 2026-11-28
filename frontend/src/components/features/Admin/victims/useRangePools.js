import { useCallback, useEffect, useRef, useState } from 'react';
import { getContestRangePools, prewarmContestRangePools } from '@/api/admin/contest';

const POOL_POLL_INTERVAL = 15000;

// useRangePools 读取比赛里靶场题目的实例池状态,并提供手动补货。
export default function useRangePools(contestId) {
  const [pools, setPools] = useState([]);
  const [loading, setLoading] = useState(true);
  const [prewarming, setPrewarming] = useState(false);
  const requestId = useRef(0);

  const load = useCallback(async () => {
    const current = ++requestId.current;
    try {
      const response = await getContestRangePools(contestId);
      if (current !== requestId.current) return;
      if (response.code === 200) setPools(response.data.pools || []);
    } catch {
      if (current === requestId.current) setPools([]);
    } finally {
      if (current === requestId.current) setLoading(false);
    }
  }, [contestId]);

  const prewarm = useCallback(async () => {
    setPrewarming(true);
    try {
      await prewarmContestRangePools(contestId);
      await load();
    } finally {
      setPrewarming(false);
    }
  }, [contestId, load]);

  useEffect(() => {
    load();
    const timer = setInterval(load, POOL_POLL_INTERVAL);
    return () => {
      clearInterval(timer);
      requestId.current += 1;
    };
  }, [load]);

  return { pools, loading, prewarming, prewarm, refresh: load };
}
