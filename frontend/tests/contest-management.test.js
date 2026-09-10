import test from "node:test";
import assert from "node:assert/strict";
import {
  createContestDraft,
  contestUpdatePayload,
  formatDateForInput,
  validateContestDraft,
} from "../src/components/features/Admin/Contests/editor/contestForm.js";
import {
  challengeDraft,
  challengeQuery,
  flagUpdatePayload,
  isFlagDirty,
  toggleChallengeSelection,
} from "../src/components/features/Admin/Contests/challenges/challengeData.js";
import {
  containerStatus,
  filterTeamFlags,
  flagFilterOptions,
  formatFileSize,
  TEAM_DETAIL_PAGE_SIZE,
} from "../src/components/features/Admin/Contests/teams/teamDetailData.js";

const contest = {
  name: "Contest",
  description: "Description",
  picture: "/cover.png",
  start: "2026-09-10T10:00:00Z",
  duration: 3600,
  users: 12,
  rules: ["Rule"],
  prizes: [{ amount: "$100", description: "Prize", id: 3 }],
  timelines: [
    {
      date: "2026-09-10T10:30:00Z",
      title: "Event",
      description: "Details",
      id: 7,
    },
  ],
  prefix: "CTF",
  size: 4,
  hidden: true,
  captcha: "invite",
  blood: false,
  victims: 2,
};

test("contest conversion preserves settings and isolates mutable draft arrays", () => {
  const draft = createContestDraft(contest);
  assert.equal(draft.endTime, "2026-09-10T11:00:00.000Z");
  assert.equal(draft.participants, 12);
  assert.equal(draft.image, "/cover.png");
  assert.equal(draft.hidden, true);
  assert.equal(draft.blood, false);
  draft.rules.push("Draft rule");
  draft.prizes[0].amount = "$200";
  draft.timeline[0].title = "Draft event";
  assert.deepEqual(contest.rules, ["Rule"]);
  assert.equal(contest.prizes[0].amount, "$100");
  assert.equal(contest.timelines[0].title, "Event");
});

test("contest payload excludes immediate image upload and display-only fields", () => {
  const draft = createContestDraft(contest);
  draft.image = "/uploaded-immediately.png";
  draft.endTime = "2026-09-10T11:00:00.999Z";
  assert.deepEqual(contestUpdatePayload(draft), {
    name: "Contest",
    description: "Description",
    prefix: "CTF",
    size: 4,
    hidden: true,
    captcha: "invite",
    blood: false,
    victims: 2,
    start: "2026-09-10T10:00:00.000Z",
    duration: 3600,
    rules: ["Rule"],
    prizes: [{ amount: "$100", description: "Prize" }],
    timelines: [
      {
        date: "2026-09-10T10:30:00.000Z",
        title: "Event",
        description: "Details",
      },
    ],
  });
});

test("empty timeline dates and default arrays remain supported", () => {
  const draft = createContestDraft({
    ...contest,
    rules: null,
    prizes: null,
    timelines: null,
  });
  assert.deepEqual(draft.rules, []);
  assert.deepEqual(draft.prizes, []);
  draft.timeline.push({ date: "", title: "Undated", description: "" });
  assert.equal(contestUpdatePayload(draft).timelines[0].date, "");
  assert.equal(createContestDraft().blood, true);
  assert.notEqual(createContestDraft().rules, createContestDraft().rules);
});

test("contest validation rejects missing, invalid and reversed schedule values", () => {
  const draft = createContestDraft(contest);
  assert.deepEqual(validateContestDraft(draft), {});
  assert.equal(
    validateContestDraft({ ...draft, title: "  " }).title,
    "titleRequired",
  );
  assert.equal(
    validateContestDraft({ ...draft, startTime: "" }).startTime,
    "invalidDate",
  );
  assert.equal(
    validateContestDraft({ ...draft, endTime: "invalid" }).endTime,
    "invalidDate",
  );
  assert.equal(
    validateContestDraft({ ...draft, endTime: draft.startTime }).endTime,
    "endTimeAfterStart",
  );
  assert.equal(
    validateContestDraft({ ...draft, endTime: "2020-01-01" }).endTime,
    "endTimeAfterStart",
  );
});

test("date input formatting uses local calendar fields and handles invalid dates", () => {
  const date = new Date(2026, 8, 10, 9, 5);
  assert.equal(formatDateForInput(date.toISOString()), "2026-09-10T09:05");
  assert.equal(formatDateForInput("bad date"), "");
  assert.equal(formatDateForInput(""), "");
});

