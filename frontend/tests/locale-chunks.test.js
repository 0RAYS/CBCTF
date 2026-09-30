import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import { createInstance } from 'i18next';
import test from 'node:test';
import localeChunks from '../build/localeChunks.js';

test('public/admin locale chunks partition both dictionaries without losing messages or placeholders', async () => {
  const resources = {};
  const adminResources = {};
  for (const language of ['en', 'zh-CN']) {
    const file = fileURLToPath(new URL(`../src/i18n/locales/${language}.json`, import.meta.url));
    const original = JSON.parse(await readFile(file, 'utf8'));
    const watched = new Set();
    const load = async (scope) => {
      const result = await localeChunks().load.call({ addWatchFile: (path) => watched.add(path) }, `${file}?scope=${scope}`);
      return (await import(`data:text/javascript,${encodeURIComponent(result.code)}`)).default;
    };
    const publicMessages = await load('public');
    const adminMessages = await load('admin');
    assert.equal(publicMessages.admin, undefined);
    assert.deepEqual(Object.keys(adminMessages), ['admin']);
    assert.deepEqual({ ...publicMessages, ...adminMessages }, original);
    assert.deepEqual([...watched], [file], 'translation edits must trigger dev-server invalidation');
    resources[language] = { translation: publicMessages };
    adminResources[language] = adminMessages;
  }
  const i18n = createInstance();
  await i18n.init({ lng: 'zh-CN', fallbackLng: 'en', resources, interpolation: { escapeValue: false } });
  assert.equal(i18n.t('common.notices.important'), '重要');
  assert.equal(i18n.exists('admin.title'), false);
  for (const language of ['en', 'zh-CN']) {
    i18n.addResourceBundle(language, 'translation', adminResources[language], true, true);
  }
  for (const language of ['en', 'zh-CN']) {
    await i18n.changeLanguage(language);
    assert.ok(i18n.exists('admin.title'));
    assert.ok(i18n.exists('common.notices.fetchFailed'));
    assert.ok(!i18n.t('admin.topbar.managementTitle', { label: 'Example' }).includes('{{'));
  }
});
