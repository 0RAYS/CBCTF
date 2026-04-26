import { motion } from 'motion/react';
import { IconEdit, IconTrash, IconPlus, IconSearch } from '@tabler/icons-react';
import { Button, Pagination, List, Chip, Input } from '../../../common';
import { useTranslation } from 'react-i18next';
import { getChallengeCategoryChipClass, getChallengeTypeChipClass } from '../../../../config/challengeChips';

/**
 * 赛事内部赛题管理组件
 * @param {Object} props
 * @param {Array} props.challenges - 赛题列表
 * @param {number} props.totalCount - 赛题总数
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {Array} props.categories - 分类列表
 * @param {string} props.filterCategory - 当前选中分类
 * @param {string} props.filterType - 当前选中类型
 * @param {Function} props.onPageChange - 页码改变回调
 * @param {Function} props.onAddChallenge - 添加赛题回调
 * @param {Function} props.onEditChallenge - 编辑赛题回调
 * @param {Function} props.onDeleteChallenge - 删除赛题回调
 * @param {Function} props.onFilterTypeChange - 过滤类型改变回调
 * @param {Function} props.onFilterCategoryChange - 过滤分类改变回调
 */
function AdminContestChallenges({
  challenges = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 10,
  categories = [],
  filterCategory = 'all',
  filterType = 'all',
  onPageChange,
  onAddChallenge,
  onEditChallenge,
  onDeleteChallenge,
  onFilterTypeChange,
  onFilterCategoryChange,
  nameQuery = '',
  onNameChange,
}) {
  const { t } = useTranslation();
  // 自定义标签渲染
  const renderTags = (tags) => {
    if (!tags || tags.length === 0)
      return <span className="text-neutral-500 font-mono text-sm">{t('common.none')}</span>;

    return (
      <div className="flex flex-wrap gap-1">
        {tags.map((tag, index) => (
          <Chip
            key={index}
            variant="tag"
            size="sm"
            label={tag.length > 4 ? tag.slice(0, 4) + '…' : tag}
            colorClass="border-geek-600/30 text-geek-300"
            className="rounded-full"
          />
        ))}
      </div>
    );
  };

  // 提示信息渲染
  const renderHints = (hints) => {
    if (!hints || hints.length === 0)
      return <span className="text-neutral-500 font-mono text-sm">{t('common.none')}</span>;

    return (
      <div className="flex flex-wrap gap-1">
        {hints.map((hint, index) => (
          <Chip
            key={index}
            variant="tag"
            size="sm"
            label={hint.length > 8 ? hint.slice(0, 8) + '…' : hint}
            colorClass="border-yellow-400/30 text-yellow-300"
            title={hint}
          />
        ))}
      </div>
    );
  };

  const columns = [
    {
      key: 'name',
      label: t('admin.contests.challenges.table.name'),
      width: '15%',
    },
    {
      key: 'category',
      label: t('admin.contests.challenges.table.category'),
      width: '10%',
    },
    {
      key: 'type',
      label: t('admin.contests.challenges.table.type'),
      width: '10%',
    },
    {
      key: 'metrics',
      label: t('admin.contests.challenges.table.metrics'),
      width: '10%',
    },
    {
      key: 'tags',
      label: t('admin.contests.challenges.table.tags'),
      width: '10%',
    },
    {
      key: 'hints',
      label: t('admin.contests.challenges.table.hints'),
      width: '10%',
    },
    {
      key: 'actions',
      label: t('admin.contests.challenges.table.actions'),
      width: '7%',
    },
  ];

  const renderCell = (challenge, column) => {
    switch (column.key) {
      case 'name':
        return (
          <div className="flex flex-col gap-1 min-w-0">
            <div className="flex items-center gap-2 min-w-0">
              <span className="text-neutral-50 font-mono truncate">{challenge.name}</span>
              {challenge.hidden && (
                <Chip size="sm" label={t('admin.contests.challenges.hidden')} colorClass="bg-red-400/20 text-red-400" />
              )}
            </div>
            <span className="text-xs text-neutral-400 line-clamp-1 break-all">{challenge.description || '-'}</span>
          </div>
        );
      case 'category':
        return challenge.category ? (
          <Chip label={challenge.category} colorClass={getChallengeCategoryChipClass(challenge.category)} />
        ) : (
          <span className="text-neutral-500">-</span>
        );
      case 'type': {
        const typeLabels = {
          static: t('admin.challenge.types.static'),
          dynamic: t('admin.challenge.types.dynamic'),
          pods: t('admin.challenge.types.pods'),
        };
        return (
          <Chip
            label={typeLabels[challenge.type] || challenge.type}
            colorClass={getChallengeTypeChipClass(challenge.type)}
          />
        );
      }
      case 'metrics':
        return (
          <div className="flex flex-col gap-1 text-xs font-mono text-neutral-400">
            <span>
              {t('admin.contests.challenges.attemptsLabel')} {challenge.attempt || 0}
            </span>
            <span>{t('common.points', { count: challenge.score || 0 })}</span>
            <span>
              {t('admin.contests.challenges.solves', {
                count: challenge.solvers || 0,
              })}
            </span>
          </div>
        );
      case 'tags':
        return <div className="whitespace-normal">{renderTags(challenge.tags)}</div>;
      case 'hints':
        return <div className="whitespace-normal">{renderHints(challenge.hints)}</div>;
      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!text-geek-400 hover:!text-geek-300"
              onClick={() => onEditChallenge(challenge)}
            >
              <IconEdit size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!text-red-400 hover:!text-red-300"
              onClick={() => onDeleteChallenge(challenge)}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        );
      default:
        return challenge[column.key];
    }
  };

  return (
    <div className="w-full mx-auto">
      <div className="flex justify-end items-center mb-8">
        <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={18} />} onClick={onAddChallenge}>
          {t('admin.challenge.actions.add')}
        </Button>
      </div>

      <div className="mb-6">
        <label className="block text-sm font-mono text-neutral-400 mb-2">
          {t('admin.contests.challenges.search.label')}
        </label>
        <Input
          type="search"
          value={nameQuery}
          placeholder={t('admin.contests.challenges.search.placeholder')}
          onChange={(e) => onNameChange?.(e.target.value)}
          icon={<IconSearch size={16} />}
        />
      </div>

      {/* 过滤器 */}
      <div className="mb-6 border-b border-neutral-700">
        <div className="flex justify-between items-center">
          <div className="flex items-center">
            <div className="flex gap-8">
              <Button
                variant="ghost"
                className={`pb-3 px-2 relative font-mono text-sm ${filterType === 'all' ? 'text-geek-400' : 'text-neutral-400'}`}
                onClick={() => onFilterTypeChange('all')}
              >
                {t('admin.challenge.filters.all')}
                {filterType === 'all' && (
                  <motion.div
                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-geek-400"
                    layoutId="filterTypeIndicator"
                  />
                )}
              </Button>
              <Button
                variant="ghost"
                className={`pb-3 px-2 relative font-mono text-sm ${filterType === 'static' ? 'text-geek-400' : 'text-neutral-400'}`}
                onClick={() => onFilterTypeChange('static')}
              >
                {t('admin.challenge.types.static')}
                {filterType === 'static' && (
                  <motion.div
                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-geek-400"
                    layoutId="filterTypeIndicator"
                  />
                )}
              </Button>
              <Button
                variant="ghost"
                className={`pb-3 px-2 relative font-mono text-sm ${filterType === 'dynamic' ? 'text-geek-400' : 'text-neutral-400'}`}
                onClick={() => onFilterTypeChange('dynamic')}
              >
                {t('admin.challenge.types.dynamic')}
                {filterType === 'dynamic' && (
                  <motion.div
                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-geek-400"
                    layoutId="filterTypeIndicator"
                  />
                )}
              </Button>
              <Button
                variant="ghost"
                className={`pb-3 px-2 relative font-mono text-sm ${filterType === 'pods' ? 'text-geek-400' : 'text-neutral-400'}`}
                onClick={() => onFilterTypeChange('pods')}
              >
                {t('admin.challenge.types.pods')}
                {filterType === 'pods' && (
                  <motion.div
                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-geek-400"
                    layoutId="filterTypeIndicator"
                  />
                )}
              </Button>
            </div>
          </div>
          <div className="w-48">
            <select
              value={filterCategory}
              onChange={(e) => onFilterCategoryChange(e.target.value)}
              className="select-custom select-custom-sm"
            >
              <option value="all">{t('admin.challenge.category.all')}</option>
              {categories.map((category) => (
                <option key={category} value={category}>
                  {category}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* 赛题列表 */}
      <List
        columns={columns}
        data={challenges}
        renderCell={renderCell}
        empty={challenges.length === 0}
        emptyContent={t('common.noData')}
        footer={
          totalCount > pageSize ? (
            <Pagination
              total={Math.ceil(totalCount / pageSize)}
              current={currentPage}
              onChange={onPageChange}
              showTotal={true}
              totalItems={totalCount}
            />
          ) : null
        }
      />
    </div>
  );
}

export default AdminContestChallenges;
