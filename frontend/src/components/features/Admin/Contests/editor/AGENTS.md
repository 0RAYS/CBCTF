# Contest Editor

- `ContestEditor.jsx` owns the form draft; Identity, Schedule, Rules, Prizes and Timeline edit only their fields.
- Cover uploads persist immediately. Never replace the draft after uploading a cover.
- Keep API/draft/date conversions pure in `contestForm.js`; test with `node --test tests/contest-management.test.js` from frontend.
