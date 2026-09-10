import assert from 'node:assert/strict';
import test from 'node:test';
import { transformContestData } from '../src/components/features/CTFGame/OverView/contestModel.js';
import { DEFAULT_CONTEST_IMAGE } from '../src/config/contest.js';

const contest = Object.freeze({
  prefix: 'CTF',
  name: 'Example',
  description: 'Contest description',
  picture: '/contest.png',
  start: '2024-01-01T08:00:00+08:00',
  duration: 3600,
  users: 12,
  size: 3,
  teams: 4,
  notices: 2,
  blood: false,
  hidden: true,
});

test('a null or absent top-level response remains no contest', () => {
  assert.equal(transformContestData(null), null);
  assert.equal(transformContestData(undefined), null);
});

for (const [name, content] of [
  ['absent', {}],
  ['null', { rules: null, prizes: null, timelines: null }],
  ['empty', { rules: Object.freeze([]), prizes: Object.freeze([]), timelines: Object.freeze([]) }],
]) {
  test(`${name} contest collections remain empty without sample content`, () => {
    const input = Object.freeze({ ...contest, ...content });
    const before = structuredClone(input);
    const result = transformContestData(input);
    assert.deepEqual(result.rules, []);
    assert.deepEqual(result.prizes, []);
    assert.deepEqual(result.timeline, []);
    assert.deepEqual(input, before);
  });
}

test('populated content and metadata are preserved without mutating the response', () => {
  const rules = Object.freeze(['First rule', 'Second rule']);
  const prizes = Object.freeze([Object.freeze({ amount: '$0', description: 'Recognition' })]);
  const timelines = Object.freeze([
    Object.freeze({ date: '2024-01-01T00:00:00Z', title: 'Start', description: 'Challenges open' }),
  ]);
  const input = Object.freeze({ ...contest, rules, prizes, timelines });
  const before = structuredClone(input);

  assert.deepEqual(transformContestData(input), {
    title: 'CTF Example',
    description: 'Contest description',
    image: '/contest.png',
    status: 'ended',
    startTime: '2024-01-01T00:00:00.000Z',
    endTime: '2024-01-01T01:00:00.000Z',
    participants: 12,
    rules,
    prizes,
    timeline: timelines,
    teamSize: 3,
    teamsCount: 4,
    noticesCount: 2,
    isBlood: false,
    isHidden: true,
  });
  assert.deepEqual(input, before);
});

test('collection normalization is independent and retains image and participant fallbacks', () => {
  const result = transformContestData({ ...contest, rules: ['Keep this'], prizes: null, picture: '', users: 0 });
  assert.deepEqual(result.rules, ['Keep this']);
  assert.deepEqual(result.prizes, []);
  assert.deepEqual(result.timeline, []);
  assert.equal(result.image, DEFAULT_CONTEST_IMAGE);
  assert.equal(result.participants, 0);
  assert.equal(transformContestData({ ...contest, users: undefined }).participants, 0);
});
