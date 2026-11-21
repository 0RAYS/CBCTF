import { useEffect, useRef, useState } from 'react';
import { toast } from '@/utils/toast';
import { downloadBlobResponse } from '@/utils/fileDownload';
import AdminChallenge from '@/components/features/Admin/AdminChallenge.jsx';
import ChallengeEditorDialog from '@/components/features/Admin/challenges/editor/ChallengeEditorDialog.jsx';
import ChallengeCopyDialog from '@/components/features/Admin/challenges/editor/ChallengeCopyDialog.jsx';
import { buildChallengePayload } from '@/components/features/Admin/challenges/editor/challengePayload.js';
import { applyChallengeType, challengeDraft } from '@/components/features/Admin/challenges/editor/challengeData.js';
import AdminChallengeTestModal from '@/components/features/Admin/challenges/testing/ChallengeTestDialog.jsx';
import {
  copyChallenge,
  createChallenge,
  deleteChallenge,
  downloadChallengeFile,
  getChallenge,
  getChallengeCategories,
  getChallengeList,
  updateChallenge,
  uploadChallengeFile,
} from '@/api/admin/challenge';
import { getContestList } from '@/api/admin/contest';
import { useDebounce } from '@/hooks';
import { useTranslation } from 'react-i18next';
import { DEFAULT_CHALLENGE_CATEGORIES, mergeChallengeCategories } from '@/config/challenges';

