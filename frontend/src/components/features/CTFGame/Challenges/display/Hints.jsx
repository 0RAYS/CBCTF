import { useState } from 'react';
import { useTranslation } from 'react-i18next';

function HintItem({ hint, index }) {
  const [isExpanded, setIsExpanded] = useState(false);
  const { t } = useTranslation();
  return (
    <div
      className={`border border-neutral-300/30 rounded-md overflow-hidden transition-colors duration-200 ${
        isExpanded ? 'bg-neutral-900/50' : 'bg-black/20'
      }`}
    >
      <button
        type="button"
        aria-expanded={isExpanded}
        className="flex w-full items-center justify-between p-2 text-left cursor-pointer hover:bg-neutral-800/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-geek-400"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center gap-2">
          <span className="text-geek-400 font-mono text-sm">#{index + 1}</span>
          <span className="text-neutral-300 font-mono text-sm">
            {isExpanded ? t('game.challengeModal.hint.hide') : t('game.challengeModal.hint.show')}
          </span>
        </div>
        <span className="text-neutral-400 text-[10px]" aria-hidden="true">
          {'\u25bc'}
        </span>
      </button>
      <div hidden={!isExpanded}>
        <div className="p-3 border-t border-neutral-300/10">
          <span className="break-words text-neutral-400 font-mono text-sm">{hint}</span>
        </div>
      </div>
    </div>
  );
}

export default function Hints({ challenge }) {
  const { t } = useTranslation();
  if (!challenge.hints?.length) return null;
  return (
    <div className="space-y-1.5">
      <h3 className="text-neutral-400 font-mono text-sm">{t('game.challengeModal.sections.hints')}</h3>
      <div className="space-y-2">
        {challenge.hints.map((hint, index) => (
          <HintItem key={`hint-${challenge.id}-${index}`} hint={hint} index={index} />
        ))}
      </div>
    </div>
  );
}
