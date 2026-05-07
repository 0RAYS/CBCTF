import { IconFilter, IconSearch, IconTarget, IconUsers } from '@tabler/icons-react';
import { motion } from 'motion/react';
import { Button } from '../../../../components/common';

export function ContainerFilters({
  t,
  filters,
  searchResults,
  searchLoading,
  usersSearchRef,
  teamsSearchRef,
  challengesSearchRef,
  onResetFilters,
  onFilterChange,
  onSearch,
  onSetSearchResults,
  onSetSearchLoading,
}) {
  const runSearch = (model, value, key) => {
    onSearch(
      model,
      value,
      (results) => onSetSearchResults((prev) => ({ ...prev, [key]: results })),
      (loading) => onSetSearchLoading((prev) => ({ ...prev, [key]: loading }))
    );
  };

  return (
    <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
      <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <IconFilter size={18} className="text-neutral-400" />
            <h3 className="text-base font-mono text-neutral-50">{t('admin.contests.containers.filters.title')}</h3>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={onResetFilters}
            className="!text-neutral-400 hover:!text-neutral-300 !text-xs !h-6 !px-2"
          >
            {t('admin.contests.containers.filters.reset')}
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          <SearchFilter
            t={t}
            label={t('admin.contests.containers.filters.userName')}
            placeholder={t('admin.contests.containers.filters.searchUserPlaceholder')}
            icon={
              <IconSearch size={14} className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400" />
            }
            results={searchResults.users}
            loading={searchLoading.users}
            containerRef={usersSearchRef}
            onSearch={(value) => runSearch('User', value, 'users')}
            onSelect={(user) => {
              onFilterChange('user_id', user.id.toString());
              onSetSearchResults((prev) => ({ ...prev, users: [] }));
            }}
            getLabel={(user) =>
              user.name || user.username || t('admin.contests.containers.filters.userFallback', { id: user.id })
            }
          />

          <SearchFilter
            t={t}
            label={t('admin.contests.containers.filters.teamName')}
            placeholder={t('admin.contests.containers.filters.searchTeamPlaceholder')}
            icon={
              <IconUsers size={14} className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400" />
            }
            results={searchResults.teams}
            loading={searchLoading.teams}
            containerRef={teamsSearchRef}
            onSearch={(value) => runSearch('Team', value, 'teams')}
            onSelect={(team) => {
              onFilterChange('team_id', team.id.toString());
              onSetSearchResults((prev) => ({ ...prev, teams: [] }));
            }}
            getLabel={(team) => team.name || t('admin.contests.containers.filters.teamFallback', { id: team.id })}
          />

          <SearchFilter
            t={t}
            label={t('admin.contests.containers.filters.challengeName')}
            placeholder={t('admin.contests.containers.filters.searchChallengePlaceholder')}
            icon={
              <IconTarget size={14} className="absolute left-2 top-1/2 transform -translate-y-1/2 text-neutral-400" />
            }
            results={searchResults.challenges}
            loading={searchLoading.challenges}
            containerRef={challengesSearchRef}
            onSearch={(value) => runSearch('Challenge', value, 'challenges')}
            onSelect={(challenge) => {
              onFilterChange('challenge_id', challenge.id.toString());
              onSetSearchResults((prev) => ({ ...prev, challenges: [] }));
            }}
            getLabel={(challenge) =>
              challenge.name || t('admin.contests.containers.filters.challengeFallback', { id: challenge.id })
            }
          />
        </div>

        {(filters.user_id || filters.team_id || filters.challenge_id) && (
          <div className="mt-3 pt-3 border-t border-neutral-300/20">
            <div className="flex flex-wrap gap-2">
              {filters.user_id && (
                <FilterChip
                  label={t('admin.contests.containers.filters.userIdLabel')}
                  value={filters.user_id}
                  onClear={() => onFilterChange('user_id', '')}
                />
              )}
              {filters.team_id && (
                <FilterChip
                  label={t('admin.contests.containers.filters.teamIdLabel')}
                  value={filters.team_id}
                  onClear={() => onFilterChange('team_id', '')}
                />
              )}
              {filters.challenge_id && (
                <FilterChip
                  label={t('admin.contests.containers.filters.challengeIdLabel')}
                  value={filters.challenge_id}
                  className="bg-green-400/20 text-green-400 border-green-400/30"
                  onClear={() => onFilterChange('challenge_id', '')}
                />
              )}
            </div>
          </div>
        )}
      </div>
    </motion.div>
  );
}

function SearchFilter({ label, placeholder, icon, results, loading, containerRef, onSearch, onSelect, getLabel }) {
  return (
    <div className="relative" ref={containerRef}>
      <label className="block text-xs font-mono text-neutral-400 mb-1">{label}</label>
      <div className="relative">
        {icon}
        <input
          type="text"
          placeholder={placeholder}
          onChange={(e) => onSearch(e.target.value)}
          className="w-full h-8 pl-7 pr-2 bg-black/20 border border-neutral-300/30 rounded-md text-xs text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 focus:shadow-focus transition-all duration-200"
        />
        {loading && (
          <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
            <div className="w-3 h-3 border border-geek-400 border-t-transparent rounded-full animate-spin" />
          </div>
        )}
      </div>
      {results.length > 0 && (
        <div className="dropdown-custom max-h-32">
          {results.map((item) => (
            <div key={item.id} className="dropdown-option text-xs" onClick={() => onSelect(item)}>
              {getLabel(item)}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function FilterChip({ label, value, onClear, className = 'bg-geek-400/20 text-geek-400 border-geek-400/30' }) {
  return (
    <span className={`px-2 py-1 text-xs font-mono rounded border ${className}`}>
      {label}: {value}
      <button type="button" onClick={onClear} className="ml-1 hover:text-red-400" aria-label={`${label}: ${value}`}>
        x
      </button>
    </span>
  );
}
