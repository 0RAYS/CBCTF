import { useState, useEffect, useRef } from 'react';
import { toast } from '../../utils/toast';
import { getVictims, stopVictims } from '../../api/admin/victims';
import { Modal } from '../../components/common';
import ModalButton from '../../components/common/ModalButton';
import { Button, Pagination, Card, EmptyState, StatCard } from '../../components/common';
import { motion } from 'motion/react';
import {
  IconPlayerPlay,
  IconBan,
  IconFilter,
  IconRefresh,
  IconTable,
  IconServer,
  IconSearch,
  IconTarget,
  IconUsers,
  IconTrash,
  IconClockPlay,
} from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { searchModels } from '../../api/admin/search.js';

const VICTIM_STATUS_STYLES = {
  waiting: 'bg-yellow-400/10 text-yellow-400 border-yellow-400/30',
  pending: 'bg-geek-400/10 text-geek-400 border-geek-400/30',
  running: 'bg-green-400/10 text-green-400 border-green-400/30',
  stopped: 'bg-neutral-500/10 text-neutral-400 border-neutral-500/30',
};

function VictimStatusBadge({ status, t }) {
  const style = VICTIM_STATUS_STYLES[status] ?? VICTIM_STATUS_STYLES.stopped;
  return (
    <span className={`inline-block px-2 py-0.5 rounded border text-xs font-mono ${style}`}>
      {t(`admin.victims.statusBadge.${status}`, status)}
    </span>
  );
}

