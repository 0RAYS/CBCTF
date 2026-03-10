import { motion } from 'motion/react';
import { IconEdit, IconTrash, IconPlus, IconUpload, IconDownload, IconSearch, IconFlask } from '@tabler/icons-react';
import { Button, Pagination, Input, Spinner, Chip } from '../../../components/common';
import { List } from '../../common';
import { useTranslation } from 'react-i18next';
import { getChallengeCategoryChipClass, getChallengeTypeChipClass } from '../../../config/challengeChips';

/**
 * 题目管理组件
 * @param {Object} props
 * @param {Array} props.challenges - 题目列表
 * @param {number} props.totalCount - 题目总数
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {Array} props.categories - 分类列表
 * @param {string} props.filterCategory - 当前选中分类
 * @param {string} props.filterType - 当前选中类型
 * @param {Function} props.onPageChange - 页码改变回调
 * @param {Function} props.onAddChallenge - 添加题目回调
 * @param {Function} props.onEditChallenge - 编辑题目回调
 * @param {Function} props.onDeleteChallenge - 删除题目回调
 * @param {Function} props.onUploadAttachment - 上传附件回调
 * @param {Function} props.onDownloadAttachment - 下载附件回调
 * @param {Function} props.onTestChallenge - 测试题目回调
 * @param {Function} props.onFilterTypeChange - 过滤类型改变回调
 * @param {Function} props.onFilterCategoryChange - 过滤分类改变回调
 * @param {string} props.searchQuery - 搜索查询
 * @param {boolean} props.searchLoading - 搜索加载状态
 * @param {boolean} props.isSearchMode - 是否处于搜索模式
 * @param {Function} props.onSearchChange - 搜索查询变化回调
 */
