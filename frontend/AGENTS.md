# AGENTS.md

Frontend is the embedded React app for CBCTF. Trust executable config over stale `frontend/README.md`.

## Commands

- Use pnpm 11.8.0: `pnpm install`, `pnpm dev`, `pnpm lint`, `pnpm build`, `pnpm preview`.
- `pnpm lint` runs `eslint . --fix`; expect it to rewrite files.
- `pnpm build` outputs `dist/`, which is embedded by `dist.go` into the Go binary.
- `pnpm test` runs the committed Node test suite (models, payloads, sessions, polling, and architecture). No extra test dependency is required.
- `pnpm lint:check` checks without editing and enables `react-hooks/rules-of-hooks`; prefer this for verification.

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

## Structure

- See `README.md` for the domain map and component contracts. Pages adapt routes/scopes; business components, hooks, loaders and pure models live together in `components/features`.
- Shared scoreboard presentation lives in `features/Scoreboard`. Admin detail hooks, IP lookup, traffic, images, generators and victims live in their Admin domains, not `common`, global `hooks`, or contest-only directories.
- Keep `common` business-independent, never import pages from components, and avoid synchronous dependency cycles. Import heavy Markdown/EChart adapters directly, preserving lazy chart/editor boundaries.
- Use `Modal` for the shell, `ConfirmModal` for confirmation actions, `Checkbox` for checkbox fields, and `MarkdownContent` for GFM. Do not restore deleted compatibility aliases or old-path re-export wrappers.
- Tests enforce a 500-line limit for each source JS/JSX file and exact-case import resolution. Split by responsibility rather than moving a large page into a large hook. New pure modules must be added to the reviewed architecture-test dependency set.
- Preserve independent request scope, query and selection generations. Polling must not supersede an in-flight request for the same query, and cleanup must invalidate late responses.
- Use the real toast service contract in tests: `toast.danger({ description })`, not a fabricated `toast.error` method or a string argument.
