const TEAM_COLORS = [
  '#597ef7',
  '#EF4444',
  '#10B981',
  '#F59E0B',
  '#8B5CF6',
  '#F97316',
  '#06B6D4',
  '#84CC16',
  '#EC4899',
  '#FF64F7',
];

export function timelineTeamColor(index) {
  return TEAM_COLORS[index % TEAM_COLORS.length];
}

export function buildTimelineOption(timePoints, chartData, locale = 'en-US') {
  const formatScore = (value) => value.toLocaleString();
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
      data: timePoints.map((time) =>
        new Date(time).toLocaleTimeString(locale, {
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
        })
      ),
      axisLine: { lineStyle: { color: '#374151' } },
      axisLabel: {
        color: '#9CA3AF',
        fontSize: 12,
        alignMinLabel: 'left',
        alignMaxLabel: 'right',
        hideOverlap: true,
      },
    },
    yAxis: {
      type: 'value',
      axisLine: { lineStyle: { color: '#374151' } },
      axisLabel: {
        color: '#9CA3AF',
        fontSize: 12,
        formatter: formatScore,
      },
      splitLine: { lineStyle: { color: '#374151', type: 'dashed' } },
    },
    series: chartData.map(({ team, index, data }) => ({
      id: `team_${team.id}`,
      name: `#${team.rank} ${team.name}`,
      type: 'line',
      data,
      smooth: true,
      lineStyle: { color: timelineTeamColor(index), width: 2 },
      itemStyle: { color: timelineTeamColor(index) },
      symbol: 'circle',
      symbolSize: 4,
      emphasis: {
        itemStyle: {
          color: timelineTeamColor(index),
          borderColor: timelineTeamColor(index),
          borderWidth: 2,
          symbolSize: 6,
        },
      },
    })),
  };
}
