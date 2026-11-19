// Returns the "ip:port" part of an address entry.
export function copyableAddress(target) {
  const value = typeof target === 'string' ? target.trim() : '';
  if (!value) return '';
  const address = value.split(/\s+/).pop();
  return address.replace(/^[a-z][a-z0-9+.-]*:\/\//i, '');
}
