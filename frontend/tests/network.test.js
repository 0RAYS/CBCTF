import assert from 'node:assert/strict';
import test from 'node:test';
import { injectClickableIps, isPublicIp } from '../src/components/features/Admin/network/logIps.js';
import { ansiToHtml } from '../src/utils/ansi.js';

test('public IPv4 classification excludes private, reserved, multicast, and invalid addresses', () => {
  for (const ip of ['1.1.1.1', '8.8.8.8', '100.63.255.255', '100.128.0.1', '172.15.0.1', '172.32.0.1']) {
    assert.equal(isPublicIp(ip), true, ip);
  }
  for (const ip of [
    '0.0.0.0', '10.1.2.3', '127.0.0.1', '172.16.0.1', '172.31.255.255', '192.168.1.1',
    '169.254.1.1', '100.64.0.1', '100.127.255.255', '198.18.0.1', '198.19.0.1',
    '192.0.0.1', '192.0.2.1', '192.88.99.1', '198.51.100.1', '203.0.113.1',
    '224.0.0.1', '239.255.255.255', '240.0.0.1', '255.255.255.255',
    '256.1.1.1', '1.2.3', '1.2.3.4.5', '1.2.3.-1', '1.2.3.4x', ' 8.8.8.8', '1..2.3', '', null,
  ]) {
    assert.equal(isPublicIp(ip), false, String(ip));
  }
});

test('IP decoration preserves ANSI tags and attributes and adds keyboard semantics only to text', () => {
  const html = '<span style="color:red" title="8.8.8.8">1.1.1.1 10.0.0.1</span>';
  const decorated = injectClickableIps(html);
  assert.ok(decorated.startsWith('<span style="color:red" title="8.8.8.8">'));
  assert.match(decorated, /data-ip="1\.1\.1\.1" role="button" tabindex="0"/);
  assert.doesNotMatch(decorated, /data-ip="(?:8\.8\.8\.8|10\.0\.0\.1)"/);
  assert.ok(decorated.endsWith(' 10.0.0.1</span>'));
});

test('IP decoration does not match partial invalid addresses and supports ports and punctuation', () => {
  assert.doesNotMatch(injectClickableIps('1.2.3.4.5 999.8.8.8 8.8.8.888 a8.8.8.8 8.8.8.8x'), /data-ip/);
  assert.equal((injectClickableIps('(8.8.8.8:443), 1.1.1.1.').match(/data-ip=/g) || []).length, 2);
});

test('untrusted log markup remains escaped through ANSI conversion and IP decoration', () => {
  const html = injectClickableIps(ansiToHtml('\u001b[31m<img src=x onerror="alert(1)"> 8.8.8.8\u001b[0m'));
  assert.match(html, /&lt;img src=x onerror=&quot;alert\(1\)&quot;&gt;/);
  assert.doesNotMatch(html, /<img/);
  assert.match(html, /style="color:#ff6b6b;"/);
  assert.match(html, /data-ip="8\.8\.8\.8"/);
});
