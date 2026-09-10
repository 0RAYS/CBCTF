/**
 * 公告页面组件
 * @param {Object} props
 * @param {Array} props.notices - 公告列表
 * @example
 * const notices = [{
 *   id: 1,
 *   title: "Challenge 'Web Injection' Updated",
 *   content: "We've updated the challenge description...",
 *   type: "update",
 *   timestamp: "2024-03-15 14:30:22"
 * }]
 */

import { motion } from 'motion/react';
import { useId, useState } from 'react';
import { useTranslation } from 'react-i18next';

const typeColors = {
  important: 'border-red-400 bg-red-400/5',
  update: 'border-yellow-400 bg-yellow-400/5',
  normal: 'border-neutral-300 bg-black/30',
};

const typeIcons = {
  important: '⚠️',
  update: '📝',
  normal: '📢',
};

function Notice({ notices }) {
  const [expandedId, setExpandedId] = useState(null);
  const id = useId();
  const { t } = useTranslation();
  const typeLabels = {
    important: t('admin.contests.notices.types.important'),
    update: t('admin.contests.notices.types.update'),
    normal: t('admin.contests.notices.types.normal'),
  };

  return (
    <div className="contest-container mx-auto">
      {/* 整体容器 */}
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        {/* 标题栏 */}
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

        {/* 公告列表 */}
        <div className="space-y-4">
          {notices.map((notice) => (
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
                  <span aria-hidden="true" className="shrink-0 text-xl select-none">
                    {typeIcons[notice.type] || typeIcons.normal}
                  </span>
                  <span className="min-w-0 flex-1 [overflow-wrap:anywhere]">
                    <span className="block text-xs uppercase tracking-wider text-neutral-400 font-mono">
                      {typeLabels[notice.type] || notice.type || typeLabels.normal}
                    </span>
                    <span className="block font-mono text-neutral-50 mt-1">{notice.title}</span>
                    <span className="block text-neutral-400 text-xs font-mono mt-2">{notice.timestamp}</span>
                  </span>
                  <svg
                    aria-hidden="true"
                    viewBox="0 0 20 20"
                    fill="none"
                    className={`mt-1 size-5 shrink-0 text-neutral-400 transition-transform ${expandedId === notice.id ? 'rotate-180' : ''}`}
                  >
                    <path d="m5 7.5 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
                  </svg>
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
          ))}
        </div>
      </motion.div>
    </div>
  );
}

export default Notice;
