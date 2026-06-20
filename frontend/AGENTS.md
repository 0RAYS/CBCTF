# AGENTS.md

Frontend is the embedded React app for CBCTF. Trust executable config over stale `frontend/README.md`.

## Commands

- Use pnpm 11.8.0: `pnpm install`, `pnpm dev`, `pnpm lint`, `pnpm build`, `pnpm preview`.
- `pnpm lint` runs `eslint . --fix`; expect it to rewrite files.
- `pnpm build` outputs `dist/`, which is embedded by `dist.go` into the Go binary.

## App Wiring

- Vite serves the app at `base: '/platform/'`; backend redirects `/` to `/platform` and serves embedded static files there.
- Routing uses `HashRouter`; route definitions live in `src/routes`.
- Redux stores live in `src/store`; API wrappers live in `src/api`.
- API base URL is hardcoded in `src/api/config.js`; there is no Vite proxy config.

## Code Rules

- Current stack is React 19, Vite 8, Tailwind CSS 4, React Router 7, Redux Toolkit 2, i18next 26, ECharts 6, Monaco, and motion 12.
- ESLint disables hook rules (`rules-of-hooks`, `exhaustive-deps`, `set-state-in-effect`); do not rely on lint to catch invalid hook usage.
- Prettier is enforced as ESLint warnings: single quotes, semicolons, trailing commas, width 120, `endOfLine: auto`.
- Keep UI bilingual for `en` and `zh-CN`; check CJK text length and wrapping.
- If repo-root `.impeccable.md` exists, follow its frontend design direction: dark theme, Maple UI/Maple Mono, geek blue `#597ef7`, restrained motion, no matrix/glitch/cyan-purple-gradient aesthetics.
