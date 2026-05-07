export const DEFAULT_BRANDING = {
  code: 'default',
  site_name: { zh_cn: 'CBCTF', en: 'CBCTF' },
  admin_name: { zh_cn: 'DASHBOARD', en: 'DASHBOARD' },
  browser_title: { zh_cn: 'CBCTF', en: 'CBCTF' },
  browser_description: {
    zh_cn: 'DEEP DIVE CTF competition platform',
    en: 'DEEP DIVE CTF competition platform',
  },
  footer_copyright: { zh_cn: '© 2026 CBCTF', en: '© 2026 CBCTF' },
  footer_icp_number: '陕ICP备2025076339号-1',
  footer_icp_link: 'https://beian.miit.gov.cn/',
  footer_contact_email: 'support@0rays.club',
  footer_github_url: 'https://github.com/0RAYS/CBCTF',
  home_logo: '/platform/logo.png',
  home_logo_alt: { zh_cn: 'CBCTF 首页 Logo', en: 'CBCTF home logo' },
  home: {
    hero: {
      title_prefix: { zh_cn: 'Dive Deep into the', en: 'Dive Deep into the' },
      title_highlight: { zh_cn: 'Cyber Security', en: 'Cyber Security' },
      title_suffix: { zh_cn: 'Challenge', en: 'Challenge' },
      subtitle: {
        zh_cn:
          'Join the elite community of hackers, compete in real-world challenges, and master the art of cybersecurity through hands-on experience.',
        en: 'Join the elite community of hackers, compete in real-world challenges, and master the art of cybersecurity through hands-on experience.',
      },
      primary_action: { zh_cn: 'START HACKING', en: 'START HACKING' },
      secondary_action: { zh_cn: 'LEARN MORE', en: 'LEARN MORE' },
    },
    challenge_types: {
      title_prefix: { zh_cn: 'Master All Aspects of', en: 'Master All Aspects of' },
      title_highlight: { zh_cn: 'Cyber Security', en: 'Cyber Security' },
      subtitle: {
        zh_cn: 'From web exploitation to binary analysis, our challenges cover every major domain of cybersecurity',
        en: 'From web exploitation to binary analysis, our challenges cover every major domain of cybersecurity',
      },
    },
    upcoming: {
      title_prefix: { zh_cn: 'Upcoming', en: 'Upcoming' },
      title_highlight: { zh_cn: 'Competitions', en: 'Competitions' },
      subtitle: {
        zh_cn: 'Register now for our upcoming CTF events and compete with the best',
        en: 'Register now for our upcoming CTF events and compete with the best',
      },
      action: { zh_cn: 'VIEW ALL CONTESTS', en: 'VIEW ALL CONTESTS' },
    },
    leaderboard: {
      title_prefix: { zh_cn: 'Top', en: 'Top' },
      title_highlight: { zh_cn: 'Performers', en: 'Performers' },
      subtitle: {
        zh_cn: 'Meet the elite teams leading our global leaderboard',
        en: 'Meet the elite teams leading our global leaderboard',
      },
      action: { zh_cn: 'VIEW SCOREBOAR', en: 'VIEW SCOREBOARD' },
    },
  },
};

export function getBrandingLanguageKey(language) {
  return language?.toLowerCase().startsWith('zh') ? 'zh_cn' : 'en';
}

export function resolveLocalizedText(value, language, fallback = '') {
  if (!value || typeof value !== 'object') return fallback;
  const languageKey = getBrandingLanguageKey(language);
  return value[languageKey] || value.en || value.zh_cn || fallback;
}

