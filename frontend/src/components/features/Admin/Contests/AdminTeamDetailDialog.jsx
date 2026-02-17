import { useState } from 'react';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { Modal, Button, Card, EmptyState, Avatar } from '../../../common';
import AdminContestTeamDetail from './AdminContestTeamDetail';
import AdminUserDetailDialog from '../AdminUserDetailDialog';
import { getUserInfo } from '../../../../api/admin/user';
import { toast } from '../../../../utils/toast';

const TABS = ['info', 'flags', 'submissions', 'writeups', 'containers'];

function AdminTeamDetailDialog({
  isOpen,
  onClose,
  team,
  activeTab = 'info',
  onTabChange,
  // info tab
  members = [],
  membersLoading = false,
  // submissions/writeups/containers tab data
  detailSubmissions = [],
  detailSubmissionCount = 0,
  detailSubmissionPage = 1,
  detailWriteups = [],
  detailWriteupCount = 0,
  detailWriteupPage = 1,
  detailContainers = [],
  detailContainerCount = 0,
  detailContainerPage = 1,
  detailLoading = { submissions: false, writeups: false, traffic: false },
  onDetailPageChange,
  onDetailDownloadTraffic,
  onDetailDownloadWriteup,
  onViewTrafficGraph,
  // flags tab data
  detailFlags = [],
  detailFlagsLoading = false,
}) {
  const { t } = useTranslation();

  // User detail dialog state (self-managed)
  const [showUserDetail, setShowUserDetail] = useState(false);
  const [userDetailData, setUserDetailData] = useState(null);

  if (!team) return null;

  // Map dialog tab key to AdminContestTeamDetail tab key
  const mapTabKey = (tab) => {
    if (tab === 'containers') return 'traffic';
    return tab;
  };

  const handleTeamDetailTabChange = (teamDetailTab) => {
    if (teamDetailTab === 'traffic') {
      onTabChange('containers');
    } else {
      onTabChange(teamDetailTab);
    }
  };

  const handleTeamDetailPageChange = (type, page) => {
    if (type === 'traffic') {
      onDetailPageChange('containers', page);
    } else {
      onDetailPageChange(type, page);
    }
  };

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

  const handleUserDetailClose = () => {
    setShowUserDetail(false);
    setUserDetailData(null);
  };

  const renderInfoTab = () => (
    <div className="space-y-6">
      {/* Team header */}
      <div className="flex items-start gap-4">
        <Avatar src={team.picture} name={team.name} size="lg" className="border border-neutral-300/30" />
        <div className="min-w-0 flex-1">
          <h3 className="text-lg font-mono text-neutral-50 font-medium">{team.name}</h3>
          <p className="text-sm text-neutral-400 mt-1">
            {team.description || t('admin.contests.teams.detail.info.noDescription')}
          </p>
          <div className="flex items-center gap-2 mt-2">
            {team.banned && (
              <span className="px-2 py-1 rounded text-xs font-mono bg-red-400/20 text-red-400">
                {t('admin.contests.teams.status.banned')}
              </span>
            )}
            {team.hidden && (
              <span className="px-2 py-1 rounded text-xs font-mono bg-yellow-400/20 text-yellow-400">
                {t('admin.contests.teams.status.hidden')}
              </span>
            )}
          </div>
        </div>
      </div>

      {/* Info grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card variant="default" padding="md">
          <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.contests.teams.detail.info.score')}</div>
          <div className="text-lg font-mono text-neutral-50">{team.score?.toLocaleString() || 0}</div>
        </Card>
        <Card
          variant="default"
          padding="md"
          onClick={team.captain_id ? () => handleUserClick(team.captain_id) : undefined}
        >
          <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.contests.teams.detail.info.captain')}</div>
          <div className={`text-lg font-mono ${team.captain_id ? 'text-geek-400' : 'text-neutral-50'}`}>
            {team.captain_id || '-'}
          </div>
        </Card>
        <Card variant="default" padding="md">
          <div className="text-xs font-mono text-neutral-400 mb-1">
            {t('admin.contests.teams.detail.info.inviteCode')}
          </div>
          <div className="text-sm font-mono text-neutral-50 break-all">{team.captcha || '-'}</div>
        </Card>
        <Card variant="default" padding="md">
          <div className="text-xs font-mono text-neutral-400 mb-1">{t('admin.contests.teams.detail.info.members')}</div>
          <div className="text-lg font-mono text-neutral-50">{team.users || 0}</div>
        </Card>
      </div>

      {/* Members table */}
      <div>
        <h4 className="text-sm font-mono text-neutral-400 mb-3">{t('admin.contests.teams.detail.info.members')}</h4>
        {membersLoading ? (
          <Card variant="default" padding="md" className="flex justify-center items-center h-24">
            <div className="animate-spin w-6 h-6 border-2 border-geek-500 rounded-full border-t-transparent"></div>
          </Card>
        ) : members.length === 0 ? (
          <Card variant="default" padding="md">
            <EmptyState title={t('admin.contests.teams.detail.empty.members')} />
          </Card>
        ) : (
          <Card variant="default" padding="none" className="overflow-hidden">
            <table className="w-full">
              <thead className="bg-neutral-800/50">
                <tr>
                  <th className="px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                    {t('admin.contests.teams.detail.info.memberId')}
                  </th>
                  <th className="px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                    {t('admin.contests.teams.detail.info.memberName')}
                  </th>
                  <th className="px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                    {t('admin.contests.teams.detail.info.memberRole')}
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-neutral-700">
                {members.map((member) => (
                  <tr
                    key={member.id}
                    className="hover:bg-neutral-800/30 transition-colors cursor-pointer"
                    onClick={() => handleUserClick(member.id)}
                  >
                    <td className="px-4 py-2 text-sm font-mono text-geek-400">{member.id}</td>
                    <td className="px-4 py-2 text-sm font-mono text-geek-400">{member.name}</td>
                    <td className="px-4 py-2">
                      {member.id === team.captain_id ? (
                        <span className="px-2 py-0.5 rounded text-xs font-mono bg-geek-500/20 text-geek-400">
                          {t('admin.contests.teams.detail.info.captainBadge')}
                        </span>
                      ) : (
                        <span className="px-2 py-0.5 rounded text-xs font-mono bg-neutral-400/20 text-neutral-400">
                          {t('admin.contests.teams.detail.info.memberBadge')}
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </Card>
        )}
      </div>
    </div>
  );

  const renderDataTab = () => (
    <AdminContestTeamDetail
      hideTabs
      activeTab={mapTabKey(activeTab)}
      onTabChange={handleTeamDetailTabChange}
      detailFlags={detailFlags}
      detailFlagsLoading={detailFlagsLoading}
      recentSubmissions={detailSubmissions}
      submissionCount={detailSubmissionCount}
      currentSubmissionPage={detailSubmissionPage}
      teamWriteups={detailWriteups}
      writeupCount={detailWriteupCount}
      currentWriteupPage={detailWriteupPage}
      containerTraffic={detailContainers}
      trafficCount={detailContainerCount}
      currentTrafficPage={detailContainerPage}
      loading={detailLoading}
      onPageChange={handleTeamDetailPageChange}
      onDownloadTraffic={onDetailDownloadTraffic}
      onDownloadWriteup={onDetailDownloadWriteup}
      onViewTrafficGraph={onViewTrafficGraph}
      onUserClick={handleUserClick}
    />
  );

  const renderTabContent = () => {
    if (activeTab === 'info') return renderInfoTab();
    return renderDataTab();
  };

  return (
    <>
      <Modal isOpen={isOpen} onClose={onClose} title={t('admin.contests.teams.detail.title')} size="2xl">
        {/* Tab bar */}
        <div className="mb-6 border-b border-neutral-700">
          <div className="flex gap-8">
            {TABS.map((tab) => (
              <Button
                key={tab}
                variant="ghost"
                className={`pb-1 px-2 relative font-mono text-sm ${
                  activeTab === tab ? 'text-geek-400' : 'text-neutral-400'
                }`}
                onClick={() => onTabChange(tab)}
              >
                {t(`admin.contests.teams.detail.tabs.${tab}`)}
                {activeTab === tab && (
                  <motion.div
                    className="absolute bottom-0 left-0 right-0 h-0.5 bg-geek-400"
                    layoutId="teamDetailTabIndicator"
                  />
                )}
              </Button>
            ))}
          </div>
        </div>

        {/* Tab content */}
        {renderTabContent()}
      </Modal>

      <AdminUserDetailDialog isOpen={showUserDetail} onClose={handleUserDetailClose} user={userDetailData} />
    </>
  );
}

export default AdminTeamDetailDialog;
