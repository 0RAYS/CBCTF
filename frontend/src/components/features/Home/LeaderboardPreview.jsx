import { motion } from 'motion/react';
import { ScrollingText } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function LeaderboardPreview(props) {
  const topUsers = props?.topUsers;
  const { t } = useTranslation();

  return (
    <div className="py-20 px-8 bg-black/30">
      <div className="w-full max-w-[1200px] mx-auto">
        {/* 标题 */}
        <motion.div
          className="text-center mb-12"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
        >
          <h2 className="text-3xl font-mono text-neutral-50 mb-4">
            {t('home.leaderboard.titlePrefix')}{' '}
            <span className="text-geek-400">{t('home.leaderboard.titleHighlight')}</span>
          </h2>
          <p className="text-neutral-300">{t('home.leaderboard.subtitle')}</p>
        </motion.div>

        {/* 排行榜 */}
        <motion.div
          className="border border-neutral-300 rounded-md overflow-hidden"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
        >
          {topUsers.length > 0 &&
            topUsers.map((team, index) => (
              <motion.div
                key={index}
                className={`flex items-center justify-between p-6 bg-black/30 backdrop-blur-[2px]
                                ${index !== topUsers.length - 1 ? 'border-b border-neutral-300/30' : ''}
                                hover:bg-black/50 transition-all duration-200 cursor-pointer group`}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
              >
                {/* 排名和队伍信息 */}
                <div className="flex items-center gap-6">
                  <span
                    className={`text-2xl font-mono ${
                      index === 0
                        ? 'text-yellow-400'
                        : index === 1
                          ? 'text-neutral-300'
                          : index === 2
                            ? 'text-yellow-600'
                            : 'text-neutral-400'
                    }`}
                  >
                    #{index + 1}
                  </span>
                  <div className="flex items-center gap-2">
                    <ScrollingText
                      text={team.name}
                      className="text-neutral-50 font-mono group-hover:text-geek-400 transition-colors duration-200"
                      maxWidth={200}
                      speed={15}
                    />
                  </div>
                </div>

                {/* 分数和解题数 */}
                <div className="flex items-center gap-8">
                  <div className="flex items-center gap-2">
                    <span className="text-neutral-400">{t('common.score')}:</span>
                    <span className="text-geek-400 font-mono">{team.score}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-neutral-400">{t('common.solved')}:</span>
                    <span className="text-neutral-50 font-mono">{team.solved}</span>
                  </div>
                </div>
              </motion.div>
            ))}
        </motion.div>

        {/* 查看更多按钮 */}
        {/*<motion.div*/}
        {/*  className="flex justify-center mt-8"*/}
        {/*  initial={{ opacity: 0 }}*/}
        {/*  whileInView={{ opacity: 1 }}*/}
        {/*  viewport={{ once: true }}*/}
        {/*  transition={{ delay: 0.5 }}*/}
        {/*>*/}
        {/*  <motion.button*/}
        {/*    className="px-8 h-[40px] border border-neutral-300 rounded-md*/}
        {/*                    text-neutral-300 font-mono tracking-wider*/}
        {/*                    hover:border-geek-400 hover:text-geek-400 hover:bg-geek-400/10*/}
        {/*                    transition-all duration-200"*/}
        {/*    whileHover={{ scale: 1.05 }}*/}
        {/*    whileTap={{ scale: 0.95 }}*/}
        {/*  >*/}
        {/*    VIEW FULL LEADERBOARD*/}
        {/*  </motion.button>*/}
        {/*</motion.div>*/}
      </div>
    </div>
  );
}

export default LeaderboardPreview;
