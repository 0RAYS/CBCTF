/**
 * 赛题展示面板组件
 * @param {Object} props
 * @param {Array<string>} props.categories - 题目分类列表
 * @param {string} props.selectedCategory - 当前选中的分类（空字符串表示不过滤）
 * @param {Function} props.onCategoryChange - 分类切换回调 (category: string) => void
 * @param {Array<Object>} props.challenges - 题目列表
 * @param {number} props.challenges[].id - 题目ID
 * @param {string} props.challenges[].category - 题目类别
 * @param {string} props.challenges[].title - 题目标题
 * @param {number} props.challenges[].score - 题目分值
 * @param {boolean} props.challenges[].isInitialized - 是否已初始化
 * @param {number} props.challenges[].solves - 解题人数
 * @param {boolean} props.challenges[].hasInstance - 是否有靶机
 * @param {boolean} props.challenges[].instanceRunning - 靶机是否运行中
 * @param {boolean} [props.challenges[].solved] - 是否已解决
 * @param {Function} props.onChallengeClick - 题目点击回调 (challenge: Object) => void
 * @param {Object} props.teamInfo - 团队信息
 * @param {string} props.teamInfo.name - 团队名称
 * @param {Array<{picture: string, name: string}>} props.teamInfo.members - 团队成员头像列表
 */

import { motion } from 'motion/react';
import { Button, Pagination, Card, Avatar, EmptyState, ChallengeSkeleton } from '../../../../components/common';
import { useTranslation } from 'react-i18next';
import { EASE_T2, staggerDelay } from '../../../../config/motion';

