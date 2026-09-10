import { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import {
  getTeamMembers,
  getContestTeamSubmissions,
  getContestTeamWriteups,
  getTeamContainers,
  downloadContainerTraffic,
  downloadContestTeamWriteup,
  getContestTeamFlags,
} from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import { downloadBlobResponse } from '../../../../utils/fileDownload';
import TeamDetailDialog from './TeamDetailDialog';
import TrafficGraphDialog from '../traffic/TrafficGraphDialog';
import { TEAM_DETAIL_PAGE_SIZE } from '../Contests/teams/teamDetailData.js';

export default function TeamDetailSession({ team, contestId, onClose }) {
  const routes = useSelector((state) => state.user.routes);
  const canViewTraffic = routes.includes('GET /admin/contests/:contestID/teams/:teamID/victims');
  const [activeTab, setActiveTab] = useState('info');
  const [page, setPage] = useState(1);
  const [data, setData] = useState({ items: [], count: 0 });
  const [loading, setLoading] = useState(true);
  const [graph, setGraph] = useState(null);
  useEffect(() => {
    if (activeTab === 'containers' && !canViewTraffic) {
      setActiveTab('info');
      setGraph(null);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setData({ items: [], count: 0 });
    const params = { limit: TEAM_DETAIL_PAGE_SIZE, offset: (page - 1) * TEAM_DETAIL_PAGE_SIZE };
    const requests = {
      info: () => getTeamMembers(contestId, team.id),
      flags: () => getContestTeamFlags(contestId, team.id),
      submissions: () => getContestTeamSubmissions(contestId, team.id, params),
      writeups: () => getContestTeamWriteups(contestId, team.id, params),
      containers: () => getTeamContainers(contestId, team.id, params),
    };
    requests[activeTab]()
      .then((response) => {
        if (cancelled || response.code !== 200) return;
        const key = { submissions: 'submissions', writeups: 'writeups', containers: 'victims' }[activeTab];
        setData({ items: key ? response.data?.[key] || [] : response.data || [], count: response.data?.count || 0 });
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message });
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [contestId, team.id, activeTab, page, canViewTraffic]);
  const changeTab = (tab) => {
    if (tab === activeTab) return;
    if (tab === 'containers' && !canViewTraffic) return;
    setData({ items: [], count: 0 });
    setLoading(true);
    setPage(1);
    setActiveTab(tab);
  };
  const changePage = (type, nextPage) => {
    if (nextPage === page) return;
    if (type !== activeTab || (type === 'containers' && !canViewTraffic)) return;
    setLoading(true);
    setPage(nextPage);
  };
  const downloadTraffic = async (container) => {
    if (!canViewTraffic) return;
    try {
      const response = await downloadContainerTraffic(contestId, team.id, container.id);
      if (response.headers?.file === 'true') downloadBlobResponse(response);
    } catch (error) {
      toast.danger({ description: error.message });
    }
  };
  const downloadWriteup = async (writeup) => {
    try {
      const response = await downloadContestTeamWriteup(contestId, team.id, writeup.id);
      if (response.headers?.file === 'true') downloadBlobResponse(response);
    } catch (error) {
      toast.danger({ description: error.message });
    }
  };
  return (
    <>
      <TeamDetailDialog
        isOpen
        onClose={onClose}
        team={team}
        activeTab={activeTab}
        onTabChange={changeTab}
        members={activeTab === 'info' ? data.items : []}
        membersLoading={loading}
        detailSubmissions={activeTab === 'submissions' ? data.items : []}
        detailSubmissionCount={data.count}
        detailSubmissionPage={page}
        detailWriteups={activeTab === 'writeups' ? data.items : []}
        detailWriteupCount={data.count}
        detailWriteupPage={page}
        detailContainers={activeTab === 'containers' ? data.items : []}
        detailContainerCount={data.count}
        detailContainerPage={page}
        detailFlags={activeTab === 'flags' ? data.items : []}
        detailFlagsLoading={loading}
        detailLoading={{ submissions: loading, writeups: loading, traffic: loading }}
        onDetailPageChange={changePage}
        onDetailDownloadTraffic={downloadTraffic}
        onDetailDownloadWriteup={downloadWriteup}
        onViewTrafficGraph={canViewTraffic ? setGraph : undefined}
        canViewTraffic={canViewTraffic}
      />
      {canViewTraffic && (
        <TrafficGraphDialog
          isOpen={!!graph}
          onClose={() => setGraph(null)}
          container={graph}
          contestId={contestId}
          teamId={team.id}
        />
      )}
    </>
  );
}
