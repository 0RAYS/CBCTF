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
import { Modal } from '../../../components/common';
import ModalButton from '../../../components/common/ModalButton';
import { Button, Pagination, Card, EmptyState } from '../../../components/common';
import { motion } from 'motion/react';
import {
  IconPlayerPlay,
  IconBan,
  IconFilter,
  IconTable,
  IconServer,
  IconUsers,
  IconTarget,
  IconSearch,
  IconRefresh,
} from '@tabler/icons-react';
import { useWebSocket } from '../../../components/common/WebSocketProvider.jsx';
import { useTranslation } from 'react-i18next';
import { searchModels } from '../../../api/admin/search.js';

function ContestContainers() {
  const { id: contestId } = useParams();

  // 容器列表相关状态
  const [containers, setContainers] = useState([]);
  const [runningCount, setRunningCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);

  // 过滤参数
  const [filters, setFilters] = useState({
    user_id: '',
    team_id: '',
    challenge_id: '',
    limit: 20,
    offset: 0,
  });

  // 搜索相关状态
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

  // 搜索输入框refs
  const usersSearchRef = useRef(null);
  const teamsSearchRef = useRef(null);
  const challengesSearchRef = useRef(null);

  // 选中的容器
  const [selectedContainers, setSelectedContainers] = useState([]);

  // 模态框状态
  const [isStartModalOpen, setIsStartModalOpen] = useState(false);
  const [isStopModalOpen, setIsStopModalOpen] = useState(false);

  // 开启容器相关状态
  const [teams, setTeams] = useState([]);
  const [challenges, setChallenges] = useState([]);
  const [selectedTeams, setSelectedTeams] = useState([]);
  const [selectedChallenges, setSelectedChallenges] = useState([]);
  const [randomTeamPercentage, setRandomTeamPercentage] = useState(50); // 随机选择队伍的百分比

  // 统计信息
  const [stats, setStats] = useState({
    totalContainers: 0,
    runningContainers: 0,
    stoppedContainers: 0,
  });
  const { addMessageHandler } = useWebSocket();
  const { t, i18n } = useTranslation();

  const pageSize = 20; // 增加每页显示数量

  // 获取容器列表
  const fetchContainers = async () => {
    try {
      const params = {
        ...filters,
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      };
      // 清除空值
      Object.keys(params).forEach((key) => {
        if (params[key] === '') {
          delete params[key];
        }
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

  // 获取团队列表
  const fetchTeams = async () => {
    try {
      const response = await getContestTeams(parseInt(contestId), { limit: 20, offset: 0 });
      if (response.code === 200) {
        setTeams(response.data.teams || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchTeamsFailed') });
    }
  };

  // 获取题目列表
  const fetchChallenges = async () => {
    try {
      const response = await getContestChallenges(parseInt(contestId));
      if (response.code === 200) {
        setChallenges(response.data.challenges || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchChallengesFailed') });
    }
  };

  useEffect(() => {
    return addMessageHandler((data) => {
      if (data.type === 'start_victim' || data.type === 'stop_victim') {
        switch (data.level) {
          case 'error':
            toast.danger({ title: data.title, description: data.msg });
            break;
          case 'warning':
            toast.warning({ title: data.title, description: data.msg });
            break;
          case 'success':
            toast.success({ title: data.title, description: data.msg });
            break;
          case 'info':
            toast.info({ title: data.title, description: data.msg });
            break;
          default:
            toast.default({ title: data.title, description: data.msg });
            break;
        }
        fetchContainers();
      }
    });
  }, [addMessageHandler]);

  useEffect(() => {
    fetchContainers();
    fetchTeams();
    fetchChallenges();
  }, [contestId]);

  useEffect(() => {
    fetchContainers();
  }, [currentPage, filters.user_id, filters.team_id, filters.challenge_id]);

  // 点击外部关闭搜索结果
  useEffect(() => {
    const handleClickOutside = (event) => {
      const isOutsideUsers = usersSearchRef.current && !usersSearchRef.current.contains(event.target);
      const isOutsideTeams = teamsSearchRef.current && !teamsSearchRef.current.contains(event.target);
      const isOutsideChallenges = challengesSearchRef.current && !challengesSearchRef.current.contains(event.target);

      if (isOutsideUsers && searchResults.users.length > 0) {
        setSearchResults((prev) => ({ ...prev, users: [] }));
      }
      if (isOutsideTeams && searchResults.teams.length > 0) {
        setSearchResults((prev) => ({ ...prev, teams: [] }));
      }
      if (isOutsideChallenges && searchResults.challenges.length > 0) {
        setSearchResults((prev) => ({ ...prev, challenges: [] }));
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [searchResults]);

  // 处理过滤器变更
  const handleFilterChange = (key, value) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
    }));
    setCurrentPage(1);
  };

  // 搜索函数
  const handleSearch = async (model, name, setResults, setLoading) => {
    if (!name || name.trim() === '') {
      setResults([]);
      return;
    }

    setLoading(true);
    try {
      const response = await searchModels({
        model,
        name: name.trim(),
        limit: 10,
        offset: 0,
      });

      if (response.code === 200) {
        let results = response.data.results || [];

        // 如果是搜索团队，需要过滤contest_id
        if (model === 'Team') {
          results = results.filter((item) => item.contest_id === parseInt(contestId));
        }

        setResults(results);
      } else {
        setResults([]);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.searchFailed') });
      setResults([]);
    }
  };

  // 防抖搜索函数
  const debounceTimerRef = useRef(null);
  const debouncedSearch = (model, name, setResults, setLoading) => {
    clearTimeout(debounceTimerRef.current);
    debounceTimerRef.current = setTimeout(() => {
      handleSearch(model, name, setResults, setLoading);
    }, 300); // 300ms 防抖延迟
  };

  // 重置过滤器
  const handleResetFilters = () => {
    setFilters({
      user_id: '',
      team_id: '',
      challenge_id: '',
      limit: 20,
      offset: 0,
    });
    setSearchResults({
      users: [],
      teams: [],
      challenges: [],
    });
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

  // 监听随机选择百分比变化，自动执行随机选择
  useEffect(() => {
    if (teams.length > 0) {
      const percentage = randomTeamPercentage / 100;
      const count = Math.max(1, Math.floor(teams.length * percentage));

      // 随机打乱队伍数组并选择前count个
      const shuffledTeams = [...teams].sort(() => Math.random() - 0.5);
      const randomTeamIds = shuffledTeams.slice(0, count).map((team) => team.id);

      setSelectedTeams(randomTeamIds);
    }
  }, [randomTeamPercentage, teams]);

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
    if (selectedChallenges.length === 0 || selectedTeams.length === 0) {
      toast.warning({ description: t('admin.contests.containers.toast.selectStartRequired') });
      return;
    }

    try {
      const response = await startContestVictims(parseInt(contestId), selectedChallenges, selectedTeams);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.containers.toast.taskDispatched') });
        setSelectedChallenges([]);
        setSelectedTeams([]);
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
          <div className="flex items-center justify-between mb-4">
            <div>
              <p className="text-neutral-400 font-mono">{t('admin.contests.containers.page.subtitle')}</p>
            </div>
            <Button
              variant="primary"
              size="sm"
              align="icon-left"
              icon={<IconRefresh size={16} />}
              onClick={fetchContainers}
            >
              {t('common.refresh')}
            </Button>
          </div>

          {/* 统计卡片 */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <motion.div
              className="border border-neutral-300/30 rounded-md bg-black/30 backdrop-blur-[2px] p-4"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-blue-400/20 rounded-md flex items-center justify-center">
                  <IconServer size={20} className="text-blue-400" />
                </div>
                <div>
                  <p className="text-sm font-mono text-neutral-400">{t('admin.contests.containers.stats.total')}</p>
                  <p className="text-2xl font-mono text-neutral-50">{stats.totalContainers}</p>
                </div>
              </div>
            </motion.div>

            <motion.div
              className="border border-neutral-300/30 rounded-md bg-black/30 backdrop-blur-[2px] p-4"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-green-400/20 rounded-md flex items-center justify-center">
                  <IconPlayerPlay size={20} className="text-green-400" />
                </div>
                <div>
                  <p className="text-sm font-mono text-neutral-400">{t('admin.contests.containers.stats.running')}</p>
                  <p className="text-2xl font-mono text-green-400">{stats.runningContainers}</p>
                </div>
              </div>
            </motion.div>

            <motion.div
              className="border border-neutral-300/30 rounded-md bg-black/30 backdrop-blur-[2px] p-4"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2 }}
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-red-400/20 rounded-md flex items-center justify-center">
                  <IconBan size={20} className="text-red-400" />
                </div>
                <div>
                  <p className="text-sm font-mono text-neutral-400">{t('admin.contests.containers.stats.stopped')}</p>
                  <p className="text-2xl font-mono text-red-400">{stats.stoppedContainers}</p>
                </div>
              </div>
            </motion.div>
          </div>
        </div>

        {/* 快速操作和过滤条件 */}
        <motion.div
          className="grid grid-cols-1 lg:grid-cols-2 gap-6"
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
        >
          {/* 快速操作区域 */}
          <div className="border border-neutral-300/30 rounded-md bg-black/30 backdrop-blur-[2px] p-4">
            <div className="flex items-center gap-2 mb-3">
              <IconPlayerPlay size={18} className="text-neutral-400" />
              <h3 className="text-base font-mono text-neutral-50">
                {t('admin.contests.containers.quickActions.title')}
              </h3>
            </div>

            <div className="grid grid-cols-2 gap-4">
              {/* 选择题目 */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
                    <IconTarget size={14} />
                    {t('admin.contests.containers.quickActions.selectChallenges')}
                  </label>
                  <div className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() =>
                        setSelectedChallenges(challenges.filter((c) => c.type === 'pods').map((c) => c.id))
                      }
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
                <div className="max-h-24 overflow-y-auto border border-neutral-300/30 rounded-md bg-black/10">
                  {challenges
                    .filter((challenge) => challenge.type === 'pods')
                    .map((challenge) => (
                      <div key={challenge.id} className="flex items-center p-1 hover:bg-black/30 transition-colors">
                        <input
                          type="checkbox"
                          id={`challenge-${challenge.id}`}
                          checked={selectedChallenges.includes(challenge.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedChallenges((prev) => [...prev, challenge.id]);
                            } else {
                              setSelectedChallenges((prev) => prev.filter((id) => id !== challenge.id));
                            }
                          }}
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
                    ))}
                </div>
              </div>

              {/* 选择队伍 */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
                    <IconUsers size={14} />
                    <span className="text-xs font-mono text-neutral-400">
                      {t('admin.contests.containers.quickActions.randomTeams')}
                    </span>
                    <span className="text-xs font-mono text-geek-400">{randomTeamPercentage}%</span>
                  </label>
                  <div className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setSelectedTeams(teams.map((t) => t.id))}
                      className="!text-xs !h-5 !px-1"
                    >
                      {t('admin.contests.containers.quickActions.selectAll')}
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setSelectedTeams([])}
                      className="!text-xs !h-5 !px-1"
                    >
                      {t('admin.contests.containers.quickActions.clear')}
                    </Button>
                  </div>
                </div>

                {/* 随机选择拖动条 */}
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

                <div className="max-h-24 overflow-y-auto border border-neutral-300/30 rounded-md bg-black/10">
                  {teams.map((team) => (
                    <div key={team.id} className="flex items-center p-1 hover:bg-black/30 transition-colors">
                      <input
                        type="checkbox"
                        id={`team-${team.id}`}
                        checked={selectedTeams.includes(team.id)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedTeams((prev) => [...prev, team.id]);
                          } else {
                            setSelectedTeams((prev) => prev.filter((id) => id !== team.id));
                          }
                        }}
                        className="w-3 h-3 rounded border-neutral-300/30 text-geek-400
                                focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                      />
                      <label
                        htmlFor={`team-${team.id}`}
                        className="ml-2 text-xs font-mono text-neutral-300 cursor-pointer flex-1 truncate"
                      >
                        {team.name}
                      </label>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* 开启容器按钮 */}
            <div className="mt-3 flex justify-end">
              <Button
                variant="primary"
                size="sm"
                align="icon-left"
                icon={<IconPlayerPlay size={14} />}
                onClick={() => setIsStartModalOpen(true)}
                disabled={selectedChallenges.length === 0 || selectedTeams.length === 0}
                className="!text-xs !h-7 !px-3"
              >
                {t('admin.contests.containers.quickActions.startButton', {
                  challenges: selectedChallenges.length,
                  teams: selectedTeams.length,
                })}
              </Button>
            </div>
          </div>

          {/* 过滤器 */}
          <div className="border border-neutral-300/30 rounded-md bg-black/30 backdrop-blur-[2px] p-4">
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
            <div className="grid grid-cols-3 gap-3">
              {/* 用户搜索 */}
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
                    className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md
                            text-xs text-neutral-50 placeholder-neutral-500
                            focus:outline-none focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]
                            transition-all duration-200"
                  />
                  {searchLoading.users && (
                    <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
                      <div className="w-3 h-3 border border-geek-400 border-t-transparent rounded-full animate-spin"></div>
                    </div>
                  )}
                </div>
                {/* 搜索结果下拉 */}
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

              {/* 团队搜索 */}
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
                    className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md
                            text-xs text-neutral-50 placeholder-neutral-500
                            focus:outline-none focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]
                            transition-all duration-200"
                  />
                  {searchLoading.teams && (
                    <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
                      <div className="w-3 h-3 border border-geek-400 border-t-transparent rounded-full animate-spin"></div>
                    </div>
                  )}
                </div>
                {/* 搜索结果下拉 */}
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

              {/* 题目搜索 */}
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
                    className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md
                            text-xs text-neutral-50 placeholder-neutral-500
                            focus:outline-none focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]
                            transition-all duration-200"
                  />
                  {searchLoading.challenges && (
                    <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
                      <div className="w-3 h-3 border border-geek-400 border-t-transparent rounded-full animate-spin"></div>
                    </div>
                  )}
                </div>
                {/* 搜索结果下拉 */}
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

            {/* 当前选中的过滤条件显示 */}
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
                    <span className="px-2 py-1 bg-blue-400/20 text-blue-400 text-xs font-mono rounded border border-blue-400/30">
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

        {/* 容器列表 */}
        <Card variant="default" padding="none" className="overflow-hidden">
          {/* 列表头部 */}
          <div className="p-4 bg-black/20 border-b border-neutral-300/30 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <IconTable size={20} className="text-neutral-400" />
              <h3 className="text-lg font-mono text-neutral-50">{t('admin.contests.containers.table.title')}</h3>
              <span className="text-sm font-mono text-neutral-400">
                {t('admin.contests.containers.table.total', { count: runningCount })}
              </span>
            </div>
            <div className="flex items-center gap-3">
              <span className="text-sm font-mono text-neutral-400">
                {t('admin.contests.containers.table.selectedCount', { count: selectedContainers.length })}
              </span>
              <Button
                variant="danger"
                size="sm"
                align="icon-left"
                icon={<IconBan size={16} />}
                onClick={() => setIsStopModalOpen(true)}
                disabled={selectedContainers.length === 0}
              >
                {t('admin.contests.containers.table.stopButton')}
              </Button>
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
                </tr>
              </thead>
              <tbody>
                {containers.length === 0 ? (
                  <tr>
                    <td colSpan="9">
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
                  {t('admin.contests.containers.modals.selectedTeams', { count: selectedTeams.length })}
                </h4>
                <div className="max-h-32 overflow-y-auto border border-neutral-300/30 rounded-md bg-black/10 p-2">
                  {selectedTeams.map((teamId) => {
                    const team = teams.find((t) => t.id === teamId);
                    return team ? (
                      <div key={teamId} className="text-sm font-mono text-geek-400 py-1">
                        • {team.name}
                      </div>
                    ) : null;
                  })}
                </div>
              </div>
            </div>

            <div className="bg-neutral-800/50 border border-neutral-600/30 rounded-md p-3">
              <p className="text-xs font-mono text-neutral-400">
                {t('admin.contests.containers.modals.summaryPrefix')}
                <span className="text-geek-400">{selectedChallenges.length}</span>
                {t('admin.contests.containers.modals.summaryMiddle')}
                <span className="text-geek-400">{selectedTeams.length}</span>
                {t('admin.contests.containers.modals.summaryEquals')}
                <span className="text-green-400"> {selectedChallenges.length * selectedTeams.length}</span>
                {t('admin.contests.containers.modals.summarySuffix')}
              </p>
            </div>
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
