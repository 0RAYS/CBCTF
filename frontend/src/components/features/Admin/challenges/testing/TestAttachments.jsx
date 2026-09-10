import { useLayoutEffect, useRef, useState } from 'react';
import { motion } from 'motion/react';
import { IconLoader2, IconPaperclip } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { downloadTestAttachment } from '../../../../../api/admin/challenge';
import { downloadBlobResponse } from '../../../../../utils/fileDownload';
import { toast } from '../../../../../utils/toast';

export default function TestAttachments({ challenge, file }) {
  const { t } = useTranslation();
  const [downloading, setDownloading] = useState(false);
  const lifetime = useRef(null);
  useLayoutEffect(() => {
    const session = { active: true, downloading: false };
    lifetime.current = session;
    return () => {
      session.active = false;
    };
  }, []);

  async function download() {
    const session = lifetime.current;
    if (!challenge?.id || !session?.active || session.downloading) return;
    session.downloading = true;
    setDownloading(true);
    try {
      const response = await downloadTestAttachment(challenge.id);
      if (!session.active) return;
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, 'attachment.zip', 'application/octet-stream');
        toast.success({ description: t('admin.challenge.testModal.toast.downloadSuccess') });
      }
    } catch (error) {
      if (session.active) {
        toast.danger({ description: error.message || t('admin.challenge.testModal.toast.downloadFailed') });
      }
    } finally {
      session.downloading = false;
      if (session.active) setDownloading(false);
    }
  }

  if (challenge?.type !== 'dynamic' && !file) return null;

  return (
    <section className="space-y-1.5">
      <h3 className="text-neutral-400 font-mono text-sm">{t('admin.challenge.testModal.sections.attachments')}</h3>
      <motion.button
        type="button"
        className="flex w-full flex-wrap items-center gap-2 p-2 border border-neutral-300/30 rounded-md text-left
          text-neutral-300 hover:text-geek-400 hover:border-geek-400 transition-colors duration-200
          disabled:cursor-wait disabled:opacity-60"
        whileHover={downloading ? undefined : { x: 5 }}
        onClick={download}
        disabled={downloading || !challenge?.id}
        aria-busy={downloading}
      >
        {downloading ? (
          <IconLoader2 size={16} className="shrink-0 animate-spin" aria-hidden="true" />
        ) : (
          <IconPaperclip size={16} className="shrink-0" aria-hidden="true" />
        )}
        <span className="min-w-0 font-mono text-sm [overflow-wrap:anywhere]">{file}</span>
        <span className="text-neutral-400 text-sm ml-auto">
          {downloading
            ? t('admin.challenge.testModal.attachments.generating')
            : challenge?.type === 'dynamic'
              ? t('admin.challenge.testModal.attachments.generateAndDownload')
              : t('admin.challenge.testModal.attachments.download')}
        </span>
      </motion.button>
    </section>
  );
}
