import { motion } from 'motion/react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

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
    { label: t('footer.github'), href: 'https://github.com/your-repo', isExternal: true },
  ];
  // 渲染链接的辅助函数
  const renderLink = (link, index) => {
    if (link.isExternal) {
      return (
        <motion.a
          key={index}
          href={link.href}
          target="_blank"
          rel="noopener noreferrer"
          className="text-neutral-300 text-sm tracking-wider hover:text-geek-400 transition-colors duration-200"
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          {link.label}
        </motion.a>
      );
    }

    return (
      <motion.div key={index} whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
        <Link
          to={link.href}
          className="text-neutral-300 text-sm tracking-wider hover:text-geek-400 transition-colors duration-200"
        >
          {link.label}
        </Link>
      </motion.div>
    );
  };

  return (
    <motion.div
      className="fixed bottom-0 left-0 w-full h-[60px] bg-black/30 backdrop-blur-[2px] border-t border-neutral-300"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2 }}
    >
      <div className="h-full max-w-[1200px] mx-auto px-8 flex items-center justify-between">
        {/* 左侧装饰和备案信息 */}
        <div className="relative flex items-center">
          {/* 装饰线条 */}
          <div className="absolute left-0 top-0 w-[100px] h-[2px] bg-gradient-to-r from-neutral-300 to-transparent"></div>
          <div className="absolute left-0 bottom-0 w-[60px] h-[2px] bg-gradient-to-r from-neutral-300 to-transparent"></div>

          {/* 备案信息 */}
          <div className="flex items-center space-x-4 ml-8">
            <span className="text-neutral-300 text-sm tracking-wider font-mono">{resolvedCopyright}</span>
            {resolvedIcp && (
              <>
                <div className="w-[1px] h-[14px] bg-neutral-300/30"></div>
                <a
                  href={resolvedIcp.link}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-neutral-300 text-sm tracking-wider hover:text-geek-400 transition-colors duration-200"
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
            <div className="absolute right-0 top-0 w-[80px] h-[2px] bg-gradient-to-l from-neutral-300 to-transparent"></div>
            <div className="absolute right-0 bottom-0 w-[40px] h-[2px] bg-gradient-to-l from-neutral-300 to-transparent"></div>

            {/* 链接区域 */}
            <div className="flex items-center space-x-6 mr-8">
              {resolvedLinks.map((link, index) => renderLink(link, index))}
            </div>
          </div>
        )}
      </div>
    </motion.div>
  );
}

export default Footer;
