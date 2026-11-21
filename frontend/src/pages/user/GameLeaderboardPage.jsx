import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import LeaderboardRanking from '@/components/features/Leaderboard/LeaderboardRanking';
import LeaderboardTable from '@/components/features/Leaderboard/LeaderboardTable';
import LeaderboardTimeline from '@/components/features/Leaderboard/LeaderboardTimeline';
import { collectChallenges, toRankingTeam } from '@/components/features/Leaderboard/leaderboardModel.js';
import {
  getContestInfo,
  getContestLeaderboard,
  getContestLeaderboardChallenges,
  getContestLeaderboardTimeline,
} from '@/api/contest';
import { getTeamInfo } from '@/api/game/team';
import Button from '@/components/common/Button';
import Pagination from '@/components/common/Pagination';
import { IconList, IconTable } from '@tabler/icons-react';
import LeaderboardStats from '@/components/features/Leaderboard/LeaderboardStats';
import { toast } from '@/utils/toast.js';
import { useTranslation } from 'react-i18next';
import { useHiddenVisibility } from '@/hooks/useHiddenVisibility.js';
import Checkbox from '@/components/common/Checkbox';

function GameLeaderboardPage() {
  const { contestId } = useParams();
  // A contest owns its pages, cached responses and timeline selection.
  return <ContestLeaderboard key={contestId} contestId={contestId} />;
}

