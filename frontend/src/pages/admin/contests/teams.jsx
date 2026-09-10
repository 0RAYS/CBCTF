import { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import AdminTeams from '../../../components/features/Admin/Contests/teams/Teams';
import {
  getContestTeams,
  updateTeamInfo,
  updateTeamPicture,
  deleteTeam,
  kickTeamMember,
  getTeamMembers,
} from '../../../api/admin/contest';
import { useTranslation } from 'react-i18next';
import { useDebounce } from '../../../hooks/useDebounce';
import { useTeamDetailDialog } from '../../../components/features/Admin/details/useTeamDetailDialog';

function ContestTeamsManagement({ id }) {
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
  const [kickUserId, setKickUserId] = useState('');
  const [teamMembers, setTeamMembers] = useState([]);
  const pageSize = 20;
  const { t } = useTranslation();
  const listRequest = useRef(0);
  const memberRequest = useRef(0);
  const modalSession = useRef(0);
  const [revision, setRevision] = useState(0);
  const refreshTeams = () => setRevision((previous) => previous + 1);
  useEffect(
    () => () => {
      listRequest.current += 1;
      memberRequest.current += 1;
      modalSession.current += 1;
    },
    []
  );

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
      setSearchLoading(false);
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
  }, [id, debouncedName, debouncedDesc, revision]);

  const fetchTeams = async () => {
    const version = ++listRequest.current;
    try {
      const response = await getContestTeams(parseInt(id), { limit: pageSize, offset: (currentPage - 1) * pageSize });
      if (version === listRequest.current && response.code === 200) {
        setTeams(response.data.teams);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      if (version === listRequest.current)
        toast.danger({ description: error.message || t('admin.contests.teams.toast.fetchFailed') });
    }
  };

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
        refreshTeams();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teams.toast.pictureUpdateFailed') });
    }
  };

  useEffect(() => {
    if (!isSearchMode) {
      fetchTeams();
    }
    return () => {
      listRequest.current += 1;
    };
  }, [id, currentPage, isSearchMode, revision]);

  const fetchTeamMembers = async (team) => {
    const version = ++memberRequest.current;
    try {
      const response = await getTeamMembers(parseInt(id), team.id);
      if (version === memberRequest.current && response.code === 200) {
        setTeamMembers(response.data || []);
      }
    } catch (error) {
      if (version === memberRequest.current)
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
    modalSession.current += 1;
    setTeamMembers([]);
    setSelectedTeam(team);
    setSelectedUserId(team.captain_id || '');
    setKickUserId('');
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
    modalSession.current += 1;
    memberRequest.current += 1;
    setSelectedTeam(team);
    setModalMode('delete');
    setShowModal(true);
  };

  const handleKickSubmit = async () => {
    const session = modalSession.current;
    if (!kickUserId) {
      toast.warning({ description: t('admin.contests.teams.toast.selectMember') });
      return;
    }
    try {
      const response = await kickTeamMember(parseInt(id), selectedTeam.id, parseInt(kickUserId));
      if (session !== modalSession.current) return;
      if (response.code === 200) {
        toast.success({ description: t('admin.contests.teams.toast.kickSuccess') });
        setKickUserId('');
        fetchTeamMembers(selectedTeam);
        refreshTeams();
      }
    } catch (error) {
      if (session === modalSession.current)
        toast.danger({ description: error.message || t('admin.contests.teams.toast.actionFailed') });
    }
  };

  const handleModalClose = () => {
    modalSession.current += 1;
    memberRequest.current += 1;
    setShowModal(false);
  };

  const handleFormChange = (form) => {
    setEditForm(form);
  };

  const handleUserSelect = (userId) => {
    setSelectedUserId(userId);
    setEditForm((previous) => ({ ...previous, captain_id: Number(userId) }));
  };

  const handleModalSubmit = async () => {
    const session = modalSession.current;
    try {
      if (modalMode === 'edit') {
        const response = await updateTeamInfo(parseInt(id), selectedTeam.id, editForm);
        if (session !== modalSession.current) return;
        if (response.code !== 200) throw new Error(response.msg || t('admin.contests.teams.toast.actionFailed'));
        if (response.code === 200) {
          toast.success({ description: t('admin.contests.teams.toast.updateSuccess') });
        }
      } else if (modalMode === 'delete') {
        const response = await deleteTeam(parseInt(id), selectedTeam.id);
        if (session !== modalSession.current) return;
        if (response.code !== 200) throw new Error(response.msg || t('admin.contests.teams.toast.actionFailed'));
        if (response.code === 200) {
          toast.success({ description: t('admin.contests.teams.toast.deleteSuccess') });
        }
      }
      handleModalClose();
      refreshTeams();
    } catch (error) {
      if (session === modalSession.current)
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
        kickUserId={kickUserId}
        showModal={showModal}
        modalMode={modalMode}
        selectedTeam={selectedTeam}
        onPageChange={setCurrentPage}
        onEditTeam={handleEditTeam}
        onDeleteTeam={handleDeleteTeam}
        onModalClose={handleModalClose}
        onModalSubmit={handleModalSubmit}
        onFormChange={handleFormChange}
        onUserSelect={handleUserSelect}
        onKickUserSelect={setKickUserId}
        onKickSubmit={handleKickSubmit}
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

export default function AdminContestTeams() {
  const { id } = useParams();
  return <ContestTeamsManagement key={id} id={id} />;
}
