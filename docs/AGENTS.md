# AGENTS.md

Docs are an Rspress 2 site. Prefer executable config and `docs/README.md` over stale assumptions.

## Commands

- Use pnpm, matching `pnpm-lock.yaml` and CI: `pnpm install`, `pnpm dev`, `pnpm build`, `pnpm preview`.
- Use `DOCS_BASE=/CBCTF/ pnpm build` for GitHub Project Pages-style subpaths.
- `pnpm inspect` writes the final builder config when Rspress/Rsbuild behavior is unclear.

## Content

- Source pages live under repo-root `docs/docs` (`docs/` from this directory); build output is `doc_build/` and must not be committed.
- Every page needs `title` and `description` frontmatter for SEO, local search, and generated `llms.txt` output.
- Navigation/sidebar order comes from `_nav.json` and `_meta.json`, not filename prefixes.
- `rspress.config.ts` strips MDX `style` attributes during rendering; prefer classes/global CSS for styling.

## Docs

- Rspress: https://rspress.rs/llms.txt
- Rsbuild: https://rsbuild.rs/llms.txt
- Rspack: https://rspack.rs/llms.txt