function ContestLeaderboard({ contestId }) {
  const { t, i18n } = useTranslation();
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const [totalCount, setTotalCount] = useState(0);
  const [viewMode, setViewMode] = useState('ranking');
  const [tableCurrentPage, setTableCurrentPage] = useState(1);
  const tablePageSize = 20;
  const [tableData, setTableData] = useState({ page: null, teams: [], challenges: [], count: 0 });
  const [rankingLoading, setRankingLoading] = useState(true);
  const tableLoading = viewMode === 'table' && tableData.page !== tableCurrentPage;
  // null distinguishes an unfetched timeline from a successfully fetched empty result.
  const [timelineData, setTimelineData] = useState(null);
  const [timelineLoading, setTimelineLoading] = useState(true);
  const [timelineError, setTimelineError] = useState(false);
  const [timelineVersion, setTimelineVersion] = useState(0);
  const [showHidden, setShowHidden] = useHiddenVisibility('cbctf.leaderboard.showHiddenTeams');
  const [leaderboardData, setLeaderboardData] = useState({
    stats: { totalTeams: 0, totalSolves: 0, highestScore: 0, totalPlayers: 0 },
    teams: [],
    userTeam: null,
  });

  useEffect(() => {
    let ignore = false;
    const teamTransform = (teamData, index = 0) =>
      toRankingTeam(teamData, index, currentPage, pageSize, i18n.language || 'en-US');

    setRankingLoading(true);
    Promise.all([
      getContestLeaderboard(contestId, pageSize, (currentPage - 1) * pageSize, showHidden),
      getTeamInfo(contestId),
      getContestInfo(contestId),
    ])
      .then(([rankResponse, teamResponse, contestInfoResponse]) => {
        if (ignore) return;
        if (rankResponse.code === 200 && teamResponse.code === 200 && contestInfoResponse.code === 200) {
          setLeaderboardData({
            stats: {
              totalTeams: contestInfoResponse.data.teams,
              totalPlayers: contestInfoResponse.data.users,
              totalSolves: contestInfoResponse.data.solved,
              highestScore: contestInfoResponse.data.highest,
            },
            teams: (rankResponse.data.teams || []).map(teamTransform),
            userTeam: teamResponse.data ? teamTransform(teamResponse.data) : null,
          });
          setTotalCount(rankResponse.data.count || 0);
        }
      })
      .catch((error) => {
        if (!ignore) toast.danger({ description: error.message || t('game.leaderboard.toast.fetchFailed') });
      })
      .finally(() => {
        if (!ignore) setRankingLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, [contestId, currentPage, showHidden, i18n.language, t]);

  useEffect(() => {
    if (viewMode !== 'table' || tableData.page === tableCurrentPage) return;
    let ignore = false;
    setTableData({ page: null, teams: [], challenges: [], count: 0 });
    getContestLeaderboardChallenges(contestId, {
      limit: tablePageSize,
      offset: (tableCurrentPage - 1) * tablePageSize,
      show_hidden: showHidden ? 1 : 0,
    })
      .then((response) => {
        if (ignore || response.code !== 200) return;
        setTableData({
          page: tableCurrentPage,
          teams: response.data.teams || [],
          challenges: collectChallenges(response.data.teams || []),
          count: response.data.count || 0,
        });
      })
      .catch((error) => {
        if (!ignore) toast.danger({ description: error.message || t('game.leaderboard.toast.fetchFailed') });
      });
    return () => {
      ignore = true;
    };
  }, [contestId, viewMode, tableCurrentPage, showHidden]);

  useEffect(() => {
    let ignore = false;
    setTimelineLoading(true);
    setTimelineError(false);
    getContestLeaderboardTimeline(contestId, showHidden)
      .then((response) => {
        if (ignore) return;
        if (response.code === 200) setTimelineData(response.data || []);
        else setTimelineError(true);
      })
      .catch((error) => {
        if (!ignore) {
          setTimelineError(true);
          toast.danger({ description: error.message || t('game.leaderboard.toast.fetchFailed') });
        }
      })
      .finally(() => {
        if (!ignore) setTimelineLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, [contestId, showHidden, timelineVersion]);

  const handleShowHiddenChange = (event) => {
    setShowHidden(event.target.checked);
    setCurrentPage(1);
    setTableCurrentPage(1);
  };

  return (
    <div className="contest-container mx-auto min-w-0 space-y-4">
      <div className="flex flex-wrap justify-between items-center gap-3">
        <h1 className="text-xl sm:text-3xl font-mono text-neutral-50 tracking-wider">{t('nav.leaderboard')}</h1>
        <div className="flex flex-wrap items-center gap-3">
          <div
            role="group"
            aria-label={t('nav.leaderboard')}
            className="flex items-center gap-2 p-1 bg-black/30 border border-neutral-300/30 rounded-md"
          >
            {[
              { mode: 'ranking', label: t('common.rank'), icon: <IconList size={16} /> },
              { mode: 'table', label: t('common.leaderboard.headers.challenges'), icon: <IconTable size={16} /> },
            ].map(({ mode, label, icon }) => (
              <Button
                key={mode}
                variant={viewMode === mode ? 'primary' : 'ghost'}
                size="sm"
                align="icon-left"
                icon={icon}
                aria-label={label}
                title={label}
                aria-pressed={viewMode === mode}
                onClick={() => setViewMode(mode)}
              />
            ))}
          </div>
          <Checkbox
            checked={showHidden}
            onChange={handleShowHiddenChange}
            label={t('game.leaderboard.showHiddenTeams')}
            className="text-xs text-neutral-400"
          />
        </div>
      </div>
      <LeaderboardStats {...leaderboardData.stats} />
      <LeaderboardTimeline
        timelineData={timelineData || []}
        loading={timelineLoading}
        error={timelineError}
        onRetry={() => setTimelineVersion((version) => version + 1)}
      />
      {viewMode === 'table' ? (
        <LeaderboardTable
          teams={tableData.teams}
          challenges={tableData.challenges}
          totalCount={tableData.count}
          currentPage={tableCurrentPage}
          pageSize={tablePageSize}
          loading={tableLoading}
          onPageChange={setTableCurrentPage}
        />
      ) : (
        <LeaderboardRanking
          teams={leaderboardData.teams}
          loading={rankingLoading}
          locale={i18n.language || 'en-US'}
          labels={{
            rank: t('common.leaderboard.headers.rank'),
            team: t('common.leaderboard.headers.team'),
            score: t('common.leaderboard.headers.score'),
            challenges: t('common.leaderboard.headers.challenges'),
            lastSubmit: t('common.leaderboard.headers.lastSubmit'),
            total: t('common.leaderboard.total'),
          }}
          emptyMessage={t('common.noData')}
          footer={
            totalCount > pageSize ? (
              <div className="mt-6 flex justify-center">
                <Pagination
                  current={currentPage}
                  total={Math.ceil(totalCount / pageSize)}
                  onChange={setCurrentPage}
                  showTotal={true}
                  totalItems={totalCount}
                  animate
                />
              </div>
            ) : null
          }
        />
      )}
    </div>
  );
}

export default GameLeaderboardPage;
