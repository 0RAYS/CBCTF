# Scoreboard

- Shared ranking, table, statistics, timeline and solve-progress presentation belongs here, not under `CTFGame` or `Admin/Contests`.
- `AdminRanking.jsx` adapts admin translation keys, pagination presentation and row clicks to `ScoreboardRanking`. Both pages use `ScoreboardTable` directly; do not restore a table wrapper or pagination-component injection.
- `scoreboardModel.js`, `timelineModel.js` and `timelineOption.js` are pure JavaScript with no React, API or ECharts imports. Run their tests from `frontend` with `node --test tests/scoreboard.test.js`.
- Table challenge collection and per-team indexes keep the first matching challenge ID, including unsolved duplicates. Build indexes once per teams update, not with a search per cell. Hover styling must not rerender the entire table.
- Timeline normalization stably sorts teams by numeric rank and keeps the first event at duplicate timestamps. The selector and chart share this model. Use a forward cursor over each visible team's points, retaining all teams' shared axis and unfiltered rank-order palette indices when hiding teams.
- Preserve `axisLabel.alignMinLabel`, `alignMaxLabel`, `hideOverlap`, lazy loading of `common/EChart` and `replaceMerge={['series']}`. Do not add a runtime ECharts dependency to the option builder.
- The user and admin scoreboard pages own their distinct API calls, pagination and request cleanup. User table data caches its fetched page; admin view switches refetch and reset pagination. The timeline is always above both views, loads independently on contest entry and supports retry; list switches/pages must not cancel or repeat its request. Timeline `null` means unfetched while `[]` is a cached empty response.
- Keep the contest-keyed page boundary and stale-request guards. Admin team details use `features/Admin/details/useTeamDetailDialog.jsx`; do not move detail fetching into shared presentation.
