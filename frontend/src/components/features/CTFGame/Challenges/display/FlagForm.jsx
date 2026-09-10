import { useTranslation } from 'react-i18next';
import { Button } from '../../../../common';

export default function FlagForm({
  challenge,
  prefix,
  flag,
  onChange,
  onSubmit,
  disabled,
  submitting,
  flagError,
  error,
}) {
  const { t } = useTranslation();
  return (
    <div className="min-w-0 w-full space-y-1.5">
      <label htmlFor="challenge-flag" className="text-neutral-400 font-mono text-sm">
        {t('game.challengeModal.sections.submitFlag')}
      </label>
      <form onSubmit={onSubmit} className="flex flex-col gap-2 sm:flex-row" aria-busy={submitting}>
        <input
          id="challenge-flag"
          type="text"
          value={flag}
          autoComplete="off"
          autoCapitalize="none"
          spellCheck={false}
          readOnly={disabled}
          aria-invalid={!!flagError}
          aria-describedby={flagError ? 'challenge-flag-error' : undefined}
          onChange={(event) => onChange(event.target.value)}
          placeholder={`${prefix}{...}`}
          className="min-w-0 w-full sm:flex-1 h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4 text-neutral-50 placeholder-neutral-400 focus:border-geek-400 focus:shadow-focus transition-colors duration-200"
        />
        <Button
          type="submit"
          variant="primary"
          size="action"
          className="!h-auto min-h-10 whitespace-normal break-words sm:shrink-0"
          loading={submitting}
          disabled={disabled || !flag.trim()}
        >
          {submitting
            ? t('game.challengeModal.submit.submitting')
            : t('game.challengeModal.submit.button', {
                status: challenge.isSolved ? t('common.solved') : t('common.submit'),
                attempts: challenge.attempts,
                max: challenge.maxAttempts || '\u221e',
              })}
        </Button>
      </form>
      {(flagError || error) && (
        <div className="max-h-20 space-y-1 overflow-y-auto text-red-400 text-sm [overflow-wrap:anywhere]">
          {flagError && (
            <div id="challenge-flag-error" role="alert">
              {flagError}
            </div>
          )}
          {error && <div role="alert">{error}</div>}
        </div>
      )}
    </div>
  );
}
