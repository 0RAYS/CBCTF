import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import ContestDetail from '../../components/features/CTFGame/OverView/ContestDetail';
import { transformContestData } from '../../components/features/CTFGame/OverView/contestModel.js';
import { toast } from '../../utils/toast';
import { getContestInfo } from '../../api/contest';
import Loading from '../../components/common/Loading';
import { useTranslation } from 'react-i18next';

function GameDetailPage() {
  const { contestId } = useParams();
  const [contest, setContest] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const { t } = useTranslation();

  const handleJoinContest = () => {
    navigate(`/contests/${contestId}/challenges`);
  };

  useEffect(() => {
    const fetchContestDetail = async () => {
      try {
        const response = await getContestInfo(contestId);
        if (response.code === 200) {
          const transformedData = transformContestData(response.data);
          setContest(transformedData);
        } else {
          throw new Error(response.msg || t('game.detail.toast.fetchFailed'));
        }
      } catch (error) {
        toast.danger({ title: t('game.detail.toast.fetchFailed'), description: error.message });
      } finally {
        setLoading(false);
      }
    };

    fetchContestDetail();
  }, [contestId, t]);

  if (loading) {
    return <Loading />;
  }

  return <ContestDetail contest={contest} handleJoinContest={handleJoinContest} />;
}

export default GameDetailPage;
