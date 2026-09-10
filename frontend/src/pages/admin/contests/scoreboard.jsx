import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import {
  exportContestWriteups,
  getContestInfo,
  getContestRank,
  getContestScoreboard,
  getContestTimeline,
} from '../../../api/admin/contest';
import { downloadBlobResponse } from '../../../utils/fileDownload';
import AdminRanking from '../../../components/features/Scoreboard/AdminRanking';
import ScoreboardTable from '../../../components/features/Scoreboard/ScoreboardTable';
import ScoreboardTimeline from '../../../components/features/Scoreboard/ScoreboardTimeline';
import { collectChallenges, toRankingTeam } from '../../../components/features/Scoreboard/scoreboardModel.js';
import Button from '../../../components/common/Button';
import { IconTable, IconList, IconChartLine } from '@tabler/icons-react';
import ScoreboardStats from '../../../components/features/Scoreboard/ScoreboardStats.jsx';
import { useTranslation } from 'react-i18next';
import { useTeamDetailDialog } from '../../../components/features/Admin/details/useTeamDetailDialog.jsx';

function AdminContestScoreboard(props) {
  const { id } = useParams();
  return <ContestScoreboard key={id} id={id} {...props} />;
}

function ContestScoreboard({ id, viewMode: externalViewMode, onViewModeChange: externalOnViewModeChange }) {
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
  const [timelineData, setTimelineData] = useState(null);
  const [timelineLoading, setTimelineLoading] = useState(false);
  const { t, i18n } = useTranslation();

  const { openTeamDetail, renderTeamDetailDialog } = useTeamDetailDialog(parseInt(id));

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

  const teamTransform = (page, teamData, index) => {
    return {
      ...toRankingTeam(teamData, index, page, pageSize, i18n.language || 'en-US'),
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

  // 获取时间线数据
  const fetchTimelineData = async (isCurrent) => {
    setTimelineLoading(true);
    try {
      const response = await getContestTimeline(id);
      if (isCurrent() && response.code === 200) {
        setTimelineData(response.data || []);
      }
    } catch (error) {
      if (isCurrent())
        toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTimelineFailed') });
    } finally {
      if (isCurrent()) setTimelineLoading(false);
    }
  };

  const fetchRankings = async (isCurrent) => {
    try {
      const response = await getContestRank(parseInt(id), {
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      });
      if (isCurrent() && response.code === 200) {
        setTeams(rankTransform(currentPage, response.data));
        setTotalCount(response.data.count);
      }
    } catch (error) {
      if (isCurrent())
        toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    }
  };

  const fetchScoreboardTable = async (isCurrent) => {
    setTableTeams([]);
    setChallenges([]);
    try {
      const response = await getContestScoreboard(parseInt(id), {
        limit: tablePageSize,
        offset: (tableCurrentPage - 1) * tablePageSize,
      });

      if (isCurrent() && response.code === 200) {
        setTableTeams(response.data.teams || []);
        setTableTotalCount(response.data.count || 0);

        setChallenges(collectChallenges(response.data.teams || []));
      }
    } catch (error) {
      if (isCurrent())
        toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchScoreboardFailed') });
    }
  };

  // Statistics belong to the contest, not the active view or ranking page.
  useEffect(() => {
    let ignore = false;
    getContestInfo(parseInt(id))
      .then((response) => {
        if (ignore || response.code !== 200) return;
        setStats({
          totalTeams: response.data.teams || 0,
          totalSolves: response.data.solved || 0,
          highestScore: response.data.highest || 0,
          totalPlayers: response.data.users || 0,
        });
      })
      .catch((error) => {
        if (!ignore)
          toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchScoreboardFailed') });
      });
    return () => {
      ignore = true;
    };
  }, [id]);

  useEffect(() => {
    let ignore = false;
    const isCurrent = () => !ignore;
    if (viewMode === 'ranking') {
      fetchRankings(isCurrent);
    } else if (viewMode === 'table') {
      fetchScoreboardTable(isCurrent);
    } else if (viewMode === 'timeline' && timelineData === null) {
      fetchTimelineData(isCurrent);
    }
    return () => {
      ignore = true;
    };
  }, [id, viewMode, currentPage, tableCurrentPage, i18n.language]);

  // 同步外部视图模式
  useEffect(() => {
    if (externalViewMode && externalViewMode !== viewMode) {
      setViewMode(externalViewMode);
    }
  }, [externalViewMode]);

  const handleExportScoreboard = async () => {
    try {
      const response = await exportContestWriteups(parseInt(id));
      downloadBlobResponse(response, `contest-${id}.zip`, 'application/zip');
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.exportWriteupsFailed') });
    }
  };

  const handleRowClick = (team) => {
    if (!team.id) return;
    openTeamDetail(team);
  };

  return (
    <div className="w-full mx-auto space-y-6">
      {/* 头部和视图切换 */}
      <div className="flex flex-wrap justify-between items-center gap-3">
        <h1 className="text-xl font-mono text-neutral-50">{t('nav.scoreboard')}</h1>
        <div className="flex flex-wrap items-center gap-3">
          {/* 视图切换按钮 */}
          <div
            role="group"
            aria-label={t('nav.scoreboard')}
            className="flex items-center gap-2 p-1 bg-black/30 border border-neutral-300/30 rounded-md"
          >
            <Button
              variant={viewMode === 'ranking' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconList size={16} />}
              aria-label={t('common.rank')}
              title={t('common.rank')}
              aria-pressed={viewMode === 'ranking'}
              onClick={() => handleViewModeChange('ranking')}
            />
            <Button
              variant={viewMode === 'table' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconTable size={16} />}
              aria-label={t('game.scoreboard.headers.challenges')}
              title={t('game.scoreboard.headers.challenges')}
              aria-pressed={viewMode === 'table'}
              onClick={() => handleViewModeChange('table')}
            />
            <Button
              variant={viewMode === 'timeline' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconChartLine size={16} />}
              aria-label={t('game.detail.labels.timeline')}
              title={t('game.detail.labels.timeline')}
              aria-pressed={viewMode === 'timeline'}
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
        <AdminRanking
          teams={teams}
          currentPage={currentPage}
          pageSize={pageSize}
          totalCount={totalCount}
          onPageChange={setCurrentPage}
          onRowClick={handleRowClick}
        />
      ) : viewMode === 'table' ? (
        <ScoreboardTable
          teams={tableTeams}
          challenges={challenges}
          totalCount={tableTotalCount}
          currentPage={tableCurrentPage}
          pageSize={tablePageSize}
          onPageChange={setTableCurrentPage}
        />
      ) : (
        <ScoreboardTimeline timelineData={timelineData || []} loading={timelineLoading} />
      )}

      {renderTeamDetailDialog()}
    </div>
  );
}

export default AdminContestScoreboard;
