import { useTranslation } from 'react-i18next';
import { Button, Input, List, Modal, Pagination } from '../../../common';
import Checkbox from '../../../common/Checkbox';
import useGroupMembers from './useGroupMembers';

export default function GroupMembersDialog({ group, onClose, onChanged }) {
  const { t } = useTranslation();
  const members = useGroupMembers(group, onChanged);
  const close = () => {
    if (!members.pending) onClose();
  };
  const columns = [
    { key: 'id', label: 'ID', width: '15%' },
    { key: 'name', label: t('admin.rbac.groups.columns.userName'), width: '30%' },
    { key: 'email', label: t('admin.rbac.groups.columns.email'), width: '35%' },
  ];
  return (
    <Modal
      isOpen
      onClose={close}
      title={`${t('admin.rbac.groups.modal.usersTitle')} - ${group.name}`}
      size="xl"
      footer={
        <>
          <Button size="sm" variant="ghost" onClick={close} disabled={members.pending}>
            {t('common.confirm')}
          </Button>
          <Button
            size="sm"
            variant="primary"
            onClick={members.assignSelected}
            disabled={members.selectedIds.length === 0 || members.loadingCandidates || members.pending}
          >
            {t('admin.rbac.groups.form.addSelected')}
          </Button>
        </>
      }
    >
      <div className="space-y-6">
        <div className="space-y-3">
          <div className="flex items-center justify-between gap-3">
            <h3 className="text-sm font-mono text-neutral-200">{t('admin.rbac.groups.modal.currentUsersTitle')}</h3>
            <span className="text-xs text-neutral-500">
              {t('admin.rbac.groups.form.currentUsersCount', { count: members.userCount })}
            </span>
          </div>
          <List
            minWidth={720}
            columns={[...columns, { key: 'actions', label: t('admin.rbac.groups.columns.actions'), width: '20%' }]}
            data={members.users}
            loading={members.loadingUsers}
            empty={!members.loadingUsers && members.users.length === 0}
            emptyContent={t('admin.rbac.groups.empty.currentUsers')}
            animate={false}
            renderCell={(item, column) =>
              column.key === 'actions' ? (
                <Button size="sm" variant="danger" disabled={members.pending} onClick={() => members.removeUser(item)}>
                  {t('admin.rbac.groups.form.remove')}
                </Button>
              ) : (
                (item[column.key] ?? '-')
              )
            }
            footer={
              members.userCount > members.pageSize && (
                <Pagination
                  total={Math.ceil(members.userCount / members.pageSize)}
                  current={members.userPage}
                  onChange={members.setUserPage}
                  showTotal
                  totalItems={members.userCount}
                />
              )
            }
          />
        </div>
        <div className="border-t border-neutral-300/10 pt-6 space-y-4">
          <div className="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
            <div className="grid flex-1 gap-3 md:grid-cols-3">
              {Object.entries({ name: 'userNameSearch', email: 'userEmailSearch', description: 'userDescSearch' }).map(
                ([key, label]) => (
                  <Input
                    key={key}
                    label={t(`admin.rbac.groups.form.${label}`)}
                    value={members.queries[key]}
                    onChange={(e) => members.changeQuery(key, e.target.value)}
                    placeholder={t(`admin.rbac.groups.form.${label}Placeholder`)}
                  />
                )
              )}
            </div>
            <div className="pb-1 text-sm text-neutral-400">
              {t('admin.rbac.groups.form.selectedUsers', { count: members.selectedIds.length })}
            </div>
          </div>
          <Checkbox
            label={t('common.selectAll')}
            checked={members.allVisibleSelected}
            onChange={members.toggleAllCandidates}
            disabled={members.candidates.length === 0 || members.pending}
          />
          <List
            minWidth={800}
            columns={[
              { key: 'select', label: '', width: '10%' },
              { key: 'id', label: 'ID', width: '12%' },
              { key: 'name', label: t('admin.rbac.groups.columns.userName'), width: '24%' },
              { key: 'email', label: t('admin.rbac.groups.columns.email'), width: '28%' },
              { key: 'description', label: t('admin.rbac.groups.columns.description'), width: '26%' },
            ]}
            data={members.candidates}
            loading={members.loadingCandidates}
            empty={!members.loadingCandidates && members.candidates.length === 0}
            emptyContent={t(`admin.rbac.groups.empty.${members.searching ? 'searchCandidates' : 'availableUsers'}`)}
            animate={false}
            renderCell={(item, column) =>
              column.key === 'select' ? (
                <Checkbox
                  aria-label={`${item.name} (${item.id})`}
                  checked={members.selectedIds.includes(item.id)}
                  disabled={members.pending}
                  onChange={() => members.toggleCandidate(item.id)}
                />
              ) : (
                item[column.key] || '-'
              )
            }
            footer={
              members.candidateCount > members.pageSize && (
                <Pagination
                  total={Math.ceil(members.candidateCount / members.pageSize)}
                  current={members.candidatePage}
                  onChange={members.setCandidatePage}
                  showTotal
                  totalItems={members.candidateCount}
                />
              )
            }
          />
        </div>
      </div>
    </Modal>
  );
}
