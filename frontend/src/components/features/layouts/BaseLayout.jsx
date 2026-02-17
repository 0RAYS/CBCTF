import NavBar from './NavBar';
import Footer from './Footer';
import { useTranslation } from 'react-i18next';

function BaseLayout({
  children,
  tabs = [{ id: 'GAMES', label: 'GAMES' }],
  activeTab = 'GAMES',
  onTabChange = () => {},
  onLogoClick = () => {},
  onPictureClick = () => {},
  logo = '',
  pictureSrc = '',
  userName = '',
}) {
  const { t } = useTranslation();
  const resolvedLogo = logo || t('branding.main');
  const footerConfig = {
    copyright: t('footer.copyright'),
    icp: {
      number: t('footer.icp'),
      link: 'https://beian.miit.gov.cn/',
    },
    links: [
      { label: t('footer.support'), href: '/support', isExternal: false },
      { label: t('footer.contact'), href: '/contact', isExternal: false },
      { label: t('footer.github'), href: 'https://github.com/0RAYS/CBCTF', isExternal: true },
    ],
  };

  return (
    <div className="h-full w-full overflow-x-hidden">
      <div className="fixed top-0 left-0 w-full h-full bg-black">
        {/* <Squares
          speed={0.03}
          hoverFillColor="#434343"
          borderColor="#606060"
          direction="down"
          gradientConfig={{
            enabled: true,
            innerColor: 'rgba(0, 0, 0, 0)',
            outerColor: '#060606',
            opacity: 0.8,
          }}
        /> */}
      </div>
      <div className="relative z-50">
        <NavBar
          tabs={tabs}
          activeTab={activeTab}
          onTabChange={onTabChange}
          onLogoClick={onLogoClick}
          onPictureClick={onPictureClick}
          logo={resolvedLogo}
          pictureSrc={pictureSrc}
          userName={userName}
        />
      </div>

      {/* 内容区域 */}
      <div className="w-full min-h-full pt-[110px] pb-[80px] px-4 md:px-8 relative z-1">{children}</div>

      <div className="relative z-2">
        <Footer {...footerConfig} />
      </div>
    </div>
  );
}

export default BaseLayout;
