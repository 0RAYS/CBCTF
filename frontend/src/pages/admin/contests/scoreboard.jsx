import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import { downloadBlobResponse } from '../../../utils/fileDownload';
import {
  getContestInfo,
  getContestRank,
  getContestScoreboard,
  getContestTimeline,
  getTeamMembers,
  getContestTeamSubmissions,
  getContestTeamWriteups,
  getTeamContainers,
  downloadContainerTraffic,
  downloadContestTeamWriteup,
  getContestTeamFlags,
} from '../../../api/admin/contest';
import AdminScoreboard from '../../../components/features/Admin/Contests/AdminScoreboard';
import AdminScoreboardTable from '../../../components/features/Admin/Contests/AdminScoreboardTable';
import ScoreboardTimeline from '../../../components/features/CTFGame/Scoreboard/ScoreboardTimeline';
import { Button } from '../../../components/common';
import { IconTable, IconList, IconChartLine } from '@tabler/icons-react';
import ScoreboardStats from '../../../components/features/CTFGame/Scoreboard/ScoreboardStats.jsx';
import { useTranslation } from 'react-i18next';

function AdminContestScoreboard({ viewMode: externalViewMode, onViewModeChange: externalOnViewModeChange }) {
  const { id } = useParams();

  // 视图状态
  const [viewMode, setViewMode] = useState(externalViewMode || 'ranking'); // 'ranking' | 'table' | 'timeline'

  // 排名视图状态
  const [teams, setTeams] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const [stats, setStats] = useState({
    totalTeams: 0,
    totalSolves: 0,
    highestScore: 0,
    totalPlayers: 0,
  });

  // 表格视图状态
  const [tableTeams, setTableTeams] = useState([]);
  const [challenges, setChallenges] = useState([]);
  const [tableCurrentPage, setTableCurrentPage] = useState(1);
  const [tableTotalCount, setTableTotalCount] = useState(0);
  const tablePageSize = 20;

  // 时间线相关状态
  const [timelineData, setTimelineData] = useState([]);
  const [timelineLoading, setTimelineLoading] = useState(false);
  const { t, i18n } = useTranslation();

  // Detail dialog state
  const [showDetailDialog, setShowDetailDialog] = useState(false);
  const [detailTeam, setDetailTeam] = useState(null);
  const [detailTab, setDetailTab] = useState('info');
  const [detailMembers, setDetailMembers] = useState([]);
  const [detailMembersLoading, setDetailMembersLoading] = useState(false);

  const [detailSubmissions, setDetailSubmissions] = useState([]);
  const [detailSubmissionCount, setDetailSubmissionCount] = useState(0);
  const [detailSubmissionPage, setDetailSubmissionPage] = useState(1);

  const [detailWriteups, setDetailWriteups] = useState([]);
  const [detailWriteupCount, setDetailWriteupCount] = useState(0);
  const [detailWriteupPage, setDetailWriteupPage] = useState(1);

  const [detailContainers, setDetailContainers] = useState([]);
  const [detailContainerCount, setDetailContainerCount] = useState(0);
  const [detailContainerPage, setDetailContainerPage] = useState(1);

  const [detailLoading, setDetailLoading] = useState({
    submissions: false,
    writeups: false,
    traffic: false,
  });

  const [detailFlags, setDetailFlags] = useState([]);
  const [detailFlagsLoading, setDetailFlagsLoading] = useState(false);

  const detailPageSize = 20;

  // 处理视图模式变化
  const handleViewModeChange = (mode) => {
    const newMode = mode;
    setViewMode(newMode);

    // 通知外部组件
    if (externalOnViewModeChange) {
      externalOnViewModeChange(newMode);
    }

    if (newMode === 'table') {
      setTableCurrentPage(1);
    } else if (newMode === 'ranking') {
      setCurrentPage(1);
    }
  };

  // 获取时间线数据
  const fetchTimelineData = async () => {
    setTimelineLoading(true);
    try {
      const response = await getContestTimeline(id);
      if (response.code === 200) {
        setTimelineData(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTimelineFailed') });
    } finally {
      setTimelineLoading(false);
    }
  };

  useEffect(() => {
    if (viewMode === 'ranking') {
      fetchRankings();
    } else if (viewMode === 'table') {
      fetchScoreboardTable();
    } else if (viewMode === 'timeline' && timelineData.length === 0) {
      fetchTimelineData();
    }
  }, [id, viewMode, currentPage, tableCurrentPage]);

  // 同步外部视图模式
  useEffect(() => {
    if (externalViewMode && externalViewMode !== viewMode) {
      setViewMode(externalViewMode);
    }
  }, [externalViewMode]);

  const teamTransform = (page, teamData, index) => {
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return '-';
      const date = new Date(dateString);

      // 格式化为：YYYY-MM-DD HH:MM:SS
      return date
        .toLocaleString(i18n.language || 'en-US', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
          hour12: false,
        })
        .replace(/\//g, '-');
    };

    return {
      id: teamData.id,
      rank: (page - 1) * pageSize + index + 1,
      name: teamData.name,
      picture: teamData.picture,
      score: teamData.score,
      solved: teamData.solved,
      totalSolved: teamData.solved.reduce((total, category) => total + category.solved, 0),
      lastSubmit: formatDate(teamData.last),
      captain_id: teamData.captain_id,
      captcha: teamData.captcha,
      description: teamData.description,
      users: teamData.users,
      banned: teamData.banned,
      hidden: teamData.hidden,
    };
  };

  const rankTransform = (page, rankData) => {
    return rankData.teams.map((v, index) => {
      return teamTransform(page, v, index);
    });
  };

  const fetchRankings = async (noLoading = false) => {
    try {
      const contestInfoResponse = await getContestInfo(parseInt(id));
      const response = await getContestRank(
        parseInt(id),
        {
          limit: pageSize,
          offset: (currentPage - 1) * pageSize,
        },
        noLoading
      );
      if (response.code === 200 && contestInfoResponse.code === 200) {
        setStats({
          totalTeams: contestInfoResponse.data.teams || 0,
          totalSolves: contestInfoResponse.data.solved || 0,
          highestScore: contestInfoResponse.data.highest || 0,
          totalPlayers: contestInfoResponse.data.users || 0,
        });

        setTeams(rankTransform(currentPage, response.data));
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    }
  };

  const fetchScoreboardTable = async () => {
    try {
      const response = await getContestScoreboard(parseInt(id), {
        limit: tablePageSize,
        offset: (tableCurrentPage - 1) * tablePageSize,
      });

      if (response.code === 200) {
        setTableTeams(response.data.teams || []);
        setTableTotalCount(response.data.count || 0);

        // 提取所有题目并按分类分组
        const allChallenges = [];
        const challengeMap = new Map();

        response.data.teams?.forEach((team) => {
          team.challenges?.forEach((challenge) => {
            if (!challengeMap.has(challenge.id)) {
              challengeMap.set(challenge.id, challenge);
              allChallenges.push(challenge);
            }
          });
        });

        // 按分类分组排序
        allChallenges.sort((a, b) => {
          if (a.category !== b.category) {
            return a.category.localeCompare(b.category);
          }
          return a.name.localeCompare(b.name);
        });

        setChallenges(allChallenges);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchScoreboardFailed') });
    }
  };

  const handleExportScoreboard = async () => {
    toast.info({ description: t('admin.contests.scoreboard.toast.exportUnavailable') });
  };

  // === Detail dialog handlers ===

  const fetchDetailSubmissions = async (teamId, page = 1) => {
    setDetailLoading((prev) => ({ ...prev, submissions: true }));
    try {
      const response = await getContestTeamSubmissions(parseInt(id), teamId, {
        limit: detailPageSize,
        offset: (page - 1) * detailPageSize,
      });
      if (response.code === 200) {
        setDetailSubmissions(response.data.submissions || []);
        setDetailSubmissionCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    } finally {
      setDetailLoading((prev) => ({ ...prev, submissions: false }));
    }
  };

  const fetchDetailWriteups = async (teamId, page = 1) => {
    setDetailLoading((prev) => ({ ...prev, writeups: true }));
    try {
      const response = await getContestTeamWriteups(parseInt(id), teamId, {
        limit: detailPageSize,
        offset: (page - 1) * detailPageSize,
      });
      if (response.code === 200) {
        setDetailWriteups(response.data.writeups || []);
        setDetailWriteupCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    } finally {
      setDetailLoading((prev) => ({ ...prev, writeups: false }));
    }
  };

  const fetchDetailContainers = async (teamId, page = 1) => {
    setDetailLoading((prev) => ({ ...prev, traffic: true }));
    try {
      const response = await getTeamContainers(parseInt(id), teamId, {
        limit: detailPageSize,
        offset: (page - 1) * detailPageSize,
      });
      if (response.code === 200) {
        setDetailContainers(response.data.victims || []);
        setDetailContainerCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    } finally {
      setDetailLoading((prev) => ({ ...prev, traffic: false }));
    }
  };

  const fetchDetailFlags = async (teamId) => {
    setDetailFlagsLoading(true);
    try {
      const response = await getContestTeamFlags(parseInt(id), teamId);
      if (response.code === 200) {
        setDetailFlags(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    } finally {
      setDetailFlagsLoading(false);
    }
  };

  const handleRowClick = async (team) => {
    if (!team.id) return;
    setDetailTeam(team);
    setDetailTab('info');
    setShowDetailDialog(true);
    setDetailMembersLoading(true);
    try {
      const response = await getTeamMembers(parseInt(id), team.id);
      if (response.code === 200) {
        setDetailMembers(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    } finally {
      setDetailMembersLoading(false);
    }
  };

  const handleDetailClose = () => {
    setShowDetailDialog(false);
    setDetailTeam(null);
    setDetailTab('info');
    setDetailMembers([]);
    setDetailSubmissions([]);
    setDetailSubmissionCount(0);
    setDetailSubmissionPage(1);
    setDetailWriteups([]);
    setDetailWriteupCount(0);
    setDetailWriteupPage(1);
    setDetailContainers([]);
    setDetailContainerCount(0);
    setDetailContainerPage(1);
    setDetailLoading({ submissions: false, writeups: false, traffic: false });
    setDetailFlags([]);
    setDetailFlagsLoading(false);
  };

  const handleDetailTabChange = (tab) => {
    setDetailTab(tab);
    if (!detailTeam) return;

    if (tab === 'submissions') {
      setDetailSubmissionPage(1);
      fetchDetailSubmissions(detailTeam.id, 1);
    } else if (tab === 'writeups') {
      setDetailWriteupPage(1);
      fetchDetailWriteups(detailTeam.id, 1);
    } else if (tab === 'containers') {
      setDetailContainerPage(1);
      fetchDetailContainers(detailTeam.id, 1);
    } else if (tab === 'flags') {
      fetchDetailFlags(detailTeam.id);
    }
  };

  const handleDetailPageChange = (type, page) => {
    if (!detailTeam) return;
    if (type === 'submissions') {
      setDetailSubmissionPage(page);
      fetchDetailSubmissions(detailTeam.id, page);
    } else if (type === 'writeups') {
      setDetailWriteupPage(page);
      fetchDetailWriteups(detailTeam.id, page);
    } else if (type === 'containers') {
      setDetailContainerPage(page);
      fetchDetailContainers(detailTeam.id, page);
    }
  };

  const handleDetailDownloadTraffic = async (container) => {
    if (!detailTeam) return;
    try {
      const response = await downloadContainerTraffic(parseInt(id), detailTeam.id, container.id);
      downloadBlobResponse(response, `traffic_${container.id}.zip`);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    }
  };

  const handleDetailDownloadWriteup = async (writeup) => {
    if (!detailTeam) return;
    try {
      const response = await downloadContestTeamWriteup(parseInt(id), detailTeam.id, writeup.id);
      downloadBlobResponse(response, writeup.filename);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    }
  };

  return (
    <div className="w-full mx-auto space-y-6">
      {/* 头部和视图切换 */}
      <div className="flex justify-end items-center">
        <div className="flex items-center gap-4">
          {/* 视图切换按钮 */}
          <div className="flex items-center gap-2 p-1 bg-black/30 border border-neutral-300/30 rounded-md">
            <Button
              variant={viewMode === 'ranking' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconList size={16} />}
              onClick={() => handleViewModeChange('ranking')}
            />
            <Button
              variant={viewMode === 'table' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconTable size={16} />}
              onClick={() => handleViewModeChange('table')}
            />
            <Button
              variant={viewMode === 'timeline' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconChartLine size={16} />}
              onClick={() => handleViewModeChange('timeline')}
            />
          </div>

          {/* 导出按钮 */}
          <Button variant="outline" size="sm" onClick={handleExportScoreboard}>
            {t('admin.contests.scoreboard.export')}
          </Button>
        </div>
      </div>

      {/* 视图内容 */}
      <ScoreboardStats {...stats} />
      {viewMode === 'ranking' ? (
        <AdminScoreboard
          teams={teams}
          currentPage={currentPage}
          pageSize={pageSize}
          totalCount={totalCount}
          onPageChange={setCurrentPage}
          onExportScoreboard={handleExportScoreboard}
          viewMode={viewMode}
          onViewModeChange={handleViewModeChange}
          onRowClick={handleRowClick}
          showDetailDialog={showDetailDialog}
          detailTeam={detailTeam}
          detailTab={detailTab}
          detailMembers={detailMembers}
          detailMembersLoading={detailMembersLoading}
          detailSubmissions={detailSubmissions}
          detailSubmissionCount={detailSubmissionCount}
          detailSubmissionPage={detailSubmissionPage}
          detailWriteups={detailWriteups}
          detailWriteupCount={detailWriteupCount}
          detailWriteupPage={detailWriteupPage}
          detailContainers={detailContainers}
          detailContainerCount={detailContainerCount}
          detailContainerPage={detailContainerPage}
          detailLoading={detailLoading}
          onDetailClose={handleDetailClose}
          onDetailTabChange={handleDetailTabChange}
          onDetailPageChange={handleDetailPageChange}
          onDetailDownloadTraffic={handleDetailDownloadTraffic}
          onDetailDownloadWriteup={handleDetailDownloadWriteup}
          detailFlags={detailFlags}
          detailFlagsLoading={detailFlagsLoading}
        />
      ) : viewMode === 'table' ? (
        <AdminScoreboardTable
          teams={tableTeams}
          challenges={challenges}
          totalCount={tableTotalCount}
          currentPage={tableCurrentPage}
          pageSize={tablePageSize}
          onPageChange={setTableCurrentPage}
        />
      ) : (
        <ScoreboardTimeline timelineData={timelineData} loading={timelineLoading} />
      )}
    </div>
  );
}

export default AdminContestScoreboard;
