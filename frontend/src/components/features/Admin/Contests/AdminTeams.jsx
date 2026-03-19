/**
 * 团队管理组件
 * @param {Object} props
 * @param {Array} props.teams - 团队列表
 * @param {number} props.totalCount - 团队总数
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {Array} props.teamMembers - 当前选中团队的成员列表
 * @param {Object} props.editForm - 编辑表单数据
 * @param {string} props.selectedUserId - 选中的用户ID
 * @param {boolean} props.showModal - 是否显示模态框
 * @param {string} props.modalMode - 模态框模式
 * @param {Object} props.selectedTeam - 选中的团队
 * @param {Function} props.onPageChange - 页码改变回调
 * @param {Function} props.onEditTeam - 编辑团队回调
 * @param {Function} props.onDeleteTeam - 删除团队回调
 * @param {Function} props.onKickMember - 踢出成员回调
 * @param {Function} props.onViewDetails - 查看详情回调
 * @param {Function} props.onModalClose - 关闭模态框回调
 * @param {Function} props.onModalSubmit - 模态框提交回调
 * @param {Function} props.onFormChange - 表单改变回调
 * @param {Function} props.onUserSelect - 用户选择回调
 * @param {string} props.searchQuery - 搜索查询
 * @param {Array} props.searchResults - 搜索结果
 * @param {boolean} props.searchLoading - 搜索加载状态
 * @param {Function} props.onSearchChange - 搜索输入变化回调
 * @param {Function} props.onSearchResultSelect - 搜索结果选择回调
 * @param {Function} props.onResetSearch - 重置搜索回调
 * @param {Object} props.searchRef - 搜索框ref
 * @param {boolean} props.isSearchMode - 是否处于搜索模式
 */

import { IconEdit, IconTrash, IconUserMinus, IconContainer, IconSearch } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import {
  Avatar,
  Button,
  EmptyState,
  FormField,
  FormSwitch,
  Input,
  List,
  Modal,
  ModalFooter,
  Pagination,
  ScrollingText,
  Select,
  Spinner,
  StatusTag,
  Textarea,
} from '../../../common';

