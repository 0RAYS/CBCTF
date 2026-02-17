import { useTranslation } from 'react-i18next';
import { Modal, Card, StatusTag, Avatar } from '../../common';

function AdminUserDetailDialog({ isOpen, onClose, user }) {
  const { t } = useTranslation();

  if (!user) return null;

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={t('admin.users.detail.title')} size="md">
      <div className="space-y-6">
        {/* User header */}
        <div className="flex items-start gap-4">
          <Avatar
            src={user.picture}
            name={user.name}
            size="lg"
            shape="circle"
            className="border border-neutral-300/30"
          />
          <div className="min-w-0 flex-1">
            <h3 className="text-lg font-mono text-neutral-50 font-medium">{user.name}</h3>
            <p className="text-sm text-neutral-400 mt-1">
              {user.description || t('admin.users.detail.info.noDescription')}
            </p>
            <div className="flex items-center gap-2 mt-2">
              {user.verified && <StatusTag type="success" text={t('admin.users.status.verified')} />}
              {user.banned && <StatusTag type="error" text={t('admin.users.status.banned')} />}
              {user.hidden && <StatusTag type="warning" text={t('admin.users.status.hidden')} />}
            </div>
          </div>
        </div>

        {/* Info grid */}
        <div className="grid grid-cols-2 gap-4">
          <Card variant="default" padding="md">
            <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.users.detail.info.id')}</div>
            <div className="text-lg font-mono text-neutral-50">{user.id}</div>
          </Card>
          <Card variant="default" padding="md">
            <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.users.detail.info.email')}</div>
            <div className="text-sm font-mono text-neutral-50 truncate" title={user.email}>
              {user.email || '-'}
            </div>
          </Card>
          <Card variant="default" padding="md">
            <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.users.detail.info.score')}</div>
            <div className="text-lg font-mono text-neutral-50">{user.score ?? 0}</div>
          </Card>
          <Card variant="default" padding="md">
            <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.users.detail.info.solved')}</div>
            <div className="text-lg font-mono text-neutral-50">{user.solved ?? 0}</div>
          </Card>
          <Card variant="default" padding="md">
            <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.users.detail.info.contests')}</div>
            <div className="text-lg font-mono text-neutral-50">{user.contests ?? 0}</div>
          </Card>
          <Card variant="default" padding="md">
            <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.users.detail.info.teams')}</div>
            <div className="text-lg font-mono text-neutral-50">{user.teams ?? 0}</div>
          </Card>
          {user.provider && (
            <Card variant="default" padding="md" className="col-span-2">
              <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.users.detail.info.provider')}</div>
              <div className="text-sm font-mono text-neutral-50">{user.provider}</div>
            </Card>
          )}
        </div>
      </div>
    </Modal>
  );
}

export default AdminUserDetailDialog;
