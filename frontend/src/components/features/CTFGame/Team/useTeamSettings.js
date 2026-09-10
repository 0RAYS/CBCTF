import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  updateTeamInfo,
  uploadTeamPicture,
  deleteTeam,
  kickTeamMember,
  updateTeamCaptcha,
} from '../../../../api/game/team';
import { toast } from '../../../../utils/toast';
import { loadTeamSettings } from './loader';
import { teamUpdatePayload, unchangedTeamFields } from './model';

// The page keys this hook's owner by contest and user, isolating drafts and pending work.
export default function useTeamSettings(contestId, navigate) {
  const { t } = useTranslation();
  const [team, setTeam] = useState(null);
  const [loading, setLoading] = useState(true);
  const scopeRef = useRef(null);

  const patchTeam = (scope, patch) => {
    if (!scope.active) return;
    for (const key of Object.keys(patch)) scope.versions[key] = (scope.versions[key] || 0) + 1;
    setTeam((prev) => ({ ...prev, ...patch }));
  };
  const refresh = async (scope, kind = 'team') => {
    if (!scope.active) return false;
    const request = (scope.requests[kind] || 0) + 1;
    scope.requests[kind] = request;
    const versions = { ...scope.versions };
    const data = await loadTeamSettings(contestId, kind);
    if (!scope.active || request !== scope.requests[kind] || !data) return false;
    patchTeam(scope, unchangedTeamFields(data, versions, scope.versions));
    return true;
  };

  useEffect(() => {
    const scope = { active: true, requests: {}, versions: {} };
    scopeRef.current = scope;
    refresh(scope)
      .catch((error) => {
        if (scope.active) toast.danger({ title: t('game.team.toast.fetchFailed'), description: error.message });
      })
      .finally(() => {
        if (scope.active) setLoading(false);
      });
    return () => {
      scope.active = false;
    };
  }, [contestId]);

  const mutate = async (request, onSuccess, failure, title = true) => {
    const scope = scopeRef.current;
    if (!scope?.active) return;
    try {
      const response = await request();
      if (scope.active && response.code === 200) await onSuccess(scope, response);
    } catch (error) {
      if (!scope.active) return;
      toast.danger(
        title ? { title: t(failure), description: error.message } : { description: error.message || t(failure) }
      );
    }
  };

  return {
    team,
    loading,
    onCopyCode: () => {
      navigator.clipboard.writeText(team.inviteCode);
      toast.success({ title: t('game.team.toast.inviteCopied') });
    },
    onRefreshCode: () =>
      mutate(
        () => updateTeamCaptcha(contestId),
        (scope, response) => {
          patchTeam(scope, { inviteCode: response.data });
          toast.success({ title: t('game.team.toast.inviteRefreshed') });
        },
        'game.team.toast.inviteRefreshFailed'
      ),
    onEditTeam: (updatedTeam) =>
      mutate(
        () => updateTeamInfo(contestId, teamUpdatePayload(updatedTeam)),
        (scope) => {
          patchTeam(scope, teamUpdatePayload(updatedTeam));
          toast.success({ title: t('game.team.toast.updateSuccess') });
        },
        'game.team.toast.updateFailed'
      ),
    onKickMember: (memberName) => {
      const member = team.members.find((member) => member.name === memberName);
      if (!member?.id) return;
      return mutate(
        () => kickTeamMember(contestId, member.id),
        async (scope) => {
          toast.success({
            title: t('game.team.toast.memberRemoved'),
            description: t('game.team.toast.memberRemovedDescription', { name: memberName }),
          });
          await refresh(scope);
        },
        'game.team.toast.removeFailed',
        false
      );
    },
    onDisbandTeam: () =>
      mutate(
        () => deleteTeam(contestId),
        () => {
          toast.success({ title: t('game.team.toast.disbanded') });
          navigate(`/contests/${contestId}`);
        },
        'game.team.toast.disbandFailed',
        false
      ),
    onPictureUpload: (file) =>
      mutate(
        () => uploadTeamPicture(contestId, file),
        async (scope) => {
          if (await refresh(scope, 'picture')) toast.success({ title: t('game.team.toast.avatarUpdated') });
        },
        'game.team.toast.avatarUploadFailed'
      ),
  };
}
