// The address list shows "name: protocol://ip:port"; only the ip:port part is
// worth pasting into a client.
export function copyableAddress(target) {
  const value = typeof target === 'string' ? target.trim() : '';
  if (!value) return '';
  const address = value.split(/\s+/).pop();
  return address.replace(/^[a-z][a-z0-9+.-]*:\/\//i, '');
}
