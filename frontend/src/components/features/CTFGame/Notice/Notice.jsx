import { motion } from 'motion/react';
import { useId, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconAlertTriangle, IconChevronDown, IconNote, IconSpeakerphone } from '@tabler/icons-react';

const typeColors = {
  important: 'border-red-400 bg-red-400/5',
  update: 'border-yellow-400 bg-yellow-400/5',
  normal: 'border-neutral-300 bg-black/30',
};

const typeIcons = {
  important: { Icon: IconAlertTriangle, colorClass: 'text-red-400' },
  update: { Icon: IconNote, colorClass: 'text-yellow-400' },
  normal: { Icon: IconSpeakerphone, colorClass: 'text-neutral-400' },
};

/** 公告页面组件 */
function Notice({ notices }) {
  const [expandedId, setExpandedId] = useState(null);
  const id = useId();
  const { t } = useTranslation();
  const typeLabels = {
    important: t('common.notices.important'),
    update: t('common.notices.update'),
    normal: t('common.notices.normal'),
  };

  return (
    <div className="contest-container mx-auto">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <div className="flex items-center justify-end mb-6">
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-red-400" />
              <span className="text-neutral-400 text-sm">{typeLabels.important}</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-yellow-400" />
              <span className="text-neutral-400 text-sm">{typeLabels.update}</span>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          {notices.map((notice) => {
            const { Icon, colorClass } = typeIcons[notice.type] || typeIcons.normal;
            return (
              <motion.article
                key={notice.id}
                className={`border rounded-md overflow-hidden transition-colors duration-200
                                ${typeColors[notice.type] || typeColors.normal}`}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.2 }}
              >
                <h3>
                  <button
                    type="button"
                    id={`${id}-title-${notice.id}`}
                    aria-expanded={expandedId === notice.id}
                    aria-controls={`${id}-content-${notice.id}`}
                    onClick={() => setExpandedId((current) => (current === notice.id ? null : notice.id))}
                    className="flex w-full items-start gap-3 p-4 text-left transition-colors hover:bg-white/5 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-geek-400"
                  >
                    <Icon aria-hidden="true" size={20} stroke={1.5} className={`shrink-0 mt-0.5 ${colorClass}`} />
                    <span className="min-w-0 flex-1 [overflow-wrap:anywhere]">
                      <span className="block text-xs uppercase tracking-wider text-neutral-400 font-mono">
                        {typeLabels[notice.type] || notice.type || typeLabels.normal}
                      </span>
                      <span className="block font-mono text-neutral-50 mt-1">{notice.title}</span>
                      <span className="block text-neutral-400 text-xs font-mono mt-2">{notice.timestamp}</span>
                    </span>
                    <IconChevronDown
                      aria-hidden="true"
                      size={20}
                      stroke={1.5}
                      className={`mt-1 shrink-0 text-neutral-400 transition-transform ${expandedId === notice.id ? 'rotate-180' : ''}`}
                    />
                  </button>
                </h3>
                <div
                  id={`${id}-content-${notice.id}`}
                  role="region"
                  aria-labelledby={`${id}-title-${notice.id}`}
                  hidden={expandedId !== notice.id}
                  className="border-t border-neutral-300/10 p-4 text-neutral-300 whitespace-pre-wrap leading-relaxed [overflow-wrap:anywhere]"
                >
                  {notice.content}
                </div>
              </motion.article>
            );
          })}
        </div>
      </motion.div>
    </div>
  );
}

export default Notice;
