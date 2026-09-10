import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Card, Chip, EmptyState, Input, List, Spinner, StatusTag } from '../../../../common';
import { filterTeamFlags, flagFilterOptions } from './teamDetailData.js';

export default function FlagsPanel({ flags = [], loading }) {
  const { t } = useTranslation();
  const [filters, setFilters] = useState({ name: '', type: '', category: '', solved: '' });
  const options = flagFilterOptions(flags);
  const filtered = filterTeamFlags(flags, filters);
  const label = (key) => t(`admin.contests.teams.detail.flags.${key}`);
  const columns = [
    ['value', 'value'],
    ['template', 'template'],
    ['current_score', 'currentScore'],
    ['init_score', 'initScore'],
    ['min_score', 'minScore'],
    ['decay', 'decay'],
    ['solvers', 'solvers'],
    ['solved', 'solved'],
  ].map(([key, name]) => ({ key, label: label(name) }));
  const renderCell = (flag, { key }) =>
    key === 'solved' ? (
      <StatusTag type={flag.solved ? 'success' : 'default'} text={label(flag.solved ? 'solvedYes' : 'solvedNo')} />
    ) : (
      <span className="block max-w-xs truncate" title={String(flag[key] ?? '-')}>
        {flag[key] ?? '-'}
      </span>
    );
  return (
    <section>
      <h2 className="text-xl font-mono text-neutral-50 mb-4">{t('admin.contests.teamDetail.sections.flags')}</h2>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-3 mb-4">
        {['name', 'type', 'category', 'solved'].map((key) => {
          const suffix = key[0].toUpperCase() + key.slice(1);
          const change = (e) => setFilters((previous) => ({ ...previous, [key]: e.target.value }));
          return (
            <label key={key} className="block text-xs font-mono text-neutral-400">
              <span className="block mb-1">{label(`filter${suffix}`)}</span>
              {key === 'name' ? (
                <Input value={filters.name} onChange={change} placeholder={label('filterNamePlaceholder')} />
              ) : (
                <select value={filters[key]} onChange={change} className="select-custom select-custom-sm w-full">
                  <option value="">{label(`filter${suffix}Placeholder`)}</option>
                  {(key === 'solved' ? ['true', 'false'] : options[key === 'type' ? 'types' : 'categories']).map(
                    (value) => (
                      <option key={value} value={value}>
                        {key === 'solved' ? label(value === 'true' ? 'filterSolvedYes' : 'filterSolvedNo') : value}
                      </option>
                    )
                  )}
                </select>
              )}
            </label>
          );
        })}
      </div>
      {loading ? (
        <Card className="flex justify-center p-8">
          <Spinner />
        </Card>
      ) : !filtered.length ? (
        <EmptyState title={t('admin.contests.teamDetail.empty.flags')} />
      ) : (
        <div className="space-y-4">
          {filtered.map((challenge) => (
            <Card key={challenge.id} padding="none" className="overflow-hidden">
              <div className="px-4 py-3 bg-neutral-800/50 flex gap-2 items-center flex-wrap">
                <span className="text-sm font-mono text-neutral-50">{challenge.name}</span>
                {challenge.category && <Chip label={challenge.category} />}
                {challenge.type && <Chip label={challenge.type} />}
                {challenge.hidden && <StatusTag type="warning" text={label('hidden')} />}
              </div>
              <List
                className="[&_table]:min-w-[960px]"
                columns={columns}
                data={challenge.flags || []}
                renderCell={renderCell}
              />
            </Card>
          ))}
        </div>
      )}
    </section>
  );
}
