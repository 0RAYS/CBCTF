import { readFile } from 'node:fs/promises';

// Keep one editable dictionary per language, but ship admin copy only with admin layouts.
export default function localeChunks() {
  return {
    name: 'locale-chunks',
    enforce: 'pre',
    async load(id) {
      const match = /[/\\]i18n[/\\]locales[/\\](?:en|zh-CN)\.json\?scope=(public|admin)$/.exec(id);
      if (!match) return;
      const file = id.slice(0, id.indexOf('?'));
      this.addWatchFile(file);
      const { admin, ...publicMessages } = JSON.parse(await readFile(file, 'utf8'));
      const messages = match[1] === 'admin' ? { admin } : publicMessages;
      return { code: `export default ${JSON.stringify(messages)};`, moduleType: 'js' };
    },
  };
}
