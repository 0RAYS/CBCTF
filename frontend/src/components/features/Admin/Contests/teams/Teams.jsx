/** 团队管理组件 */

import { IconEdit, IconEye, IconScale, IconSearch, IconTrash, IconUserMinus } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import {
  Avatar,
  Button,
  EmptyState,
  Input,
  List,
  Modal,
  Pagination,
  ScrollingText,
  Select,
  Spinner,
  StatusTag,
  Textarea,
} from '@/components/common';
import Checkbox from '@/components/common/Checkbox';

function AdminTeams({
  teams = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 6,
  teamMembers = [],
  editForm = {},
  selectedUserId = '',
  kickUserId = '',
  showModal = false,
  modalMode = 'edit',
  selectedTeam = null,
  scoreForm = { score: '', reason: '' },
  adjustments = [],
  loading = false,
  onPageChange,
  onEditTeam,
  onDeleteTeam,
  onAdjustScore,
  onModalClose,
  onModalSubmit,
  onFormChange,
  onScoreFormChange,
  onUserSelect,
  onKickUserSelect,
  onKickSubmit,
  nameQuery = '',
  descQuery = '',
  searchResults = [],
  searchLoading = false,
  onNameChange,
  onDescChange,
  searchRef,
  isSearchMode = false,
  onViewDetail,
  onPictureUpload,
}) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language || 'en-US';

  const columns = [
    { key: 'team', label: t('admin.contests.teams.table.team'), width: '25%' },
    {
      key: 'score',
      label: t('admin.contests.teams.table.score'),
      width: '10%',
    },
    {
      key: 'members',
      label: t('admin.contests.teams.table.members'),
      width: '10%',
    },
    {
      key: 'lastSubmit',
      label: t('admin.contests.teams.table.lastSubmit'),
      width: '20%',
    },
    {
      key: 'status',
      label: t('admin.contests.teams.table.status'),
      width: '20%',
    },
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
        return (
          <div className="flex flex-col">
            <span className="text-neutral-300 font-mono">{team.score?.toLocaleString(locale) || 0}</span>
            {team.adjustment ? (
              <span className="text-xs text-yellow-400 font-mono">
                {t('admin.contests.teams.score.adjusted', {
                  value: `${team.adjustment > 0 ? '+' : ''}${team.adjustment.toLocaleString(locale)}`,
                })}
              </span>
            ) : null}
          </div>
        );
      case 'members':
        return (
          <span className="text-neutral-300 font-mono">
            {t('admin.contests.teams.memberCount', {
              count: team.users == null ? t('common.notAvailable') : team.users.toLocaleString(locale),
            })}
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
          {t('admin.contests.teams.modal.deletePrompt', {
            name: selectedTeam?.name || '',
          })}
        </p>
      );
    }

    if (modalMode === 'score') {
      return (
        <div className="space-y-4 [&>div>label]:font-mono [&>div>label]:mb-2">
          <div className="flex items-center justify-between rounded-md border border-neutral-700 px-4 py-3">
            <span className="text-sm text-neutral-400">{t('admin.contests.teams.score.current')}</span>
            <span className="font-mono text-neutral-50">{(selectedTeam?.score || 0).toLocaleString(locale)}</span>
          </div>
          <Input
            label={t('admin.contests.teams.score.delta')}
            type="number"
            step="0.01"
            value={scoreForm.score}
            onChange={(e) => onScoreFormChange({ ...scoreForm, score: e.target.value })}
            placeholder={t('admin.contests.teams.score.deltaPlaceholder')}
          />
          <Textarea
            label={t('admin.contests.teams.score.reason')}
            value={scoreForm.reason}
            onChange={(e) => onScoreFormChange({ ...scoreForm, reason: e.target.value })}
            rows={3}
            placeholder={t('admin.contests.teams.score.reasonPlaceholder')}
          />
          {adjustments.length > 0 && (
            <div className="border-t border-neutral-700 pt-4 space-y-2">
              <p className="text-sm font-mono text-neutral-400">{t('admin.contests.teams.score.history')}</p>
              {adjustments.map((adjustment) => (
                <div key={adjustment.id} className="flex items-start justify-between gap-3 text-sm">
                  <span className="text-neutral-300 [overflow-wrap:anywhere]">{adjustment.reason}</span>
                  <span className="font-mono text-yellow-400 shrink-0">
                    {adjustment.score > 0 ? '+' : ''}
                    {adjustment.score.toLocaleString(locale)}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      );
    }

    return (
      <div className="space-y-4 [&>div>label]:font-mono [&>div>label]:mb-2">
        <Input
          label={t('admin.contests.teams.form.name')}
          type="text"
          value={editForm.name}
          onChange={(e) => onFormChange({ ...editForm, name: e.target.value })}
        />
        <Textarea
          label={t('admin.contests.teams.form.description')}
          value={editForm.description}
          onChange={(e) => onFormChange({ ...editForm, description: e.target.value })}
          rows={3}
        />
        <Select
          label={t('admin.contests.teams.form.captain')}
          value={selectedUserId}
          onChange={(e) => onUserSelect(e.target.value)}
          options={teamMembers.map((member) => ({
            value: member.id,
            label: member.name,
          }))}
          placeholder={t('admin.contests.teams.form.selectCaptain')}
        />
        <Input
          label={t('admin.contests.teams.form.inviteCode')}
          type="text"
          value={editForm.captcha}
          onChange={(e) => onFormChange({ ...editForm, captcha: e.target.value })}
        />
        <div className="flex gap-6">
          <Checkbox
            id="banned"
            checked={editForm.banned}
            onChange={(e) => onFormChange({ ...editForm, banned: e.target.checked })}
            label={t('admin.contests.teams.labels.banned')}
            className="font-mono text-sm text-neutral-400"
          />
          <Checkbox
            id="hidden"
            checked={editForm.hidden}
            onChange={(e) => onFormChange({ ...editForm, hidden: e.target.checked })}
            label={t('admin.contests.teams.labels.hidden')}
            className="font-mono text-sm text-neutral-400"
          />
        </div>

        <div className="border-t border-neutral-700 pt-4">
          <p className="text-sm font-mono text-neutral-400 mb-3">
            {t('admin.contests.teams.modal.kickPrompt', {
              name: selectedTeam?.name || '',
            })}
          </p>
          {teamMembers.length === 0 ? (
            <EmptyState title={t('admin.contests.teams.empty.noMembers')} />
          ) : (
            <div className="flex items-end gap-3">
              <div className="flex-1 [&_label]:font-mono [&_label]:mb-2">
                <Select
                  label={t('admin.contests.teams.modal.selectMemberLabel')}
                  value={kickUserId}
                  onChange={(e) => onKickUserSelect(e.target.value)}
                  options={teamMembers.map((member) => ({
                    value: member.id,
                    label: member.name,
                  }))}
                  placeholder={t('admin.contests.teams.modal.selectMemberPlaceholder')}
                />
              </div>
              <Button
                variant="ghost"
                size="sm"
                className="!text-yellow-400 hover:!text-yellow-300 shrink-0 mb-0.5"
                onClick={onKickSubmit}
                disabled={!kickUserId}
              >
                <IconUserMinus size={16} className="mr-1" />
                {t('admin.contests.teams.modal.kickAction')}
              </Button>
            </div>
          )}
        </div>
      </div>
    );
  };

  const renderModalFooter = () => (
    <>
      <Button size="sm" variant="ghost" onClick={onModalClose}>
        {t('common.cancel')}
      </Button>
      <Button size="sm" variant={modalMode === 'delete' ? 'danger' : 'primary'} onClick={onModalSubmit}>
        {t('common.confirm')}
      </Button>
    </>
  );

  return (
    <div className="w-full mx-auto">
      <div className="mb-8" />

      <div className="mb-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div ref={searchRef} className="[&_label]:font-mono [&_label]:text-neutral-400 [&_label]:mb-2">
            <Input
              label={t('admin.contests.teams.search.label')}
              type="text"
              value={nameQuery}
              placeholder={t('admin.contests.teams.search.placeholder')}
              onChange={(e) => onNameChange(e.target.value)}
              icon={<IconSearch size={16} />}
              iconRight={searchLoading && <Spinner size="sm" />}
            />
          </div>
          <div className="[&_label]:font-mono [&_label]:text-neutral-400 [&_label]:mb-2">
            <Input
              label={t('admin.contests.teams.search.descLabel')}
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
        minWidth={1000}
        columns={columns}
        data={displayTeams}
        renderCell={renderCell}
        getRowActions={(team) => [
          {
            key: 'detail',
            label: t('common.viewDetails'),
            icon: <IconEye size={18} />,
            hidden: !onViewDetail,
            onClick: () => onViewDetail(team),
          },
          {
            key: 'edit',
            label: t('common.edit'),
            icon: <IconEdit size={18} />,
            onClick: () => onEditTeam(team),
          },
          {
            key: 'score',
            label: t('admin.contests.teams.actions.adjustScore'),
            icon: <IconScale size={18} />,
            onClick: () => onAdjustScore(team),
          },
          {
            key: 'delete',
            label: t('common.delete'),
            icon: <IconTrash size={18} />,
            danger: true,
            onClick: () => onDeleteTeam(team),
          },
        ]}
        loading={loading}
        empty={!loading && displayTeams.length === 0}
        emptyContent={emptyContent}
        onRowClick={onViewDetail || onEditTeam}
      />

      {!isSearchMode && totalCount > pageSize && (
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

      <Modal
        isOpen={showModal}
        onClose={onModalClose}
        title={
          modalMode === 'delete'
            ? t('admin.contests.teams.modal.title.delete')
            : modalMode === 'score'
              ? t('admin.contests.teams.modal.title.score')
              : t('admin.contests.teams.modal.title.edit')
        }
        size={modalMode === 'delete' ? 'sm' : 'lg'}
        footer={renderModalFooter()}
      >
        {renderModalContent()}
      </Modal>
    </div>
  );
}

export default AdminTeams;
