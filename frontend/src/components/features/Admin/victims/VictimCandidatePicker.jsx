import { IconTarget, IconSearch, IconArrowsMaximize, IconChevronLeft, IconChevronRight } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../common';
import { getChallengeCategoryChipClass, getChallengeTypeChipClass } from '../../../../config/challengeChips';
import { VictimCandidateDialog } from './VictimStartDialogs';

export default function VictimCandidatePicker({ candidates }) {
  const { t } = useTranslation();
  const { challenges, selectedChallenges, challengeSearch, challengePage, challengePageSize, challengeTotal } =
    candidates;
  return (
    <>
      <div className="flex flex-col min-h-0 flex-1">
        <div className="flex justify-between items-center mb-2">
          <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
            <IconTarget size={14} />
            {t('admin.contests.containers.quickActions.selectChallenges')}
          </label>
          <div className="flex gap-1">
            <Button variant="ghost" size="sm" onClick={candidates.openDetails} className="!text-xs !h-5 !px-1">
              <span className="inline-flex items-center gap-1">
                <IconArrowsMaximize size={12} />
                {t('admin.contests.containers.quickActions.expand')}
              </span>
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => candidates.setSelectedChallenges(challenges.map((challenge) => challenge.id))}
              className="!text-xs !h-5 !px-1"
            >
              {t('admin.contests.containers.quickActions.selectAll')}
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => candidates.setSelectedChallenges([])}
              className="!text-xs !h-5 !px-1"
            >
              {t('admin.contests.containers.quickActions.clear')}
            </Button>
          </div>
        </div>
        <div className="flex-1 min-h-[260px] overflow-y-auto border border-neutral-300/30 rounded-md bg-black/10">
          <div className="p-2 border-b border-neutral-300/20">
            <div className="relative">
              <IconSearch
                size={12}
                className="absolute left-2 top-1/2 -translate-y-1/2 text-neutral-500 pointer-events-none"
              />
              <input
                type="text"
                value={challengeSearch}
                onChange={(event) => candidates.changeSearch(event.target.value)}
                placeholder={t('admin.contests.containers.quickActions.searchPlaceholder')}
                className="w-full h-7 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md text-xs text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 transition-all duration-200"
              />
            </div>
          </div>
          {challenges.length > 0 ? (
            <div className="grid grid-cols-3 gap-0.5 p-1">
              {challenges.map((challenge) => (
                <div
                  key={challenge.id}
                  className="flex items-center p-1 hover:bg-black/30 rounded transition-colors min-w-0"
                >
                  <input
                    type="checkbox"
                    id={`challenge-${challenge.id}`}
                    checked={selectedChallenges.includes(challenge.id)}
                    onChange={(event) => candidates.selectChallenge(challenge.id, event.target.checked)}
                    className="w-3 h-3 shrink-0 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
                  />
                  <label
                    htmlFor={`challenge-${challenge.id}`}
                    className="ml-1.5 text-xs font-mono text-neutral-300 cursor-pointer truncate min-w-0"
                  >
                    {challenge.name}
                  </label>
                </div>
              ))}
            </div>
          ) : (
            <div className="p-3 text-xs font-mono text-neutral-500">
              {t('admin.contests.containers.quickActions.noChallenges')}
            </div>
          )}
        </div>
        {Math.ceil(challengeTotal / challengePageSize) > 1 && (
          <div className="flex items-center justify-between gap-2 mt-2 px-1">
            <span className="text-[11px] font-mono text-geek-400/80 whitespace-nowrap">
              {t('admin.contests.containers.quickActions.pageHint')}
            </span>
            <div className="flex items-center gap-2 ml-auto">
              <button
                disabled={challengePage === 1}
                onClick={() => candidates.setChallengePage((page) => page - 1)}
                className="inline-flex items-center justify-center w-7 h-7 rounded-md border border-geek-400/40 bg-geek-400/10 text-geek-300 hover:bg-geek-400/20 hover:text-geek-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                aria-label={t('common.previous')}
              >
                <IconChevronLeft size={15} />
              </button>
              <span className="text-xs font-mono text-neutral-300 min-w-[56px] text-center">
                {challengePage} / {Math.ceil(challengeTotal / challengePageSize)}
              </span>
              <button
                disabled={challengePage >= Math.ceil(challengeTotal / challengePageSize)}
                onClick={() => candidates.setChallengePage((page) => page + 1)}
                className="inline-flex items-center justify-center w-7 h-7 rounded-md border border-geek-400/40 bg-geek-400/10 text-geek-300 hover:bg-geek-400/20 hover:text-geek-200 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
                aria-label={t('common.next')}
              >
                <IconChevronRight size={15} />
              </button>
            </div>
          </div>
        )}
      </div>
      <VictimCandidateDialog
        t={t}
        isOpen={candidates.detailsOpen}
        onClose={candidates.closeDetails}
        detailChallenges={candidates.detailChallenges}
        detailChallengeTotal={candidates.detailChallengeTotal}
        detailChallengePage={candidates.detailChallengePage}
        challengePageSize={challengePageSize}
        challengeSearch={challengeSearch}
        selectedChallenges={selectedChallenges}
        typeLabels={{
          static: t('admin.challenge.types.static'),
          dynamic: t('admin.challenge.types.dynamic'),
          pods: t('admin.challenge.types.pods'),
        }}
        onChallengeSearchChange={candidates.changeSearch}
        onChallengeSelectionChange={candidates.selectChallenge}
        onPageChange={candidates.setDetailChallengePage}
        getChallengeCategoryChipClass={getChallengeCategoryChipClass}
        getChallengeTypeChipClass={getChallengeTypeChipClass}
      />
    </>
  );
}
