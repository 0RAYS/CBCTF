import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { IconWorld, IconLock, IconMicroscope, IconTerminal2, IconRotateClockwise, IconLink } from '@tabler/icons-react';
import { useBranding } from '../../../hooks/useBranding';

function ChallengeTypes() {
  const { t } = useTranslation();
  const { home } = useBranding();
  const types = [
    {
      name: t('home.challengeTypes.web.name'),
      Icon: IconWorld,
      description: t('home.challengeTypes.web.description'),
      iconColor: 'text-geek-400',
    },
    {
      name: t('home.challengeTypes.crypto.name'),
      Icon: IconLock,
      description: t('home.challengeTypes.crypto.description'),
      iconColor: 'text-geek-400',
    },
    {
      name: t('home.challengeTypes.forensics.name'),
      Icon: IconMicroscope,
      description: t('home.challengeTypes.forensics.description'),
      iconColor: 'text-green-400',
    },
    {
      name: t('home.challengeTypes.binary.name'),
      Icon: IconTerminal2,
      description: t('home.challengeTypes.binary.description'),
      iconColor: 'text-geek-400',
    },
    {
      name: t('home.challengeTypes.reverse.name'),
      Icon: IconRotateClockwise,
      description: t('home.challengeTypes.reverse.description'),
      iconColor: 'text-yellow-400',
    },
    {
      name: t('home.challengeTypes.blockchain.name'),
      Icon: IconLink,
      description: t('home.challengeTypes.blockchain.description'),
      iconColor: 'text-green-400',
    },
  ];

  return (
    <div className="py-14 md:py-20 px-4 md:px-8 bg-neutral-800/20 border-y border-neutral-700/40">
      <div className="w-full max-w-[1200px] mx-auto">
        <motion.div
          className="mb-10"
          initial={{ opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ ease: [0.16, 1, 0.3, 1], duration: 0.4 }}
        >
          <div className="flex items-center gap-4 mb-3">
            <div className="w-6 h-[2px] bg-geek-400" />
            <span className="text-xs font-mono text-geek-400 tracking-[0.2em] uppercase">
              {home.challengeTypes.titlePrefix}
            </span>
          </div>
          <h2 className="text-2xl md:text-3xl font-mono text-neutral-50">{home.challengeTypes.titleHighlight}</h2>
          <p className="text-neutral-400 text-sm mt-2 max-w-[60ch]">{home.challengeTypes.subtitle}</p>
        </motion.div>

        {/* Type cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          {types.map((type, index) => (
            <motion.div
              key={index}
              className="p-6 border border-neutral-600/50 rounded-md bg-neutral-800/40
                         hover:bg-neutral-800/60 hover:border-geek-400/50 transition-colors duration-200 cursor-pointer group"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: index * 0.1, ease: [0.25, 1, 0.5, 1], duration: 0.4 }}
            >
              <div className="flex items-start gap-4">
                <type.Icon className={`w-8 h-8 shrink-0 mt-0.5 ${type.iconColor}`} />
                <div>
                  <h3 className="text-xl font-mono mb-2 text-neutral-50 group-hover:text-geek-400 transition-colors duration-200">
                    {type.name}
                  </h3>
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
