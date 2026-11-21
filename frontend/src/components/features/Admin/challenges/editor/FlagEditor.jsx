import { useTranslation } from 'react-i18next';
import { IconUsers } from '@tabler/icons-react';
import { Button, Chip, Input, Select } from '@/components/common';
import ScoreCurveChart from '@/components/features/Admin/Contests/ScoreCurveChart';

export default function FlagEditor({ flag, index, onChange, onSave, onViewSolvers, disabled, dirty }) {
  const { t } = useTranslation();
  const label = (key) => t(`admin.challengeModal.labels.${key}`);
  return (
    <fieldset disabled={disabled} className="p-3 border border-neutral-700 rounded-md bg-black/20 space-y-4 min-w-0">
      <legend className="font-mono text-geek-400">Flag #{index + 1}</legend>
      <div className="flex flex-wrap items-center justify-end gap-2">
        <div className="flex flex-wrap items-center gap-2">
          {flag.generation > 1 && (
            <Chip
              size="sm"
              label={t('admin.challengeModal.labels.generation', { generation: flag.generation })}
              colorClass="border-geek-600/30 text-geek-300"
            />
          )}
          {flag.active === false && (
            <Chip
              size="sm"
              label={t('admin.challengeModal.labels.archived')}
              colorClass="bg-yellow-400/20 text-yellow-300"
            />
          )}
        </div>
        <div className="flex flex-wrap gap-2">
          {flag.id && (
            <Button
              size="sm"
              variant="ghost"
              disabled={disabled}
              icon={<IconUsers size={15} />}
              onClick={onViewSolvers}
            >
              {t('admin.challengeModal.actions.viewSolvers')}
            </Button>
          )}
          <Button size="sm" variant="primary" disabled={disabled || !dirty} onClick={onSave}>
            {t('common.saveChanges')}
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <Select
          label={label('scoreCurve')}
          value={flag.score_type}
          onChange={(e) => onChange({ ...flag, score_type: Number(e.target.value) })}
          options={['static', 'linear', 'log'].map((type, value) => ({
            value,
            label: t(`admin.challengeModal.scoreCurve.${type}`),
          }))}
        />
        {[
          ['score', 'initialScore', false],
          ['current_score', 'currentScore', true],
          ...(flag.score_type !== 0
            ? [
                ['decay', 'decay', false],
                ['min_score', 'minScore', false],
              ]
            : []),
          ['solvers', 'solvers', true],
        ].map(([field, name, readOnly]) => (
          <Input
            key={field}
            label={label(name)}
            type="number"
            step="0.1"
            value={flag[field] ?? 0}
            disabled={readOnly}
            onChange={(e) => onChange({ ...flag, [field]: parseFloat(e.target.value) || 0 })}
          />
        ))}
      </div>
      <div className={disabled ? 'pointer-events-none' : ''}>
        <ScoreCurveChart
          scoreType={flag.score_type}
          score={flag.score}
          decay={flag.decay}
          minScore={flag.min_score}
          onChange={disabled ? undefined : (patch) => onChange({ ...flag, ...patch })}
        />
      </div>
    </fieldset>
  );
}
