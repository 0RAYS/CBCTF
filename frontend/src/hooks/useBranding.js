import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { DEFAULT_BRANDING, mergeBranding, resolveLocalizedText } from '../config/branding';

export function useBranding() {
  const { i18n, t } = useTranslation();
  const brandingState = useSelector((state) => state.branding?.data);
  const branding = mergeBranding(brandingState || DEFAULT_BRANDING);
  const resolve = (value, fallback = '') =>
    resolveLocalizedText(value, i18n.resolvedLanguage || i18n.language, fallback);

  return {
    branding,
    resolve,
    siteName: resolve(branding.site_name, t('branding.main')),
    adminName: resolve(branding.admin_name, t('branding.admin')),
    browserTitle: resolve(branding.browser_title, 'DEEP DIVE CTF'),
    browserDescription: resolve(branding.browser_description, 'DEEP DIVE CTF competition platform'),
    footerCopyright: resolve(branding.footer_copyright, t('footer.copyright')),
    footerICPNumber: branding.footer_icp_number,
    footerICPLink: branding.footer_icp_link,
    footerContactEmail: branding.footer_contact_email,
    footerGithubURL: branding.footer_github_url,
    homeLogo: branding.home_logo || DEFAULT_BRANDING.home_logo,
    homeLogoAlt: resolve(branding.home_logo_alt, 'DEEP DIVE CTF'),
    home: {
      hero: {
        titlePrefix: resolve(branding.home?.hero?.title_prefix, t('home.hero.titlePrefix')),
        titleHighlight: resolve(branding.home?.hero?.title_highlight, t('home.hero.titleHighlight')),
        titleSuffix: resolve(branding.home?.hero?.title_suffix, t('home.hero.titleSuffix')),
        subtitle: resolve(branding.home?.hero?.subtitle, t('home.hero.subtitle')),
        primaryAction: resolve(branding.home?.hero?.primary_action, t('home.hero.start')),
        secondaryAction: resolve(branding.home?.hero?.secondary_action, t('home.hero.learnMore')),
      },
      challengeTypes: {
        titlePrefix: resolve(branding.home?.challenge_types?.title_prefix, t('home.challengeTypes.titlePrefix')),
        titleHighlight: resolve(
          branding.home?.challenge_types?.title_highlight,
          t('home.challengeTypes.titleHighlight')
        ),
        subtitle: resolve(branding.home?.challenge_types?.subtitle, t('home.challengeTypes.subtitle')),
      },
      upcoming: {
        titlePrefix: resolve(branding.home?.upcoming?.title_prefix, t('home.upcoming.titlePrefix')),
        titleHighlight: resolve(branding.home?.upcoming?.title_highlight, t('home.upcoming.titleHighlight')),
        subtitle: resolve(branding.home?.upcoming?.subtitle, t('home.upcoming.subtitle')),
        action: resolve(branding.home?.upcoming?.action, t('home.upcoming.viewAll')),
      },
      leaderboard: {
        titlePrefix: resolve(branding.home?.leaderboard?.title_prefix, t('home.leaderboard.titlePrefix')),
        titleHighlight: resolve(branding.home?.leaderboard?.title_highlight, t('home.leaderboard.titleHighlight')),
        subtitle: resolve(branding.home?.leaderboard?.subtitle, t('home.leaderboard.subtitle')),
        action: resolve(branding.home?.leaderboard?.action, t('home.leaderboard.viewScoreboard')),
      },
    },
  };
}
