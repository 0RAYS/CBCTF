import assert from 'node:assert/strict';
import test from 'node:test';
import {
  buildWebhookPayload,
  cleanHeaders,
  renameHeader,
  webhookForm,
} from '../src/components/features/Admin/webhook/payload.js';
import { buildSmtpPayload, smtpForm } from '../src/components/features/Admin/smtp/payload.js';
import { buildOAuthPayload, oauthForm, parseScopes } from '../src/components/features/Admin/oauth/payload.js';

test('webhook defaults and editing isolate mutable collections from server data', () => {
  assert.deepEqual(webhookForm(), {
    name: '',
    url: '',
    method: 'POST',
    headers: {},
    timeout: 0,
    retry: 0,
    events: [],
    on: false,
  });
  const original = { ...webhookForm(), headers: { Authorization: 'token' }, events: ['login'] };
  const form = webhookForm(original);
  form.headers.Authorization = 'new';
  form.events.push('logout');
  assert.deepEqual(original.headers, { Authorization: 'token' });
  assert.deepEqual(original.events, ['login']);
  assert.deepEqual(webhookForm({ ...original, headers: null, events: null }).headers, {});
});

test('webhook headers trim names and values, omit blanks, and keep last trimmed duplicate', () => {
  const headers = { ' X-Token ': ' token ', Empty: '', Blank: '  ', '  ': 'value', 'X-Token': 'last' };
  assert.deepEqual(cleanHeaders(headers), { 'X-Token': 'last' });
  assert.equal(headers[' X-Token '], ' token ');
  assert.deepEqual(buildWebhookPayload({ ...webhookForm(), headers }).headers, { 'X-Token': 'last' });
});

test('header rename preserves latest values and other rows without mutating the source', () => {
  const headers = { First: 'latest value', Second: 'second value' };
  const renamed = renameHeader(headers, 'First', ' Renamed ');
  assert.deepEqual(renamed, { Renamed: 'latest value', Second: 'second value' });
  assert.deepEqual(Object.keys(renamed), ['Renamed', 'Second']);
  assert.deepEqual(headers, { First: 'latest value', Second: 'second value' });
  assert.deepEqual(buildWebhookPayload({ ...webhookForm(), headers: renamed }).headers, renamed);
});

test('header rename rejects blanks, duplicates, unchanged keys and deleted rows', () => {
  const headers = { First: 'one', Second: 'two' };
  for (const next of ['', '  ', 'First', ' First ', 'Second', ' Second ']) {
    assert.equal(renameHeader(headers, 'First', next), headers);
  }
  assert.equal(renameHeader(headers, 'Removed', 'New'), headers);
});

test('headers named like object prototype properties remain ordinary payload data', () => {
  const headers = JSON.parse('{"__proto__":" value ","constructor":" ctor "}');
  const cleaned = cleanHeaders(headers);
  assert.equal(Object.getPrototypeOf(cleaned), Object.prototype);
  assert.equal(Object.hasOwn(cleaned, '__proto__'), true);
  assert.equal(cleaned.__proto__, 'value');
  assert.deepEqual(renameHeader({ A: 'one' }, 'A', '__proto__'), JSON.parse('{"__proto__":"one"}'));
});

test('webhook unchanged form sends an empty delta', () => {
  const original = { ...webhookForm(), id: 12, headers: { Token: 'value' }, events: ['login'] };
  assert.deepEqual(buildWebhookPayload(webhookForm(original), original), {});
});

test('webhook delta retains explicit empty, false and zero changes', () => {
  const original = {
    ...webhookForm(),
    name: 'old',
    url: 'https://old',
    method: 'GET',
    timeout: 30,
    retry: 3,
    on: true,
    headers: { Token: 'value' },
    events: ['login'],
  };
  assert.deepEqual(buildWebhookPayload(webhookForm(), original), {
    name: '',
    url: '',
    method: 'POST',
    timeout: 0,
    retry: 0,
    on: false,
    headers: {},
    events: [],
  });
});

test('webhook compares cleaned headers and preserves ordered event comparison', () => {
  const original = { ...webhookForm(), headers: { Token: 'value' }, events: ['a', 'b'] };
  const form = { ...webhookForm(original), headers: { ' Token ': ' value ', Blank: '' }, events: ['b', 'a'] };
  assert.deepEqual(buildWebhookPayload(form, original), { events: ['b', 'a'] });
  assert.deepEqual(form.headers, { ' Token ': ' value ', Blank: '' });
});

