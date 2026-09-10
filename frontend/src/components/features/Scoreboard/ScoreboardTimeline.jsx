import { Suspense, lazy, useMemo, useState } from 'react';
import Button from '../../common/Button';
import Card from '../../common/Card';
import { useTranslation } from 'react-i18next';
import { buildTimelineModel, buildTimelineChartData } from './timelineModel.js';
import { buildTimelineOption, timelineTeamColor } from './timelineOption.js';

const ReactECharts = lazy(() => import('../../common/EChart'));

/**
 * 分数时间线图表组件
 * @param {Object} props
 * @param {Array} props.timelineData - 时间线数据
 */
function ScoreboardTimeline({ timelineData = [], loading = false }) {
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

  if (loading || !timelineData || timelineData.length === 0) {
    return (
      <div className="flex justify-center items-center py-20">
        <div className="text-neutral-400" role="status">
          {t(loading ? 'common.loading' : 'common.noData')}
        </div>
      </div>
    );
  }

  return (
    <Card variant="default" padding="md" animate>
      {/* 队伍选择器 */}
      <div className="flex flex-wrap gap-2 mb-6">
        {timelineData.map((team, index) => (
          <Button
            key={team.id}
            variant={!hiddenTeams.has(team.id) ? 'primary' : 'ghost'}
            aria-pressed={!hiddenTeams.has(team.id)}
            size="sm"
            onClick={() => toggleTeam(team.id)}
            className="!text-xs"
            style={{
              borderColor: !hiddenTeams.has(team.id) ? timelineTeamColor(index) : undefined,
              color: !hiddenTeams.has(team.id) ? timelineTeamColor(index) : undefined,
            }}
          >
            #{team.rank} {team.name}
          </Button>
        ))}
      </div>

      {/* 图表 */}
      <div className="h-96">
        <Suspense
          fallback={
            <div className="h-full flex items-center justify-center text-neutral-400">{t('common.loading')}</div>
          }
        >
          <ReactECharts option={chartOption} replaceMerge={['series']} style={{ height: '100%', width: '100%' }} />
        </Suspense>
      </div>
    </Card>
  );
}

export default ScoreboardTimeline;
