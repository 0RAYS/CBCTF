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
  increaseContainerTime,
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
import EmptyState from '../../components/common/EmptyState';
import { Button } from '../../components/common';
import { useTranslation } from 'react-i18next';

// 计算比赛状态
const getContestStatus = (contest, teamInfo) => {
  const now = new Date().getTime();
  const start = new Date(contest.start).getTime();
  const end = start + contest.duration * 1000; // duration 是秒，转换为毫秒

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
  return {
    id: challenge.id,
    type: challenge.type,
    category: challenge.category,
    title: challenge.name,
    score: challenge.score,
    description: challenge.description,
    attachments: [], // 默认为空，点击时再获取
    hasInstance: challenge.type === 'pods',
    hasAttachments: challenge.type === 'dynamic',
    instanceRunning: challenge.remote.status === 'running',
    instancePending: challenge.remote.status === 'pending',
    instanceIP: challenge.remote.target || [''],
    instanceDuration: challenge.remote.remaining || 0,
    instanceTimeLeft: challenge.remote.remaining || 0,
    solves: challenge.solvers || 0,
    isInitialized: challenge.init,
    solved: challenge.solved, // 是否已解决
    attempts: challenge.attempts, // 尝试次数
    maxAttempts: challenge.attempt, // 最大尝试次数
    hints: challenge.hints,
    tags: challenge.tags,
    hidden: challenge.hidden,
    options: challenge.options || [], // 添加选项字段，用于question类型
  };
};

