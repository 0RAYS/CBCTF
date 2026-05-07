import NavBar from './NavBar';
import Footer from './Footer';
import { useTranslation } from 'react-i18next';
import { getFooterConfig } from '../../../config/footer';
import { useBranding } from '../../../hooks/useBranding';

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
  backgroundImage = '',
}) {
  const { t } = useTranslation();
  const { footerCopyright, footerICPNumber, footerICPLink, footerContactEmail, footerGithubURL } = useBranding();
  const resolvedLogo = logo || t('branding.main');
  const footerConfig = getFooterConfig(t, {
    copyright: footerCopyright,
    icpNumber: footerICPNumber,
    icpLink: footerICPLink,
    contactEmail: footerContactEmail,
    githubURL: footerGithubURL,
  });

  return (
    <div className="h-full w-full overflow-x-hidden">
      <div className="fixed inset-0 bg-neutral-900">
        {backgroundImage && (
          <>
            <div
              className="absolute inset-0 bg-cover bg-center bg-no-repeat opacity-45"
              style={{ backgroundImage: `url(${backgroundImage})` }}
              aria-hidden="true"
            />
            <div className="absolute inset-0 bg-neutral-900/65 backdrop-blur-[2px]" aria-hidden="true" />
            <div
              className="absolute inset-0 bg-gradient-to-b from-neutral-900/45 via-neutral-900/55 to-neutral-900/90"
              aria-hidden="true"
            />
          </>
        )}
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
      <main id="main-content" className="w-full min-h-full pt-[110px] pb-[80px] px-4 md:px-8 relative z-1">
        {children}
      </main>

      <div className="relative z-2">
        <Footer {...footerConfig} />
      </div>
    </div>
  );
}

export default BaseLayout;
