import { motion } from 'motion/react';
import { IconEdit, IconTrash, IconPlus } from '@tabler/icons-react';
import { Button, Pagination, List } from '../../../common';
import { useTranslation } from 'react-i18next';

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
}) {
  const { t } = useTranslation();
  // 分类标签渲染
  const renderCategoryChip = (category) => {
    const categoryColors = {
      web: 'bg-blue-400/20 text-blue-400',
      crypto: 'bg-purple-400/20 text-purple-400',
      pwn: 'bg-red-400/20 text-red-400',
      reverse: 'bg-green-400/20 text-green-400',
      misc: 'bg-yellow-400/20 text-yellow-400',
    };

    const colorClass = categoryColors[category.toLowerCase()] || 'bg-neutral-400/20 text-neutral-400';

    return <span className={`px-2 py-1 rounded-full text-xs font-mono ${colorClass}`}>{category}</span>;
  };

  // 类型标签渲染
  const renderTypeChip = (type) => {
    const typeMap = {
      static: { label: t('admin.challenge.types.static'), class: 'bg-geek-400/20 text-geek-400' },
      question: { label: t('admin.challenge.types.question'), class: 'bg-green-400/20 text-green-400' },
      dynamic: { label: t('admin.challenge.types.dynamic'), class: 'bg-orange-400/20 text-orange-400' },
      pods: { label: t('admin.challenge.types.pods'), class: 'bg-cyan-400/20 text-cyan-400' },
    };

    const typeInfo = typeMap[type] || { label: type, class: 'bg-neutral-400/20 text-neutral-400' };

    return <span className={`px-2 py-1 rounded-full text-xs font-mono ${typeInfo.class}`}>{typeInfo.label}</span>;
  };

  // 自定义标签渲染
  const renderTags = (tags) => {
    if (!tags || tags.length === 0)
      return <span className="text-neutral-500 font-mono text-sm">{t('common.none')}</span>;

    return (
      <div className="flex flex-wrap gap-1">
        {tags.map((tag, index) => (
          <span
            key={index}
            className="px-2 py-0.5 bg-black/30 border border-geek-600/30 rounded-full text-xs font-mono text-geek-300"
          >
            {tag}
          </span>
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
          <span
            key={index}
            className="px-2 py-0.5 bg-black/30 border border-yellow-400/30 rounded text-xs font-mono text-yellow-300"
            title={hint}
          >
            {t('admin.contests.challenges.hintLabel', { index: index + 1 })}
          </span>
        ))}
      </div>
    );
  };

  const columns = [
    { key: 'name', label: t('admin.contests.challenges.table.name'), width: '28%' },
    { key: 'category', label: t('admin.contests.challenges.table.category'), width: '10%' },
    { key: 'type', label: t('admin.contests.challenges.table.type'), width: '10%' },
    { key: 'metrics', label: t('admin.contests.challenges.table.metrics'), width: '16%' },
    { key: 'tags', label: t('admin.contests.challenges.table.tags'), width: '16%' },
    { key: 'hints', label: t('admin.contests.challenges.table.hints'), width: '12%' },
    { key: 'actions', label: t('admin.contests.challenges.table.actions'), width: '8%' },
  ];

  const renderCell = (challenge, column) => {
    switch (column.key) {
      case 'name':
        return (
          <div className="flex flex-col gap-1 min-w-0">
            <div className="flex items-center gap-2 min-w-0">
              <span className="text-neutral-50 font-mono truncate">{challenge.name}</span>
              {challenge.hidden && (
                <span className="px-2 py-0.5 rounded-full text-xs font-mono bg-red-400/20 text-red-400">
                  {t('admin.contests.challenges.hidden')}
                </span>
              )}
            </div>
            <span className="text-xs text-neutral-400 line-clamp-1 break-all">{challenge.description || '-'}</span>
          </div>
        );
      case 'category':
        return challenge.category ? (
          renderCategoryChip(challenge.category)
        ) : (
          <span className="text-neutral-500">-</span>
        );
      case 'type':
        return renderTypeChip(challenge.type);
      case 'metrics':
        return (
          <div className="flex flex-col gap-1 text-xs font-mono text-neutral-400">
            <span>
              {t('admin.contests.challenges.attemptsLabel')} {challenge.attempt || 0}
            </span>
            <span>{t('common.points', { count: challenge.score || 0 })}</span>
            <span>{t('admin.contests.challenges.solves', { count: challenge.solvers || 0 })}</span>
          </div>
        );
      case 'tags':
        return <div className="whitespace-normal">{renderTags(challenge.tags)}</div>;
      case 'hints':
        return <div className="whitespace-normal">{renderHints(challenge.hints)}</div>;
      case 'actions':
        return (
          <div className="flex gap-2">
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
                className={`pb-3 px-2 relative font-mono text-sm ${filterType === 'question' ? 'text-geek-400' : 'text-neutral-400'}`}
                onClick={() => onFilterTypeChange('question')}
              >
                {t('admin.challenge.types.question')}
                {filterType === 'question' && (
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
