import { motion } from 'motion/react';
import { Button } from '../../../components/common';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

function HeroSection() {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <div className="min-h-[600px] flex items-center justify-center px-4 md:px-8">
      <div className="w-full max-w-[1200px] flex flex-col md:flex-row items-center md:justify-between gap-8">
        {/* 左侧文本区域 */}
        <motion.div
          className="max-w-[600px] space-y-6"
          initial={{ opacity: 0, x: -50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5 }}
        >
          <h1 className="text-5xl font-mono text-neutral-50 leading-tight">
            {t('home.hero.titlePrefix')}
            <span className="text-geek-400"> {t('home.hero.titleHighlight')} </span>
            {t('home.hero.titleSuffix')}
          </h1>
          <p className="text-neutral-300 text-lg">{t('home.hero.subtitle')}</p>
          <div className="flex gap-4">
            <Button variant="primary" size="lg" className="shadow-focus-strong" onClick={() => navigate('/games')}>
              {t('home.hero.start')}
            </Button>
            <Button variant="outline" size="lg" onClick={() => navigate('/support')}>
              {t('home.hero.learnMore')}
            </Button>
          </div>
        </motion.div>

        {/* 右侧装饰图形 — 仅在 md+ 显示 */}
        <motion.div
          className="relative w-[400px] h-[400px] hidden md:block"
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          {/* 这里可以添加一些科技感的SVG动画或3D模型 */}
          <div className="absolute inset-0 border border-neutral-300/30 rounded-md">
            <img alt="CBCTF" src="./logo.png" width="400" height="400" />
          </div>
        </motion.div>
      </div>
    </div>
  );
}

export default HeroSection;