function ChallengeBoard({
  categories,
  selectedCategory,
  onCategoryChange,
  unsolvedOnly = false,
  onSolvedFilterChange,
  challenges,
  onChallengeClick,
  teamInfo,
  totalCount = 0,
  currentPage = 1,
  pageSize = 12,
  onPageChange,
  isLoading = false,
}) {
  const { t } = useTranslation();
  const normalizedCategories = Array.isArray(categories) ? categories.filter(Boolean) : [];
  const members = Array.isArray(teamInfo?.members) ? teamInfo.members : [];
  const teamName = teamInfo?.name || '-';

  return (
    <Card variant="default" padding="lg" animate className="">
      {/* 分类和团队信息 */}
      <div className="flex flex-col gap-4 mb-8 sm:flex-row sm:justify-between sm:items-center">
        <div className="flex flex-col gap-3">
          {/* 分类标签 */}
          <div className="flex flex-wrap gap-2 md:gap-4">
            {normalizedCategories.map((category) => (
              <Button
                key={category}
                variant={selectedCategory === category ? 'primary' : 'ghost'}
                size="sm"
                className={`min-w-0 px-4 py-1 ${
                  selectedCategory === category ? '' : 'text-neutral-400 hover:text-neutral-200'
                }`}
                onClick={() => onCategoryChange(category)}
              >
                {category}
              </Button>
            ))}
          </div>
        </div>

        {/* 团队信息 */}
        <div className="flex flex-col gap-3 min-w-0 sm:items-end">
          <div className="flex items-center gap-4 min-w-0">
            <div className="flex -space-x-2 flex-shrink-0">
              {members.map((member, i) => (
                <Avatar key={i} src={member.picture} name={member.name} size="xs" className="border-2 border-black" />
              ))}
            </div>
            <span className="text-neutral-400 font-mono truncate max-w-[160px]" title={teamName}>
              {teamName}
            </span>
          </div>
          <label className="flex items-center gap-2 text-sm font-mono text-neutral-300 whitespace-nowrap cursor-pointer">
            <input
              type="checkbox"
              checked={unsolvedOnly}
              onChange={() => onSolvedFilterChange?.()}
              className="h-4 w-4 rounded border border-neutral-500/60 bg-black/20 text-geek-400 focus:ring-geek-400/70"
            />
            <span>{t('game.challengeBoard.filters.hideSolved')}</span>
          </label>
        </div>
      </div>

      {/* 赛题列表 */}
      {isLoading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {Array.from({ length: pageSize }).map((_, i) => (
            <ChallengeSkeleton key={i} />
          ))}
        </div>
      ) : challenges.length === 0 ? (
        <EmptyState title={t('game.noChallenges')} description={t('game.noChallengesDescription')} className="py-12" />
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {challenges.map((challenge, index) => (
            <motion.div
              key={challenge.id}
              className={`p-4 border rounded-md transition-colors duration-200 cursor-pointer
                ${
                  challenge.solved
                    ? 'border-geek-400/40 bg-geek-400/8 hover:bg-geek-400/12'
                    : 'border-neutral-600/50 bg-neutral-800/40 hover:bg-neutral-800/60'
                }`}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: staggerDelay(index), ease: EASE_T2, duration: 0.22 }}
              whileHover={{ y: -2 }}
              onClick={() => onChallengeClick(challenge)}
            >
              {/* 标题栏 */}
              <div className="mb-3 min-w-0">
                {/* Category — small label tier */}
                <div className="flex items-center justify-between mb-1.5">
                  <span
                    className="text-[10px] font-mono tracking-[0.2em] uppercase text-geek-400/80 flex-shrink-0 truncate max-w-[100px]"
                    title={challenge.category}
                  >
                    {challenge.category}
                  </span>
                  <div className="flex items-center gap-2 flex-shrink-0 ml-2">
                    <span className="text-neutral-500 text-xs">{t('game.challengeBoard.solves')}</span>
                    <span className="text-neutral-300 font-mono text-xs tabular-nums">{challenge.solves}</span>
                  </div>
                </div>
                {/* Title — primary size tier */}
                <div className="flex items-end justify-between gap-2">
                  <h3
                    className="text-neutral-50 font-mono text-base leading-snug truncate min-w-0 font-medium"
                    title={challenge.title}
                  >
                    {challenge.title}
                  </h3>
                  <span className="text-yellow-400 font-mono text-sm flex-shrink-0 tabular-nums">
                    {t('common.points', { count: challenge.score })}
                  </span>
                </div>
              </div>

              {/* 标签和状态区域 */}
              <div className="flex items-center justify-between">
                {/* 标签列表 */}
                <div className="flex items-center gap-2 min-w-0 flex-1 overflow-hidden">
                  {challenge.tags &&
                    challenge.tags.map((tag, index) => (
                      <span
                        key={index}
                        className={`
                      px-2 py-0.5 rounded
                      text-xs font-mono
                      border
                      truncate max-w-[96px]
                      ${
                        challenge.solved
                          ? 'border-geek-400/30 bg-geek-400/15 text-geek-300'
                          : 'border-neutral-600/60 bg-neutral-700/40 text-neutral-400'
                      }
                    `}
                        title={tag}
                      >
                        {tag}
                      </span>
                    ))}
                </div>

                {/* 状态指示器 */}
                <div className="flex items-center gap-3">
                  {!challenge.isInitialized && (
                    <div className="flex items-center gap-1.5 px-2 py-1 bg-black/50 border border-yellow-400/30 rounded">
                      <span className="w-1.5 h-1.5 rounded-full bg-yellow-400 animate-pulse"></span>
                      <span className="text-yellow-400 text-xs font-mono">
                        {t('game.challengeBoard.status.notInitialized')}
                      </span>
                    </div>
                  )}

                  {challenge.hasInstance && challenge.instanceRunning && (
                    <div className="flex items-center gap-1.5 px-2 py-1 bg-black/50 border border-green-400/30 rounded">
                      <span className="w-1.5 h-1.5 rounded-full bg-green-400 animate-pulse"></span>
                      <span className="text-green-400 text-xs font-mono">
                        {t('game.challengeBoard.status.instanceRunning')}
                      </span>
                    </div>
                  )}

                  {challenge.solved && (
                    <div className="flex items-center gap-1.5 px-2 py-1 bg-black/50 border border-geek-400/30 rounded">
                      <span className="w-1.5 h-1.5 rounded-full bg-geek-400"></span>
                      <span className="text-geek-400 text-xs font-mono">{t('common.solved')}</span>
                    </div>
                  )}
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      )}

      {/* 分页 */}
      {totalCount > 0 && (
        <div className="mt-6 flex justify-center">
          <Pagination
            total={Math.ceil(totalCount / pageSize)}
            current={currentPage}
            onChange={onPageChange}
            showTotal={true}
            totalItems={totalCount}
          />
        </div>
      )}
    </Card>
  );
}

export default ChallengeBoard;
