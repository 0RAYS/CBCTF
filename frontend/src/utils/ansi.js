// 简单 ANSI 转 HTML（保留常见颜色）, 不引入外部依赖
// 仅处理 SGR（Select Graphic Rendition）码, 如 \u001b[31m, \u001b[0m 等

const COLOR_MAP = {
  // standard (brightened for dark background)
  30: '#9ca3af', // black -> gray-400
  31: '#ff6b6b', // red
  32: '#22c55e', // green
  33: '#facc15', // yellow
  34: '#60a5fa', // blue
  35: '#f472b6', // magenta
  36: '#22d3ee', // cyan
  37: '#ffffff', // white
  // bright
  90: '#cbd5e1', // bright black -> slate-300
  91: '#ff8585',
  92: '#4ade80',
  93: '#fde047',
  94: '#93c5fd',
  95: '#f9a8d4',
  96: '#67e8f9',
  97: '#ffffff',
};

const BGCOLOR_MAP = {
  40: 'black',
  41: 'red',
  42: 'green',
  43: 'yellow',
  44: 'blue',
  45: 'magenta',
  46: 'cyan',
  47: 'white',
  100: 'gray',
  101: 'red',
  102: 'green',
  103: 'yellow',
  104: 'blue',
  105: 'magenta',
  106: 'cyan',
  107: 'white',
};

const escapeHtml = (str) =>
  str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;');

export function ansiToHtml(input) {
  if (!input) return '';
  const text = escapeHtml(input);

  let openTags = [];
  const closeAll = () => (openTags.length ? '</span>'.repeat(openTags.length) : '');

  // 通过字符串拼接构造匹配 ESC 的正则, 避免源码中直接出现控制字符
  const ESC = String.fromCharCode(27);
  const sgrPattern = new RegExp(ESC + '\\[[0-9;]*m', 'g');
  const ctrlPattern = new RegExp(ESC + '\\[[0-9;]*[A-Za-z]', 'g');

  const converted =
    text
      // 将 ESC[...m 序列转换为 span
      .replace(sgrPattern, (match) => {
        const codes = match
          .slice(2) // remove ESC[
          .slice(0, -1) // remove trailing m
          .split(';')
          .filter(Boolean)
          .map((n) => parseInt(n, 10));

        // 重置
        if (codes.length === 0 || codes.includes(0)) {
          const closing = closeAll();
          openTags = [];
          return closing;
        }

        let style = '';
        codes.forEach((code) => {
          if (code === 1) style += 'font-weight:700;'; // bold
          if (code === 2) style += 'opacity:0.8;'; // faint
          if (code === 3) style += 'font-style:italic;';
          if (code === 4) style += 'text-decoration:underline;';
          if (code >= 30 && code <= 37) style += `color:${COLOR_MAP[code]};`;
          if (code >= 90 && code <= 97) style += `color:${COLOR_MAP[code]};`;
          if (code >= 40 && code <= 47) style += `background-color:${BGCOLOR_MAP[code]};`;
          if (code >= 100 && code <= 107) style += `background-color:${BGCOLOR_MAP[code]};`;
        });

        if (!style) return '';
        openTags.push('span');
        return `<span style="${style}">`;
      })
      // 移除其余不可见控制符
      .replace(ctrlPattern, '') + closeAll();

  // 用白色作为默认前景色, 子 span 的 color 会覆盖
  return `<span style="color:#ffffff">${converted}</span>`;
}

export default ansiToHtml;
