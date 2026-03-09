import { useState, useEffect } from 'react';
import { Button, Card } from '../../../common';
import { useTranslation } from 'react-i18next';
import { IconSettings, IconBook, IconUsers } from '@tabler/icons-react';

function ContestCountdown({ startTime, joined, onJoin }) {
  const { t } = useTranslation();
  const [timeLeft, setTimeLeft] = useState({ days: 0, hours: 0, minutes: 0, seconds: 0 });

  useEffect(() => {
    const calculateTimeLeft = () => {
      const difference = new Date(startTime) - new Date();
      if (difference > 0) {
        setTimeLeft({
          days: Math.floor(difference / (1000 * 60 * 60 * 24)),
          hours: Math.floor((difference / (1000 * 60 * 60)) % 24),
          minutes: Math.floor((difference / 1000 / 60) % 60),
          seconds: Math.floor((difference / 1000) % 60),
        });
      }
    };
    calculateTimeLeft();
    const timer = setInterval(calculateTimeLeft, 1000);
    return () => clearInterval(timer);
  }, [startTime]);

  const units = [
    { label: t('game.countdown.days'), value: timeLeft.days },
    { label: t('game.countdown.hours'), value: timeLeft.hours },
    { label: t('game.countdown.minutes'), value: timeLeft.minutes },
    { label: t('game.countdown.seconds'), value: timeLeft.seconds },
  ];

  const tips = [
    { Icon: IconSettings, text: t('game.countdown.tip1') },
    { Icon: IconBook, text: t('game.countdown.tip2') },
    { Icon: IconUsers, text: t('game.countdown.tip3') },
  ];

  return (
    <Card variant="default" padding="lg" animate>
      <div className="flex flex-col items-center justify-center space-y-8">
        {/* Title */}
        <div className="text-center">
          <h2 className="text-2xl text-neutral-50 font-mono mb-2">{t('game.countdown.title')}</h2>
          <p className="text-neutral-400">{t('game.countdown.subtitle')}</p>
        </div>

        {/* Countdown digits */}
        <div className="flex flex-wrap justify-center gap-4 sm:gap-8">
          {units.map((item, index) => (
            <div key={index} className="flex flex-col items-center">
              <div
                className="w-[72px] h-[72px] sm:w-[100px] sm:h-[100px] border border-geek-400 rounded-md
                           bg-black/50 flex items-center justify-center
                           text-2xl sm:text-4xl text-geek-400 font-mono shadow-glow-primary"
              >
                {String(item.value).padStart(2, '0')}
              </div>
              <span className="mt-2 text-neutral-400 text-xs sm:text-sm font-mono">{item.label}</span>
            </div>
          ))}
        </div>

        {/* Join status */}
        <div className="flex flex-col items-center">
          {joined ? (
            <div
              className="px-8 py-2 border border-green-400 rounded-md
                         bg-green-400/5 text-green-400 font-mono shadow-glow-success"
            >
              {t('game.countdown.joined')}
            </div>
          ) : (
            <Button variant="primary" size="lg" className="!text-lg" onClick={onJoin}>
              {t('game.countdown.joinNow')}
            </Button>
          )}
        </div>

        {/* While-you-wait tips */}
        <div className="w-full border-t border-neutral-600 pt-6">
          <p className="text-xs uppercase tracking-widest text-neutral-400 mb-4 text-center">
            {t('game.countdown.tipTitle')}
          </p>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            {tips.map((tip, i) => (
              <div key={i} className="flex items-start gap-3 text-sm text-neutral-400">
                <tip.Icon className="w-4 h-4 shrink-0 mt-0.5 text-geek-400" />
                <span>{tip.text}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </Card>
  );
}

export default ContestCountdown;
