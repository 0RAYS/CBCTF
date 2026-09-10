import test from 'node:test';
import assert from 'node:assert/strict';
import {
  normalizeInstanceStatus,
  isInstanceTransitioning,
  normalizeCategories,
  mapChallengeStatusToViewModel,
  getContestOverview,
  formatTimeLeft,
} from '../src/components/features/CTFGame/Challenges/models/challengeViewModel.js';

test('instance status normalization and transition classification', () => {
  for (const status of ['waiting', 'pending', 'terminating', 'running']) {
    assert.equal(normalizeInstanceStatus(status.toUpperCase()), status);
    assert.equal(isInstanceTransitioning(status), status !== 'running');
  }
  for (const status of [null, undefined, 42, {}, 'stopped', 'unknown']) {
    assert.equal(normalizeInstanceStatus(status), '');
  }
  assert.equal(isInstanceTransitioning(''), false);
});

test('category normalization handles absent categories without mutating input', () => {
  assert.deepEqual(normalizeCategories(null), []);
  assert.deepEqual(normalizeCategories({}), []);
  const categories = Object.freeze(['web', '', null, 'pwn']);
  assert.deepEqual(normalizeCategories(categories), ['web', 'pwn']);
});

test('list model maps API fields and preserves challenge metadata', () => {
  const challenge = Object.freeze({
    id: 3,
    name: 'Example',
    type: 'pods',
    category: 'web',
    score: 200,
    init: true,
    solved: true,
    solvers: 2,
    attempt: 5,
    attempts: 1,
    hints: ['hint'],
    tags: ['tag'],
    hidden: false,
    description: 'description',
  });
  const model = mapChallengeStatusToViewModel(challenge);
  assert.equal(model.title, 'Example');
  assert.equal(model.hasInstance, true);
  assert.equal(model.hasAttachments, false);
  assert.equal(model.isInitialized, true);
  assert.equal(model.isSolved, true);
  assert.equal(model.solves, 2);
  assert.equal(model.maxAttempts, 5);
  assert.deepEqual(model.attachments, []);
  assert.deepEqual(model.options, []);
  assert.deepEqual(model.hints, ['hint']);
  assert.deepEqual(model.tags, ['tag']);
  assert.equal(model.hidden, false);
  assert.equal(model.description, 'description');
  assert.equal(challenge.title, undefined);
  assert.equal(mapChallengeStatusToViewModel({ type: 'dynamic' }).hasAttachments, true);
});

test('status response false, zero and empty values override previous state', () => {
  const challenge = Object.freeze({
    id: 1,
    title: 'Title',
    name: 'Name',
    isInitialized: true,
    solved: true,
    isSolved: true,
    attempts: 4,
    attachment: 'old.zip',
    hasInstance: false,
    hasAttachments: false,
  });
  const status = Object.freeze({
    init: false,
    solved: false,
    attempts: 0,
    file: '',
  });
  const model = mapChallengeStatusToViewModel(challenge, status);
  assert.equal(model.isInitialized, false);
  assert.equal(model.isSolved, false);
  assert.equal(model.solved, false);
  assert.equal(model.attempts, 0);
  assert.equal(model.attachment, '');
  assert.equal(model.title, 'Title');
  assert.equal(model.hasInstance, false);
  assert.equal(model.hasAttachments, false);
});

test('incorrect flag status refresh updates attempts without dropping description or attachment', () => {
  const challenge = {
    id: 1,
    description: 'Keep me',
    attachment: 'challenge.zip',
    attempts: 1,
    isSolved: false,
  };
  const model = mapChallengeStatusToViewModel(challenge, {
    attempts: 2,
    solved: false,
  });
  assert.equal(model.attempts, 2);
  assert.equal(model.description, 'Keep me');
  assert.equal(model.attachment, 'challenge.zip');
  assert.equal(model.isSolved, false);
});

test('running duration uses server duration or a non-shrinking remaining-time fallback', () => {
  const initial = mapChallengeStatusToViewModel({
    remote: { status: 'RUNNING', remaining: '120' },
  });
  assert.equal(initial.instanceRunning, true);
  assert.equal(initial.instanceTimeLeft, 120);
  assert.equal(initial.instanceDuration, 120);
  const refreshed = mapChallengeStatusToViewModel(initial, {
    remote: { status: 'running', remaining: 100 },
  });
  assert.equal(refreshed.instanceDuration, 120);
  const extended = mapChallengeStatusToViewModel(refreshed, {
    remote: { status: 'running', remaining: 200 },
  });
  assert.equal(extended.instanceDuration, 200);
  const server = mapChallengeStatusToViewModel(extended, {
    remote: {
      status: 'running',
      remaining: '100',
      duration: '300',
      target: ['127.0.0.1'],
    },
  });
  assert.equal(server.instanceDuration, 300);
  assert.deepEqual(server.instanceIP, ['127.0.0.1']);
});

test('transition and stopped models clear incompatible status booleans', () => {
  const challenge = { instanceDuration: 300, instanceRunning: true };
  for (const status of ['waiting', 'pending', 'terminating', 'stopped']) {
    const model = mapChallengeStatusToViewModel(challenge, {
      remote: { status },
    });
    assert.equal(model.instanceRunning, false);
    assert.equal(model.instanceWaiting, status === 'waiting');
    assert.equal(model.instancePending, status === 'pending');
    assert.equal(model.instanceTerminating, status === 'terminating');
    assert.equal(model.instanceTimeLeft, 0);
    assert.equal(model.instanceDuration, 300);
  }
});

test('contest overview preserves time boundaries, seconds conversion and team totals', () => {
  const start = Date.parse('2026-09-10T00:00:00Z');
  const contest = {
    start: new Date(start).toISOString(),
    duration: 3600,
    prefix: 'flag',
    teams: 10,
  };
  const team = {
    score: 42,
    solved: [
      { solved: 1, all: 3 },
      { solved: 2, all: 4 },
    ],
  };
  const overview = getContestOverview(contest, team, start);
  assert.equal(getContestOverview(contest, team, start - 1).status, 'upcoming');
  assert.equal(overview.status, 'running');
  assert.equal(getContestOverview(contest, team, start + 3600000).status, 'running');
  assert.equal(getContestOverview(contest, team, start + 3600001).status, 'ended');
  assert.equal(overview.endTime, '2026-09-10T01:00:00.000Z');
  assert.equal(overview.duration, 1);
  assert.equal(overview.prefix, 'flag');
  assert.equal(overview.teams, 10);
  assert.deepEqual(overview.team, { score: 42, rank: 0, solved: 3 });
  assert.equal(overview.totalChallenges, 7);
});

test('countdown formatting retains hours and floors fractional seconds', () => {
  assert.equal(formatTimeLeft(0), '00:00:00');
  assert.equal(formatTimeLeft(3661.9), '01:01:01');
  assert.equal(formatTimeLeft(360000), '100:00:00');
});
