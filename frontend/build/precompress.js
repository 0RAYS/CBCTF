import { writeFile } from 'node:fs/promises';
import { join } from 'node:path';
import { gzipSync } from 'node:zlib';

export default function precompress() {
  return {
    name: 'precompress-static-assets',
    apply: 'build',
    async writeBundle(output, bundle) {
      await Promise.all(
        Object.values(bundle).map(async (asset) => {
          if (!/\.(?:js|css|html|svg)$/.test(asset.fileName)) return;
          const content = Buffer.from(asset.type === 'chunk' ? asset.code : asset.source);
          if (content.length < 1024) return;
          const compressed = gzipSync(content, { level: 9 });
          if (compressed.length >= content.length * 0.9) return;
          // Build-time work: the Go server streams these bytes without per-request compression.
          await writeFile(join(output.dir, `${asset.fileName}.gz`), compressed);
        })
      );
    },
  };
}
