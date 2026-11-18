# Contest Challenges

- Picker selections span pages within a mounted session; reopening starts with clean filters and selection.
- Challenge metadata is a draft. Flag updates are independent immediate API transactions, not rolled back on Cancel.
- Save Changes waits for dirty flags before saving metadata. Never refresh every flag draft after saving one row.
- Contest challenges are independent copies of bank challenges. Never write a contest edit back to the bank challenge.
- Flag template changes keep existing solves; only "Reset Solves" clears them, and it always confirms first.
- Use common Modal for every dialog. Guard async responses against closed or replaced sessions.
- Pure query, selection and payload helpers belong in `challengeData.js` and `tests/contest-management.test.js`.
