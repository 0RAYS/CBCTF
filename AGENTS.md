# AGENTS.md

CBCTF is a Kubernetes-native CTF platform, not a challenge repo: Go/Gin backend plus React frontend embedded into the Go binary.

## Commands

- Frontend build must run before local Go release builds: `cd frontend && pnpm install && pnpm build` creates `frontend/dist/`, which `frontend/dist.go` embeds.
- Backend release build requires Go 1.26.x plus CGO/libpcap: `CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o CBCTF .`.
- Docker performs both build stages with Node 24, Go 1.26, `build-base`, and `libpcap-dev`: `docker build -t cbctf .`.
- Frontend uses pnpm 11.8.0: `cd frontend && pnpm dev`, `pnpm lint` (`eslint . --fix`), `pnpm build`.
- Docs use pnpm, not npm: `cd docs && pnpm install && pnpm build`; use `DOCS_BASE=/CBCTF/ pnpm build` for GitHub Project Pages-style subpaths.
- The Helm chart depends on the upstream `valkey` chart (`https://valkey.io/valkey-helm/`). `chart/Chart.lock` is committed but `chart/charts/` is not, so run `helm dependency build chart/` before `helm template`, `helm install`, or `go test ./chart/...`; those tests skip with that instruction when the dependency is absent. The chart release workflow runs `helm dependency update` itself.
- No committed test suite is present, and root `.gitignore` ignores `*_test.go`; adding Go tests requires changing ignore rules. Use focused `go test ./...`, `pnpm lint`, and builds only when relevant.

## Backend

- Entry flow is `main.go` -> `cmd.Cmd()` -> `preInit()` -> `run()` in `cmd/run.go`.
- Startup is intentionally two-phase: config/log/db from file first, then DB-backed system settings reload log/db before Redis, Kubernetes clients, email, webhook, Asynq tasks, cron, and Gin router start.
- Routes are registered centrally in `internal/router/router.go`; per-domain router files call service packages, which use GORM models in `internal/model`.
- DB init in `internal/db/db.go` auto-runs `AutoMigrate`, sets join tables, ensures `pg_trgm`, then seeds settings, branding, permissions, roles, groups, cron jobs, admin, and default OAuth providers.
- Asynq workers live in `internal/task`; queue names/concurrency come from `asynq.queues` config.
- Kubernetes code assumes either in-cluster config or `k8s.config` kubeconfig, then initializes Kubernetes, Multus NetworkAttachmentDefinition, Kube-OVN, and KubeVirt clients.

## Config And Secrets

- Default config is embedded from `internal/config/default.yaml`; `./CBCTF -c <path>` selects a config file.
- If the config file is missing, startup writes defaults with mode `0600` and exits so the operator can edit and restart.
- Environment overrides use Viper prefix `CBCTF_` with dots mapped to underscores, e.g. `CBCTF_GIN_PORT`.
- Do not commit local runtime material: root `config.yaml`, `admin.yaml`, `docker-compose.yaml`, `build.sh`, `data/`, `logs/`, `uploads/`, and `frontend/dist/*` are gitignored.
- Root `AGENTS.md` applies to the whole repository; subdirectory `AGENTS.md` files apply only inside their directories.

## Frontend

- Trust `frontend/package.json`, not stale `frontend/README.md`: current stack is React 19, Vite 8, Tailwind CSS 4, React Router 7, Redux Toolkit 2, i18next 26, ECharts 6, Monaco, and motion 12.
- Vite serves the platform under `base: '/platform/'`; the Go router redirects `/` to `/platform` and serves embedded static files at `/platform`.
- App routing uses `HashRouter`; route definitions are in `frontend/src/routes`, with Redux stores in `frontend/src/store` and API wrappers in `frontend/src/api`.
- API wrappers use `frontend/src/api/config.js`; there is no Vite proxy config, and `API_CONFIG.BASE_URL` is currently hardcoded.
- ESLint runs with `--fix`; Prettier is enforced as warnings with single quotes, semicolons, trailing commas, width 120, and `endOfLine: auto`.
- Hooks rules are disabled in the current ESLint config; do not assume lint catches invalid hook usage.
- UI must remain bilingual: frontend i18n supports `en` and `zh-CN`; CJK text length and wrapping are first-class constraints.

