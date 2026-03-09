/**
 * 解题进度组件
 * @param {Object} props
 * @param {Array} props.solved - 各类型解题数据
 * @param {number} props.solved[].all - 该类型总题目数
 * @param {string} props.solved[].category - 题目类型
 * @param {number} props.solved[].solved - 该类型已解题数
 * @param {number} props.totalSolved - 总解题数
 * @param {string} props.totalLabel - 总计标签
 */

import { motion } from 'motion/react';
import { ScrollingText } from '../../common/index.js';

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
    <div className="flex items-center justify-center gap-3">
      {solved.map(({ category, solved: solvedCount, all }) => (
        <div key={category} className="flex flex-col items-center" title={`${category}: ${solvedCount}/${all}`}>
          <ScrollingText
            text={category}
            className="text-[10px] font-mono text-neutral-400 mb-1"
            maxWidth={50}
            speed={15}
          />
          <div className={`h-6 w-1.5 rounded-full bg-neutral-700 relative overflow-hidden`}>
            <motion.div
              className={`absolute bottom-0 left-0 right-0 ${categoryColors[category.toUpperCase()] || 'bg-neutral-400'}`}
              initial={{ height: 0 }}
              animate={{
                height: `${(solvedCount / all) * 100}%`,
              }}
              transition={{ duration: 1 }}
            />
          </div>
          <div className="text-[10px] font-mono text-neutral-300 mt-1">{solvedCount}</div>
        </div>
      ))}
      <div className="h-full w-[1px] bg-neutral-300/10 mx-2" />
      <div className="flex flex-col items-center justify-center">
        <div className="text-2xl font-mono text-geek-400">{totalSolved}</div>
        <div className="text-[10px] text-neutral-400">{totalLabel}</div>
      </div>
    </div>
  );
}

export default ChallengeSolves;
