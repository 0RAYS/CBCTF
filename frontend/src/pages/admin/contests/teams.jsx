import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import AdminTeams from '../../../components/features/Admin/Contests/AdminTeams';
import {
  getContestTeams,
  updateTeamInfo,
  updateTeamPicture,
  deleteTeam,
  kickTeamMember,
  getTeamMembers,
} from '../../../api/admin/contest';
import { useTranslation } from 'react-i18next';
import { useDebounce, useTeamDetailDialog } from '../../../hooks';

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
  const [nameQuery, setNameQuery] = useState('');
  const [descQuery, setDescQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState(null);

  const debouncedName = useDebounce(nameQuery, 300);
  const debouncedDesc = useDebounce(descQuery, 300);

  const isSearchMode = !!(nameQuery.trim() || descQuery.trim()) && !searchError;

  const { openTeamDetail, renderTeamDetailDialog } = useTeamDetailDialog(parseInt(id));

  useEffect(() => {
    let cancelled = false;
    if (!debouncedName.trim() && !debouncedDesc.trim()) {
      setSearchResults([]);
      setSearchError(null);
      return;
    }
    const doSearch = async () => {
      setSearchLoading(true);
      setSearchError(null);
      try {
        const params = { limit: 10, offset: 0 };
        if (debouncedName.trim()) params.name = debouncedName.trim();
        if (debouncedDesc.trim()) params.description = debouncedDesc.trim();
        const response = await getContestTeams(parseInt(id, 10), params);
        if (response.code !== 200) {
          throw new Error(response.msg || t('admin.contests.teams.toast.searchFailed'));
        }
        if (!cancelled) {
          setSearchResults(response.data.teams || []);
        }
      } catch (error) {
        if (!cancelled) {
          setSearchError(error);
          toast.danger({ description: error.message || t('admin.contests.teams.toast.searchFailed') });
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
  }, [debouncedName, debouncedDesc]);

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
  const handleNameChange = (value) => {
    setNameQuery(value);
    setSearchError(null);
  };
  const handleDescChange = (value) => {
    setDescQuery(value);
    setSearchError(null);
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

  const handleRowClick = (team) => openTeamDetail(team);

  return (
    <>
      <input
        type="file"
        ref={fileInputRef}
        className="hidden"
        accept="image/png,image/jpeg,image/jpg,image/gif"
        onChange={handleFileChange}
      />
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
        nameQuery={nameQuery}
        descQuery={descQuery}
        searchResults={searchResults}
        searchLoading={searchLoading}
        onNameChange={handleNameChange}
        onDescChange={handleDescChange}
        searchRef={searchRef}
        isSearchMode={isSearchMode}
        onRowClick={handleRowClick}
        onPictureUpload={handlePictureUpload}
      />
      {renderTeamDetailDialog()}
    </>
  );
}

export default AdminContestTeams;