function AdminChallenge({
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
  onUploadAttachment,
  onDownloadAttachment,
  onTestChallenge,
  onFilterTypeChange,
  onFilterCategoryChange,
  nameQuery = '',
  descQuery = '',
  searchLoading = false,
  isSearchMode = false,
  onNameChange,
  onDescChange,
}) {
  const { t } = useTranslation();

  // Flag 展示
  const renderFlags = (flags) => {
    if (!flags || flags.length === 0)
      return <span className="text-neutral-500 font-mono text-sm">{t('common.none')}</span>;

    return (
      <div className="flex flex-wrap gap-1">
        {flags.map((flag, index) => (
          <Chip
            key={index}
            variant="tag"
            label={flag.value.length > 15 ? `${flag.value.substring(0, 15)}...` : flag.value}
            title={flag.value}
          />
        ))}
      </div>
    );
  };

  const columns = [
    { key: 'name', label: t('admin.challenge.table.name'), width: '24%' },
    { key: 'category', label: t('admin.challenge.table.category'), width: '12%' },
    { key: 'type', label: t('admin.challenge.table.type'), width: '12%' },
    { key: 'flags', label: t('admin.challenge.table.flags'), width: '22%' },
    { key: 'file', label: t('admin.challenge.table.file'), width: '14%' },
    { key: 'actions', label: t('admin.challenge.table.actions'), width: '16%' },
  ];

  const renderCell = (challenge, column) => {
    switch (column.key) {
      case 'name':
        return (
          <div className="flex flex-col gap-1 min-w-0 whitespace-normal">
            <span className="text-neutral-50 font-mono truncate">{challenge.name}</span>
            <span className="text-xs text-neutral-400 line-clamp-2">{challenge.description || '-'}</span>
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
          question: t('admin.challenge.types.question'),
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
      case 'flags':
        if (challenge.type === 'pods' || challenge.type === 'question') {
          return <span className="text-neutral-500 font-mono text-sm">{t('common.notAvailable')}</span>;
        }
        return <div className="whitespace-normal">{renderFlags(challenge.flags)}</div>;
      case 'file':
        return challenge.file ? (
          <Chip variant="tag" label={challenge.file} title={challenge.file} className="whitespace-normal" />
        ) : (
          <span className="text-neutral-500 font-mono text-sm">{t('common.none')}</span>
        );
      case 'actions':
        return (
          <div className="flex flex-wrap gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-neutral-400 hover:!text-neutral-300"
              onClick={() => onUploadAttachment(challenge)}
            >
              <IconUpload size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-neutral-400 hover:!text-neutral-300"
              disabled={!challenge.file}
              onClick={() => onDownloadAttachment(challenge)}
            >
              <IconDownload size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-orange-400 hover:!text-orange-300"
              onClick={() => onTestChallenge(challenge)}
              title={t('admin.challenge.actions.test')}
            >
              <IconFlask size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-geek-400 hover:!text-geek-300"
              onClick={() => onEditChallenge(challenge)}
            >
              <IconEdit size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-red-400 hover:!text-red-300"
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

  const paginationComponent = !isSearchMode && totalCount > pageSize && (
    <Pagination
      total={Math.ceil(totalCount / pageSize)}
      current={currentPage}
      onChange={onPageChange}
      showTotal
      totalItems={totalCount}
      showJumpTo
    />
  );

  return (
    <div className="w-full mx-auto">
      <div className="flex justify-end items-center mb-8">
        <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onAddChallenge}>
          {t('admin.challenge.actions.add')}
        </Button>
      </div>

      {/* 搜索框 */}
      <div className="mb-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">{t('admin.challenge.search.label')}</label>
            <Input
              type="text"
              value={nameQuery}
              placeholder={t('admin.challenge.search.placeholder')}
              onChange={(e) => onNameChange?.(e.target.value)}
              icon={<IconSearch size={16} />}
              iconRight={searchLoading && <Spinner size="sm" />}
            />
          </div>
          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">{t('admin.challenge.search.descLabel')}</label>
            <Input
              type="text"
              value={descQuery}
              placeholder={t('admin.challenge.search.descPlaceholder')}
              onChange={(e) => onDescChange?.(e.target.value)}
              icon={<IconSearch size={16} />}
            />
          </div>
        </div>
      </div>

      {/* 过滤器 */}
      <div className="mb-6 border-b border-neutral-700">
        <div className="flex justify-between items-center">
          <div className="flex items-center">
            <div className="flex gap-8">
              {[
                { id: 'all', label: t('admin.challenge.filters.all') },
                { id: 'static', label: t('admin.challenge.filters.static') },
                { id: 'question', label: t('admin.challenge.filters.question') },
                { id: 'dynamic', label: t('admin.challenge.filters.dynamic') },
                { id: 'pods', label: t('admin.challenge.filters.pods') },
              ].map((type) => (
                <div key={type.id} className="relative">
                  <Button
                    variant={filterType === type.id ? 'primary' : 'ghost'}
                    size="sm"
                    className={`!px-2 !min-w-0 !border-0 pb-3 ${
                      filterType === type.id ? '!bg-transparent !text-geek-400' : '!text-neutral-400'
                    }`}
                    onClick={() => onFilterTypeChange(type.id)}
                  >
                    {type.label}
                  </Button>
                  {filterType === type.id && (
                    <motion.div
                      className="absolute bottom-0 left-0 right-0 h-0.5 bg-geek-400"
                      layoutId="filterTypeIndicator"
                    />
                  )}
                </div>
              ))}
            </div>
          </div>
          <div className="flex items-center gap-4">
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
      </div>

      {/* 题目列表 */}
      <List
        columns={columns}
        data={challenges}
        renderCell={renderCell}
        empty={challenges.length === 0}
        emptyContent={
          isSearchMode ? (
            <div className="flex flex-col items-center justify-center space-y-2">
              <svg className="w-12 h-12 text-neutral-300/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={1.5}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
              <span className="font-mono text-neutral-400">{t('admin.challenge.empty.search')}</span>
            </div>
          ) : undefined
        }
        footer={paginationComponent}
      />
    </div>
  );
}

export default AdminChallenge;
