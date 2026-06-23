export const REPO_URL = 'https://github.com/0RAYS/CBCTF';

export function getFooterConfig(t, footerBranding = {}) {
  const branding = footerBranding && typeof footerBranding === 'object' ? footerBranding : {};
  const contactEmail = branding.contactEmail ?? 'support@0rays.club';
  const icpNumber = branding.icpNumber ?? t('footer.icp');
  const icpLink = branding.icpLink ?? 'https://beian.miit.gov.cn/';
  const githubURL = branding.githubURL ?? REPO_URL;

  return {
    copyright: branding.copyright ?? t('footer.copyright'),
    icp: icpNumber ? { number: icpNumber, link: icpLink } : null,
    links: [
      contactEmail ? { label: contactEmail, href: `mailto:${contactEmail}`, isExternal: true } : null,
      githubURL ? { label: t('footer.github'), href: githubURL, isExternal: true } : null,
    ].filter(Boolean),
  };
}
