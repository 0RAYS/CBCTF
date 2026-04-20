
/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
      extend: {
        fontFamily: {
          sans: [
            '"Maple UI"',
            '"Source Han Sans SC"',
            'ui-sans-serif', 'system-ui', 'sans-serif',
            '"Apple Color Emoji"', '"Segoe UI Emoji"', '"Segoe UI Symbol"', '"Noto Color Emoji"',
          ],
          mono: [
            '"Maple Mono"',
            '"Source Han Sans SC"',
            'ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas',
            '"Liberation Mono"', '"Courier New"', 'monospace',
          ],
        },
        colors: {
          // 中性灰系 — 纯灰，不偏色，保持黑色基调
          'neutral': {
            50:  '#f9f9f9',
            100: '#f0f0f0',
            200: '#d9d9d9',
            300: '#b3b3b3', // 次要文字 — contrast 7.5:1 on neutral-900 (WCAG AA ✓)
            400: '#8a8a8a', // 三级文字 — contrast 4.7:1 on neutral-900 (WCAG AA ✓)
            500: '#4d4d4d',
            600: '#2a2a2a',
            700: '#1a1a1a',
            800: '#0f0f0f',
            900: '#0a0a0a', // 页面底色 — 极深黑，比纯黑 #000 稍有层次
          },
          geek: {
            50: '#f0f5ff',
            100: '#d6e4ff',
            200: '#adc6ff',
            300: '#85a5ff',
            400: '#597ef7',
            500: '#2f54eb',
            600: '#1d39c4',
            700: '#10239e',
            800: '#061178',
            900: '#030852',
          },
          // Semantic rank tokens (scoreboard gold/silver/bronze)
          'rank': {
            'gold':   '#c9a227',
            'silver': '#8a9bb0',
            'bronze': '#a06535',
          },
        },
        boxShadow: {
          // Design token: focus state glow (geek-400 accent) — inputs & interactive elements
          'focus':        '0 0 0 2px oklch(62% 0.22 265 / 0.35), 0 0 16px oklch(62% 0.22 265 / 0.18)',
          'focus-strong': '0 0 0 2px oklch(62% 0.22 265 / 0.5),  0 0 20px oklch(62% 0.22 265 / 0.28)',
          // Design token: error state glow (red-400) — input error state
          'error':        '0 0 0 2px oklch(65% 0.21 25 / 0.35),  0 0 16px oklch(65% 0.21 25 / 0.18)',
          // Design tokens: status glow — toast notifications & status indicators
          'glow-primary': '0 0 18px oklch(62% 0.22 265 / 0.22)',
          'glow-success': '0 0 18px oklch(72% 0.18 155 / 0.22)',
          'glow-warning': '0 0 18px oklch(80% 0.17 75  / 0.22)',
          'glow-danger':  '0 0 18px oklch(65% 0.21 25  / 0.22)',
          'glow-info':    '0 0 18px oklch(75% 0.14 210 / 0.22)',
          'glow-muted':   '0 0 18px oklch(55% 0.04 265 / 0.12)',
        },
      },
    },
    plugins: [],
  }
