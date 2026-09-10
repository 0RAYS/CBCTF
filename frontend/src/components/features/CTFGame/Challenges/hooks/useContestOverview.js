import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getContestInfo, getContestNotices } from '../../../../../api/contest';
import { getTeamInfo, getTeamMembers } from '../../../../../api/game/team';
import { toast } from '../../../../../utils/toast';
import { getContestOverview } from '../models/challengeViewModel';

export default function useContestOverview(contestId) {
  const { t } = useTranslation();
  const [contestStatus, setContestStatus] = useState({});
  const [teamInfo, setTeamInfo] = useState(null);
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const requestRef = useRef(0);

  const refresh = async () => {
    const request = ++requestRef.current;
    const isCurrent = () => request === requestRef.current;
    setError(null);
    if (!contestStatus.status) setLoading(true);

    // Optional notices must not block the contest on latency or failure.
    getContestNotices(contestId)
      .then((response) => {
        if (!isCurrent()) return;
        if (response.code !== 200) throw new Error(response.msg || t('errors.requestFailed'));
        setNotifications(
          (response.data?.notices || []).map((notice) => ({
            type: notice.type || 'info',
            title: notice.title,
            message: notice.content,
          }))
        );
      })
      .catch((error) => {
        if (isCurrent()) toast.danger({ description: error.message || t('errors.requestFailed') });
      });

    try {
      const responses = await Promise.all([
        getContestInfo(contestId),
        getTeamMembers(contestId),
        getTeamInfo(contestId),
      ]);
      if (!isCurrent()) return;
      const failed = responses.find((response) => response.code !== 200);
      if (failed) throw new Error(failed.msg || t('game.challenges.toast.fetchFailed'));
      const [contest, members, team] = responses;
      setContestStatus(getContestOverview(contest.data, team.data, Date.now()));
      setTeamInfo({
        members: members.data.map(({ picture, name }) => ({ picture, name })),
        name: team.data.name,
      });
    } catch (error) {
      if (!isCurrent()) return;
      setError(error.message || t('game.challenges.toast.fetchFailed'));
      if (contestStatus.status) {
        toast.danger({ title: t('game.challenges.toast.fetchFailed'), description: error.message });
      }
    } finally {
      if (isCurrent()) setLoading(false);
    }
  };

  useEffect(() => {
    refresh();
    return () => {
      requestRef.current += 1;
    };
  }, [contestId]);

  return { contestStatus, setContestStatus, teamInfo, notifications, loading, error, refresh };
}
