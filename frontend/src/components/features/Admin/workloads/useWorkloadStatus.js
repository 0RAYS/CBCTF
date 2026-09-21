import { useEffect, useEffectEvent, useRef, useState } from 'react';

// One request at a time; refresh never supersedes an in-flight request for this
// resource. Cleanup aborts the transport and invalidates late responses.
export default function useWorkloadStatus(resourceId, load) {
  const [state, setState] = useState({
    pods: [],
    refreshing: false,
    statusError: false,
    updatedAt: null,
  });
  const refresh = useRef(() => {});
  const fetchStatus = useEffectEvent((signal) => load(resourceId, signal));
  useEffect(() => {
    let active = true;
    let timer;
    let inFlight = false;
    const controller = new AbortController();
    setState({
      pods: [],
      refreshing: false,
      statusError: false,
      updatedAt: null,
    });
    const poll = async () => {
      if (!active || inFlight || !resourceId) return;
      clearTimeout(timer);
      inFlight = true;
      setState((previous) => ({ ...previous, refreshing: true }));
      try {
        const response = await fetchStatus(controller.signal);
        if (!active) return;
        if (response?.code !== 200) throw new Error(response?.msg);
        setState({
          pods: response.data?.pods ?? [],
          refreshing: false,
          statusError: false,
          updatedAt: new Date().toLocaleTimeString(),
        });
      } catch {
        if (active)
          setState((previous) => ({
            ...previous,
            refreshing: false,
            statusError: true,
          }));
      } finally {
        inFlight = false;
        if (active) timer = setTimeout(poll, 5000);
      }
    };
    refresh.current = poll;
    poll();
    return () => {
      active = false;
      controller.abort();
      clearTimeout(timer);
      refresh.current = () => {};
    };
  }, [resourceId]);
  return { ...state, refreshStatus: () => refresh.current() };
}
