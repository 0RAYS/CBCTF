import { lazy, Suspense, useMemo, useState } from 'react';
import Button from '../../common/Button';
import Card from '../../common/Card';
import { useTranslation } from 'react-i18next';
import { buildTimelineChartData, buildTimelineModel } from './timelineModel.js';
import { buildTimelineOption, timelineTeamColor } from './timelineOption.js';
import TruncatedText from '../../common/TruncatedText';

const ReactECharts = lazy(() => import('../../common/EChart'));

/**
 * 分数时间线图表组件
 * @param {Object} props
 * @param {Array} props.timelineData - 时间线数据
 */
function ScoreboardTimeline({ timelineData = [], loading = false, error = false, onRetry }) {
  const { t, i18n } = useTranslation();
  const [hiddenTeams, setHiddenTeams] = useState(() => new Set());

  const model = useMemo(() => buildTimelineModel(timelineData), [timelineData]);
  const chartData = useMemo(() => buildTimelineChartData(model, hiddenTeams), [model, hiddenTeams]);
  const chartOption = buildTimelineOption(model.timePoints, chartData, i18n.language || 'en-US');

  // 切换队伍显示
  const toggleTeam = (teamId) => {
    setHiddenTeams((prev) => {
      const next = new Set(prev);
      if (next.has(teamId)) next.delete(teamId);
      else next.add(teamId);
      return next;
    });
  };

  return (
    <Card variant="default" padding="sm" className="min-w-0">
      <h2 className="mb-2 text-xs font-mono text-neutral-300">{t('game.scoreboard.timeline.title')}</h2>
      {/* 队伍选择器 */}
      {!loading && !error && model.teams.length > 0 && (
        <div
          className="flex max-h-20 flex-wrap gap-1.5 overflow-y-auto mb-2"
          role="group"
          aria-label={t('game.scoreboard.timeline.teams')}
        >
          {model.teams.map((team, index) => (
            <Button
              key={team.id}
              variant={!hiddenTeams.has(team.id) ? 'primary' : 'ghost'}
              aria-pressed={!hiddenTeams.has(team.id)}
              size="sm"
              onClick={() => toggleTeam(team.id)}
              className="!h-8 !px-2 !text-xs min-w-0 max-w-44"
              title={`#${team.rank} ${team.name}`}
              aria-label={`#${team.rank} ${team.name}`}
              style={{
                borderColor: !hiddenTeams.has(team.id) ? timelineTeamColor(index) : undefined,
                color: !hiddenTeams.has(team.id) ? timelineTeamColor(index) : undefined,
              }}
            >
              <TruncatedText>{`#${team.rank} ${team.name}`}</TruncatedText>
            </Button>
          ))}
        </div>
      )}

      {/* 图表 */}
      <div className="h-44 sm:h-56">
        {loading || error || model.teams.length === 0 ? (
          <div
            className="h-full flex flex-col gap-2 items-center justify-center text-sm text-neutral-400"
            role="status"
          >
            <span>{t(loading ? 'common.loading' : error ? 'game.scoreboard.timeline.failed' : 'common.noData')}</span>
            {error && !loading && onRetry && (
              <Button size="sm" variant="ghost" onClick={onRetry}>
                {t('common.retry')}
              </Button>
            )}
          </div>
        ) : (
          <Suspense
            fallback={
              <div className="h-full flex items-center justify-center text-neutral-400">{t('common.loading')}</div>
            }
          >
            <ReactECharts option={chartOption} replaceMerge={['series']} style={{ height: '100%', width: '100%' }} />
          </Suspense>
        )}
      </div>
    </Card>
  );
}

export default ScoreboardTimeline;
