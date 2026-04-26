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
import { useDebounce } from '../../hooks';
import { useTranslation } from 'react-i18next';
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
    };
  };

  // 使用防抖搜索
  const [nameQuery, setNameQuery] = useState('');
  const [descQuery, setDescQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [searchLoading, setSearchLoading] = useState(false);

  const debouncedName = useDebounce(nameQuery, 300);
  const debouncedDesc = useDebounce(descQuery, 300);

  // 搜索结果和模式
  const isSearchMode = !!(nameQuery.trim() || descQuery.trim());

  useEffect(() => {
    let cancelled = false;
    if (!debouncedName.trim() && !debouncedDesc.trim()) {
      setSearchResults([]);
      return;
    }
    const doSearch = async () => {
      setSearchLoading(true);
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
        if (!cancelled) setSearchLoading(false);
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
      toast.danger({
        description: error.message || t('admin.challenge.toast.fetchListFailed'),
      });
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

  const handleEditChallenge = (challenge) => {
    setMode('edit');
    setSelectedChallenge(challenge);
    setEditChallenge({
      ...challenge,
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

    try {
      const response = await uploadChallengeFile(selectedChallenge.id, file);
      if (response.code === 200) {
        toast.success({
          description: t('admin.challenge.toast.uploadSuccess'),
        });
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.challenge.toast.uploadFailed'),
      });
    }
  };

  const handleChallengeChange = (updatedChallenge) => {
    // 如果题目类型发生变化, 需要相应调整flags
    if (updatedChallenge.type !== editChallenge.type) {
      if (updatedChallenge.type === 'pods') {
        // 切换到pods类型时清空flags
        updatedChallenge.flags = [];
      } else if (editChallenge.type === 'pods') {
        // 从pods类型切换到其他类型时, 确保有默认flags
        updatedChallenge.flags = [{ id: 0, value: '' }];
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

        // 处理flags格式, 根据模式和类型确定正确的格式
        let processedFlags;

        if (challenge.type === 'pods') {
          // pods类型不需要flags, 由后端自动生成
          processedFlags = [];
        } else if (mode === 'add') {
          // 创建模式: flags是字符串数组
          processedFlags = challenge.flags.map((flag) => {
            return typeof flag === 'string' ? flag : flag.value || '';
          });
        } else if (mode === 'edit') {
          // 编辑模式: flags是对象数组, 包含id和value
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
        }

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
        nameQuery={nameQuery}
        descQuery={descQuery}
        searchLoading={searchLoading}
        isSearchMode={isSearchMode}
        onNameChange={setNameQuery}
        onDescChange={setDescQuery}
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
      />

      <AdminChallengeTestModal challenge={selectedChallenge} isOpen={isTestModalOpen} onClose={handleCloseTestModal} />

      <input type="file" ref={fileInputRef} className="hidden" onChange={handleFileChange} />
    </>
  );
}

export default ChallengesManagement;
