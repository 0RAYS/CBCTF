import { motion } from 'motion/react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useState } from 'react';
import { Button, Card } from '../../../../components/common';

/**
 * @param {Object} contest - 比赛详情对象
 * @param handleJoinContest
 * @param {string} contest.title - 比赛标题 (例: "CTF Winter Challenge 2024")
 * @param {string} contest.description - 比赛描述 (例: "Join our winter themed CTF challenge...")
 * @param {string} contest.image - 比赛背景图片URL (例: "/images/winter-ctf.jpg")
 * @param {string} contest.status - 比赛状态 ("upcoming" | "active" | "ended")
 * @param {string} contest.startTime - 开始时间 (例: "2024-01-01T00:00:00Z")
 * @param {string} contest.endTime - 结束时间 (例: "2024-01-07T23:59:59Z")
 * @param {number} contest.participants - 参与人数 (例: 256)
 * @param {Array<string>} contest.rules - 比赛规则列表 (例: ["No sharing of flags...", "Maximum team size is 4..."])
 * @param {Array<Object>} contest.prizes - 奖励列表
 * @param {string} contest.prizes[].amount - 奖金数额 (例: "$1000")
 * @param {string} contest.prizes[].description - 奖励描述 (例: "First place prize includes...")
 * @param {Array<Object>} contest.timeline - 时间线
 * @param {string} contest.timeline[].date - 日期 (例: "2025-02-16T15:00:00+08:00")
 * @param {string} contest.timeline[].title - 事件标题 (例: "Registration Opens")
 * @param {string} contest.timeline[].description - 事件描述 (例: "Team registration begins...")
 */

