import { useLayoutEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getContestTeamTraffic, downloadContainerTraffic } from '../../../../api/admin/contest.js';
import { downloadVictimTraffic } from '../../../../api/admin/victims.js';
import { downloadBlobResponse } from '../../../../utils/fileDownload';
import { toast } from '../../../../utils/toast';
import { createDemoTopology } from './trafficDemo.js';
import { sanitizeSlice } from './trafficPresentation.js';

export default function useTrafficSession({ isOpen, container, contestId, teamId, fetchTraffic: customFetchTraffic }) {
  const { t } = useTranslation();
  const scopeKey = JSON.stringify([contestId ?? null, teamId ?? null, container?.id ?? null, !!customFetchTraffic]);
  const scopeRef = useRef(null);
  const requestSequenceRef = useRef(0);
  const [topology, setTopology] = useState(null);
  const [universeNodes, setUniverseNodes] = useState(null);
  const [demoMode, setDemoMode] = useState(false);
  const [isFetching, setIsFetching] = useState(false);
  const [selectionVersion, setSelectionVersion] = useState(0);

  useLayoutEffect(() => {
    const scope = { active: isOpen, universePending: false, universeFetched: false };
    scopeRef.current = scope;
    setTopology(null);
    setUniverseNodes(null);
    setDemoMode(false);
    setIsFetching(false);
    setSelectionVersion((current) => current + 1);
    return () => {
      // Invalidate both request lanes on close, scope change, and unmount.
      scope.active = false;
      requestSequenceRef.current += 1;
    };
  }, [isOpen, scopeKey]);

  const fetchData = async ({ nextShift, nextSlice, forceLive = false }) => {
    const scope = scopeRef.current;
    if (!container?.id || !scope?.active) return;
    const resolvedSlice = sanitizeSlice(nextSlice);
    const requestId = ++requestSequenceRef.current;
    const isCurrent = () => scope.active && requestSequenceRef.current === requestId;
    const request = (params) =>
      customFetchTraffic
        ? customFetchTraffic(container, params)
        : getContestTeamTraffic(contestId, teamId, container.id, params);

    const fetchUniverseNodes = async (totalDuration) => {
      if (totalDuration <= 0 || scope.universePending || scope.universeFetched) return;
      scope.universePending = true;
      try {
        const response = await request({ time_shift: 0, duration: totalDuration });
        // A newer frame may finish first; the universe belongs to the scope, not a frame.
        if (!scope.active || response.code !== 200) return;
        scope.universeFetched = true;
        setUniverseNodes(response.data?.nodes || []);
      } catch {
        // Keep current-frame layout when the full-duration request fails.
      } finally {
        scope.universePending = false;
      }
    };

    setIsFetching(true);
    try {
      const response = await request({ time_shift: nextShift, duration: resolvedSlice });
      if (!isCurrent()) return;
      if (response.code !== 200) throw new Error(t('admin.contests.trafficGraph.toast.fetchFailed'));
      setTopology(response.data);
      setDemoMode(false);
      void fetchUniverseNodes(response.data?.total_duration || 0);
    } catch {
      if (!isCurrent()) return;
      if (!forceLive) toast.warning({ description: t('admin.contests.trafficGraph.toast.demoFallback') });
      setTopology(createDemoTopology(nextShift, resolvedSlice, container.id));
      setDemoMode(true);
    } finally {
      if (isCurrent()) {
        setIsFetching(false);
        setSelectionVersion((current) => current + 1);
      }
    }
  };

  const invalidateFrame = () => {
    requestSequenceRef.current += 1;
  };

  const downloadTraffic = async () => {
    if (!container?.id) return;
    try {
      const response =
        contestId && teamId
          ? await downloadContainerTraffic(contestId, teamId, container.id)
          : await downloadVictimTraffic(container.id);
      if (response.headers?.['file'] === 'true') downloadBlobResponse(response, `traffic_${container.id}.zip`);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.downloadTrafficFailed') });
    }
  };

  return {
    scopeKey,
    topology,
    universeNodes,
    demoMode,
    isFetching,
    selectionVersion,
    fetchData,
    invalidateFrame,
    downloadTraffic,
  };
}
