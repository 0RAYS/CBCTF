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
  const [lines, setLines] = useState(1000);
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);
  const request = useRef(0);
  const previousSelection = useRef('');
  const victimId = victim?.id;
  const fetchPods = useEffectEvent(() => loadPods(victim.id));
  const fetchLogs = useEffectEvent(() => loadLogs(victim.id, podName, containerName, lines));
  const reportError = useEffectEvent((key) => toast.danger({ description: t(`${translationKey}.logs.${key}`) }));

  useEffect(() => {
    let active = true;
    setPods([]);
    setPodsLoading(true);
    setPodsVictimId(null);
    setPodName('');
    setContainerName('');
    setLines(1000);
    setContent('');
    setLoading(false);
    previousSelection.current = '';
    if (!victimId) return;
    fetchPods()
      .then((response) => {
        if (!active) return;
        if (response?.code !== 200) throw new Error(response?.msg);
        const available = response.data?.pods ?? [];
        setPods(available);
        setPodsVictimId(victimId);
        setPodName(available[0]?.name ?? '');
        setContainerName(available[0]?.containers?.[0] ?? '');
      })
      .catch(() => {
        if (active) reportError('fetchPodsFailed');
      })
      .finally(() => {
        if (active) setPodsLoading(false);
      });
    return () => {
      active = false;
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
  }, [victimId, podsVictimId, podName, containerName, lines]);

  return {
    pods,
    podsLoading,
    podName,
    containerName,
    lines,
    content,
    loading,
    setLines,
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
