export const normalizeInstanceStatus = (status) => {
  const normalized = typeof status === 'string' ? status.toLowerCase() : '';
  return ['waiting', 'pending', 'terminating', 'running'].includes(normalized) ? normalized : '';
};

export const isInstanceTransitioning = (status) => ['waiting', 'pending', 'terminating'].includes(status);

export const normalizeCategories = (categories) => (Array.isArray(categories) ? categories.filter(Boolean) : []);

export function mapChallengeStatusToViewModel(challenge, statusData = null) {
  const remote = statusData?.remote || challenge.remote || {};
  const instanceStatus = normalizeInstanceStatus(remote.status);
  const timeLeft = Number(remote.remaining) || 0;
  const remoteDuration = Number(remote.duration) || 0;
  const previousDuration = Number(challenge.instanceDuration) || 0;
  const instanceDuration =
    remoteDuration > 0
      ? remoteDuration
      : instanceStatus === 'running'
        ? Math.max(previousDuration, timeLeft)
        : previousDuration;

  return {
    ...challenge,
    title: challenge.title || challenge.name,
    attachments: challenge.attachments || [],
    attachment: statusData?.file ?? challenge.attachment ?? '',
    hasInstance: challenge.hasInstance ?? challenge.type === 'pods',
    hasAttachments: challenge.hasAttachments ?? challenge.type === 'dynamic',
    instanceStatus,
    instanceRunning: instanceStatus === 'running',
    instancePending: instanceStatus === 'pending',
    instanceWaiting: instanceStatus === 'waiting',
    instanceTerminating: instanceStatus === 'terminating',
    instanceIP: remote.target || challenge.instanceIP || [''],
    instanceDuration,
    instanceTimeLeft: timeLeft,
    solves: challenge.solves ?? challenge.solvers ?? 0,
    isInitialized: statusData?.init ?? challenge.isInitialized ?? challenge.init,
    isSolved: statusData?.solved ?? challenge.isSolved ?? challenge.solved ?? false,
    solved: statusData?.solved ?? challenge.solved ?? challenge.isSolved ?? false,
    attempts: statusData?.attempts ?? challenge.attempts,
    maxAttempts: challenge.maxAttempts ?? challenge.attempt,
    options: challenge.options || [],
  };
}

export function getContestOverview(contest, team, now) {
  const start = new Date(contest.start).getTime();
  const end = start + contest.duration * 1000;
  return {
    status: now < start ? 'upcoming' : now > end ? 'ended' : 'running',
    startTime: new Date(start).toISOString(),
    endTime: new Date(end).toISOString(),
    joined: true,
    duration: contest.duration / 3600,
    prefix: contest.prefix,
    team: {
      score: team.score,
      rank: team.rank || 0,
      solved: team.solved.reduce((sum, category) => sum + category.solved, 0),
    },
    teams: contest.teams,
    totalChallenges: team.solved.reduce((sum, category) => sum + category.all, 0),
  };
}

export function formatTimeLeft(seconds) {
  seconds = Math.floor(seconds);
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;
  return [hours, minutes, secs].map((part) => part.toString().padStart(2, '0')).join(':');
}
