export const REPO_URL = 'https://github.com/0RAYS/CBCTF';

export function getFooterConfig(t, footerBranding = {}) {
  const contactEmail = footerBranding.contactEmail || 'support@0rays.club';

  return {
    copyright: footerBranding.copyright || t('footer.copyright'),
    icp: {
      number: footerBranding.icpNumber || t('footer.icp'),
      link: footerBranding.icpLink || 'https://beian.miit.gov.cn/',
    },
    links: [
      { label: t('footer.support'), href: '/support', isExternal: false },
      { label: contactEmail, href: `mailto:${contactEmail}`, isExternal: true },
      { label: t('footer.github'), href: footerBranding.githubURL || REPO_URL, isExternal: true },
    ],
  };
}
