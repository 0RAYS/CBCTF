import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { toast } from '../../../utils/toast';
import { downloadBlobResponse } from '../../../utils/fileDownload';
import {
  getTeamContainers,
  getContainerTraffic,
  downloadContainerTraffic,
  getContestTeamSubmissions,
  getContestTeamWriteups,
  downloadContestTeamWriteup,
  getContestTeamFlags,
} from '../../../api/admin/contest';
import { hexToUtf8 } from '../../../utils/hex';
import AdminContestTeamDetail from '../../../components/features/Admin/Contests/AdminContestTeamDetail';
import { Modal } from '../../../components/common';
import { useUserDetailDialog } from '../../../hooks';
import TrafficGraphModal from '../../../components/features/Admin/Contests/TrafficGraphModal';
import { motion } from 'motion/react';
import { Pagination, EmptyState } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function LoadingSpinner({ className = '' }) {
  const { t } = useTranslation();
  return (
    <div className={`flex justify-center items-center ${className}`}>
      <div className="relative w-12 h-12">
        <div className="absolute inset-0 border-2 border-t-geek-400 border-r-geek-400 border-b-transparent border-l-transparent rounded-full animate-spin"></div>
        <div className="absolute inset-1 border-2 border-t-transparent border-r-transparent border-b-geek-400 border-l-geek-400 rounded-full animate-spin animate-delay-150"></div>
      </div>
      <span className="ml-3 text-neutral-300 font-mono">{t('common.loading')}</span>
    </div>
  );
}

