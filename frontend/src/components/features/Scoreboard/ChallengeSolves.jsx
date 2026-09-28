/**
 * 积分榜解题进度组件
 * @param {Object} props
 * @param {Array} props.solved - 各类型解题数据
 * @param {number} props.solved[].all - 该类型总题目数
 * @param {string} props.solved[].category - 题目类型
 * @param {number} props.solved[].solved - 该类型已解题数
 * @param {number} props.totalSolved - 总解题数
 * @param {string} props.totalLabel - 总计标签
 */

import { motion } from 'motion/react';
import TruncatedText from '../../common/TruncatedText';

const categoryColors = {
  WEB: 'bg-geek-400',
  CRYPTO: 'bg-purple-400',
  PWN: 'bg-red-400',
  REVERSE: 'bg-green-400',
  MISC: 'bg-yellow-400',
};

function ChallengeSolves({ solved = [], totalSolved = 0, totalLabel = 'TOTAL' }) {
  if (!Array.isArray(solved) || solved.length === 0) return null;

  return (
    <div className="inline-flex max-w-full min-w-0 items-center gap-2">
      <div className="flex min-w-0 flex-wrap items-center gap-x-1 gap-y-1.5 sm:gap-x-2">
        {solved.map(({ category, solved: solvedCount, all }) => (
          <div key={category} className="w-13 sm:w-16 min-w-0 space-y-1" title={`${category}: ${solvedCount}/${all}`}>
            <div className="flex items-baseline justify-between gap-1 font-mono text-[10px]">
              <TruncatedText className="text-neutral-400">{category}</TruncatedText>
              <TruncatedText maxWidth={28} className="shrink-0 text-neutral-200 tabular-nums">
                {solvedCount}
              </TruncatedText>
            </div>
            <div className="h-1 rounded-full bg-neutral-700 relative overflow-hidden" aria-hidden="true">
              <motion.div
                className={`absolute inset-y-0 left-0 ${categoryColors[category.toUpperCase()] || 'bg-neutral-400'}`}
                initial={{ width: 0 }}
                animate={{
                  width: `${all > 0 ? Math.min(100, Math.max(0, (solvedCount / all) * 100)) : 0}%`,
                }}
                transition={{ duration: 0.3 }}
              />
            </div>
          </div>
        ))}
      </div>
      <div className="flex shrink-0 max-w-20 flex-col items-center justify-center border-l border-neutral-600/40 pl-2 font-mono">
        <TruncatedText className="text-[10px] text-neutral-400">{totalLabel}</TruncatedText>
        <TruncatedText className="text-sm font-semibold leading-4 tabular-nums text-geek-400">
          {totalSolved}
        </TruncatedText>
      </div>
    </div>
  );
}

export default ChallengeSolves;
