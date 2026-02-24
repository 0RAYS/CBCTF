export const REPO_URL = 'https://github.com/0RAYS/CBCTF';

export function getFooterConfig(t) {
  return {
    copyright: t('footer.copyright'),
    icp: {
      number: t('footer.icp'),
      link: 'https://beian.miit.gov.cn/',
    },
    links: [
      { label: t('footer.support'), href: '/support', isExternal: false },
      { label: t('footer.contact'), href: '/contact', isExternal: false },
      { label: t('footer.github'), href: REPO_URL, isExternal: true },
    ],
  };
}
