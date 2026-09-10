import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import ScoreboardRanking from '../../components/features/Scoreboard/ScoreboardRanking';
import ScoreboardTable from '../../components/features/Scoreboard/ScoreboardTable';
import ScoreboardTimeline from '../../components/features/Scoreboard/ScoreboardTimeline';
import { collectChallenges, toRankingTeam } from '../../components/features/Scoreboard/scoreboardModel.js';
import { getContestRank, getContestInfo, getContestScoreboard, getContestTimeline } from '../../api/contest';
import { getTeamInfo } from '../../api/game/team';
import Button from '../../components/common/Button';
import Pagination from '../../components/common/Pagination';
import { IconList, IconTable, IconChartLine } from '@tabler/icons-react';
import ScoreboardStats from '../../components/features/Scoreboard/ScoreboardStats';
import { toast } from '../../utils/toast.js';
import { useTranslation } from 'react-i18next';

function GameScoreBoardPage() {
  const { contestId } = useParams();
  // A contest owns its pages, cached responses and timeline selection.
  return <ContestScoreboard key={contestId} contestId={contestId} />;
}

function ContestScoreboard({ contestId }) {
  const { t, i18n } = useTranslation();
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;
  const [totalCount, setTotalCount] = useState(0);
  const [viewMode, setViewMode] = useState('ranking');
  const [tableCurrentPage, setTableCurrentPage] = useState(1);
  const tablePageSize = 20;
  const [tableData, setTableData] = useState({ page: null, teams: [], challenges: [], count: 0 });
  // null distinguishes an unfetched timeline from a successfully fetched empty result.
  const [timelineData, setTimelineData] = useState(null);
  const [timelineLoading, setTimelineLoading] = useState(false);
  const [scoreboardData, setScoreboardData] = useState({
    stats: { totalTeams: 0, totalSolves: 0, highestScore: 0, totalPlayers: 0 },
    teams: [],
    userTeam: null,
  });

  useEffect(() => {
    let ignore = false;
    const teamTransform = (teamData, index = 0) =>
      toRankingTeam(teamData, index, currentPage, pageSize, i18n.language || 'en-US');

    Promise.all([
      getContestRank(contestId, pageSize, (currentPage - 1) * pageSize),
      getTeamInfo(contestId),
      getContestInfo(contestId),
    ])
      .then(([rankResponse, teamResponse, contestInfoResponse]) => {
        if (ignore) return;
        if (rankResponse.code === 200 && teamResponse.code === 200 && contestInfoResponse.code === 200) {
          setScoreboardData({
            stats: {
              totalTeams: contestInfoResponse.data.teams || 0,
              totalPlayers: contestInfoResponse.data.users || 0,
              totalSolves: contestInfoResponse.data.solved || 0,
              highestScore: contestInfoResponse.data.highest || 0,
            },
            teams: (rankResponse.data.teams || []).map(teamTransform),
            userTeam: teamResponse.data ? teamTransform(teamResponse.data) : null,
          });
          setTotalCount(rankResponse.data.count || 0);
        }
      })
      .catch((error) => {
        if (!ignore) toast.danger({ description: error.message || t('game.scoreboard.toast.fetchFailed') });
      });
    return () => {
      ignore = true;
    };
  }, [contestId, currentPage, i18n.language, t]);

  useEffect(() => {
    if (viewMode !== 'table' || tableData.page === tableCurrentPage) return;
    let ignore = false;
    setTableData({ page: null, teams: [], challenges: [], count: 0 });
    getContestScoreboard(contestId, {
      limit: tablePageSize,
      offset: (tableCurrentPage - 1) * tablePageSize,
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
        if (!ignore) toast.danger({ description: error.message || t('game.scoreboard.toast.fetchFailed') });
      });
    return () => {
      ignore = true;
    };
  }, [contestId, viewMode, tableCurrentPage]);

  useEffect(() => {
    if (viewMode !== 'timeline' || timelineData !== null) return;
    let ignore = false;
    setTimelineLoading(true);
    getContestTimeline(contestId)
      .then((response) => {
        if (!ignore && response.code === 200) setTimelineData(response.data || []);
      })
      .catch((error) => {
        if (!ignore) toast.danger({ description: error.message || t('game.scoreboard.toast.fetchFailed') });
      })
      .finally(() => {
        if (!ignore) setTimelineLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, [contestId, viewMode]);

  return (
    <div className="contest-container mx-auto min-w-0 space-y-6">
      <div className="flex flex-wrap justify-between items-center gap-3">
        <h1 className="text-xl sm:text-3xl font-mono text-neutral-50 tracking-wider">{t('nav.scoreboard')}</h1>
        <div
          role="group"
          aria-label={t('nav.scoreboard')}
          className="flex items-center gap-2 p-1 bg-black/30 border border-neutral-300/30 rounded-md"
        >
          {[
            { mode: 'ranking', label: t('common.rank'), icon: <IconList size={16} /> },
            { mode: 'table', label: t('game.scoreboard.headers.challenges'), icon: <IconTable size={16} /> },
            { mode: 'timeline', label: t('game.detail.labels.timeline'), icon: <IconChartLine size={16} /> },
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
      </div>
      <ScoreboardStats {...scoreboardData.stats} />
      {viewMode === 'timeline' ? (
        <ScoreboardTimeline timelineData={timelineData || []} loading={timelineLoading} />
      ) : viewMode === 'table' ? (
        <ScoreboardTable
          teams={tableData.teams}
          challenges={tableData.challenges}
          totalCount={tableData.count}
          currentPage={tableCurrentPage}
          pageSize={tablePageSize}
          onPageChange={setTableCurrentPage}
        />
      ) : (
        <ScoreboardRanking
          teams={scoreboardData.teams}
          locale={i18n.language || 'en-US'}
          labels={{
            rank: t('game.scoreboard.headers.rank'),
            team: t('game.scoreboard.headers.team'),
            score: t('game.scoreboard.headers.score'),
            challenges: t('game.scoreboard.headers.challenges'),
            lastSubmit: t('game.scoreboard.headers.lastSubmit'),
            total: t('game.scoreboard.total'),
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

export default GameScoreBoardPage;