## Comments

- Comments carry only what the code cannot: intent, non-obvious constraints, and pitfalls. Everything else is noise.
- Never narrate the edit history: no "used to", "no longer", "now we", "kept for compatibility" remarks about how the code changed. State the current behavior instead.
- Keep comments short. One or two lines at the point of use beats a paragraph block; split long explanations into bullets only when the content is genuinely a list of rules.
- No agent-style instruction blocks: do not address the reader imperatively ("you must", "ensure that", "never do X"), and do not explain steps the code already shows.
- No large pasted payloads: do not inline JSON request/response samples or API schemas in comments. Keep a single descriptive line per exported function.
- No commented-out code. Delete it; version control keeps the history.
- No boilerplate doc comments on interface implementations (e.g. `// SourceID 实现 GeneratorSource`). Document the interface once, then let the compiler speak.
- Keep the long English explanations of non-obvious mechanics (Kubernetes capability requirements, kubelet image-list truncation, cache-key derivation, traffic-analysis invariants). Those are the comments worth reading.
- Frontend API wrappers (`frontend/src/api`) get one short line per function; do not duplicate response shapes there.

## Functions And Structure

- Do not wrap a single expression in a function that only restates it; inline the expression at the call site.
- Keep the call chain short: the behavior should be reachable in one or two hops, not through a chain of one-line helpers.
- A helper earns its place when it names a domain concept, removes duplication, or isolates logic worth testing.
- Do not add a `OldFunctionWithNewValue` variant when only the returned data changed. If the old function is fully replaced by the new one, change the old function in place instead of introducing a second name.
- Chained builder methods (`WithX`) are the exception: they exist to configure a value, so keep them.

## Docs

- Docs are an Rspress 2 site in `docs/`; source pages are under `docs/docs`, build output is `docs/doc_build/` and should not be committed.
- Every docs page needs `title` and `description` frontmatter for SEO, local search, and generated `llms.txt` output.
- Navigation/sidebar order is controlled by `_nav.json` and `_meta.json`, not filename prefixes.
- `docs/AGENTS.md` applies inside `docs/`; keep it in sync with root docs guidance if present.

## CTF/Kubernetes Domain Rules

- Dynamic container modes are inferred from compose YAML: no `networks` field means shared Pod network; `networks` means Kube-OVN VPC mode with static IPs; `x-kubevirt: true` enables VM-backed services.
- Compose parsing uses `compose-go` with known extensions from `internal/model/spec.go` (`x-volumes`, `x-kubevirt`, `x-boot`, `x-cloudinit`); preserve those names in challenge YAML examples.
- Dynamic attachment generators must provide `sleep` and `unzip`, run `/root/run.sh <team_id> <base64_encoded_flags>`, and write `/root/mnt/attachments/{id}.zip`.
- Flag raw forms are `static{...}`, `leet{...}`, and `uuid{}`; leet mutation is implemented in `internal/utils/flag.go`.
- Platform runtime depends on PostgreSQL, Valkey, Kubernetes, Kube-OVN, Multus, and optionally KubeVirt/FRP; the Helm chart in `chart/` deploys app, PostgreSQL 17, and Valkey 9 through the upstream `valkey` chart at `https://valkey.io/valkey-helm/`.

## CI And Design

- GitHub Actions publish `ghcr.io/0rays/cbctf:latest` on `main` pushes that touch backend/frontend build paths; docs and chart releases are separate `main` workflows scoped to `docs/**` and `chart/**`.
- `chart/Chart.yaml` is the chart version source of truth; update it with chart changes.
