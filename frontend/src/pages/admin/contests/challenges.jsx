import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { toast } from '../../../utils/toast';
import { getContestChallenges, removeContestChallenge } from '../../../api/admin/contest';
import { getContestChallengeCategories } from '../../../api/admin/challenge';
import { DEFAULT_CHALLENGE_CATEGORIES, mergeChallengeCategories } from '../../../config/challenges';
import ContestChallenges from '../../../components/features/Admin/Contests/challenges/ContestChallenges';
import ChallengePickerSession from '../../../components/features/Admin/Contests/challenges/ChallengePickerSession';
import ChallengeEditorDialog from '../../../components/features/Admin/Contests/challenges/ChallengeEditorDialog';
import { challengeQuery } from '../../../components/features/Admin/Contests/challenges/challengeData.js';

function ContestChallengeManagement({ contestId }) {
  const { t } = useTranslation();
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [categories, setCategories] = useState(DEFAULT_CHALLENGE_CATEGORIES);
  const [filters, setFilters] = useState({ page: 1, type: 'all', category: 'all', name: '' });
  const [picker, setPicker] = useState(false);
  const [editing, setEditing] = useState(null);
  const [revision, setRevision] = useState(0);
  const refresh = () => setRevision((previous) => previous + 1);
  useEffect(() => {
    let cancelled = false;
    getContestChallengeCategories(contestId)
      .then((response) => {
        if (!cancelled && response.code === 200) setCategories(mergeChallengeCategories(response.data));
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message });
      });
    return () => {
      cancelled = true;
    };
  }, [contestId]);
  useEffect(() => {
    let cancelled = false;
    getContestChallenges(contestId, challengeQuery(filters))
      .then((response) => {
        if (!cancelled && response.code === 200) {
          setChallenges(response.data?.challenges || []);
          setTotalCount(response.data?.count || 0);
        }
      })
      .catch((error) => {
        if (!cancelled)
          toast.danger({ description: error.message || t('admin.contests.challenges.toast.fetchListFailed') });
      });
    return () => {
      cancelled = true;
    };
  }, [contestId, filters, revision, t]);
  const changeFilter = (field, value) => setFilters((previous) => ({ ...previous, page: 1, [field]: value }));
  const remove = async (challenge) => {
    try {
      const response = await removeContestChallenge(contestId, challenge.id);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.challenges.toast.removeSuccess') });
        if (challenges.length === 1 && filters.page > 1) changeFilter('page', filters.page - 1);
        else refresh();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.removeFailed') });
    }
  };
  return (
    <>
      <ContestChallenges
        challenges={challenges}
        totalCount={totalCount}
        currentPage={filters.page}
        pageSize={10}
        categories={categories}
        filterCategory={filters.category}
        filterType={filters.type}
        nameQuery={filters.name}
        onPageChange={(page) => changeFilter('page', page)}
        onAddChallenge={() => setPicker(true)}
        onEditChallenge={setEditing}
        onDeleteChallenge={remove}
        onFilterTypeChange={(value) => changeFilter('type', value)}
        onFilterCategoryChange={(value) => changeFilter('category', value)}
        onNameChange={(value) => changeFilter('name', value)}
      />
      {picker && (
        <ChallengePickerSession
          contestId={contestId}
          categories={categories}
          onClose={() => setPicker(false)}
          onAdded={refresh}
        />
      )}
      {editing && (
        <ChallengeEditorDialog
          key={editing.id}
          contestId={contestId}
          challenge={editing}
          onClose={() => setEditing(null)}
          onSaved={refresh}
        />
      )}
    </>
  );
}

export default function AdminContestChallengesPage() {
  const { id } = useParams();
  return <ContestChallengeManagement key={id} contestId={Number(id)} />;
}
