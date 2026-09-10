export const PAGE_SIZE = 20;

export const TASK_STATUSES = {
  history: ['', 'success', 'failed'],
  live: ['', 'active', 'pending', 'scheduled', 'retry', 'archived', 'completed'],
};

export function createTaskQuery() {
  return { page: 1, filters: { task_id: '', queue: '', status: '' } };
}

export function taskQueryReducer(query, action) {
  switch (action.type) {
    case 'page':
      return { ...query, page: action.page };
    case 'filter':
      return { page: 1, filters: { ...query.filters, [action.key]: action.value } };
    case 'reset':
      return createTaskQuery();
    default:
      return query;
  }
}

export function taskQueryParams({ page, filters }) {
  return { ...filters, limit: PAGE_SIZE, offset: (page - 1) * PAGE_SIZE };
}

export function normalizeTaskPayload(payload) {
  return {
    rows: Array.isArray(payload?.tasks) ? payload.tasks : [],
    totalCount: payload?.count || 0,
    queues: Array.isArray(payload?.queues) ? payload.queues : [],
  };
}

export function taskRowKey(item, mode) {
  return JSON.stringify([mode === 'history' ? item.id : item.queue, item.task_id]);
}

export function taskTimestamp(item, mode) {
  return mode === 'history' ? item.processed_at : item.next_process_at || item.completed_at || item.last_failed_at;
}

export function livePollingDelay(active, intervalSeconds) {
  return active && intervalSeconds > 0 ? intervalSeconds * 1000 : null;
}

export function formatDateTime(value, language) {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString(language || 'en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

export function formatPayload(value) {
  if (value === null || value === undefined || value === '') return '-';
  if (typeof value === 'string') return value;
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}
