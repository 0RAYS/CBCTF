import request from '../request';

export function getContestList({ limit = 6, offset = 0 }) {
  return request({
    url: '/admin/contests',
    method: 'GET',
    params: {
      limit,
      offset,
    },
  });
}

export function createContest(data) {
  return request({
    url: '/admin/contests',
    method: 'POST',
    data,
  });
}

export function deleteContest(contestId) {
  return request({
    url: `/admin/contests/${contestId}`,
    method: 'DELETE',
  });
}

export function updateContestPicture(contestId, file) {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/admin/contests/${contestId}/picture`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

export function getContestInfo(contestId) {
  return request({
    url: `/admin/contests/${contestId}`,
    method: 'GET',
  });
}

export function getContestChallenges(contestId, params = { limit: 20, offset: 0 }) {
  return request({
    url: `/admin/contests/${contestId}/challenges`,
    method: 'GET',
    params,
  });
}

export function addContestChallenge(contestId, challenge_ids) {
  return request({
    url: `/admin/contests/${contestId}/challenges`,
    method: 'POST',
    noToast: true,
    data: {
      challenge_ids: challenge_ids,
    },
  });
}

export function updateContestChallenge(contestId, challengeId, data) {
  return request({
    url: `/admin/contests/${contestId}/challenges/${challengeId}`,
    method: 'PUT',
    data,
  });
}

export function removeContestChallenge(contestId, challengeId) {
  return request({
    url: `/admin/contests/${contestId}/challenges/${challengeId}`,
    method: 'DELETE',
  });
}

export function resetContestChallengeSolves(contestId, challengeId) {
  return request({
    url: `/admin/contests/${contestId}/challenges/${challengeId}/reset-solves`,
    method: 'POST',
  });
}

export function updateContestInfo(contestId, data) {
  return request({
    url: `/admin/contests/${contestId}`,
    method: 'PUT',
    data,
  });
}

export function getContestTeams(contestId, params = { limit: 20, offset: 0 }) {
  return request({
    url: `/admin/contests/${contestId}/teams`,
    method: 'GET',
    params,
  });
}

export function getContestTeam(contestId, teamId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}`,
    method: 'GET',
  });
}

export function getContestLeaderboard(contestId, params = { limit: 20, offset: 0 }, noLoading = false) {
  return request({
    url: `/admin/contests/${contestId}/leaderboard`,
    method: 'GET',
    params,
    noLoading,
  });
}

export function updateTeamPicture(contestId, teamId, file) {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/picture`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

export function updateTeamInfo(contestId, teamId, data) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}`,
    method: 'PUT',
    data,
  });
}

export function deleteTeam(contestId, teamId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}`,
    method: 'DELETE',
  });
}

export function kickTeamMember(contestId, teamId, userId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/kick`,
    method: 'POST',
    data: {
      user_id: userId,
    },
  });
}

export function getTeamMembers(contestID, teamID) {
  return request({
    url: `/admin/contests/${contestID}/teams/${teamID}/users`,
    method: 'GET',
  });
}

export const getTeamContainers = async (contestId, teamId, params) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/victims`,
    method: 'GET',
    params,
  });
};

export function downloadContainerTraffic(contestId, teamId, containerId) {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/victims/${containerId}/traffic/download`,
    method: 'GET',
    responseType: 'blob',
  });
}

export const getContestNotices = (contestId, params = {}) => {
  return request({
    url: `/admin/contests/${contestId}/notices`,
    method: 'GET',
    params,
  });
};

export const createContestNotice = (contestId, data) => {
  return request({
    url: `/admin/contests/${contestId}/notices`,
    method: 'POST',
    data,
  });
};

export const updateContestNotice = (contestId, noticeId, data) => {
  return request({
    url: `/admin/contests/${contestId}/notices/${noticeId}`,
    method: 'PUT',
    data,
  });
};

export const deleteContestNotice = (contestId, noticeId) => {
  return request({
    url: `/admin/contests/${contestId}/notices/${noticeId}`,
    method: 'DELETE',
  });
};

export const getContestTeamSubmissions = (contestId, teamId, params) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/submissions`,
    method: 'GET',
    params,
  });
};

export const getContestTeamWriteups = (contestId, teamId, params) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/writeups`,
    method: 'GET',
    params,
  });
};

export const downloadContestTeamWriteup = (contestId, teamId, writeupId) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/writeups/${writeupId}`,
    method: 'GET',
    responseType: 'blob',
  }).then((response) => {
    return response;
  });
};

export const exportContestWriteups = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/writeups/export`,
    method: 'GET',
    responseType: 'blob',
  });
};

