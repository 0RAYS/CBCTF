const IPV4_RE = /(?<![\w.])(?:\d{1,3}\.){3}\d{1,3}(?!\w|\.\d)/g;

export function isPublicIp(ip) {
  if (typeof ip !== 'string' || !/^(?:\d{1,3}\.){3}\d{1,3}$/.test(ip)) return false;
  const parts = ip.split('.').map(Number);
  if (parts.some((part) => part > 255)) return false;
  const [a, b, c] = parts;
  return !(
    a === 10 ||
    a === 127 ||
    (a === 172 && b >= 16 && b <= 31) ||
    (a === 192 && b === 168) ||
    (a === 169 && b === 254) ||
    (a === 100 && b >= 64 && b <= 127) ||
    a === 0 ||
    (a === 198 && (b === 18 || b === 19)) ||
    (a === 192 && b === 0 && (c === 0 || c === 2)) ||
    (a === 192 && b === 88 && c === 99) ||
    (a === 198 && b === 51 && c === 100) ||
    (a === 203 && b === 0 && c === 113) ||
    a >= 224
  );
}

// Input is escaped ANSI HTML. Only decorate text; AnsiLog sanitizes the result afterwards.
export function injectClickableIps(html) {
  return html.replace(/(<[^>]+>)|([^<]+)/g, (match, tag, text) => {
    if (tag || !text) return match;
    return text.replace(IPV4_RE, (ip) => {
      if (!isPublicIp(ip)) return ip;
      return `<span data-ip="${ip}" role="button" tabindex="0" class="ip-lookup-trigger focus-visible:outline-2 focus-visible:outline-geek-400" style="color:#597ef7;cursor:pointer;text-decoration:underline;text-decoration-style:dotted;">${ip}</span>`;
    });
  });
}
