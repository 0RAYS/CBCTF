import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { motion } from 'motion/react';
import { IconFilter, IconEdit, IconEye, IconTrash, IconShieldCheck } from '@tabler/icons-react';
import { toast } from '../../../utils/toast';
import {
  getContestCheats,
  updateContestCheat,
  deleteContestCheat,
  deleteAllContestCheats,
  checkContestCheats,
  getContestTeam,
  getTeamMembers,
  getContestTeamSubmissions,
  getContestTeamWriteups,
  getTeamContainers,
  downloadContainerTraffic,
  downloadContestTeamWriteup,
  getContestTeamFlags,
} from '../../../api/admin/contest';
import { getUserInfo } from '../../../api/admin/user';
import { downloadBlobResponse } from '../../../utils/fileDownload';
import { List, StatusTag, IpAddress } from '../../../components/common';
import { Modal } from '../../../components/common';
import ModalButton from '../../../components/common/ModalButton';
import { Button, Pagination, EmptyState } from '../../../components/common';
import AdminUserDetailDialog from '../../../components/features/Admin/AdminUserDetailDialog';
import AdminTeamDetailDialog from '../../../components/features/Admin/Contests/AdminTeamDetailDialog';
import { useTranslation } from 'react-i18next';

function AdminContestCheats() {
  const { id } = useParams();
  const [cheats, setCheats] = useState([]);
  const [total, setTotal] = useState(0);
  const [checked, setChecked] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(20);
  const [filterType, setFilterType] = useState('');
  const [filterReasonType, setFilterReasonType] = useState('');
  const [showFilter, setShowFilter] = useState(false);
  const [selectedCheat, setSelectedCheat] = useState(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDetailModal, setShowDetailModal] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [showDeleteAllConfirm, setShowDeleteAllConfirm] = useState(false);
  const [showCheckConfirm, setShowCheckConfirm] = useState(false);
  const [checkLoading, setCheckLoading] = useState(false);
  const [editForm, setEditForm] = useState({
    reason: '',
    type: '',
    checked: false,
    comment: '',
  });
  const { t, i18n } = useTranslation();

  // User detail dialog state
  const [showUserDetail, setShowUserDetail] = useState(false);
  const [userDetailData, setUserDetailData] = useState(null);

  // Team detail dialog state
  const [showTeamDetail, setShowTeamDetail] = useState(false);
  const [teamDetailData, setTeamDetailData] = useState(null);
  const [teamDetailTab, setTeamDetailTab] = useState('info');
  const [teamDetailMembers, setTeamDetailMembers] = useState([]);
  const [teamDetailMembersLoading, setTeamDetailMembersLoading] = useState(false);
  const [teamDetailSubmissions, setTeamDetailSubmissions] = useState([]);
  const [teamDetailSubmissionCount, setTeamDetailSubmissionCount] = useState(0);
  const [teamDetailSubmissionPage, setTeamDetailSubmissionPage] = useState(1);
  const [teamDetailWriteups, setTeamDetailWriteups] = useState([]);
  const [teamDetailWriteupCount, setTeamDetailWriteupCount] = useState(0);
  const [teamDetailWriteupPage, setTeamDetailWriteupPage] = useState(1);
  const [teamDetailContainers, setTeamDetailContainers] = useState([]);
  const [teamDetailContainerCount, setTeamDetailContainerCount] = useState(0);
  const [teamDetailContainerPage, setTeamDetailContainerPage] = useState(1);
  const [teamDetailLoading, setTeamDetailLoading] = useState({ submissions: false, writeups: false, traffic: false });
  const [teamDetailFlags, setTeamDetailFlags] = useState([]);
  const [teamDetailFlagsLoading, setTeamDetailFlagsLoading] = useState(false);
  const teamDetailPageSize = 20;

  useEffect(() => {
    fetchCheats();
  }, [id, currentPage, filterType, filterReasonType]);

  const fetchCheats = async () => {
    try {
      const params = {
        limit: pageSize,
        offset: (currentPage - 1) * pageSize,
      };

      if (filterType) {
        params.type = filterType;
      }

      if (filterReasonType) {
        params.reason_type = filterReasonType;
      }

      const response = await getContestCheats(parseInt(id), params);
      if (response.code === 200) {
        setCheats(response.data.cheats || []);
        setChecked(response.data.checked || 0);
        setTotal(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.toast.fetchFailed') });
    }
  };

  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  const handleFilterChange = (type) => {
    setFilterType(type);
    setCurrentPage(1);
  };

  const handleReasonTypeFilterChange = (reasonType) => {
    setFilterReasonType(reasonType);
    setCurrentPage(1);
  };

  const handleEditCheat = (cheat) => {
    setSelectedCheat(cheat);
    setEditForm({
      reason: cheat.reason || '',
      type: cheat.type || '',
      checked: cheat.checked || false,
      comment: cheat.comment || '',
    });
    setShowEditModal(true);
  };

  const handleViewDetail = (cheat) => {
    setSelectedCheat(cheat);
    setShowDetailModal(true);
  };

  const handleSaveEdit = async () => {
    try {
      const updateData = {};

      // 只更新有变化的字段
      if (editForm.reason !== selectedCheat.reason) {
        updateData.reason = editForm.reason;
      }
      if (editForm.type !== selectedCheat.type) {
        updateData.type = editForm.type;
      }
      if (editForm.checked !== selectedCheat.checked) {
        updateData.checked = editForm.checked;
      }
      if (editForm.comment !== selectedCheat.comment) {
        updateData.comment = editForm.comment;
      }

      if (Object.keys(updateData).length === 0) {
        setShowEditModal(false);
        return;
      }

      const response = await updateContestCheat(parseInt(id), selectedCheat.id, updateData);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.cheats.toast.updateSuccess') });
        setShowEditModal(false);
        fetchCheats();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.toast.updateFailed') });
    }
  };

  const handleDeleteCheat = async () => {
    if (!selectedCheat) return;
    try {
      const response = await deleteContestCheat(parseInt(id), selectedCheat.id);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.cheats.toast.deleteSuccess') });
        setShowDeleteConfirm(false);
        setSelectedCheat(null);
        fetchCheats();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.toast.deleteFailed') });
    }
  };

  const handleDeleteAllCheats = async () => {
    try {
      const response = await deleteAllContestCheats(parseInt(id));
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.cheats.toast.deleteAllSuccess') });
        setShowDeleteAllConfirm(false);
        fetchCheats();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.toast.deleteAllFailed') });
    }
  };

  const handleCheckCheats = async () => {
    setCheckLoading(true);
    try {
      const response = await checkContestCheats(parseInt(id));
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.cheats.toast.checkSuccess') });
        setShowCheckConfirm(false);
        fetchCheats();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.toast.checkFailed') });
    } finally {
      setCheckLoading(false);
    }
  };

  const formatTime = (timestamp) => {
    if (!timestamp) return '-';
    return new Date(timestamp).toLocaleString(i18n.language || 'en-US');
  };

  // === User detail handlers ===
  const handleUserClick = async (userId) => {
    if (!userId) return;
    try {
      const response = await getUserInfo(userId);
      if (response.code === 200) {
        setUserDetailData(response.data);
        setShowUserDetail(true);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.users.toast.fetchFailed') });
    }
  };

  // === Team detail handlers ===
  const handleTeamClick = async (teamId) => {
    if (!teamId) return;
    try {
      const response = await getContestTeam(parseInt(id), teamId);
      if (response.code === 200) {
        setTeamDetailData(response.data);
        setTeamDetailTab('info');
        setShowTeamDetail(true);
        setTeamDetailMembersLoading(true);
        try {
          const membersRes = await getTeamMembers(parseInt(id), teamId);
          if (membersRes.code === 200) {
            setTeamDetailMembers(membersRes.data || []);
          }
        } finally {
          setTeamDetailMembersLoading(false);
        }
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.cheats.toast.fetchFailed') });
    }
  };

  const handleTeamDetailClose = () => {
    setShowTeamDetail(false);
    setTeamDetailData(null);
    setTeamDetailTab('info');
    setTeamDetailMembers([]);
    setTeamDetailSubmissions([]);
    setTeamDetailSubmissionCount(0);
    setTeamDetailSubmissionPage(1);
    setTeamDetailWriteups([]);
    setTeamDetailWriteupCount(0);
    setTeamDetailWriteupPage(1);
    setTeamDetailContainers([]);
    setTeamDetailContainerCount(0);
    setTeamDetailContainerPage(1);
    setTeamDetailLoading({ submissions: false, writeups: false, traffic: false });
    setTeamDetailFlags([]);
    setTeamDetailFlagsLoading(false);
  };

  const fetchTeamDetailSubmissions = async (teamId, page = 1) => {
    setTeamDetailLoading((prev) => ({ ...prev, submissions: true }));
    try {
      const response = await getContestTeamSubmissions(parseInt(id), teamId, {
        limit: teamDetailPageSize,
        offset: (page - 1) * teamDetailPageSize,
      });
      if (response.code === 200) {
        setTeamDetailSubmissions(response.data.submissions || []);
        setTeamDetailSubmissionCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message });
    } finally {
      setTeamDetailLoading((prev) => ({ ...prev, submissions: false }));
    }
  };

  const fetchTeamDetailWriteups = async (teamId, page = 1) => {
    setTeamDetailLoading((prev) => ({ ...prev, writeups: true }));
    try {
      const response = await getContestTeamWriteups(parseInt(id), teamId, {
        limit: teamDetailPageSize,
        offset: (page - 1) * teamDetailPageSize,
      });
      if (response.code === 200) {
        setTeamDetailWriteups(response.data.writeups || []);
        setTeamDetailWriteupCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message });
    } finally {
      setTeamDetailLoading((prev) => ({ ...prev, writeups: false }));
    }
  };

  const fetchTeamDetailContainers = async (teamId, page = 1) => {
    setTeamDetailLoading((prev) => ({ ...prev, traffic: true }));
    try {
      const response = await getTeamContainers(parseInt(id), teamId, {
        limit: teamDetailPageSize,
        offset: (page - 1) * teamDetailPageSize,
      });
      if (response.code === 200) {
        setTeamDetailContainers(response.data.victims || []);
        setTeamDetailContainerCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message });
    } finally {
      setTeamDetailLoading((prev) => ({ ...prev, traffic: false }));
    }
  };

  const fetchTeamDetailFlags = async (teamId) => {
    setTeamDetailFlagsLoading(true);
    try {
      const response = await getContestTeamFlags(parseInt(id), teamId);
      if (response.code === 200) {
        setTeamDetailFlags(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message });
    } finally {
      setTeamDetailFlagsLoading(false);
    }
  };

  const handleTeamDetailTabChange = (tab) => {
    setTeamDetailTab(tab);
    if (!teamDetailData) return;
    if (tab === 'submissions') {
      setTeamDetailSubmissionPage(1);
      fetchTeamDetailSubmissions(teamDetailData.id, 1);
    } else if (tab === 'writeups') {
      setTeamDetailWriteupPage(1);
      fetchTeamDetailWriteups(teamDetailData.id, 1);
    } else if (tab === 'containers') {
      setTeamDetailContainerPage(1);
      fetchTeamDetailContainers(teamDetailData.id, 1);
    } else if (tab === 'flags') {
      fetchTeamDetailFlags(teamDetailData.id);
    }
  };

  const handleTeamDetailPageChange = (type, page) => {
    if (!teamDetailData) return;
    if (type === 'submissions') {
      setTeamDetailSubmissionPage(page);
      fetchTeamDetailSubmissions(teamDetailData.id, page);
    } else if (type === 'writeups') {
      setTeamDetailWriteupPage(page);
      fetchTeamDetailWriteups(teamDetailData.id, page);
    } else if (type === 'containers') {
      setTeamDetailContainerPage(page);
      fetchTeamDetailContainers(teamDetailData.id, page);
    }
  };

  const handleTeamDetailDownloadTraffic = async (container) => {
    if (!teamDetailData) return;
    try {
      const response = await downloadContainerTraffic(parseInt(id), teamDetailData.id, container.id);
      downloadBlobResponse(response);
    } catch (error) {
      toast.danger({ description: error.message });
    }
  };

  const handleTeamDetailDownloadWriteup = async (writeup) => {
    if (!teamDetailData) return;
    try {
      const response = await downloadContestTeamWriteup(parseInt(id), teamDetailData.id, writeup.id);
      downloadBlobResponse(response);
    } catch (error) {
      toast.danger({ description: error.message });
    }
  };

  const getCheatTypeLabel = (type) => {
    const typeMap = {
      cheater: t('admin.contests.cheats.types.cheater'),
      suspicious: t('admin.contests.cheats.types.suspicious'),
      pass: t('admin.contests.cheats.types.pass'),
    };
    return typeMap[type] || type;
  };

  const getCheatTypeColor = (type) => {
    const colorMap = {
      cheater: 'error',
      suspicious: 'warning',
      pass: 'default',
    };
    return colorMap[type] || 'default';
  };

  const getReasonTypeLabel = (reasonType) => {
    if (!reasonType) return '-';
    const key = `admin.contests.cheats.reasonTypes.${reasonType}`;
    const label = t(key);
    return label === key ? reasonType : label;
  };

  const columns = [
    { key: 'id', label: t('admin.contests.cheats.columns.id'), width: '6%' },
    { key: 'model', label: t('admin.contests.cheats.columns.model'), width: '18%' },
    { key: 'type', label: t('admin.contests.cheats.columns.type'), width: '10%' },
    { key: 'reason_type', label: t('admin.contests.cheats.columns.reasonType'), width: '10%' },
    { key: 'reason', label: t('admin.contests.cheats.columns.reason'), width: '14%' },
    { key: 'ip', label: t('admin.contests.cheats.columns.ipOrDevice'), width: '10%' },
    { key: 'checked', label: t('admin.contests.cheats.columns.status'), width: '8%' },
    { key: 'time', label: t('admin.contests.cheats.columns.time'), width: '14%' },
    { key: 'actions', label: t('admin.contests.cheats.columns.actions'), width: '10%' },
  ];

  const isValidIp = (str) => {
    if (!str) return false;
    return /^(\d{1,3}\.){3}\d{1,3}(\/\d+)?$/.test(str) || str.includes(':');
  };

  const renderCell = (item, column) => {
    switch (column.key) {
      case 'model':
        if (!item.model || Object.keys(item.model).length === 0) {
          return '-';
        }
        return (
          <div className="flex flex-wrap gap-1">
            {Object.entries(item.model).map(([key, values]) => {
              const ids = Array.isArray(values) ? values : [values];
              if (key === 'User') {
                return ids.map((uid) => (
                  <span
                    key={`${key}-${uid}`}
                    className="text-geek-400 cursor-pointer hover:underline"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleUserClick(uid);
                    }}
                  >
                    {key}-{uid};
                  </span>
                ));
              }
              if (key === 'Team') {
                return ids.map((tid) => (
                  <span
                    key={`${key}-${tid}`}
                    className="text-geek-400 cursor-pointer hover:underline"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleTeamClick(tid);
                    }}
                  >
                    {key}-{tid};
                  </span>
                ));
              }
              return (
                <span key={key} className="text-neutral-300">
                  {key}: {ids.join(', ')}
                </span>
              );
            })}
          </div>
        );
      case 'type':
        return <StatusTag type={getCheatTypeColor(item.type)} text={getCheatTypeLabel(item.type)} />;
      case 'reason_type':
        return <span className="text-neutral-300 font-mono text-xs">{getReasonTypeLabel(item.reason_type)}</span>;
      case 'reason':
        return (
          <span className="max-w-50 truncate block" title={item.reason}>
            {item.reason || '-'}
          </span>
        );
      case 'ip':
        if (item.ip && isValidIp(item.ip)) {
          return <IpAddress ip={item.ip} className="text-xs" />;
        }
        return <span className="font-mono text-xs text-neutral-300">{item.ip || item.magic}</span>;
      case 'checked':
        return (
          <StatusTag
            type={item.checked ? 'success' : 'warning'}
            text={
              item.checked ? t('admin.contests.cheats.status.processed') : t('admin.contests.cheats.status.unprocessed')
            }
          />
        );
      case 'time':
        return formatTime(item.time);
      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                handleViewDetail(item);
              }}
              className="p-1! h-6! w-6!"
            >
              <IconEye size={14} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                handleEditCheat(item);
              }}
              className="p-1! h-6! w-6!"
            >
              <IconEdit size={14} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                setSelectedCheat(item);
                setShowDeleteConfirm(true);
              }}
              className="p-1! h-6! w-6! text-red-400! hover:text-red-300!"
            >
              <IconTrash size={14} />
            </Button>
          </div>
        );
      default:
        return item[column.key] || '-';
    }
  };

  const filterOptions = [
    { value: 'cheater', label: t('admin.contests.cheats.types.cheater') },
    { value: 'suspicious', label: t('admin.contests.cheats.types.suspicious') },
    { value: 'pass', label: t('admin.contests.cheats.types.pass') },
  ];

  const reasonTypeFilterOptions = [
    { value: 'same_device', label: t('admin.contests.cheats.reasonTypes.same_device') },
    { value: 'same_web_ip', label: t('admin.contests.cheats.reasonTypes.same_web_ip') },
    { value: 'same_victim_ip', label: t('admin.contests.cheats.reasonTypes.same_victim_ip') },
    { value: 'wrong_flag', label: t('admin.contests.cheats.reasonTypes.wrong_flag') },
    { value: 'token_magic', label: t('admin.contests.cheats.reasonTypes.token_magic') },
  ];

  return (
    <div className="w-full mx-auto space-y-6">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        {/* 头部 */}
        <div className="flex items-center justify-end gap-2 mb-6">
          <Button variant="ghost" size="sm" onClick={() => setShowCheckConfirm(true)} className="bg-black/30!">
            <IconShieldCheck size={16} />
            {t('admin.contests.cheats.actions.check')}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setShowDeleteAllConfirm(true)}
            className="bg-black/30! text-red-400! hover:text-red-300!"
          >
            <IconTrash size={16} />
            {t('admin.contests.cheats.actions.deleteAll')}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setShowFilter(!showFilter)}
            className={`bg-black/30! ${showFilter ? 'bg-geek-400/20!' : ''}`}
          >
            <IconFilter size={16} />
            {t('admin.contests.cheats.filterButton')}
          </Button>
        </div>

        {/* 筛选器 */}
        {showFilter && (
          <motion.div
            className="border border-neutral-600 rounded-md bg-neutral-900 p-4 mb-6"
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
          >
            <div className="flex items-center gap-4">
              <span className="text-neutral-400 font-mono">{t('admin.contests.cheats.filterLabel')}</span>
              <div className="flex flex-wrap gap-2">
                {filterOptions.map((option) => (
                  <Button
                    key={option.value}
                    variant={filterType === option.value ? 'primary' : 'ghost'}
                    size="sm"
                    onClick={() => handleFilterChange(option.value)}
                    className={filterType === option.value ? '' : 'bg-black/30!'}
                  >
                    {option.label}
                  </Button>
                ))}
              </div>
            </div>
            <div className="flex items-center gap-4 mt-3">
              <span className="text-neutral-400 font-mono">{t('admin.contests.cheats.reasonTypeFilterLabel')}</span>
              <div className="flex flex-wrap gap-2">
                {reasonTypeFilterOptions.map((option) => (
                  <Button
                    key={option.value}
                    variant={filterReasonType === option.value ? 'primary' : 'ghost'}
                    size="sm"
                    onClick={() => handleReasonTypeFilterChange(option.value)}
                    className={filterReasonType === option.value ? '' : 'bg-black/30!'}
                  >
                    {option.label}
                  </Button>
                ))}
              </div>
            </div>
          </motion.div>
        )}

        {/* 统计信息 */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4">
            <h3 className="text-sm font-mono text-neutral-400 mb-2">{t('admin.contests.cheats.stats.total')}</h3>
            <p className="text-2xl text-neutral-50">{total}</p>
          </div>
          <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4">
            <h3 className="text-sm font-mono text-neutral-400 mb-2">{t('admin.contests.cheats.stats.processed')}</h3>
            <p className="text-2xl text-green-400">{checked}</p>
          </div>
          <div className="border border-neutral-600 rounded-md bg-neutral-900 p-4">
            <h3 className="text-sm font-mono text-neutral-400 mb-2">{t('admin.contests.cheats.stats.unprocessed')}</h3>
            <p className="text-2xl text-yellow-400">{total - checked}</p>
          </div>
        </div>

        {/* 作弊事件列表 */}
        <List
          columns={columns}
          data={cheats}
          renderCell={renderCell}
          empty={cheats.length === 0}
          emptyContent={<EmptyState title={t('admin.contests.cheats.empty')} />}
          footer={
            total > pageSize ? (
              <Pagination
                total={Math.ceil(total / pageSize)}
                current={currentPage}
                onChange={handlePageChange}
                showTotal={true}
                totalItems={total}
              />
            ) : null
          }
        />
      </motion.div>

      {/* 编辑模态框 */}
      <Modal
        isOpen={showEditModal}
        onClose={() => setShowEditModal(false)}
        title={t('admin.contests.cheats.modals.editTitle')}
        size="md"
        footer={
          <>
            <ModalButton variant="default" onClick={() => setShowEditModal(false)}>
              {t('common.cancel')}
            </ModalButton>
            <ModalButton variant="primary" onClick={handleSaveEdit}>
              {t('common.save')}
            </ModalButton>
          </>
        }
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.cheats.form.verdict')}
            </label>
            <select
              value={editForm.type}
              onChange={(e) => setEditForm({ ...editForm, type: e.target.value })}
              className="w-full px-3 py-2 bg-neutral-800 border border-neutral-700 rounded text-neutral-300 focus:outline-none focus:ring-2 focus:ring-geek-400"
            >
              {filterOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.cheats.form.reason')}
            </label>
            <textarea
              value={editForm.reason}
              onChange={(e) => setEditForm({ ...editForm, reason: e.target.value })}
              className="w-full px-3 py-2 bg-neutral-800 border border-neutral-700 rounded text-neutral-300 focus:outline-none focus:ring-2 focus:ring-geek-400"
              rows={3}
              placeholder={t('admin.contests.cheats.form.reasonPlaceholder')}
            />
          </div>

          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.cheats.form.status')}
            </label>
            <div className="flex items-center gap-4">
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={editForm.checked}
                  onChange={(e) => setEditForm({ ...editForm, checked: e.target.checked })}
                  className="w-4 h-4 text-geek-400 bg-neutral-800 border-neutral-700 rounded focus:ring-geek-400"
                />
                <span className="text-neutral-300">{t('admin.contests.cheats.status.processed')}</span>
              </label>
            </div>
          </div>

          <div>
            <label className="block text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.cheats.form.comment')}
            </label>
            <textarea
              value={editForm.comment}
              onChange={(e) => setEditForm({ ...editForm, comment: e.target.value })}
              className="w-full px-3 py-2 bg-neutral-800 border border-neutral-700 rounded text-neutral-300 focus:outline-none focus:ring-2 focus:ring-geek-400"
              rows={3}
              placeholder={t('admin.contests.cheats.form.commentPlaceholder')}
            />
          </div>
        </div>
      </Modal>

      {/* 详情模态框 */}
      <Modal
        isOpen={showDetailModal}
        onClose={() => setShowDetailModal(false)}
        title={t('admin.contests.cheats.modals.detailTitle')}
        size="lg"
        footer={
          <ModalButton variant="default" onClick={() => setShowDetailModal(false)}>
            {t('admin.contests.cheats.actions.close')}
          </ModalButton>
        }
      >
        {selectedCheat && (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.detail.eventId')}
                </label>
                <p className="text-neutral-300">{selectedCheat.id}</p>
              </div>
              {selectedCheat.model && Object.keys(selectedCheat.model).length > 0 ? (
                Object.entries(selectedCheat.model).map(([key, values]) => {
                  const ids = Array.isArray(values) ? values : [values];
                  return (
                    <div key={key}>
                      <label className="block text-sm font-mono text-neutral-400 mb-1">{key}</label>
                      {key === 'User' ? (
                        <div className="flex flex-wrap gap-2">
                          {ids.map((uid) => (
                            <p
                              key={uid}
                              className="text-geek-400 cursor-pointer hover:underline w-fit"
                              onClick={() => handleUserClick(uid)}
                            >
                              {uid}
                            </p>
                          ))}
                        </div>
                      ) : key === 'Team' ? (
                        <div className="flex flex-wrap gap-2">
                          {ids.map((tid) => (
                            <p
                              key={tid}
                              className="text-geek-400 cursor-pointer hover:underline w-fit"
                              onClick={() => handleTeamClick(tid)}
                            >
                              {tid}
                            </p>
                          ))}
                        </div>
                      ) : (
                        <p className="text-neutral-300">{ids.join(', ')}</p>
                      )}
                    </div>
                  );
                })
              ) : (
                <div>
                  <label className="block text-sm font-mono text-neutral-400 mb-1">
                    {t('admin.contests.cheats.detail.model')}
                  </label>
                  <p className="text-neutral-300">-</p>
                </div>
              )}
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.detail.ip')}
                </label>
                {selectedCheat.ip && isValidIp(selectedCheat.ip) ? (
                  <IpAddress ip={selectedCheat.ip} />
                ) : (
                  <p className="text-neutral-300">{selectedCheat.ip || '-'}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.detail.time')}
                </label>
                <p className="text-neutral-300">{formatTime(selectedCheat.time)}</p>
              </div>
            </div>

            <div>
              <label className="block text-sm font-mono text-neutral-400 mb-1">
                {t('admin.contests.cheats.detail.type')}
              </label>
              <StatusTag type={getCheatTypeColor(selectedCheat.type)} text={getCheatTypeLabel(selectedCheat.type)} />
            </div>

            <div>
              <label className="block text-sm font-mono text-neutral-400 mb-1">
                {t('admin.contests.cheats.detail.reasonType')}
              </label>
              <p className="text-neutral-300">{getReasonTypeLabel(selectedCheat.reason_type)}</p>
            </div>

            <div>
              <label className="block text-sm font-mono text-neutral-400 mb-1">
                {t('admin.contests.cheats.detail.reason')}
              </label>
              <p className="text-neutral-300 bg-neutral-800 p-3 rounded">{selectedCheat.reason || '-'}</p>
            </div>

            <div>
              <label className="block text-sm font-mono text-neutral-400 mb-1">
                {t('admin.contests.cheats.detail.status')}
              </label>
              <StatusTag
                type={selectedCheat.checked ? 'success' : 'warning'}
                text={
                  selectedCheat.checked
                    ? t('admin.contests.cheats.status.processed')
                    : t('admin.contests.cheats.status.unprocessed')
                }
              />
            </div>

            {selectedCheat.comment && (
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.detail.comment')}
                </label>
                <p className="text-neutral-300 bg-neutral-800 p-3 rounded">{selectedCheat.comment}</p>
              </div>
            )}

            {selectedCheat.magic && (
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.detail.device')}
                </label>
                <p className="text-neutral-300 font-mono text-sm bg-neutral-800 p-3 rounded break-all">
                  {selectedCheat.magic}
                </p>
              </div>
            )}

            {selectedCheat.hash && (
              <div>
                <label className="block text-sm font-mono text-neutral-400 mb-1">
                  {t('admin.contests.cheats.detail.hash')}
                </label>
                <p className="text-neutral-300 font-mono text-sm bg-neutral-800 p-3 rounded break-all">
                  {selectedCheat.hash}
                </p>
              </div>
            )}
          </div>
        )}
      </Modal>

      {/* 删除确认模态框 */}
      <Modal
        isOpen={showDeleteConfirm}
        onClose={() => setShowDeleteConfirm(false)}
        title={t('admin.contests.cheats.actions.confirmDeleteTitle')}
        size="sm"
        footer={
          <>
            <ModalButton variant="default" onClick={() => setShowDeleteConfirm(false)}>
              {t('common.cancel')}
            </ModalButton>
            <ModalButton variant="danger" onClick={handleDeleteCheat}>
              {t('admin.contests.cheats.actions.delete')}
            </ModalButton>
          </>
        }
      >
        <p className="text-neutral-300">{t('admin.contests.cheats.actions.confirmDeletePrompt')}</p>
      </Modal>

      {/* 清空全部确认模态框 */}
      <Modal
        isOpen={showDeleteAllConfirm}
        onClose={() => setShowDeleteAllConfirm(false)}
        title={t('admin.contests.cheats.actions.confirmDeleteAllTitle')}
        size="sm"
        footer={
          <>
            <ModalButton variant="default" onClick={() => setShowDeleteAllConfirm(false)}>
              {t('common.cancel')}
            </ModalButton>
            <ModalButton variant="danger" onClick={handleDeleteAllCheats}>
              {t('admin.contests.cheats.actions.deleteAll')}
            </ModalButton>
          </>
        }
      >
        <p className="text-neutral-300">{t('admin.contests.cheats.actions.confirmDeleteAllPrompt')}</p>
      </Modal>

      {/* 手动检测确认模态框 */}
      <Modal
        isOpen={showCheckConfirm}
        onClose={() => setShowCheckConfirm(false)}
        title={t('admin.contests.cheats.actions.confirmCheckTitle')}
        size="sm"
        footer={
          <>
            <ModalButton variant="default" onClick={() => setShowCheckConfirm(false)}>
              {t('common.cancel')}
            </ModalButton>
            <ModalButton variant="primary" onClick={handleCheckCheats} disabled={checkLoading}>
              {t('admin.contests.cheats.actions.check')}
            </ModalButton>
          </>
        }
      >
        <p className="text-neutral-300">{t('admin.contests.cheats.actions.confirmCheckPrompt')}</p>
      </Modal>

      {/* 用户信息对话框 */}
      <AdminUserDetailDialog
        isOpen={showUserDetail}
        onClose={() => {
          setShowUserDetail(false);
          setUserDetailData(null);
        }}
        user={userDetailData}
      />

      {/* 队伍信息对话框 */}
      <AdminTeamDetailDialog
        isOpen={showTeamDetail}
        onClose={handleTeamDetailClose}
        team={teamDetailData}
        activeTab={teamDetailTab}
        onTabChange={handleTeamDetailTabChange}
        members={teamDetailMembers}
        membersLoading={teamDetailMembersLoading}
        detailSubmissions={teamDetailSubmissions}
        detailSubmissionCount={teamDetailSubmissionCount}
        detailSubmissionPage={teamDetailSubmissionPage}
        detailWriteups={teamDetailWriteups}
        detailWriteupCount={teamDetailWriteupCount}
        detailWriteupPage={teamDetailWriteupPage}
        detailContainers={teamDetailContainers}
        detailContainerCount={teamDetailContainerCount}
        detailContainerPage={teamDetailContainerPage}
        detailLoading={teamDetailLoading}
        onDetailPageChange={handleTeamDetailPageChange}
        onDetailDownloadTraffic={handleTeamDetailDownloadTraffic}
        onDetailDownloadWriteup={handleTeamDetailDownloadWriteup}
        detailFlags={teamDetailFlags}
        detailFlagsLoading={teamDetailFlagsLoading}
      />
    </div>
  );
}

export default AdminContestCheats;
