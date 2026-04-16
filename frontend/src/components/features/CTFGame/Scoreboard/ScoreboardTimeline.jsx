import { useState, useEffect } from 'react';
import ReactECharts from 'echarts-for-react';
import { Button, Card } from '../../../../components/common';
import { useTranslation } from 'react-i18next';

/**
 * 分数时间线图表组件
 * @param {Object} props
 * @param {Array} props.timelineData - 时间线数据
 */
function ScoreboardTimeline({ timelineData = [] }) {
  const { t, i18n } = useTranslation();
  const [chartData, setChartData] = useState([]);
  const [selectedTeams, setSelectedTeams] = useState([]);
  const [chartKey, setChartKey] = useState(0);

  // 处理时间线数据, 转换为图表格式
  useEffect(() => {
    if (!timelineData || timelineData.length === 0) {
      setChartData([]);
      return;
    }

    // 获取所有时间点
    const allTimePoints = new Set();
    timelineData.forEach((team) => {
      team.timeline.forEach((point) => {
        allTimePoints.add(point.time);
      });
    });

    // 按时间排序
    const sortedTimePoints = Array.from(allTimePoints).sort();

    // 构建图表数据
    const processedData = sortedTimePoints.map((time) => {
      const dataPoint = { time };

      timelineData.forEach((team) => {
        // 找到该时间点或之前最近的分数
        const teamScore = team.timeline
          .filter((point) => point.time <= time)
          .sort((a, b) => new Date(b.time) - new Date(a.time))[0];

        // 如果找到分数, 使用该分数; 否则使用0（从0开始）
        dataPoint[`team_${team.id}`] = teamScore ? teamScore.score : 0;
      });

      return dataPoint;
    });

    setChartData(processedData);

    setSelectedTeams(timelineData.map((team) => team.id));
  }, [timelineData]);

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
      '#3B82F6', // blue
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
    if (!chartData.length) return {};

    const timePoints = chartData.map((item) => formatTime(item.time));
    const series = [];

    timelineData.forEach((team, index) => {
      if (!selectedTeams.includes(team.id)) return;

      series.push({
        name: `#${team.rank} ${team.name}`,
        type: 'line',
        data: chartData.map((item) => item[`team_${team.id}`]),
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
        formatter: function (params) {
          let result = `<div style="margin-bottom: 8px; font-family: 'Maple Mono', 'Source Han Sans SC', ui-monospace, monospace; color: #9CA3AF;">${params[0].axisValue}</div>`;
          params.forEach((param) => {
            // 显示所有队伍, 包括0分
            const team = timelineData.find((t) => `#${t.rank} ${t.name}` === param.seriesName);
            if (team) {
              result += `<div style="color: ${param.color}; margin: 2px 0;">${param.seriesName}: ${formatScore(param.value)}</div>`;
            }
          });
          return result;
        },
      },
      xAxis: {
        type: 'category',
        data: timePoints,
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
    setSelectedTeams((prev) => (prev.includes(teamId) ? prev.filter((id) => id !== teamId) : [...prev, teamId]));
    // 强制图表重新渲染
    setChartKey((prev) => prev + 1);
  };

  if (!timelineData || timelineData.length === 0) {
    return (
      <div className="flex justify-center items-center py-20">
        <div className="text-neutral-400">{t('common.noData')}</div>
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
            variant={selectedTeams.includes(team.id) ? 'primary' : 'ghost'}
            size="sm"
            onClick={() => toggleTeam(team.id)}
            className="!text-xs"
            style={{
              borderColor: selectedTeams.includes(team.id) ? generateColor(index) : undefined,
              color: selectedTeams.includes(team.id) ? generateColor(index) : undefined,
            }}
          >
            #{team.rank} {team.name}
          </Button>
        ))}
      </div>

      {/* 图表 */}
      <div className="h-96">
        <ReactECharts key={chartKey} option={chartOption} style={{ height: '100%', width: '100%' }} />
      </div>
    </Card>
  );
}

export default ScoreboardTimeline;
