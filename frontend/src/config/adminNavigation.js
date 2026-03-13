/**
 * Admin sidebar navigation definitions.
 *
 * Each item has an optional `apiRoute` ("METHOD /path") that maps to the
 * backend RoutePermissions table.  When a user's accessible-routes list is
 * available the sidebar filters out items the user cannot reach.
 * Items without an `apiRoute` are always shown (e.g. Settings which uses
 * the /me endpoint the user already has access to).
 */

export const getAdminNavSections = (t, routes = null) => {
  const items = [
    { id: 'dashboard', label: t('nav.dashboard'), path: '/admin/dashboard', apiRoute: 'GET /admin/system/status' },
    { id: 'contests', label: t('nav.contests'), path: '/admin/contests', apiRoute: 'GET /admin/contests' },
    { id: 'rbac', label: t('nav.rbac'), path: '/admin/rbac', apiRoute: 'GET /admin/roles' },
    { id: 'challenges', label: t('nav.challenges'), path: '/admin/challenges', apiRoute: 'GET /admin/challenges' },
    { id: 'victims', label: t('nav.victims'), path: '/admin/victims', apiRoute: 'GET /admin/victims' },
    { id: 'generators', label: t('nav.generators'), path: '/admin/generators', apiRoute: 'GET /admin/generators' },
    { id: 'oauth', label: t('nav.oauth'), path: '/admin/oauth', apiRoute: 'GET /admin/oauth' },
    { id: 'smtp', label: t('nav.smtp'), path: '/admin/smtp', apiRoute: 'GET /admin/smtp' },
    { id: 'webhook', label: t('nav.webhook'), path: '/admin/webhook', apiRoute: 'GET /admin/webhook' },
    { id: 'files', label: t('nav.files'), path: '/admin/files', apiRoute: 'GET /admin/files' },
    { id: 'system', label: t('nav.system'), path: '/admin/system', apiRoute: 'GET /admin/system/config' },
    { id: 'logs', label: t('nav.logs'), path: '/admin/logs', apiRoute: 'GET /admin/logs' },
  ];

  const routeSet = routes ? new Set(routes) : null;
  const visibleItems = routeSet ? items.filter((item) => !item.apiRoute || routeSet.has(item.apiRoute)) : items;

  return [{ id: 'admin-main', title: t('admin.navigation'), items: visibleItems }];
};

export const getAdminContestNavSections = (t, contestId, routes = null) => {
  const items = [
    {
      id: 'overview',
      label: t('nav.settings'),
      path: `/admin/contests/${contestId}`,
      apiRoute: 'GET /admin/contests/:contestID',
    },
    {
      id: 'notices',
      label: t('nav.notices'),
      path: `/admin/contests/${contestId}/notices`,
      apiRoute: 'GET /admin/contests/:contestID/notices',
    },
    {
      id: 'challenges',
      label: t('nav.challenges'),
      path: `/admin/contests/${contestId}/challenges`,
      apiRoute: 'GET /admin/contests/:contestID/challenges',
    },
    {
      id: 'scoreboard',
      label: t('nav.scoreboard'),
      path: `/admin/contests/${contestId}/scoreboard`,
      apiRoute: 'GET /admin/contests/:contestID/scoreboard',
    },
    {
      id: 'teams',
      label: t('nav.team'),
      path: `/admin/contests/${contestId}/teams`,
      apiRoute: 'GET /admin/contests/:contestID/teams',
    },
    {
      id: 'images',
      label: t('nav.images'),
      path: `/admin/contests/${contestId}/images`,
      apiRoute: 'GET /admin/contests/:contestID/images',
    },
    {
      id: 'victims',
      label: t('nav.victims'),
      path: `/admin/contests/${contestId}/victims`,
      apiRoute: 'GET /admin/contests/:contestID/victims',
    },
    {
      id: 'generators',
      label: t('nav.generators'),
      path: `/admin/contests/${contestId}/generators`,
      apiRoute: 'GET /admin/contests/:contestID/generators',
    },
    {
      id: 'cheats',
      label: t('nav.cheats'),
      path: `/admin/contests/${contestId}/cheats`,
      apiRoute: 'GET /admin/contests/:contestID/cheats',
    },
  ];

  const routeSet = routes ? new Set(routes) : null;
  const visibleItems = routeSet ? items.filter((item) => !item.apiRoute || routeSet.has(item.apiRoute)) : items;

  return [{ id: 'admin-contest', title: t('admin.contest'), items: visibleItems }];
};
