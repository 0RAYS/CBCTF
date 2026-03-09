/**
 * 赛题展示面板组件
 * @param {Object} props
 * @param {Array<string>} props.categories - 题目分类列表
 * @param {string} props.selectedCategory - 当前选中的分类 ('ALL' 或具体分类)
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
import { Button, Pagination, Card, Avatar } from '../../../../components/common';
import { useTranslation } from 'react-i18next';

function ChallengeBoard({
  categories,
  selectedCategory,
  onCategoryChange,
  challenges,
  onChallengeClick,
  teamInfo,
  totalCount = 0,
  currentPage = 1,
  pageSize = 12,
  onPageChange,
}) {
  const { t } = useTranslation();

  return (
    <Card variant="default" padding="lg" animate className="">
      {/* 分类和团队信息 */}
      <div className="flex justify-between items-center mb-8">
        {/* 分类标签 */}
        <div className="flex gap-4">
          {categories.map((category) => (
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

        {/* 团队信息 */}
        <div className="flex items-center gap-4 min-w-0">
          <div className="flex -space-x-2 flex-shrink-0">
            {teamInfo.members.map((member, i) => (
              <Avatar key={i} src={member.picture} name={member.name} size="xs" className="border-2 border-black" />
            ))}
          </div>
          <span className="text-neutral-400 font-mono truncate max-w-[160px]" title={teamInfo.name}>
            {teamInfo.name}
          </span>
        </div>
      </div>

      {/* 赛题列表 */}
      <div className="grid grid-cols-2 gap-4">
        {challenges.map((challenge) => (
          <motion.div
            key={challenge.id}
            className={`p-4 border rounded-md transition-colors duration-200 cursor-pointer backdrop-blur-none
                ${
                  challenge.solved
                    ? 'border-geek-400/50 bg-geek-400/5 hover:bg-geek-400/10'
                    : 'border-neutral-300/30 bg-black/30 hover:bg-black/50'
                }`}
            whileHover={{ y: -2 }}
            onClick={() => onChallengeClick(challenge)}
          >
            {/* 标题栏 */}
            <div className="flex items-center justify-between mb-3 min-w-0">
              <div className="flex items-center gap-3 min-w-0 flex-1">
                <span
                  className="text-geek-400 font-mono flex-shrink-0 truncate max-w-[80px]"
                  title={challenge.category}
                >
                  {challenge.category}
                </span>
                <h3 className="text-neutral-50 font-mono truncate min-w-0" title={challenge.title}>
                  {challenge.title}
                </h3>
                <span className="text-yellow-400 font-mono text-sm flex-shrink-0">
                  {t('common.points', { count: challenge.score })}
                </span>
              </div>
              <div className="flex items-center gap-2 flex-shrink-0 ml-2">
                <span className="text-neutral-400 text-sm">{t('game.challengeBoard.solves')}</span>
                <span className="text-neutral-50 font-mono">{challenge.solves}</span>
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
                      border backdrop-blur-none
                      truncate max-w-[96px]
                      ${
                        challenge.solved
                          ? 'border-geek-400/30 bg-geek-400/20 text-neutral-50'
                          : 'border-neutral-600 bg-neutral-900 text-neutral-300'
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
