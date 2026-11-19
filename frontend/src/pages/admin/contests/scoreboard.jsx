import { useEffect, useState } from 'react';
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
import { IconList, IconTable } from '@tabler/icons-react';
import ScoreboardStats from '../../../components/features/Scoreboard/ScoreboardStats.jsx';
import { useTranslation } from 'react-i18next';
import { useTeamDetailDialog } from '../../../components/features/Admin/details/useTeamDetailDialog.jsx';

function AdminContestScoreboard(props) {
  const { id } = useParams();
  return <ContestScoreboard key={id} id={id} {...props} />;
}

function ContestScoreboard({ id, viewMode: externalViewMode, onViewModeChange: externalOnViewModeChange }) {
  const [viewMode, setViewMode] = useState(externalViewMode === 'table' ? 'table' : 'ranking');

  const [teams, setTeams] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [rankingLoading, setRankingLoading] = useState(true);
  const pageSize = 20;
  const [stats, setStats] = useState({
    totalTeams: 0,
    totalSolves: 0,
    highestScore: 0,
    totalPlayers: 0,
  });

  const [tableTeams, setTableTeams] = useState([]);
  const [challenges, setChallenges] = useState([]);
  const [tableCurrentPage, setTableCurrentPage] = useState(1);
  const [tableTotalCount, setTableTotalCount] = useState(0);
  const [tableLoading, setTableLoading] = useState(false);
  const tablePageSize = 20;

  const [timelineData, setTimelineData] = useState(null);
  const [timelineLoading, setTimelineLoading] = useState(true);
  const [timelineError, setTimelineError] = useState(false);
  const [timelineVersion, setTimelineVersion] = useState(0);
  const { t, i18n } = useTranslation();

  const { openTeamDetail, renderTeamDetailDialog } = useTeamDetailDialog(parseInt(id));

  const handleViewModeChange = (mode) => {
    const newMode = mode;
    setViewMode(newMode);

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

  const fetchTimelineData = async (isCurrent) => {
    setTimelineLoading(true);
    setTimelineError(false);
    try {
      const response = await getContestTimeline(id);
      if (isCurrent()) {
        if (response.code === 200) setTimelineData(response.data || []);
        else setTimelineError(true);
      }
    } catch (error) {
      if (isCurrent()) {
        setTimelineError(true);
        toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTimelineFailed') });
      }
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
        setTeams(response.data.teams.map((team, index) => teamTransform(currentPage, team, index)));
        setTotalCount(response.data.count);
      }
    } catch (error) {
      if (isCurrent())
        toast.danger({ description: error.message || t('admin.contests.scoreboard.toast.fetchTeamsFailed') });
    } finally {
      if (isCurrent()) setRankingLoading(false);
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
    } finally {
      if (isCurrent()) setTableLoading(false);
    }
  };

  // Statistics belong to the contest, not the active view or ranking page.
  useEffect(() => {
    let ignore = false;
    getContestInfo(parseInt(id))
      .then((response) => {
        if (ignore || response.code !== 200) return;
        setStats({
          totalTeams: response.data.teams,
          totalSolves: response.data.solved,
          highestScore: response.data.highest,
          totalPlayers: response.data.users,
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
      setRankingLoading(true);
      fetchRankings(isCurrent);
    } else if (viewMode === 'table') {
      setTableLoading(true);
      fetchScoreboardTable(isCurrent);
    }
    return () => {
      ignore = true;
    };
  }, [id, viewMode, currentPage, tableCurrentPage, i18n.language]);

  // The shared chart has its own request lifetime; list switches and pagination cannot cancel it.
  useEffect(() => {
    let ignore = false;
    fetchTimelineData(() => !ignore);
    return () => {
      ignore = true;
    };
  }, [id, timelineVersion]);

  useEffect(() => {
    if (['ranking', 'table'].includes(externalViewMode) && externalViewMode !== viewMode) {
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
    <div className="w-full min-w-0 mx-auto space-y-4">
      <div className="flex flex-wrap justify-between items-center gap-3">
        <h1 className="text-xl font-mono text-neutral-50">{t('nav.scoreboard')}</h1>
        <div className="flex flex-wrap items-center gap-3">
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
              aria-label={t('common.scoreboard.headers.challenges')}
              title={t('common.scoreboard.headers.challenges')}
              aria-pressed={viewMode === 'table'}
              onClick={() => handleViewModeChange('table')}
            />
          </div>

          <Button variant="outline" size="sm" onClick={handleExportScoreboard}>
            {t('admin.contests.scoreboard.export')}
          </Button>
        </div>
      </div>

      <ScoreboardStats {...stats} />
      <ScoreboardTimeline
        timelineData={timelineData || []}
        loading={timelineLoading}
        error={timelineError}
        onRetry={() => setTimelineVersion((version) => version + 1)}
      />
      {viewMode === 'ranking' ? (
        <AdminRanking
          teams={teams}
          currentPage={currentPage}
          pageSize={pageSize}
          totalCount={totalCount}
          loading={rankingLoading}
          onPageChange={setCurrentPage}
          onRowClick={handleRowClick}
        />
      ) : (
        <ScoreboardTable
          teams={tableTeams}
          challenges={challenges}
          totalCount={tableTotalCount}
          currentPage={tableCurrentPage}
          pageSize={tablePageSize}
          loading={tableLoading}
          onPageChange={setTableCurrentPage}
        />
      )}

      {renderTeamDetailDialog()}
    </div>
  );
}

export default AdminContestScoreboard;
