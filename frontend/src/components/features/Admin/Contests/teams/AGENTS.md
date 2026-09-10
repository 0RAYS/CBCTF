# Contest Teams

- `Teams.jsx` is the management view. Detail panels are separate Submissions, Writeups, Traffic and Flags components.
- Detail pagination uses 20 items. Flags filter locally and reset when their panel unmounts.
- Traffic access must be checked before fetching, rendering or downloading; dialogs and sessions live in `Admin/details`.
- Use `Admin/network/IpLookup` for IP lookup and `common/Checkbox` for boolean fields.