function AdminTeams({
  teams = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 6,
  teamMembers = [],
  editForm = {},
  selectedUserId = '',
  showModal = false,
  modalMode = 'edit',
  selectedTeam = null,
  onPageChange,
  onEditTeam,
  onDeleteTeam,
  onKickMember,
  onViewDetails,
  onModalClose,
  onModalSubmit,
  onFormChange,
  onUserSelect,
  nameQuery = '',
  descQuery = '',
  searchResults = [],
  searchLoading = false,
  onNameChange,
  onDescChange,
  searchRef,
  isSearchMode = false,
  onRowClick,
  onPictureUpload,
}) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language || 'en-US';

  const columns = [
    { key: 'team', label: t('admin.contests.teams.table.team'), width: '8%' },
    { key: 'score', label: t('admin.contests.teams.table.score'), width: '5%' },
    { key: 'members', label: t('admin.contests.teams.table.members'), width: '5%' },
    { key: 'lastSubmit', label: t('admin.contests.teams.table.lastSubmit'), width: '8%' },
    { key: 'status', label: t('admin.contests.teams.table.status'), width: '5%' },
    { key: 'actions', label: t('admin.contests.teams.table.actions'), width: '7%' },
  ];

  const renderCell = (team, column) => {
    switch (column.key) {
      case 'team':
        return (
          <div className="flex items-center gap-3">
            <div
              className="relative w-8 h-8 rounded-lg overflow-hidden cursor-pointer group/avatar"
              onClick={(e) => {
                e.stopPropagation();
                onPictureUpload?.(team);
              }}
            >
              <Avatar src={team.picture} name={team.name} size="xs" className="border border-neutral-300/30" />
              <div
                className={[
                  'absolute inset-0 bg-black/50 flex items-center justify-center opacity-0',
                  'group-hover/avatar:opacity-100 transition-opacity',
                ].join(' ')}
              >
                <span className="text-neutral-300 text-[10px]">{t('admin.contests.teams.picture.replace')}</span>
              </div>
            </div>
            <div className="min-w-0">
              <ScrollingText text={team.name} className="text-neutral-50 font-medium" maxWidth={240} speed={15} />
              {team.description && (
                <div className="text-sm text-neutral-400 mt-0.5 line-clamp-1" title={team.description}>
                  {team.description}
                </div>
              )}
            </div>
          </div>
        );
      case 'score':
        return <span className="text-neutral-300 font-mono">{team.score?.toLocaleString(locale) || 0}</span>;
      case 'members':
        return (
          <span className="text-neutral-300 font-mono">
            {t('admin.contests.teams.memberCount', { count: (team.users || 0).toLocaleString(locale) })}
          </span>
        );
      case 'lastSubmit':
        return (
          <span className="text-neutral-300 font-mono text-sm">
            {team.last ? new Date(team.last).toLocaleString(locale) : t('admin.contests.teams.lastSubmit.none')}
          </span>
        );
      case 'status':
        return (
          <div className="flex items-center gap-2">
            <StatusTag
              type={team.banned ? 'error' : 'default'}
              text={team.banned ? t('admin.contests.teams.status.banned') : t('admin.contests.teams.status.normal')}
            />
            <StatusTag
              type={team.hidden ? 'warning' : 'default'}
              text={team.hidden ? t('admin.contests.teams.status.hidden') : t('admin.contests.teams.status.visible')}
            />
          </div>
        );
      case 'actions':
        return (
          <div className="flex items-center gap-3" onClick={(e) => e.stopPropagation()}>
            <Button
              variant="ghost"
              size="icon"
              className="!text-geek-400 hover:!text-geek-300"
              onClick={() => onEditTeam(team)}
            >
              <IconEdit size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!text-yellow-400 hover:!text-yellow-300"
              onClick={() => onKickMember(team)}
            >
              <IconUserMinus size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!text-neutral-400 hover:!text-neutral-300"
              onClick={() => onViewDetails(team)}
            >
              <IconContainer size={18} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!text-red-400 hover:!text-red-300"
              onClick={() => onDeleteTeam(team)}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        );
      default:
        return team[column.key];
    }
  };

  const displayTeams = isSearchMode ? searchResults : teams;
  const emptyContent = isSearchMode
    ? t('admin.contests.teams.empty.noResults')
    : t('admin.contests.teams.empty.noData');

  const renderModalContent = () => {
    if (modalMode === 'delete') {
      return (
        <p className="text-neutral-300">
          {t('admin.contests.teams.modal.deletePrompt', { name: selectedTeam?.name || '' })}
        </p>
      );
    }

    if (modalMode === 'kick') {
      return (
        <div className="space-y-4">
          <p className="text-neutral-300">
            {t('admin.contests.teams.modal.kickPrompt', { name: selectedTeam?.name || '' })}
          </p>
          {teamMembers.length === 0 ? (
            <EmptyState title={t('admin.contests.teams.empty.noMembers')} />
          ) : (
            <FormField
              label={t('admin.contests.teams.modal.selectMemberLabel')}
              className="[&_label]:font-mono [&_label]:mb-2"
            >
              <Select
                value={selectedUserId}
                onChange={(e) => onUserSelect(e.target.value)}
                options={teamMembers.map((member) => ({
                  value: member.id,
                  label: member.name,
                }))}
                placeholder={t('admin.contests.teams.modal.selectMemberPlaceholder')}
              />
            </FormField>
          )}
        </div>
      );
    }

    return (
      <div className="space-y-4">
        <FormField label={t('admin.contests.teams.form.name')} className="[&_label]:font-mono [&_label]:mb-2">
          <Input
            type="text"
            value={editForm.name}
            onChange={(e) => onFormChange({ ...editForm, name: e.target.value })}
          />
        </FormField>
        <FormField
          label={t('admin.contests.teams.form.description')}
          className="[&_label]:font-mono [&_label]:mb-2"
        >
          <Textarea
            value={editForm.description}
            onChange={(e) => onFormChange({ ...editForm, description: e.target.value })}
            rows={3}
          />
        </FormField>
        <FormField label={t('admin.contests.teams.form.captain')} className="[&_label]:font-mono [&_label]:mb-2">
          <Select
            value={selectedUserId}
            onChange={(e) => onUserSelect(e.target.value)}
            options={teamMembers.map((member) => ({
              value: member.id,
              label: member.name,
            }))}
            placeholder={t('admin.contests.teams.form.selectCaptain')}
          />
        </FormField>
        <FormField label={t('admin.contests.teams.form.inviteCode')} className="[&_label]:font-mono [&_label]:mb-2">
          <Input
            type="text"
            value={editForm.captcha}
            onChange={(e) => onFormChange({ ...editForm, captcha: e.target.value })}
          />
        </FormField>
        <div className="flex gap-6">
          <FormSwitch
            id="banned"
            checked={editForm.banned}
            onChange={(e) => onFormChange({ ...editForm, banned: e.target.checked })}
            label={t('admin.contests.teams.labels.banned')}
            className="font-mono text-sm text-neutral-400"
          />
          <FormSwitch
            id="hidden"
            checked={editForm.hidden}
            onChange={(e) => onFormChange({ ...editForm, hidden: e.target.checked })}
            label={t('admin.contests.teams.labels.hidden')}
            className="font-mono text-sm text-neutral-400"
          />
        </div>
      </div>
    );
  };

  const renderModalFooter = () => (
    <ModalFooter
      onCancel={onModalClose}
      onSubmit={onModalSubmit}
      cancelLabel={t('common.cancel')}
      submitLabel={t('common.confirm')}
      submitVariant={modalMode === 'delete' ? 'danger' : 'primary'}
    />
  );

  return (
    <div className="w-full mx-auto">
      <div className="mb-8" />

      {/* 搜索框 */}
      <div className="mb-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div ref={searchRef}>
            <label className="block text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.teams.search.label')}
            </label>
            <Input
              type="text"
              value={nameQuery}
              placeholder={t('admin.contests.teams.search.placeholder')}
              onChange={(e) => onNameChange(e.target.value)}
              icon={<IconSearch size={16} />}
              iconRight={searchLoading && <Spinner size="sm" />}
            />
          </div>
          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.teams.search.descLabel')}
            </label>
            <Input
              type="text"
              value={descQuery}
              placeholder={t('admin.contests.teams.search.descPlaceholder')}
              onChange={(e) => onDescChange(e.target.value)}
              icon={<IconSearch size={16} />}
            />
          </div>
        </div>
      </div>

      <List
        columns={columns}
        data={displayTeams}
        renderCell={renderCell}
        empty={displayTeams.length === 0}
        emptyContent={emptyContent}
        onRowClick={onRowClick}
      />

      {/* 分页 - 只在非搜索模式下显示 */}
      {!isSearchMode && totalCount > pageSize && (
        <div className="mt-6">
          <Pagination
            total={Math.ceil(totalCount / pageSize)}
            current={currentPage}
            pageSize={pageSize}
            onChange={onPageChange}
            showTotal
            totalItems={totalCount}
          />
        </div>
      )}

      {/* 模态框 */}
      <Modal
        isOpen={showModal}
        onClose={onModalClose}
        title={
          modalMode === 'edit'
            ? t('admin.contests.teams.modal.title.edit')
            : modalMode === 'delete'
              ? t('admin.contests.teams.modal.title.delete')
              : t('admin.contests.teams.modal.title.kick')
        }
        size={modalMode === 'edit' ? 'lg' : 'sm'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>
    </div>
  );
}

export default AdminTeams;
