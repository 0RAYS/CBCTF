import { useTranslation } from 'react-i18next';
import { Button } from '../../../../common';
import SubmissionsPanel from './SubmissionsPanel';
import WriteupsPanel from './WriteupsPanel';
import TrafficPanel from './TrafficPanel';
import FlagsPanel from './FlagsPanel';

export default function TeamDetailPanels({
  recentSubmissions = [],
  teamWriteups = [],
  containerTraffic = [],
  submissionCount = 0,
  writeupCount = 0,
  trafficCount = 0,
  currentSubmissionPage = 1,
  currentWriteupPage = 1,
  currentTrafficPage = 1,
  onPageChange,
  onViewTrafficGraph,
  onDownloadTraffic,
  onDownloadWriteup,
  loading = {},
  activeTab = 'submissions',
  onTabChange,
  hideTabs = false,
  onUserClick,
  detailFlags = [],
  detailFlagsLoading = false,
  canViewTraffic = true,
}) {
  const { t } = useTranslation();
  return (
    <div className="w-full mx-auto">
      {!hideTabs && (
        <div className="flex flex-wrap gap-3 border-b border-neutral-700 mb-6">
          {['flags', 'submissions', 'writeups', 'traffic']
            .filter((tab) => tab !== 'traffic' || canViewTraffic)
            .map((tab) => (
              <Button
                key={tab}
                variant="ghost"
                size="sm"
                className={activeTab === tab ? 'text-geek-400' : 'text-neutral-400'}
                onClick={() => onTabChange(tab)}
              >
                {t(`admin.contests.teamDetail.tabs.${tab}`)}
              </Button>
            ))}
        </div>
      )}
      {activeTab === 'flags' && <FlagsPanel flags={detailFlags} loading={detailFlagsLoading} />}
      {activeTab === 'submissions' && (
        <SubmissionsPanel
          submissions={recentSubmissions}
          count={submissionCount}
          page={currentSubmissionPage}
          loading={loading.submissions}
          onPageChange={(page) => onPageChange('submissions', page)}
          onUserClick={onUserClick}
        />
      )}
      {activeTab === 'writeups' && (
        <WriteupsPanel
          writeups={teamWriteups}
          count={writeupCount}
          page={currentWriteupPage}
          loading={loading.writeups}
          onPageChange={(page) => onPageChange('writeups', page)}
          onUserClick={onUserClick}
          onDownload={onDownloadWriteup}
        />
      )}
      {canViewTraffic && activeTab === 'traffic' && (
        <TrafficPanel
          containers={containerTraffic}
          count={trafficCount}
          page={currentTrafficPage}
          loading={loading.traffic}
          onPageChange={(page) => onPageChange('traffic', page)}
          onUserClick={onUserClick}
          onDownload={onDownloadTraffic}
          onViewGraph={onViewTrafficGraph}
        />
      )}
    </div>
  );
}
