import { useEffect, useEffectEvent, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from '../../../../utils/toast';

export default function useVictimLogSession({ victim, loadPods, loadLogs, translationKey }) {
  const { t } = useTranslation();
  const [pods, setPods] = useState([]);
  const [podsLoading, setPodsLoading] = useState(true);
  const [podsVictimId, setPodsVictimId] = useState(null);
  const [podName, setPodName] = useState('');
  const [containerName, setContainerName] = useState('');
  const [logRevision, setLogRevision] = useState(0);
  const [lines, setLines] = useState(1000);
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);
  const [statusError, setStatusError] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [updatedAt, setUpdatedAt] = useState(null);
  const refreshStatus = useRef(() => {});
  const request = useRef(0);
  const previousSelection = useRef('');
  const victimId = victim?.id;
  const fetchPods = useEffectEvent((signal) => loadPods(victim.id, signal));
  const fetchLogs = useEffectEvent(() => loadLogs(victim.id, podName, containerName, lines));
  const reportError = useEffectEvent((key) => toast.danger({ description: t(`${translationKey}.logs.${key}`) }));

  const applyPods = useEffectEvent((available) => {
    const selected = available.find((pod) => pod.name === podName) ?? available[0];
    const nextPod = selected?.name ?? '';
    const nextContainer = selected?.containers?.includes(containerName)
      ? containerName
      : (selected?.containers?.[0] ?? '');
    if (nextPod !== podName || nextContainer !== containerName) {
      setContent('');
      setLoading(false);
    }
    setPods(available);
    setPodsVictimId(victimId);
    setPodName(nextPod);
    setContainerName(nextContainer);
  });

  useEffect(() => {
    let active = true;
    let timer;
    let inFlight = false;
    let first = true;
    const controller = new AbortController();
    setPods([]);
    setPodsLoading(!!victimId);
    setPodsVictimId(null);
    setPodName('');
    setContainerName('');
    setLines(1000);
    setContent('');
    setLoading(false);
    setStatusError(false);
    setUpdatedAt(null);
    setRefreshing(false);
    previousSelection.current = '';
    refreshStatus.current = () => {};
    if (!victimId) return;
    const poll = async () => {
      if (!active || inFlight) return;
      clearTimeout(timer);
      inFlight = true;
      setRefreshing(true);
      try {
        const response = await fetchPods(controller.signal);
        if (!active) return;
        if (response?.code !== 200) throw new Error(response?.msg);
        applyPods(response.data?.pods ?? []);
        setStatusError(false);
        setUpdatedAt(new Date().toLocaleTimeString());
      } catch {
        if (active) {
          setStatusError(true);
          if (first) reportError('fetchPodsFailed');
        }
      } finally {
        inFlight = false;
        if (active) {
          first = false;
          setPodsLoading(false);
          setRefreshing(false);
          timer = setTimeout(poll, 5000);
        }
      }
    };
    refreshStatus.current = poll;
    poll();
    return () => {
      active = false;
      controller.abort();
      clearTimeout(timer);
      refreshStatus.current = () => {};
    };
  }, [victimId]);

  useEffect(() => {
    if (!victimId || podsVictimId !== victimId || !podName || !containerName) return;
    const session = ++request.current;
    const selection = JSON.stringify([podName, containerName]);
    const delay = previousSelection.current === selection ? 500 : 0;
    previousSelection.current = selection;
    const timer = setTimeout(async () => {
      setLoading(true);
      try {
        const response = await fetchLogs();
        if (response?.code !== 200) throw new Error(response?.msg);
        if (session === request.current) setContent(response.data?.logs ?? '');
      } catch {
        if (session === request.current) reportError('fetchLogsFailed');
      } finally {
        if (session === request.current) setLoading(false);
      }
    }, delay);
    return () => {
      clearTimeout(timer);
      request.current += 1;
    };
  }, [victimId, podsVictimId, podName, containerName, lines, logRevision]);

  return {
    pods,
    podsLoading,
    statusError,
    refreshing,
    updatedAt,
    refreshStatus: () => refreshStatus.current(),
    podName,
    containerName,
    lines,
    content,
    loading,
    setLines,
    refreshLogs: () => setLogRevision((value) => value + 1),
    selectPod: (name) => {
      setPodName(name);
      setContainerName(pods.find((pod) => pod.name === name)?.containers?.[0] ?? '');
      setContent('');
      setLoading(false);
    },
    selectContainer: (name) => {
      setContainerName(name);
      setContent('');
      setLoading(false);
    },
  };
}
