import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import { downloadBlobResponse } from '../../../utils/fileDownload';
import AdminTeams from '../../../components/features/Admin/Contests/AdminTeams';
import {
  getContestTeams,
  updateTeamInfo,
  updateTeamPicture,
  deleteTeam,
  kickTeamMember,
  getTeamMembers,
  getContestTeamSubmissions,
  getContestTeamWriteups,
  getTeamContainers,
  downloadContainerTraffic,
  downloadContestTeamWriteup,
  getContestTeamFlags,
} from '../../../api/admin/contest';
import { useTranslation } from 'react-i18next';
import { searchModels } from '../../../api/admin/search.js';
import { useDebounceSearch } from '../../../hooks';

function AdminContestTeams() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [teams, setTeams] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [showModal, setShowModal] = useState(false);
  const [selectedTeam, setSelectedTeam] = useState(null);
  const [modalMode, setModalMode] = useState('edit');
  const [editForm, setEditForm] = useState({
    name: '',
    description: '',
    hidden: false,
    banned: false,
    captcha: '',
    captain_id: '',
  });
  const [selectedUserId, setSelectedUserId] = useState('');
  const [teamMembers, setTeamMembers] = useState([]);
  const pageSize = 20;
  const { t } = useTranslation();

  // 搜索相关状态
  const searchRef = useRef(null);
  const {
    query: searchQuery,
    setQuery: setSearchQuery,
    results: rawSearchResults,
    loading: searchLoading,
    error: searchError,
  } = useDebounceSearch(
    async (name) => {
      if (!name || name.trim() === '') return [];

      try {
        const response = await searchModels({
          model: 'Team',
          'search[name]': name.trim(),
          limit: 10,
          offset: 0,
        });

        if (response.code !== 200) {
          const businessError = new Error(response.msg || t('admin.contests.teams.toast.searchFailed'));
          businessError.__business = true;
          throw businessError;
        }

        const contestId = parseInt(id, 10);
        const results = response.data.models || [];
        return Number.isFinite(contestId) ? results.filter((item) => item.contest_id === contestId) : results;
      } catch (error) {
        if (!error?.__business) {
          toast.danger({ description: error.message || t('admin.contests.teams.toast.searchFailed') });
        }
        throw error;
      }
    },
    { delay: 300, minLength: 1 }
  );

  const searchResults = rawSearchResults || [];
  const isSearchMode = searchQuery.trim().length > 0 && !searchError;

  // Detail dialog state
  const [showDetailDialog, setShowDetailDialog] = useState(false);
  const [detailTeam, setDetailTeam] = useState(null);
  const [detailTab, setDetailTab] = useState('info');
  const [detailMembers, setDetailMembers] = useState([]);
  const [detailMembersLoading, setDetailMembersLoading] = useState(false);

  const [detailSubmissions, setDetailSubmissions] = useState([]);
  const [detailSubmissionCount, setDetailSubmissionCount] = useState(0);
  const [detailSubmissionPage, setDetailSubmissionPage] = useState(1);

  const [detailWriteups, setDetailWriteups] = useState([]);
  const [detailWriteupCount, setDetailWriteupCount] = useState(0);
  const [detailWriteupPage, setDetailWriteupPage] = useState(1);

  const [detailContainers, setDetailContainers] = useState([]);
  const [detailContainerCount, setDetailContainerCount] = useState(0);
  const [detailContainerPage, setDetailContainerPage] = useState(1);

  const [detailLoading, setDetailLoading] = useState({
    submissions: false,
    writeups: false,
    traffic: false,
  });

  const [detailFlags, setDetailFlags] = useState([]);
  const [detailFlagsLoading, setDetailFlagsLoading] = useState(false);

  const detailPageSize = 20;

  // 头像上传
  const fileInputRef = useRef(null);
  const [pictureUploadTeam, setPictureUploadTeam] = useState(null);

  const handlePictureUpload = (team) => {
    setPictureUploadTeam(team);
    fileInputRef.current?.click();
  };

  const handleFileChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file || !pictureUploadTeam) return;

    try {
      const response = await updateTeamPicture(parseInt(id), pictureUploadTeam.id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.teams.toast.pictureUpdated') });
        fetchTeams();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.pictureUpdateFailed') });
    }
  };

  useEffect(() => {
    if (!isSearchMode) {
      fetchTeams();
    }
  }, [id, currentPage, isSearchMode]);

  const fetchTeams = async () => {
    try {
      const response = await getContestTeams(parseInt(id), { limit: pageSize, offset: (currentPage - 1) * pageSize });
      if (response.code === 200) {
        setTeams(response.data.teams);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchFailed') });
    }
  };

  const fetchTeamMembers = async (team) => {
    try {
      const response = await getTeamMembers(parseInt(id), team.id);
      if (response.code === 200) {
        setTeamMembers(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchMembersFailed') });
    }
  };

  // 处理搜索输入变化
  const handleSearchChange = (value) => {
    setSearchQuery(value);
  };

  const handleEditTeam = (team) => {
    setSelectedTeam(team);
    setSelectedUserId('');
    setEditForm({
      name: team.name,
      description: team.description || '',
      hidden: team.hidden,
      banned: team.banned,
      captcha: team.captcha || '',
      captain_id: team.captain_id || '',
    });
    setModalMode('edit');
    fetchTeamMembers(team);
    setShowModal(true);
  };

  const handleDeleteTeam = (team) => {
    setSelectedTeam(team);
    setModalMode('delete');
    setShowModal(true);
  };

  const handleKickMember = (team) => {
    setSelectedTeam(team);
    setSelectedUserId('');
    setModalMode('kick');
    fetchTeamMembers(team);
    setShowModal(true);
  };

  const handleViewDetails = (team) => {
    navigate(`/admin/contests/${id}/teams/${team.id}/details`);
  };

  const handleModalClose = () => {
    setShowModal(false);
  };

  const handleFormChange = (form) => {
    setEditForm(form);
  };

  const handleUserSelect = (userId) => {
    setSelectedUserId(userId);
  };

  const handleModalSubmit = async () => {
    try {
      if (modalMode === 'edit') {
        const response = await updateTeamInfo(parseInt(id), selectedTeam.id, editForm);
        if (response.code === 200) {
          toast.success({ description: t('admin.contests.teams.toast.updateSuccess') });
        }
      } else if (modalMode === 'delete') {
        const response = await deleteTeam(parseInt(id), selectedTeam.id);
        if (response.code === 200) {
          toast.success({ description: t('admin.contests.teams.toast.deleteSuccess') });
        }
      } else if (modalMode === 'kick') {
        if (!selectedUserId) {
          toast.warning({ description: t('admin.contests.teams.toast.selectMember') });
          return;
        }
        const response = await kickTeamMember(parseInt(id), selectedTeam.id, parseInt(selectedUserId));
        if (response.code === 200) {
          toast.success({ description: t('admin.contests.teams.toast.kickSuccess') });
        }
      }
      setShowModal(false);
      fetchTeams();
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.actionFailed') });
    }
  };

  // === Detail dialog handlers ===

  const fetchDetailSubmissions = async (teamId, page = 1) => {
    setDetailLoading((prev) => ({ ...prev, submissions: true }));
    try {
      const response = await getContestTeamSubmissions(parseInt(id), teamId, {
        limit: detailPageSize,
        offset: (page - 1) * detailPageSize,
      });
      if (response.code === 200) {
        setDetailSubmissions(response.data.submissions || []);
        setDetailSubmissionCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchFailed') });
    } finally {
      setDetailLoading((prev) => ({ ...prev, submissions: false }));
    }
  };

  const fetchDetailWriteups = async (teamId, page = 1) => {
    setDetailLoading((prev) => ({ ...prev, writeups: true }));
    try {
      const response = await getContestTeamWriteups(parseInt(id), teamId, {
        limit: detailPageSize,
        offset: (page - 1) * detailPageSize,
      });
      if (response.code === 200) {
        setDetailWriteups(response.data.writeups || []);
        setDetailWriteupCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchFailed') });
    } finally {
      setDetailLoading((prev) => ({ ...prev, writeups: false }));
    }
  };

  const fetchDetailContainers = async (teamId, page = 1) => {
    setDetailLoading((prev) => ({ ...prev, traffic: true }));
    try {
      const response = await getTeamContainers(parseInt(id), teamId, {
        limit: detailPageSize,
        offset: (page - 1) * detailPageSize,
      });
      if (response.code === 200) {
        setDetailContainers(response.data.victims || []);
        setDetailContainerCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchFailed') });
    } finally {
      setDetailLoading((prev) => ({ ...prev, traffic: false }));
    }
  };

  const fetchDetailFlags = async (teamId) => {
    setDetailFlagsLoading(true);
    try {
      const response = await getContestTeamFlags(parseInt(id), teamId);
      if (response.code === 200) {
        setDetailFlags(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchFailed') });
    } finally {
      setDetailFlagsLoading(false);
    }
  };

  const handleRowClick = async (team) => {
    setDetailTeam(team);
    setDetailTab('info');
    setShowDetailDialog(true);
    setDetailMembersLoading(true);
    try {
      const response = await getTeamMembers(parseInt(id), team.id);
      if (response.code === 200) {
        setDetailMembers(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchMembersFailed') });
    } finally {
      setDetailMembersLoading(false);
    }
  };

  const handleDetailClose = () => {
    setShowDetailDialog(false);
    setDetailTeam(null);
    setDetailTab('info');
    setDetailMembers([]);
    setDetailSubmissions([]);
    setDetailSubmissionCount(0);
    setDetailSubmissionPage(1);
    setDetailWriteups([]);
    setDetailWriteupCount(0);
    setDetailWriteupPage(1);
    setDetailContainers([]);
    setDetailContainerCount(0);
    setDetailContainerPage(1);
    setDetailLoading({ submissions: false, writeups: false, traffic: false });
    setDetailFlags([]);
    setDetailFlagsLoading(false);
  };

  const handleDetailTabChange = (tab) => {
    setDetailTab(tab);
    if (!detailTeam) return;

    if (tab === 'submissions') {
      setDetailSubmissionPage(1);
      fetchDetailSubmissions(detailTeam.id, 1);
    } else if (tab === 'writeups') {
      setDetailWriteupPage(1);
      fetchDetailWriteups(detailTeam.id, 1);
    } else if (tab === 'containers') {
      setDetailContainerPage(1);
      fetchDetailContainers(detailTeam.id, 1);
    } else if (tab === 'flags') {
      fetchDetailFlags(detailTeam.id);
    }
  };

  const handleDetailPageChange = (type, page) => {
    if (!detailTeam) return;
    if (type === 'submissions') {
      setDetailSubmissionPage(page);
      fetchDetailSubmissions(detailTeam.id, page);
    } else if (type === 'writeups') {
      setDetailWriteupPage(page);
      fetchDetailWriteups(detailTeam.id, page);
    } else if (type === 'containers') {
      setDetailContainerPage(page);
      fetchDetailContainers(detailTeam.id, page);
    }
  };

  const handleDetailDownloadTraffic = async (container) => {
    if (!detailTeam) return;
    try {
      const response = await downloadContainerTraffic(parseInt(id), detailTeam.id, container.id);
      downloadBlobResponse(response, `traffic_${container.id}.zip`);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.actionFailed') });
    }
  };

  const handleDetailDownloadWriteup = async (writeup) => {
    if (!detailTeam) return;
    try {
      const response = await downloadContestTeamWriteup(parseInt(id), detailTeam.id, writeup.id);
      downloadBlobResponse(response, writeup.filename);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.actionFailed') });
    }
  };

  return (
    <>
      <input type="file" ref={fileInputRef} className="hidden" accept="image/*" onChange={handleFileChange} />
      <AdminTeams
        teams={teams}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        teamMembers={teamMembers}
        editForm={editForm}
        selectedUserId={selectedUserId}
        showModal={showModal}
        modalMode={modalMode}
        selectedTeam={selectedTeam}
        onPageChange={setCurrentPage}
        onEditTeam={handleEditTeam}
        onDeleteTeam={handleDeleteTeam}
        onKickMember={handleKickMember}
        onViewDetails={handleViewDetails}
        onModalClose={handleModalClose}
        onModalSubmit={handleModalSubmit}
        onFormChange={handleFormChange}
        onUserSelect={handleUserSelect}
        searchQuery={searchQuery}
        searchResults={searchResults}
        searchLoading={searchLoading}
        onSearchChange={handleSearchChange}
        searchRef={searchRef}
        isSearchMode={isSearchMode}
        onRowClick={handleRowClick}
        showDetailDialog={showDetailDialog}
        detailTeam={detailTeam}
        detailTab={detailTab}
        detailMembers={detailMembers}
        detailMembersLoading={detailMembersLoading}
        detailSubmissions={detailSubmissions}
        detailSubmissionCount={detailSubmissionCount}
        detailSubmissionPage={detailSubmissionPage}
        detailWriteups={detailWriteups}
        detailWriteupCount={detailWriteupCount}
        detailWriteupPage={detailWriteupPage}
        detailContainers={detailContainers}
        detailContainerCount={detailContainerCount}
        detailContainerPage={detailContainerPage}
        detailLoading={detailLoading}
        onDetailClose={handleDetailClose}
        onDetailTabChange={handleDetailTabChange}
        onDetailPageChange={handleDetailPageChange}
        onDetailDownloadTraffic={handleDetailDownloadTraffic}
        onDetailDownloadWriteup={handleDetailDownloadWriteup}
        onPictureUpload={handlePictureUpload}
        detailFlags={detailFlags}
        detailFlagsLoading={detailFlagsLoading}
      />
    </>
  );
}

export default AdminContestTeams;
