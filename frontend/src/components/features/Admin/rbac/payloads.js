export function userUpdatePayload(form) {
  const data = { ...form };
  if (!data.password) delete data.password;
  return data;
}

export function successfulAssignmentIds(ids, results) {
  return ids.filter((_, index) => results[index]?.status === 'fulfilled' && results[index].value?.code === 200);
}

export function toggleVisibleCandidates(selectedIds, users) {
  const visibleIds = users.map((user) => user.id);
  if (visibleIds.length > 0 && visibleIds.every((id) => selectedIds.includes(id))) {
    return selectedIds.filter((id) => !visibleIds.includes(id));
  }
  return [...new Set([...selectedIds, ...visibleIds])];
}

export function lastAvailablePage(count, pageSize, page) {
  return Math.min(page, Math.max(1, Math.ceil(count / pageSize)));
}
