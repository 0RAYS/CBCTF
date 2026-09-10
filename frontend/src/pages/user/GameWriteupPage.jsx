import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { toast } from '../../utils/toast';
import { getContestInfo } from '../../api/contest';
import WriteupUpload from '../../components/features/CTFGame/Writeup/WriteupUpload';
import useContestWriteups from '../../components/features/CTFGame/Writeup/useContestWriteups';
import Loading from '../../components/common/Loading';
import { useTranslation } from 'react-i18next';
import { getContestStatus } from '../../config/contest';

function ContestWriteups({ contestId }) {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const { t } = useTranslation();
  const { writeups, refresh, upload } = useContestWriteups(contestId);

  useEffect(() => {
    let active = true;
    const checkContestAndFetch = async () => {
      try {
        const response = await getContestInfo(contestId);
        if (!active) return;
        if (response.code === 200) {
          const contest = response.data;
          if (getContestStatus(contest.start, contest.duration) === 'upcoming') {
            navigate(`/contests/${contestId}/challenges`, { replace: true });
            return;
          }
        }
        await refresh();
      } catch (error) {
        if (active) toast.danger({ description: error.message || t('game.challenges.toast.fetchFailed') });
      } finally {
        if (active) setLoading(false);
      }
    };
    checkContestAndFetch();
    return () => {
      active = false;
    };
  }, [contestId]);

  if (loading) return <Loading />;

  return (
    <div className="contest-container mx-auto">
      <WriteupUpload onUploadWriteup={upload} writeups={writeups} />
    </div>
  );
}

export default function GameWriteupPage() {
  const { contestId } = useParams();
  return <ContestWriteups key={contestId} contestId={contestId} />;
}
