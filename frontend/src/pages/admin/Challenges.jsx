import { useState, useEffect, useRef } from 'react';
import { toast } from '../../utils/toast';
import { downloadBlobResponse } from '../../utils/fileDownload';
import AdminChallenge from '../../components/features/Admin/AdminChallenge.jsx';
import AdminChallengeModal from '../../components/features/Admin/AdminChallengeModal.jsx';
import AdminChallengeTestModal from '../../components/features/Admin/AdminChallengeTestModal.jsx';
import {
  getChallengeCategories,
  getChallengeList,
  createChallenge,
  updateChallenge,
  deleteChallenge,
  uploadChallengeFile,
  downloadChallengeFile,
} from '../../api/admin/challenge';
import { generateUUID } from '../../utils/uuid';
import { useDebounceSearch } from '../../hooks';
import { useTranslation } from 'react-i18next';
import { searchModels } from '../../api/admin/search.js';
import { DEFAULT_CHALLENGE_CATEGORIES, mergeChallengeCategories } from '../../config/challenges';

function ChallengesManagement() {
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedType, setSelectedType] = useState('all');
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isTestModalOpen, setIsTestModalOpen] = useState(false);
  const [mode, setMode] = useState('add'); // 'add' | 'edit' | 'delete'
  const [selectedChallenge, setSelectedChallenge] = useState(null);
  const [editChallenge, setEditChallenge] = useState({
    name: '',
    description: '',
    category: '',
    type: 'static',
    flags: [{ id: 0, value: '' }],
    options: [],
    docker_compose: '',
    network_policies: [],
  });
  const [initDockerCompose, setInitDockerCompose] = useState([]);
  const fileInputRef = useRef(null);
  const pageSize = 20;
  const { t } = useTranslation();

  // 处理 flags 字段的辅助函数
  const processFlags = (challenge) => {
    let flags;
    if (challenge.flags && Array.isArray(challenge.flags)) {
      flags = challenge.flags;
    } else if (challenge.flag) {
      flags = [challenge.flag];
    } else {
      flags = [''];
    }
    return {
      ...challenge,
      flags,
      options: challenge.options || [],
    };
  };

  // 使用防抖搜索 Hook
  const {
    query: searchQuery,
    setQuery: setSearchQuery,
    results: rawSearchResults,
    loading: searchLoading,
  } = useDebounceSearch(
    async (name) => {
      if (!name || name.trim() === '') return [];

      try {
        const params = {
          model: 'Challenge',
          'search[name]': name.trim(),
          limit: 10,
          offset: 0,
        };
        if (selectedType !== 'all') params['search[type]'] = selectedType;
        if (selectedCategory !== 'all') params['search[category]'] = selectedCategory;
        const response = await searchModels(params);

        if (response.code === 200) {
          const results = response.data.models || [];
          return results.map(processFlags);
        }
        return [];
      } catch (error) {
        toast.danger({ description: error.message || t('admin.challenge.toast.searchFailed') });
        return [];
      }
    },
    { delay: 300, minLength: 1 }
  );

  // 搜索结果和模式
  const searchResults = rawSearchResults || [];
  const isSearchMode = searchQuery.trim().length > 0;

  const fetchCategories = async () => {
    try {
      const response = await getChallengeCategories();
      if (response.code === 200) {
        setCategories(mergeChallengeCategories(response.data));
      } else {
        setCategories(DEFAULT_CHALLENGE_CATEGORIES);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.challenge.toast.fetchCategoriesFailed') });
      // 出现异常时仍然设置默认分类
      setCategories(DEFAULT_CHALLENGE_CATEGORIES);
    }
  };

  const fetchChallenges = async () => {
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
      toast.danger({ description: error.message || t('admin.challenge.toast.fetchListFailed') });
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
    setSearchQuery('');
  };

  const handleFilterCategoryChange = (category) => {
    setSelectedCategory(category);
    setCurrentPage(1);
    setSearchQuery('');
  };

  const handleAddChallenge = () => {
    setMode('add');
    setEditChallenge({
      name: '',
      description: '',
      category: '',
      type: 'static',
      flags: [{ id: 0, value: '' }], // 静态类型默认需要flags
      options: [], // question类型需要选项
      docker_compose: '',
      network_policies: [],
    });
    setIsModalOpen(true);
  };

  const handleEditChallenge = (challenge) => {
    setMode('edit');
    setSelectedChallenge(challenge);
    setEditChallenge({
      ...challenge,
      options: challenge.options || [], // 确保options字段存在
    });
    setIsModalOpen(true);
  };

  const handleDeleteChallenge = (challenge) => {
    setMode('delete');
    setSelectedChallenge(challenge);
    setIsModalOpen(true);
  };

  const handleUploadAttachment = (challenge) => {
    fileInputRef.current?.click();
    setSelectedChallenge(challenge);
  };

  const handleDownloadAttachment = async (challenge) => {
    try {
      const response = await downloadChallengeFile(challenge.id);
      downloadBlobResponse(response, 'attachment.zip', 'application/octet-stream');
    } catch (error) {
      toast.danger({ description: error.message || t('admin.challenge.toast.downloadFailed') });
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

    try {
      const response = await uploadChallengeFile(selectedChallenge.id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.challenge.toast.uploadSuccess') });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.challenge.toast.uploadFailed') });
    }
  };

  const handleChallengeChange = (updatedChallenge) => {
    // 如果题目类型发生变化，需要相应调整flags和options
    if (updatedChallenge.type !== editChallenge.type) {
      if (updatedChallenge.type === 'pods') {
        // 切换到pods类型时清空flags
        updatedChallenge.flags = [];
      } else if (editChallenge.type === 'pods') {
        // 从pods类型切换到其他类型时，确保有默认flags
        updatedChallenge.flags = [{ id: 0, value: '' }];
      }

      // 如果切换到question类型，初始化options
      if (updatedChallenge.type === 'question') {
        updatedChallenge.options = [];
      }
    }

    setEditChallenge(updatedChallenge);
  };

  const handleAddFlag = () => {
    // pods类型不需要flags
    if (editChallenge.type === 'pods') return;

    setEditChallenge({
      ...editChallenge,
      flags: [...(editChallenge.flags || []), { id: 0, value: '' }],
    });
  };

  const handleRemoveFlag = (index) => {
    // pods类型不需要flags
    if (editChallenge.type === 'pods') return;

    const newFlags = [...(editChallenge.flags || [])];
    newFlags.splice(index, 1);
    setEditChallenge({
      ...editChallenge,
      flags: newFlags,
    });
  };

  const handleFlagChange = (index, value) => {
    // pods类型不需要flags
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

  // 选项相关处理函数
  const handleAddOption = () => {
    // 只有question类型需要选项
    if (editChallenge.type !== 'question') return;

    const newOptions = [...(editChallenge.options || [])];
    newOptions.push({
      rand_id: generateUUID(),
      content: '',
      correct: false,
    });
    setEditChallenge({
      ...editChallenge,
      options: newOptions,
    });
  };

  const handleRemoveOption = (index) => {
    // 只有question类型需要选项
    if (editChallenge.type !== 'question') return;

    const newOptions = [...(editChallenge.options || [])];
    newOptions.splice(index, 1);
    setEditChallenge({
      ...editChallenge,
      options: newOptions,
    });
  };

  const handleOptionChange = (index, value) => {
    // 只有question类型需要选项
    if (editChallenge.type !== 'question') return;

    const newOptions = [...(editChallenge.options || [])];
    newOptions[index] = {
      ...newOptions[index],
      content: value,
    };
    setEditChallenge({
      ...editChallenge,
      options: newOptions,
    });
  };

  const handleCorrectOptionChange = (index) => {
    // 只有question类型需要选项
    if (editChallenge.type !== 'question') return;

    const newOptions = [...(editChallenge.options || [])];
    // 切换当前选项的正确状态（允许多选）
    newOptions[index].correct = !newOptions[index].correct;
    setEditChallenge({
      ...editChallenge,
      options: newOptions,
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
        // 删除操作直接调用删除API
        response = await deleteChallenge(selectedChallenge.id);
      } else {
        // 添加和编辑操作需要构建challenge对象
        const apiChallenge = {
          name: challenge.name,
          description: challenge.description,
          category: challenge.category,
          type: challenge.type,
        };

        // 如果是question类型，添加options字段
        if (challenge.type === 'question') {
          apiChallenge.options = challenge.options || [];
        }

        // 处理flags格式，根据模式和类型确定正确的格式
        let processedFlags;

        if (challenge.type === 'pods') {
          // pods类型不需要flags，由后端自动生成
          processedFlags = [];
        } else if (mode === 'add') {
          // 创建模式：flags是字符串数组
          processedFlags = challenge.flags.map((flag) => {
            return typeof flag === 'string' ? flag : flag.value || '';
          });
        } else if (mode === 'edit') {
          // 编辑模式：flags是对象数组，包含id和value
          processedFlags = challenge.flags.map((flag) => {
            if (typeof flag === 'string') {
              return { id: 0, value: flag };
            } else {
              return { id: flag.id || 0, value: flag.value || '' };
            }
          });
        }

        if (challenge.type === 'static') {
          apiChallenge.flags = processedFlags || [];
        } else if (challenge.type === 'dynamic') {
          apiChallenge.flags = processedFlags || [];
          apiChallenge.generator_image = challenge.generator_image || '';
        } else if (challenge.type === 'pods') {
          // docker-compose 内容发生变化时才传递
          if (mode === 'add') {
            apiChallenge.docker_compose = challenge.docker_compose || '';
          } else {
            initDockerCompose.map((docker_compose) => {
              if (docker_compose.id === challenge.id && docker_compose.value !== challenge.docker_compose) {
                apiChallenge.docker_compose = challenge.docker_compose;
              }
            });
          }
          apiChallenge.network_policies =
            challenge.network_policies?.map((policy) => ({
              from: policy.from || [],
              to: policy.to || [],
            })) || [];
        } else if (challenge.type === 'question') {
          // 对于question类型，flags不需要，直接使用options
          apiChallenge.flags = []; // 清空flags
          apiChallenge.options = challenge.options || [];
        }

        if (mode === 'add') {
          response = await createChallenge(apiChallenge);
        } else if (mode === 'edit') {
          response = await updateChallenge(selectedChallenge.id, apiChallenge);
        }
      }

      if (response.code === 200) {
        const actionKey = mode === 'add' ? 'create' : mode === 'edit' ? 'update' : 'delete';
        toast.success({ description: t(`admin.challenge.toast.${actionKey}Success`) });
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

  // 确定要显示的题目数据
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
        onPageChange={setCurrentPage}
        onAddChallenge={handleAddChallenge}
        onEditChallenge={handleEditChallenge}
        onDeleteChallenge={handleDeleteChallenge}
        onUploadAttachment={handleUploadAttachment}
        onDownloadAttachment={handleDownloadAttachment}
        onTestChallenge={handleTestChallenge}
        onFilterTypeChange={handleFilterTypeChange}
        onFilterCategoryChange={handleFilterCategoryChange}
        searchQuery={searchQuery}
        searchLoading={searchLoading}
        isSearchMode={isSearchMode}
        onSearchChange={setSearchQuery}
      />

      <AdminChallengeModal
        isOpen={isModalOpen}
        mode={mode}
        challenge={editChallenge}
        categories={categories}
        onClose={handleCloseModal}
        onSubmit={handleSubmitChallenge}
        onChange={handleChallengeChange}
        onAddFlag={handleAddFlag}
        onRemoveFlag={handleRemoveFlag}
        onFlagChange={handleFlagChange}
        onAddOption={handleAddOption}
        onRemoveOption={handleRemoveOption}
        onOptionChange={handleOptionChange}
        onCorrectOptionChange={handleCorrectOptionChange}
      />

      <AdminChallengeTestModal challenge={selectedChallenge} isOpen={isTestModalOpen} onClose={handleCloseTestModal} />

      <input type="file" ref={fileInputRef} className="hidden" onChange={handleFileChange} />
    </>
  );
}

export default ChallengesManagement;
