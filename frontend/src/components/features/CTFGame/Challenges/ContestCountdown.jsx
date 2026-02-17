/**
 * 比赛倒计时组件
 * @param {Object} props
 * @param {string} props.startTime - 比赛开始时间的ISO字符串
 * @param {boolean} props.joined - 是否已加入比赛
 * @param {Function} props.onJoin - 加入比赛的回调函数
 */

import { useState, useEffect } from 'react';
import { Button, Card } from '../../../common';

function ContestCountdown({ startTime, joined, onJoin }) {
  const [timeLeft, setTimeLeft] = useState({
    days: 0,
    hours: 0,
    minutes: 0,
    seconds: 0,
  });

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

  return (
    <Card variant="default" padding="lg" animate>
      <div className="flex flex-col items-center justify-center space-y-8">
        {/* 倒计时标题 */}
        <div className="text-center">
          <h2 className="text-2xl text-neutral-50 font-mono mb-2">Contest Starting In</h2>
          <p className="text-neutral-400">Get ready for the challenge</p>
        </div>

        {/* 倒计时显示 */}
        <div className="flex gap-8">
          {[
            { label: 'DAYS', value: timeLeft.days },
            { label: 'HOURS', value: timeLeft.hours },
            { label: 'MINUTES', value: timeLeft.minutes },
            { label: 'SECONDS', value: timeLeft.seconds },
          ].map((item, index) => (
            <div key={index} className="flex flex-col items-center">
              <div
                className="w-[100px] h-[100px] border border-geek-400 rounded-md 
                                bg-black/50 flex items-center justify-center
                                text-4xl text-geek-400 font-mono
                                shadow-[0_0_15px_rgba(89,126,247,0.1)]"
              >
                {String(item.value).padStart(2, '0')}
              </div>
              <span className="mt-2 text-neutral-400 text-sm font-mono">{item.label}</span>
            </div>
          ))}
        </div>

        {/* 加入状态 */}
        <div className="flex flex-col items-center space-y-4">
          {joined ? (
            <div
              className="px-8 py-2 border border-green-400 rounded-md
                            bg-green-400/5 text-green-400 font-mono
                            shadow-[0_0_15px_rgba(74,222,128,0.1)]"
            >
              JOINED
            </div>
          ) : (
            <Button variant="primary" size="lg" className="!text-lg" onClick={onJoin}>
              JOIN NOW
            </Button>
          )}
        </div>
      </div>
    </Card>
  );
}

export default ContestCountdown;
