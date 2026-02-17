/**
 * 预热镜像管理组件
 * @param {Object} props
 * @param {Array} props.images - 镜像列表
 * @param {Array} props.selectedImages - 选中的镜像列表
 * @param {string} props.pullPolicy - 拉取策略
 * @param {boolean} props.loading - 是否正在加载
 * @param {boolean} props.submitting - 是否正在提交
 * @param {Function} props.onImageToggle - 镜像选择切换回调
 * @param {Function} props.onSelectAll - 全选/取消全选回调
 * @param {Function} props.onPullPolicyChange - 拉取策略改变回调
 * @param {Function} props.onWarmup - 执行预热回调
 * @param {Function} props.onRefresh - 刷新状态回调
 */

import { IconRefresh, IconDownload, IconServer } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Card, EmptyState, List } from '../../../common';

function AdminImagesWarmup({
  images = [],
  selectedImages = [],
  pullPolicy = 'IfNotPresent',
  loading = false,
  submitting = false,
  onImageToggle,
  onSelectAll,
  onPullPolicyChange,
  onWarmup,
  onRefresh,
}) {
  const { t } = useTranslation();

  const getStatusBadge = (status) => {
    if (status) {
      return (
        <span className="px-2 py-1 text-xs font-mono rounded border bg-green-400/20 text-green-400 border-green-400/30">
          {t('admin.contests.imagesWarmup.status.loaded')}
        </span>
      );
    }
    return (
      <span className="px-2 py-1 text-xs font-mono rounded border bg-neutral-400/20 text-neutral-400 border-neutral-400/30">
        {t('admin.contests.imagesWarmup.status.notLoaded')}
      </span>
    );
  };

  const pullPolicyOptions = [
    { value: 'Always', label: t('admin.contests.imagesWarmup.pullPolicy.always') },
    { value: 'IfNotPresent', label: t('admin.contests.imagesWarmup.pullPolicy.ifNotPresent') },
    { value: 'Never', label: t('admin.contests.imagesWarmup.pullPolicy.never') },
  ];

  const selectAllLabel = (
    <div className="flex items-center gap-2">
      <input
        type="checkbox"
        checked={selectedImages.length === images.length && images.length > 0}
        onChange={onSelectAll}
        className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
      />
      <span className="text-xs font-mono text-neutral-400">{t('admin.contests.imagesWarmup.actions.selectAll')}</span>
    </div>
  );

  const columns = [
    { key: 'select', label: selectAllLabel, width: '12%' },
    { key: 'image', label: t('admin.contests.imagesWarmup.table.image'), width: '38%' },
    { key: 'nodes', label: t('admin.contests.imagesWarmup.table.nodes'), width: '50%' },
  ];

  const renderCell = (imageObj, column) => {
    const imageName = Object.keys(imageObj)[0];
    const nodes = imageObj[imageName];
    const isSelected = selectedImages.includes(imageName);

    switch (column.key) {
      case 'select':
        return (
          <input
            type="checkbox"
            checked={isSelected}
            onChange={() => onImageToggle(imageName)}
            className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
          />
        );

      case 'image':
        return <span className="text-sm font-mono text-neutral-50 break-all">{imageName}</span>;

      case 'nodes':
        return (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-2">
            {nodes.map((node, nodeIndex) => (
              <div
                key={nodeIndex}
                className="flex items-center gap-2 p-2 bg-black/20 border border-neutral-300/20 rounded-md"
              >
                <span className="text-xs font-mono text-neutral-300 flex-1 truncate">{node.node}</span>
                {getStatusBadge(node.status)}
              </div>
            ))}
          </div>
        );

      default:
        return null;
    }
  };

  if (loading) {
    return (
      <div className="w-full mx-auto">
        <div className="flex justify-center items-center h-64">
          <div className="flex items-center gap-3 text-neutral-400">
            <div className="animate-spin">
              <IconRefresh size={20} />
            </div>
            <span className="font-mono">{t('common.loading')}</span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full mx-auto space-y-6">
      <div className="flex justify-end items-center">
        <Button variant="primary" size="sm" align="icon-left" icon={<IconRefresh size={16} />} onClick={onRefresh}>
          {t('common.refresh')}
        </Button>
      </div>

      <Card variant="default" padding="md" animate>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-mono text-neutral-50">{t('admin.contests.imagesWarmup.control.title')}</h2>
          <div className="flex items-center gap-2 text-neutral-400 font-mono text-sm">
            <IconServer size={16} />
            <span>{t('admin.contests.imagesWarmup.control.subtitle')}</span>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.imagesWarmup.labels.pullPolicy')}
            </label>
            <select
              value={pullPolicy}
              onChange={(e) => onPullPolicyChange(e.target.value)}
              className="select-custom select-custom-md"
            >
              {pullPolicyOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-end">
            <Button variant="outline" onClick={onSelectAll} className="w-full">
              {selectedImages.length === images.length
                ? t('admin.contests.imagesWarmup.actions.deselectAll')
                : t('admin.contests.imagesWarmup.actions.selectAll')}
            </Button>
          </div>

          <div className="flex items-end">
            <Button
              variant="primary"
              onClick={onWarmup}
              loading={submitting}
              disabled={selectedImages.length === 0}
              className="w-full"
              align="icon-left"
              icon={<IconDownload size={18} />}
            >
              {t('admin.contests.imagesWarmup.actions.warmup', { count: selectedImages.length })}
            </Button>
          </div>
        </div>
      </Card>

      <Card variant="default" padding="none" animate className="overflow-hidden">
        <div className="px-6 py-4 border-b border-neutral-300/30">
          <h2 className="text-lg font-mono text-neutral-50">{t('admin.contests.imagesWarmup.statusTitle')}</h2>
        </div>

        {images.length === 0 ? (
          <EmptyState title={t('admin.contests.imagesWarmup.empty')} />
        ) : (
          <div className="p-4">
            <List columns={columns} data={images} renderCell={renderCell} />
          </div>
        )}
      </Card>
    </div>
  );
}

export default AdminImagesWarmup;
