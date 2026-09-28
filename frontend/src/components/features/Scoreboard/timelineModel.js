function rankPosition(team) {
  const rank = Number(team.rank);
  return Number.isInteger(rank) && rank > 0 ? rank : Infinity;
}

export function buildTimelineModel(timelineData) {
  const allTimePoints = new Set();
  const teams = (timelineData || []).map((team) => {
    const points = (team.timeline || [])
      .toSorted((a, b) => (a.time < b.time ? -1 : a.time > b.time ? 1 : 0))
      // Stable sorting preserves the first event at each duplicate timestamp.
      .filter((point, index, points) => index === 0 || point.time !== points[index - 1].time);
    points.forEach((point) => allTimePoints.add(point.time));
    return { ...team, timeline: points };
  });
  // Stable numeric order also handles string ranks; unknown ranks follow ranked teams.
  return {
    timePoints: Array.from(allTimePoints).toSorted(),
    teams: teams.toSorted((a, b) => rankPosition(a) - rankPosition(b)),
  };
}

export function buildTimelineChartData({ teams, timePoints }, hiddenTeams) {
  return teams.flatMap((team, index) => {
    if (hiddenTeams.has(team.id)) return [];
    let cursor = 0;
    let score = 0;
    // Keep the shared axis and original color index; allocate only visible series.
    const data = timePoints.map((time) => {
      while (cursor < team.timeline.length && team.timeline[cursor].time <= time) {
        score = team.timeline[cursor++].score;
      }
      return score;
    });
    return [{ team, index, data }];
  });
}
