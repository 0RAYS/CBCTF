import { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import {
  getContestVictims,
  stopContestVictims,
  startContestVictims,
  getContestTeams,
  getContestChallenges,
} from '../../../api/admin/contest';
import { getUserList } from '../../../api/admin/user';
import { Modal } from '../../../components/common';
import ModalButton from '../../../components/common/ModalButton';
import { Button, Pagination, Card, EmptyState, StatCard, Chip } from '../../../components/common';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import {
  IconPlayerPlay,
  IconBan,
  IconFilter,
  IconTable,
  IconServer,
  IconUsers,
  IconTarget,
  IconSearch,
  IconTrash,
  IconClockPlay,
  IconRefresh,
  IconArrowsMaximize,
  IconChevronLeft,
  IconChevronRight,
} from '@tabler/icons-react';
import { getChallengeCategoryChipClass, getChallengeTypeChipClass } from '../../../config/challengeChips';

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
      {t(`admin.contests.containers.statusBadge.${status}`, status)}
    </span>
  );
}

function ContestContainers() {
  const { id: contestId } = useParams();

  // 容器列表相关状态
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

  // 选中的容器
  const [selectedContainers, setSelectedContainers] = useState([]);

  // 模态框状态
  const [isStartModalOpen, setIsStartModalOpen] = useState(false);
  const [isStopModalOpen, setIsStopModalOpen] = useState(false);
  const [isChallengeDetailsOpen, setIsChallengeDetailsOpen] = useState(false);

  // 开启容器相关状态
  const [challenges, setChallenges] = useState([]);
  const [detailChallenges, setDetailChallenges] = useState([]);
  const [selectedChallenges, setSelectedChallenges] = useState([]);
  const [challengeSearch, setChallengeSearch] = useState('');
  const [randomTeamPercentage, setRandomTeamPercentage] = useState(50); // 随机选择队伍的百分比
  const [victimDurationInput, setVictimDurationInput] = useState('7200');

  const challengePageSize = 20;
  const [challengePage, setChallengePage] = useState(1);
  const [challengeTotal, setChallengeTotal] = useState(0);
  const [detailChallengePage, setDetailChallengePage] = useState(1);
  const [detailChallengeTotal, setDetailChallengeTotal] = useState(0);
  const [teamTotal, setTeamTotal] = useState(0);

  // 统计信息
  const [stats, setStats] = useState({
    totalContainers: 0,
    runningContainers: 0,
    stoppedContainers: 0,
  });
  const [showDeleted, setShowDeleted] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState(10);
  const { t, i18n } = useTranslation();

  const pageSize = 20; // 增加每页显示数量
  const totalTeamCount = teamTotal;

  // 获取容器列表
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

      const response = await getContestVictims(parseInt(contestId), params);

      if (response.code === 200) {
        setContainers(response.data.victims || []);
        setRunningCount(response.data.running || 0);
        const total = response.data.count || 0;
        const running = response.data.running || 0;
        setStats({
          totalContainers: total,
          runningContainers: running,
          stoppedContainers: total - running,
        });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchContainersFailed') });
    }
  };

  const toggleShowDeleted = () => {
    const next = !showDeleted;
    setShowDeleted(next);
    setCurrentPage(1);
    setSelectedContainers([]);
    fetchContainers(1, next, filtersRef.current);
  };

  // 获取团队列表
  const fetchTeams = async () => {
    try {
      const response = await getContestTeams(parseInt(contestId), { limit: 1, offset: 0 });
      if (response.code === 200) {
        setTeamTotal(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchTeamsFailed') });
    }
  };

  const fetchChallenges = async (page = challengePage, query = challengeSearch) => {
    try {
      const params = {
        type: 'pods',
        limit: challengePageSize,
        offset: (page - 1) * challengePageSize,
      };
      if (query.trim() !== '') {
        params['search[name]'] = query.trim();
      }
      const response = await getContestChallenges(parseInt(contestId), params);
      if (response.code === 200) {
        setChallenges(response.data.challenges || []);
        setChallengeTotal(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchChallengesFailed') });
    }
  };

  const fetchDetailChallenges = async (page = detailChallengePage, query = challengeSearch) => {
    try {
      const params = {
        type: 'pods',
        limit: challengePageSize,
        offset: (page - 1) * challengePageSize,
      };
      if (query.trim() !== '') {
        params['search[name]'] = query.trim();
      }
      const response = await getContestChallenges(parseInt(contestId), params);
      if (response.code === 200) {
        setDetailChallenges(response.data.challenges || []);
        setDetailChallengeTotal(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchChallengesFailed') });
    }
  };

  const openChallengeDetails = () => {
    setDetailChallengePage(1);
    setIsChallengeDetailsOpen(true);
  };

  const updateChallengeSelection = (challengeId, checked) => {
    setSelectedChallenges((prev) => {
      if (checked) {
        return prev.includes(challengeId) ? prev : [...prev, challengeId];
      }
      return prev.filter((id) => id !== challengeId);
    });
  };

  const handleChallengeSearchChange = (value) => {
    setChallengeSearch(value);
    setChallengePage(1);
    setDetailChallengePage(1);
  };

  useEffect(() => {
    fetchContainers();
    fetchTeams();
    fetchChallenges(1, challengeSearch);
    setChallengePage(1);
    setDetailChallengePage(1);
  }, [contestId]);

  useEffect(() => {
    fetchContainers();
  }, [currentPage, filters.user_id, filters.team_id, filters.challenge_id]);

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
    if (refreshInterval <= 0) return;
    const id = setInterval(
      () => fetchContainers(currentPageRef.current, showDeletedRef.current, filtersRef.current),
      refreshInterval * 1000
    );
    return () => clearInterval(id);
  }, [refreshInterval]);

  useEffect(() => {
    fetchChallenges(challengePage, challengeSearch);
  }, [challengePage, challengeSearch]);

  useEffect(() => {
    if (!isChallengeDetailsOpen) return;
    fetchDetailChallenges(detailChallengePage, challengeSearch);
  }, [isChallengeDetailsOpen, detailChallengePage, challengeSearch]);

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
      const keyword = name.trim();
      if (model === 'User') {
        const response = await getUserList({ name: keyword, limit: 10, offset: 0 });
        setResults(response.code === 200 ? response.data.users || [] : []);
        return;
      }
      if (model === 'Team') {
        const response = await getContestTeams(parseInt(contestId, 10), { name: keyword, limit: 10, offset: 0 });
        setResults(response.code === 200 ? response.data.teams || [] : []);
        return;
      }
      if (model === 'Challenge') {
        const response = await getContestChallenges(parseInt(contestId, 10), {
          name: keyword,
          type: 'pods',
          limit: 10,
          offset: 0,
        });
        if (response.code !== 200) {
          setResults([]);
          return;
        }
        const results = (response.data.challenges || [])
          .filter((challenge) => challenge.name?.toLowerCase().includes(keyword.toLowerCase()))
          .slice(0, 10);
        setResults(results);
        return;
      }
      setResults([]);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.searchFailed') });
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

  // 处理页面切换
  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  // 处理容器选择
  const handleContainerSelect = (containerId) => {
    setSelectedContainers((prev) => {
      if (prev.includes(containerId)) {
        return prev.filter((id) => id !== containerId);
      } else {
        return [...prev, containerId];
      }
    });
  };

  // 全选/取消全选
  const handleSelectAll = () => {
    if (selectedContainers.length === containers.length) {
      setSelectedContainers([]);
    } else {
      setSelectedContainers(containers.map((c) => c.id));
    }
  };

  const isTeamRatioValid = randomTeamPercentage > 0 && randomTeamPercentage < 100;
  const victimDurationSeconds = Number.parseInt(victimDurationInput, 10) || 0;
  const isVictimDurationValid = victimDurationSeconds > 0;
  const selectedTeamCount =
    totalTeamCount > 0 && isTeamRatioValid ? Math.max(1, Math.floor((totalTeamCount * randomTeamPercentage) / 100)) : 0;
  const typeLabels = {
    static: t('admin.challenge.types.static'),
    question: t('admin.challenge.types.question'),
    dynamic: t('admin.challenge.types.dynamic'),
    pods: t('admin.challenge.types.pods'),
  };

  const formatVictimDuration = (seconds) => {
    if (!seconds || seconds <= 0) {
      return t('admin.contests.containers.quickActions.invalidDuration');
    }
    if (seconds < 60) {
      return t('utils.time.units.second', { count: seconds });
    }
    if (seconds < 3600) {
      const minutes = Math.floor(seconds / 60);
      const remainingSeconds = seconds % 60;
      return remainingSeconds > 0
        ? `${t('utils.time.units.minute', { count: minutes })}${t('utils.time.units.second', { count: remainingSeconds })}`
        : t('utils.time.units.minute', { count: minutes });
    }
    if (seconds < 86400) {
      const hours = Math.floor(seconds / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);
      return minutes > 0
        ? `${t('utils.time.units.hour', { count: hours })}${t('utils.time.units.minute', { count: minutes })}`
        : t('utils.time.units.hour', { count: hours });
    }
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    return hours > 0
      ? `${t('utils.time.units.day', { count: days })}${t('utils.time.units.hour', { count: hours })}`
      : t('utils.time.units.day', { count: days });
  };

  // 停止容器
  const handleStopContainers = async () => {
    if (selectedContainers.length === 0) {
      toast.warning({ description: t('admin.contests.containers.toast.selectStopRequired') });
      return;
    }

    try {
      const response = await stopContestVictims(parseInt(contestId), selectedContainers);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.containers.toast.taskDispatched') });
        setSelectedContainers([]);
        fetchContainers();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.taskDispatchFailed') });
    }
    setIsStopModalOpen(false);
  };

  // 开启容器
  const handleStartContainers = async () => {
    if (!isVictimDurationValid) {
      toast.warning({ description: t('admin.contests.containers.toast.invalidDuration') });
      return;
    }
    if (selectedChallenges.length === 0 || selectedTeamCount === 0) {
      toast.warning({ description: t('admin.contests.containers.toast.selectStartRequired') });
      return;
    }

    try {
      const response = await startContestVictims(
        parseInt(contestId),
        selectedChallenges,
        randomTeamPercentage / 100,
        victimDurationSeconds
      );
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.containers.toast.taskDispatched') });
        setSelectedChallenges([]);
        fetchContainers();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.taskDispatchFailed') });
    }
    setIsStartModalOpen(false);
  };

  // 格式化时间
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

  // 格式化剩余时间
  const formatRemaining = (remaining) => {
    if (!remaining || remaining <= 0) return t('admin.contests.containers.status.stopped');
    const hours = Math.floor(remaining / 3600);
    const minutes = Math.floor((remaining % 3600) / 60);
    const seconds = Math.floor(remaining % 60);
    return `${hours}h ${minutes}m ${seconds}s`;
  };

  // 获取容器状态样式
  const getContainerStatusStyle = (remaining) => {
    if (!remaining || remaining <= 0) {
      return 'text-red-400 bg-red-400/10 border-red-400/30';
    }
    return 'text-green-400 bg-green-400/10 border-green-400/30';
  };

  return (
    <>
      <style>
        {`
          .slider::-webkit-slider-thumb {
            appearance: none;
            height: 12px;
            width: 12px;
            border-radius: 50%;
            background: #597ef7;
            cursor: pointer;
            border: 2px solid #1f2937;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
          }
          
          .slider::-moz-range-thumb {
            height: 12px;
            width: 12px;
            border-radius: 50%;
            background: #597ef7;
            cursor: pointer;
            border: 2px solid #1f2937;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
          }
          
          .slider::-webkit-slider-track {
            background: transparent;
          }
          
          .slider::-moz-range-track {
            background: transparent;
          }
        `}
      </style>
      <div className="w-full mx-auto space-y-6">
        {/* 页面标题和统计信息 */}
        <div className="mb-8">
          <div className="mb-4">
            <p className="text-neutral-400 font-mono">{t('admin.contests.containers.page.subtitle')}</p>
          </div>

          {/* 统计卡片 */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <StatCard
              title={t('admin.contests.containers.stats.total')}
              value={stats.totalContainers}
              valueColor="text-neutral-50"
              icon={<IconServer size={20} className="text-geek-400" />}
            />
            <StatCard
              title={t('admin.contests.containers.stats.running')}
              value={stats.runningContainers}
              valueColor="text-green-400"
              icon={<IconPlayerPlay size={20} className="text-green-400" />}
              iconBgClass="bg-green-400/20"
              delay={0.1}
            />
            <StatCard
              title={t('admin.contests.containers.stats.stopped')}
              value={stats.stoppedContainers}
              valueColor="text-red-400"
              icon={<IconBan size={20} className="text-red-400" />}
              iconBgClass="bg-red-400/20"
              delay={0.2}
            />
          </div>
        </div>

        {/* 快速操作 */}
        <motion.div
          className="grid grid-cols-1 gap-6 items-stretch"
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
        >
          {/* 快速操作区域 */}
          <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4 min-h-[460px] flex flex-col">
            <div className="flex items-center gap-2 mb-3">
              <IconPlayerPlay size={18} className="text-neutral-400" />
              <h3 className="text-base font-mono text-neutral-50">
                {t('admin.contests.containers.quickActions.title')}
              </h3>
            </div>

            <div className="flex flex-col gap-4 flex-1 min-h-0">
              {/* 选择队伍 */}
              <div className="border border-neutral-300/20 rounded-md bg-black/10 p-4">
                <div className="flex justify-between items-center mb-2">
                  <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
                    <IconUsers size={14} />
                    <span className="text-xs font-mono text-neutral-400">
                      {t('admin.contests.containers.quickActions.randomTeams')}
                    </span>
                    <span className="text-xs font-mono text-geek-400">{randomTeamPercentage}%</span>
                  </label>
                  <span className="text-xs font-mono text-neutral-500">
                    {t('admin.contests.containers.quickActions.estimatedTeams', { count: selectedTeamCount })}
                  </span>
                </div>

                <div className="mb-3 p-2 border border-neutral-300/20 rounded-md bg-black/10">
                  <div className="relative">
                    <input
                      type="range"
                      min="0"
                      max="100"
                      value={randomTeamPercentage}
                      onChange={(e) => setRandomTeamPercentage(parseInt(e.target.value))}
                      className="w-full h-1 bg-neutral-700 rounded-lg appearance-none cursor-pointer slider"
                      style={{
                        background: `linear-gradient(to right, #597ef7 0%, #597ef7 ${randomTeamPercentage}%, #374151 ${randomTeamPercentage}%, #374151 100%)`,
                      }}
                    />
                  </div>
                </div>

                <div className="border border-neutral-300/30 rounded-md bg-black/10 p-3">
                  <p className="text-xs font-mono text-neutral-400">
                    {t('admin.contests.containers.quickActions.teamSelectionHint', { total: totalTeamCount })}
                  </p>
                </div>
              </div>

              <div className="border border-neutral-300/20 rounded-md bg-black/10 p-4">
                <div className="flex justify-between items-center gap-3 mb-2">
                  <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
                    <IconClockPlay size={14} />
                    <span>{t('common.duration')}</span>
                  </label>
                  <span className="text-xs font-mono text-geek-400">
                    {t('admin.contests.containers.quickActions.durationPreview', {
                      value: formatVictimDuration(victimDurationSeconds),
                    })}
                  </span>
                </div>

                <div className="flex items-center gap-3">
                  <input
                    type="number"
                    min="1"
                    step="1"
                    value={victimDurationInput}
                    onChange={(e) => setVictimDurationInput(e.target.value)}
                    className="w-full h-9 px-3 bg-black/20 border border-neutral-300/30 rounded-md text-sm text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 transition-all duration-200"
                  />
                  <span className="text-xs font-mono text-neutral-400 whitespace-nowrap">
                    {t('admin.contests.containers.quickActions.durationUnit')}
                  </span>
                </div>

                <p className="mt-2 text-xs font-mono text-neutral-400">
                  {t('admin.contests.containers.quickActions.durationHint')}
                </p>
              </div>

              {/* 选择题目 */}
              <div className="flex flex-col min-h-0 flex-1">
                <div className="flex justify-between items-center mb-2">
                  <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
                    <IconTarget size={14} />
                    {t('admin.contests.containers.quickActions.selectChallenges')}
                  </label>
                  <div className="flex gap-1">
                    <Button variant="ghost" size="sm" onClick={openChallengeDetails} className="!text-xs !h-5 !px-1">
                      <span className="inline-flex items-center gap-1">
                        <IconArrowsMaximize size={12} />
                        {t('admin.contests.containers.quickActions.expand')}
                      </span>
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setSelectedChallenges(challenges.map((c) => c.id))}
                      className="!text-xs !h-5 !px-1"
                    >
                      {t('admin.contests.containers.quickActions.selectAll')}
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setSelectedChallenges([])}
                      className="!text-xs !h-5 !px-1"
                    >
                      {t('admin.contests.containers.quickActions.clear')}
                    </Button>
                  </div>
                </div>
                <div className="flex-1 min-h-[260px] overflow-y-auto border border-neutral-300/30 rounded-md bg-black/10">
                  <div className="p-2 border-b border-neutral-300/20">
                    <div className="relative">
                      <IconSearch
                        size={12}
                        className="absolute left-2 top-1/2 -translate-y-1/2 text-neutral-500 pointer-events-none"
                      />
                      <input
                        type="text"
                        value={challengeSearch}
                        onChange={(e) => handleChallengeSearchChange(e.target.value)}
                        placeholder={t('admin.contests.containers.quickActions.searchPlaceholder')}
                        className="w-full h-7 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md text-xs text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 transition-all duration-200"
                      />
                    </div>
                  </div>
                  {challenges.length > 0 ? (
                    challenges.map((challenge) => (
                      <div key={challenge.id} className="flex items-center p-1 hover:bg-black/30 transition-colors">
                        <input
                          type="checkbox"
                          id={`challenge-${challenge.id}`}
                          checked={selectedChallenges.includes(challenge.id)}
                          onChange={(e) => updateChallengeSelection(challenge.id, e.target.checked)}
                          className="w-3 h-3 rounded border-neutral-300/30 text-geek-400
                              focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                        />
                        <label
                          htmlFor={`challenge-${challenge.id}`}
                          className="ml-2 text-xs font-mono text-neutral-300 cursor-pointer flex-1 truncate"
                        >
                          {challenge.name}
                        </label>
                      </div>
                    ))
                  ) : (
                    <div className="p-3 text-xs font-mono text-neutral-500">
                      {t('admin.contests.containers.quickActions.noChallenges')}
                    </div>
                  )}
                </div>
                {Math.ceil(challengeTotal / challengePageSize) > 1 && (
                  <div className="flex items-center justify-between gap-2 mt-2 px-1">
                    <span className="text-[11px] font-mono text-geek-400/80 whitespace-nowrap">
                      {t('admin.contests.containers.quickActions.pageHint')}
                    </span>
                    <div className="flex items-center gap-2 ml-auto">
                      <button
                        disabled={challengePage === 1}
                        onClick={() => setChallengePage((p) => p - 1)}
                        className="inline-flex items-center justify-center w-7 h-7 rounded-md border border-geek-400/40 bg-geek-400/10 text-geek-300 hover:bg-geek-400/20 hover:text-geek-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                        aria-label={t('common.previous')}
                      >
                        <IconChevronLeft size={15} />
                      </button>
                      <span className="text-xs font-mono text-neutral-300 min-w-[56px] text-center">
                        {challengePage} / {Math.ceil(challengeTotal / challengePageSize)}
                      </span>
                      <button
                        disabled={challengePage >= Math.ceil(challengeTotal / challengePageSize)}
                        onClick={() => setChallengePage((p) => p + 1)}
                        className="inline-flex items-center justify-center w-7 h-7 rounded-md border border-geek-400/40 bg-geek-400/10 text-geek-300 hover:bg-geek-400/20 hover:text-geek-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                        aria-label={t('common.next')}
                      >
                        <IconChevronRight size={15} />
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* 开启容器按钮 */}
            <div className="mt-4 flex justify-end">
              <Button
                variant="primary"
                size="sm"
                align="icon-left"
                icon={<IconPlayerPlay size={14} />}
                onClick={() => setIsStartModalOpen(true)}
                disabled={selectedChallenges.length === 0 || selectedTeamCount === 0 || !isVictimDurationValid}
                className="!text-xs !h-7 !px-3"
              >
                {t('admin.contests.containers.quickActions.startButton', {
                  challenges: selectedChallenges.length,
                  teams: selectedTeamCount,
                })}
              </Button>
            </div>
          </div>
        </motion.div>

        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
          <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <IconFilter size={18} className="text-neutral-400" />
                <h3 className="text-base font-mono text-neutral-50">{t('admin.contests.containers.filters.title')}</h3>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleResetFilters}
                className="!text-neutral-400 hover:!text-neutral-300 !text-xs !h-6 !px-2"
              >
                {t('admin.contests.containers.filters.reset')}
              </Button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              <div className="relative" ref={usersSearchRef}>
                <label className="block text-xs font-mono text-neutral-400 mb-1">
                  {t('admin.contests.containers.filters.userName')}
                </label>
                <div className="relative">
                  <IconSearch
                    size={14}
                    className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400"
                  />
                  <input
                    type="text"
                    placeholder={t('admin.contests.containers.filters.searchUserPlaceholder')}
                    onChange={(e) => {
                      const value = e.target.value;
                      debouncedSearch(
                        'User',
                        value,
                        (results) => setSearchResults((prev) => ({ ...prev, users: results })),
                        (loading) => setSearchLoading((prev) => ({ ...prev, users: loading }))
                      );
                    }}
                    className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md text-xs text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 focus:shadow-focus transition-all duration-200"
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
                        {user.name ||
                          user.username ||
                          t('admin.contests.containers.filters.userFallback', { id: user.id })}
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div className="relative" ref={teamsSearchRef}>
                <label className="block text-xs font-mono text-neutral-400 mb-1">
                  {t('admin.contests.containers.filters.teamName')}
                </label>
                <div className="relative">
                  <IconUsers
                    size={14}
                    className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400"
                  />
                  <input
                    type="text"
                    placeholder={t('admin.contests.containers.filters.searchTeamPlaceholder')}
                    onChange={(e) => {
                      const value = e.target.value;
                      debouncedSearch(
                        'Team',
                        value,
                        (results) => setSearchResults((prev) => ({ ...prev, teams: results })),
                        (loading) => setSearchLoading((prev) => ({ ...prev, teams: loading }))
                      );
                    }}
                    className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md text-xs text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 focus:shadow-focus transition-all duration-200"
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
                        {team.name || t('admin.contests.containers.filters.teamFallback', { id: team.id })}
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div className="relative" ref={challengesSearchRef}>
                <label className="block text-xs font-mono text-neutral-400 mb-1">
                  {t('admin.contests.containers.filters.challengeName')}
                </label>
                <div className="relative">
                  <IconTarget
                    size={14}
                    className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400"
                  />
                  <input
                    type="text"
                    placeholder={t('admin.contests.containers.filters.searchChallengePlaceholder')}
                    onChange={(e) => {
                      const value = e.target.value;
                      debouncedSearch(
                        'Challenge',
                        value,
                        (results) => setSearchResults((prev) => ({ ...prev, challenges: results })),
                        (loading) => setSearchLoading((prev) => ({ ...prev, challenges: loading }))
                      );
                    }}
                    className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md text-xs text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 focus:shadow-focus transition-all duration-200"
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
                          handleFilterChange('challenge_id', challenge.id.toString());
                          setSearchResults((prev) => ({ ...prev, challenges: [] }));
                        }}
                      >
                        {challenge.name ||
                          t('admin.contests.containers.filters.challengeFallback', { id: challenge.id })}
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
                      {t('admin.contests.containers.filters.userIdLabel')}: {filters.user_id}
                      <button onClick={() => handleFilterChange('user_id', '')} className="ml-1 hover:text-red-400">
                        ×
                      </button>
                    </span>
                  )}
                  {filters.team_id && (
                    <span className="px-2 py-1 bg-geek-400/20 text-geek-400 text-xs font-mono rounded border border-geek-400/30">
                      {t('admin.contests.containers.filters.teamIdLabel')}: {filters.team_id}
                      <button onClick={() => handleFilterChange('team_id', '')} className="ml-1 hover:text-red-400">
                        ×
                      </button>
                    </span>
                  )}
                  {filters.challenge_id && (
                    <span className="px-2 py-1 bg-green-400/20 text-green-400 text-xs font-mono rounded border border-green-400/30">
                      {t('admin.contests.containers.filters.challengeIdLabel')}: {filters.challenge_id}
                      <button
                        onClick={() => handleFilterChange('challenge_id', '')}
                        className="ml-1 hover:text-red-400"
                      >
                        ×
                      </button>
                    </span>
                  )}
                </div>
              </div>
            )}
          </div>
        </motion.div>

        <Card variant="default" padding="none" className="overflow-hidden">
          {/* 列表头部 */}
          <div className="p-4 bg-black/20 border-b border-neutral-300/30 space-y-4">
            <div className="flex flex-wrap items-center gap-2">
              <IconTable size={20} className="text-neutral-400" />
              <h3 className="text-lg font-mono text-neutral-50">{t('admin.contests.containers.table.title')}</h3>
              <span className="text-sm font-mono text-neutral-400">
                {t('admin.contests.containers.table.total', { count: runningCount })}
              </span>
            </div>

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
                {t('admin.contests.containers.showDeleted')}
              </Button>
              {selectedContainers.length > 0 && (
                <Button
                  variant="danger"
                  size="sm"
                  leftIcon={<IconBan size={14} />}
                  onClick={() => setIsStopModalOpen(true)}
                >
                  {t('admin.contests.containers.table.stopButton')} ({selectedContainers.length})
                </Button>
              )}
            </div>
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
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.id')}
                  </th>
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.challenge')}
                  </th>
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.team')}
                  </th>
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.user')}
                  </th>
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.remote')}
                  </th>
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.startTime')}
                  </th>
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.status')}
                  </th>
                  <th className="p-4 text-left text-neutral-400 font-mono whitespace-nowrap">
                    {t('admin.contests.containers.table.columns.remaining')}
                  </th>
                </tr>
              </thead>
              <tbody>
                {containers.length === 0 ? (
                  <tr>
                    <td colSpan="10">
                      <EmptyState title={t('admin.contests.containers.table.empty')} />
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
                      whileHover={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}
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
                      <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{container.challenge}</td>
                      <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{container.team}</td>
                      <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">{container.user}</td>
                      <td className="p-4 text-neutral-300 font-mono">
                        {container.remote && container.remote.length > 0 ? (
                          <div className="space-y-1">
                            {container.remote.map((addr, index) => (
                              <div
                                key={index}
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
                      <td className="p-4 text-neutral-300 font-mono whitespace-nowrap">
                        {formatTime(container.start)}
                      </td>
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

          {/* 分页 */}
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

        {/* 开启容器确认模态框 */}
        <Modal
          isOpen={isStartModalOpen}
          onClose={() => setIsStartModalOpen(false)}
          title={t('admin.contests.containers.modals.startTitle')}
          footer={
            <>
              <ModalButton onClick={() => setIsStartModalOpen(false)}>{t('common.cancel')}</ModalButton>
              <ModalButton variant="primary" onClick={handleStartContainers}>
                {t('admin.contests.containers.modals.startConfirm')}
              </ModalButton>
            </>
          }
        >
          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <IconPlayerPlay size={20} className="text-geek-400" />
              <p className="text-neutral-300 font-mono">{t('admin.contests.containers.modals.startPrompt')}</p>
            </div>

            <div className="space-y-4">
              <div>
                <h4 className="text-sm font-mono text-neutral-400 mb-2">
                  {t('admin.contests.containers.modals.selectedChallenges', { count: selectedChallenges.length })}
                </h4>
                <div className="max-h-32 overflow-y-auto border border-neutral-300/30 rounded-md bg-black/10 p-2">
                  {selectedChallenges.map((challengeId) => {
                    const challenge = challenges.find((c) => c.id === challengeId);
                    return challenge ? (
                      <div key={challengeId} className="text-sm font-mono text-geek-400 py-1">
                        • {challenge.name}
                      </div>
                    ) : null;
                  })}
                </div>
              </div>

              <div>
                <h4 className="text-sm font-mono text-neutral-400 mb-2">
                  {t('admin.contests.containers.modals.teamRatioTitle')}
                </h4>
                <div className="border border-neutral-300/30 rounded-md bg-black/10 p-3 space-y-2">
                  <p className="text-sm font-mono text-geek-400">
                    {t('admin.contests.containers.modals.teamRatioValue', { ratio: randomTeamPercentage })}
                  </p>
                  <p className="text-xs font-mono text-neutral-400">
                    {t('admin.contests.containers.modals.teamRatioHint', {
                      count: selectedTeamCount,
                      total: totalTeamCount,
                    })}
                  </p>
                </div>
              </div>

              <div>
                <h4 className="text-sm font-mono text-neutral-400 mb-2">
                  {t('admin.contests.containers.modals.durationTitle')}
                </h4>
                <div className="border border-neutral-300/30 rounded-md bg-black/10 p-3 space-y-2">
                  <p className="text-sm font-mono text-geek-400">
                    {t('admin.contests.containers.modals.durationValue', { seconds: victimDurationSeconds })}
                  </p>
                  <p className="text-xs font-mono text-neutral-400">
                    {t('admin.contests.containers.modals.durationHint', {
                      value: formatVictimDuration(victimDurationSeconds),
                    })}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-neutral-800/50 border border-neutral-600/30 rounded-md p-3">
              <p className="text-xs font-mono text-neutral-400">
                {t('admin.contests.containers.modals.summaryPrefix')}
                <span className="text-geek-400">{selectedChallenges.length}</span>
                {t('admin.contests.containers.modals.summaryMiddle')}
                <span className="text-geek-400">{selectedTeamCount}</span>
                {t('admin.contests.containers.modals.summaryEquals')}
                <span className="text-green-400"> {selectedChallenges.length * selectedTeamCount}</span>
                {t('admin.contests.containers.modals.summarySuffix')}
              </p>
            </div>

            <div className="border border-amber-400/40 rounded-md bg-amber-400/10 p-3">
              <p className="text-xs font-mono text-amber-200">{t('admin.contests.containers.modals.startWarning')}</p>
            </div>
          </div>
        </Modal>

        <Modal
          isOpen={isChallengeDetailsOpen}
          onClose={() => setIsChallengeDetailsOpen(false)}
          title={t('admin.contests.containers.modals.challengeDetailsTitle')}
          size="xl"
          footer={
            <>
              <ModalButton onClick={() => setIsChallengeDetailsOpen(false)}>{t('common.cancel')}</ModalButton>
            </>
          }
        >
          <div className="space-y-4">
            <div className="relative">
              <IconSearch
                size={14}
                className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-500 pointer-events-none"
              />
              <input
                type="text"
                value={challengeSearch}
                onChange={(e) => handleChallengeSearchChange(e.target.value)}
                placeholder={t('admin.contests.containers.quickActions.searchPlaceholder')}
                className="w-full h-10 pl-10 pr-3 bg-black/20 border border-neutral-300/30 rounded-md text-sm text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 transition-all duration-200"
              />
            </div>

            <p className="text-sm font-mono text-neutral-400">
              {t('admin.contests.containers.modals.challengeDetailsHint', { count: detailChallengeTotal })}
            </p>

            {detailChallenges.length === 0 ? (
              <div className="border border-neutral-300/20 rounded-md bg-black/10 p-4 text-sm font-mono text-neutral-500">
                {t('admin.contests.containers.modals.challengeDetailsEmpty')}
              </div>
            ) : (
              <div className="space-y-4">
                {detailChallenges.map((challenge) => {
                  const isSelected = selectedChallenges.includes(challenge.id);
                  return (
                    <div
                      key={challenge.id}
                      className="border border-neutral-300/20 rounded-md bg-black/10 p-4 space-y-4"
                    >
                      <div className="flex items-start justify-between gap-4">
                        <div className="min-w-0 space-y-2">
                          <div className="flex items-center gap-2 flex-wrap">
                            <h4 className="text-base font-mono text-neutral-50 break-all">{challenge.name}</h4>
                            {challenge.category ? (
                              <Chip
                                size="sm"
                                label={challenge.category}
                                colorClass={getChallengeCategoryChipClass(challenge.category)}
                              />
                            ) : null}
                            {challenge.type ? (
                              <Chip
                                size="sm"
                                label={typeLabels[challenge.type] || challenge.type}
                                colorClass={getChallengeTypeChipClass(challenge.type)}
                              />
                            ) : null}
                            {challenge.hidden ? (
                              <Chip
                                size="sm"
                                label={t('admin.contests.challenges.hidden')}
                                colorClass="bg-red-400/20 text-red-400"
                              />
                            ) : null}
                          </div>
                          <div className="flex items-center gap-3 flex-wrap text-xs font-mono text-neutral-400">
                            <span>
                              {t('admin.contests.containers.modals.challengeScore', { score: challenge.score || 0 })}
                            </span>
                            <span>
                              {t('admin.contests.containers.modals.challengeSolvers', {
                                count: challenge.solvers || 0,
                              })}
                            </span>
                            <span>
                              {t('admin.contests.containers.modals.challengeAttempts', {
                                count: challenge.attempt || 0,
                              })}
                            </span>
                            <span>ID: {challenge.id}</span>
                          </div>
                        </div>

                        <label className="inline-flex items-center gap-2 text-sm font-mono text-neutral-300 cursor-pointer shrink-0">
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={(e) => updateChallengeSelection(challenge.id, e.target.checked)}
                            className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                          />
                          {t('admin.contests.containers.modals.selectChallenge')}
                        </label>
                      </div>

                      {challenge.type === 'question' ? (
                        <div className="space-y-2">
                          <h5 className="text-xs font-mono uppercase tracking-wide text-neutral-400">
                            {t('admin.contests.containers.modals.challengeOptions')}
                          </h5>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                            {(challenge.options || []).length > 0 ? (
                              challenge.options.map((option, index) => (
                                <div
                                  key={`${challenge.id}-option-${option.rand_id || index}`}
                                  className="rounded-md border border-neutral-300/20 bg-black/20 p-3 text-sm text-neutral-300 break-words"
                                >
                                  <div className="text-xs font-mono text-neutral-500 mb-1">
                                    {t('admin.contests.containers.modals.challengeOptionLabel', { index: index + 1 })}
                                  </div>
                                  <div>{option.content}</div>
                                </div>
                              ))
                            ) : (
                              <span className="text-xs font-mono text-neutral-500">
                                {t('admin.contests.containers.modals.challengeOptionsEmpty')}
                              </span>
                            )}
                          </div>
                        </div>
                      ) : null}
                    </div>
                  );
                })}
              </div>
            )}

            {Math.ceil(detailChallengeTotal / challengePageSize) > 1 ? (
              <div className="pt-2 border-t border-neutral-300/20">
                <Pagination
                  total={Math.ceil(detailChallengeTotal / challengePageSize)}
                  current={detailChallengePage}
                  pageSize={challengePageSize}
                  onChange={setDetailChallengePage}
                  showTotal
                  totalItems={detailChallengeTotal}
                />
              </div>
            ) : null}
          </div>
        </Modal>

        {/* 停止容器确认模态框 */}
        <Modal
          isOpen={isStopModalOpen}
          onClose={() => setIsStopModalOpen(false)}
          title={t('admin.contests.containers.modals.stopTitle')}
          footer={
            <>
              <ModalButton onClick={() => setIsStopModalOpen(false)}>{t('common.cancel')}</ModalButton>
              <ModalButton variant="danger" onClick={handleStopContainers}>
                {t('admin.contests.containers.modals.stopConfirm')}
              </ModalButton>
            </>
          }
        >
          <div className="flex items-center gap-3">
            <IconBan size={20} className="text-red-400" />
            <p className="text-neutral-300 font-mono">
              {t('admin.contests.containers.modals.stopPrompt', { count: selectedContainers.length })}
            </p>
          </div>
        </Modal>
      </div>
    </>
  );
}

export default ContestContainers;
