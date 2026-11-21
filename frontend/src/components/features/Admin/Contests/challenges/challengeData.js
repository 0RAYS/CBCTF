export function challengeQuery({
  page = 1,
  pageSize = 10,
  type = 'all',
  category = 'all',
  name = '',
  description = '',
  contest = 'all',
}) {
  return {
    limit: pageSize,
    offset: (page - 1) * pageSize,
    ...(type !== 'all' ? { type } : {}),
    ...(category !== 'all' ? { category } : {}),
    ...(name.trim() ? { name: name.trim() } : {}),
    ...(description.trim() ? { description: description.trim() } : {}),
    ...(contest !== 'all' ? { contest_id: Number(contest) } : {}),
  };
}

export function toggleChallengeSelection(selected, challenge) {
  return selected.some((item) => item.id === challenge.id)
    ? selected.filter((item) => item.id !== challenge.id)
    : [...selected, challenge];
}
