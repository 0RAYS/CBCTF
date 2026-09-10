# Admin Details

- User and team detail dialogs remain separate domain-specific components, not a configurable universal dialog.
- Import hooks and dialogs directly from concrete modules here, never from the global hooks barrel.
- `useTeamDetailDialog` opens a keyed `TeamDetailSession`; close, reopen and contest changes invalidate prior requests.
- Keep traffic permission checks at fetch, render and download boundaries. Use `Admin/traffic/TrafficGraphDialog`.
