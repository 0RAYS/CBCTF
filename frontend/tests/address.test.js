import assert from 'node:assert/strict';
import test from 'node:test';
import { copyableAddress } from '../src/utils/address.js';

test('the address list keeps the service name and protocol out of the clipboard', () => {
  assert.equal(copyableAddress('web: tcp://127.0.0.1:32206'), '127.0.0.1:32206');
  assert.equal(copyableAddress('tcp://127.0.0.1:32206'), '127.0.0.1:32206');
  assert.equal(copyableAddress('web: 127.0.0.1:32206'), '127.0.0.1:32206');
  assert.equal(copyableAddress('127.0.0.1:32206'), '127.0.0.1:32206');
});

test('ipv6 addresses keep their brackets', () => {
  assert.equal(copyableAddress('web: tcp://[fd00::1]:32206'), '[fd00::1]:32206');
});

test('an empty or unusable target copies nothing', () => {
  assert.equal(copyableAddress(''), '');
  assert.equal(copyableAddress('   '), '');
  assert.equal(copyableAddress(undefined), '');
});