test("challenge queries retain filters across pages and omit all/blank filters", () => {
  assert.deepEqual(challengeQuery({}), { limit: 10, offset: 0 });
  assert.deepEqual(
    challengeQuery({
      page: 3,
      pageSize: 20,
      type: "pods",
      category: "web",
      name: " name ",
      description: " desc ",
    }),
    {
      limit: 20,
      offset: 40,
      type: "pods",
      category: "web",
      name: "name",
      description: "desc",
    },
  );
  assert.deepEqual(challengeQuery({ name: " ", description: "\n" }), {
    limit: 10,
    offset: 0,
  });
});

test("picker selection persists across pages and toggles by identity", () => {
  const first = { id: 1, name: "First" };
  const second = { id: 2, name: "Second" };
  const selected = [first];
  const both = toggleChallengeSelection(selected, second);
  assert.deepEqual(both, [first, second]);
  assert.deepEqual(selected, [first]);
  assert.deepEqual(toggleChallengeSelection(both, { id: 1, name: "Renamed" }), [
    second,
  ]);
});

test("challenge draft clones hints/tags without carrying pool or flag fields", () => {
  const original = {
    name: "Challenge",
    hidden: true,
    tags: ["web"],
    hints: ["hint"],
    flags: [{ id: 1 }],
  };
  const draft = challengeDraft(original);
  draft.tags.push("draft");
  draft.hints[0] = "draft";
  assert.deepEqual(original.tags, ["web"]);
  assert.deepEqual(original.hints, ["hint"]);
  assert.equal("flags" in draft, false);
  assert.equal(draft.attempt, 0);
  assert.deepEqual(challengeDraft({}).tags, []);
});

test("flag payload contains only writable fields and preserves zeroes", () => {
  const flag = {
    id: 3,
    value: "",
    score: 0,
    score_type: 0,
    decay: 0,
    min_score: 0,
    current_score: 200,
    solvers: 2,
    blood: [1, 2, 3],
  };
  assert.deepEqual(flagUpdatePayload(flag), {
    value: "",
    score: 0,
    score_type: 0,
    decay: 0,
    min_score: 0,
  });
  assert.equal(
    isFlagDirty(flag, { ...flag, current_score: 100, solvers: 8 }),
    false,
  );
  for (const [field, value] of Object.entries({
    value: "flag",
    score: 100,
    score_type: 1,
    decay: 20,
    min_score: 50,
  })) {
    assert.equal(isFlagDirty({ ...flag, [field]: value }, flag), true, field);
  }
});

const flags = [
  {
    id: 1,
    name: "Mixed WEB",
    type: "pods",
    category: "web",
    flags: [{ solved: false }, { solved: true }],
  },
  {
    id: 2,
    name: "Crypto",
    type: "static",
    category: "crypto",
    flags: [{ solved: false }],
  },
  { id: 3, name: "No flags", flags: null },
];

test("team flags combine case-insensitive names, type/category and any-solved semantics", () => {
  assert.deepEqual(
    filterTeamFlags(flags, {
      name: "web",
      type: "pods",
      category: "web",
      solved: "true",
    }),
    [flags[0]],
  );
  assert.deepEqual(filterTeamFlags(flags, { solved: "false" }), [
    flags[1],
    flags[2],
  ]);
  assert.deepEqual(filterTeamFlags(flags, { name: "missing" }), []);
  assert.deepEqual(
    filterTeamFlags(flags, { type: "static", category: "web" }),
    [],
  );
  assert.deepEqual(flagFilterOptions([...flags, flags[0]]), {
    types: ["pods", "static"],
    categories: ["crypto", "web"],
  });
});

test("traffic status keeps inclusive running boundaries and team page size", () => {
  const start = "2026-09-10T10:00:00Z";
  const now = new Date(start).getTime();
  assert.equal(containerStatus(start, 60, now - 1), "upcoming");
  assert.equal(containerStatus(start, 60, now), "running");
  assert.equal(containerStatus(start, 60, now + 60000), "running");
  assert.equal(containerStatus(start, 60, now + 60001), "ended");
  assert.equal(TEAM_DETAIL_PAGE_SIZE, 20);
  assert.equal(formatFileSize(0), "0 B");
  assert.equal(formatFileSize(1024), "1.00 KB");
  assert.equal(formatFileSize(1048576), "1.00 MB");
});
