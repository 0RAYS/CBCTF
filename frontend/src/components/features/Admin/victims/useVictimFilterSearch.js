import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from '../../../../utils/toast';

export default function useVictimFilterSearch(scope) {
  const { t } = useTranslation();
  const [searchResults, setSearchResults] = useState({ users: [], teams: [], challenges: [] });
  const [searchLoading, setSearchLoading] = useState({ users: false, teams: false, challenges: false });
  const usersSearchRef = useRef(null);
  const teamsSearchRef = useRef(null);
  const challengesSearchRef = useRef(null);
  const timers = useRef({});
  const requests = useRef({});

  const clearResults = (key) => {
    requests.current[key] = (requests.current[key] || 0) + 1;
    clearTimeout(timers.current[key]);
    setSearchResults((previous) => ({ ...previous, [key]: [] }));
    setSearchLoading((previous) => ({ ...previous, [key]: false }));
  };

  useEffect(() => {
    const refs = { users: usersSearchRef, teams: teamsSearchRef, challenges: challengesSearchRef };
    const closeOutside = (event) => {
      Object.entries(refs).forEach(([key, ref]) => {
        if (ref.current && !ref.current.contains(event.target)) clearResults(key);
      });
    };
    document.addEventListener('mousedown', closeOutside);
    return () => {
      document.removeEventListener('mousedown', closeOutside);
      Object.values(timers.current).forEach(clearTimeout);
      Object.keys(requests.current).forEach((key) => {
        requests.current[key] += 1;
      });
    };
  }, []);

  const onSearch = (model, name) => {
    const key = { User: 'users', Team: 'teams', Challenge: 'challenges' }[model];
    clearResults(key);
    const keyword = name.trim();
    if (!keyword) return;
    const session = requests.current[key];
    timers.current[key] = setTimeout(async () => {
      setSearchLoading((previous) => ({ ...previous, [key]: true }));
      try {
        const response = await scope.search[key](keyword);
        if (requests.current[key] !== session) return;
        let results = response.code === 200 ? response.data[key] || [] : [];
        if (key === 'challenges' && scope.contestId) {
          results = results
            .filter((challenge) => challenge.name?.toLowerCase().includes(keyword.toLowerCase()))
            .slice(0, 10);
        }
        setSearchResults((previous) => ({ ...previous, [key]: results }));
      } catch (error) {
        if (requests.current[key] === session)
          toast.danger({ description: error.message || t(`${scope.translationKey}.toast.searchFailed`) });
      } finally {
        if (requests.current[key] === session) setSearchLoading((previous) => ({ ...previous, [key]: false }));
      }
    }, 300);
  };

  return {
    searchResults,
    searchLoading,
    usersSearchRef,
    teamsSearchRef,
    challengesSearchRef,
    onSearch,
    clearResults,
    resetSearch: () => ['users', 'teams', 'challenges'].forEach(clearResults),
  };
}