export const exportContestLeaderboard = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/leaderboard/export`,
    method: 'GET',
    responseType: 'blob',
  });
};

export const getContestTeamTraffic = (contestId, teamId, victimID, params, signal) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/victims/${victimID}/traffic`,
    method: 'GET',
    params,
    signal,
    noLoading: true,
    noToast: true,
  });
};

export const getContestPullImages = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/images`,
    method: 'GET',
  });
};

export const pullContestImages = (contestId, data) => {
  return request({
    url: `/admin/contests/${contestId}/images`,
    method: 'POST',
    noToast: true,
    data,
  });
};

export const detectContestImages = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/images/detect`,
    method: 'POST',
    noToast: true,
  });
};

export const getContestLeaderboardChallenges = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/contests/${contestId}/leaderboard/challenges`,
    method: 'GET',
    params,
  });
};

export const getContestLeaderboardTimeline = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/leaderboard/timeline`,
    method: 'GET',
  });
};

export const getContestVictims = (contestId, params = {}) => {
  return request({
    url: `/admin/contests/${contestId}/victims`,
    method: 'GET',
    params,
  });
};

export const stopContestVictims = (contestId, victimIds) => {
  return request({
    url: `/admin/contests/${contestId}/victims`,
    method: 'DELETE',
    noToast: true,
    data: {
      victims: victimIds,
    },
  });
};

export const startContestVictims = (contestId, challenges, teamRatio, duration) => {
  return request({
    url: `/admin/contests/${contestId}/victims`,
    method: 'POST',
    noToast: true,
    data: {
      challenges,
      team_ratio: teamRatio,
      duration,
    },
  });
};

export const getContestVictimPods = (contestId, victimId, signal) => {
  return request({
    url: `/admin/contests/${contestId}/victims/${victimId}/pods`,
    signal,
    noLoading: true,
    noToast: true,
    method: 'GET',
  });
};

export const getContestVictimPodLogs = (contestId, victimId, podName, container, lines = 1000) => {
  return request({
    url: `/admin/contests/${contestId}/victims/${victimId}/pods/logs`,
    method: 'GET',
    params: { pod_name: podName, container, lines },
  });
};

export const getContestCheats = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/contests/${contestId}/cheats`,
    method: 'GET',
    params,
  });
};

export const getContestCheat = (contestId, cheatId) => {
  return request({
    url: `/admin/contests/${contestId}/cheats/${cheatId}`,
    method: 'GET',
  });
};

export const getContestTeamFlags = (contestId, teamId) => {
  return request({
    url: `/admin/contests/${contestId}/teams/${teamId}/flags`,
    method: 'GET',
  });
};

export const updateContestCheat = (contestId, cheatId, data) => {
  return request({
    url: `/admin/contests/${contestId}/cheats/${cheatId}`,
    method: 'PUT',
    data,
  });
};

export const deleteContestCheat = (contestId, cheatId) => {
  return request({
    url: `/admin/contests/${contestId}/cheats/${cheatId}`,
    method: 'DELETE',
  });
};

export const deleteAllContestCheats = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/cheats`,
    method: 'DELETE',
  });
};

export const checkContestCheats = (contestId) => {
  return request({
    url: `/admin/contests/${contestId}/cheats`,
    method: 'POST',
    noToast: true,
  });
};

export const getIpInfo = (ip) => {
  return request({
    url: '/admin/ip',
    method: 'GET',
    params: { ip },
  });
};

export const getContestGenerators = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/admin/contests/${contestId}/generators`,
    method: 'GET',
    params,
  });
};

export const startContestGenerators = (contestId, challenges) => {
  return request({
    url: `/admin/contests/${contestId}/generators`,
    method: 'POST',
    noToast: true,
    data: { challenges },
  });
};

export const stopContestGenerators = (contestId, generatorIds) => {
  return request({
    url: `/admin/contests/${contestId}/generators`,
    method: 'DELETE',
    noToast: true,
    data: { generators: generatorIds },
  });
};

export const getContestGeneratorLogs = (contestId, generatorId, lines = 1000) => {
  return request({
    url: `/admin/contests/${contestId}/generators/${generatorId}/logs`,
    method: 'GET',
    params: { lines },
  });
};

export const getContestGeneratorStatus = (contestId, generatorId, signal) =>
  request({
    url: `/admin/contests/${contestId}/generators/${generatorId}/status`,
    method: 'GET',
    signal,
    noLoading: true,
    noToast: true,
  });