function TeamDetails() {
  const { id: contestId, teamId } = useParams();
  const routes = useSelector((state) => state.user.routes);
  const canViewTraffic = routes.includes('GET /admin/contests/:contestID/teams/:teamID/victims');

  // 模态框状态
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isGraphModalOpen, setIsGraphModalOpen] = useState(false);

  const [selectedContainer, setSelectedContainer] = useState(null);
  const [trafficData, setTrafficData] = useState([]);
  const [trafficTotalCount, setTrafficTotalCount] = useState(0);
  const [trafficCurrentPage, setTrafficCurrentPage] = useState(1);

  // 当前打开的负载详情索引
  const [expandedPayloadIndex, setExpandedPayloadIndex] = useState(null);

  // 当前活动的标签页
  const [activeTab, setActiveTab] = useState('submissions');

  // AdminContestTeamDetail所需数据
  const [recentSubmissions, setRecentSubmissions] = useState([]);
  const [submissionCount, setSubmissionCount] = useState(0);
  const [currentSubmissionPage, setCurrentSubmissionPage] = useState(1);

  const [teamWriteups, setTeamWriteups] = useState([]);
  const [writeupCount, setWriteupCount] = useState(0);
  const [currentWriteupPage, setCurrentWriteupPage] = useState(1);

  const [containerTraffic, setContainerTraffic] = useState([]);
  const [trafficCount, setTrafficCount] = useState(0);
  const [currentTrafficPage, setCurrentTrafficPage] = useState(1);

  const [loading, setLoading] = useState({
    submissions: false,
    writeups: false,
    traffic: false,
    modalTraffic: false,
  });

  const [detailFlags, setDetailFlags] = useState([]);
  const [detailFlagsLoading, setDetailFlagsLoading] = useState(false);

  const { openUserDetail, renderUserDetailDialog } = useUserDetailDialog();

  const pageSize = 20;
  const trafficLimit = 20;
  const { t, i18n } = useTranslation();

  // 监听标签页切换
  useEffect(() => {
    if (!canViewTraffic && activeTab === 'traffic') {
      setActiveTab('submissions');
      return;
    }
    switch (activeTab) {
      case 'submissions':
        fetchSubmissions();
        break;
      case 'writeups':
        fetchWriteups();
        break;
      case 'traffic':
        fetchContainers();
        break;
      case 'flags':
        fetchFlags();
        break;
      default:
        break;
    }
  }, [activeTab, canViewTraffic]);

  // 监听分页变化，重新获取数据
  useEffect(() => {
    if (activeTab === 'submissions') {
      fetchSubmissions();
    }
  }, [currentSubmissionPage]);

  useEffect(() => {
    if (activeTab === 'writeups') {
      fetchWriteups();
    }
  }, [currentWriteupPage]);

  useEffect(() => {
    if (activeTab === 'traffic') {
      fetchContainers();
    }
  }, [currentTrafficPage]);

  // 获取队伍提交记录
  const fetchSubmissions = async () => {
    setLoading((prev) => ({ ...prev, submissions: true }));
    try {
      const response = await getContestTeamSubmissions(parseInt(contestId), parseInt(teamId), {
        limit: pageSize,
        offset: (currentSubmissionPage - 1) * pageSize,
      });

      if (response.code === 200) {
        setRecentSubmissions(response.data.submissions || []);
        setSubmissionCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.fetchSubmissionsFailed') });
    } finally {
      setLoading((prev) => ({ ...prev, submissions: false }));
    }
  };

  // 获取队伍题解
  const fetchWriteups = async () => {
    setLoading((prev) => ({ ...prev, writeups: true }));
    try {
      // DEBUG
      const response = await getContestTeamWriteups(parseInt(contestId), parseInt(teamId), {
        limit: pageSize,
        offset: (currentWriteupPage - 1) * pageSize,
      });

      if (response.code === 200) {
        setTeamWriteups(response.data.writeups || []);
        setWriteupCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.fetchWriteupsFailed') });
    } finally {
      setLoading((prev) => ({ ...prev, writeups: false }));
    }
  };

  // 获取容器列表
  const fetchContainers = async () => {
    setLoading((prev) => ({ ...prev, traffic: true }));
    try {
      const response = await getTeamContainers(parseInt(contestId), parseInt(teamId), {
        limit: pageSize,
        offset: (currentTrafficPage - 1) * pageSize,
      });

      if (response.code === 200) {
        setContainerTraffic(response.data.victims || []);
        setTrafficCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.fetchContainersFailed') });
    } finally {
      setLoading((prev) => ({ ...prev, traffic: false }));
    }
  };

  // 获取队伍Flag
  const fetchFlags = async () => {
    setDetailFlagsLoading(true);
    try {
      const response = await getContestTeamFlags(parseInt(contestId), parseInt(teamId));
      if (response.code === 200) {
        setDetailFlags(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.fetchContainersFailed') });
    } finally {
      setDetailFlagsLoading(false);
    }
  };

  const handleViewTrafficGraph = (container) => {
    setSelectedContainer(container);
    setIsGraphModalOpen(true);
  };

  const fetchTrafficData = async (containerId) => {
    try {
      // DEBUG
      const trafficResponse = await getContainerTraffic(parseInt(contestId), parseInt(teamId), containerId, {
        limit: trafficLimit,
        offset: (trafficCurrentPage - 1) * trafficLimit,
      });
      if (trafficResponse.code === 200) {
        setTrafficData(trafficResponse.data.traffics || []);
        setTrafficTotalCount(trafficResponse.data.count || 0);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.fetchTrafficFailed') });
    }
  };

  useEffect(() => {
    if (selectedContainer && isModalOpen) {
      fetchTrafficData(selectedContainer.id);
    }
  }, [trafficCurrentPage, selectedContainer, isModalOpen]);

  const handlePageChange = (type, page) => {
    switch (type) {
      case 'submissions':
        setCurrentSubmissionPage(page);
        break;
      case 'writeups':
        setCurrentWriteupPage(page);
        break;
      case 'traffic':
        setCurrentTrafficPage(page);
        break;
      default:
        break;
    }
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
    setExpandedPayloadIndex(null);
  };

  const handleGraphModalClose = () => {
    setIsGraphModalOpen(false);
  };

  const handleTabChange = (tab) => {
    setActiveTab(tab);

    // 切换标签页时重置对应的页码
    switch (tab) {
      case 'submissions':
        setCurrentSubmissionPage(1);
        break;
      case 'writeups':
        setCurrentWriteupPage(1);
        break;
      case 'traffic':
        setCurrentTrafficPage(1);
        break;
      default:
        break;
    }
  };

  const handleDownloadTraffic = async (container) => {
    try {
      const response = await downloadContainerTraffic(parseInt(contestId), parseInt(teamId), container.id);
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, `traffic_${container.id}.zip`);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.downloadTrafficFailed') });
    }
  };

  const handleDownloadWriteup = async (writeup) => {
    try {
      const response = await downloadContestTeamWriteup(parseInt(contestId), parseInt(teamId), writeup.id);
      if (response.headers?.['file'] === 'true') {
        downloadBlobResponse(response, writeup.filename);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.downloadWriteupFailed') });
    }
  };

  // 处理负载数据点击展开/收起
  const togglePayload = (index) => {
    if (expandedPayloadIndex === index) {
      setExpandedPayloadIndex(null);
    } else {
      setExpandedPayloadIndex(index);
    }
  };

  // 渲染负载内容
  const renderPayload = (payload, index) => {
    if (!payload) return <span className="text-neutral-400">{t('common.none')}</span>;

    try {
      const decoded = hexToUtf8(payload);
      const isExpanded = expandedPayloadIndex === index;

      return (
        <div>
          <div
            onClick={() => togglePayload(index)}
            className="cursor-pointer font-mono text-sm hover:text-geek-400 transition-colors"
          >
            <div className="flex items-center">
              <span className={`line-clamp-2 ${isExpanded ? 'text-geek-400' : ''}`}>
                {decoded.substring(0, 100)}
                {decoded.length > 100 ? '...' : ''}
              </span>
              <span className="ml-2 text-xs text-neutral-400">
                {isExpanded
                  ? t('admin.contests.teamContainers.payload.collapse')
                  : t('admin.contests.teamContainers.payload.expand')}
              </span>
            </div>
          </div>

          {isExpanded && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.2 }}
              className="mt-2 border border-neutral-800 rounded-md p-3 bg-black/50 max-h-[300px] overflow-auto"
            >
              <pre className="whitespace-pre-wrap font-mono text-sm">{decoded}</pre>
            </motion.div>
          )}
        </div>
      );
    } catch (error) {
      return (
        <span className="text-red-400">
          {t('admin.contests.teamContainers.payload.parseError', { message: error.message })}
        </span>
      );
    }
  };

  return (
    <>
      <AdminContestTeamDetail
        recentSubmissions={recentSubmissions}
        teamWriteups={teamWriteups}
        containerTraffic={containerTraffic}
        submissionCount={submissionCount}
        writeupCount={writeupCount}
        trafficCount={trafficCount}
        currentSubmissionPage={currentSubmissionPage}
        currentWriteupPage={currentWriteupPage}
        currentTrafficPage={currentTrafficPage}
        onPageChange={handlePageChange}
        onViewTrafficGraph={handleViewTrafficGraph}
        onDownloadTraffic={handleDownloadTraffic}
        onDownloadWriteup={handleDownloadWriteup}
        loading={loading}
        onTabChange={handleTabChange}
        activeTab={activeTab}
        onUserClick={openUserDetail}
        detailFlags={detailFlags}
        detailFlagsLoading={detailFlagsLoading}
        canViewTraffic={canViewTraffic}
      />

      {/* 使用AdminModal组件替代heroui的Modal */}
      <Modal
        isOpen={isModalOpen}
        onClose={handleModalClose}
        title={t('admin.contests.teamContainers.modal.title')}
        size="lg"
      >
        {loading.modalTraffic ? (
          <LoadingSpinner className="h-40" />
        ) : trafficData.length === 0 ? (
          <EmptyState title={t('admin.contests.teamContainers.modal.empty')} />
        ) : (
          <div className="space-y-4 max-h-[500px] overflow-y-auto pr-2">
            {trafficData.map((traffic, index) => (
              <div key={traffic.id || index} className="border border-neutral-700 rounded-md bg-black/30 p-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-neutral-400 font-mono">
                        {t('admin.contests.teamContainers.labels.type')}
                      </span>
                      <span
                        className={`px-2 py-1 rounded-md text-xs font-mono ${
                          traffic.type === 'request'
                            ? 'bg-geek-500/20 text-geek-400'
                            : traffic.type === 'response'
                              ? 'bg-green-400/20 text-green-400'
                              : 'bg-geek-400/20 text-geek-400'
                        }`}
                      >
                        {traffic.type === 'request'
                          ? t('admin.contests.teamContainers.types.request')
                          : traffic.type === 'response'
                            ? t('admin.contests.teamContainers.types.response')
                            : t('admin.contests.teamContainers.types.tcp')}
                      </span>
                    </div>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-neutral-400 font-mono">
                        {t('admin.contests.teamContainers.labels.sourceIp')}
                      </span>
                      <span className="text-white font-mono">{traffic.src_ip}</span>
                    </div>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-neutral-400 font-mono">
                        {t('admin.contests.teamContainers.labels.sourcePort')}
                      </span>
                      <span className="text-white font-mono">{traffic.src_port}</span>
                    </div>
                  </div>
                  <div>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-neutral-400 font-mono">
                        {t('admin.contests.teamContainers.labels.timestamp')}
                      </span>
                      <span className="text-white font-mono">
                        {new Date(traffic.timestamp).toLocaleString(i18n.language || 'en-US')}
                      </span>
                    </div>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-neutral-400 font-mono">
                        {t('admin.contests.teamContainers.labels.destinationIp')}
                      </span>
                      <span className="text-white font-mono">{traffic.dst_ip}</span>
                    </div>
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-neutral-400 font-mono">
                        {t('admin.contests.teamContainers.labels.destinationPort')}
                      </span>
                      <span className="text-white font-mono">{traffic.dst_port}</span>
                    </div>
                  </div>
                </div>
                <div className="mt-4">
                  <span className="text-neutral-400 font-mono">
                    {t('admin.contests.teamContainers.labels.payload')}
                  </span>
                  <div className="mt-2">{renderPayload(traffic.payload, index)}</div>
                </div>
              </div>
            ))}
          </div>
        )}

        {trafficTotalCount > trafficLimit && (
          <div className="flex justify-center mt-4">
            <Pagination
              total={Math.ceil(trafficTotalCount / trafficLimit)}
              current={trafficCurrentPage}
              onChange={setTrafficCurrentPage}
              showTotal={true}
              totalItems={trafficTotalCount}
            />
          </div>
        )}
      </Modal>

      {/* 流量关系图弹窗 */}
      {canViewTraffic && (
        <TrafficGraphModal
          isOpen={isGraphModalOpen}
          onClose={handleGraphModalClose}
          container={selectedContainer}
          contestId={parseInt(contestId)}
          teamId={parseInt(teamId)}
        />
      )}

      {renderUserDetailDialog()}
    </>
  );
}

export default TeamDetails;
