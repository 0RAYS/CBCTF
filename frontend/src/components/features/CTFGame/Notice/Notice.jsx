/**
 * 公告页面组件
 * @param {Object} props
 * @param {Array} props.notices - 公告列表
 * @example
 * const notices = [{
 *   id: 1,
 *   title: "Challenge 'Web Injection' Updated",
 *   content: "We've updated the challenge description...",
 *   type: "update", // update/important/normal
 *   timestamp: "2024-03-15 14:30:22"
 * }]
 */

import { motion } from 'motion/react';
import { useState } from 'react';

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

  return (
    <div className="contest-container mx-auto">
      {/* 整体容器 */}
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        {/* 标题栏 */}
        <div className="flex items-center justify-end mb-6">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-red-400" />
              <span className="text-neutral-400 text-sm">Important</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-yellow-400" />
              <span className="text-neutral-400 text-sm">Update</span>
            </div>
          </div>
        </div>

        {/* 公告列表 */}
        <div className="space-y-4">
          {notices.map((notice) => (
            <motion.div
              key={notice.id}
              className={`border rounded-md overflow-hidden transition-colors duration-200
                                ${typeColors[notice.type]}
                                hover:border-neutral-100`}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.2 }}
              onClick={() => setExpandedId(expandedId === notice.id ? null : notice.id)}
            >
              <div className="p-4 cursor-pointer transition-colors duration-200 hover:bg-white/5">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-3">
                    <span className="text-xl select-none">{typeIcons[notice.type]}</span>
                    <h3 className="font-mono text-neutral-50 transition-colors duration-200">{notice.title}</h3>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-neutral-400 text-sm font-mono">{notice.timestamp}</span>
                  </div>
                </div>

                <motion.div
                  initial={false}
                  animate={{
                    height: expandedId === notice.id ? 'auto' : 0,
                    opacity: expandedId === notice.id ? 1 : 0,
                  }}
                  transition={{ duration: 0.2 }}
                  className="overflow-hidden"
                >
                  <div className="py-2 text-neutral-300 whitespace-pre-wrap">{notice.content}</div>
                </motion.div>
              </div>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </div>
  );
}

export default Notice;
