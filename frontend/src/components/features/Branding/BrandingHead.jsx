import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useBranding } from '../../../hooks/useBranding';

const resolveRouteTitle = (pathname, t) => {
  const matchers = [
    ['/admin/branding', 'nav.branding'],
    ['/admin/dashboard', 'nav.dashboard'],
    ['/admin/contests', 'nav.contests'],
    ['/admin/rbac', 'nav.rbac'],
    ['/admin/challenges', 'nav.challenges'],
    ['/admin/oauth', 'nav.oauth'],
    ['/admin/smtp', 'nav.smtp'],
    ['/admin/cronjobs', 'nav.cronjobs'],
    ['/admin/webhook', 'nav.webhook'],
    ['/admin/files', 'nav.files'],
    ['/admin/tasks', 'nav.tasks'],
    ['/admin/system', 'nav.system'],
    ['/admin/logs', 'nav.logs'],
    ['/admin/images', 'nav.images'],
    ['/admin/generators', 'nav.generators'],
    ['/admin/victims', 'nav.victims'],
    ['/settings', 'nav.settings'],
    ['/games', 'nav.games'],
    ['/login', 'auth.login'],
    ['/support', 'footer.support'],
    ['/contact', 'footer.contact'],
  ];
  const exact = matchers.find(([prefix]) => pathname.startsWith(prefix));
  if (exact) return t(exact[1]);
  if (pathname.startsWith('/contests/')) return t('nav.overview');
  return '';
};

function BrandingHead() {
  const location = useLocation();
  const { t, i18n } = useTranslation();
  const { browserTitle, browserDescription } = useBranding();

  useEffect(() => {
    const titleSuffix = resolveRouteTitle(location.pathname, t);
    document.title = titleSuffix ? `${browserTitle} - ${titleSuffix}` : browserTitle;

    const description = document.querySelector('meta[name="description"]');
    if (description) {
      description.setAttribute('content', browserDescription);
    }
  }, [browserDescription, browserTitle, i18n.resolvedLanguage, location.pathname, t]);

  return null;
}

export default BrandingHead;
