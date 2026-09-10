import { IconClockPlay } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';

/** Controlled UI only: value and onChange(value: number) use seconds; 0 means off. No polling side effects. */
export default function AutoRefreshControl({ value, onChange }) {
  const { t } = useTranslation();

  return (
    <label className="flex items-center gap-1 px-2 h-8 rounded-md border border-neutral-700 bg-neutral-900">
      <IconClockPlay size={13} className="text-neutral-400 shrink-0" aria-hidden="true" />
      <span className="text-xs text-neutral-400 shrink-0">{t('common.autoRefresh')}</span>
      <select
        value={value}
        onChange={(event) => onChange(Number(event.target.value))}
        className="bg-transparent text-xs text-neutral-300 cursor-pointer focus-visible:outline-2 focus-visible:outline-geek-400"
      >
        {[5, 10, 30, 60].map((seconds) => (
          <option key={seconds} value={seconds} className="bg-neutral-900">
            {seconds}s
          </option>
        ))}
        <option value={0} className="bg-neutral-900">
          {t('common.autoRefreshOff')}
        </option>
      </select>
    </label>
  );
}
