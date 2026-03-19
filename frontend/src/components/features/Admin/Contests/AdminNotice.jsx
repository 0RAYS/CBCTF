/**
 * 比赛公告管理组件
 * @param {Object} props
 * @param {Array} props.notices - 公告列表
 * @param {number} props.totalCount - 公告总数
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {Function} props.onPageChange - 页码改变回调
 * @param {Function} props.onCreateNotice - 创建公告回调
 * @param {Function} props.onUpdateNotice - 更新公告回调
 * @param {Function} props.onDeleteNotice - 删除公告回调
 */

import { useState } from 'react';
import { IconEdit, IconTrash, IconPlus } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, FormField, Input, List, Modal, ModalFooter, Pagination, Select, Textarea } from '../../../common';

function AdminNotice({
  notices = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  onPageChange,
  onCreateNotice,
  onUpdateNotice,
  onDeleteNotice,
}) {
  const { t } = useTranslation();

  const [showModal, setShowModal] = useState(false);
  const [mode, setMode] = useState('create'); // 'create' | 'edit' | 'delete'
  const [selectedNotice, setSelectedNotice] = useState(null);
  const [form, setForm] = useState({
    title: '',
    content: '',
    type: '',
  });

  const columns = [
    { key: 'title', label: t('admin.contests.notices.table.title'), width: '15%' },
    { key: 'type', label: t('admin.contests.notices.table.type'), width: '10%' },
    { key: 'content', label: t('admin.contests.notices.table.content'), width: '48%' },
    { key: 'actions', label: t('admin.contests.notices.table.actions'), width: '7%' },
  ];

  const typeLabels = {
    normal: t('admin.contests.notices.types.normal'),
    important: t('admin.contests.notices.types.important'),
    update: t('admin.contests.notices.types.update'),
  };

  const renderCell = (notice, column) => {
    switch (column.key) {
      case 'title':
        return (
          <div className="min-w-0">
            <span className="text-neutral-50 font-mono truncate block">{notice.title}</span>
          </div>
        );
      case 'type':
        return <span className="text-neutral-300 font-mono text-sm">{typeLabels[notice.type] || notice.type}</span>;
      case 'content':
        return <div className="text-neutral-300 line-clamp-2 whitespace-normal">{notice.content}</div>;
      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!text-geek-400 hover:!text-geek-300"
              onClick={() => handleEdit(notice)}
            >
              <IconEdit size={16} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!text-red-400 hover:!text-red-300"
              onClick={() => handleDelete(notice)}
            >
              <IconTrash size={16} />
            </Button>
          </div>
        );
      default:
        return notice[column.key];
    }
  };

  const handleCreate = () => {
    setMode('create');
    setForm({ title: '', content: '', type: 'normal' });
    setShowModal(true);
  };

  const handleEdit = (notice) => {
    setMode('edit');
    setSelectedNotice(notice);
    setForm({
      title: notice.title,
      content: notice.content,
      type: notice.type,
    });
    setShowModal(true);
  };

  const handleDelete = (notice) => {
    setMode('delete');
    setSelectedNotice(notice);
    setShowModal(true);
  };

  const handleSubmit = () => {
    if (mode === 'create') {
      onCreateNotice(form);
    } else if (mode === 'edit') {
      onUpdateNotice(selectedNotice.id, form);
    } else {
      onDeleteNotice(selectedNotice.id);
    }
    setShowModal(false);
  };

  const renderModalContent = () => {
    if (mode === 'delete') {
      return (
        <p className="text-neutral-300">
          {t('admin.contests.notices.modal.deletePrompt', { title: selectedNotice?.title || '' })}
        </p>
      );
    }

    return (
      <div className="space-y-4">
        <FormField label={t('admin.contests.notices.form.title')} className="[&_label]:font-mono [&_label]:mb-2">
          <Input
            type="text"
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            placeholder={t('admin.contests.notices.form.titlePlaceholder')}
            required
          />
        </FormField>
        <FormField label={t('admin.contests.notices.form.type')} className="[&_label]:font-mono [&_label]:mb-2">
          <Select
            value={form.type}
            onChange={(e) => setForm({ ...form, type: e.target.value })}
            options={[
              { value: 'normal', label: t('admin.contests.notices.types.normal') },
              { value: 'important', label: t('admin.contests.notices.types.important') },
              { value: 'update', label: t('admin.contests.notices.types.update') },
            ]}
            required
          />
        </FormField>
        <FormField label={t('admin.contests.notices.form.content')} className="[&_label]:font-mono [&_label]:mb-2">
          <Textarea
            required
            value={form.content}
            onChange={(e) => setForm({ ...form, content: e.target.value })}
            placeholder={t('admin.contests.notices.form.contentPlaceholder')}
            rows={5}
          />
        </FormField>
      </div>
    );
  };

  const renderModalFooter = () => (
    <ModalFooter
      onCancel={() => setShowModal(false)}
      onSubmit={handleSubmit}
      cancelLabel={t('common.cancel')}
      submitLabel={t('common.confirm')}
      submitVariant={mode === 'delete' ? 'danger' : 'primary'}
    />
  );

  return (
    <div className="w-full mx-auto">
      {/* 头部 */}
      <div className="flex justify-end items-center mb-8">
        <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={18} />} onClick={handleCreate}>
          {t('admin.contests.notices.actions.create')}
        </Button>
      </div>

      <List
        columns={columns}
        data={notices}
        renderCell={renderCell}
        empty={notices.length === 0}
        emptyContent={t('common.noData')}
        footer={
          totalCount > pageSize ? (
            <Pagination
              total={Math.ceil(totalCount / pageSize)}
              current={currentPage}
              pageSize={pageSize}
              onChange={onPageChange}
              showTotal
              totalItems={totalCount}
            />
          ) : null
        }
      />

      {/* 模态框 */}
      <Modal
        isOpen={showModal}
        onClose={() => setShowModal(false)}
        title={
          mode === 'create'
            ? t('admin.contests.notices.modal.title.create')
            : mode === 'edit'
              ? t('admin.contests.notices.modal.title.edit')
              : t('admin.contests.notices.modal.title.delete')
        }
        size={mode === 'delete' ? 'sm' : 'md'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>
    </div>
  );
}

export default AdminNotice;
