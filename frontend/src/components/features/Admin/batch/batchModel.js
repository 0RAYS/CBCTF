const batchStates = new Set(['success', 'partial', 'failed', 'cancelled']);

// These are the two current API contracts: direct batch data, or the contest
// import envelope's explicitly selected `batch` field.
export function getBatchResult(response, nested = false) {
  const batch = nested ? response?.data?.batch : response?.data;
  return batch && batchStates.has(batch.status) && Array.isArray(batch.items) ? batch : null;
}

export function hasBatchProgress(batch) {
  return (
    !!batch && (batch.succeeded > 0 || batch.skipped > 0 || batch.items.some((item) => hasBatchProgress(item.details)))
  );
}

export function batchTone(batch) {
  if (batch.status === 'success') return 'success';
  return batch.status === 'partial' || batch.status === 'cancelled' || hasBatchProgress(batch) ? 'warning' : 'danger';
}

export function remainingBatchIds(ids, batch) {
  const completed = new Set(
    batch.items.filter((item) => ['success', 'skipped'].includes(item.status)).map((item) => item.id)
  );
  return ids.filter((id) => !completed.has(String(id)));
}

export function remainingGeneratorCounts(submitted, batch) {
  const completed = batch.items.filter((item) => ['success', 'skipped'].includes(item.status));
  const counts = {};
  submitted.forEach((id, index) => {
    const prefix = `challenge:${id}/instance:${index + 1}`;
    if (!completed.some((item) => item.id === prefix || item.id.startsWith(`${prefix}/generator:`))) {
      counts[id] = (counts[id] || 0) + 1;
    }
  });
  return counts;
}

export function batchRows(batch, parents = []) {
  return batch.items.flatMap((item) => {
    const path = [...parents, item.id];
    const row = { ...item, path: path.join(' / ') };
    return item.details ? [row, ...batchRows(item.details, path)] : [row];
  });
}
