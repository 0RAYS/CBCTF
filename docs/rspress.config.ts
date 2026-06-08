import * as path from 'node:path';
import { defineConfig } from '@rspress/core';

export default defineConfig({
  root: path.join(__dirname, 'docs'),
  lang: 'zh',
  title: 'CBCTF',
  description: '基于 Kubernetes 的现代化 CTF 竞赛平台',
  icon: '/logo.png',
  logo: {
    light: '/logo.png',
    dark: '/logo.png',
  },
  themeConfig: {
    socialLinks: [
      {
        icon: 'github',
        mode: 'link',
        content: 'https://github.com/0RAYS/CBCTF',
      },
    ],
    footer: {
      message: 'AGPL-3.0 Licensed | © 0RAYS',
    },
  },
});
