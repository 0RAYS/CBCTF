// formatGap 把秒数格式化成便于阅读的时间间隔。
export function formatGap(seconds) {
  if (!seconds || seconds <= 0) return '-';
  if (seconds < 60) return `${seconds.toFixed(1)}s`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ${Math.round(seconds % 60)}s`;
  const hours = Math.floor(minutes / 60);
  return `${hours}h ${minutes % 60}m`;
}

// formatDateTime 统一线索证据里的时间展示格式。
export function formatDateTime(value, language = 'en-US') {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleString(language, { hour12: false });
}
