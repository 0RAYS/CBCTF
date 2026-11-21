import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconSearch } from '@tabler/icons-react';
import { Input, Spinner } from '@/components/common';
import { getContestList } from '@/api/admin/contest';
import { useDebounce } from '@/hooks';

/** 比赛搜索选择器：按名称查询比赛，从结果中选择一项。 */
export default function ContestSearchSelect({
  value = 'all',
  onChange,
  allowAll = false,
  allowUnassigned = false,
  selectedName = '',
}) {
  const { t } = useTranslation();
  const [query, setQuery] = useState('');
  const [options, setOptions] = useState([]);
  const [picked, setPicked] = useState(null);
  const [loading, setLoading] = useState(false);
  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    getContestList({ limit: 20, offset: 0, name: debouncedQuery })
      .then((response) => {
        if (!cancelled && response.code === 200) setOptions(response.data?.contests || []);
      })
      .catch(() => {
        if (!cancelled) setOptions([]);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [debouncedQuery]);

  const allLabel = t('admin.challenge.contestSearch.all');
  const unassignedLabel = t('admin.challengeModal.contest.unassigned');
  const pick = (option) => {
    setPicked(option);
    onChange(option.id);
  };
  const selectedOption = options.find((contest) => String(contest.id) === value);
  const selectedLabel =
    value === 'all'
      ? allLabel
      : value === '0'
        ? unassignedLabel
        : (selectedOption?.name ?? (picked?.id === value ? picked.name : '')) || selectedName || `#${value}`;
  const rowClass = (optionValue) =>
    `flex w-full items-center justify-between gap-2 px-3 py-2 text-left font-mono text-sm transition-colors ${
      value === optionValue ? 'bg-geek-400/15 text-geek-300' : 'text-neutral-300 hover:bg-neutral-800/60'
    }`;

  return (
    <div className="space-y-2">
      <Input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder={t('admin.challenge.contestSearch.placeholder')}
        icon={<IconSearch size={16} />}
        iconRight={loading && <Spinner size="sm" />}
      />
      {value !== 'all' && (
        <div className="flex items-center justify-between gap-2 rounded-md border border-geek-600/40 bg-geek-400/10 px-3 py-2">
          <span className="truncate font-mono text-sm text-geek-300">{selectedLabel}</span>
          <span className="shrink-0 font-mono text-xs text-neutral-500">
            {t('admin.challenge.contestSearch.selected')}
          </span>
        </div>
      )}
      <div className="max-h-56 overflow-y-auto rounded-md border border-neutral-700/70">
        {allowAll && (
          <button type="button" className={rowClass('all')} onClick={() => pick({ id: 'all', name: allLabel })}>
            {allLabel}
          </button>
        )}
        {allowUnassigned && (
          <button type="button" className={rowClass('0')} onClick={() => pick({ id: '0', name: unassignedLabel })}>
            {unassignedLabel}
          </button>
        )}
        {options.map((contest) => (
          <button
            key={contest.id}
            type="button"
            className={rowClass(String(contest.id))}
            onClick={() => pick({ id: String(contest.id), name: contest.name })}
          >
            <span className="truncate">{contest.name}</span>
            <span className="shrink-0 text-xs text-neutral-500">#{contest.id}</span>
          </button>
        ))}
        {options.length === 0 && !loading && (
          <p className="px-3 py-3 font-mono text-sm text-neutral-500">{t('admin.challenge.contestSearch.empty')}</p>
        )}
      </div>
    </div>
  );
}
