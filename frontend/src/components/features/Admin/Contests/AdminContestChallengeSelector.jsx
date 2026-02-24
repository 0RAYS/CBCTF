import { motion } from 'motion/react';
import { IconX, IconSearch } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Pagination, Card, EmptyState } from '../../../../components/common';
import { getChallengeCategoryChipClass, getChallengeTypeChipClass } from '../../../../config/challengeChips';

/**
 * 赛事赛题选择弹窗组件
 * @param {Object} props
 * @param {boolean} props.isOpen - 是否显示弹窗
 * @param {Array} props.challenges - 赛题库列表
 * @param {Array} props.selectedChallenges - 已选中的赛题
 * @param {Array} props.categories - 分类列表
 * @param {number} props.totalCount - 赛题总数
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {Function} props.onClose - 关闭弹窗回调
 * @param {Function} props.onSearch - 搜索赛题回调
 * @param {Function} props.onSelect - 选择赛题回调
 * @param {Function} props.onConfirm - 确认选择回调
 * @param {Function} props.onPageChange - 页码变更回调
 * @param {Function} props.onFilterCategoryChange - 过滤分类改变回调
 * @param {Function} props.onFilterTypeChange - 过滤类型改变回调
 */
function AdminContestChallengeSelector({
  isOpen = false,
  challenges = [],
  selectedChallenges = [],
  categories = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 10,
  onClose,
  onSearch,
  onSelect,
  onConfirm,
  onPageChange,
  onFilterCategoryChange,
  onFilterTypeChange,
}) {
  const { t } = useTranslation();

  // 统一输入框样式
  const inputBaseClass =
    'w-full bg-black/20 border border-neutral-300/30 rounded-md p-3 text-neutral-50 font-mono focus:border-geek-400 focus:outline-none transition-colors duration-200';

  // 分类标签渲染
  const renderCategoryChip = (category) => {
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-mono ${getChallengeCategoryChipClass(category)}`}>
        {category}
      </span>
    );
  };

  // 类型标签渲染
  const renderTypeChip = (type) => {
    const typeMap = {
      static: { label: t('admin.contests.challengeSelector.types.static') },
      question: {
        label: t('admin.contests.challengeSelector.types.question'),
      },
      dynamic: {
        label: t('admin.contests.challengeSelector.types.dynamic'),
      },
      pods: { label: t('admin.contests.challengeSelector.types.pods') },
    };

    const typeInfo = typeMap[type] || { label: type };

    return (
      <span className={`px-2 py-1 rounded-full text-xs font-mono ${getChallengeTypeChipClass(type)}`}>
        {typeInfo.label}
      </span>
    );
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <motion.div
        className="w-full max-w-5xl h-[80vh] bg-neutral-900 border border-neutral-300 rounded-md overflow-hidden flex flex-col"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.2 }}
      >
        {/* 标题栏 */}
        <div className="flex justify-between items-center p-4 border-b border-neutral-700">
          <h2 className="text-xl font-mono text-neutral-50">{t('admin.contests.challengeSelector.title')}</h2>
          <Button variant="ghost" size="icon" className="!text-neutral-400 hover:!text-neutral-300" onClick={onClose}>
            <IconX size={18} />
          </Button>
        </div>

        {/* 搜索和过滤 */}
        <div className="p-4 border-b border-neutral-700">
          <div className="flex gap-4">
            <div className="flex-1 relative">
              <IconSearch size={18} className="absolute left-3 top-1/2 transform -translate-y-1/2 text-neutral-400" />
              <input
                type="text"
                className={`${inputBaseClass} pl-10`}
                placeholder={t('admin.contests.challengeSelector.search.placeholder')}
                onChange={(e) => onSearch(e.target.value)}
              />
            </div>

            <div className="w-48">
              <select
                onChange={(e) => onFilterCategoryChange(e.target.value)}
                className="select-custom select-custom-lg"
              >
                <option value="all">{t('admin.contests.challengeSelector.filters.categoryAll')}</option>
                {categories.map((category) => (
                  <option key={category} value={category}>
                    {category}
                  </option>
                ))}
              </select>
            </div>

            <div className="w-48">
              <select onChange={(e) => onFilterTypeChange(e.target.value)} className="select-custom select-custom-lg">
                <option value="all">{t('admin.contests.challengeSelector.filters.typeAll')}</option>
                <option value="question">{t('admin.contests.challengeSelector.types.question')}</option>
                <option value="static">{t('admin.contests.challengeSelector.types.static')}</option>
                <option value="dynamic">{t('admin.contests.challengeSelector.types.dynamic')}</option>
                <option value="pods">{t('admin.contests.challengeSelector.types.pods')}</option>
              </select>
            </div>
          </div>
        </div>

        {/* 赛题列表 */}
        <div className="flex-1 overflow-y-auto p-4">
          {challenges.length === 0 ? (
            <Card variant="default" padding="md" className="flex justify-center items-center h-full">
              <EmptyState title={t('admin.contests.challengeSelector.empty.noMatch')} />
            </Card>
          ) : (
            <div className="space-y-3">
              {challenges.map((challenge) => {
                const isSelected = selectedChallenges.some((c) => c.id === challenge.id);

                return (
                  <motion.div
                    key={challenge.id || challenge.name}
                    className={`border rounded-md bg-black/30 backdrop-blur-[2px] overflow-hidden transition-colors duration-200 ${
                      isSelected ? 'border-geek-400' : 'border-neutral-300/30'
                    }`}
                    whileHover={{ y: -1, boxShadow: '0 4px 15px rgba(0,0,0,0.2)' }}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                  >
                    <div className="p-3 flex items-start gap-3">
                      <div
                        className={`w-5 h-5 border-2 cursor-pointer transition-all duration-200 flex items-center justify-center flex-shrink-0 mt-1 ${
                          isSelected
                            ? 'bg-geek-400 border-geek-400 text-white'
                            : 'border-neutral-400 text-transparent hover:border-geek-400 hover:bg-geek-400/10'
                        }`}
                        onClick={() => onSelect(challenge)}
                      >
                        {isSelected && (
                          <svg
                            width="12"
                            height="12"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="3"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                          >
                            <polyline points="20,6 9,17 4,12"></polyline>
                          </svg>
                        )}
                      </div>

                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <h3 className="text-lg font-mono text-neutral-50">{challenge.name}</h3>
                          <div className="flex gap-2">
                            {renderCategoryChip(challenge.category)}
                            {renderTypeChip(challenge.type)}
                          </div>
                        </div>

                        <p className="text-neutral-300 text-sm font-mono line-clamp-2">{challenge.description}</p>
                      </div>
                    </div>
                  </motion.div>
                );
              })}
            </div>
          )}

          {/* 分页 */}
          {totalCount > pageSize && (
            <div className="mt-6">
              <Pagination
                total={Math.ceil(totalCount / pageSize)}
                current={currentPage}
                pageSize={pageSize}
                onChange={onPageChange}
                showTotal={true}
                totalItems={totalCount}
              />
            </div>
          )}
        </div>

        {/* 已选赛题计数和确认按钮 */}
        <div className="flex justify-between items-center p-4 border-t border-neutral-700">
          <div className="text-neutral-300 font-mono">
            {t('admin.contests.challengeSelector.selectedPrefix')}
            <span className="text-geek-400 font-bold">{selectedChallenges.length}</span>
            {t('admin.contests.challengeSelector.selectedSuffix')}
          </div>
          <div className="flex gap-4">
            <Button variant="ghost" onClick={onClose}>
              {t('common.cancel')}
            </Button>
            <Button variant="primary" onClick={onConfirm} disabled={selectedChallenges.length === 0}>
              {t('admin.contests.challengeSelector.actions.confirmAdd')}
            </Button>
          </div>
        </div>
      </motion.div>
    </div>
  );
}

export default AdminContestChallengeSelector;
