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
}) {
  const { t } = useTranslation();
  const { footerCopyright } = useBranding();
  const resolvedLogo = logo || t('branding.main');
  const footerConfig = getFooterConfig(t, footerCopyright);

  return (
    <div className="h-full w-full overflow-x-hidden">
      {/* Skip-to-content — keyboard accessibility */}
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-[9999] focus:px-4 focus:py-2 focus:bg-neutral-800 focus:border focus:border-geek-400 focus:rounded-md focus:text-geek-400 focus:font-mono focus:text-sm focus:outline-none"
      >
        Skip to main content
      </a>
      <div className="fixed top-0 left-0 w-full h-full bg-neutral-900">
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
      <main id="main-content" className="w-full min-h-full pt-[110px] pb-[80px] px-4 md:px-8 relative z-1">{children}</main>

      <div className="relative z-2">
        <Footer {...footerConfig} />
      </div>
    </div>
  );
}

export default BaseLayout;
