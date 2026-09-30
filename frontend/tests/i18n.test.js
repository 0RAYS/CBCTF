import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import { createElement } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { createInstance } from 'i18next';
import { Trans } from 'react-i18next';
import { auditMessages } from '../scripts/i18n-audit.mjs';
import en from '../src/i18n/locales/en.json' with { type: 'json' };
import zh from '../src/i18n/locales/zh-CN.json' with { type: 'json' };

test('translation audit checks JSX, conditional keys, aliases and interpolation without counting comments', () => {
  const report = auditMessages(
    {
      en: { common: { active: 'Active', unused: 'Unused', alias: '$t(common.active)', name: 'Hello {{name}}' } },
      'zh-CN': { common: { active: '启用', unused: '未使用', alias: '$t(common.active)', name: '你好 {{user}}' } },
    },
    new Map([
      [
        'src/example.jsx',
        `
    // t('common.unused')
    const label = t(ok ? 'common.active' : 'typo.missing');
    const message = <Trans i18nKey="common.name" />;
    const alias = t('common.alias');
  `,
      ],
    ])
  );
  assert.deepEqual(report.unused, ['common.unused']);
  assert.deepEqual(
    report.missingReferences.map(({ key }) => key),
    ['typo.missing']
  );
  assert.deepEqual(report.placeholderMismatch, ['common.name']);
  assert.deepEqual(report.aliases, ['common.alias']);
});

test('translation audit preserves dynamic scopes, metadata references and plural variants', () => {
  const messages = {
    domain: {
      state: { running: 'Running', stopped: 'Stopped' },
      title: 'Title',
      count_one: '{{count}} item',
      count_other: '{{count}} items',
      unused: 'Unused',
    },
  };
  const report = auditMessages(
    { en: messages, 'zh-CN': messages },
    new Map([
      [
        'src/example.js',
        `
    const metadata = { titleKey: 'domain.title' };
    t(metadata.titleKey);
    t('domain.count', { count: 2 });
    t(\`domain.state.\${status}\`);
  `,
      ],
    ])
  );
  assert.deepEqual(report.unused, ['domain.unused']);
  assert.deepEqual(report.missingReferences, []);
  assert.deepEqual(report.mismatch, []);
});

test('application locales have matching keys/parameters, valid references and no unused candidates or aliases', () => {
  const script = fileURLToPath(new URL('../scripts/i18n-audit.mjs', import.meta.url));
  const report = JSON.parse(execFileSync(process.execPath, [script, '--check'], { encoding: 'utf8' }));
  assert.equal(report.keys.en, report.keys['zh-CN']);
  assert.equal(report.unusedCandidates, 0);
  assert.deepEqual(report.missingReferences, []);
});

test('shared workload translations and the complete batch summary render in both languages', async () => {
  const i18n = createInstance();
  await i18n.init({
    lng: 'en',
    fallbackLng: false,
    resources: { en: { translation: en }, 'zh-CN': { translation: zh } },
    interpolation: { escapeValue: false },
  });
  for (const lng of ['en', 'zh-CN']) {
    await i18n.changeLanguage(lng);
    for (const status of ['waiting', 'pending', 'terminating', 'running', 'stopped']) {
      assert.ok(i18n.exists(`admin.generators.status.${status}`));
      assert.ok(i18n.exists(`admin.victims.statusBadge.${status}`));
    }
    for (const key of ['teamName', 'searchTeamPlaceholder', 'teamFallback', 'teamIdLabel']) {
      assert.ok(i18n.exists(`admin.victims.filters.${key}`));
    }
    assert.ok(i18n.exists('admin.victims.table.selectContainer'));
    assert.ok(i18n.exists('admin.cronjobs.columns.runOnStart'));
    const html = renderToStaticMarkup(
      createElement(Trans, {
        i18n,
        i18nKey: 'admin.contests.containers.modals.summary',
        values: { challenges: 3, teams: 4, total: 12 },
        components: { amount: createElement('span'), total: createElement('strong') },
      })
    );
    assert.ok(html.includes('<span>3</span>'));
    assert.ok(html.includes('<span>4</span>'));
    assert.ok(html.includes('<strong>12</strong>'));
    assert.ok(!html.includes('{{') && !html.includes('admin.contests'));
    assert.equal(i18n.exists('admin.contests.generators'), false);
    assert.equal(i18n.exists('admin.contests.teamContainers'), false);
    assert.equal(i18n.exists('game.detail.defaultRules'), false);
  }
});
