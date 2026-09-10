import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { toast } from '../../utils/toast';
import { downloadBlobResponse } from '../../utils/fileDownload';
import StatusPanel from '../../components/features/CTFGame/StatusPanel';
import ChallengeBoard from '../../components/features/CTFGame/Challenges/ChallengeBoard';
import ContestCountdown from '../../components/features/CTFGame/Challenges/ContestCountdown';
import ContestEnded from '../../components/features/CTFGame/Challenges/ContestEnded';
import ChallengeModal from '../../components/features/CTFGame/Challenges/ChallengeModal';
import { getContestInfo } from '../../api/contest';
import {
  getChallengeList,
  getChallengeStatus,
  initChallenge,
  startRemoteTarget,
  extendContainerTime,
  stopContainer,
  submitFlag,
  resetChallenge,
  uploadWriteup,
  getWriteups,
  getChallengeCategories,
} from '../../api/challenge';
import { getTeamMembers, getTeamInfo } from '../../api/game/team';
import { downloadChallengeAttachment } from '../../api/challenge';
import { getContestNotices } from '../../api/contest';
import Loading from '../../components/common/Loading';
import { Button, EmptyState } from '../../components/common';
import { useTranslation } from 'react-i18next';

const normalizeInstanceStatus = (status) => {
  const normalizedStatus = typeof status === 'string' ? status.toLowerCase() : '';
  if (
    normalizedStatus === 'waiting' ||
    normalizedStatus === 'pending' ||
    normalizedStatus === 'terminating' ||
    normalizedStatus === 'running'
  ) {
    return normalizedStatus;
  }
  return '';
};

const mapChallengeStatusToViewModel = (challenge, statusData = null) => {
  const remote = statusData?.remote || challenge.remote || {};
  const instanceStatus = normalizeInstanceStatus(remote.status);
  const timeLeft = Number(remote.remaining) || 0;
  const remoteDuration = Number(remote.duration) || 0;
  const previousDuration = Number(challenge.instanceDuration) || 0;
  const instanceDuration =
    remoteDuration > 0
      ? remoteDuration
      : instanceStatus === 'running'
        ? Math.max(previousDuration, timeLeft)
        : previousDuration;

  return {
    ...challenge,
    id: challenge.id,
    type: challenge.type,
    category: challenge.category,
    title: challenge.title || challenge.name,
    score: challenge.score,
    description: challenge.description,
    attachments: challenge.attachments || [],
    attachment: statusData?.file ?? challenge.attachment ?? '',
    hasInstance: challenge.hasInstance ?? challenge.type === 'pods',
    hasAttachments: challenge.hasAttachments ?? challenge.type === 'dynamic',
    instanceStatus,
    instanceRunning: instanceStatus === 'running',
    instancePending: instanceStatus === 'pending',
    instanceWaiting: instanceStatus === 'waiting',
    instanceTerminating: instanceStatus === 'terminating',
    instanceIP: remote.target || challenge.instanceIP || [''],
    instanceDuration,
    instanceTimeLeft: timeLeft,
    solves: challenge.solves ?? challenge.solvers ?? 0,
    isInitialized: statusData?.init ?? challenge.isInitialized ?? challenge.init,
    isSolved: statusData?.solved ?? challenge.isSolved ?? challenge.solved ?? false,
    solved: statusData?.solved ?? challenge.solved ?? challenge.isSolved ?? false,
    attempts: statusData?.attempts ?? challenge.attempts,
    maxAttempts: challenge.maxAttempts ?? challenge.attempt,
    hints: challenge.hints,
    tags: challenge.tags,
    hidden: challenge.hidden,
    options: challenge.options || [],
  };
};

