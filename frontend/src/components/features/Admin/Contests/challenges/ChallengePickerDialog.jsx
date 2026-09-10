import { motion } from 'motion/react';
import { IconSearch } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Pagination, Card, EmptyState, Chip, Modal, Spinner } from '../../../../common';
import Checkbox from '../../../../common/Checkbox';
import { getChallengeCategoryChipClass, getChallengeTypeChipClass } from '../../../../../config/challengeChips';

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
  loading = false,
  saving = false,
  searchQuery = '',
  descQuery = '',
  category = 'all',
  type = 'all',
  onClose,
  onSearch,
  onDescSearch,
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

  // 分类标签 / 类型标签由 Chip 组件统一渲染

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={t('admin.contests.challengeSelector.title')}
      size="xl"
      footer={
        <>
          <span className="mr-auto text-neutral-300 font-mono">
            {t('admin.contests.challengeSelector.selectedPrefix')}{' '}
            <span className="text-geek-400">{selectedChallenges.length}</span>{' '}
            {t('admin.contests.challengeSelector.selectedSuffix')}
          </span>
          <Button size="sm" variant="ghost" onClick={onClose} disabled={saving}>
            {t('common.cancel')}
          </Button>
          <Button size="sm" onClick={onConfirm} disabled={saving || selectedChallenges.length === 0}>
            {t('admin.contests.challengeSelector.actions.confirmAdd')}
          </Button>
        </>
      }
    >
      {/* 搜索和过滤 */}
      <div className="p-4 border-b border-neutral-700">
        <div className="flex flex-col sm:flex-row gap-4 mb-3">
          <div className="flex-1 relative">
            <IconSearch size={18} className="absolute left-3 top-1/2 transform -translate-y-1/2 text-neutral-400" />
            <input
              type="text"
              value={searchQuery}
              className={`${inputBaseClass} pl-10`}
              placeholder={t('admin.contests.challengeSelector.search.placeholder')}
              onChange={(e) => onSearch(e.target.value)}
            />
          </div>
          <div className="flex-1 relative">
            <IconSearch size={18} className="absolute left-3 top-1/2 transform -translate-y-1/2 text-neutral-400" />
            <input
              type="text"
              value={descQuery}
              className={`${inputBaseClass} pl-10`}
              placeholder={t('admin.contests.challengeSelector.search.descPlaceholder')}
              onChange={(e) => onDescSearch?.(e.target.value)}
            />
          </div>
        </div>
        <div className="flex flex-wrap gap-4">
          <div className="w-48">
            <select
              value={category}
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
            <select
              value={type}
              onChange={(e) => onFilterTypeChange(e.target.value)}
              className="select-custom select-custom-lg"
            >
              <option value="all">{t('admin.contests.challengeSelector.filters.typeAll')}</option>
              <option value="static">{t('admin.contests.challengeSelector.types.static')}</option>
              <option value="dynamic">{t('admin.contests.challengeSelector.types.dynamic')}</option>
              <option value="pods">{t('admin.contests.challengeSelector.types.pods')}</option>
            </select>
          </div>
        </div>
      </div>

      {/* 赛题列表 */}
      <div className="flex-1 overflow-y-auto p-4">
        {loading ? (
          <div className="flex justify-center p-8">
            <Spinner />
          </div>
        ) : challenges.length === 0 ? (
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
                  className={`border rounded-md bg-neutral-900 overflow-hidden transition-colors duration-200 ${
                    isSelected ? 'border-geek-400' : 'border-neutral-300/30'
                  }`}
                  whileHover={{
                    y: -1,
                    boxShadow: '0 4px 15px rgba(0,0,0,0.2)',
                  }}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                >
                  <div className="p-3 flex items-start gap-3">
                    <Checkbox
                      checked={isSelected}
                      onChange={() => onSelect(challenge)}
                      aria-label={challenge.name}
                      disabled={saving}
                    />

                    <div className="flex-1 min-w-0">
                      <div className="flex flex-wrap items-center gap-3 mb-2">
                        <h3 className="text-lg font-mono text-neutral-50">{challenge.name}</h3>
                        <div className="flex gap-2">
                          <Chip
                            label={challenge.category}
                            colorClass={getChallengeCategoryChipClass(challenge.category)}
                          />
                          <Chip
                            label={t(`admin.contests.challengeSelector.types.${challenge.type}`, {
                              defaultValue: challenge.type,
                            })}
                            colorClass={getChallengeTypeChipClass(challenge.type)}
                          />
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
    </Modal>
  );
}

export default AdminContestChallengeSelector;
