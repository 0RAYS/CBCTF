import assert from 'node:assert/strict';
import test from 'node:test';
import { mapTeamSettings, teamUpdatePayload, unchangedTeamFields } from '../src/components/features/CTFGame/Team/model.js';

test('maps the captain separately and preserves member IDs without mutating API data', () => {
  const info = Object.freeze({ name: 'Team', picture: '/team.png', description: 'Bio', captain_id: 2 });
  const members = Object.freeze([
    Object.freeze({ id: 1, name: 'Member', picture: '/member.png', email: 'member@example.com' }),
    Object.freeze({ id: 2, name: 'Captain', picture: '/captain.png', email: 'captain@example.com' }),
  ]);
  const result = mapTeamSettings(info, members, 'invite-code');
  assert.deepEqual(result, {
    name: 'Team', picture: '/team.png', description: 'Bio', captainId: 2, inviteCode: 'invite-code',
    leader: { name: 'Captain', picture: '/captain.png', email: 'captain@example.com' },
    members: [{ id: 1, name: 'Member', picture: '/member.png', email: 'member@example.com' }],
  });
  assert.equal(result instanceof Promise, false);
  assert.notEqual(result.members[0], members[0]);
});

test('a missing captain does not crash mapping or remove other members', () => {
  const result = mapTeamSettings({ captain_id: 5 }, [{ id: 1, name: 'Member' }], '');
  assert.deepEqual(result.leader, { name: '', picture: '', email: '' });
  assert.equal(result.members.length, 1);
  assert.equal(result.inviteCode, '');
});

test('an empty member list remains valid', () => {
  assert.deepEqual(mapTeamSettings({ captain_id: 5 }, [], undefined).members, []);
});

test('saving only sends name and description, not the unimplemented leadership selection', () => {
  const description = 'x'.repeat(120);
  assert.deepEqual(teamUpdatePayload({ name: '', description, newLeader: 'Member', captainId: 2, inviteCode: 'secret' }), {
    name: '', description,
  });
});

test('refresh applies untouched fields but preserves edits made during the request', () => {
  const data = Object.freeze({ name: 'Old name', description: 'Old bio', inviteCode: 'old-code', picture: '/new.png', members: [] });
  assert.deepEqual(unchangedTeamFields(data, {}, { name: 1, description: 1, inviteCode: 1 }), {
    picture: '/new.png', members: [],
  });
});

test('refresh accepts fields whose versions have not changed, including earlier local edits', () => {
  assert.deepEqual(unchangedTeamFields({ name: 'Server name', inviteCode: 'new-code' }, { name: 2 }, { name: 2 }), {
    name: 'Server name', inviteCode: 'new-code',
  });
});
