# Scheduling review - 2026-09-21

## Scope and approach

Repository inventory and cross-layer review of startup/config, models/repositories,
routers/services, Asynq workers, Redis leases, cron reconciliation, client-go,
Kube-OVN/Multus/KubeVirt/FRP, admin victim/generator UI, chart and build workflows.
This is a targeted scheduling review, not a claim that every unrelated feature has
been exhaustively audited. Runtime secrets and generated/vendor assets are excluded.

## Implementation sequence

1. Fix generator pool key/lease consistency and contest isolation.
2. Fix Pod readiness/watch handling, parallel error propagation and cleanup order.
3. Harden task state transitions, cancellation and fail-closed reconciliation.
4. Add bounded, authorized workload diagnostics and bilingual admin presentation.
5. Scope Helm RBAC and improve startup/shutdown protection.
6. Reject unscoped deletion and harden Kubernetes client/image-pull job defaults.

Each item is committed separately with regression coverage where feasible.

## Confirmed findings

- Lua uses `generator:<id>`/`:locked`, but Go registers `generators:<id>` and
  `generators:locked:<id>`: pool discovery removes healthy members and locks diverge.
- Global generator pool rebuild includes contest-owned generators; unregistering
  also deletes a lease potentially owned by an executing attachment.
- CreatePod performs a redundant GET and deletes on any GET error, treats Running
  as Ready and hand-written watches fail on expiration/transient disconnects.
- Victim parallel creation shares `ret`; failures lacking an Error attribute are
  silently ignored, potentially dereferencing nil Kubernetes objects.
- Cleanup tears down network isolation before workloads stop and releases FRP
  ports before proxy Pods disappear. Generator deletion does not wait for absence.
- Generator startup ignores Asynq cancellation and lacks conditional transitions.
- ForceStopVictim rollback restores terminating, not the actual original state.
- Expiry checks include not-yet-started victims; orphan cleanup treats DB errors
  as proof of absence. Both can destroy otherwise valid workloads.
- VM creation currently reports API acceptance, not guest readiness.
- Admin Pod views omit readiness, container waiting/termination state and restarts;
  log reads lack a byte cap and use detached request contexts.
- Chart binds namespaced CRUD through a ClusterRoleBinding, granting it in all
  namespaces; slow initialization lacks a startup probe.

## Validation environment

Actual executable configs require Go 1.27.1 and pnpm 12.4.1 (older AGENTS prose is
stale). The selected Go 1.27.1 installation initially contains only bin tools;
`go test` fails with `go: no such tool "vet"`. Do not mistake skipped/incomplete
checks for successful tests. No live Kubernetes/Redis/PostgreSQL cluster is assumed.

## Follow-up / integration checks

- Exercise real Kube-OVN CNI cleanup, VM readiness and FRP handoff under load.
- Recover pending records after hard worker/process loss using a reconciler and
  durable operation ownership/outbox; do not blindly retry partially completed creates.
- Consider namespace-scoped shared informers once resource/label ownership and
  cache lifecycle are defined; avoid an unsafe broad cache refactor in a fix batch.
- Network isolation of generator workloads, configurable generator/capture budgets,
  image pull policy and node placement need operator-compatible deployment policy.
- Kubernetes event history and durable operation timelines need retention/redaction
  policy; never expose challenge environment/flags or full Pod specs to clients.

## Completed: pool and Kubernetes lifecycle

- Shared Lua/Go key prefixes, preserved owner leases and global/contest pool isolation.
- Removed destructive preflight GET/delete; Ready condition and terminal container
  checks, UID-pinned client-go relist/watch recovery, detailed timeout reasons.
- Parallel Kubernetes operations now use local results and the errgroup context.
- VM readiness is awaited; foreground VM deletion and workload disappearance precede
  removal of network isolation. FRP Pods are included and ports released last.
- Generator cleanup waits for Pod deletion and scopes services by generator ID.
- Downloaded the complete official Go 1.27.1 toolchain into ignored `.gocache/`;
  existing global tool installations were not modified. Focused Go tests now run.
- Go fake-client tests cover Ready vs Running, exits, UID replacement, disappearance,
  timeout diagnostics, expired watches and non-destructive creation/deletion.
- Redis integration test is opt-in (`CBCTF_TEST_REDIS_ADDR`); not run without Redis.

