import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import StatusPanel from '../../components/features/CTFGame/StatusPanel';
import ChallengeBoard from '../../components/features/CTFGame/Challenges/ChallengeBoard';
import ContestCountdown from '../../components/features/CTFGame/Challenges/ContestCountdown';
import ContestEnded from '../../components/features/CTFGame/Challenges/ContestEnded';
import ChallengeModal from '../../components/features/CTFGame/Challenges/ChallengeModal';
import useChallengeList from '../../components/features/CTFGame/Challenges/hooks/useChallengeList';
import useChallengeSession from '../../components/features/CTFGame/Challenges/hooks/useChallengeSession';
import useContestOverview from '../../components/features/CTFGame/Challenges/hooks/useContestOverview';
import useContestWriteups from '../../components/features/CTFGame/Writeup/useContestWriteups';
import Loading from '../../components/common/Loading';
import { Button, EmptyState } from '../../components/common';

function ContestChallenges({ contestId }) {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [showChallengesAfterEnd, setShowChallengesAfterEnd] = useState(false);
  const [uploading, setUploading] = useState(false);
  const overview = useContestOverview(contestId);
  const list = useChallengeList(contestId);
  const writeups = useContestWriteups(contestId);
  const session = useChallengeSession(contestId, {
    updateChallenge: list.updateChallenge,
    onSolved: () => {
      void overview.refresh();
      list.refresh();
    },
  });
  const { contestStatus } = overview;

  useEffect(() => {
    if (contestStatus.status === 'ended') void writeups.refresh();
  }, [contestStatus.status]);

  useEffect(() => {
    if (contestStatus.status === 'ended' && !showChallengesAfterEnd) session.closeChallenge();
  }, [contestStatus.status, showChallengesAfterEnd]);

  const handleUploadWriteup = async (file) => {
    setUploading(true);
    try {
      await writeups.upload(file);
    } finally {
      setUploading(false);
    }
  };

  if (overview.loading || uploading) return <Loading />;

  if (!contestStatus.status) {
    return (
      <div role="alert">
        <EmptyState
          title={t('game.challenges.toast.fetchFailed')}
          description={overview.error}
          action={
            <Button
              onClick={() => {
                void overview.refresh();
                void list.refreshCategories();
                list.refresh();
              }}
            >
              {t('common.refresh')}
            </Button>
          }
        />
      </div>
    );
  }

  const challengeContent = (
    <>
      <ChallengeBoard {...list.boardProps} onChallengeClick={session.openChallenge} teamInfo={overview.teamInfo} />
      <ChallengeModal key={session.selectedChallenge?.id || 'closed'} {...session.modalProps} contest={contestStatus} />
    </>
  );

  return (
    <div className="contest-container mx-auto space-y-6">
      <h1 className="sr-only">{contestStatus.name || t('game.challenges.title')}</h1>
      <StatusPanel
        contestStatus={contestStatus}
        onStatusExpired={(status) => overview.setContestStatus((previous) => ({ ...previous, status }))}
        notifications={overview.notifications}
      />
      {contestStatus.status === 'upcoming' ? (
        <ContestCountdown startTime={contestStatus.startTime} joined={contestStatus.joined} onJoin={() => {}} />
      ) : contestStatus.status === 'ended' ? (
        showChallengesAfterEnd ? (
          <div>
            <div className="mb-4">
              <Button
                variant="secondary"
                size="sm"
                onClick={() => {
                  setShowChallengesAfterEnd(false);
                  session.closeChallenge();
                }}
              >
                {t('game.challenges.backToSummary')}
              </Button>
            </div>
            {challengeContent}
          </div>
        ) : (
          <ContestEnded
            contestInfo={{
              duration: `${contestStatus.duration}h`,
              totalTeams: contestStatus.teams,
              totalChallenges: contestStatus.totalChallenges,
              teamRank: contestStatus.team.rank,
              teamScore: contestStatus.team.score,
              teamSolved: contestStatus.team.solved,
            }}
            onViewScoreboard={() => navigate(`/contests/${contestId}/scoreboard`)}
            onUploadWriteup={handleUploadWriteup}
            onViewChallenges={() => setShowChallengesAfterEnd(true)}
            writeups={writeups.writeups}
          />
        )
      ) : (
        <div>{challengeContent}</div>
      )}
    </div>
  );
}

export default function GameChallengesPage() {
  const { contestId } = useParams();
  return <ContestChallenges key={contestId} contestId={contestId} />;
}
