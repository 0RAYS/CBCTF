import { useEffect, useEffectEvent, useRef, useState } from 'react';
import { toast } from '../../../../utils/toast';
import {
  GENERATOR_PAGE_SIZE,
  expandStartCounts,
  getGeneratorResponseData,
  isGeneratorStoppable,
  retainStoppableSelection,
} from './generatorUtils.js';

export default function useGeneratorSession(api, text) {
  const [generators, setGenerators] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedIds, setSelectedIds] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showDeleted, setShowDeleted] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState(10);
  const [revision, setRevision] = useState(0);
  const [dynamicChallenges, setDynamicChallenges] = useState([]);
  const [startModalOpen, setStartModalOpen] = useState(false);
  const [pendingOperation, setPendingOperation] = useState(null);
  const sessionRef = useRef(null);

  useEffect(() => {
    const session = { active: true, operation: null };
    sessionRef.current = session;
    return () => {
      session.active = false;
    };
  }, []);

  const queryGenerators = useEffectEvent((params) => api.list(params));
  const queryChallenges = useEffectEvent(() => api.challenges({ limit: 100, offset: 0, type: 'dynamic' }));
  const reportFetchError = useEffectEvent(() => toast.danger({ description: text('toast.fetchFailed') }));

  useEffect(() => {
    let active = true;
    let requestId = 0;
    let inFlight = false;
    const fetchPage = async () => {
      // Polling must not supersede a slow request for this query.
      if (inFlight) return;
      inFlight = true;
      const request = ++requestId;
      setLoading(true);
      try {
        const response = await queryGenerators({
          limit: GENERATOR_PAGE_SIZE,
          offset: (currentPage - 1) * GENERATOR_PAGE_SIZE,
          ...(showDeleted && { deleted: true }),
        });
        if (!active || request !== requestId) return;
        const data = getGeneratorResponseData(response);
        const next = data.generators ?? [];
        setGenerators(next);
        setTotalCount(data.count ?? 0);
        setSelectedIds((ids) => retainStoppableSelection(ids, next));
      } catch {
        if (!active || request !== requestId) return;
        setGenerators([]);
        setTotalCount(0);
        setSelectedIds([]);
        reportFetchError();
      } finally {
        inFlight = false;
        if (active && request === requestId) setLoading(false);
      }
    };
    fetchPage();
    const timer = refreshInterval > 0 ? setInterval(fetchPage, refreshInterval * 1000) : null;
    return () => {
      active = false;
      clearInterval(timer);
    };
  }, [currentPage, showDeleted, refreshInterval, revision]);

  useEffect(() => {
    let active = true;
    const fetchChallenges = async () => {
      try {
        const response = await queryChallenges();
        if (active) setDynamicChallenges(getGeneratorResponseData(response).challenges ?? []);
      } catch {
        // The start dialog retains the existing empty-list fallback.
        if (active) setDynamicChallenges([]);
      }
    };
    fetchChallenges();
    return () => {
      active = false;
    };
  }, []);

  const refresh = () => setRevision((value) => value + 1);

  const changePage = (page) => {
    setCurrentPage(page);
    setSelectedIds([]);
  };

  const toggleShowDeleted = () => {
    setShowDeleted((value) => !value);
    changePage(1);
  };

  const toggleSelect = (id) => {
    if (!isGeneratorStoppable(generators.find((generator) => generator.id === id))) return;
    setSelectedIds((ids) => (ids.includes(id) ? ids.filter((selected) => selected !== id) : [...ids, id]));
  };

  const toggleSelectAll = () => {
    const ids = generators.filter(isGeneratorStoppable).map((generator) => generator.id);
    setSelectedIds((selected) => (ids.every((id) => selected.includes(id)) ? [] : ids));
  };

  const closeStart = () => {
    if (sessionRef.current?.operation !== 'start') setStartModalOpen(false);
  };

  const start = async (counts) => {
    const session = sessionRef.current;
    if (!session?.active || session.operation) return;
    const challenges = expandStartCounts(counts);
    if (challenges.length === 0) {
      toast.warning({ description: text('toast.selectRequired') });
      return;
    }
    session.operation = 'start';
    setPendingOperation('start');
    try {
      const response = await api.start(challenges);
      if (!session.active) return;
      getGeneratorResponseData(response);
      setStartModalOpen(false);
      changePage(1);
      refresh();
    } catch {
      if (session.active) toast.danger({ description: text('toast.startFailed') });
    } finally {
      if (session.active) {
        session.operation = null;
        setPendingOperation(null);
      }
    }
  };

  const stop = async () => {
    const session = sessionRef.current;
    const ids = retainStoppableSelection(selectedIds, generators);
    if (!session?.active || session.operation || ids.length === 0) return;
    session.operation = 'stop';
    setPendingOperation('stop');
    try {
      const response = await api.stop(ids);
      if (!session.active) return;
      getGeneratorResponseData(response);
      toast.success({ description: text('toast.stopSuccess') });
      setSelectedIds((selected) => selected.filter((id) => !ids.includes(id)));
      refresh();
    } catch {
      if (session.active) toast.danger({ description: text('toast.stopFailed') });
    } finally {
      if (session.active) {
        session.operation = null;
        setPendingOperation(null);
      }
    }
  };

  return {
    generators,
    totalCount,
    currentPage,
    selectedIds,
    loading,
    showDeleted,
    refreshInterval,
    setRefreshInterval,
    dynamicChallenges,
    startModalOpen,
    openStart: () => setStartModalOpen(true),
    closeStart,
    pendingOperation,
    refresh,
    changePage,
    toggleShowDeleted,
    toggleSelect,
    toggleSelectAll,
    start,
    stop,
  };
}
