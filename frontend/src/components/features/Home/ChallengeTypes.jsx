import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';

function ChallengeTypes() {
  const { t } = useTranslation();
  const types = [
    {
      name: t('home.challengeTypes.web.name'),
      icon: '🌐',
      description: t('home.challengeTypes.web.description'),
      color: 'border-yellow-400',
      hoverColor: 'hover:border-yellow-400',
      textColor: 'text-yellow-400',
    },
    {
      name: t('home.challengeTypes.crypto.name'),
      icon: '🔐',
      description: t('home.challengeTypes.crypto.description'),
      color: 'border-green-400',
      hoverColor: 'hover:border-green-400',
      textColor: 'text-green-400',
    },
    {
      name: t('home.challengeTypes.forensics.name'),
      icon: '🔍',
      description: t('home.challengeTypes.forensics.description'),
      color: 'border-purple-400',
      hoverColor: 'hover:border-purple-400',
      textColor: 'text-purple-400',
    },
    {
      name: t('home.challengeTypes.binary.name'),
      icon: '💻',
      description: t('home.challengeTypes.binary.description'),
      color: 'border-red-400',
      hoverColor: 'hover:border-red-400',
      textColor: 'text-red-400',
    },
    {
      name: t('home.challengeTypes.reverse.name'),
      icon: '🚩',
      description: t('home.challengeTypes.reverse.description'),
      color: 'border-geek-400',
      hoverColor: 'hover:border-geek-400',
      textColor: 'text-geek-400',
    },
    {
      name: t('home.challengeTypes.blockchain.name'),
      icon: '💰',
      description: t('home.challengeTypes.blockchain.description'),
      color: 'border-gray-400',
      hoverColor: 'hover:border-gray-400',
      textColor: 'text-gray-400',
    },
  ];

  return (
    <div className="py-20 px-8 bg-black/30">
      <div className="w-full max-w-[1200px] mx-auto">
        {/* 标题 */}
        <div className="text-center mb-12">
          <motion.h2
            className="text-3xl font-mono text-neutral-50 mb-4"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
          >
            {t('home.challengeTypes.titlePrefix')}
            <span className="text-geek-400"> {t('home.challengeTypes.titleHighlight')}</span>
          </motion.h2>
          <motion.p
            className="text-neutral-300 max-w-[600px] mx-auto"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
          >
            {t('home.challengeTypes.subtitle')}
          </motion.p>
        </div>

        {/* 类型卡片 */}
        <div className="grid grid-cols-2 gap-6">
          {types.map((type, index) => (
            <motion.div
              key={index}
              className={`p-6 border ${type.color} rounded-md bg-neutral-900
                                hover:bg-neutral-800 transition-all duration-200 cursor-pointer group`}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1 }}
              whileHover={{ y: -5 }}
            >
              <div className="flex items-start gap-4">
                <span className="text-4xl">{type.icon}</span>
                <div>
                  <h3 className={`text-xl font-mono mb-2 ${type.textColor}`}>{type.name}</h3>
                  <p className="text-neutral-300">{type.description}</p>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default ChallengeTypes;