test('SMTP creation preserves default port and supplied credentials', () => {
  assert.deepEqual(smtpForm(), { address: '', host: '', port: 587, pwd: '', on: false });
  const form = { ...smtpForm(), pwd: ' secret with spaces ' };
  assert.deepEqual(buildSmtpPayload(form), form);
  assert.notEqual(buildSmtpPayload(form), form);
});

test('SMTP editing never preloads a password and omits a blank password', () => {
  const original = { id: 3, ...smtpForm(), pwd: 'server-secret' };
  const form = smtpForm(original);
  assert.equal(form.pwd, '');
  assert.deepEqual(buildSmtpPayload(form, original), {});
  assert.equal(Object.hasOwn(buildSmtpPayload(form, original), 'pwd'), false);
  assert.deepEqual(buildSmtpPayload({ ...form, pwd: 'replacement' }, original), { pwd: 'replacement' });
});

test('SMTP delta retains empty strings, zero and disabling without server-only fields', () => {
  const original = { id: 1, address: 'a@example.com', host: 'smtp.example.com', port: 587, on: true, success: 8 };
  const form = { ...smtpForm(original), address: '', host: '', port: 0, on: false };
  assert.deepEqual(buildSmtpPayload(form, original), { address: '', host: '', port: 0, on: false });
});

test('OAuth editing formats scopes, normalizes optional fields and never copies secrets', () => {
  const form = oauthForm({ protocol: '', scopes: ['openid', 'profile'], client_secret: 'server-secret', id: 10 });
  assert.equal(form.protocol, 'oauth2');
  assert.equal(form.scopes, 'openid, profile');
  assert.equal(form.client_secret, '');
  assert.equal(form.callback_url, '');
  assert.equal(form.default_group, 0);
  assert.equal(Object.hasOwn(form, 'id'), false);
  assert.equal(oauthForm({ scopes: null }).scopes, '');
});

test('OAuth scopes trim and remove empty entries without changing order or deduplicating', () => {
  assert.deepEqual(parseScopes(' openid, , profile ,openid, '), ['openid', 'profile', 'openid']);
  assert.deepEqual(parseScopes(' , '), []);
});

test('OAuth2 payload retains protocol for the API query and omits blank client secret', () => {
  const form = { ...oauthForm(), scopes: 'openid, profile', token_url: '/token', client_id: 'client' };
  const payload = buildOAuthPayload(form);
  assert.equal(payload.protocol, 'oauth2');
  assert.equal(payload.token_url, '/token');
  assert.equal(payload.client_id, 'client');
  assert.deepEqual(payload.scopes, ['openid', 'profile']);
  assert.equal(Object.hasOwn(payload, 'client_secret'), false);
  assert.equal(buildOAuthPayload({ ...form, client_secret: ' new-secret ' }).client_secret, ' new-secret ');
  assert.equal(form.scopes, 'openid, profile');
});

test('CAS payload excludes OAuth-only fields even after switching a filled OAuth2 form', () => {
  const form = {
    ...oauthForm(),
    protocol: 'cas',
    provider: 'CAS',
    uri: 'cas',
    auth_url: '/login',
    user_info_url: '/validate',
    callback_url: '/callback',
    picture: '/logo.png',
    id_claim: 'id',
    name_claim: 'name',
    email_claim: 'email',
    picture_claim: 'picture',
    description_claim: 'description',
    groups_claim: 'groups',
    admin_group: 'admins',
    default_group: 42,
    on: true,
    token_url: '/token',
    client_id: 'client',
    client_secret: 'secret',
    scopes: 'openid',
  };
  assert.deepEqual(buildOAuthPayload(form), {
    protocol: 'cas',
    provider: 'CAS',
    uri: 'cas',
    auth_url: '/login',
    user_info_url: '/validate',
    callback_url: '/callback',
    picture: '/logo.png',
    id_claim: 'id',
    name_claim: 'name',
    email_claim: 'email',
    picture_claim: 'picture',
    description_claim: 'description',
    groups_claim: 'groups',
    admin_group: 'admins',
    default_group: 42,
    on: true,
  });
  assert.equal(buildOAuthPayload({ ...form, protocol: 'oauth2' }).client_secret, 'secret');
});

test('OAuth updates remain full protocol payloads, including clears and disabled/default group values', () => {
  const payload = buildOAuthPayload(oauthForm());
  assert.equal(payload.callback_url, '');
  assert.equal(payload.groups_claim, '');
  assert.equal(payload.admin_group, '');
  assert.equal(payload.default_group, 0);
  assert.equal(payload.on, false);
  assert.deepEqual(payload.scopes, []);
});
