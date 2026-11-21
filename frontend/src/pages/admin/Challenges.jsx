import { useEffect, useRef, useState } from 'react';
import { toast } from '@/utils/toast';
import { downloadBlobResponse } from '@/utils/fileDownload';
import AdminChallenge from '@/components/features/Admin/AdminChallenge.jsx';
import ChallengeEditorDialog from '@/components/features/Admin/challenges/editor/ChallengeEditorDialog.jsx';
import { buildChallengePayload } from '@/components/features/Admin/challenges/editor/challengePayload.js';
import AdminChallengeTestModal from '@/components/features/Admin/challenges/testing/ChallengeTestDialog.jsx';
import {
  createChallenge,
  deleteChallenge,
  downloadChallengeFile,
  getChallenge,
  getChallengeCategories,
  getChallengeList,
  getChallengeUsages,
  updateChallenge,
  uploadChallengeFile,
} from '@/api/admin/challenge';
import { useDebounce } from '@/hooks';
import { useTranslation } from 'react-i18next';
import { DEFAULT_CHALLENGE_CATEGORIES, mergeChallengeCategories } from '@/config/challenges';

function ChallengesManagement() {
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [selectedType, setSelectedType] = useState('all');
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isTestModalOpen, setIsTestModalOpen] = useState(false);
  const [mode, setMode] = useState('add'); // 'add' | 'edit' | 'delete'
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  // 删除题库题目之前先查清引用范围
  const [deleteUsages, setDeleteUsages] = useState([]);
  const [editChallenge, setEditChallenge] = useState({
    name: '',
    description: '',
    category: '',
    type: 'static',
    flags: [{ id: 0, value: '' }],
    docker_compose: '',
    network_policies: [],
  });
  const [initDockerCompose, setInitDockerCompose] = useState([]);
  const fileInputRef = useRef(null);
  const uploadTargetRef = useRef(null);
  const pageSize = 20;
  const { t } = useTranslation();

  const processFlags = (challenge) => {
    let flags;
    if (challenge.flags && Array.isArray(challenge.flags)) {
      flags = challenge.flags;
    } else {
      flags = [''];
    }
    return {
      ...challenge,
      flags,
    };
  };

  const [nameQuery, setNameQuery] = useState('');
  const [descQuery, setDescQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [resolvedSearchSignature, setResolvedSearchSignature] = useState('');

  const debouncedName = useDebounce(nameQuery, 300);
  const debouncedDesc = useDebounce(descQuery, 300);

  // 搜索结果对应的查询签名，未匹配时说明这批结果还在加载
  const searchSignature = `${debouncedName}|${debouncedDesc}|${selectedType}|${selectedCategory}`;
  const searchLoading = searchSignature !== resolvedSearchSignature;

  const isSearchMode = !!(nameQuery.trim() || descQuery.trim());

  useEffect(() => {
    let cancelled = false;
    if (!debouncedName.trim() && !debouncedDesc.trim()) {
      setSearchResults([]);
      setResolvedSearchSignature(searchSignature);
      return;
    }
    const doSearch = async () => {
      try {
        const params = { limit: 10, offset: 0 };
        if (debouncedName.trim()) params.name = debouncedName.trim();
        if (debouncedDesc.trim()) params.description = debouncedDesc.trim();
        if (selectedType !== 'all') params.type = selectedType;
        if (selectedCategory !== 'all') params.category = selectedCategory;
        const response = await getChallengeList(params);
        if (!cancelled && response.code === 200) {
          setSearchResults((response.data.challenges || []).map(processFlags));
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
  }, [debouncedName, debouncedDesc, selectedType, selectedCategory]);

  const fetchCategories = async () => {
    try {
      const response = await getChallengeCategories();
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
        setChallenges(response.data.challenges.map(processFlags));
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
    fetchCategories();
  }, []);

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
    setEditChallenge({
      name: '',
      description: '',
      category: '',
      type: 'static',
      flags: [{ id: 0, value: '' }], // 静态类型默认需要flags
      docker_compose: '',
      network_policies: [],
    });
    setIsModalOpen(true);
  };

  const handleEditChallenge = async (challenge) => {
    try {
      const response = await getChallenge(challenge.id);
      if (response.code === 200) {
        const detail = processFlags(response.data);
        setMode('edit');
        setSelectedChallenge(detail);
        setEditChallenge({
          ...detail,
        });
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
    setDeleteUsages([]);
    setIsModalOpen(true);
    getChallengeUsages(challenge.id)
      .then((response) => {
        if (response.code === 200) setDeleteUsages(response.data?.usages || []);
      })
      .catch(() => {
        // 引用范围查询失败不阻断删除
      });
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
    if (updatedChallenge.type !== editChallenge.type) {
      if (updatedChallenge.type === 'pods') {
        updatedChallenge.flags = [];
      } else if (editChallenge.type === 'pods') {
        updatedChallenge.flags = [{ id: 0, value: '' }];
      }
    }

    setEditChallenge(updatedChallenge);
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
        const keptCount = response.data?.contest_challenge_count || 0;
        if (mode === 'delete' && keptCount > 0) {
          toast.success({
            description: t('admin.challenge.toast.deleteKeepsContestCopies', { count: keptCount }),
          });
        }
        setIsModalOpen(false);
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
        onUploadAttachment={handleUploadAttachment}
        onDownloadAttachment={handleDownloadAttachment}
        onTestChallenge={handleTestChallenge}
        onFilterTypeChange={handleFilterTypeChange}
        onFilterCategoryChange={handleFilterCategoryChange}
        nameQuery={nameQuery}
        descQuery={descQuery}
        searchLoading={searchLoading}
        isSearchMode={isSearchMode}
        onNameChange={setNameQuery}
        onDescChange={setDescQuery}
      />

      <ChallengeEditorDialog
        isOpen={isModalOpen}
        mode={mode}
        challenge={editChallenge}
        usages={deleteUsages}
        categories={categories}
        onClose={handleCloseModal}
        onSubmit={handleSubmitChallenge}
        onChange={handleChallengeChange}
        onAddFlag={handleAddFlag}
        onRemoveFlag={handleRemoveFlag}
        onFlagChange={handleFlagChange}
      />

      <AdminChallengeTestModal challenge={selectedChallenge} isOpen={isTestModalOpen} onClose={handleCloseTestModal} />

      <input type="file" ref={fileInputRef} className="hidden" onChange={handleFileChange} />
    </>
  );
}

export default ChallengesManagement;
