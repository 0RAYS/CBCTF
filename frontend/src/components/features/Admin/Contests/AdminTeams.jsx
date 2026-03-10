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
import ModalButton from '../../../common/ModalButton';
import { Button, Pagination, ScrollingText, Modal, List, EmptyState, Avatar, Input, Spinner } from '../../../common';
import AdminTeamDetailDialog from './AdminTeamDetailDialog';

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
  showDetailDialog = false,
  detailTeam = null,
  detailTab = 'info',
  detailMembers = [],
  detailMembersLoading = false,
  detailSubmissions = [],
  detailSubmissionCount = 0,
  detailSubmissionPage = 1,
  detailWriteups = [],
  detailWriteupCount = 0,
  detailWriteupPage = 1,
  detailContainers = [],
  detailContainerCount = 0,
  detailContainerPage = 1,
  detailLoading = { submissions: false, writeups: false, traffic: false },
  onDetailClose,
  onDetailTabChange,
  onDetailPageChange,
  onDetailDownloadTraffic,
  onDetailDownloadWriteup,
  onPictureUpload,
  detailFlags = [],
  detailFlagsLoading = false,
}) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language || 'en-US';

  const columns = [
    { key: 'team', label: t('admin.contests.teams.table.team'), width: '28%' },
    { key: 'score', label: t('admin.contests.teams.table.score'), width: '12%' },
    { key: 'members', label: t('admin.contests.teams.table.members'), width: '12%' },
    { key: 'lastSubmit', label: t('admin.contests.teams.table.lastSubmit'), width: '20%' },
    { key: 'status', label: t('admin.contests.teams.table.status'), width: '16%' },
    { key: 'actions', label: t('admin.contests.teams.table.actions'), width: '12%' },
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
            <span
              className={`px-2 py-1 rounded text-xs font-mono ${
                team.banned ? 'bg-red-400/20 text-red-400' : 'bg-neutral-400/20 text-neutral-400'
              }`}
            >
              {team.banned ? t('admin.contests.teams.status.banned') : t('admin.contests.teams.status.normal')}
            </span>
            <span
              className={`px-2 py-1 rounded text-xs font-mono ${
                team.hidden ? 'bg-yellow-400/20 text-yellow-400' : 'bg-neutral-400/20 text-neutral-400'
              }`}
            >
              {team.hidden ? t('admin.contests.teams.status.hidden') : t('admin.contests.teams.status.visible')}
            </span>
          </div>
        );
      case 'actions':
        return (
          <div className="flex justify-end gap-3" onClick={(e) => e.stopPropagation()}>
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
            <div>
              <label className="block text-sm font-mono text-neutral-400 mb-2">
                {t('admin.contests.teams.modal.selectMemberLabel')}
              </label>
              <select
                value={selectedUserId}
                onChange={(e) => onUserSelect(e.target.value)}
                className="select-custom select-custom-md"
              >
                <option value="">{t('admin.contests.teams.modal.selectMemberPlaceholder')}</option>
                {teamMembers.map((member) => (
                  <option key={member.id} value={member.id}>
                    {member.name}
                  </option>
                ))}
              </select>
            </div>
          )}
        </div>
      );
    }

    return (
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-mono text-neutral-400 mb-2">{t('admin.contests.teams.form.name')}</label>
          <input
            type="text"
            value={editForm.name}
            onChange={(e) => onFormChange({ ...editForm, name: e.target.value })}
            className="w-full h-[40px] px-4 bg-black/20 border border-neutral-300/30 rounded-md
                            text-neutral-50 placeholder-neutral-500
                            focus:outline-none focus:border-geek-400 focus:shadow-focus
                            transition-all duration-200"
          />
        </div>
        <div>
          <label className="block text-sm font-mono text-neutral-400 mb-2">
            {t('admin.contests.teams.form.description')}
          </label>
          <textarea
            value={editForm.description}
            onChange={(e) => onFormChange({ ...editForm, description: e.target.value })}
            rows={3}
            className="w-full p-4 bg-black/20 border border-neutral-300/30 rounded-md
                            text-neutral-50 placeholder-neutral-500
                            focus:outline-none focus:border-geek-400 focus:shadow-focus
                            transition-all duration-200 resize-none"
          />
        </div>
        <div>
          <label className="block text-sm font-mono text-neutral-400 mb-2">
            {t('admin.contests.teams.form.captain')}
          </label>
          <select
            value={selectedUserId}
            onChange={(e) => onUserSelect(e.target.value)}
            className="select-custom select-custom-md"
          >
            <option value="">{t('admin.contests.teams.form.selectCaptain')}</option>
            {teamMembers.map((member) => (
              <option key={member.id} value={member.id}>
                {member.name}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-mono text-neutral-400 mb-2">
            {t('admin.contests.teams.form.inviteCode')}
          </label>
          <input
            type="text"
            value={editForm.captcha}
            onChange={(e) => onFormChange({ ...editForm, captcha: e.target.value })}
            className="w-full h-[40px] px-4 bg-black/20 border border-neutral-300/30 rounded-md
                            text-neutral-50 placeholder-neutral-500
                            focus:outline-none focus:border-geek-400 focus:shadow-focus
                            transition-all duration-200"
          />
        </div>
        <div className="flex gap-6">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="banned"
              checked={editForm.banned}
              onChange={(e) => onFormChange({ ...editForm, banned: e.target.checked })}
              className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 
                                focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
            />
            <label htmlFor="banned" className="text-sm font-mono text-neutral-400">
              {t('admin.contests.teams.labels.banned')}
            </label>
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="hidden"
              checked={editForm.hidden}
              onChange={(e) => onFormChange({ ...editForm, hidden: e.target.checked })}
              className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 
                                focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
            />
            <label htmlFor="hidden" className="text-sm font-mono text-neutral-400">
              {t('admin.contests.teams.labels.hidden')}
            </label>
          </div>
        </div>
      </div>
    );
  };

  const renderModalFooter = () => (
    <>
      <ModalButton onClick={onModalClose}>{t('common.cancel')}</ModalButton>
      <ModalButton variant={modalMode === 'delete' ? 'danger' : 'primary'} onClick={onModalSubmit}>
        {t('common.confirm')}
      </ModalButton>
    </>
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

      {/* 详情对话框 */}
      <AdminTeamDetailDialog
        isOpen={showDetailDialog}
        onClose={onDetailClose}
        team={detailTeam}
        activeTab={detailTab}
        onTabChange={onDetailTabChange}
        members={detailMembers}
        membersLoading={detailMembersLoading}
        detailSubmissions={detailSubmissions}
        detailSubmissionCount={detailSubmissionCount}
        detailSubmissionPage={detailSubmissionPage}
        detailWriteups={detailWriteups}
        detailWriteupCount={detailWriteupCount}
        detailWriteupPage={detailWriteupPage}
        detailContainers={detailContainers}
        detailContainerCount={detailContainerCount}
        detailContainerPage={detailContainerPage}
        detailLoading={detailLoading}
        onDetailPageChange={onDetailPageChange}
        onDetailDownloadTraffic={onDetailDownloadTraffic}
        onDetailDownloadWriteup={onDetailDownloadWriteup}
        detailFlags={detailFlags}
        detailFlagsLoading={detailFlagsLoading}
      />
    </div>
  );
}

export default AdminTeams;
