import { useTranslation } from 'react-i18next';
import Avatar from '../../../common/Avatar';
import { Button, Card } from '../../../common';

export default function MembersCard({ leader, members, isLeader, onKickMember }) {
  const { t } = useTranslation();
  return (
    <Card variant="default" padding="none" animate className="overflow-hidden">
      <div className="p-4 border-b border-neutral-300/30">
        <h2 className="text-neutral-50 font-mono">{t('game.team.settings.teamMembers')}</h2>
      </div>
      <div className="divide-y divide-neutral-300/10">
        <div className="p-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div className="flex min-w-0 items-center gap-4">
            <Avatar src={leader.picture} name={leader.name} size="sm" className="shrink-0 border-2 border-yellow-400" />
            <div className="min-w-0 [overflow-wrap:anywhere]">
              <div className="text-neutral-50 font-mono">{leader.name}</div>
              <div className="text-yellow-400 text-sm font-mono">{t('game.team.settings.leader')}</div>
            </div>
          </div>
          <div className="min-w-0 text-neutral-400 text-sm [overflow-wrap:anywhere]">{leader.email}</div>
        </div>
        {members.map((member) => (
          <div key={member.id} className="p-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
            <div className="flex min-w-0 items-center gap-4">
              <Avatar
                src={member.picture}
                name={member.name}
                size="sm"
                className="shrink-0 border-2 border-neutral-300"
              />
              <div className="min-w-0 [overflow-wrap:anywhere]">
                <div className="text-neutral-50 font-mono">{member.name}</div>
                <div className="text-neutral-400 text-sm font-mono">{t('game.team.settings.member')}</div>
              </div>
            </div>
            <div className="flex min-w-0 items-center justify-between gap-3">
              <div className="min-w-0 text-neutral-400 text-sm [overflow-wrap:anywhere]">{member.email}</div>
              {isLeader && (
                <Button
                  variant="ghost"
                  size="icon"
                  className="shrink-0 text-red-400! hover:bg-red-400/10!"
                  aria-label={`${t('game.team.settings.kickModal.title')}: ${member.name}`}
                  onClick={() => onKickMember(member.name)}
                >
                  ✕
                </Button>
              )}
            </div>
          </div>
        ))}
      </div>
    </Card>
  );
}
