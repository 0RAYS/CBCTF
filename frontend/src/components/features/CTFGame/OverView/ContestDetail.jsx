import { motion } from 'motion/react';
import MarkdownContent from '@/components/common/MarkdownContent';
import { useState } from 'react';
import { Button, Card } from '@/components/common';
import { useTranslation } from 'react-i18next';

/** 比赛详情组件 */
function ContestDetail({ contest, handleJoinContest }) {
  const [hoveredPrize, setHoveredPrize] = useState(null);
  const [hoveredTimeline, setHoveredTimeline] = useState(null);
  const { t, i18n } = useTranslation();

  if (!contest) return null;

  return (
    <div className="contest-container mx-auto space-y-6 [overflow-wrap:anywhere]">
      <motion.div
        className="relative w-full min-h-[250px] md:min-h-[300px] border border-neutral-600 rounded-md overflow-hidden bg-black/30"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <div
          className="absolute inset-0 opacity-30 bg-cover bg-center bg-no-repeat"
          style={{ backgroundImage: `url(${contest.image})` }}
        />

        <div className="relative min-w-0 p-4 md:p-8 flex flex-col gap-8">
          <div>
            <div className="flex flex-wrap items-center gap-2 md:gap-4 mb-4">
              <h1 className="min-w-0 text-2xl md:text-4xl font-mono text-neutral-50 tracking-wider">{contest.title}</h1>
              <span
                className={`px-3 py-1 border rounded-md text-sm ${
                  contest.status === 'upcoming'
                    ? 'text-yellow-400 border-yellow-400'
                    : contest.status === 'active'
                      ? 'text-green-400 border-green-400'
                      : contest.status === 'paused'
                        ? 'text-yellow-400 border-yellow-400'
                        : 'text-neutral-400 border-neutral-400'
                }`}
              >
                {t(`game.status.${contest.status}`)}
              </span>
            </div>
            <MarkdownContent className="min-w-0 text-neutral-300 max-w-[800px] text-base sm:text-lg">
              {contest.description}
            </MarkdownContent>
          </div>

          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div className="min-w-0 flex flex-wrap gap-4 md:gap-8">
              <div>
                <span className="text-neutral-400">{t('game.detail.labels.startTime')}</span>
                <div className="text-neutral-50 font-mono">
                  {new Date(contest.startTime).toLocaleString(i18n.language)}
                </div>
              </div>
              <div>
                <span className="text-neutral-400">{t('game.detail.labels.endTime')}</span>
                <div className="text-neutral-50 font-mono">
                  {new Date(contest.endTime).toLocaleString(i18n.language)}
                </div>
              </div>
              <div>
                <span className="text-neutral-400">{t('game.detail.labels.participants')}</span>
                <div className="text-neutral-50 font-mono">{contest.participants ?? t('common.notAvailable')}</div>
              </div>
            </div>

            <Button
              variant="primary"
              size="lg"
              className="w-full lg:w-auto shrink-0 min-h-[50px] h-auto! py-2"
              disabled={contest.status === 'ended'}
              onClick={handleJoinContest}
            >
              {contest.status === 'ended' ? t('game.status.ended') : t('game.team.joinModal.title.select')}
            </Button>
          </div>
        </div>
      </motion.div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="min-w-0 col-span-1 lg:col-span-2" variant="default" padding="lg">
          <h2 className="text-2xl font-mono text-neutral-50 tracking-wider mb-6">
            {t('game.detail.labels.contestRules')}
          </h2>
          <div className="space-y-4 text-neutral-300">
            {contest.rules?.map((rule, index) => (
              <motion.div
                key={index}
                className="flex gap-4 p-2 rounded-md hover:bg-neutral-300/5 transition-colors duration-200 cursor-default"
              >
                <span className="shrink-0 text-geek-400 font-mono">{(index + 1).toString().padStart(2, '0')}</span>
                <p className="min-w-0">{rule}</p>
              </motion.div>
            ))}
          </div>
        </Card>

        <Card variant="default" padding="lg" className="min-w-0">
          <h2 className="text-2xl font-mono text-neutral-50 tracking-wider mb-6">{t('game.detail.labels.prizes')}</h2>
          <div className="space-y-6">
            {contest.prizes?.map((prize, index) => (
              <motion.div
                key={index}
                className={`group p-4 rounded-md transition-all duration-200 cursor-default
                                    ${hoveredPrize === index ? 'bg-neutral-300/5' : 'hover:bg-neutral-300/5'}`}
                onMouseEnter={() => setHoveredPrize(index)}
                onMouseLeave={() => setHoveredPrize(null)}
              >
                <motion.div
                  className="flex flex-wrap items-center gap-3 text-neutral-50 mb-2"
                  animate={{
                    opacity: hoveredPrize === index ? 0.9 : 1,
                    transition: { duration: 0.2 },
                  }}
                >
                  <span
                    className={`text-lg font-mono
                                        ${
                                          index === 0
                                            ? 'text-yellow-400'
                                            : index === 1
                                              ? 'text-neutral-300'
                                              : index === 2
                                                ? 'text-orange-400'
                                                : 'text-neutral-400'
                                        }
                                    `}
                  >
                    {index === 0 ? '1ST' : index === 1 ? '2ND' : index === 2 ? '3RD' : `${index + 1}TH`}
                  </span>
                  <span className="text-xl font-mono">{prize.amount}</span>
                </motion.div>
                <motion.div
                  className="text-neutral-400 text-sm prose prose-sm max-w-none"
                  animate={{
                    opacity: hoveredPrize === index ? 0.9 : 1,
                    transition: { duration: 0.2, delay: 0.1 },
                  }}
                >
                  <MarkdownContent className="prose-sm text-neutral-400">{prize.description}</MarkdownContent>
                </motion.div>
              </motion.div>
            ))}
          </div>
        </Card>
      </div>

      <Card variant="default" padding="lg">
        <h2 className="text-2xl font-mono text-neutral-50 tracking-wider mb-6">{t('game.detail.labels.timeline')}</h2>
        <div className="overflow-x-auto">
          <div className="relative flex flex-col md:flex-row items-stretch gap-6 md:gap-8 md:min-w-max">
            <div className="hidden md:block absolute top-[30px] left-0 right-0 h-[2px] bg-neutral-300/20" />

            {contest.timeline?.map((item, index) => (
              <motion.div
                key={index}
                className={`min-w-0 md:w-64 md:flex-none relative p-4 rounded-md transition-colors duration-200 cursor-default
                                ${hoveredTimeline === index ? 'bg-neutral-300/5' : 'hover:bg-neutral-300/5'}`}
                onMouseEnter={() => setHoveredTimeline(index)}
                onMouseLeave={() => setHoveredTimeline(null)}
              >
                <motion.div
                  animate={{
                    opacity: hoveredTimeline === index ? 0.9 : 1,
                    transition: { duration: 0.2 },
                  }}
                  className="border-none"
                >
                  <div className="text-geek-400 font-mono mb-2">{new Date(item.date).toISOString().slice(0, 10)}</div>
                  <div className="text-neutral-50 font-mono mb-1">{item.title}</div>
                  <MarkdownContent className="text-neutral-400 text-sm prose-sm">{item.description}</MarkdownContent>
                </motion.div>
              </motion.div>
            ))}
          </div>
        </div>
      </Card>
    </div>
  );
}

export default ContestDetail;