// 计算比赛状态
const getContestStatus = (contest, teamInfo) => {
  const now = new Date().getTime();
  const start = new Date(contest.start).getTime();
  const end = start + contest.duration * 1000; // duration 是秒, 转换为毫秒

  return {
    status: now < start ? 'upcoming' : now > end ? 'ended' : 'running',

    startTime: new Date(start).toISOString(),
    endTime: new Date(end).toISOString(),
    joined: true,
    duration: contest.duration / 3600,
    prefix: contest.prefix,
    team: {
      score: teamInfo.data.score,
      rank: teamInfo.data.rank || 0,
      solved: teamInfo.data.solved.reduce((acc, curr) => acc + curr.solved, 0),
    },
    teams: contest.teams,
    totalChallenges: teamInfo.data.solved.reduce((acc, curr) => acc + curr.all, 0),
  };
};

// 修改 transformChallengeData 函数
const transformChallengeData = (challenge) => {
  return mapChallengeStatusToViewModel(challenge);
};

const normalizeCategories = (categoryList) => (Array.isArray(categoryList) ? categoryList.filter(Boolean) : []);
const isInstanceTransitioning = (status) => ['waiting', 'pending', 'terminating'].includes(status);

function ContestChallenges({ contestId }) {
  const navigate = useNavigate();
  const [contestStatus, setContestStatus] = useState({});
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [unsolvedOnly, setUnsolvedOnly] = useState(false);
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const [loading, setLoading] = useState(true);
  const [contestError, setContestError] = useState(null);
  const [listError, setListError] = useState(null);
  const [listLoading, setListLoading] = useState(true);
  const [listRevision, setListRevision] = useState(0);
  const [teamInfo, setTeamInfo] = useState(null);
  const [notifications, setNotifications] = useState([]);
  const [writeups, setWriteups] = useState([]);
  const [showChallengesAfterEnd, setShowChallengesAfterEnd] = useState(false);
  const { t } = useTranslation();

  // 分页配置
  const pageSize = 20;

  const selectedChallengeRef = useRef(null);
  const pollingIntervalRef = useRef(null);
  const pollingTimeoutRef = useRef(null);
  const pollingGenerationRef = useRef(0);
  const selectionRef = useRef(0);
  const statusRequestRef = useRef(0);
  const scopeRef = useRef(null);
  const contestRequestRef = useRef(0);
  const categoriesRequestRef = useRef(0);
  const writeupRequestRef = useRef(0);

  const stopPolling = () => {
    pollingGenerationRef.current += 1;
    if (pollingIntervalRef.current) {
      clearTimeout(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
    if (pollingTimeoutRef.current) {
      clearTimeout(pollingTimeoutRef.current);
      pollingTimeoutRef.current = null;
    }
  };

  const startPolling = (challengeId, targetStatus, selection) => {
    if (selectionRef.current !== selection || selectedChallengeRef.current?.id !== challengeId) return;
    stopPolling();
    const generation = pollingGenerationRef.current;
    const isCurrent = () => generation === pollingGenerationRef.current && selection === selectionRef.current;
    // Schedule only after completion, so slow responses cannot overlap the next poll.
    const poll = async () => {
      if (!isCurrent()) return;
      const request = ++statusRequestRef.current;
      try {
        const statusRes = await getChallengeStatus(contestId, challengeId);
        if (!isCurrent()) return;
        if (statusRes.code === 200 && request === statusRequestRef.current) {
          const updatedChallenge = mapChallengeStatusToViewModel(selectedChallengeRef.current, statusRes.data);
          selectedChallengeRef.current = updatedChallenge;
          setSelectedChallenge(updatedChallenge);
          setChallenges((prev) => prev.map((c) => (c.id === updatedChallenge.id ? updatedChallenge : c)));

          if (
            (targetStatus === 'running' && updatedChallenge.instanceStatus === 'running') ||
            (targetStatus === 'stopped' &&
              updatedChallenge.instanceStatus !== 'running' &&
              !isInstanceTransitioning(updatedChallenge.instanceStatus))
          ) {
            stopPolling();
          }
        }
      } catch {
        // Silently ignore polling errors
      }
      if (isCurrent()) pollingIntervalRef.current = setTimeout(poll, 5000);
    };
    pollingIntervalRef.current = setTimeout(poll, 5000);
    pollingTimeoutRef.current = setTimeout(stopPolling, 3 * 60 * 1000);
  };

  const closeChallenge = () => {
    stopPolling();
    selectionRef.current += 1;
    selectedChallengeRef.current = null;
    setSelectedChallenge(null);
  };

  // 获取比赛信息和题目列表
  useEffect(() => {
    scopeRef.current = { contestId };
    setLoading(true);
    setContestStatus({});
    setTeamInfo(null);
    setNotifications([]);
    setCategories([]);
    setWriteups([]);
    setShowChallengesAfterEnd(false);
    closeChallenge();
    fetchContestInfo();
    fetchCategories();
    return () => {
      scopeRef.current = null;
      selectionRef.current += 1;
      selectedChallengeRef.current = null;
      stopPolling();
    };
  }, [contestId]);

  // 当当前页变化时, 重新获取题目数据
  useEffect(() => {
    let active = true;
    setListLoading(true);
    setListError(null);
    fetchChallengesWithFilters(currentPage, selectedCategory, unsolvedOnly, () => active);
    return () => {
      active = false;
    };
  }, [currentPage, contestId, selectedCategory, unsolvedOnly, listRevision]);

  // 获取已上传的题解
  useEffect(() => {
    if (contestStatus?.status === 'ended') {
      fetchWriteups();
    }
  }, [contestId, contestStatus?.status]);

  useEffect(() => {
    if (contestStatus.status === 'ended' && !showChallengesAfterEnd) closeChallenge();
  }, [contestStatus.status, showChallengesAfterEnd]);

  const fetchContestInfo = async () => {
    const scope = scopeRef.current;
    const request = ++contestRequestRef.current;
    const isCurrent = () => scopeRef.current === scope && request === contestRequestRef.current;
    setContestError(null);
    if (!contestStatus.status) setLoading(true);

    // Notices are optional: neither their latency nor failure should block the contest.
    getContestNotices(contestId)
      .then((response) => {
        if (!isCurrent()) return;
        if (response.code !== 200) throw new Error(response.msg || t('errors.requestFailed'));
        setNotifications(
          (response.data?.notices || []).map((notice) => ({
            type: notice.type || 'info',
            title: notice.title,
            message: notice.content,
          }))
        );
      })
      .catch((error) => {
        if (isCurrent()) toast.danger({ description: error.message || t('errors.requestFailed') });
      });

    try {
      const [contestRes, teamMembersRes, teamInfoRes] = await Promise.all([
        getContestInfo(contestId),
        getTeamMembers(contestId),
        getTeamInfo(contestId),
      ]);
      if (!isCurrent()) return;

      if (contestRes.code === 200 && teamMembersRes.code === 200 && teamInfoRes.code === 200) {
        // 设置比赛状态
        const status = getContestStatus(contestRes.data, teamInfoRes);
        const team = {
          members: teamMembersRes.data.map((member) => ({
            picture: member.picture,
            name: member.name,
          })),
          name: teamInfoRes.data.name,
        };
        setContestStatus(status);
        setTeamInfo(team);
      } else {
        const failedResponse = [contestRes, teamMembersRes, teamInfoRes].find((response) => response.code !== 200);
        throw new Error(failedResponse.msg || t('game.challenges.toast.fetchFailed'));
      }
    } catch (error) {
      if (!isCurrent()) return;
      setContestError(error.message || t('game.challenges.toast.fetchFailed'));
      if (contestStatus.status) {
        toast.danger({
          title: t('game.challenges.toast.fetchFailed'),
          description: error.message,
        });
      }
    } finally {
      if (isCurrent()) setLoading(false);
    }
  };

  // 获取所有分类
  const fetchCategories = async () => {
    const scope = scopeRef.current;
    const request = ++categoriesRequestRef.current;
    const isCurrent = () => scopeRef.current === scope && request === categoriesRequestRef.current;
    try {
      // 获取所有题目来提取分类
      const response = await getChallengeCategories(contestId);
      if (!isCurrent()) return;
      if (response.code !== 200) throw new Error(response.msg || t('game.challenges.toast.fetchCategoriesFailed'));
      if (response.code === 200) {
        setCategories(normalizeCategories(response.data));
      }
    } catch (error) {
      if (!isCurrent()) return;
      setCategories([]);
      toast.danger({
        description: error.message || t('game.challenges.toast.fetchCategoriesFailed'),
      });
    }
  };

  // 获取带过滤器的题目数据
  const fetchChallengesWithFilters = async (page, category, unsolved, isCurrent) => {
    try {
      const params = {
        limit: pageSize,
        offset: (page - 1) * pageSize,
        unsolved,
      };

      // 添加过滤参数
      if (category) {
        params.category = category;
      }

      const response = await getChallengeList(contestId, params);
      if (!isCurrent()) return;
      if (response.code !== 200) throw new Error(response.msg || t('game.challenges.toast.fetchListFailed'));
      if (response.code === 200) {
        const challengeData = response.data === null ? [] : response.data;

        // 转换并设置题目列表
        const transformedChallenges = (challengeData.challenges || []).map((challenge) =>
          transformChallengeData(challenge)
        );
        setChallenges(transformedChallenges);
        setTotalCount(challengeData.count || 0);
        const lastPage = Math.max(1, Math.ceil((challengeData.count || 0) / pageSize));
        if (page > lastPage) setCurrentPage(lastPage);
      }
    } catch (error) {
      if (!isCurrent()) return;
      setChallenges([]);
      setTotalCount(0);
      setListError(error.message || t('game.challenges.toast.fetchListFailed'));
    } finally {
      if (isCurrent()) setListLoading(false);
    }
  };

  // 获取已上传的题解数据
  const fetchWriteups = async () => {
    const scope = scopeRef.current;
    const request = ++writeupRequestRef.current;
    const isCurrent = () => scopeRef.current === scope && request === writeupRequestRef.current;
    try {
      const response = await getWriteups(contestId);
      if (!isCurrent()) return;

      if (response.code === 200 && response.data.writeups) {
        setWriteups(response.data.writeups);
      }
    } catch (error) {
      if (!isCurrent()) return;
      toast.danger({
        description: error.message || t('game.challenges.toast.fetchWriteupsFailed'),
      });
    }
  };

  // 刷新当前选中题目的状态
  const refreshChallengeStatus = async (challengeId, selection) => {
    if (selectionRef.current !== selection || selectedChallengeRef.current?.id !== challengeId) return;
    const request = ++statusRequestRef.current;
    const isCurrent = () => selectionRef.current === selection && request === statusRequestRef.current;

    try {
      const statusRes = await getChallengeStatus(contestId, challengeId);
      if (!isCurrent()) return;

      if (statusRes.code === 200) {
        const updatedChallenge = mapChallengeStatusToViewModel(selectedChallengeRef.current, statusRes.data);
        selectedChallengeRef.current = updatedChallenge;
        setSelectedChallenge(updatedChallenge);

        // 更新题目列表中的对应题目
        setChallenges((prev) => prev.map((c) => (c.id === updatedChallenge.id ? updatedChallenge : c)));
      }
    } catch (error) {
      if (!isCurrent()) return;
      toast.danger({
        description: error.message || t('game.challenges.toast.refreshStatusFailed'),
      });
    }
  };

  // 处理题目初始化
  const handleInitialize = async (challengeId) => {
    const selection = selectionRef.current;
    try {
      const res = await initChallenge(contestId, challengeId);
      if (selectionRef.current !== selection) return false;
      if (res.code !== 200) throw new Error(res.msg || t('game.challenges.toast.initFailed'));
      if (res.code === 200) {
        toast.success({
          title: res.msg || t('game.challenges.toast.initSuccess'),
        });
        await refreshChallengeStatus(challengeId, selection);
        return true;
      }
    } catch (error) {
      if (selectionRef.current !== selection) return false;
      toast.danger({
        title: t('game.challenges.toast.initFailed'),
        description: error.message,
      });
    }
    return false;
  };

  const handleReset = async (challengeId) => {
    const selection = selectionRef.current;
    try {
      const res = await resetChallenge(contestId, challengeId);
      if (selectionRef.current !== selection) return false;
      if (res.code !== 200) throw new Error(res.msg || t('game.challenges.toast.resetFailed'));
      if (res.code === 200) {
        toast.success({ title: t('game.challenges.toast.resetSuccess') });
        await refreshChallengeStatus(challengeId, selection);
        return true;
      }
    } catch (error) {
      if (selectionRef.current !== selection) return false;
      toast.danger({
        title: t('game.challenges.toast.resetFailed'),
        description: error.message,
      });
    }
    return false;
  };

  // 处理启动靶机
  const handleLaunchInstance = async (challengeId) => {
    const selection = selectionRef.current;
    try {
      const res = await startRemoteTarget(contestId, challengeId);
      if (selectionRef.current !== selection) return false;
      if (res.code !== 200) throw new Error(res.msg || t('game.challenges.toast.launchFailed'));
      if (res.code === 200) {
        statusRequestRef.current += 1;
        selectedChallengeRef.current = {
          ...selectedChallengeRef.current,
          instanceStatus: 'waiting',
          instanceRunning: false,
          instancePending: false,
          instanceWaiting: true,
        };
        setSelectedChallenge(selectedChallengeRef.current);
        setChallenges((prev) =>
          prev.map((challenge) =>
            challenge.id === challengeId
              ? {
                  ...challenge,
                  instanceStatus: 'waiting',
                  instanceRunning: false,
                  instancePending: false,
                  instanceWaiting: true,
                }
              : challenge
          )
        );
        toast.success({ title: t('game.challenges.toast.launchSuccess') });
        startPolling(challengeId, 'running', selection);
        return true;
      }
    } catch (error) {
      if (selectionRef.current !== selection) return false;
      toast.danger({
        title: t('game.challenges.toast.launchFailed'),
        description: error.message,
      });
    }
    return false;
  };

  // 处理延长靶机时间
  const handleExtendInstance = async (challengeId) => {
    const selection = selectionRef.current;
    try {
      const res = await extendContainerTime(contestId, challengeId);
      if (selectionRef.current !== selection) return false;
      if (res.code !== 200) throw new Error(res.msg || t('game.challenges.toast.extendFailed'));
      if (res.code === 200) {
        toast.success({
          title: res.msg || t('game.challenges.toast.extendSuccess'),
        });
        await refreshChallengeStatus(challengeId, selection);
        return true;
      }
    } catch (error) {
      if (selectionRef.current !== selection) return false;
      toast.danger({
        title: t('game.challenges.toast.extendFailed'),
        description: error.message,
      });
    }
    return false;
  };

  // 处理销毁靶机
  const handleDestroyInstance = async (challengeId) => {
    const selection = selectionRef.current;
    try {
      const res = await stopContainer(contestId, challengeId);
      if (selectionRef.current !== selection) return false;
      if (res.code !== 200) throw new Error(res.msg || t('game.challenges.toast.destroyFailed'));
      if (res.code === 200) {
        toast.success({
          title: res.msg || t('game.challenges.toast.destroySuccess'),
        });
        await refreshChallengeStatus(challengeId, selection);
        startPolling(challengeId, 'stopped', selection);
        return true;
      }
    } catch (error) {
      if (selectionRef.current !== selection) return false;
      toast.danger({
        title: t('game.challenges.toast.destroyFailed'),
        description: error.message,
      });
    }
    return false;
  };

  // 处理提交flag
  const handleSubmitFlag = async (challengeId, value) => {
    const selection = selectionRef.current;
    try {
      const data = { flag: value };
      const res = await submitFlag(contestId, challengeId, data);
      if (selectionRef.current !== selection) return { success: false };
      // Refresh attempts even for an incorrect flag, without replacing the modal or its input.
      void refreshChallengeStatus(challengeId, selection);
      if (res.code === 200) {
        void fetchContestInfo();
        setListRevision((revision) => revision + 1);
        toast.success({
          title: res.msg || t('game.challenges.toast.submitSuccess'),
        });
        return {
          success: true,
          message: t('game.challenges.toast.submitSuccessMessage'),
        };
      }
      return { success: false, message: res.msg || t('errors.requestFailed') };
    } catch (error) {
      return {
        success: false,
        message: error.message || t('errors.requestFailed'),
      };
    }
  };

  // 处理分类切换
  const handleCategoryChange = (category) => {
    const nextCategory = selectedCategory === category ? '' : category;
    setSelectedCategory(nextCategory);
    setCurrentPage(1); // 重置到第一页
  };

  const handleSolvedFilterChange = () => {
    const nextUnsolvedOnly = !unsolvedOnly;
    setUnsolvedOnly(nextUnsolvedOnly);
    setCurrentPage(1);
  };

  // 处理题目点击
  const handleChallengeClick = async (challenge) => {
    stopPolling();
    const selection = ++selectionRef.current;
    selectedChallengeRef.current = null;
    setSelectedChallenge(null);
    try {
      // 获取题目状态, 包含附件、靶机信息和初始化状态
      const statusRes = await getChallengeStatus(contestId, challenge.id);
      if (selectionRef.current !== selection) return;
      if (statusRes.code !== 200) throw new Error(statusRes.msg || t('game.challenges.toast.fetchStatusFailed'));
      if (statusRes.code === 200) {
        const updatedChallenge = mapChallengeStatusToViewModel(challenge, statusRes.data);

        selectedChallengeRef.current = updatedChallenge;
        setSelectedChallenge(updatedChallenge);

        // 页面刷新后 Pod 仍在排队或启动中 → 自动开始轮询
        if (isInstanceTransitioning(updatedChallenge.instanceStatus)) {
          startPolling(
            challenge.id,
            updatedChallenge.instanceStatus === 'terminating' ? 'stopped' : 'running',
            selection
          );
        }
      }
    } catch (error) {
      if (selectionRef.current !== selection) return;
      selectedChallengeRef.current = challenge;
      setSelectedChallenge(challenge); // 即使失败也显示题目
      toast.danger({
        description: error.message || t('game.challenges.toast.fetchStatusFailed'),
      });
    }
  };

  // 处理附件下载
  const handleDownloadAttachment = async (attachment) => {
    try {
      const response = await downloadChallengeAttachment(contestId, selectedChallenge.id);
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, attachment, 'application/octet-stream');
      }

      toast.success({ title: t('game.challenges.toast.downloadSuccess') });
    } catch (error) {
      toast.danger({
        description: error.message || t('game.challenges.toast.downloadFailed'),
      });
    }
  };

  // 处理上传题解
  const handleUploadWriteup = async (file) => {
    const scope = scopeRef.current;
    try {
      setLoading(true);
      const res = await uploadWriteup(contestId, file);
      if (scopeRef.current !== scope) return;

      if (res.code === 200) {
        toast.success({
          title: t('game.challenges.toast.uploadSuccess'),
          description: t('game.challenges.toast.uploadThanks'),
        });
        // 上传成功后刷新题解列表
        fetchWriteups();
      }
    } catch (error) {
      if (scopeRef.current !== scope) return;
      toast.danger({
        description: error.message || t('game.challenges.toast.uploadFailed'),
      });
    } finally {
      if (scopeRef.current === scope) setLoading(false);
    }
  };

  // 处理查看题目
  const handleViewChallenges = () => {
    setShowChallengesAfterEnd(true);
  };

  // 处理返回比赛结束页面
  const handleBackToContestEnd = () => {
    setShowChallengesAfterEnd(false);
    closeChallenge();
  };

  if (loading) {
    return <Loading />;
  }

  if (!contestStatus.status) {
    return (
      <div role="alert">
        <EmptyState
          title={t('game.challenges.toast.fetchFailed')}
          description={contestError}
          action={
            <Button
              onClick={() => {
                fetchContestInfo();
                fetchCategories();
                setListRevision((revision) => revision + 1);
              }}
            >
              {t('common.refresh')}
            </Button>
          }
        />
      </div>
    );
  }

  return (
    <div className="contest-container mx-auto space-y-6">
      <h1 className="sr-only">{contestStatus?.name || t('game.challenges.title')}</h1>
      <StatusPanel
        contestStatus={contestStatus}
        onStatusExpired={(newStatus) => {
          // 这里更新比赛状态
          setContestStatus((prev) => ({
            ...prev,
            status: newStatus,
          }));
        }}
        notifications={notifications}
      />
      {contestStatus?.status === 'upcoming' ? (
        <ContestCountdown
          startTime={contestStatus.startTime}
          joined={contestStatus.joined}
          onJoin={() => {}} // 需要实现加入比赛的逻辑
        />
      ) : contestStatus?.status === 'ended' ? (
        showChallengesAfterEnd ? (
          <div>
            <div className="mb-4">
              <Button variant="secondary" size="sm" onClick={handleBackToContestEnd}>
                {t('game.challenges.backToSummary')}
              </Button>
            </div>
            <ChallengeBoard
              categories={categories}
              selectedCategory={selectedCategory}
              onCategoryChange={handleCategoryChange}
              unsolvedOnly={unsolvedOnly}
              onSolvedFilterChange={handleSolvedFilterChange}
              challenges={challenges}
              onChallengeClick={handleChallengeClick}
              teamInfo={teamInfo}
              totalCount={totalCount}
              currentPage={currentPage}
              pageSize={pageSize}
              onPageChange={setCurrentPage}
              isLoading={listLoading}
              error={listError}
              onRetry={() => setListRevision((revision) => revision + 1)}
            />
            <ChallengeModal
              key={selectedChallenge?.id || 'closed'}
              challenge={selectedChallenge}
              contest={contestStatus}
              isOpen={!!selectedChallenge}
              onClose={closeChallenge}
              onInitialize={handleInitialize}
              onReset={handleReset}
              onLaunchInstance={handleLaunchInstance}
              onExtendInstance={handleExtendInstance}
              onDestroyInstance={handleDestroyInstance}
              onSubmitFlag={handleSubmitFlag}
              onDownloadAttachment={handleDownloadAttachment}
            />
          </div>
        ) : (
          <ContestEnded
            contestInfo={{
              duration: `${contestStatus.duration}h`,
              totalTeams: contestStatus.teams,
              totalChallenges: contestStatus.totalChallenges,
              teamRank: contestStatus.team.rank,
              teamScore: contestStatus.team.score,
              teamSolved: contestStatus.team.solved,
            }}
            onViewScoreboard={() => {
              navigate(`/contests/${contestId}/scoreboard`);
            }}
            onUploadWriteup={handleUploadWriteup}
            onViewChallenges={handleViewChallenges}
            writeups={writeups}
          />
        )
      ) : (
        <div>
          <ChallengeBoard
            categories={categories}
            selectedCategory={selectedCategory}
            onCategoryChange={handleCategoryChange}
            unsolvedOnly={unsolvedOnly}
            onSolvedFilterChange={handleSolvedFilterChange}
            challenges={challenges}
            onChallengeClick={handleChallengeClick}
            teamInfo={teamInfo}
            totalCount={totalCount}
            currentPage={currentPage}
            pageSize={pageSize}
            onPageChange={setCurrentPage}
            isLoading={listLoading}
            error={listError}
            onRetry={() => setListRevision((revision) => revision + 1)}
          />
          <ChallengeModal
            key={selectedChallenge?.id || 'closed'}
            challenge={selectedChallenge}
            contest={contestStatus}
            isOpen={!!selectedChallenge}
            onClose={closeChallenge}
            onInitialize={handleInitialize}
            onReset={handleReset}
            onLaunchInstance={handleLaunchInstance}
            onExtendInstance={handleExtendInstance}
            onDestroyInstance={handleDestroyInstance}
            onSubmitFlag={handleSubmitFlag}
            onDownloadAttachment={handleDownloadAttachment}
          />
        </div>
      )}
    </div>
  );
}

export default function GameChallengesPage() {
  const { contestId } = useParams();
  return <ContestChallenges key={contestId} contestId={contestId} />;
}
