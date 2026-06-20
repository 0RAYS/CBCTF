import * as path from 'node:path';
import { defineConfig } from '@rspress/core';

const siteBase = process.env.DOCS_BASE || '/';

function sanitizeMdxStyleAttributes() {
  return (tree: any) => {
    const visit = (node: any) => {
      if (node?.properties && 'style' in node.properties) {
        delete node.properties.style;
      }

      if (Array.isArray(node?.children)) {
        node.children.forEach(visit);
      }
    };

    visit(tree);
  };
}

export default defineConfig({
  root: path.join(__dirname, 'docs'),
  outDir: 'doc_build',
  base: siteBase,
  lang: 'zh',
  title: 'CBCTF',
  description: 'CBCTF 平台部署、运维与功能使用文档',
  icon: '/logo.png',
  logo: {
    light: '/logo.png',
    dark: '/logo.png',
  },
  llms: true,
  search: {
    codeBlocks: true,
    versioned: false,
  },
  globalStyles: path.join(__dirname, 'styles/global.css'),
  markdown: {
    rehypePlugins: [sanitizeMdxStyleAttributes],
  },
  builderConfig: {
    html: {
      tags: [
        {
          tag: 'script',
          children: "window.RSPRESS_THEME = 'dark';",
        },
      ],
    },
  },
  themeConfig: {
    lastUpdated: true,
    editLink: {
      docRepoBaseUrl: 'https://github.com/0RAYS/CBCTF/tree/main/docs/docs',
    },
    llmsUI: {
      placement: 'outline',
    },
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
