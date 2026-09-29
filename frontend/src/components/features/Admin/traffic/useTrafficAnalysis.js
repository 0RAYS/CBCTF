import { useEffect, useState } from 'react';
import { getTrafficAnalysis, getTrafficOverlaps } from '../../../../api/admin/traffic.js';

export default function useTrafficAnalysis({ isOpen, victimId, contestId, teamId }) {
  const [revision, setRevision] = useState(0);
  const [state, setState] = useState({});
  const scope = JSON.stringify([victimId, contestId, teamId]);
  useEffect(() => {
    if (!isOpen || !victimId) return;
    const controller = new AbortController();
    const { signal } = controller;
    setState({ scope, loading: true });
    const load = async () => {
      try {
        const response = await getTrafficAnalysis({ contestId, teamId, victimId }, signal);
        if (signal.aborted) return;
        if (response.code !== 200) throw new Error(response.msg);
        setState({
          scope,
          data: response.data,
          loading: false,
          overlapsLoading: !!contestId,
        });
        // The analysis request publishes live access evidence before querying overlaps.
        if (contestId) {
          try {
            const overlaps = await getTrafficOverlaps(contestId, signal);
            if (signal.aborted) return;
            if (overlaps.code !== 200) throw new Error(overlaps.msg);
            setState((current) => ({
              ...current,
              overlaps: overlaps.data || [],
              overlapsLoading: false,
            }));
          } catch {
            if (!signal.aborted)
              setState((current) => ({
                ...current,
                overlapsError: true,
                overlapsLoading: false,
              }));
          }
        }
      } catch {
        if (!signal.aborted) setState({ scope, error: true, loading: false });
      }
    };
    void load();
    return () => controller.abort();
  }, [isOpen, scope, victimId, contestId, teamId, revision]);
  return {
    ...(state.scope === scope ? state : { loading: true }),
    refresh: () => setRevision((value) => value + 1),
  };
}
