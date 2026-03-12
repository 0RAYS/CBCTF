import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  getContestTeam,
  getTeamMembers,
  getContestTeamSubmissions,
  getContestTeamWriteups,
  getTeamContainers,
  downloadContainerTraffic,
  downloadContestTeamWriteup,
  getContestTeamFlags,
} from '../api/admin/contest';
import { toast } from '../utils/toast';
import { downloadBlobResponse } from '../utils/fileDownload';
import AdminTeamDetailDialog from '../components/features/Admin/Contests/AdminTeamDetailDialog';

const PAGE_SIZE = 20;

export function useTeamDetailDialog(contestId) {
  const { t } = useTranslation();

  const [show, setShow] = useState(false);
  const [teamData, setTeamData] = useState(null);
  const [activeTab, setActiveTab] = useState('info');
  const [members, setMembers] = useState([]);
  const [membersLoading, setMembersLoading] = useState(false);

  const [submissions, setSubmissions] = useState([]);
  const [submissionCount, setSubmissionCount] = useState(0);
  const [submissionPage, setSubmissionPage] = useState(1);

  const [writeups, setWriteups] = useState([]);
  const [writeupCount, setWriteupCount] = useState(0);
  const [writeupPage, setWriteupPage] = useState(1);

  const [containers, setContainers] = useState([]);
  const [containerCount, setContainerCount] = useState(0);
  const [containerPage, setContainerPage] = useState(1);

  const [loading, setLoading] = useState({ submissions: false, writeups: false, traffic: false });
  const [flags, setFlags] = useState([]);
  const [flagsLoading, setFlagsLoading] = useState(false);

  const fetchMembers = async (teamId) => {
    setMembersLoading(true);
    try {
      const res = await getTeamMembers(contestId, teamId);
      if (res.code === 200) setMembers(res.data || []);
    } finally {
      setMembersLoading(false);
    }
  };

  const fetchSubmissions = async (teamId, page = 1) => {
    setLoading((prev) => ({ ...prev, submissions: true }));
    try {
      const res = await getContestTeamSubmissions(contestId, teamId, {
        limit: PAGE_SIZE,
        offset: (page - 1) * PAGE_SIZE,
      });
      if (res.code === 200) {
        setSubmissions(res.data.submissions || []);
        setSubmissionCount(res.data.count || 0);
      }
    } catch (err) {
      toast.danger({ description: err.message });
    } finally {
      setLoading((prev) => ({ ...prev, submissions: false }));
    }
  };

  const fetchWriteups = async (teamId, page = 1) => {
    setLoading((prev) => ({ ...prev, writeups: true }));
    try {
      const res = await getContestTeamWriteups(contestId, teamId, {
        limit: PAGE_SIZE,
        offset: (page - 1) * PAGE_SIZE,
      });
      if (res.code === 200) {
        setWriteups(res.data.writeups || []);
        setWriteupCount(res.data.count || 0);
      }
    } catch (err) {
      toast.danger({ description: err.message });
    } finally {
      setLoading((prev) => ({ ...prev, writeups: false }));
    }
  };

  const fetchContainers = async (teamId, page = 1) => {
    setLoading((prev) => ({ ...prev, traffic: true }));
    try {
      const res = await getTeamContainers(contestId, teamId, {
        limit: PAGE_SIZE,
        offset: (page - 1) * PAGE_SIZE,
      });
      if (res.code === 200) {
        setContainers(res.data.victims || []);
        setContainerCount(res.data.count || 0);
      }
    } catch (err) {
      toast.danger({ description: err.message });
    } finally {
      setLoading((prev) => ({ ...prev, traffic: false }));
    }
  };

  const fetchFlags = async (teamId) => {
    setFlagsLoading(true);
    try {
      const res = await getContestTeamFlags(contestId, teamId);
      if (res.code === 200) setFlags(res.data || []);
    } catch (err) {
      toast.danger({ description: err.message });
    } finally {
      setFlagsLoading(false);
    }
  };

  const openTeamDetail = async (teamOrId) => {
    let team;
    if (teamOrId !== null && typeof teamOrId === 'object') {
      team = teamOrId;
    } else {
      if (!teamOrId) return;
      try {
        const res = await getContestTeam(contestId, teamOrId);
        if (res.code !== 200) return;
        team = res.data;
      } catch (err) {
        toast.danger({ description: err.message || t('admin.contests.cheats.toast.fetchFailed') });
        return;
      }
    }
    setTeamData(team);
    setActiveTab('info');
    setShow(true);
    fetchMembers(team.id);
  };

  const handleClose = () => {
    setShow(false);
    setTeamData(null);
    setActiveTab('info');
    setMembers([]);
    setSubmissions([]);
    setSubmissionCount(0);
    setSubmissionPage(1);
    setWriteups([]);
    setWriteupCount(0);
    setWriteupPage(1);
    setContainers([]);
    setContainerCount(0);
    setContainerPage(1);
    setLoading({ submissions: false, writeups: false, traffic: false });
    setFlags([]);
    setFlagsLoading(false);
  };

  const handleTabChange = (tab) => {
    setActiveTab(tab);
    if (!teamData) return;
    if (tab === 'submissions') {
      setSubmissionPage(1);
      fetchSubmissions(teamData.id, 1);
    } else if (tab === 'writeups') {
      setWriteupPage(1);
      fetchWriteups(teamData.id, 1);
    } else if (tab === 'containers') {
      setContainerPage(1);
      fetchContainers(teamData.id, 1);
    } else if (tab === 'flags') {
      fetchFlags(teamData.id);
    }
  };

  const handlePageChange = (type, page) => {
    if (!teamData) return;
    if (type === 'submissions') {
      setSubmissionPage(page);
      fetchSubmissions(teamData.id, page);
    } else if (type === 'writeups') {
      setWriteupPage(page);
      fetchWriteups(teamData.id, page);
    } else if (type === 'containers') {
      setContainerPage(page);
      fetchContainers(teamData.id, page);
    }
  };

  const handleDownloadTraffic = async (container) => {
    if (!teamData) return;
    try {
      const res = await downloadContainerTraffic(contestId, teamData.id, container.id);
      downloadBlobResponse(res);
    } catch (err) {
      toast.danger({ description: err.message });
    }
  };

  const handleDownloadWriteup = async (writeup) => {
    if (!teamData) return;
    try {
      const res = await downloadContestTeamWriteup(contestId, teamData.id, writeup.id);
      downloadBlobResponse(res);
    } catch (err) {
      toast.danger({ description: err.message });
    }
  };

  const renderTeamDetailDialog = () => (
    <AdminTeamDetailDialog
      isOpen={show}
      onClose={handleClose}
      team={teamData}
      activeTab={activeTab}
      onTabChange={handleTabChange}
      members={members}
      membersLoading={membersLoading}
      detailSubmissions={submissions}
      detailSubmissionCount={submissionCount}
      detailSubmissionPage={submissionPage}
      detailWriteups={writeups}
      detailWriteupCount={writeupCount}
      detailWriteupPage={writeupPage}
      detailContainers={containers}
      detailContainerCount={containerCount}
      detailContainerPage={containerPage}
      detailLoading={loading}
      onDetailPageChange={handlePageChange}
      onDetailDownloadTraffic={handleDownloadTraffic}
      onDetailDownloadWriteup={handleDownloadWriteup}
      detailFlags={flags}
      detailFlagsLoading={flagsLoading}
    />
  );

  return { openTeamDetail, renderTeamDetailDialog };
}
