import { motion } from 'motion/react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { REPO_URL } from '../../../config/footer';

function Footer({ copyright, icp, links }) {
  const { t } = useTranslation();

  const resolvedCopyright = copyright ?? t('footer.copyright');
  const resolvedIcp = icp ?? {
    number: t('footer.icp'),
    link: 'https://beian.miit.gov.cn/',
  };
  const resolvedLinks = links ?? [
    { label: t('footer.support'), href: '/support', isExternal: false },
    { label: t('footer.contact'), href: '/contact', isExternal: false },
    { label: t('footer.github'), href: REPO_URL, isExternal: true },
  ];
  // 渲染链接的辅助函数
  const renderLink = (link, index) => {
    if (link.isExternal) {
      return (
        <a
          key={index}
          href={link.href}
          target="_blank"
          rel="noopener noreferrer"
          className="text-neutral-300 text-sm tracking-wider hover:text-geek-400 transition-colors duration-200"
        >
          {link.label}
        </a>
      );
    }

    return (
      <Link
        key={index}
        to={link.href}
        className="text-neutral-300 text-sm tracking-wider hover:text-geek-400 transition-colors duration-200"
      >
        {link.label}
      </Link>
    );
  };

  return (
    <motion.div
      className="fixed bottom-0 left-0 w-full bg-neutral-900/80 backdrop-blur-[4px] border-t border-neutral-600/50"
      style={{ minHeight: '60px' }}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2 }}
    >
      <div className="min-h-[60px] max-w-[1200px] mx-auto px-4 md:px-8 flex flex-col sm:flex-row items-center justify-between gap-2 py-2 sm:py-0">
        {/* 左侧装饰和备案信息 */}
        <div className="relative flex items-center">
          {/* 装饰线条 */}
          <div className="absolute left-0 top-0 w-[100px] h-[1px] bg-gradient-to-r from-neutral-300/60 to-transparent hidden sm:block"></div>
          <div className="absolute left-0 bottom-0 w-[60px] h-[1px] bg-gradient-to-r from-neutral-300/60 to-transparent hidden sm:block"></div>

          {/* 备案信息 */}
          <div className="flex flex-wrap items-center gap-x-4 gap-y-1 sm:ml-8">
            <span className="text-neutral-300 text-xs tracking-wider font-mono">{resolvedCopyright}</span>
            {resolvedIcp && (
              <>
                <div className="w-[1px] h-[12px] bg-neutral-600/60 hidden sm:block"></div>
                <a
                  href={resolvedIcp.link}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-neutral-400 text-xs tracking-wider hover:text-geek-400 transition-colors duration-200"
                >
                  {resolvedIcp.number}
                </a>
              </>
            )}
          </div>
        </div>

        {/* 右侧装饰和链接 */}
        {resolvedLinks.length > 0 && (
          <div className="relative flex items-center">
            {/* 装饰线条 */}
            <div className="absolute right-0 top-0 w-[80px] h-[1px] bg-gradient-to-l from-neutral-300/60 to-transparent hidden sm:block"></div>
            <div className="absolute right-0 bottom-0 w-[40px] h-[1px] bg-gradient-to-l from-neutral-300/60 to-transparent hidden sm:block"></div>

            {/* 链接区域 */}
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1 sm:mr-8">
              {resolvedLinks.map((link, index) => renderLink(link, index))}
            </div>
          </div>
        )}
      </div>
    </motion.div>
  );
}

export default Footer;
