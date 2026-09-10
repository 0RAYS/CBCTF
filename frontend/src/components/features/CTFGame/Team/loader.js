import { getTeamInfo, getTeamMembers, getTeamCaptcha } from '../../../../api/game/team';
import { mapTeamSettings } from './model';

export async function loadTeamSettings(contestId, kind = 'team') {
  if (kind === 'picture') {
    const response = await getTeamInfo(contestId);
    return response.code === 200 ? { picture: response.data.picture } : null;
  }
  const [teamInfo, members] = await Promise.all([getTeamInfo(contestId), getTeamMembers(contestId)]);
  if (teamInfo.code !== 200 || members.code !== 200) return null;
  const captcha = await getTeamCaptcha(contestId);
  return mapTeamSettings(teamInfo.data, members.data, captcha.data);
}