function ContestDetail({ contest, handleJoinContest }) {
  const [hoveredPrize, setHoveredPrize] = useState(null);
  const [hoveredTimeline, setHoveredTimeline] = useState(null);

  if (!contest) return null;

  return (
    <div className="contest-container mx-auto space-y-6">
      {/* 头部信息区域 - 移除悬停效果 */}
      <motion.div
        className="relative w-full h-[300px] border border-neutral-300 rounded-md overflow-hidden bg-black/30"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        {/* 背景图片 */}
        <div
          className="absolute inset-0 opacity-30"
          style={{
            backgroundImage: `url(${contest.image})`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
          }}
        />

        {/* 内容 */}
        <div className="relative h-full p-8 flex flex-col justify-between">
          <div>
            <div className="flex items-center gap-4 mb-4">
              <h1 className="text-4xl font-mono text-neutral-50 tracking-wider">{contest.title}</h1>
              <span
                className={`px-3 py-1 border rounded-md text-sm ${
                  contest.status === 'upcoming'
                    ? 'text-yellow-400 border-yellow-400'
                    : contest.status === 'active'
                      ? 'text-green-400 border-green-400'
                      : 'text-neutral-400 border-neutral-400'
                }`}
              >
                {contest.status.toUpperCase()}
              </span>
            </div>
            <div className="text-neutral-300 max-w-[800px] text-lg prose prose-invert">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>{contest.description || ''}</ReactMarkdown>
            </div>
          </div>

          <div className="flex items-center justify-between">
            <div className="flex gap-8">
              <div>
                <span className="text-neutral-400">Start Time</span>
                <div className="text-neutral-50 font-mono">{new Date(contest.startTime).toLocaleString()}</div>
              </div>
              <div>
                <span className="text-neutral-400">End Time</span>
                <div className="text-neutral-50 font-mono">{new Date(contest.endTime).toLocaleString()}</div>
              </div>
              <div>
                <span className="text-neutral-400">Participants</span>
                <div className="text-neutral-50 font-mono">{contest.participants || 0}</div>
              </div>
            </div>

            {/* 使用Button组件重构 */}
            <Button variant="primary" size="lg" disabled={contest.status === 'ended'} onClick={handleJoinContest}>
              {contest.status === 'upcoming' || contest.status === 'active' ? 'JOIN CONTEST' : 'CONTEST ENDED'}
            </Button>
          </div>
        </div>
      </motion.div>

      {/* 详细信息区域 */}
      <div className="grid grid-cols-3 gap-6">
        {/* 比赛规则 - 添加规则项悬停效果 */}
        <Card className="col-span-2" variant="default" padding="lg">
          <h2 className="text-2xl font-mono text-neutral-50 tracking-wider mb-6">Contest Rules</h2>
          <div className="space-y-4 text-neutral-300">
            {contest.rules?.map((rule, index) => (
              <motion.div
                key={index}
                className="flex gap-4 p-2 rounded-md hover:bg-neutral-300/5 transition-colors duration-200 cursor-default"
                whileHover={{ x: 10 }}
              >
                <span className="text-geek-400 font-mono">{(index + 1).toString().padStart(2, '0')}</span>
                <p>{rule}</p>
              </motion.div>
            ))}
          </div>
        </Card>

        {/* 奖励信息 - 添加奖项悬停效果 */}
        <Card variant="default" padding="lg">
          <h2 className="text-2xl font-mono text-neutral-50 tracking-wider mb-6">Prizes</h2>
          <div className="space-y-6">
            {contest.prizes?.map((prize, index) => (
              <motion.div
                key={index}
                className={`group p-4 rounded-md transition-all duration-200 cursor-default
                                    ${hoveredPrize === index ? 'bg-neutral-300/5' : 'hover:bg-neutral-300/5'}`}
                onMouseEnter={() => setHoveredPrize(index)}
                onMouseLeave={() => setHoveredPrize(null)}
                whileHover={{ scale: 1.02 }}
              >
                <motion.div
                  className="flex items-center gap-3 text-neutral-50 mb-2"
                  animate={{
                    x: hoveredPrize === index ? 10 : 0,
                    transition: { duration: 0.2 },
                  }}
                >
                  <span
                    className={`text-lg font-mono
                                        ${
                                          index === 0
                                            ? 'text-yellow-400'
                                            : index === 1
                                              ? 'text-neutral-300'
                                              : index === 2
                                                ? 'text-orange-400'
                                                : 'text-neutral-400'
                                        }
                                    `}
                  >
                    {index === 0 ? '1ST' : index === 1 ? '2ND' : index === 2 ? '3RD' : `${index + 1}TH`}
                  </span>
                  <span className="text-xl font-mono">{prize.amount}</span>
                </motion.div>
                <motion.div
                  className="text-neutral-400 text-sm prose prose-invert prose-sm max-w-none"
                  animate={{
                    x: hoveredPrize === index ? 10 : 0,
                    transition: { duration: 0.2, delay: 0.1 },
                  }}
                >
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>{prize.description || ''}</ReactMarkdown>
                </motion.div>
              </motion.div>
            ))}
          </div>
        </Card>
      </div>

      {/* 时间线 - 简化动画效果 */}
      <Card variant="default" padding="lg">
        <h2 className="text-2xl font-mono text-neutral-50 tracking-wider mb-6">Timeline</h2>
        <div className="relative flex items-start gap-8">
          {/* 连接线 */}
          <div className="absolute top-[30px] left-0 right-0 h-[2px] bg-neutral-300/20" />

          {contest.timeline?.map((item, index) => (
            <motion.div
              key={index}
              className={`flex-1 relative p-4 rounded-md transition-all duration-200 cursor-default
                                ${hoveredTimeline === index ? 'bg-neutral-300/5' : 'hover:bg-neutral-300/5'}`}
              onMouseEnter={() => setHoveredTimeline(index)}
              onMouseLeave={() => setHoveredTimeline(null)}
            >
              <motion.div
                animate={{
                  y: hoveredTimeline === index ? -5 : 0,
                  transition: { duration: 0.2 },
                }}
                className="border-none"
              >
                <div className="text-geek-400 font-mono mb-2">{new Date(item.date).toISOString().slice(0, 10)}</div>
                <div className="text-neutral-50 font-mono mb-1">{item.title}</div>
                <div className="text-neutral-400 text-sm prose prose-invert prose-sm max-w-none">
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>{item.description || ''}</ReactMarkdown>
                </div>
              </motion.div>
            </motion.div>
          ))}
        </div>
      </Card>
    </div>
  );
}

export default ContestDetail;
