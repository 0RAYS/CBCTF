export const getAdminNavSections = (t) => [
  {
    id: 'admin-main',
    title: t('admin.navigation'),
    items: [
      { id: 'dashboard', label: t('nav.dashboard'), path: '/admin/dashboard' },
      { id: 'contests', label: t('nav.contests'), path: '/admin/contests' },
      { id: 'users', label: t('nav.users'), path: '/admin/users' },
      { id: 'challenges', label: t('nav.challenges'), path: '/admin/challenges' },
      { id: 'oauth', label: t('nav.oauth'), path: '/admin/oauth' },
      { id: 'smtp', label: t('nav.smtp'), path: '/admin/smtp' },
      { id: 'webhook', label: t('nav.webhook'), path: '/admin/webhook' },
      { id: 'files', label: t('nav.files'), path: '/admin/files' },
      { id: 'system', label: t('nav.system'), path: '/admin/system' },
      { id: 'logs', label: t('nav.logs'), path: '/admin/logs' },
    ],
  },
];

export const getAdminContestNavSections = (t, contestId) => [
  {
    id: 'admin-contest',
    title: t('admin.contest'),
    items: [
      { id: 'overview', label: t('nav.settings'), path: `/admin/contests/${contestId}` },
      { id: 'notices', label: t('nav.notices'), path: `/admin/contests/${contestId}/notices` },
      { id: 'challenges', label: t('nav.challenges'), path: `/admin/contests/${contestId}/challenges` },
      { id: 'scoreboard', label: t('nav.scoreboard'), path: `/admin/contests/${contestId}/scoreboard` },
      { id: 'teams', label: t('nav.team'), path: `/admin/contests/${contestId}/teams` },
      { id: 'containers', label: t('nav.containers'), path: `/admin/contests/${contestId}/containers` },
      { id: 'images', label: t('nav.images'), path: `/admin/contests/${contestId}/images` },
      { id: 'cheats', label: t('nav.cheats'), path: `/admin/contests/${contestId}/cheats` },
    ],
  },
];
