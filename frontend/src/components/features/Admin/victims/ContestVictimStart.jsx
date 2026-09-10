import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { motion } from 'motion/react';
import { IconClockPlay, IconPlayerPlay, IconUsers } from '@tabler/icons-react';
import { Button } from '../../../common';
import { getContestTeams, startContestVictims } from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import { buildVictimStartPayload, estimateVictimTeams } from './victimPayload';
import useVictimCandidates from './useVictimCandidates';
import VictimCandidatePicker from './VictimCandidatePicker';
import { VictimStartDialog } from './VictimStartDialogs';

export default function ContestVictimStart({ contestId, onStarted }) {
  const { t } = useTranslation();
  const candidates = useVictimCandidates(contestId);
  const [totalTeamCount, setTotalTeamCount] = useState(0);
  const [randomTeamPercentage, setRandomTeamPercentage] = useState(50);
  const [victimDurationInput, setVictimDurationInput] = useState('7200');
  const [startOpen, setStartOpen] = useState(false);
  const payload = buildVictimStartPayload(candidates.selectedChallenges, randomTeamPercentage, victimDurationInput);
  const victimDurationSeconds = payload.duration;
  const isVictimDurationValid = victimDurationSeconds > 0;
  const selectedTeamCount = estimateVictimTeams(totalTeamCount, randomTeamPercentage);

  useEffect(() => {
    let active = true;
    getContestTeams(contestId, { limit: 1, offset: 0 })
      .then((response) => {
        if (active && response.code === 200) setTotalTeamCount(response.data.count || 0);
      })
      .catch((error) => {
        if (active)
          toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchTeamsFailed') });
      });
    return () => {
      active = false;
    };
  }, [contestId, t]);

  const start = async () => {
    if (!isVictimDurationValid) {
      toast.warning({ description: t('admin.contests.containers.toast.invalidDuration') });
      return;
    }
    if (!payload.challenges.length || !selectedTeamCount) {
      toast.warning({ description: t('admin.contests.containers.toast.selectStartRequired') });
      return;
    }
    try {
      const response = await startContestVictims(contestId, payload.challenges, payload.team_ratio, payload.duration);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.containers.toast.taskDispatched') });
        candidates.setSelectedChallenges([]);
        onStarted();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.containers.toast.taskDispatchFailed') });
    }
    setStartOpen(false);
  };

  const formatVictimDuration = (seconds) => {
    if (!seconds || seconds <= 0) return t('admin.contests.containers.quickActions.invalidDuration');
    if (seconds < 60) return t('utils.time.units.second', { count: seconds });
    if (seconds < 3600) {
      const minutes = Math.floor(seconds / 60);
      const remainingSeconds = seconds % 60;
      return remainingSeconds > 0
        ? `${t('utils.time.units.minute', { count: minutes })}${t('utils.time.units.second', { count: remainingSeconds })}`
        : t('utils.time.units.minute', { count: minutes });
    }
    if (seconds < 86400) {
      const hours = Math.floor(seconds / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);
      return minutes > 0
        ? `${t('utils.time.units.hour', { count: hours })}${t('utils.time.units.minute', { count: minutes })}`
        : t('utils.time.units.hour', { count: hours });
    }
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    return hours > 0
      ? `${t('utils.time.units.day', { count: days })}${t('utils.time.units.hour', { count: hours })}`
      : t('utils.time.units.day', { count: days });
  };

  return (
    <>
      <style>{`
        .slider::-webkit-slider-thumb {
          appearance: none; height: 12px; width: 12px; border-radius: 50%; background: #597ef7;
          cursor: pointer; border: 2px solid #1f2937; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
        }
        .slider::-moz-range-thumb {
          height: 12px; width: 12px; border-radius: 50%; background: #597ef7;
          cursor: pointer; border: 2px solid #1f2937; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
        }
        .slider::-webkit-slider-track { background: transparent; }
        .slider::-moz-range-track { background: transparent; }
      `}</style>
      <motion.div
        className="grid grid-cols-1 gap-6 items-stretch"
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4 min-h-[460px] flex flex-col">
          <div className="flex items-center gap-2 mb-3">
            <IconPlayerPlay size={18} className="text-neutral-400" />
            <h3 className="text-base font-mono text-neutral-50">{t('admin.contests.containers.quickActions.title')}</h3>
          </div>
          <div className="flex flex-col gap-4 flex-1 min-h-0">
            <div className="border border-neutral-300/20 rounded-md bg-black/10 p-4">
              <div className="flex justify-between items-center mb-2">
                <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
                  <IconUsers size={14} />
                  <span className="text-xs font-mono text-neutral-400">
                    {t('admin.contests.containers.quickActions.randomTeams')}
                  </span>
                  <span className="text-xs font-mono text-geek-400">{randomTeamPercentage}%</span>
                </label>
                <span className="text-xs font-mono text-neutral-500">
                  {t('admin.contests.containers.quickActions.estimatedTeams', { count: selectedTeamCount })}
                </span>
              </div>
              <div className="mb-3 p-2 border border-neutral-300/20 rounded-md bg-black/10">
                <div className="relative">
                  <input
                    type="range"
                    min="0"
                    max="100"
                    value={randomTeamPercentage}
                    onChange={(event) => setRandomTeamPercentage(Number.parseInt(event.target.value, 10))}
                    className="w-full h-1 bg-neutral-700 rounded-lg appearance-none cursor-pointer slider"
                    style={{
                      background: `linear-gradient(to right, #597ef7 0%, #597ef7 ${randomTeamPercentage}%, #374151 ${randomTeamPercentage}%, #374151 100%)`,
                    }}
                  />
                </div>
              </div>
              <div className="border border-neutral-300/30 rounded-md bg-black/10 p-3">
                <p className="text-xs font-mono text-neutral-400">
                  {t('admin.contests.containers.quickActions.teamSelectionHint', { total: totalTeamCount })}
                </p>
              </div>
            </div>
            <div className="border border-neutral-300/20 rounded-md bg-black/10 p-4">
              <div className="flex justify-between items-center gap-3 mb-2">
                <label className="text-xs font-mono text-neutral-400 flex items-center gap-1">
                  <IconClockPlay size={14} />
                  <span>{t('common.duration')}</span>
                </label>
                <span className="text-xs font-mono text-geek-400">
                  {t('admin.contests.containers.quickActions.durationPreview', {
                    value: formatVictimDuration(victimDurationSeconds),
                  })}
                </span>
              </div>
              <div className="flex items-center gap-3">
                <input
                  type="number"
                  min="1"
                  step="1"
                  value={victimDurationInput}
                  onChange={(event) => setVictimDurationInput(event.target.value)}
                  className="w-full h-9 px-3 bg-black/20 border border-neutral-300/30 rounded-md text-sm text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 transition-all duration-200"
                />
                <span className="text-xs font-mono text-neutral-400 whitespace-nowrap">
                  {t('admin.contests.containers.quickActions.durationUnit')}
                </span>
              </div>
              <p className="mt-2 text-xs font-mono text-neutral-400">
                {t('admin.contests.containers.quickActions.durationHint')}
              </p>
            </div>
            <VictimCandidatePicker candidates={candidates} />
          </div>
          <div className="mt-4 flex justify-end">
            <Button
              variant="primary"
              size="sm"
              align="icon-left"
              icon={<IconPlayerPlay size={14} />}
              onClick={() => setStartOpen(true)}
              disabled={!payload.challenges.length || !selectedTeamCount || !isVictimDurationValid}
              className="!text-xs !h-7 !px-3"
            >
              {t('admin.contests.containers.quickActions.startButton', {
                challenges: payload.challenges.length,
                teams: selectedTeamCount,
              })}
            </Button>
          </div>
        </div>
      </motion.div>
      <VictimStartDialog
        t={t}
        isOpen={startOpen}
        onClose={() => setStartOpen(false)}
        onConfirm={start}
        selectedChallenges={candidates.selectedChallenges}
        challenges={candidates.selectedChallengeDetails}
        randomTeamPercentage={randomTeamPercentage}
        selectedTeamCount={selectedTeamCount}
        totalTeamCount={totalTeamCount}
        victimDurationSeconds={victimDurationSeconds}
        formatVictimDuration={formatVictimDuration}
      />
    </>
  );
}
