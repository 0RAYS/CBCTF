# Scheduling review ? 2026-09-21

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
