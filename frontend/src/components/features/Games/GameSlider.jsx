/**
 * 游戏轮播组件
 * @param {Object} props
 * @param {Array} props.games - 游戏列表
 * @param {Object} props.games[].id - 游戏ID
 * @param {string} props.games[].title - 游戏标题
 * @param {string} props.games[].description - 游戏描述
 * @param {string} props.games[].status - 游戏状态 (upcoming/active/ended)
 * @param {string} props.games[].startTime - 开始时间
 * @param {string} props.games[].endTime - 结束时间
 * @param {string} props.games[].image - 背景图片URL
 * @param {number} [props.currentIndex] - 当前显示的游戏索引（用于外部控制）
 * @param {Function} [props.onIndexChange] - 索引改变回调
 * @param {Function} [props.onGameAction] - 游戏操作回调（注册/加入/查看）
 * @example
 * const games = [{
 *   id: 1,
 *   title: "AI Innovation Challenge",
 *   description: "Join our AI challenge...",
 *   status: "active",
 *   startTime: "2024-03-01T00:00:00Z",
 *   endTime: "2024-04-30T23:59:59Z",
 *   image: "https://example.com/ai-challenge.jpg"
 * }]
 *
 * <GameSlider
 *   games={games}
 *   currentIndex={0}
 *   onIndexChange={(index) => console.log('Current index:', index)}
 *   onGameAction={(gameId, action) => console.log('Game action:', gameId, action)}
 * />
 */

import { motion, AnimatePresence } from 'motion/react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useState, useEffect } from 'react';
import Button from '../../common/Button';
import { useTranslation } from 'react-i18next';

function GameSlider({ games = [], currentIndex: externalIndex, onIndexChange, onGameAction, user }) {
  const [internalIndex, setInternalIndex] = useState(0);
  const { t } = useTranslation();

  // 使用外部索引或内部索引
  const currentIndex = typeof externalIndex === 'number' ? externalIndex : internalIndex;
  const currentGame = games[currentIndex];

  // 同步外部索引变化
  useEffect(() => {
    if (typeof externalIndex === 'number') {
      setInternalIndex(externalIndex);
    }
  }, [externalIndex]);

  const getStatusColor = (status) => {
    switch (status) {
      case 'upcoming':
        return 'text-yellow-400';
      case 'active':
        return 'text-green-400';
      case 'ended':
        return 'text-neutral-400';
      default:
        return 'text-neutral-300';
    }
  };

  const handleNext = () => {
    const nextIndex = (currentIndex + 1) % games.length;
    setInternalIndex(nextIndex);
    onIndexChange?.(nextIndex);
  };

  const handlePrev = () => {
    const prevIndex = (currentIndex - 1 + games.length) % games.length;
    setInternalIndex(prevIndex);
    onIndexChange?.(prevIndex);
  };

  const handleGameAction = () => {
    if (!currentGame) return;

    const action = currentGame.status === 'ended' ? 'view' : 'join';

    onGameAction?.(currentGame.id, action);
  };

  if (!currentGame) return null;

  const statusLabel = t(`game.status.${currentGame.status}`);

  return (
    <div className="w-full max-w-[1200px] mx-auto relative">
      <AnimatePresence mode="wait">
        <motion.div
          key={currentIndex}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.2 }}
          className="relative"
        >
          {/* 导航按钮 */}
          <div
            className={`absolute left-4 top-1/2 -translate-y-1/2 z-10 w-12 h-12
                            border rounded-md overflow-hidden
                            bg-neutral-900 group cursor-pointer
                            transition-all duration-200
                            ${
                              games.length <= 1
                                ? 'border-neutral-600 opacity-50 cursor-not-allowed'
                                : 'border-neutral-300 hover:border-geek-400'
                            }`}
            onClick={() => games.length > 1 && handlePrev()}
          >
            <motion.div
              className="w-full h-full flex items-center justify-center"
              whileTap={games.length > 1 ? { scale: 0.9 } : {}}
            >
              <span
                className={`text-xl transition-colors duration-200
                                ${
                                  games.length <= 1 ? 'text-neutral-600' : 'text-neutral-300 group-hover:text-geek-400'
                                }`}
              >
                &larr;
              </span>
            </motion.div>
          </div>

          <div
            className={`absolute right-4 top-1/2 -translate-y-1/2 z-10 w-12 h-12 
                            border rounded-md overflow-hidden
                            bg-neutral-900 group cursor-pointer
                            transition-all duration-200
                            ${
                              games.length <= 1
                                ? 'border-neutral-600 opacity-50 cursor-not-allowed'
                                : 'border-neutral-300 hover:border-geek-400'
                            }`}
            onClick={() => games.length > 1 && handleNext()}
          >
            <motion.div
              className="w-full h-full flex items-center justify-center"
              whileTap={games.length > 1 ? { scale: 0.9 } : {}}
            >
              <span
                className={`text-xl transition-colors duration-200
                                ${
                                  games.length <= 1 ? 'text-neutral-600' : 'text-neutral-300 group-hover:text-geek-400'
                                }`}
              >
                &rarr;
              </span>
            </motion.div>
          </div>

          {/* 游戏展示区域 */}
          <motion.div
            initial={{ opacity: 0, x: 100 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -100 }}
            transition={{ duration: 0.3 }}
            className="relative h-[500px] border border-neutral-300 rounded-md overflow-hidden bg-black/30"
          >
            {/* 背景图片 */}
            <div
              className="absolute inset-0 opacity-40"
              style={{
                backgroundImage: `url(${currentGame.image})`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
              }}
            />

            {/* 内容区域 */}
            <div className="relative h-full p-8 flex flex-col justify-between">
              {/* 标题和状态 */}
              <div>
                <div className="flex items-center gap-4 mb-4">
                  <h2 className="text-3xl font-mono text-neutral-50 tracking-wider">{currentGame.title}</h2>
                  <span className={`px-3 py-1 border rounded-md text-sm ${getStatusColor(currentGame.status)}`}>
                    {statusLabel}
                  </span>
                </div>
                <div className="text-neutral-300 max-w-[600px] prose prose-invert prose-sm line-clamp-3">
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>{currentGame.description || ''}</ReactMarkdown>
                </div>
              </div>

              {/* 时间信息和操作按钮 */}
              <div className="flex items-end justify-between">
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-neutral-400">{t('common.start')}:</span>
                    <span className="text-neutral-200 font-mono">
                      {new Date(currentGame.startTime).toLocaleString()}
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-neutral-400">{t('common.end')}:</span>
                    <span className="text-neutral-200 font-mono">{new Date(currentGame.endTime).toLocaleString()}</span>
                  </div>
                </div>

                <Button variant="primary" size="action" onClick={handleGameAction}>
                  {user?.user
                    ? currentGame.status === 'ended'
                      ? t('game.actions.results')
                      : t('game.actions.joinNow')
                    : currentGame.status === 'ended'
                      ? t('game.actions.ended')
                      : t('auth.login')}
                </Button>
              </div>
            </div>

            {/* 进度指示器 - 只有多于一张图片时显示 */}
            {games.length > 1 && (
              <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                {games.map((_, index) => (
                  <div
                    key={index}
                    className={`w-2 h-2 rounded-full transition-colors duration-200 
                                            ${index === currentIndex ? 'bg-geek-400' : 'bg-neutral-300'}`}
                  />
                ))}
              </div>
            )}
          </motion.div>
        </motion.div>
      </AnimatePresence>
    </div>
  );
}

export default GameSlider;
