import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getNotInContestChallengeList } from '../../../../../api/admin/challenge';
import { addContestChallenge } from '../../../../../api/admin/contest';
import { toast } from '../../../../../utils/toast';
import ChallengePickerDialog from './ChallengePickerDialog';
import { challengeQuery, toggleChallengeSelection } from './challengeData.js';
import useBatchAction from '../../batch/useBatchAction.js';
import { remainingBatchIds } from '../../batch/batchModel.js';

export default function ChallengePickerSession({ contestId, categories, onClose, onAdded }) {
  const { t } = useTranslation();
  const [filters, setFilters] = useState({ page: 1, type: 'all', category: 'all', name: '', description: '' });
  const [challenges, setChallenges] = useState([]);
  const [count, setCount] = useState(0);
  const [selected, setSelected] = useState([]);
  const [loading, setLoading] = useState(true);
  const action = useBatchAction(contestId);
  const [revision, setRevision] = useState(0);
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
  }, [contestId, filters, revision, t]);
  const changeFilter = (key, value) => {
    setLoading(true);
    setFilters((previous) => ({ ...previous, page: 1, [key]: value }));
  };
  const add = async () => {
    if (!selected.length) return;
    const ids = selected.map((item) => item.id);
    await action.run(() => addContestChallenge(contestId, ids), {
      nested: true,
      successMessage: t('admin.contests.challenges.toast.addSuccess'),
      failureMessage: t('admin.contests.challenges.toast.addFailed'),
      onResult: (batch) => {
        const remaining = batch ? remainingBatchIds(ids, batch) : [];
        setSelected((current) => current.filter((item) => !ids.includes(item.id) || remaining.includes(item.id)));
        setRevision((value) => value + 1);
        onAdded();
        if (!batch || batch.status === 'success') onClose();
      },
    });
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
      saving={action.pending}
      batchResult={action.result}
      error={action.error}
      searchQuery={filters.name}
      descQuery={filters.description}
      type={filters.type}
      category={filters.category}
      onClose={() => {
        if (!action.pending) onClose();
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