export function mergeBranding(branding = {}) {
  return {
    ...DEFAULT_BRANDING,
    ...branding,
    site_name: { ...DEFAULT_BRANDING.site_name, ...(branding.site_name || {}) },
    admin_name: { ...DEFAULT_BRANDING.admin_name, ...(branding.admin_name || {}) },
    browser_title: { ...DEFAULT_BRANDING.browser_title, ...(branding.browser_title || {}) },
    browser_description: { ...DEFAULT_BRANDING.browser_description, ...(branding.browser_description || {}) },
    footer_copyright: { ...DEFAULT_BRANDING.footer_copyright, ...(branding.footer_copyright || {}) },
    footer_icp_number: branding.footer_icp_number || DEFAULT_BRANDING.footer_icp_number,
    footer_icp_link: branding.footer_icp_link || DEFAULT_BRANDING.footer_icp_link,
    footer_contact_email: branding.footer_contact_email || DEFAULT_BRANDING.footer_contact_email,
    footer_github_url: branding.footer_github_url || DEFAULT_BRANDING.footer_github_url,
    home_logo_alt: { ...DEFAULT_BRANDING.home_logo_alt, ...(branding.home_logo_alt || {}) },
    home: {
      ...DEFAULT_BRANDING.home,
      ...(branding.home || {}),
      hero: {
        ...DEFAULT_BRANDING.home.hero,
        ...(branding.home?.hero || {}),
        title_prefix: {
          ...DEFAULT_BRANDING.home.hero.title_prefix,
          ...(branding.home?.hero?.title_prefix || {}),
        },
        title_highlight: {
          ...DEFAULT_BRANDING.home.hero.title_highlight,
          ...(branding.home?.hero?.title_highlight || {}),
        },
        title_suffix: {
          ...DEFAULT_BRANDING.home.hero.title_suffix,
          ...(branding.home?.hero?.title_suffix || {}),
        },
        subtitle: {
          ...DEFAULT_BRANDING.home.hero.subtitle,
          ...(branding.home?.hero?.subtitle || {}),
        },
        primary_action: {
          ...DEFAULT_BRANDING.home.hero.primary_action,
          ...(branding.home?.hero?.primary_action || {}),
        },
        secondary_action: {
          ...DEFAULT_BRANDING.home.hero.secondary_action,
          ...(branding.home?.hero?.secondary_action || {}),
        },
      },
      challenge_types: {
        ...DEFAULT_BRANDING.home.challenge_types,
        ...(branding.home?.challenge_types || {}),
        title_prefix: {
          ...DEFAULT_BRANDING.home.challenge_types.title_prefix,
          ...(branding.home?.challenge_types?.title_prefix || {}),
        },
        title_highlight: {
          ...DEFAULT_BRANDING.home.challenge_types.title_highlight,
          ...(branding.home?.challenge_types?.title_highlight || {}),
        },
        subtitle: {
          ...DEFAULT_BRANDING.home.challenge_types.subtitle,
          ...(branding.home?.challenge_types?.subtitle || {}),
        },
      },
      upcoming: {
        ...DEFAULT_BRANDING.home.upcoming,
        ...(branding.home?.upcoming || {}),
        title_prefix: {
          ...DEFAULT_BRANDING.home.upcoming.title_prefix,
          ...(branding.home?.upcoming?.title_prefix || {}),
        },
        title_highlight: {
          ...DEFAULT_BRANDING.home.upcoming.title_highlight,
          ...(branding.home?.upcoming?.title_highlight || {}),
        },
        subtitle: {
          ...DEFAULT_BRANDING.home.upcoming.subtitle,
          ...(branding.home?.upcoming?.subtitle || {}),
        },
        action: {
          ...DEFAULT_BRANDING.home.upcoming.action,
          ...(branding.home?.upcoming?.action || {}),
        },
      },
      leaderboard: {
        ...DEFAULT_BRANDING.home.leaderboard,
        ...(branding.home?.leaderboard || {}),
        title_prefix: {
          ...DEFAULT_BRANDING.home.leaderboard.title_prefix,
          ...(branding.home?.leaderboard?.title_prefix || {}),
        },
        title_highlight: {
          ...DEFAULT_BRANDING.home.leaderboard.title_highlight,
          ...(branding.home?.leaderboard?.title_highlight || {}),
        },
        subtitle: {
          ...DEFAULT_BRANDING.home.leaderboard.subtitle,
          ...(branding.home?.leaderboard?.subtitle || {}),
        },
        action: {
          ...DEFAULT_BRANDING.home.leaderboard.action,
          ...(branding.home?.leaderboard?.action || {}),
        },
      },
    },
  };
}
