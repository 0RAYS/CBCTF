export const DEFAULT_BRANDING = {
  code: 'default',
  site_name: { zh_cn: '深潜 CTF', en: 'DEEP DIVE CTF' },
  admin_name: { zh_cn: '深潜管理台', en: 'DEEP DIVE Admin' },
  browser_title: { zh_cn: '深潜 CTF', en: 'DEEP DIVE CTF' },
  browser_description: {
    zh_cn: '深潜 CTF 网络安全竞赛平台',
    en: 'DEEP DIVE CTF competition platform',
  },
  footer_copyright: { zh_cn: '© 2025 深潜 CTF', en: '© 2025 DEEP DIVE CTF' },
  home_logo: '/platform/logo.png',
  home_logo_alt: { zh_cn: '深潜 CTF 首页 Logo', en: 'DEEP DIVE CTF home logo' },
  home: {
    hero: {
      title_prefix: { zh_cn: '深入探索', en: 'Dive Deep into the' },
      title_highlight: { zh_cn: '网络安全', en: 'Cyber Security' },
      title_suffix: { zh_cn: '挑战', en: 'Challenge' },
      subtitle: {
        zh_cn: '加入深潜 CTF 社区，在真实场景中对抗、练习与成长。',
        en: 'Join the elite community of hackers, compete in real-world challenges, and master the art of cybersecurity through hands-on experience.',
      },
      primary_action: { zh_cn: '立即参赛', en: 'START HACKING' },
      secondary_action: { zh_cn: '了解更多', en: 'LEARN MORE' },
    },
    challenge_types: {
      title_prefix: { zh_cn: '掌握', en: 'Master All Aspects of' },
      title_highlight: { zh_cn: '网络安全', en: 'Cyber Security' },
      subtitle: {
        zh_cn: '从 Web 渗透到二进制分析，覆盖网络安全核心方向。',
        en: 'From web exploitation to binary analysis, our challenges cover every major domain of cybersecurity',
      },
    },
    upcoming: {
      title_prefix: { zh_cn: '近期', en: 'Upcoming' },
      title_highlight: { zh_cn: '赛事', en: 'Competitions' },
      subtitle: {
        zh_cn: '报名即将开启的 CTF 赛事，与高手同场竞技。',
        en: 'Register now for our upcoming CTF events and compete with the best',
      },
      action: { zh_cn: '查看全部赛事', en: 'VIEW ALL CONTESTS' },
    },
    leaderboard: {
      title_prefix: { zh_cn: '顶尖', en: 'Top' },
      title_highlight: { zh_cn: '战队', en: 'Performers' },
      subtitle: {
        zh_cn: '查看当前积分榜上的领先队伍。',
        en: 'Meet the elite teams leading our global leaderboard',
      },
      action: { zh_cn: '查看排行榜', en: 'VIEW SCOREBOARD' },
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
