import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getChallengeCategories, getChallengeList } from '../../../../../api/challenge';
import { toast } from '../../../../../utils/toast';
import { mapChallengeStatusToViewModel, normalizeCategories } from '../models/challengeViewModel';

export default function useChallengeList(contestId) {
  const { t } = useTranslation();
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [unsolvedOnly, setUnsolvedOnly] = useState(false);
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [revision, setRevision] = useState(0);
  const categoriesRequestRef = useRef(0);
  const pageSize = 20;

  const refreshCategories = async () => {
    const request = ++categoriesRequestRef.current;
    try {
      const response = await getChallengeCategories(contestId);
      if (request !== categoriesRequestRef.current) return;
      if (response.code !== 200) throw new Error(response.msg || t('game.challenges.toast.fetchCategoriesFailed'));
      setCategories(normalizeCategories(response.data));
    } catch (error) {
      if (request !== categoriesRequestRef.current) return;
      setCategories([]);
      toast.danger({ description: error.message || t('game.challenges.toast.fetchCategoriesFailed') });
    }
  };

  useEffect(() => {
    refreshCategories();
    return () => {
      categoriesRequestRef.current += 1;
    };
  }, [contestId]);

  useEffect(() => {
    let active = true;
    setIsLoading(true);
    setError(null);
    const fetchList = async () => {
      try {
        const params = { limit: pageSize, offset: (currentPage - 1) * pageSize, unsolved: unsolvedOnly };
        if (selectedCategory) params.category = selectedCategory;
        const response = await getChallengeList(contestId, params);
        if (!active) return;
        if (response.code !== 200) throw new Error(response.msg || t('game.challenges.toast.fetchListFailed'));
        const data = response.data || {};
        setChallenges((data.challenges || []).map((challenge) => mapChallengeStatusToViewModel(challenge)));
        setTotalCount(data.count || 0);
        const lastPage = Math.max(1, Math.ceil((data.count || 0) / pageSize));
        if (currentPage > lastPage) setCurrentPage(lastPage);
      } catch (error) {
        if (!active) return;
        setChallenges([]);
        setTotalCount(0);
        setError(error.message || t('game.challenges.toast.fetchListFailed'));
      } finally {
        if (active) setIsLoading(false);
      }
    };
    fetchList();
    return () => {
      active = false;
    };
  }, [contestId, currentPage, selectedCategory, unsolvedOnly, revision]);

  const refresh = () => setRevision((value) => value + 1);
  const updateChallenge = (challenge) => {
    setChallenges((previous) => previous.map((item) => (item.id === challenge.id ? challenge : item)));
  };

  return {
    refresh,
    refreshCategories,
    updateChallenge,
    boardProps: {
      categories,
      selectedCategory,
      onCategoryChange: (category) => {
        setSelectedCategory((previous) => (previous === category ? '' : category));
        setCurrentPage(1);
      },
      unsolvedOnly,
      onSolvedFilterChange: () => {
        setUnsolvedOnly((previous) => !previous);
        setCurrentPage(1);
      },
      challenges,
      totalCount,
      currentPage,
      pageSize,
      onPageChange: setCurrentPage,
      isLoading,
      error,
      onRetry: refresh,
    },
  };
}
