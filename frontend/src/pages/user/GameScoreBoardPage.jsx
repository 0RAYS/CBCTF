import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Scoreboard from '../../components/features/CTFGame/Scoreboard/Scoreboard';
import ScoreboardTimeline from '../../components/features/CTFGame/Scoreboard/ScoreboardTimeline';
import { getContestRank, getContestInfo, getContestScoreboard, getContestTimeline } from '../../api/contest';
import { getTeamInfo } from '../../api/game/team';
import { Button } from '../../components/common';
import { IconList, IconTable, IconChartLine } from '@tabler/icons-react';
import ScoreboardStats from '../../components/features/CTFGame/Scoreboard/ScoreboardStats';
import { toast } from '../../utils/toast.js';
import { useTranslation } from 'react-i18next';

function GameScoreBoardPage() {
  const { contestId } = useParams();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20); // 每页显示20条数据
  const [totalPages, setTotalPages] = useState(0);
  const [totalCount, setTotalCount] = useState(0); // 添加总数状态
  const [viewMode, setViewMode] = useState('ranking'); // 'ranking' | 'table' | 'timeline'

  // 表格视图相关状态
  const [tableTeams, setTableTeams] = useState([]);
  const [tableChallenges, setTableChallenges] = useState([]);
  const [tableCurrentPage, setTableCurrentPage] = useState(1);
  const [tableTotalCount, setTableTotalCount] = useState(0);
  const tablePageSize = 20;

  // 时间线相关状态
  const [timelineData, setTimelineData] = useState([]);
  const { t, i18n } = useTranslation();

  const [scoreboardData, setScoreboardData] = useState({
    stats: {
      totalTeams: 0,
      totalSolves: 0,
      highestScore: 0,
      totalPlayers: 0,
    },
    teams: [],
    userTeam: null,
  });

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
      rank: (page - 1) * pageSize + index + 1,
      name: teamData.name,
      picture: teamData.picture,
      score: teamData.score,
      solved: teamData.solved,
      totalSolved: teamData.solved.reduce((total, category) => total + category.solved, 0),
      lastSubmit: formatDate(teamData.last),
    };
  };

  const createChallengesFromTeams = (teams) => {
    const challengeMap = new Map();

    teams.forEach((team) => {
      if (team.solved && Array.isArray(team.solved)) {
        team.solved.forEach((category) => {
          // 为每个分类创建模拟的题目
          const categoryName = category.category;
          const totalCount = category.all || 0;

          // 创建该分类的题目
          for (let i = 0; i < totalCount; i++) {
            const challengeId = `${categoryName}-${i + 1}`;
            if (!challengeMap.has(challengeId)) {
              challengeMap.set(challengeId, {
                id: challengeId,
                name: `${categoryName} ${i + 1}`,
                category: categoryName,
                solved: 0,
              });
            }
          }
        });
      }
    });

    return Array.from(challengeMap.values()).sort((a, b) => {
      if (a.category !== b.category) {
        return a.category.localeCompare(b.category);
      }
      return a.name.localeCompare(b.name);
    });
  };

  const rankTransform = (page, rankData) => {
    return rankData.teams.map((v, index) => {
      return teamTransform(page, v, index);
    });
  };

  // 处理页码变化
  const handlePageChange = (newPage) => {
    setCurrentPage(newPage);
    fetchScoreboardData(newPage);
  };

  // 处理表格视图页码变化
  const handleTablePageChange = (newPage) => {
    setTableCurrentPage(newPage);
    fetchTableData(newPage);
  };

  // 获取表格视图数据
  const fetchTableData = async (newPage) => {
    try {
      const response = await getContestScoreboard(contestId, {
        limit: tablePageSize,
        offset: (newPage - 1) * tablePageSize,
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

        setTableChallenges(allChallenges);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.scoreboard.toast.fetchFailed') });
    }
  };

  const fetchScoreboardData = async (newPage) => {
    try {
      const rankResponse = await getContestRank(contestId, pageSize, (newPage - 1) * pageSize);
      const teamResponse = await getTeamInfo(contestId);
      const contestInfoResponse = await getContestInfo(contestId);
      if (rankResponse.code === 200 && teamResponse.code === 200) {
        setScoreboardData({
          stats: {
            totalTeams: contestInfoResponse.data.teams || 0,
            totalPlayers: contestInfoResponse.data.users || 0,
            totalSolves: contestInfoResponse.data.solved || 0,
            highestScore: contestInfoResponse.data.highest || 0,
          },
          teams: rankTransform(newPage, rankResponse.data),
          userTeam: teamTransform(newPage, teamResponse.data),
        });
        setTotalPages(Math.ceil(rankResponse.data.count / pageSize));
        setTotalCount(rankResponse.data.count); // 更新总数
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.scoreboard.toast.fetchFailed') });
    }
  };

  useEffect(() => {
    fetchScoreboardData(1);
  }, [contestId]);

  // 当视图模式改变时, 获取相应的数据
  useEffect(() => {
    let ignore = false;
    if (viewMode === 'table' && tableTeams.length === 0) {
      getContestScoreboard(contestId, {
        limit: tablePageSize,
        offset: 0,
      })
        .then((response) => {
          if (ignore) return;
          if (response.code === 200) {
            setTableTeams(response.data.teams || []);
            setTableTotalCount(response.data.count || 0);

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
            allChallenges.sort((a, b) => {
              if (a.category !== b.category) return a.category.localeCompare(b.category);
              return a.name.localeCompare(b.name);
            });
            setTableChallenges(allChallenges);
          }
        })
        .catch((error) => {
          if (!ignore) toast.danger({ description: error.message || t('game.scoreboard.toast.fetchFailed') });
        });
    } else if (viewMode === 'timeline' && timelineData.length === 0) {
      getContestTimeline(contestId)
        .then((response) => {
          if (ignore) return;
          if (response.code === 200) {
            setTimelineData(response.data || []);
          }
        })
        .catch((error) => {
          if (!ignore) toast.danger({ description: error.message || t('game.scoreboard.toast.fetchFailed') });
        });
    }
    return () => {
      ignore = true;
    };
  }, [viewMode]);

  // 准备表格视图的数据
  const challenges = createChallengesFromTeams(scoreboardData.teams);

  // 渲染时间线视图
  if (viewMode === 'timeline') {
    return (
      <div className="contest-container mx-auto space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-mono text-neutral-50 tracking-wider"></h1>

          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 p-1 bg-black/30 border border-neutral-300/30 rounded-md">
              <Button
                variant={viewMode === 'ranking' ? 'primary' : 'ghost'}
                size="sm"
                align="icon-left"
                icon={<IconList size={16} />}
                onClick={() => setViewMode('ranking')}
              />
              <Button
                variant={viewMode === 'table' ? 'primary' : 'ghost'}
                size="sm"
                align="icon-left"
                icon={<IconTable size={16} />}
                onClick={() => setViewMode('table')}
              />
              <Button
                variant={viewMode === 'timeline' ? 'primary' : 'ghost'}
                size="sm"
                align="icon-left"
                icon={<IconChartLine size={16} />}
                onClick={() => setViewMode('timeline')}
              />
            </div>
          </div>
        </div>

        {/* 总览面板 */}
        <ScoreboardStats {...scoreboardData.stats} />

        <ScoreboardTimeline timelineData={timelineData} />
      </div>
    );
  }

  return (
    <div className="contest-container mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-mono text-neutral-50 tracking-wider"></h1>

        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 p-1 bg-black/30 border border-neutral-300/30 rounded-md">
            <Button
              variant={viewMode === 'ranking' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconList size={16} />}
              onClick={() => setViewMode('ranking')}
            />
            <Button
              variant={viewMode === 'table' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconTable size={16} />}
              onClick={() => setViewMode('table')}
            />
            <Button
              variant={viewMode === 'timeline' ? 'primary' : 'ghost'}
              size="sm"
              align="icon-left"
              icon={<IconChartLine size={16} />}
              onClick={() => setViewMode('timeline')}
            />
          </div>
        </div>
      </div>

      {/* 总览面板 */}
      <ScoreboardStats {...scoreboardData.stats} />

      {/* 用户队伍状态 */}
      <Scoreboard
        {...scoreboardData}
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={handlePageChange}
        viewMode={viewMode}
        onViewModeChange={setViewMode}
        challenges={viewMode === 'table' ? tableChallenges : challenges}
        totalCount={viewMode === 'table' ? tableTotalCount : totalCount}
        teams={viewMode === 'table' ? tableTeams : scoreboardData.teams}
        tableCurrentPage={tableCurrentPage}
        onTablePageChange={handleTablePageChange}
      />
    </div>
  );
}

export default GameScoreBoardPage;