function ChallengesManagement() {
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [selectedType, setSelectedType] = useState('all');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [categories, setCategories] = useState([]);
  const [contests, setContests] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isTestModalOpen, setIsTestModalOpen] = useState(false);
  const [mode, setMode] = useState('add'); // 'add' | 'edit' | 'delete'
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const [copyTarget, setCopyTarget] = useState(null);
  const [copyContest, setCopyContest] = useState(0);
  const [flagVersion, setFlagVersion] = useState(0);
  const [editChallenge, setEditChallenge] = useState(challengeDraft());
  const [initDockerCompose, setInitDockerCompose] = useState([]);
  const fileInputRef = useRef(null);
  const uploadTargetRef = useRef(null);
  const pageSize = 20;
  const { t } = useTranslation();

  const [nameQuery, setNameQuery] = useState('');
  const [descQuery, setDescQuery] = useState('');
  const [contestQuery, setContestQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [resolvedSearchSignature, setResolvedSearchSignature] = useState('');

  const debouncedName = useDebounce(nameQuery, 300);
  const debouncedDesc = useDebounce(descQuery, 300);
  const debouncedContest = useDebounce(contestQuery, 300);

  // 搜索结果对应的查询签名，未匹配时说明这批结果还在加载
  const searchSignature = `${debouncedName}|${debouncedDesc}|${debouncedContest}|${selectedType}|${selectedCategory}`;
  const searchLoading = searchSignature !== resolvedSearchSignature;

  const isSearchMode = !!(nameQuery.trim() || descQuery.trim() || contestQuery.trim());

  const contestParams = () => (debouncedContest.trim() ? { contest_name: debouncedContest.trim() } : {});

  useEffect(() => {
    let cancelled = false;
    if (!debouncedName.trim() && !debouncedDesc.trim() && !debouncedContest.trim()) {
      setSearchResults([]);
      setResolvedSearchSignature(searchSignature);
      return;
    }
    const doSearch = async () => {
      try {
        const params = { limit: 10, offset: 0, ...contestParams() };
        if (debouncedName.trim()) params.name = debouncedName.trim();
        if (debouncedDesc.trim()) params.description = debouncedDesc.trim();
        if (selectedType !== 'all') params.type = selectedType;
        if (selectedCategory !== 'all') params.category = selectedCategory;
        const response = await getChallengeList(params);
        if (!cancelled && response.code === 200) {
          setSearchResults(response.data.challenges || []);
        }
      } catch (error) {
        if (!cancelled) {
          toast.danger({
            description: error.message || t('admin.challenge.toast.searchFailed'),
          });
          setSearchResults([]);
        }
      } finally {
        if (!cancelled) setResolvedSearchSignature(searchSignature);
      }
    };
    doSearch();
    return () => {
      cancelled = true;
    };
  }, [debouncedName, debouncedDesc, debouncedContest, selectedType, selectedCategory]);

  const fetchCategories = async () => {
    try {
      const response = await getChallengeCategories(contestParams());
      if (response.code === 200) {
        setCategories(mergeChallengeCategories(response.data));
      } else {
        setCategories(DEFAULT_CHALLENGE_CATEGORIES);
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.challenge.toast.fetchCategoriesFailed'),
      });
      setCategories(DEFAULT_CHALLENGE_CATEGORIES);
    }
  };

  const fetchChallenges = async () => {
    setLoading(true);
    try {
      const params = {
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      };

      if (selectedType !== 'all') {
        params.type = selectedType;
      }

      if (selectedCategory !== 'all') {
        params.category = selectedCategory;
      }

      const response = await getChallengeList(params);
      if (response.code === 200) {
        setInitDockerCompose(
          response.data.challenges.map((challenge) => {
            return { id: challenge.id, value: challenge.docker_compose };
          })
        );
        setChallenges(response.data.challenges);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.challenge.toast.fetchListFailed'),
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getContestList({ limit: 100, offset: 0 })
      .then((response) => {
        if (response.code === 200) setContests(response.data?.contests || []);
      })
      .catch(() => {
        // 比赛列表加载失败不阻断题库浏览
      });
  }, []);

  useEffect(() => {
    fetchCategories();
  }, [debouncedContest, selectedType]);

  useEffect(() => {
    if (!isSearchMode) {
      fetchChallenges();
    }
  }, [currentPage, selectedType, selectedCategory, isSearchMode]);

  const handleFilterTypeChange = (type) => {
    setSelectedType(type);
    setCurrentPage(1);
    setNameQuery('');
    setDescQuery('');
  };

  const handleFilterCategoryChange = (category) => {
    setSelectedCategory(category);
    setCurrentPage(1);
    setNameQuery('');
    setDescQuery('');
  };

  const handleAddChallenge = () => {
    setMode('add');
    setSelectedChallenge(null);
    setEditChallenge(challengeDraft({ contest_id: 0, hidden: false, flags: [{ id: 0, value: '' }] }));
    setIsModalOpen(true);
  };

  const handleEditChallenge = async (challenge) => {
    try {
      const response = await getChallenge(challenge.id);
      if (response.code === 200) {
        const detail = challengeDraft(response.data);
        setMode('edit');
        setSelectedChallenge(detail);
        setEditChallenge(detail);
        setInitDockerCompose((prev) => {
          const next = prev.filter((item) => item.id !== detail.id);
          next.push({ id: detail.id, value: detail.docker_compose });
          return next;
        });
        setIsModalOpen(true);
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.challenge.toast.fetchDetailFailed'),
      });
    }
  };

  const handleDeleteChallenge = (challenge) => {
    setMode('delete');
    setSelectedChallenge(challenge);
    setEditChallenge(challengeDraft(challenge));
    setIsModalOpen(true);
  };

  const handleCopyChallenge = (challenge) => {
    setCopyTarget(challenge);
    setCopyContest(challenge.contest_id || 0);
  };

  const handleConfirmCopy = async () => {
    if (!copyTarget) return;
    try {
      const response = await copyChallenge(copyTarget.id, copyContest);
      if (response.code === 200) {
        toast.success({
          description: t('admin.challenge.toast.copySuccess'),
        });
        setCopyTarget(null);
        await fetchChallenges();
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.challenge.toast.copyFailed'),
      });
    }
  };

  const handleUploadAttachment = (challenge) => {
    uploadTargetRef.current = challenge;
    setSelectedChallenge(challenge);
    fileInputRef.current?.click();
  };

  const handleDownloadAttachment = async (challenge) => {
    try {
      const response = await downloadChallengeFile(challenge.id);
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, 'attachment.zip', 'application/octet-stream');
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.challenge.toast.downloadFailed'),
      });
    }
  };

  const handleTestChallenge = (challenge) => {
    setSelectedChallenge(challenge);
    setIsTestModalOpen(true);
  };

  const handleFileChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file) return;

    const target = uploadTargetRef.current;
    if (!target) return;

    try {
      const response = await uploadChallengeFile(target.id, file);
      if (response.code === 200) {
        toast.success({
          description: t('admin.challenge.toast.uploadSuccess'),
        });
        await fetchChallenges();
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.challenge.toast.uploadFailed'),
      });
    }
  };

  const handleChallengeChange = (updatedChallenge) => {
    setEditChallenge((previous) => applyChallengeType(previous, updatedChallenge));
  };

  const handleAddFlag = () => {
    if (editChallenge.type === 'pods') return;

    setEditChallenge({
      ...editChallenge,
      flags: [...(editChallenge.flags || []), { id: 0, value: '' }],
    });
  };

  const handleRemoveFlag = (index) => {
    if (editChallenge.type === 'pods') return;

    const newFlags = [...(editChallenge.flags || [])];
    newFlags.splice(index, 1);
    setEditChallenge({
      ...editChallenge,
      flags: newFlags,
    });
  };

  const handleFlagChange = (index, value) => {
    if (editChallenge.type === 'pods') return;

    const newFlags = [...(editChallenge.flags || [])];
    newFlags[index] = {
      id: (editChallenge.flags[index] && editChallenge.flags[index].id) || 0,
      value: value,
    };
    setEditChallenge({
      ...editChallenge,
      flags: newFlags,
    });
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  const handleCloseTestModal = () => {
    setIsTestModalOpen(false);
  };

  const handleSubmitChallenge = async (challenge) => {
    try {
      let response;

      if (mode === 'delete') {
        response = await deleteChallenge(selectedChallenge.id);
      } else {
        const apiChallenge = buildChallengePayload(challenge, mode, initDockerCompose);

        if (mode === 'add') {
          response = await createChallenge(apiChallenge);
        } else if (mode === 'edit') {
          response = await updateChallenge(selectedChallenge.id, apiChallenge);
        }
      }

      if (response.code === 200) {
        const actionKey = mode === 'add' ? 'create' : mode === 'edit' ? 'update' : 'delete';
        toast.success({
          description: t(`admin.challenge.toast.${actionKey}Success`),
        });
        setIsModalOpen(false);
        if (mode === 'edit') setFlagVersion((previous) => previous + 1);
        await fetchChallenges();
      }
    } catch (error) {
      const actionKey = mode === 'add' ? 'create' : mode === 'edit' ? 'update' : 'delete';
      toast.danger({
        description: error.message || t(`admin.challenge.toast.${actionKey}Failed`),
      });
    }
  };

  const displayChallenges = isSearchMode ? searchResults : challenges;
  const displayTotalCount = isSearchMode ? searchResults.length : totalCount;

  return (
    <>
      <AdminChallenge
        challenges={displayChallenges}
        totalCount={displayTotalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        categories={categories}
        filterCategory={selectedCategory}
        filterType={selectedType}
        loading={isSearchMode ? searchLoading : loading}
        onPageChange={setCurrentPage}
        onAddChallenge={handleAddChallenge}
        onEditChallenge={handleEditChallenge}
        onDeleteChallenge={handleDeleteChallenge}
        onCopyChallenge={handleCopyChallenge}
        onUploadAttachment={handleUploadAttachment}
        onDownloadAttachment={handleDownloadAttachment}
        onTestChallenge={handleTestChallenge}
        onFilterTypeChange={handleFilterTypeChange}
        onFilterCategoryChange={handleFilterCategoryChange}
        nameQuery={nameQuery}
        descQuery={descQuery}
        contestQuery={contestQuery}
        searchLoading={searchLoading}
        isSearchMode={isSearchMode}
        onNameChange={setNameQuery}
        onDescChange={setDescQuery}
        onContestChange={setContestQuery}
      />

      <ChallengeEditorDialog
        isOpen={isModalOpen}
        mode={mode}
        challenge={editChallenge}
        categories={categories}
        contests={contests}
        flagVersion={flagVersion}
        onClose={handleCloseModal}
        onSubmit={handleSubmitChallenge}
        onChange={handleChallengeChange}
        onAddFlag={handleAddFlag}
        onRemoveFlag={handleRemoveFlag}
        onFlagChange={handleFlagChange}
      />

      <ChallengeCopyDialog
        isOpen={!!copyTarget}
        challenge={copyTarget}
        contests={contests}
        contestId={copyContest}
        onContestChange={setCopyContest}
        onClose={() => setCopyTarget(null)}
        onConfirm={handleConfirmCopy}
      />

      <AdminChallengeTestModal challenge={selectedChallenge} isOpen={isTestModalOpen} onClose={handleCloseTestModal} />

      <input type="file" ref={fileInputRef} className="hidden" onChange={handleFileChange} />
    </>
  );
}

export default ChallengesManagement;
