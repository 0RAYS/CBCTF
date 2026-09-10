import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getNotInContestChallengeList } from '../../../../../api/admin/challenge';
import { addContestChallenge } from '../../../../../api/admin/contest';
import { toast } from '../../../../../utils/toast';
import ChallengePickerDialog from './ChallengePickerDialog';
import { challengeQuery, toggleChallengeSelection } from './challengeData.js';

export default function ChallengePickerSession({ contestId, categories, onClose, onAdded }) {
  const { t } = useTranslation();
  const [filters, setFilters] = useState({ page: 1, type: 'all', category: 'all', name: '', description: '' });
  const [challenges, setChallenges] = useState([]);
  const [count, setCount] = useState(0);
  const [selected, setSelected] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const busy = useRef(false);
  const alive = useRef(true);
  useEffect(() => {
    alive.current = true;
    return () => {
      alive.current = false;
    };
  }, []);
  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    getNotInContestChallengeList(contestId, challengeQuery(filters))
      .then((response) => {
        if (cancelled) return;
        if (response.code !== 200)
          throw new Error(response.msg || t('admin.contests.challenges.toast.fetchPoolFailed'));
        setChallenges(response.data?.challenges || []);
        setCount(response.data?.count || 0);
      })
      .catch((error) => {
        if (!cancelled) {
          setChallenges([]);
          toast.danger({ description: error.message });
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [contestId, filters, t]);
  const changeFilter = (key, value) => {
    setLoading(true);
    setFilters((previous) => ({ ...previous, page: 1, [key]: value }));
  };
  const add = async () => {
    if (busy.current || !selected.length) return;
    busy.current = true;
    setSaving(true);
    try {
      const response = await addContestChallenge(
        contestId,
        selected.map((item) => item.id)
      );
      if (response.code !== 200) throw new Error(response.msg || t('admin.contests.challenges.toast.addFailed'));
      if (!alive.current) return;
      toast.success({ description: t('admin.contests.challenges.toast.addSuccess') });
      onAdded();
      onClose();
    } catch (error) {
      if (alive.current) toast.danger({ description: error.message });
    } finally {
      busy.current = false;
      if (alive.current) setSaving(false);
    }
  };
  return (
    <ChallengePickerDialog
      isOpen
      challenges={challenges}
      selectedChallenges={selected}
      categories={categories}
      totalCount={count}
      currentPage={filters.page}
      pageSize={10}
      loading={loading}
      saving={saving}
      searchQuery={filters.name}
      descQuery={filters.description}
      type={filters.type}
      category={filters.category}
      onClose={() => {
        if (!busy.current) onClose();
      }}
      onConfirm={add}
      onSearch={(value) => changeFilter('name', value)}
      onDescSearch={(value) => changeFilter('description', value)}
      onFilterCategoryChange={(value) => changeFilter('category', value)}
      onFilterTypeChange={(value) => changeFilter('type', value)}
      onPageChange={(value) => changeFilter('page', value)}
      onSelect={(challenge) => setSelected((previous) => toggleChallengeSelection(previous, challenge))}
    />
  );
}
