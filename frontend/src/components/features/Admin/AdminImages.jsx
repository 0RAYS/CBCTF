import { useState } from 'react';
import { motion } from 'motion/react';
import { IconTrash, IconDownload } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { List } from '../../common';
import { useTranslation } from 'react-i18next';

/**
 * 管理后台图片管理组件
 * @param {Object} props
 * @param {Array} props.images - 图片列表
 * @param {number} props.totalCount - 总图片数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {boolean} props.loading - 是否加载中
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onDelete - 删除图片回调
 * @param {function} props.onBatchDelete - 批量删除回调
 * @param {function} props.onDownload - 下载图片回调
 */
function AdminImages({
  images = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 12,
  loading = false,
  onPageChange,
  onDelete,
  onBatchDelete,
  onDownload,
}) {
  const { t, i18n } = useTranslation();
  const [selectedImages, setSelectedImages] = useState([]);

  // 格式化文件大小
  const formatSize = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // 处理选择图片
  const handleSelectImage = (imageId) => {
    setSelectedImages((prev) => (prev.includes(imageId) ? prev.filter((id) => id !== imageId) : [...prev, imageId]));
  };

  // 处理全选
  const handleSelectAll = () => {
    setSelectedImages((prev) => (prev.length === images.length ? [] : images.map((img) => img.id)));
  };

  const selectAllLabel = (
    <div className="flex items-center gap-2">
      <input
        type="checkbox"
        checked={selectedImages.length === images.length && images.length > 0}
        onChange={handleSelectAll}
        className="w-4 h-4 rounded border-neutral-300/30 bg-black/30 text-geek-400 focus:ring-0 focus:ring-offset-0"
      />
      <span className="text-xs font-mono text-neutral-400">{t('admin.images.selectAll')}</span>
    </div>
  );

  const columns = [
    { key: 'select', label: selectAllLabel, width: '10%' },
    { key: 'preview', label: t('admin.images.columns.preview'), width: '12%' },
    { key: 'filename', label: t('admin.images.columns.filename'), width: '26%' },
    { key: 'size', label: t('admin.images.columns.size'), width: '10%' },
    { key: 'type', label: t('admin.images.columns.type'), width: '10%' },
    { key: 'uploaded', label: t('admin.images.columns.uploaded'), width: '18%' },
    { key: 'actions', label: t('admin.images.columns.actions'), width: '14%' },
  ];

  const renderCell = (image, column) => {
    switch (column.key) {
      case 'select':
        return (
          <input
            type="checkbox"
            checked={selectedImages.includes(image.id)}
            onChange={() => handleSelectImage(image.id)}
            className="w-4 h-4 rounded border-neutral-300/30 bg-black/30 text-geek-400 focus:ring-0 focus:ring-offset-0"
          />
        );

      case 'preview':
        return (
          <div className="w-20 h-12 rounded border border-neutral-300/20 overflow-hidden bg-black/40">
            <img src={image.url} alt={image.filename} className="w-full h-full object-cover" />
          </div>
        );

      case 'filename':
        return (
          <div className="min-w-0">
            <span className="text-sm font-mono text-neutral-300 truncate block" title={image.filename}>
              {image.filename}
            </span>
          </div>
        );

      case 'size':
        return <span className="text-xs font-mono text-neutral-400">{formatSize(image.size)}</span>;

      case 'type':
        return <span className="text-xs font-mono text-neutral-400 uppercase">{image.type}</span>;

      case 'uploaded':
        return (
          <span className="text-xs font-mono text-neutral-400">
            {new Date(image.uploadTime).toLocaleDateString(i18n.language || 'en-US')}
          </span>
        );

      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-black/30 !text-neutral-300 hover:!text-geek-400"
              onClick={() => onDownload?.(image)}
            >
              <IconDownload size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-black/30 !text-red-400 hover:!text-red-300"
              onClick={() => onDelete?.(image)}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        );

      default:
        return image[column.key];
    }
  };

  return (
    <div className="w-full mx-auto">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        {selectedImages.length > 0 && (
          <div className="flex justify-end items-center mb-6">
            <Button
              variant="danger"
              size="sm"
              align="icon-left"
              icon={<IconTrash size={16} />}
              onClick={() => {
                onBatchDelete?.(selectedImages);
                setSelectedImages([]);
              }}
            >
              {t('admin.images.actions.batchDelete', { count: selectedImages.length })}
            </Button>
          </div>
        )}

        <List
          columns={columns}
          data={images}
          renderCell={renderCell}
          loading={loading}
          empty={images.length === 0}
          emptyContent={t('admin.images.empty')}
        />

        {totalCount > pageSize && (
          <div className="mt-6">
            <Pagination
              total={Math.ceil(totalCount / pageSize)}
              current={currentPage}
              onChange={onPageChange}
              showTotal
              totalItems={totalCount}
            />
          </div>
        )}
      </motion.div>
    </div>
  );
}

export default AdminImages;
