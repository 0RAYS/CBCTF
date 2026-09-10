import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import ConfirmModal from '../../../common/ConfirmModal';
import Avatar from '../../../common/Avatar';
import EditTeamModal from './EditTeamModal';
import InvitationCode from './InvitationCode';
import MembersCard from './MembersCard';
import { Button, Card } from '../../../common';
import MarkdownContent from '../../../common/MarkdownContent';

function TeamSettings({
  team,
  isLeader,
  onCopyCode = () => {},
  onRefreshCode = () => {},
  onEditTeam = () => {},
  onKickMember = () => {},
  onDisbandTeam = () => {},
  onPictureUpload = () => {},
}) {
  const [showEditModal, setShowEditModal] = useState(false);
  const [showKickModal, setShowKickModal] = useState(null);
  const [showDisbandModal, setShowDisbandModal] = useState(false);
  const { t } = useTranslation();

  const handlePictureUpload = (event) => {
    const file = event.target.files[0];
    event.target.value = '';
    if (file) {
      onPictureUpload(file);
    }
  };

  return (
    <div className="contest-container mx-auto space-y-6">
      {/* 队伍信息卡片 */}
      <Card variant="default" padding="md" animate>
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-4">
          <div className="flex min-w-0 items-center gap-4">
            <div className="relative group shrink-0">
              <Avatar src={team.picture} name={team.name} size="lg" className="border-2 border-neutral-300" />
              {isLeader && (
                <label
                  className="absolute inset-0 flex items-center justify-center 
                                    bg-black/50 opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 transition-opacity focus-within:ring-2 focus-within:ring-geek-400
                                    cursor-pointer rounded-lg"
                >
                  <input
                    type="file"
                    className="sr-only"
                    aria-label={t('game.team.settings.changeAvatar')}
                    accept="image/png,image/jpeg,image/jpg,image/gif"
                    onChange={handlePictureUpload}
                  />
                  <span className="text-neutral-50 text-sm font-mono">{t('game.team.settings.changeAvatar')}</span>
                </label>
              )}
            </div>
            <div className="min-w-0 [overflow-wrap:anywhere]">
              <div className="text-neutral-50 font-mono text-lg">{team.name}</div>
              <div className="text-neutral-400 text-sm">
                {t('game.team.settings.members', { count: team.members.length + 1 })}
              </div>
            </div>
          </div>
          <div className="shrink-0">
            {isLeader && (
              <Button variant="primary" size="sm" className="w-full sm:w-auto" onClick={() => setShowEditModal(true)}>
                {t('game.team.settings.editTeam')}
              </Button>
            )}
          </div>
        </div>

        {/* 队伍描述 */}
        <div className="mb-6">
          <div className="text-neutral-400 text-sm mb-2">{t('game.team.settings.description')}</div>
          <div className="p-3 bg-neutral-900 rounded-md">
            <MarkdownContent className="text-neutral-300 text-sm prose prose-invert prose-sm max-w-none min-w-0 [overflow-wrap:anywhere]">
              {team.description || t('game.team.settings.noDescription')}
            </MarkdownContent>
          </div>
        </div>

        {/* 邀请码部分 */}
        <InvitationCode
          code={team.inviteCode}
          isLeader={isLeader}
          onCopyCode={onCopyCode}
          onRefreshCode={onRefreshCode}
        />
      </Card>

      {/* 成员列表 */}
      <MembersCard leader={team.leader} members={team.members} isLeader={isLeader} onKickMember={setShowKickModal} />

      {isLeader && (
        <div className="border-t border-red-400/20 pt-4">
          <Button
            variant="ghost"
            textColor="text-red-400"
            size="sm"
            className="w-full sm:w-auto hover:bg-red-400/10!"
            onClick={() => setShowDisbandModal(true)}
          >
            {t('game.team.settings.disbandTeam')}
          </Button>
        </div>
      )}

      {/* 编辑模态框 */}
      <EditTeamModal
        isOpen={showEditModal}
        onClose={() => setShowEditModal(false)}
        team={team}
        onSave={(data) => {
          onEditTeam(data);
          setShowEditModal(false);
        }}
      />

      {/* 踢出确认模态框 */}
      <ConfirmModal
        isOpen={!!showKickModal}
        onClose={() => setShowKickModal(null)}
        onConfirm={() => {
          onKickMember(showKickModal);
          setShowKickModal(null);
        }}
        title={t('game.team.settings.kickModal.title')}
        message={t('game.team.settings.kickModal.message', { name: showKickModal })}
        confirmText={t('game.team.settings.kickModal.confirm')}
        type="danger"
      />

      {/* 解散确认模态框 */}
      <ConfirmModal
        isOpen={showDisbandModal}
        onClose={() => setShowDisbandModal(false)}
        onConfirm={() => {
          onDisbandTeam();
          setShowDisbandModal(false);
        }}
        title={t('game.team.settings.disbandModal.title')}
        message={t('game.team.settings.disbandModal.message')}
        confirmText={t('game.team.settings.disbandModal.confirm')}
        type="danger"
      />
    </div>
  );
}

export default TeamSettings;
