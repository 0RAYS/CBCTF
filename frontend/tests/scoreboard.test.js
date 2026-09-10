import assert from 'node:assert/strict';
import test from 'node:test';
import {
  buildScoreboardColumns,
  collectChallenges,
  indexTeamChallenges,
  toRankingTeam,
} from '../src/components/features/Scoreboard/scoreboardModel.js';
import { buildTimelineChartData, buildTimelineModel } from '../src/components/features/Scoreboard/timelineModel.js';
import { buildTimelineOption, timelineTeamColor } from '../src/components/features/Scoreboard/timelineOption.js';

const times = [0, 1, 2, 3].map((minute) => `2026-09-10T10:0${minute}:00Z`);

test('ranking keeps page offsets, solve totals and locale date formatting', () => {
  const solved = Object.freeze([
    { category: 'WEB', solved: 2 },
    { category: 'PWN', solved: 3 },
  ]);
  const team = Object.freeze({ id: 7, name: 'Team', picture: '/avatar', score: 1250, solved, last: times[0] });
  for (const locale of ['en-US', 'zh-CN']) {
    const ranking = toRankingTeam(team, 2, 3, 20, locale);
    assert.deepEqual(ranking, {
      id: 7,
      rank: 43,
      name: 'Team',
      picture: '/avatar',
      score: 1250,
      solved,
      totalSolved: 5,
      lastSubmit: new Date(times[0])
        .toLocaleString(locale, {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
          hour12: false,
        })
        .replace(/\//g, '-'),
    });
  }
});

test('ranking defaults missing solves and last submission without changing zero scores', () => {
  const ranking = toRankingTeam({ id: 1, score: 0 });
  assert.equal(ranking.rank, 1);
  assert.equal(ranking.score, 0);
  assert.deepEqual(ranking.solved, []);
  assert.equal(ranking.totalSolved, 0);
  assert.equal(ranking.lastSubmit, '-');
});

test('challenge collection retains the first matching id before category/name sorting', () => {
  const first = Object.freeze({ id: 1, category: 'WEB', name: 'Zulu', solved: 0 });
  const second = Object.freeze({ id: 2, category: 'PWN', name: 'Alpha' });
  const third = Object.freeze({ id: 3, category: 'WEB', name: 'Alpha' });
  const challenges = Object.freeze([first, second]);
  assert.deepEqual(
    collectChallenges([
      { challenges },
      {},
      { challenges: [{ id: 1, category: 'AAA', name: 'Duplicate', solved: 1 }, third] },
    ]),
    [second, third, first]
  );
  assert.deepEqual(collectChallenges(), []);
});

test('table index matches the first challenge, not the first solved duplicate', () => {
  const challenges = Object.freeze([
    Object.freeze({ id: 1, solved: 0 }),
    Object.freeze({ id: 1, solved: 3 }),
    Object.freeze({ id: '1', solved: 2 }),
    Object.freeze({ id: 2, solved: -1 }),
    Object.freeze({ id: 3, solved: 1 }),
  ]);
  const index = indexTeamChallenges(challenges);
  for (const id of [1, '1', 2, 3, 4]) {
    assert.equal(
      index.get(id),
      challenges.find((challenge) => challenge.id === id)
    );
  }
  assert.equal(index.get(1).solved > 0, false);
  assert.equal(index.get('1').solved > 0, true);
  assert.equal(index.get(4)?.solved > 0, false);
  assert.equal(indexTeamChallenges().size, 0);
});

test('table index is built once and cell lookups do not scan challenge ids', () => {
  let reads = 0;
  const challenges = Array.from({ length: 1000 }, (_, id) => ({
    get id() {
      reads++;
      return id;
    },
    solved: id % 2,
  }));
  const index = indexTeamChallenges(challenges);
  assert.ok(reads <= challenges.length * 2);
  const indexedReads = reads;
  for (let id = 0; id < challenges.length; id++) assert.equal(index.get(id).solved, id % 2);
  assert.equal(reads, indexedReads);
});

test('table columns preserve category order and clamp estimated widths', () => {
  const challenges = Object.freeze([
    { id: 1, category: 'WEB', name: '' },
    { id: 2, category: 'PWN', name: 'x'.repeat(50) },
    { id: 3, category: 'WEB', name: '12345678' },
    { id: 4, category: '__proto__', name: null },
  ]);
  const model = buildScoreboardColumns(challenges);
  assert.deepEqual(
    model.categoryEntries.map(([category]) => category),
    ['WEB', 'PWN', '__proto__']
  );
  assert.deepEqual(
    model.columns.map(({ id }) => id),
    [1, 3, 2, 4]
  );
  assert.deepEqual(model.widths, [80, 112, 200, 80]);
  assert.equal(model.widthById.get(3), 112);
  assert.deepEqual(buildScoreboardColumns().columns, []);
});

test('timeline sorting is immutable and duplicate timestamps retain their first event', () => {
  const first = Object.freeze({ time: times[2], score: 20 });
  const timeline = Object.freeze([first, { time: times[0], score: 0 }, { time: times[2], score: 99 }]);
  const input = Object.freeze([Object.freeze({ id: 1, timeline }), { id: 2 }]);
  const model = buildTimelineModel(input);
  assert.deepEqual(model.timePoints, [times[0], times[2]]);
  assert.deepEqual(model.teams[0].timeline, [timeline[1], first]);
  assert.deepEqual(model.teams[1].timeline, []);
  assert.equal(input[0].timeline[0], first);
  assert.notEqual(model.teams[0], input[0]);
});

test('timeline cursor matches the original first-match and carry-forward algorithm', () => {
  const input = Array.from({ length: 12 }, (_, id) => ({
    id,
    timeline: Array.from({ length: id + 3 }, (_, event) => ({
      time: times[(event * 7 + id) % times.length],
      score: event % 3 === 0 ? 0 : event * 10 - id,
    })).reverse(),
  }));
  input.push({ id: 12, timeline: [] });
  const model = buildTimelineModel(input);
  const chartData = buildTimelineChartData(model, new Set());
  input.forEach((team, index) => {
    let score = 0;
    const expected = model.timePoints.map((time) => {
      const point = team.timeline.find((point) => point.time === time);
      if (point) score = point.score;
      return score;
    });
    assert.deepEqual(chartData[index].data, expected);
  });
});

test('timeline starts at zero, carries scores across gaps and permits score decreases', () => {
  const model = buildTimelineModel([
    {
      id: 1,
      timeline: [
        { time: times[1], score: 20 },
        { time: times[3], score: 0 },
      ],
    },
    {
      id: 2,
      timeline: [
        { time: times[0], score: 5 },
        { time: times[2], score: 10 },
      ],
    },
  ]);
  const chartData = buildTimelineChartData(model, new Set());
  assert.deepEqual(
    chartData.map(({ data }) => data),
    [
      [0, 20, 20, 0],
      [5, 5, 10, 10],
    ]
  );
});

test('hidden teams retain the full axis and original color indices without reading their series', () => {
  const model = buildTimelineModel([
    { id: 10, timeline: [{ time: times[0], score: 10 }] },
    { id: 20, timeline: [{ time: times[1], score: 20 }] },
    { id: 30, timeline: [{ time: times[2], score: 30 }] },
  ]);
  const all = buildTimelineChartData(model, new Set());
  Object.defineProperty(model.teams[0], 'timeline', {
    get() {
      throw new Error('Hidden series must not be traversed');
    },
  });
  const visible = buildTimelineChartData(model, new Set([10, 20]));
  assert.deepEqual(visible, [all[2]]);
  assert.deepEqual(model.timePoints, times.slice(0, 3));
  assert.equal(visible[0].index, 2);
  assert.deepEqual(buildTimelineChartData(model, new Set([10, 20, 30])), []);
});

test('timeline cursor reads each visible timeline linearly rather than searching per timestamp', () => {
  let reads = 0;
  const timePoints = Array.from({ length: 1000 }, (_, index) => index);
  const timeline = timePoints.map((time) => ({
    get time() {
      reads++;
      return time;
    },
    score: time,
  }));
  const data = buildTimelineChartData({ timePoints, teams: [{ id: 1, timeline }] }, new Set());
  assert.deepEqual(data[0].data, timePoints);
  assert.ok(reads <= timePoints.length * 2);
});

test('empty and unfetched timelines produce empty models and options', () => {
  for (const input of [undefined, null, []]) {
    const model = buildTimelineModel(input);
    assert.deepEqual(model, { timePoints: [], teams: [] });
    const data = buildTimelineChartData(model, new Set());
    const option = buildTimelineOption(model.timePoints, data);
    assert.deepEqual(option.xAxis.data, []);
    assert.deepEqual(option.series, []);
  }
});

test('timeline options preserve series identity, palette and axis label boundaries', () => {
  const model = buildTimelineModel([
    { id: 10, rank: 1, name: 'First', timeline: [{ time: times[0], score: 1000 }] },
    { id: 20, rank: 2, name: 'Second', timeline: [{ time: times[1], score: 20 }] },
  ]);
  const allData = buildTimelineChartData(model, new Set());
  const all = buildTimelineOption(model.timePoints, allData);
  for (const locale of ['en-US', 'zh-CN']) {
    const visible = buildTimelineOption(model.timePoints, buildTimelineChartData(model, new Set([10])), locale);
    assert.deepEqual(visible.series, [all.series[1]]);
    assert.equal(visible.series[0].id, 'team_20');
    assert.equal(visible.series[0].name, '#2 Second');
    assert.equal(visible.series[0].lineStyle.color, timelineTeamColor(1));
    assert.deepEqual(
      visible.xAxis.data,
      model.timePoints.map((time) =>
        new Date(time).toLocaleTimeString(locale, {
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
        })
      )
    );
    assert.equal(visible.xAxis.axisLabel.alignMinLabel, 'left');
    assert.equal(visible.xAxis.axisLabel.alignMaxLabel, 'right');
    assert.equal(visible.xAxis.axisLabel.hideOverlap, true);
    assert.equal(visible.yAxis.axisLabel.formatter(1000), (1000).toLocaleString());
    assert.equal(visible.tooltip.valueFormatter(1000), (1000).toLocaleString());
  }
  assert.equal(all.series[0].data, allData[0].data);
  assert.equal(timelineTeamColor(0), '#597ef7');
  assert.equal(timelineTeamColor(10), timelineTeamColor(0));
});
