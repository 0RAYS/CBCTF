import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import test from 'node:test';
import { normalizeConfig } from '../src/components/features/Admin/SystemConfig/configNormalizer.js';
import { buildPayload } from '../src/components/features/Admin/SystemConfig/configPayloadBuilder.js';

test('system config payload covers every backend editable field', () => {
  const dto = readFileSync(new URL('../../internal/dto/setting.go', import.meta.url), 'utf8');
  const fields = [...dto.matchAll(/json:"([^"]+)"/g)].map((match) => match[1]);
  assert.deepEqual(Object.keys(buildPayload(normalizeConfig({}))).sort(), fields.sort());
});

test('workload settings round-trip false, zero and an explicitly cleared priority class', () => {
  for (const source of [
    {
      k8s_capture_enabled: true,
      k8s_priority_class_name: 'ctf-workloads',
      k8s_worker_image: 'registry.example/worker:v3',
      k8s_generator_pool_size: 4,
    },
    {
      k8s_capture_enabled: false,
      k8s_priority_class_name: '',
      k8s_worker_image: 'ghcr.io/0rays/cbctf-worker:latest',
      k8s_generator_pool_size: 0,
    },
  ]) {
    const before = structuredClone(source);
    const payload = JSON.parse(JSON.stringify(buildPayload(normalizeConfig(source))));
    for (const [key, value] of Object.entries(source)) assert.equal(payload[key], value, key);
    assert.deepEqual(source, before);
    for (const key of ['path', 'gorm_postgres_pwd', 'redis_pwd', 'gin_host', 'gin_port']) {
      assert.equal(Object.hasOwn(payload, key), false, `deployment field leaked into update: ${key}`);
    }
  }
});

test('Kubernetes form labels and descriptions exist in both languages', () => {
  const form = readFileSync(
    new URL('../src/components/features/Admin/SystemConfig/sections/K8sConfigSection.jsx', import.meta.url),
    'utf8',
  );
  for (const language of ['en', 'zh-CN']) {
    const messages = JSON.parse(readFileSync(new URL(`../src/i18n/locales/${language}.json`, import.meta.url), 'utf8'));
    const keys = [...form.matchAll(/t\(['"](admin\.system\.[^'"]+)['"]\)/g)].map((match) => match[1]);
    assert.ok(keys.length >= 4, 'expected Kubernetes form translations');
    for (const key of keys) {
      const value = key.split('.').reduce((data, part) => data?.[part], messages);
      assert.equal(typeof value, 'string', `${language}: missing ${key}`);
    }
  }
});
