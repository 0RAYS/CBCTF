
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
          // 中性灰系
          'neutral': {
            50: '#ffffff',
            100: '#f5f5f5',
            200: '#d9d9d9',
            300: '#b3b3b3',
            400: '#808080',
            500: '#4d4d4d',
            600: '#2a2a2a',
            700: '#1a1a1a',
            800: '#0d0d0d',
            900: '#000000'
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
          }
        },
        boxShadow: {
          // Design token: focus state glow (geek-400 accent) — inputs & interactive elements
          'focus': '0 0 15px rgba(89,126,247,0.3)',
          'focus-strong': '0 0 20px rgba(89,126,247,0.4)',
          // Design token: error state glow (red-400) — input error state
          'error': '0 0 15px rgba(239,68,68,0.3)',
          // Design tokens: status glow — toast notifications & status indicators (0.2 opacity)
          'glow-primary': '0 0 15px rgba(89,126,247,0.2)',
          'glow-success': '0 0 15px rgba(74,222,128,0.2)',
          'glow-warning': '0 0 15px rgba(251,191,36,0.2)',
          'glow-danger':  '0 0 15px rgba(239,68,68,0.2)',
          'glow-info':    '0 0 15px rgba(34,211,238,0.2)',
          'glow-muted':   '0 0 15px rgba(115,115,115,0.1)',
        },
      },
    },
    plugins: [],
  }
