import { useState } from 'react';
import { motion } from 'motion/react';
import { IconTrash, IconDownload, IconFile, IconPhoto, IconFileText } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { List } from '../../common';
import { useTranslation } from 'react-i18next';

/**
 * 管理后台文件管理组件
 * @param {Object} props
 * @param {Array} props.files - 文件列表
 * @param {number} props.totalCount - 总文件数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {string} props.fileType - 文件类型
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onTypeChange - 文件类型改变回调
 * @param {function} props.onDelete - 删除文件回调
 * @param {function} props.onBatchDelete - 批量删除回调
 * @param {function} props.onDownload - 下载文件回调
 */
function AdminFiles({
  files = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 12,
  fileType = 'file',
  onPageChange,
  onTypeChange,
  onDelete,
  onBatchDelete,
  onDownload,
}) {
  const { t, i18n } = useTranslation();
  const [selectedFiles, setSelectedFiles] = useState([]);

  // 文件类型选项
  const fileTypeOptions = [
    { value: 'file', label: t('admin.files.types.file'), icon: IconFile },
    { value: 'picture', label: t('admin.files.types.picture'), icon: IconPhoto },
    { value: 'traffic', label: t('admin.files.types.traffic'), icon: IconFile },
    { value: 'writeup', label: t('admin.files.types.writeup'), icon: IconFileText },
  ];

  // 格式化文件大小
  const formatSize = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // 获取文件图标
  const getFileIcon = (type) => {
    const option = fileTypeOptions.find((opt) => opt.value === type);
    return option ? option.icon : IconFile;
  };

  // 处理选择文件
  const handleSelectFile = (fileId) => {
    setSelectedFiles((prev) => (prev.includes(fileId) ? prev.filter((id) => id !== fileId) : [...prev, fileId]));
  };

  // 处理全选
  const handleSelectAll = () => {
    setSelectedFiles((prev) => (prev.length === files.length ? [] : files.map((file) => file.id)));
  };

  // 处理类型切换
  const handleTypeChange = (newType) => {
    onTypeChange?.(newType);
    onPageChange(1);
    setSelectedFiles([]);
  };

  const selectAllLabel = (
    <div className="flex items-center gap-2">
      <input
        type="checkbox"
        checked={selectedFiles.length === files.length && files.length > 0}
        onChange={handleSelectAll}
        className="w-4 h-4 rounded border-neutral-300/30 bg-black/30 text-geek-400 focus:ring-0 focus:ring-offset-0"
      />
      <span className="text-xs font-mono text-neutral-400">{t('admin.files.selectAll')}</span>
    </div>
  );

  const columns = [
    { key: 'select', label: selectAllLabel, width: '10%' },
    { key: 'name', label: t('admin.files.columns.name'), width: '28%' },
    { key: 'type', label: t('admin.files.columns.type'), width: '10%' },
    { key: 'size', label: t('admin.files.columns.size'), width: '10%' },
    { key: 'uploaded', label: t('admin.files.columns.uploaded'), width: '14%' },
    { key: 'meta', label: t('admin.files.columns.meta'), width: '18%' },
    { key: 'actions', label: t('admin.files.columns.actions'), width: '10%' },
  ];

  const renderMetaBadge = (text) => (
    <span className="text-xs font-mono text-neutral-400 bg-neutral-800/50 px-2 py-1 rounded" title={text}>
      {text}
    </span>
  );

  const renderCell = (file, column) => {
    switch (column.key) {
      case 'select':
        return (
          <input
            type="checkbox"
            checked={selectedFiles.includes(file.id)}
            onChange={() => handleSelectFile(file.id)}
            className="w-4 h-4 rounded border-neutral-300/30 bg-black/30 text-geek-400 focus:ring-0 focus:ring-offset-0"
          />
        );

      case 'name': {
        const FileIcon = getFileIcon(file.type);
        const isPicture = file.type === 'picture';
        return (
          <div className="flex items-center gap-3 min-w-0">
            <div className="w-12 h-12 flex-shrink-0">
              {isPicture ? (
                <img
                  src={file.url}
                  alt={file.filename}
                  className="w-full h-full object-cover rounded border border-neutral-300/20"
                  onError={(e) => {
                    e.target.style.display = 'none';
                    e.target.nextSibling.style.display = 'flex';
                  }}
                />
              ) : null}
              <div
                className={`w-full h-full flex items-center justify-center rounded border border-neutral-300/20 ${isPicture ? 'hidden' : 'flex'}`}
                style={{ display: isPicture ? 'none' : 'flex' }}
              >
                <FileIcon size={22} className="text-neutral-400" />
              </div>
            </div>
            <div className="min-w-0">
              <span className="text-sm font-mono text-neutral-300 truncate block" title={file.filename}>
                {file.filename}
              </span>
              <span className="text-xs font-mono text-neutral-400 uppercase bg-neutral-800/50 px-2 py-1 rounded">
                {file.suffix}
              </span>
            </div>
          </div>
        );
      }

      case 'type':
        return <span className="text-xs font-mono text-neutral-400 uppercase">{file.type}</span>;

      case 'size':
        return <span className="text-xs font-mono text-neutral-400">{formatSize(file.size)}</span>;

      case 'uploaded':
        return (
          <span className="text-xs font-mono text-neutral-400">
            {new Date(file.uploadTime).toLocaleDateString(i18n.language || 'en-US')}
          </span>
        );

      case 'meta': {
        const badges = [];
        if (file.model && file.modelId > 0)
          badges.push(renderMetaBadge(t('admin.files.meta.model', { model: file.model, id: file.modelId })));
        if (file.hash) badges.push(renderMetaBadge(t('admin.files.meta.sha256', { hash: file.hash })));

        return <div className="flex flex-wrap gap-2">{badges}</div>;
      }

      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-black/30 !text-neutral-300 hover:!text-geek-400"
              onClick={() => onDownload?.(file)}
            >
              <IconDownload size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!bg-black/30 !text-red-400 hover:!text-red-300"
              onClick={() => onDelete?.(file)}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        );

      default:
        return file[column.key];
    }
  };

  return (
    <div className="w-full mx-auto">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        {selectedFiles.length > 0 && (
          <div className="flex justify-end items-center mb-6">
            <Button
              variant="danger"
              size="sm"
              align="icon-left"
              icon={<IconTrash size={16} />}
              onClick={() => {
                onBatchDelete?.(selectedFiles);
                setSelectedFiles([]);
              }}
            >
              {t('admin.files.actions.batchDelete', { count: selectedFiles.length })}
            </Button>
          </div>
        )}

        <div className="mb-6">
          <div className="flex gap-2">
            {fileTypeOptions.map((option) => {
              const IconComponent = option.icon;
              return (
                <Button
                  key={option.value}
                  variant={fileType === option.value ? 'primary' : 'ghost'}
                  size="sm"
                  align="icon-left"
                  icon={<IconComponent size={16} />}
                  onClick={() => handleTypeChange(option.value)}
                >
                  {option.label}
                </Button>
              );
            })}
          </div>
        </div>

        <List
          columns={columns}
          data={files}
          renderCell={renderCell}
          empty={files.length === 0}
          emptyContent={t('admin.files.empty')}
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

export default AdminFiles;
