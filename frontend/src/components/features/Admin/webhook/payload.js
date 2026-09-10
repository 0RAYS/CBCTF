export function webhookForm(webhook) {
  return webhook
    ? {
        name: webhook.name,
        url: webhook.url,
        method: webhook.method,
        headers: { ...webhook.headers },
        timeout: webhook.timeout,
        retry: webhook.retry,
        events: [...(webhook.events || [])],
        on: webhook.on,
      }
    : { name: '', url: '', method: 'POST', headers: {}, timeout: 0, retry: 0, events: [], on: false };
}

export function cleanHeaders(headers) {
  return Object.fromEntries(
    Object.entries(headers)
      .filter(([key, value]) => key.trim() && value && value.trim())
      .map(([key, value]) => [key.trim(), value.trim()])
  );
}

export function renameHeader(headers, oldKey, newKey) {
  const key = newKey.trim();
  if (!key || key === oldKey || Object.hasOwn(headers, key) || !Object.hasOwn(headers, oldKey)) return headers;
  return Object.fromEntries(Object.entries(headers).map(([name, value]) => [name === oldKey ? key : name, value]));
}

export function buildWebhookPayload(form, original) {
  const payload = { ...form, headers: cleanHeaders(form.headers) };
  if (!original) return payload;
  const changes = {};
  for (const key of ['name', 'url', 'method', 'timeout', 'retry', 'on']) {
    if (payload[key] !== original[key]) changes[key] = payload[key];
  }
  for (const key of ['headers', 'events']) {
    if (JSON.stringify(payload[key]) !== JSON.stringify(original[key])) changes[key] = payload[key];
  }
  return changes;
}
