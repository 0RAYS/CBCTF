import { Suspense, lazy, useMemo, useState } from 'react';
import Button from '../../../common/Button';
import Card from '../../../common/Card';
import { useTranslation } from 'react-i18next';

const ReactECharts = lazy(() => import('../../../common/EChart'));

/**
 * 分数时间线图表组件
 * @param {Object} props
 * @param {Array} props.timelineData - 时间线数据
 */
function ScoreboardTimeline({ timelineData = [], loading = false }) {
  const { t, i18n } = useTranslation();
  const [hiddenTeams, setHiddenTeams] = useState(() => new Set());

  // 处理时间线数据, 转换为图表格式
  const { timePoints, teams } = useMemo(() => {
    const allTimePoints = new Set();
    const teams = (timelineData || []).map((team) => {
      const points = (team.timeline || [])
        .toSorted((a, b) => (a.time < b.time ? -1 : a.time > b.time ? 1 : 0))
        // Keep the first event at a duplicate timestamp, as the previous lookup did.
        .filter((point, index, points) => index === 0 || point.time !== points[index - 1].time);
      points.forEach((point) => allTimePoints.add(point.time));
      return { ...team, timeline: points };
    });
    return { timePoints: Array.from(allTimePoints).toSorted(), teams };
  }, [timelineData]);

  const chartData = useMemo(
    () =>
      teams.flatMap((team, index) => {
        if (hiddenTeams.has(team.id)) return [];
        let cursor = 0;
        let score = 0;
        // Preserve the shared category axis, but allocate values only for visible teams.
        const data = timePoints.map((time) => {
          while (cursor < team.timeline.length && team.timeline[cursor].time <= time) {
            score = team.timeline[cursor++].score;
          }
          return score;
        });
        return [{ team, index, data }];
      }),
    [teams, timePoints, hiddenTeams]
  );

  // 格式化时间显示
  const formatTime = (time) => {
    const date = new Date(time);
    return date.toLocaleTimeString(i18n.language || 'en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  // 格式化分数显示
  const formatScore = (value) => {
    return value.toLocaleString();
  };

  // 生成随机颜色
  const generateColor = (index) => {
    const colors = [
      '#597ef7', // blue
      '#EF4444', // red
      '#10B981', // green
      '#F59E0B', // yellow
      '#8B5CF6', // purple
      '#F97316', // orange
      '#06B6D4', // cyan
      '#84CC16', // lime
      '#EC4899', // pink
      '#FF64F7', // indigo
    ];
    return colors[index % colors.length];
  };

  // 生成 ECharts 配置
  const getChartOption = () => {
    const series = [];

    chartData.forEach(({ team, index, data }) => {
      series.push({
        id: `team_${team.id}`,
        name: `#${team.rank} ${team.name}`,
        type: 'line',
        data,
        smooth: true,
        lineStyle: {
          color: generateColor(index),
          width: 2,
        },
        itemStyle: {
          color: generateColor(index),
        },
        symbol: 'circle',
        symbolSize: 4,
        emphasis: {
          itemStyle: {
            color: generateColor(index),
            borderColor: generateColor(index),
            borderWidth: 2,
            symbolSize: 6,
          },
        },
      });
    });

    return {
      grid: {
        left: '3%',
        right: '4%',
        bottom: '8%',
        top: '5%',
        containLabel: true,
      },
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(0, 0, 0, 0.9)',
        borderColor: 'rgba(255, 255, 255, 0.3)',
        borderWidth: 1,
        borderRadius: 8,
        textStyle: {
          color: '#fff',
          fontSize: 12,
        },
        valueFormatter: formatScore,
      },
      xAxis: {
        type: 'category',
        data: timePoints.map(formatTime),
        axisLine: {
          lineStyle: {
            color: '#374151',
          },
        },
        axisLabel: {
          color: '#9CA3AF',
          fontSize: 12,
        },
      },
      yAxis: {
        type: 'value',
        axisLine: {
          lineStyle: {
            color: '#374151',
          },
        },
        axisLabel: {
          color: '#9CA3AF',
          fontSize: 12,
          formatter: formatScore,
        },
        splitLine: {
          lineStyle: {
            color: '#374151',
            type: 'dashed',
          },
        },
      },
      series: series,
    };
  };

  // 图表配置
  const chartOption = getChartOption();

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
              borderColor: !hiddenTeams.has(team.id) ? generateColor(index) : undefined,
              color: !hiddenTeams.has(team.id) ? generateColor(index) : undefined,
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
