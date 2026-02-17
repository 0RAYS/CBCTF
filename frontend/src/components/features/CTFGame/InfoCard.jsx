import { motion } from 'motion/react';
import { Card } from '../../../components/common';

function InfoCard({
  title,
  value,
  valueColor = 'text-neutral-50',
  subTitle,
  subValue,
  subValueColor,
  showProgress,
  progressValue,
}) {
  return (
    <Card variant="default" padding="md">
      <div className="flex items-center justify-between">
        <div className="text-neutral-400 text-sm">{title}</div>
        <div className={`font-mono text-xl ${valueColor}`}>{value}</div>
      </div>

      {subTitle && subValue && (
        <div className="flex items-center justify-between mt-2">
          <div className="text-neutral-400 text-sm">{subTitle}</div>
          <div className={`font-mono text-xl ${subValueColor}`}>{subValue}</div>
        </div>
      )}

      {showProgress && (
        <div className="h-2 bg-neutral-700 rounded-full overflow-hidden mt-2">
          <motion.div
            className="h-full bg-geek-400"
            initial={{ width: 0 }}
            animate={{ width: `${progressValue}%` }}
            transition={{ duration: 1 }}
          />
        </div>
      )}
    </Card>
  );
}

export default InfoCard;
