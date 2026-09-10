export function toRankingTeam(team, index = 0, page = 1, pageSize = 20, locale = 'en-US') {
  const solved = team.solved || [];
  return {
    id: team.id,
    rank: (page - 1) * pageSize + index + 1,
    name: team.name,
    picture: team.picture,
    score: team.score,
    solved,
    totalSolved: solved.reduce((total, category) => total + category.solved, 0),
    lastSubmit: team.last
      ? new Date(team.last)
          .toLocaleString(locale, {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false,
          })
          .replace(/\//g, '-')
      : '-',
  };
}

export function collectChallenges(teams = []) {
  const challenges = new Map();
  teams.forEach((team) => {
    team.challenges?.forEach((challenge) => {
      if (!challenges.has(challenge.id)) challenges.set(challenge.id, challenge);
    });
  });
  return Array.from(challenges.values()).toSorted((a, b) =>
    a.category !== b.category ? a.category.localeCompare(b.category) : a.name.localeCompare(b.name)
  );
}

export function buildScoreboardColumns(challenges = []) {
  const byCategory = challenges.reduce((categories, challenge) => {
    (categories[challenge.category] ??= []).push(challenge);
    return categories;
  }, Object.create(null));
  const categoryEntries = Object.entries(byCategory);
  const columns = categoryEntries.flatMap(([, items]) => items);
  const widths = columns.map((challenge) => {
    const nameLength = challenge?.name ? String(challenge.name).length : 0;
    return Math.max(80, Math.min(200, 48 + nameLength * 8));
  });
  return {
    categoryEntries,
    columns,
    widths,
    widthById: new Map(columns.map((challenge, index) => [challenge?.id, widths[index]])),
  };
}

export function indexTeamChallenges(challenges = []) {
  const byId = new Map();
  challenges.forEach((challenge) => {
    // Match Array.find's first result, including an unsolved first duplicate.
    if (!byId.has(challenge.id)) byId.set(challenge.id, challenge);
  });
  return byId;
}
