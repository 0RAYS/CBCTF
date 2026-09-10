import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getContestTeam } from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import TeamDetailSession from './TeamDetailSession';

export function useTeamDetailDialog(contestId) {
  const { t } = useTranslation();
  const [session, setSession] = useState(null);
  const request = useRef(0);
  useEffect(() => {
    setSession(null);
    return () => {
      request.current += 1;
    };
  }, [contestId]);
  const close = () => {
    request.current += 1;
    setSession(null);
  };
  const openTeamDetail = async (teamOrId) => {
    if (!teamOrId) return;
    const version = ++request.current;
    try {
      let team = teamOrId;
      if (typeof teamOrId !== 'object') {
        const response = await getContestTeam(contestId, teamOrId);
        if (response.code !== 200) return;
        team = response.data;
      }
      if (version === request.current) setSession({ team, version, contestId });
    } catch (error) {
      if (version === request.current)
        toast.danger({ description: error.message || t('admin.contests.cheats.toast.fetchFailed') });
    }
  };
  const renderTeamDetailDialog = () =>
    session && session.contestId === contestId ? (
      <TeamDetailSession key={session.version} team={session.team} contestId={contestId} onClose={close} />
    ) : null;
  return { openTeamDetail, renderTeamDetailDialog };
}
