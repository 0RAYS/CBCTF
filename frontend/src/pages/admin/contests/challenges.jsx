import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import AdminContestChallenges from '../../../components/features/Admin/Contests/AdminContestChallenges';
import AdminContestChallengeSelector from '../../../components/features/Admin/Contests/AdminContestChallengeSelector';
import AdminContestChallengeModal from '../../../components/features/Admin/Contests/AdminContestChallengeModal';
import {
  getContestChallenges,
  addContestChallenge,
  updateContestChallenge,
  removeContestChallenge,
} from '../../../api/admin/contest';
import {
  getNotInContestChallengeList,
  getChallengeFlags,
  updateChallengeFlag,
  getContestChallengeCategories,
} from '../../../api/admin/challenge';
import { useTranslation } from 'react-i18next';
import { DEFAULT_CHALLENGE_CATEGORIES, mergeChallengeCategories } from '../../../config/challenges';

function AdminContestChallengesPage() {
  const { id } = useParams();
  // 赛题列表状态
  const [challenges, setChallenges] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);

  // 过滤状态
  const [categories, setCategories] = useState([]);
  const [filterType, setFilterType] = useState('all');
  const [filterCategory, setFilterCategory] = useState('all');
  const [nameQuery, setNameQuery] = useState('');

  // 选择器状态
  const [isSelectorOpen, setSelectorOpen] = useState(false);
  const [availableChallenges, setAvailableChallenges] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [descSearchQuery, setDescSearchQuery] = useState('');
  const [selectedChallenges, setSelectedChallenges] = useState([]);
  const [modalType, setModalType] = useState('all');
  const [modalCategory, setModalCategory] = useState('all');
  const [modalCurrentPage, setModalCurrentPage] = useState(1);
  const [modalTotalCount, setModalTotalCount] = useState(0);
  const [modalLoading, setModalLoading] = useState(false);
  const modalPageSize = 10; // 模态框分页大小

  // 编辑弹窗状态
  const [isModalOpen, setModalOpen] = useState(false);
  const [editingChallenge, setEditingChallenge] = useState(null);
  const [adminContestChallengeFlags, setAdminContestChallengeFlags] = useState([
    {
      name: '',
      description: '',
      attempt: 0,
      hidden: false,
      hints: [''],
      tags: [''],
    },
  ]);
  const [editForm, setEditForm] = useState({
    name: '',
    description: '',
    attempt: 0,
    hidden: false,
    hints: [],
    tags: [],
  });
  const { t } = useTranslation();

  // 分页配置
  const pageSize = 10;

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    if (id) {
      fetchChallengesWithFilters(1);
    }
  }, [id]); // 只在id变化时获取数据

  // 当过滤器变化时应用过滤
  useEffect(() => {
    if (id) {
      filterChallenges();
    }
  }, [filterType, filterCategory, nameQuery]);

  // 当当前页变化时, 重新获取数据
  useEffect(() => {
    if (id && currentPage > 0) {
      fetchChallengesWithFilters(currentPage);
    }
  }, [currentPage, id]);

  // 获取分类列表
  const fetchCategories = async () => {
    try {
      const response = await getContestChallengeCategories(parseInt(id));
      if (response.code === 200) {
        setCategories(mergeChallengeCategories(response.data));
      } else {
        // API请求失败时仍然设置默认分类
        setCategories(DEFAULT_CHALLENGE_CATEGORIES);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.fetchCategoriesFailed') });
      // 出现异常时仍然设置默认分类
      setCategories(DEFAULT_CHALLENGE_CATEGORIES);
    }
  };

  // 过滤挑战列表
  const filterChallenges = () => {
    // 重置到第一页并重新获取数据
    setCurrentPage(1);
    fetchChallengesWithFilters(1);
  };

  // 获取带过滤器的挑战数据
  const fetchChallengesWithFilters = async (page) => {
    try {
      const params = {
        limit: pageSize,
        offset: (page - 1) * pageSize,
      };

      // 添加过滤参数
      if (filterType !== 'all') {
        params.type = filterType;
      }

      if (filterCategory !== 'all') {
        params.category = filterCategory;
      }
      if (nameQuery.trim()) {
        params.name = nameQuery.trim();
      }

      const response = await getContestChallenges(parseInt(id), params);
      if (response.code === 200) {
        const challengeData = response.data === null ? [] : response.data;
        setChallenges(challengeData.challenges || []);
        setTotalCount(challengeData.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.fetchListFailed') });
    }
  };

  // 获取可用题目列表（使用API分页）
  const fetchAvailableChallenges = async (
    page = modalCurrentPage,
    type = modalType,
    category = modalCategory,
    query = '',
    descQuery = ''
  ) => {
    setModalLoading(true);
    try {
      // 使用API分页参数
      const params = {
        limit: modalPageSize,
        offset: (page - 1) * modalPageSize,
      };

      // 添加过滤参数
      if (type !== 'all') {
        params.type = type;
      }

      if (category !== 'all') {
        params.category = category;
      }
      if (query.trim()) {
        params.name = query.trim();
      }
      if (descQuery.trim()) {
        params.description = descQuery.trim();
      }
      const response = await getNotInContestChallengeList(parseInt(id), params);
      if (response.code === 200) {
        setAvailableChallenges(response.data.challenges || []);
        setModalTotalCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.fetchPoolFailed') });
    } finally {
      setModalLoading(false);
    }
  };

  // 处理类型过滤变化
  const handleFilterTypeChange = (type) => {
    setFilterType(type);
  };

  // 处理分类过滤变化
  const handleFilterCategoryChange = (category) => {
    setFilterCategory(category);
  };

  // 处理模态框类型过滤变化
  const handleModalTypeChange = (type) => {
    setModalType(type);
    setModalCurrentPage(1); // 重置页码
    fetchAvailableChallenges(1, type, modalCategory, searchQuery, descSearchQuery);
  };

  // 处理模态框分类过滤变化
  const handleModalCategoryChange = (category) => {
    setModalCategory(category);
    setModalCurrentPage(1); // 重置页码
    fetchAvailableChallenges(1, modalType, category, searchQuery, descSearchQuery);
  };

  // 处理搜索查询变化
  const handleSearchChange = (query) => {
    setSearchQuery(query);
    setModalCurrentPage(1); // 重置页码
    fetchAvailableChallenges(1, modalType, modalCategory, query, descSearchQuery);
  };

  // 处理描述搜索查询变化
  const handleDescSearchChange = (query) => {
    setDescSearchQuery(query);
    setModalCurrentPage(1); // 重置页码
    fetchAvailableChallenges(1, modalType, modalCategory, searchQuery, query);
  };

  // 处理模态框分页变化
  const handleModalPageChange = async (page) => {
    setModalCurrentPage(page);
    fetchAvailableChallenges(page, modalType, modalCategory, searchQuery, descSearchQuery);
  };

  // 选择器相关处理
  const handleOpenSelector = async () => {
    // 重置选择器状态
    setSelectedChallenges([]);
    setModalType('all');
    setModalCategory('all');
    setModalCurrentPage(1);

    // 获取可用题目列表
    await fetchAvailableChallenges(1, 'all', 'all', '', '');

    // 打开选择器
    setSelectorOpen(true);
  };

  const handleSelectChallenge = (challenge) => {
    if (selectedChallenges.some((c) => c.id === challenge.id)) {
      setSelectedChallenges(selectedChallenges.filter((c) => c.id !== challenge.id));
    } else {
      setSelectedChallenges([...selectedChallenges, challenge]);
    }
  };

  const handleAddChallenges = async () => {
    if (selectedChallenges.length === 0) {
      toast.warning({ description: t('admin.contests.challenges.toast.selectRequired') });
      return;
    }

    try {
      const response = await addContestChallenge(
        parseInt(id),
        selectedChallenges.map((c) => c.id)
      );
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.challenges.toast.addSuccess') });
        setSelectorOpen(false);
        await fetchChallengesWithFilters(currentPage);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.addFailed') });
    }
  };

  // 编辑弹窗相关处理
  const handleEditChallenge = async (challenge) => {
    setEditingChallenge(challenge);
    setEditForm({
      name: challenge.name || '',
      description: challenge.description || '',
      attempt: challenge.attempt || 0,
      hidden: challenge.hidden || false,
      hints: challenge.hints || [],
      tags: challenge.tags || [],
    });
    const response = await getChallengeFlags(parseInt(id), challenge.id);
    if (response.code === 200) {
      setAdminContestChallengeFlags(response.data);
    }
    setModalOpen(true);
  };

  const handleFormChange = (updatedForm) => {
    setEditForm(updatedForm);
  };

  const handleAddHint = () => {
    setEditForm({
      ...editForm,
      hints: [...editForm.hints, ''],
    });
  };

  const handleRemoveHint = (index) => {
    const newHints = [...editForm.hints];
    newHints.splice(index, 1);
    setEditForm({
      ...editForm,
      hints: newHints,
    });
  };

  const handleHintChange = (index, value) => {
    const newHints = [...editForm.hints];
    newHints[index] = value;
    setEditForm({
      ...editForm,
      hints: newHints,
    });
  };

  const handleAddTag = () => {
    setEditForm({
      ...editForm,
      tags: [...editForm.tags, ''],
    });
  };

  const handleRemoveTag = (index) => {
    const newTags = [...editForm.tags];
    newTags.splice(index, 1);
    setEditForm({
      ...editForm,
      tags: newTags,
    });
  };

  const handleTagChange = (index, value) => {
    const newTags = [...editForm.tags];
    newTags[index] = value;
    setEditForm({
      ...editForm,
      tags: newTags,
    });
  };

  const handleUpdateChallenge = async (updatedForm) => {
    if (!editingChallenge) return;

    try {
      const response = await updateContestChallenge(parseInt(id), editingChallenge.id, updatedForm);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.challenges.toast.updateSuccess') });
        setModalOpen(false);
        await fetchChallengesWithFilters(currentPage);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.updateFailed') });
    }
  };

  const handleRemoveChallenge = async (challenge) => {
    try {
      const response = await removeContestChallenge(parseInt(id), challenge.id);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.challenges.toast.removeSuccess') });
        await fetchChallengesWithFilters(currentPage);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.challenges.toast.removeFailed') });
    }
  };

  const handleFlagChange = async (index, value) => {
    const response = await updateChallengeFlag(
      parseInt(id),
      editingChallenge.id,
      adminContestChallengeFlags[index].id,
      value
    );
    if (response.code === 200) {
      toast.success({ description: t('admin.contests.challenges.toast.updateFlagSuccess') });
    }
  };

  return (
    <>
      <AdminContestChallenges
        challenges={challenges}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        categories={categories}
        filterCategory={filterCategory}
        filterType={filterType}
        onPageChange={setCurrentPage}
        onAddChallenge={handleOpenSelector}
        onEditChallenge={handleEditChallenge}
        onDeleteChallenge={handleRemoveChallenge}
        onFilterTypeChange={handleFilterTypeChange}
        onFilterCategoryChange={handleFilterCategoryChange}
        nameQuery={nameQuery}
        onNameChange={setNameQuery}
      />

      {/* 选择赛题弹窗 */}
      <AdminContestChallengeSelector
        isOpen={isSelectorOpen}
        challenges={availableChallenges}
        selectedChallenges={selectedChallenges}
        categories={categories}
        totalCount={modalTotalCount}
        currentPage={modalCurrentPage}
        pageSize={modalPageSize}
        loading={modalLoading}
        onClose={() => setSelectorOpen(false)}
        onSearch={handleSearchChange}
        onDescSearch={handleDescSearchChange}
        onSelect={handleSelectChallenge}
        onConfirm={handleAddChallenges}
        onPageChange={handleModalPageChange}
        onFilterCategoryChange={handleModalCategoryChange}
        onFilterTypeChange={handleModalTypeChange}
      />

      {/* 编辑赛题弹窗 */}
      <AdminContestChallengeModal
        isOpen={isModalOpen}
        mode="edit"
        contestId={parseInt(id)}
        challengeId={editingChallenge?.id}
        challenge={editForm}
        flags={adminContestChallengeFlags}
        onClose={() => setModalOpen(false)}
        onSubmit={handleUpdateChallenge}
        onChange={handleFormChange}
        onAddHint={handleAddHint}
        onRemoveHint={handleRemoveHint}
        onHintChange={handleHintChange}
        onAddTag={handleAddTag}
        onRemoveTag={handleRemoveTag}
        onTagChange={handleTagChange}
        onFlagChange={handleFlagChange}
      />
    </>
  );
}

export default AdminContestChallengesPage;