## Completed: task ownership and reconciliation

- Start/stop workers serialize by workload using PostgreSQL session advisory locks,
  released on process death. No long-running DB transaction or Redis lease expiry
  is involved. This requires direct/session-pooled PostgreSQL, **not transaction-mode
  PgBouncer**, and uses a separate lock-only connection pool so even a one-connection TaskDB cannot
  deadlock behind its own lifecycle locks. Budget the extra PostgreSQL connections.
- Generator transitions are conditional, duplicate starts skip non-waiting records,
  startup follows task cancellation, and failed enqueue has bounded cleanup fallback.
- Failed victim startup retains partial allocations; only the worker that claimed
  startup can schedule failure cleanup. Stop enqueue rollback restores the old state.
- Expiry only applies to running victims. Orphan scans abort before any deletion on
  DB errors, deduplicate IDs, and distinguish NotFound from temporary unavailability.
- Attachment tasks re-check generator state after taking the lease. Exec uses argv
  without a shell, preserves empty arguments, and does not destroy temporarily
  unready generators. Filesystem errors are no longer silently ignored.
- Optional PostgreSQL concurrency integration test uses `CBCTF_TEST_POSTGRES_DSN`.

- Deferred: FRP allocations need tokenized ownership/durable release across DB
  failure after cleanup; synchronous administrative hard-deletes and worker crash
  recovery also need unified ownership/outbox coverage.

## Completed: workload diagnostics and admin UI

- Pod responses expose Ready/deletion, scheduling condition reasons, container state,
  restart counts and current/previous exit codes. No specs, images, environment,
  commands, free-form status messages or capture-container details are returned.
- Generator status endpoints use the existing global/contest permission and ownership
  chains. Log access validates Pod ownership; victim logs now use a scoped GET rather
  than listing all sibling Pods. Legacy generator labels remain readable.
- Kubernetes requests inherit HTTP cancellation. Log requests cap at 10,000 lines and
  1 MiB (both server request options and defensive local reads).
- Bilingual admin dialogs refresh status every 5 seconds without overlapping requests,
  abort on close/scope change, preserve selection and mark stale data after failures.
  Logs remain explicit snapshots with a refresh button and bounded line controls.
- Go tests cover diagnostic redaction/ownership/log limits and route permission mapping.
  Frontend tests cover polling, late responses, selection and disappearance: 265 passed;
  lint:check and production build passed (existing large-chunk warning remains).

## Completed: deployment permissions and startup/shutdown

- Chart 0.0.24 splits namespaced CRUD into Role/RoleBinding. Cluster grants retain only
  required cluster resources; namespace GET is restricted to the release namespace.
- Fixed SelfSubjectAccessReview resource/subresource encoding for exec/log and added
  named namespace checks. Upgrade the backend and chart together.
- Configurable five-minute startup probe and 120-second termination grace. HTTP server
  initialization precedes serving, schedulers start synchronously, and HTTP draining
  is bounded to 30 seconds; Asynq shutdown allowance is 30 seconds. Running cron work
  is still awaited and can exceed the default grace (documented operator limitation).
- Helm lint passed. Go regression tests render default and custom-service-account
  charts, compare RBAC against backend permission checks, reject wildcard/namespace
  leakage and verify configurable startup/shutdown settings. No live install performed.

## Completed: destructive API and client/job hardening

- Every batch-delete helper requires a validated non-empty workload ownership selector;
  missing/empty/ambiguous filters, role-only filters and invalid IDs fail before API calls.
  Kube-OVN IP deletion additionally permits an explicit subnet owner selector.
- Generator cleanup validates ownership before deletion, including name collisions.
  Single-Pod cleanup and orphan-generator deletion pin the observed UID; a replacement
  Pod prevents cleanup completion rather than silently releasing dependent resources.
- Core/Multus/Kube-OVN/KubeVirt clients share the existing 100-QPS / 150-burst budget
  instead of independently multiplying bursts. Requests carry a CBCTF user agent.
- Image pre-pull Jobs do not mount a service-account token or inject service links,
  have a ten-minute active deadline and retain results for five minutes before TTL GC.
  Operators with exceptionally slow image pulls may need a different policy later.
- Focused tests passed for all ten batch-delete entrypoints, ownership/UID protection,
  image-job defaults and shared client configuration. These are not load benchmarks.
