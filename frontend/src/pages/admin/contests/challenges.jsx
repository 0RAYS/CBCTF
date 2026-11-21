import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { toast } from '@/utils/toast';
import { getContestChallenges, getContestInfo } from '@/api/admin/contest';
import {
  createChallenge,
  deleteChallenge,
  getChallenge,
  getChallengeCategories,
  updateChallenge,
} from '@/api/admin/challenge';
import { DEFAULT_CHALLENGE_CATEGORIES, mergeChallengeCategories } from '@/config/challenges';
import ChallengeTable from '@/components/features/Admin/Contests/challenges/ChallengeTable';
import ChallengePickerSession from '@/components/features/Admin/Contests/challenges/ChallengePickerSession';
import { challengeQuery } from '@/components/features/Admin/Contests/challenges/challengeData.js';
import ChallengeEditorDialog from '@/components/features/Admin/challenges/editor/ChallengeEditorDialog.jsx';
import { applyChallengeType, challengeDraft } from '@/components/features/Admin/challenges/editor/challengeData.js';
import { buildChallengePayload } from '@/components/features/Admin/challenges/editor/challengePayload.js';
import { useDebounce } from '@/hooks/useDebounce';

function ChallengeManagement({ contestId }) {
  const { t } = useTranslation();
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [categories, setCategories] = useState(DEFAULT_CHALLENGE_CATEGORIES);
  const [currentContest, setCurrentContest] = useState(null);
  const [filters, setFilters] = useState({ page: 1, type: 'all', category: 'all', name: '' });
  const [nameQuery, setNameQuery] = useState('');
  const [picker, setPicker] = useState(false);
  const [editing, setEditing] = useState(null);
  const [editorMode, setEditorMode] = useState('add');
  const [initDockerCompose, setInitDockerCompose] = useState([]);
  const [flagVersion, setFlagVersion] = useState(0);
  const [revision, setRevision] = useState(0);
  const [loadedKey, setLoadedKey] = useState('');
  const refresh = () => setRevision((previous) => previous + 1);
  // 避免每次输入都发起请求
  const debouncedName = useDebounce(nameQuery, 300);
  // 刷新不会清空表格
  const filterKey = `${filters.page}|${filters.type}|${filters.category}|${filters.name}`;
  const loading = loadedKey !== filterKey;

  useEffect(() => {
    setFilters((previous) =>
      previous.name === debouncedName ? previous : { ...previous, page: 1, name: debouncedName }
    );
  }, [debouncedName]);

  useEffect(() => {
    let cancelled = false;
    getContestInfo(contestId)
      .then((response) => {
        if (!cancelled && response.code === 200) setCurrentContest(response.data);
      })
      .catch(() => {
        // 比赛信息加载失败时归属展示回退为比赛 ID
      });
    return () => {
      cancelled = true;
    };
  }, [contestId]);

  useEffect(() => {
    let cancelled = false;
    getChallengeCategories({ contest_id: contestId })
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
      })
      .finally(() => {
        if (!cancelled) setLoadedKey(filterKey);
      });
    return () => {
      cancelled = true;
    };
  }, [contestId, filters, revision, t]);
  const changeFilter = (field, value) => setFilters((previous) => ({ ...previous, page: 1, [field]: value }));

  const createInContest = () => {
    setEditorMode('add');
    setEditing(
      challengeDraft({
        contest_id: contestId,
        contest_name: currentContest?.name || '',
        hidden: true,
        flags: [{ id: 0, value: '' }],
      })
    );
    setInitDockerCompose([]);
  };

  const edit = async (challenge) => {
    try {
      const response = await getChallenge(challenge.id);
      if (response.code !== 200) throw new Error(response.msg || t('admin.contests.challenges.toast.fetchListFailed'));
      const detail = challengeDraft(response.data);
      setEditorMode('edit');
      setEditing(detail);
      setInitDockerCompose([{ id: detail.id, value: detail.docker_compose }]);
    } catch (error) {
      toast.danger({ description: error.message });
    }
  };

  const remove = async (challenge) => {
    try {
      const response = await deleteChallenge(challenge.id);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.challenges.toast.removeSuccess') });
        if (challenges.length === 1 && filters.page > 1) changeFilter('page', filters.page - 1);
        else refresh();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.removeFailed') });
    }
  };

  const save = async (challenge) => {
    try {
      const payload = buildChallengePayload(challenge, editorMode, initDockerCompose);
      const response =
        editorMode === 'add' ? await createChallenge(payload) : await updateChallenge(challenge.id, payload);
      if (response.code !== 200) throw new Error(response.msg || t('admin.contests.challenges.toast.updateFailed'));
      toast.success({ description: t('admin.contests.challenges.toast.updateSuccess') });
      setEditing(null);
      setFlagVersion((previous) => previous + 1);
      refresh();
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.updateFailed') });
    }
  };

  return (
    <>
      <ChallengeTable
        challenges={challenges}
        totalCount={totalCount}
        currentPage={filters.page}
        pageSize={10}
        categories={categories}
        filterCategory={filters.category}
        filterType={filters.type}
        nameQuery={nameQuery}
        loading={loading}
        onPageChange={(page) => changeFilter('page', page)}
        onAddChallenge={createInContest}
        onCopyChallenge={() => setPicker(true)}
        onEditChallenge={edit}
        onDeleteChallenge={remove}
        onFilterTypeChange={(value) => changeFilter('type', value)}
        onFilterCategoryChange={(value) => changeFilter('category', value)}
        onNameChange={setNameQuery}
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
          key={editing.id || 'new'}
          isOpen
          mode={editorMode}
          challenge={editing}
          categories={categories}
          contestLocked
          flagVersion={flagVersion}
          onClose={() => setEditing(null)}
          onSubmit={save}
          onChange={(updated) => setEditing((previous) => applyChallengeType(previous, updated))}
          onAddFlag={() =>
            setEditing((previous) => ({ ...previous, flags: [...previous.flags, { id: 0, value: '' }] }))
          }
          onRemoveFlag={(index) =>
            setEditing((previous) => ({ ...previous, flags: previous.flags.filter((_, i) => i !== index) }))
          }
          onFlagChange={(index, value) =>
            setEditing((previous) => ({
              ...previous,
              flags: previous.flags.map((flag, i) => (i === index ? { ...flag, value } : flag)),
            }))
          }
        />
      )}
    </>
  );
}

export default function AdminChallengePage() {
  const { id } = useParams();
  return <ChallengeManagement key={id} contestId={Number(id)} />;
}
