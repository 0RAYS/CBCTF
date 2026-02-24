import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { toast } from '../../utils/toast';
import { uploadWriteup, getWriteups } from '../../api/challenge';
import { getContestInfo } from '../../api/contest';
import WriteupUpload from '../../components/features/CTFGame/Challenges/WriteupUpload';
import Loading from '../../components/common/Loading';
import { useTranslation } from 'react-i18next';
import { isContestEnded } from '../../config/contest';

function GameWriteupPage() {
  const { contestId } = useParams();
  const navigate = useNavigate();
  const [writeups, setWriteups] = useState([]);
  const [loading, setLoading] = useState(true);
  const { t } = useTranslation();

  useEffect(() => {
    checkContestAndFetch();
  }, [contestId]);

  const checkContestAndFetch = async () => {
    try {
      const contestRes = await getContestInfo(contestId);
      if (contestRes.code === 200) {
        const contest = contestRes.data;
        if (isContestEnded(contest.start, contest.duration)) {
          navigate(`/contests/${contestId}/challenges`, { replace: true });
          return;
        }
      }
      await fetchWriteups();
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.fetchFailed') });
    } finally {
      setLoading(false);
    }
  };

  const fetchWriteups = async () => {
    try {
      const response = await getWriteups(contestId);
      if (response.code === 200 && response.data.writeups) {
        setWriteups(response.data.writeups);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.fetchWriteupsFailed') });
    }
  };

  const handleUploadWriteup = async (file) => {
    try {
      const res = await uploadWriteup(contestId, file);
      if (res.code === 200) {
        toast.success({
          title: t('game.challenges.toast.uploadSuccess'),
          description: t('game.challenges.toast.uploadThanks'),
        });
        fetchWriteups();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('game.challenges.toast.uploadFailed') });
    }
  };

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="contest-container mx-auto">
      <WriteupUpload onUploadWriteup={handleUploadWriteup} writeups={writeups} />
    </div>
  );
}

export default GameWriteupPage;
