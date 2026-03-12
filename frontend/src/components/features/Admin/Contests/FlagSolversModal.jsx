import { useEffect, useState } from 'react';
import { motion } from 'motion/react';
import { IconX, IconTrophy } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../../components/common';
import { getFlagSolvers } from '../../../../api/admin/challenge';
import { getUserInfo } from '../../../../api/admin/user';
import {
  getContestTeam,
  getTeamMembers,
  getContestTeamSubmissions,
  getContestTeamWriteups,
  getTeamContainers,
  downloadContainerTraffic,
  downloadContestTeamWriteup,
  getContestTeamFlags,
} from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import AdminUserDetailDialog from '../AdminUserDetailDialog';
import AdminTeamDetailDialog from './AdminTeamDetailDialog';
import { downloadBlobResponse } from '../../../../utils/fileDownload';

const teamDetailPageSize = 20;

function FlagSolversModal({ isOpen, onClose, flagIndex, contestId, challengeId, flagId }) {
  const { t } = useTranslation();
  const [solvers, setSolvers] = useState([]);
  const [loading, setLoading] = useState(false);

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

  useEffect(() => {
    if (!isOpen || !flagId) return;
    setLoading(true);
    setSolvers([]);
    getFlagSolvers(contestId, challengeId, flagId)
      .then((res) => {
        if (res.code === 200) {
          setSolvers(res.data?.solvers ?? []);
        }
      })
      .finally(() => setLoading(false));
  }, [isOpen, flagId, contestId, challengeId]);

  const handleUserClick = async (userId) => {
    if (!userId) return;
    try {
      const res = await getUserInfo(userId);
      if (res.code === 200) {
        setUserDetailData(res.data);
        setShowUserDetail(true);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.users.toast.fetchFailed') });
    }
  };

  const handleTeamClick = async (teamId) => {
    if (!teamId) return;
    try {
      const res = await getContestTeam(contestId, teamId);
      if (res.code === 200) {
        setTeamDetailData(res.data);
        setTeamDetailTab('info');
        setShowTeamDetail(true);
        setTeamDetailMembersLoading(true);
        try {
          const membersRes = await getTeamMembers(contestId, teamId);
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
      const res = await getContestTeamSubmissions(contestId, teamId, {
        limit: teamDetailPageSize,
        offset: (page - 1) * teamDetailPageSize,
      });
      if (res.code === 200) {
        setTeamDetailSubmissions(res.data.submissions || []);
        setTeamDetailSubmissionCount(res.data.count || 0);
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
      const res = await getContestTeamWriteups(contestId, teamId, {
        limit: teamDetailPageSize,
        offset: (page - 1) * teamDetailPageSize,
      });
      if (res.code === 200) {
        setTeamDetailWriteups(res.data.writeups || []);
        setTeamDetailWriteupCount(res.data.count || 0);
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
      const res = await getTeamContainers(contestId, teamId, {
        limit: teamDetailPageSize,
        offset: (page - 1) * teamDetailPageSize,
      });
      if (res.code === 200) {
        setTeamDetailContainers(res.data.victims || []);
        setTeamDetailContainerCount(res.data.count || 0);
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
      const res = await getContestTeamFlags(contestId, teamId);
      if (res.code === 200) {
        setTeamDetailFlags(res.data || []);
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

  const handleTeamDetailPageChange = (tab, page) => {
    if (!teamDetailData) return;
    if (tab === 'submissions') {
      setTeamDetailSubmissionPage(page);
      fetchTeamDetailSubmissions(teamDetailData.id, page);
    } else if (tab === 'writeups') {
      setTeamDetailWriteupPage(page);
      fetchTeamDetailWriteups(teamDetailData.id, page);
    } else if (tab === 'containers') {
      setTeamDetailContainerPage(page);
      fetchTeamDetailContainers(teamDetailData.id, page);
    }
  };

  const handleTeamDetailDownloadTraffic = async (containerId) => {
    if (!teamDetailData) return;
    try {
      const res = await downloadContainerTraffic(contestId, teamDetailData.id, containerId);
      downloadBlobResponse(res, `traffic-${containerId}.pcap`);
    } catch (error) {
      toast.danger({ description: error.message });
    }
  };

  const handleTeamDetailDownloadWriteup = async (writeupId) => {
    if (!teamDetailData) return;
    try {
      const res = await downloadContestTeamWriteup(contestId, teamDetailData.id, writeupId);
      downloadBlobResponse(res, `writeup-${writeupId}.pdf`);
    } catch (error) {
      toast.danger({ description: error.message });
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
      <motion.div
        className="w-full max-w-2xl bg-neutral-900 border border-neutral-700 rounded-md overflow-hidden"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.2 }}
      >
        {/* Header */}
        <div className="flex justify-between items-center p-4 border-b border-neutral-700">
          <h2 className="text-lg font-mono text-neutral-50">
            {t('admin.contests.challengeModal.solversModal.title', { index: flagIndex + 1 })}
          </h2>
          <Button variant="ghost" size="icon" className="!text-neutral-400 hover:!text-neutral-300" onClick={onClose}>
            <IconX size={18} />
          </Button>
        </div>

        {/* Body */}
        <div className="p-4 max-h-[60vh] overflow-y-auto">
          {loading ? (
            <p className="text-center text-sm font-mono text-neutral-400 py-8">
              {t('admin.contests.challengeModal.solversModal.loading')}
            </p>
          ) : solvers.length === 0 ? (
            <p className="text-center text-sm font-mono text-neutral-500 py-8">
              {t('admin.contests.challengeModal.solversModal.empty')}
            </p>
          ) : (
            <table className="w-full text-sm font-mono">
              <thead>
                <tr className="border-b border-neutral-700 text-neutral-400">
                  <th className="text-left py-2 pr-4 w-12" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.rank')}
                  </th>
                  <th className="text-left py-2 pr-4" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.user')}
                  </th>
                  <th className="text-left py-2 pr-4" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.team')}
                  </th>
                  <th className="text-right py-2 pr-4 w-20" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.score')}
                  </th>
                  <th className="text-right py-2 w-36" scope="col">
                    {t('admin.contests.challengeModal.solversModal.columns.solvedAt')}
                  </th>
                </tr>
              </thead>
              <tbody>
                {solvers.map((solver, i) => (
                  <tr key={i} className="border-b border-neutral-800 hover:bg-white/5 transition-colors">
                    <td className="py-2 pr-4">
                      {i === 0 ? (
                        <IconTrophy size={16} className="text-yellow-400" />
                      ) : i === 1 ? (
                        <IconTrophy size={16} className="text-neutral-300" />
                      ) : i === 2 ? (
                        <IconTrophy size={16} className="text-amber-600" />
                      ) : (
                        <span className="text-neutral-400">{i + 1}</span>
                      )}
                    </td>
                    <td className="py-2 pr-4">
                      {solver.user_id ? (
                        <button
                          className="text-neutral-200 hover:text-geek-400 transition-colors cursor-pointer text-left"
                          onClick={() => handleUserClick(solver.user_id)}
                        >
                          {solver.user_name || '—'}
                        </button>
                      ) : (
                        <span className="text-neutral-200">{solver.user_name || '—'}</span>
                      )}
                    </td>
                    <td className="py-2 pr-4">
                      {solver.team_id ? (
                        <button
                          className="text-neutral-300 hover:text-geek-400 transition-colors cursor-pointer text-left"
                          onClick={() => handleTeamClick(solver.team_id)}
                        >
                          {solver.team_name || '—'}
                        </button>
                      ) : (
                        <span className="text-neutral-300">{solver.team_name || '—'}</span>
                      )}
                    </td>
                    <td className="py-2 pr-4 text-right text-geek-400">{solver.score}</td>
                    <td className="py-2 text-right text-neutral-400 text-xs">
                      {solver.solved_at ? new Date(solver.solved_at).toLocaleString() : '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </motion.div>

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

export default FlagSolversModal;
