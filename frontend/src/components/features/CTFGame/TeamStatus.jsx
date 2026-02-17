/**
 * 用户队伍状态组件
 * @param {Object} props
 * @param {Object} props.team - 队伍信息
 */
import ChallengeSolves from './ChallengeSolves';
import { ScrollingText, Card, Avatar } from '../../common';
import { useTranslation } from 'react-i18next';

function TeamStatus({ team }) {
  const { t, i18n } = useTranslation();
  if (!team) return null;

  return (
    <Card variant="default" padding="md" animate>
      <div className="flex items-center justify-between">
        {/* 左侧：队伍信息 */}
        <div className="flex items-center gap-4">
          <Avatar src={team.picture} name={team.name} size="md" className="border border-neutral-300" />
          <div>
            <ScrollingText text={team.name} className="text-neutral-50 font-mono text-lg" maxWidth={200} speed={15} />
            <div className="text-neutral-400 text-sm">{t('game.teamStatus.rank', { rank: team.rank })}</div>
          </div>
        </div>
        {/* 解题进度 */}
        {team.solved && <ChallengeSolves solved={team.solved} totalSolved={team.totalSolved} />}
        {/* 右侧：分数和解题数 */}
        <div className="flex items-center gap-8">
          <div className="text-right">
            <div className="text-neutral-400 text-sm mb-1">{t('common.score')}</div>
            <div className="text-geek-400 font-mono text-2xl">
              {team.score.toLocaleString(i18n.language || 'en-US')}
            </div>
          </div>
          <div className="text-right">
            <div className="text-neutral-400 text-sm mb-1">{t('common.solved')}</div>
            <div className="text-neutral-50 font-mono text-2xl">{team.totalSolved}</div>
          </div>
        </div>
      </div>
    </Card>
  );
}

export default TeamStatus;
