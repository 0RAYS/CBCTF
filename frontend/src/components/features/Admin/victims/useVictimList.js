import { useEffect, useEffectEvent, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from '../../../../utils/toast';
import { buildVictimListParams, selectPageVictims, toggleVictimSelection } from './victimPayload';

export default function useVictimList(scope) {
  const { t } = useTranslation();
  const [containers, setContainers] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [runningCount, setRunningCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [filters, setFilters] = useState(
    scope.contestId ? { user_id: '', team_id: '', challenge_id: '' } : { user_id: '', challenge_id: '' }
  );
  const [selectedContainers, setSelectedContainers] = useState([]);
  const [showDeleted, setShowDeleted] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState(10);
  const request = useRef(0);
  const inFlight = useRef(null);
  const pageSize = 20;

  const refresh = async () => {
    const session = ++request.current;
    inFlight.current = session;
    try {
      const response = await scope.loadVictims(buildVictimListParams(filters, currentPage, showDeleted, pageSize));
      if (session !== request.current || response.code !== 200) return;
      setContainers(response.data.victims || []);
      setTotalCount(response.data.count || 0);
      setRunningCount(response.data.running || 0);
    } catch (error) {
      if (session === request.current) {
        toast.danger({ description: error.message || t(`${scope.translationKey}.toast.fetchContainersFailed`) });
      }
    } finally {
      if (inFlight.current === session) inFlight.current = null;
    }
  };
  const refreshLatest = useEffectEvent(refresh);

  useEffect(() => {
    refreshLatest();
    return () => {
      request.current += 1;
      inFlight.current = null;
    };
  }, [currentPage, filters, showDeleted, scope.contestId, scope.translationKey]);

  useEffect(() => {
    if (refreshInterval <= 0) return;
    const timer = setInterval(() => {
      // Only polling skips busy queries; manual refresh still supersedes them.
      if (inFlight.current === null) refreshLatest();
    }, refreshInterval * 1000);
    return () => clearInterval(timer);
  }, [refreshInterval]);

  return {
    containers,
    totalCount,
    runningCount,
    currentPage,
    pageSize,
    filters,
    selectedContainers,
    showDeleted,
    refreshInterval,
    setRefreshInterval,
    refresh,
    setSelectedContainers,
    stats: {
      totalContainers: totalCount,
      runningContainers: runningCount,
      stoppedContainers: totalCount - runningCount,
    },
    onFilterChange: (key, value) => {
      setFilters((previous) => ({ ...previous, [key]: value }));
      setCurrentPage(1);
    },
    onResetFilters: () => {
      setFilters((previous) => Object.fromEntries(Object.keys(previous).map((key) => [key, ''])));
      setCurrentPage(1);
    },
    onPageChange: setCurrentPage,
    onToggleShowDeleted: () => {
      setShowDeleted((previous) => !previous);
      setCurrentPage(1);
      setSelectedContainers([]);
    },
    onContainerSelect: (id) => setSelectedContainers((previous) => toggleVictimSelection(previous, containers, id)),
    onSelectAll: () => setSelectedContainers((previous) => selectPageVictims(previous, containers)),
  };
}