function AdminVictims() {
  const [containers, setContainers] = useState([]);
  const [runningCount, setRunningCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);

  const [filters, setFilters] = useState({
    user_id: '',
    team_id: '',
    challenge_id: '',
    limit: 20,
    offset: 0,
  });

  const [searchResults, setSearchResults] = useState({
    users: [],
    teams: [],
    challenges: [],
  });
  const [searchLoading, setSearchLoading] = useState({
    users: false,
    teams: false,
    challenges: false,
  });

  const usersSearchRef = useRef(null);
  const teamsSearchRef = useRef(null);
  const challengesSearchRef = useRef(null);

  const [selectedContainers, setSelectedContainers] = useState([]);
  const [isStopModalOpen, setIsStopModalOpen] = useState(false);
  const [showDeleted, setShowDeleted] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState(10);
  const [stats, setStats] = useState({ totalContainers: 0, runningContainers: 0, stoppedContainers: 0 });
  const { t, i18n } = useTranslation();

  const pageSize = 20;

  const fetchContainers = async (page = currentPage, deleted = showDeleted, activeFilters = filtersRef.current) => {
    try {
      const params = {
        ...activeFilters,
        limit: pageSize,
        offset: (page - 1) * pageSize,
        ...(deleted && { deleted: true }),
      };
      Object.keys(params).forEach((key) => {
        if (params[key] === '') delete params[key];
      });

      const response = await getVictims(params);
      if (response.code === 200) {
        setContainers(response.data.victims || []);
        setRunningCount(response.data.running || 0);
        const total = response.data.count || 0;
        const running = response.data.running || 0;
        setStats({ totalContainers: total, runningContainers: running, stoppedContainers: total - running });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.victims.toast.fetchContainersFailed') });
    }
  };

  const toggleShowDeleted = () => {
    const next = !showDeleted;
    setShowDeleted(next);
    setCurrentPage(1);
    setSelectedContainers([]);
    fetchContainers(1, next);
  };

  const currentPageRef = useRef(currentPage);
  const showDeletedRef = useRef(showDeleted);
  const filtersRef = useRef(filters);

  useEffect(() => {
    currentPageRef.current = currentPage;
  }, [currentPage]);
  useEffect(() => {
    showDeletedRef.current = showDeleted;
  }, [showDeleted]);
  useEffect(() => {
    filtersRef.current = filters;
  }, [filters]);

  useEffect(() => {
    fetchContainers();
  }, []);

  useEffect(() => {
    fetchContainers();
  }, [currentPage, filters.user_id, filters.team_id, filters.challenge_id]);

  useEffect(() => {
    if (refreshInterval <= 0) return;
    const id = setInterval(
      () => fetchContainers(currentPageRef.current, showDeletedRef.current, filtersRef.current),
      refreshInterval * 1000
    );
    return () => clearInterval(id);
  }, [refreshInterval]);

  useEffect(() => {
    const handleClickOutside = (event) => {
      const isOutsideUsers = usersSearchRef.current && !usersSearchRef.current.contains(event.target);
      const isOutsideTeams = teamsSearchRef.current && !teamsSearchRef.current.contains(event.target);
      const isOutsideChallenges = challengesSearchRef.current && !challengesSearchRef.current.contains(event.target);

      if (isOutsideUsers && searchResults.users.length > 0) setSearchResults((prev) => ({ ...prev, users: [] }));
      if (isOutsideTeams && searchResults.teams.length > 0) setSearchResults((prev) => ({ ...prev, teams: [] }));
      if (isOutsideChallenges && searchResults.challenges.length > 0)
        setSearchResults((prev) => ({ ...prev, challenges: [] }));
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [searchResults]);

  const handleFilterChange = (key, value) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setCurrentPage(1);
  };

  const handleSearch = async (model, name, setResults, setLoading) => {
    if (!name || name.trim() === '') {
      setResults([]);
      return;
    }
    setLoading(true);
    try {
      const response = await searchModels({ model, 'search[name]': name.trim(), limit: 10, offset: 0 });
      if (response.code === 200) {
        setResults(response.data.models || []);
      } else {
        setResults([]);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.victims.toast.searchFailed') });
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const debounceTimerRef = useRef(null);
  const debouncedSearch = (model, name, setResults, setLoading) => {
    clearTimeout(debounceTimerRef.current);
    debounceTimerRef.current = setTimeout(() => {
      handleSearch(model, name, setResults, setLoading);
    }, 300);
  };

  const handleResetFilters = () => {
    setFilters({ user_id: '', team_id: '', challenge_id: '', limit: 20, offset: 0 });
    setSearchResults({ users: [], teams: [], challenges: [] });
    setCurrentPage(1);
  };

  const handlePageChange = (page) => setCurrentPage(page);

  const handleContainerSelect = (containerId) => {
    setSelectedContainers((prev) =>
      prev.includes(containerId) ? prev.filter((id) => id !== containerId) : [...prev, containerId]
    );
  };

  const handleSelectAll = () => {
    if (selectedContainers.length === containers.length) {
      setSelectedContainers([]);
    } else {
      setSelectedContainers(containers.map((c) => c.id));
    }
  };

  const handleStopContainers = async () => {
    if (selectedContainers.length === 0) {
      toast.warning({ description: t('admin.victims.toast.selectStopRequired') });
      return;
    }
    try {
      const response = await stopVictims(selectedContainers);
      if (response.code === 200) {
        toast.success({ description: t('admin.victims.toast.taskDispatched') });
        setSelectedContainers([]);
        fetchContainers();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.victims.toast.taskDispatchFailed') });
    }
    setIsStopModalOpen(false);
  };

  const formatTime = (timeStr) => {
    if (!timeStr) return '-';
    return new Date(timeStr).toLocaleString(i18n.language || 'en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  const formatRemaining = (remaining) => {
    if (!remaining || remaining <= 0) return t('admin.victims.status.stopped');
    const hours = Math.floor(remaining / 3600);
    const minutes = Math.floor((remaining % 3600) / 60);
    const seconds = Math.floor(remaining % 60);
    return `${hours}h ${minutes}m ${seconds}s`;
  };

  const getContainerStatusStyle = (remaining) => {
    if (!remaining || remaining <= 0) return 'text-red-400 bg-red-400/10 border-red-400/30';
    return 'text-green-400 bg-green-400/10 border-green-400/30';
  };

  return (
    <div className="w-full mx-auto space-y-6">
      <div className="mb-8">
        <div className="mb-4">
          <p className="text-neutral-400 font-mono">{t('admin.victims.page.subtitle')}</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <StatCard
            title={t('admin.victims.stats.total')}
            value={stats.totalContainers}
            valueColor="text-neutral-50"
            icon={<IconServer size={20} className="text-geek-400" />}
          />
          <StatCard
            title={t('admin.victims.stats.running')}
            value={stats.runningContainers}
            valueColor="text-green-400"
            icon={<IconPlayerPlay size={20} className="text-green-400" />}
            iconBgClass="bg-green-400/20"
            delay={0.1}
          />
          <StatCard
            title={t('admin.victims.stats.stopped')}
            value={stats.stoppedContainers}
            valueColor="text-red-400"
            icon={<IconBan size={20} className="text-red-400" />}
            iconBgClass="bg-red-400/20"
            delay={0.2}
          />
        </div>
      </div>

      {/* Filters */}
      <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
        <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <IconFilter size={18} className="text-neutral-400" />
              <h3 className="text-base font-mono text-neutral-50">{t('admin.victims.filters.title')}</h3>
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleResetFilters}
              className="!text-neutral-400 hover:!text-neutral-300 !text-xs !h-6 !px-2"
            >
              {t('admin.victims.filters.reset')}
            </Button>
          </div>
          <div className="grid grid-cols-3 gap-3">
            {/* User search */}
            <div className="relative" ref={usersSearchRef}>
              <label className="block text-xs font-mono text-neutral-400 mb-1">
                {t('admin.victims.filters.userName')}
              </label>
              <div className="relative">
                <IconSearch size={14} className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400" />
                <input
                  type="text"
                  placeholder={t('admin.victims.filters.searchUserPlaceholder')}
                  onChange={(e) => {
                    const value = e.target.value;
                    debouncedSearch(
                      'User',
                      value,
                      (results) => setSearchResults((prev) => ({ ...prev, users: results })),
                      (loading) => setSearchLoading((prev) => ({ ...prev, users: loading }))
                    );
                  }}
                  className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md
                    text-xs text-neutral-50 placeholder-neutral-500
                    focus:outline-none focus:border-geek-400 focus:shadow-focus transition-all duration-200"
                />
                {searchLoading.users && (
                  <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
                    <div className="w-3 h-3 border border-geek-400 border-t-transparent rounded-full animate-spin" />
                  </div>
                )}
              </div>
              {searchResults.users.length > 0 && (
                <div className="dropdown-custom max-h-32">
                  {searchResults.users.map((user) => (
                    <div
                      key={user.id}
                      className="dropdown-option text-xs"
                      onClick={() => {
                        handleFilterChange('user_id', user.id.toString());
                        setSearchResults((prev) => ({ ...prev, users: [] }));
                      }}
                    >
                      {user.name || user.username || t('admin.victims.filters.userFallback', { id: user.id })}
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Team search */}
            <div className="relative" ref={teamsSearchRef}>
              <label className="block text-xs font-mono text-neutral-400 mb-1">
                {t('admin.victims.filters.teamName')}
              </label>
              <div className="relative">
                <IconUsers size={14} className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400" />
                <input
                  type="text"
                  placeholder={t('admin.victims.filters.searchTeamPlaceholder')}
                  onChange={(e) => {
                    const value = e.target.value;
                    debouncedSearch(
                      'Team',
                      value,
                      (results) => setSearchResults((prev) => ({ ...prev, teams: results })),
                      (loading) => setSearchLoading((prev) => ({ ...prev, teams: loading }))
                    );
                  }}
                  className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md
                    text-xs text-neutral-50 placeholder-neutral-500
                    focus:outline-none focus:border-geek-400 focus:shadow-focus transition-all duration-200"
                />
                {searchLoading.teams && (
                  <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
                    <div className="w-3 h-3 border border-geek-400 border-t-transparent rounded-full animate-spin" />
                  </div>
                )}
              </div>
              {searchResults.teams.length > 0 && (
                <div className="dropdown-custom max-h-32">
                  {searchResults.teams.map((team) => (
                    <div
                      key={team.id}
                      className="dropdown-option text-xs"
                      onClick={() => {
                        handleFilterChange('team_id', team.id.toString());
                        setSearchResults((prev) => ({ ...prev, teams: [] }));
                      }}
                    >
                      {team.name || t('admin.victims.filters.teamFallback', { id: team.id })}
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Challenge search */}
            <div className="relative" ref={challengesSearchRef}>
              <label className="block text-xs font-mono text-neutral-400 mb-1">
                {t('admin.victims.filters.challengeName')}
              </label>
              <div className="relative">
                <IconTarget size={14} className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400" />
                <input
                  type="text"
                  placeholder={t('admin.victims.filters.searchChallengePlaceholder')}
                  onChange={(e) => {
                    const value = e.target.value;
                    debouncedSearch(
                      'Challenge',
                      value,
                      (results) => setSearchResults((prev) => ({ ...prev, challenges: results })),
                      (loading) => setSearchLoading((prev) => ({ ...prev, challenges: loading }))
                    );
                  }}
                  className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md
                    text-xs text-neutral-50 placeholder-neutral-500
                    focus:outline-none focus:border-geek-400 focus:shadow-focus transition-all duration-200"
                />
                {searchLoading.challenges && (
                  <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
                    <div className="w-3 h-3 border border-geek-400 border-t-transparent rounded-full animate-spin" />
                  </div>
                )}
              </div>
              {searchResults.challenges.length > 0 && (
                <div className="dropdown-custom max-h-32">
                  {searchResults.challenges.map((challenge) => (
                    <div
                      key={challenge.id}
                      className="dropdown-option text-xs"
                      onClick={() => {
                        handleFilterChange('challenge_id', challenge.rand_id.toString());
                        setSearchResults((prev) => ({ ...prev, challenges: [] }));
                      }}
                    >
                      {challenge.name || t('admin.victims.filters.challengeFallback', { id: challenge.id })}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {(filters.user_id || filters.team_id || filters.challenge_id) && (
            <div className="mt-3 pt-3 border-t border-neutral-300/20">
              <div className="flex flex-wrap gap-2">
                {filters.user_id && (
                  <span className="px-2 py-1 bg-geek-400/20 text-geek-400 text-xs font-mono rounded border border-geek-400/30">
                    {t('admin.victims.filters.userIdLabel')}: {filters.user_id}
                    <button onClick={() => handleFilterChange('user_id', '')} className="ml-1 hover:text-red-400">
                      ×
                    </button>
                  </span>
                )}
                {filters.team_id && (
                  <span className="px-2 py-1 bg-geek-400/20 text-geek-400 text-xs font-mono rounded border border-geek-400/30">
                    {t('admin.victims.filters.teamIdLabel')}: {filters.team_id}
                    <button onClick={() => handleFilterChange('team_id', '')} className="ml-1 hover:text-red-400">
                      ×
                    </button>
                  </span>
                )}
                {filters.challenge_id && (
                  <span className="px-2 py-1 bg-green-400/20 text-green-400 text-xs font-mono rounded border border-green-400/30">
                    {t('admin.victims.filters.challengeIdLabel')}: {filters.challenge_id}
                    <button onClick={() => handleFilterChange('challenge_id', '')} className="ml-1 hover:text-red-400">
                      ×
                    </button>
                  </span>
                )}
              </div>
            </div>
          )}
        </div>
      </motion.div>

      {/* Container list */}

      {/* 工具栏 */}
      <div className="flex flex-wrap gap-2 items-center">
        <div className="flex items-center gap-1 px-2 h-8 rounded-md border border-neutral-700 bg-neutral-900">
          <IconClockPlay size={13} className="text-neutral-400 shrink-0" />
          <span className="text-xs text-neutral-400 shrink-0">{t('common.autoRefresh')}</span>
          <select
            value={refreshInterval}
            onChange={(e) => setRefreshInterval(Number(e.target.value))}
            className="bg-transparent text-xs text-neutral-300 outline-none cursor-pointer"
          >
            {[5, 10, 30, 60].map((s) => (
              <option key={s} value={s} className="bg-neutral-900">
                {s}s
              </option>
            ))}
            <option value={0} className="bg-neutral-900">
              {t('common.autoRefreshOff')}
            </option>
          </select>
        </div>
        <Button
          variant="ghost"
          size="sm"
          leftIcon={<IconRefresh size={14} />}
          onClick={() => fetchContainers(currentPage, showDeleted, filtersRef.current)}
        >
          {t('common.refresh')}
        </Button>
        <Button
          variant={showDeleted ? 'danger' : 'ghost'}
          size="sm"
          leftIcon={<IconTrash size={14} />}
          onClick={toggleShowDeleted}
        >
          {t('admin.victims.showDeleted')}
        </Button>
        {selectedContainers.length > 0 && (
          <Button variant="danger" size="sm" leftIcon={<IconBan size={14} />} onClick={() => setIsStopModalOpen(true)}>
            {t('admin.victims.table.stopButton')} ({selectedContainers.length})
          </Button>
        )}
      </div>

      <Card variant="default" padding="none" className="overflow-hidden">
        <div className="p-4 bg-black/20 border-b border-neutral-300/30 flex items-center gap-2">
          <IconTable size={20} className="text-neutral-400" />
          <h3 className="text-lg font-mono text-neutral-50">{t('admin.victims.table.title')}</h3>
          <span className="text-sm font-mono text-neutral-400">
            {t('admin.victims.table.total', { count: runningCount })}
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="bg-black/40">
                <th className="p-4 text-left text-neutral-400 font-mono">
                  <input
                    type="checkbox"
                    checked={selectedContainers.length === containers.length && containers.length > 0}
                    onChange={handleSelectAll}
                    className="w-4 h-4 rounded border-neutral-300/30 text-geek-400
                      focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                  />
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.id')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.contest')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.challenge')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.team')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.user')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.remote')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.startTime')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.status')}
                </th>
                <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap" scope="col">
                  {t('admin.victims.table.columns.remaining')}
                </th>
              </tr>
            </thead>
            <tbody>
              {containers.length === 0 ? (
                <tr>
                  <td colSpan="10">
                    <EmptyState title={t('admin.victims.table.empty')} />
                  </td>
                </tr>
              ) : (
                containers.map((container, index) => (
                  <motion.tr
                    key={container.id}
                    className="border-t border-neutral-300/10 hover:bg-black/40 transition-colors"
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.02 }}
                  >
                    <td className="p-4 text-neutral-300 font-mono">
                      <input
                        type="checkbox"
                        checked={selectedContainers.includes(container.id)}
                        onChange={() => handleContainerSelect(container.id)}
                        className="w-4 h-4 rounded border-neutral-300/30 text-geek-400
                          focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                      />
                    </td>
                    <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{container.id}</td>
                    <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">
                      <span className="px-2 py-0.5 rounded border border-geek-400/30 text-geek-400 text-xs">
                        {container.contest_id ?? '-'}
                      </span>
                    </td>
                    <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{container.challenge}</td>
                    <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{container.team}</td>
                    <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{container.user}</td>
                    <td className="p-4 text-neutral-300 font-mono">
                      {container.remote && container.remote.length > 0 ? (
                        <div className="space-y-1">
                          {container.remote.map((addr, i) => (
                            <div
                              key={i}
                              className="text-xs bg-black/30 text-geek-400 px-2 py-1 rounded border border-geek-400/30"
                            >
                              {addr}
                            </div>
                          ))}
                        </div>
                      ) : (
                        <span className="text-neutral-400">-</span>
                      )}
                    </td>
                    <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{formatTime(container.start)}</td>
                    <td className="p-4 whitespace-nowrap">
                      <VictimStatusBadge status={container.status} t={t} />
                    </td>
                    <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">
                      <span
                        className={`px-2 py-1 rounded-md text-xs font-mono border ${getContainerStatusStyle(container.remaining)}`}
                      >
                        {formatRemaining(container.remaining)}
                      </span>
                    </td>
                  </motion.tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {runningCount > 0 && (
          <div className="p-4 border-t border-neutral-300/30 bg-black/20 flex justify-center">
            <Pagination
              total={Math.ceil(runningCount / pageSize)}
              current={currentPage}
              pageSize={pageSize}
              onChange={handlePageChange}
              showTotal
              totalItems={runningCount}
            />
          </div>
        )}
      </Card>

      {/* Stop modal */}
      <Modal
        isOpen={isStopModalOpen}
        onClose={() => setIsStopModalOpen(false)}
        title={t('admin.victims.modals.stopTitle')}
        footer={
          <>
            <ModalButton onClick={() => setIsStopModalOpen(false)}>{t('common.cancel')}</ModalButton>
            <ModalButton variant="danger" onClick={handleStopContainers}>
              {t('admin.victims.modals.stopConfirm')}
            </ModalButton>
          </>
        }
      >
        <p className="text-neutral-300 font-mono">
          {t('admin.victims.modals.stopPrompt', { count: selectedContainers.length })}
        </p>
      </Modal>
    </div>
  );
}

export default AdminVictims;