function GameChallengesPage() {
  const { contestId } = useParams();
  const navigate = useNavigate();
  const [contestStatus, setContestStatus] = useState({});
  const [categories, setCategories] = useState(['ALL']);
  const [selectedCategory, setSelectedCategory] = useState('ALL');
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const [loading, setLoading] = useState(true);
  const [teamInfo, setTeamInfo] = useState(null);
  const [notifications, setNotifications] = useState([]);
  const [writeups, setWriteups] = useState([]);
  const [showChallengesAfterEnd, setShowChallengesAfterEnd] = useState(false);
  const { t } = useTranslation();

  // 分页配置
  const pageSize = 10;

  const selectedChallengeRef = useRef(null);
  const pollingIntervalRef = useRef(null);
  const pollingTimeoutRef = useRef(null);

  const stopPolling = () => {
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
    if (pollingTimeoutRef.current) {
      clearTimeout(pollingTimeoutRef.current);
      pollingTimeoutRef.current = null;
    }
  };

  const startPolling = () => {
    stopPolling();
    pollingIntervalRef.current = setInterval(async () => {
      if (!selectedChallengeRef.current) {
        stopPolling();
        return;
      }
      try {
        const statusRes = await getChallengeStatus(contestId, selectedChallengeRef.current.id);
        if (statusRes.code === 200 && statusRes.data.remote.status === 'running') {
          stopPolling();
          refreshChallengeStatus();
        }
      } catch {
        // Silently ignore polling errors
      }
    }, 5000);
    pollingTimeoutRef.current = setTimeout(stopPolling, 3 * 60 * 1000);
  };

  // Keep selectedChallengeRef in sync; stop polling when modal closes
  useEffect(() => {
    selectedChallengeRef.current = selectedChallenge;
    if (!selectedChallenge) {
      stopPolling();
    }
  }, [selectedChallenge]);

  // 获取比赛信息和题目列表
  useEffect(() => {
    fetchContestAndChallenges();
  }, [contestId]);

  // 当当前页变化时，重新获取题目数据
  useEffect(() => {
    if (contestId && currentPage > 0) {
      fetchChallengesWithFilters(currentPage, selectedCategory);
    }
  }, [currentPage, contestId]);

  // 获取已上传的题解
  useEffect(() => {
    if (contestStatus?.status === 'ended') {
      fetchWriteups();
    }
  }, [contestId, contestStatus?.status]);

  const fetchContestAndChallenges = async () => {
    try {
      const [contestRes, teamMembersRes, teamInfoRes, noticesRes] = await Promise.all([
        getContestInfo(contestId),
        getTeamMembers(contestId),
        getTeamInfo(contestId),
        getContestNotices(contestId),
      ]);

      if (contestRes.code === 200 && teamMembersRes.code === 200) {
        // 设置比赛状态
        const status = getContestStatus(contestRes.data, teamInfoRes);
        setContestStatus(status);
        setNotifications(
          noticesRes.data.notices.map((notice) => ({
            type: notice.type || 'info',
            title: notice.title,
            message: notice.content,
          }))
        );

        // 设置队伍信息
        setTeamInfo({
          members: teamMembersRes.data.map((member) => ({ picture: member.picture, name: member.name })),
          name: teamInfoRes.data.name,
        });

        // 获取分类和题目数据
        await fetchCategories();
        await fetchChallengesWithFilters(1, selectedCategory);
      }
    } catch (error) {
      toast.danger({ title: t('game.challenges.toast.fetchFailed'), description: error.message });
    } finally {
      setLoading(false);
    }
  };

  // 获取所有分类
  const fetchCategories = async () => {
    try {
      // 获取所有题目来提取分类
      const response = await getChallengeCategories(contestId);
      if (response.code === 200) {
        setCategories(response.data === null ? [] : response.data);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.fetchCategoriesFailed') });
    }
  };

  // 获取带过滤器的题目数据
  const fetchChallengesWithFilters = async (page, category) => {
    try {
      const params = {
        limit: pageSize,
        offset: (page - 1) * pageSize,
      };

      // 添加过滤参数
      if (category !== 'ALL') {
        params.category = category;
      }

      const response = await getChallengeList(contestId, params);
      if (response.code === 200) {
        const challengeData = response.data === null ? [] : response.data;

        // 转换并设置题目列表
        const transformedChallenges = (challengeData.challenges || []).map((challenge) =>
          transformChallengeData(challenge)
        );
        setChallenges(transformedChallenges);
        setTotalCount(challengeData.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.fetchListFailed') });
    }
  };

  // 获取已上传的题解数据
  const fetchWriteups = async () => {
    try {
      const response = await getWriteups(contestId);

      if (response.code === 200 && response.data.writeups) {
        setWriteups(response.data.writeups);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.fetchWriteupsFailed') });
    }
  };

  // 刷新当前选中题目的状态
  const refreshChallengeStatus = async () => {
    if (!selectedChallengeRef.current) return;

    try {
      const statusRes = await getChallengeStatus(contestId, selectedChallengeRef.current.id);

      if (statusRes.code === 200) {
        const updatedChallenge = {
          ...selectedChallengeRef.current,
          attachment: statusRes.data.file || '',
          isInitialized: statusRes.data.init,
          isSolved: statusRes.data.solved,
          instanceRunning: statusRes.data.remote.status === 'running',
          instancePending: statusRes.data.remote.status === 'pending',
          instanceIP: statusRes.data.remote.target || [''],
          instanceDuration: statusRes.data.remote.remaining || 0,
          instanceTimeLeft: statusRes.data.remote.remaining || 0,
        };

        setSelectedChallenge(updatedChallenge);

        // 更新题目列表中的对应题目
        setChallenges((prev) => prev.map((c) => (c.id === updatedChallenge.id ? updatedChallenge : c)));
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.refreshStatusFailed') });
    }
  };

  // 处理题目初始化
  const handleInitialize = async (challengeId) => {
    try {
      const res = await initChallenge(contestId, challengeId);
      if (res.code === 200) {
        toast.success({ title: res.msg || t('game.challenges.toast.initSuccess') });
        await refreshChallengeStatus();
        return true;
      }
    } catch (error) {
      toast.danger({ title: t('game.challenges.toast.initFailed'), description: error.message });
    }
    return false;
  };

  const handleReset = async (challengeId) => {
    try {
      const res = await resetChallenge(contestId, challengeId);
      if (res.code === 200) {
        toast.success({ title: t('game.challenges.toast.resetSuccess') });
        return true;
      }
    } catch (error) {
      toast.danger({ title: t('game.challenges.toast.resetFailed'), description: error.message });
    }
    return false;
  };

  // 处理启动靶机
  const handleLaunchInstance = async (challengeId) => {
    try {
      const res = await startRemoteTarget(contestId, challengeId);
      if (res.code === 200) {
        toast.success({ title: t('game.challenges.toast.launchSuccess') });
        startPolling();
        return true;
      }
    } catch (error) {
      toast.danger({ title: t('game.challenges.toast.launchFailed'), description: error.message });
    }
    return false;
  };

  // 处理延长靶机时间
  const handleExtendInstance = async (challengeId) => {
    try {
      const res = await increaseContainerTime(contestId, challengeId);
      if (res.code === 200) {
        toast.success({ title: res.msg || t('game.challenges.toast.extendSuccess') });
        await refreshChallengeStatus();
        return true;
      }
    } catch (error) {
      toast.danger({ title: t('game.challenges.toast.extendFailed'), description: error.message });
    }
    return false;
  };

  // 处理销毁靶机
  const handleDestroyInstance = async (challengeId) => {
    try {
      const res = await stopContainer(contestId, challengeId);
      if (res.code === 200) {
        toast.success({ title: res.msg || t('game.challenges.toast.destroySuccess') });
        await refreshChallengeStatus();
        return true;
      }
    } catch (error) {
      toast.danger({ title: t('game.challenges.toast.destroyFailed'), description: error.message });
    }
    return false;
  };

  // 处理提交flag
  const handleSubmitFlag = async (challengeId, value) => {
    try {
      const data = { flag: value };
      const res = await submitFlag(contestId, challengeId, data);
      await refreshChallengeStatus();
      await fetchContestAndChallenges();
      if (res.code === 200) {
        toast.success({ title: res.msg || t('game.challenges.toast.submitSuccess') });
        return { success: true, message: t('game.challenges.toast.submitSuccessMessage') };
      }
      return { success: true, message: res.msg };
    } catch (error) {
      return { success: false, message: error.message };
    }
  };

  // 处理分类切换
  const handleCategoryChange = (category) => {
    setSelectedCategory(category);
    setCurrentPage(1); // 重置到第一页
    fetchChallengesWithFilters(1, category);
  };

  // 处理题目点击
  const handleChallengeClick = async (challenge) => {
    try {
      // 获取题目状态，包含附件、靶机信息和初始化状态
      const statusRes = await getChallengeStatus(contestId, challenge.id);
      if (statusRes.code === 200) {
        // 更新题目状态
        const updatedChallenge = {
          ...challenge,
          attachment: statusRes.data.file || '', // 如果有附件名称则添加
          isInitialized: statusRes.data.init, // 是否已初始化
          isSolved: statusRes.data.solved,
          instanceRunning: statusRes.data.remote.status === 'running',
          instancePending: statusRes.data.remote.status === 'pending',
          instanceIP: statusRes.data.remote.target || [''],
          instanceTimeLeft: statusRes.data.remote.remaining || 0,
          instanceDuration: statusRes.data.remote.duration || 3600,
          options: challenge.options || [], // 保持选项数据
        };

        setSelectedChallenge(updatedChallenge);

        // 页面刷新后 Pod 仍在启动中 → 自动开始轮询
        if (statusRes.data.remote.status === 'pending') {
          startPolling();
        }
      }
    } catch (error) {
      setSelectedChallenge(challenge); // 即使失败也显示题目
      toast.danger({ description: error.message || t('game.challenges.toast.fetchStatusFailed') });
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
      toast.danger({ description: error.message || t('game.challenges.toast.downloadFailed') });
    }
  };

  // 处理上传题解
  const handleUploadWriteup = async (file) => {
    try {
      setLoading(true);
      const res = await uploadWriteup(contestId, file);

      if (res.code === 200) {
        toast.success({
          title: t('game.challenges.toast.uploadSuccess'),
          description: t('game.challenges.toast.uploadThanks'),
        });
        // 上传成功后刷新题解列表
        fetchWriteups();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.uploadFailed') });
    }
  };

  // 处理查看题目
  const handleViewChallenges = () => {
    setShowChallengesAfterEnd(true);
  };

  // 处理返回比赛结束页面
  const handleBackToContestEnd = () => {
    setShowChallengesAfterEnd(false);
    setSelectedChallenge(null);
  };

  if (loading) {
    return <Loading />;
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
            {challenges.length === 0 ? (
              <EmptyState title={t('game.noChallenges')} description={t('game.noChallengesDescription')} />
            ) : (
              <ChallengeBoard
                categories={categories}
                selectedCategory={selectedCategory}
                onCategoryChange={handleCategoryChange}
                challenges={challenges}
                onChallengeClick={handleChallengeClick}
                teamInfo={teamInfo}
                totalCount={totalCount}
                currentPage={currentPage}
                pageSize={pageSize}
                onPageChange={setCurrentPage}
              />
            )}
            <ChallengeModal
              challenge={selectedChallenge}
              contest={contestStatus}
              isOpen={!!selectedChallenge}
              onClose={() => setSelectedChallenge(null)}
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
          {challenges.length === 0 ? (
            <EmptyState title={t('game.noChallenges')} description={t('game.noChallengesDescription')} />
          ) : (
            <ChallengeBoard
              categories={categories}
              selectedCategory={selectedCategory}
              onCategoryChange={handleCategoryChange}
              challenges={challenges}
              onChallengeClick={handleChallengeClick}
              teamInfo={teamInfo}
              totalCount={totalCount}
              currentPage={currentPage}
              pageSize={pageSize}
              onPageChange={setCurrentPage}
            />
          )}
          <ChallengeModal
            challenge={selectedChallenge}
            contest={contestStatus}
            isOpen={!!selectedChallenge}
            onClose={() => setSelectedChallenge(null)}
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

export default GameChallengesPage;
