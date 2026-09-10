# Admin Network UI

- `IpLookup` is the default export of `IpLookup.jsx`, with `{ ip, className }` props. Empty IPs display `-`; other IPs use a native keyboard-accessible button.
- `useIpLookup(scope)` owns the single IP API request path and returns `lookup(ip)`, `close()`, and `dialogProps` for `IpLookupDialog`. Pass a stable scope identity (the IP for a row, the log level for Logs). Closing, changing scope, unmounting, or starting another lookup invalidates pending results and errors. Manual log refresh also closes the lookup.
- `IpLookupDialog` owns the shared detail UI and bilingual labels. No page should duplicate its API call or fields.
- Logs decorates escaped ANSI text via `injectClickableIps`, then the unchanged `AnsiLog` applies DOMPurify. Only Logs adds `data-ip`, `role`, and `tabindex` to its attribute allowlist. Mouse, Enter, and Space activation are delegated and restricted to public IPv4 triggers inside the log container. Log pagination responses are scoped to the current level/refresh.

## Auto Refresh Control

`common/AutoRefreshControl.jsx` is controlled UI, not a polling hook:

```jsx
<AutoRefreshControl value={refreshInterval} onChange={setRefreshInterval} />
```

- `value` is a number in **seconds**: `5`, `10`, `30`, `60`, or `0` (off).
- `onChange(value: number)` receives seconds, not a DOM event or milliseconds.
- Consumers retain their own timers, visibility rules, in-flight guards, and refresh behavior. The task, generator, and victim consumers currently all store seconds and pass their setters directly.
- A consumer storing milliseconds must adapt at its boundary: `value={intervalMs / 1000}` and `onChange={(seconds) => setIntervalMs(seconds * 1000)}`. The control never converts or starts timers itself.
