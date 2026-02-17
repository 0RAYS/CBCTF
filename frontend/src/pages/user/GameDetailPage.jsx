import { useState, useEffect, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import ContestDetail from '../../components/features/CTFGame/OverView/ContestDetail';
import { toast } from '../../utils/toast';
import { getContestInfo } from '../../api/contest';
import Loading from '../../components/common/Loading';
import { useTranslation } from 'react-i18next';

const getDefaultContestData = (t) => ({
  rules: t('game.detail.defaultRules', { returnObjects: true }),
  prizes: t('game.detail.defaultPrizes', { returnObjects: true }),
  timeline: t('game.detail.defaultTimeline', { returnObjects: true }),
});

// 转换比赛状态
const getContestStatus = (startTime, duration) => {
  const now = new Date().getTime();
  const start = new Date(startTime).getTime();
  const end = start + duration * 1000;

  if (now < start) return 'upcoming';
  if (now > end) return 'ended';
  return 'active';
};

// 转换API数据为组件所需格式
const transformContestData = (apiData, defaults) => {
  if (!apiData) return null;

  const startTime = new Date(apiData.start);
  const endTime = new Date(startTime.getTime() + apiData.duration * 1000);

  return {
    title: `${apiData.prefix} ${apiData.name}`,
    description: apiData.description,
    image:
      apiData.picture || 'https://images.unsplash.com/photo-1562813733-b31f71025d54?q=80&w=2069&auto=format&fit=crop',
    status: getContestStatus(apiData.start, apiData.duration),
    startTime: startTime.toISOString(),
    endTime: endTime.toISOString(),
    participants: apiData.users || 0,
    rules: apiData.rules || defaults.rules,
    prizes: apiData.prizes || defaults.prizes,
    timeline: apiData.timelines || defaults.timeline,
    // 额外信息，可用于后续功能扩展
    teamSize: apiData.size,
    teamsCount: apiData.teams,
    noticesCount: apiData.notices,
    isBlood: apiData.blood,
    isHidden: apiData.hidden,
  };
};

function GameDetailPage() {
  const { contestId } = useParams();
  const [contest, setContest] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const { t } = useTranslation();
  const defaultContestData = useMemo(() => getDefaultContestData(t), [t]);

  const handleJoinContest = () => {
    navigate(`/contests/${contestId}/challenges`);
  };

  useEffect(() => {
    const fetchContestDetail = async () => {
      try {
        const response = await getContestInfo(contestId);
        if (response.code === 200) {
          const transformedData = transformContestData(response.data, defaultContestData);
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
  }, [contestId, defaultContestData, t]);

  if (loading) {
    return <Loading />;
  }

  return <ContestDetail contest={contest} handleJoinContest={handleJoinContest} />;
}

export default GameDetailPage;
