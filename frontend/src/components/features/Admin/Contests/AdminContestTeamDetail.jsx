import { useState, useEffect, useMemo } from 'react';
import { motion } from 'motion/react';
import { IconDownload, IconFile, IconFileCheck, IconServer, IconGraph, IconFlag } from '@tabler/icons-react';
import { Button, Pagination, Card, EmptyState, IpAddress } from '../../../../components/common';
import { useTranslation } from 'react-i18next';

/**
 * 队伍详情展示组件，展示最近提交、队伍题解、靶机访问流量和Flag信息
 * @param {Object} props
 * @param {Array} props.recentSubmissions - 最近提交数据
 * @param {Array} props.teamWriteups - 队伍提交的题解
 * @param {Array} props.containerTraffic - 队伍对靶机的访问流量
 * @param {Function} props.onPageChange - 页码变更回调
 * @param {Function} props.onViewTrafficGraph - 查看流量关系图回调
 * @param {Function} props.onDownloadTraffic - 下载流量回调
 * @param {Function} props.onDownloadWriteup - 下载题解回调
 * @param {String} props.activeTab - 当前活动的标签页
 * @param {Function} props.onTabChange - 标签页切换回调
 * @param {Array} props.detailFlags - Flag数据
 * @param {boolean} props.detailFlagsLoading - Flag加载状态
 */
function TeamDetail({
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
  loading = {
    submissions: false,
    writeups: false,
    traffic: false,
  },
  activeTab = 'submissions',
  onTabChange,
  hideTabs = false,
  onUserClick,
  detailFlags = [],
  detailFlagsLoading = false,
}) {
  const { t, i18n } = useTranslation();
  // 每页显示条数
  const pageSize = 20;

  // Flags filter state
  const [flagFilters, setFlagFilters] = useState({ name: '', type: '', category: '', solved: '' });

  // Reset filters when tab changes away from flags
  useEffect(() => {
    if (activeTab !== 'flags') {
      setFlagFilters({ name: '', type: '', category: '', solved: '' });
    }
  }, [activeTab]);

  // Compute unique types and categories from flags data
  const flagFilterOptions = useMemo(() => {
    const types = new Set();
    const categories = new Set();
    detailFlags.forEach((challenge) => {
      if (challenge.type) types.add(challenge.type);
      if (challenge.category) categories.add(challenge.category);
    });
    return {
      types: [...types].sort(),
      categories: [...categories].sort(),
    };
  }, [detailFlags]);

  // Apply client-side filters
  const filteredFlags = useMemo(() => {
    return detailFlags.filter((challenge) => {
      if (flagFilters.name && !challenge.name?.toLowerCase().includes(flagFilters.name.toLowerCase())) {
        return false;
      }
      if (flagFilters.type && challenge.type !== flagFilters.type) {
        return false;
      }
      if (flagFilters.category && challenge.category !== flagFilters.category) {
        return false;
      }
      if (flagFilters.solved !== '') {
        const hasSolved = (challenge.flags || []).some((flag) => flag.solved);
        if (flagFilters.solved === 'true' && !hasSolved) return false;
        if (flagFilters.solved === 'false' && hasSolved) return false;
      }
      return true;
    });
  }, [detailFlags, flagFilters]);

  // 格式化日期
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleString(i18n.language || 'en-US');
  };

  // 格式化文件大小
  const formatFileSize = (sizeInBytes) => {
    if (sizeInBytes < 1024) return sizeInBytes + ' B';
    if (sizeInBytes < 1024 * 1024) return (sizeInBytes / 1024).toFixed(2) + ' KB';
    return (sizeInBytes / (1024 * 1024)).toFixed(2) + ' MB';
  };

  // 格式化持续时间（将纳秒转为可读格式）
  const formatDuration = (durationNs) => {
    // 纳秒转换为秒
    const seconds = durationNs / 1000000000;

    if (seconds < 60) {
      return t('utils.time.units.second', { count: Math.round(seconds) });
    } else if (seconds < 3600) {
      return `${t('utils.time.units.minute', { count: Math.floor(seconds / 60) })}${t('utils.time.units.second', {
        count: Math.round(seconds % 60),
      })}`;
    } else {
      const hours = Math.floor(seconds / 3600);
      const minutes = Math.floor((seconds % 3600) / 60);
      return `${t('utils.time.units.hour', { count: hours })}${minutes > 0 ? t('utils.time.units.minute', { count: minutes }) : ''}`;
    }
  };

  // 计算总页数
  const calcTotalPages = (count) => {
    return Math.ceil(count / pageSize);
  };

  const getContainerStatus = (startTime, duration) => {
    const now = new Date();
    const start = new Date(startTime);
    const durationInMilliseconds = duration / 1000000; // 纳秒转毫秒
    const end = new Date(start.getTime() + durationInMilliseconds);

    if (now < start) return t('admin.contests.teamDetail.traffic.status.upcoming');
    if (now > end) return t('admin.contests.teamDetail.traffic.status.ended');
    return t('admin.contests.teamDetail.traffic.status.running');
  };

  const tabItems = [
    { key: 'flags', label: t('admin.contests.teamDetail.tabs.flags') },
    { key: 'submissions', label: t('admin.contests.teamDetail.tabs.submissions') },
    { key: 'writeups', label: t('admin.contests.teamDetail.tabs.writeups') },
    { key: 'traffic', label: t('admin.contests.teamDetail.tabs.traffic') },
  ];

  return (
    <div className="w-full mx-auto">
      {!hideTabs && (
        <>
          <div className="mb-8" />
          <div className="mb-6 border-b border-neutral-700">
            <div className="flex gap-8">
              {tabItems.map((tab) => (
                <Button
                  key={tab.key}
                  variant="ghost"
                  className={`pb-1 px-2 relative font-mono text-sm ${
                    activeTab === tab.key ? 'text-geek-400' : 'text-neutral-400'
                  }`}
                  onClick={() => onTabChange(tab.key)}
                >
                  {tab.label}
                  {activeTab === tab.key && (
                    <motion.div
                      className="absolute bottom-0 left-0 right-0 h-0.5 bg-geek-400"
                      layoutId="filterTabIndicator"
                    />
                  )}
                </Button>
              ))}
            </div>
          </div>
        </>
      )}

      {/* 最近提交 */}
      {activeTab === 'submissions' && (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-mono text-neutral-50 flex items-center">
              <IconFileCheck size={20} className="mr-2 text-geek-400" />
              {t('admin.contests.teamDetail.sections.submissions')}
            </h2>
          </div>

          {loading.submissions ? (
            <Card variant="default" padding="md" className="flex justify-center items-center h-32">
              <div className="animate-spin w-8 h-8 border-2 border-geek-500 rounded-full border-t-transparent"></div>
            </Card>
          ) : recentSubmissions.length === 0 ? (
            <Card variant="default" padding="md">
              <EmptyState title={t('admin.contests.teamDetail.empty.submissions')} />
            </Card>
          ) : (
            <>
              <Card variant="default" padding="none" className="overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-neutral-800/50">
                      <tr>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.challengeId')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.teamId')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.userId')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.submittedFlag')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.score')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.ip')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.status')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.submissions.columns.time')}
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-neutral-700">
                      {recentSubmissions.map((submission, index) => (
                        <motion.tr
                          key={index}
                          className="hover:bg-neutral-800/30 transition-colors"
                          whileHover={{ backgroundColor: 'rgba(64, 64, 64, 0.3)' }}
                        >
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">{submission.challenge_id}</td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">{submission.team_id}</td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {onUserClick ? (
                              <span
                                className="text-geek-400 hover:text-geek-300 cursor-pointer transition-colors"
                                onClick={() => onUserClick(submission.user_id)}
                              >
                                {submission.user_id}
                              </span>
                            ) : (
                              submission.user_id
                            )}
                          </td>
                          <td
                            className="px-4 py-3 text-sm font-mono text-neutral-200 max-w-xs truncate"
                            title={submission.value}
                          >
                            {submission.value}
                          </td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">{submission.score ?? '-'}</td>
                          <td className="px-4 py-3 text-sm">
                            <IpAddress ip={submission.ip} />
                          </td>
                          <td className="px-4 py-3">
                            <span
                              className={`px-2 py-1 rounded-md text-xs font-mono ${
                                submission.solved ? 'bg-green-400/20 text-green-400' : 'bg-red-400/20 text-red-400'
                              }`}
                            >
                              {submission.solved
                                ? t('admin.contests.teamDetail.submissions.status.correct')
                                : t('admin.contests.teamDetail.submissions.status.incorrect')}
                            </span>
                          </td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {submission.created_at ? formatDate(submission.created_at) : '-'}
                          </td>
                        </motion.tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </Card>

              <div className="mt-6">
                <Pagination
                  total={calcTotalPages(submissionCount)}
                  current={currentSubmissionPage}
                  pageSize={pageSize}
                  onChange={(page) => onPageChange('submissions', page)}
                  showTotal={true}
                  totalItems={submissionCount}
                />
              </div>
            </>
          )}
        </div>
      )}

      {/* Flags */}
      {activeTab === 'flags' && (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-mono text-neutral-50 flex items-center">
              <IconFlag size={20} className="mr-2 text-geek-400" />
              {t('admin.contests.teamDetail.sections.flags')}
            </h2>
          </div>

          {/* Filter bar */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-3 mb-4">
            <div>
              <label className="block text-xs font-mono text-neutral-400 mb-1">
                {t('admin.contests.teams.detail.flags.filterName')}
              </label>
              <input
                type="text"
                value={flagFilters.name}
                onChange={(e) => setFlagFilters((prev) => ({ ...prev, name: e.target.value }))}
                placeholder={t('admin.contests.teams.detail.flags.filterNamePlaceholder')}
                className="w-full h-9 px-3 bg-black/20 border border-neutral-300/30 rounded-md
                  text-sm text-neutral-50 placeholder-neutral-500
                  focus:outline-none focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]
                  transition-all duration-200"
              />
            </div>
            <div>
              <label className="block text-xs font-mono text-neutral-400 mb-1">
                {t('admin.contests.teams.detail.flags.filterType')}
              </label>
              <select
                value={flagFilters.type}
                onChange={(e) => setFlagFilters((prev) => ({ ...prev, type: e.target.value }))}
                className="select-custom select-custom-sm w-full"
              >
                <option value="">{t('admin.contests.teams.detail.flags.filterTypePlaceholder')}</option>
                {flagFilterOptions.types.map((type) => (
                  <option key={type} value={type}>
                    {type}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-xs font-mono text-neutral-400 mb-1">
                {t('admin.contests.teams.detail.flags.filterCategory')}
              </label>
              <select
                value={flagFilters.category}
                onChange={(e) => setFlagFilters((prev) => ({ ...prev, category: e.target.value }))}
                className="select-custom select-custom-sm w-full"
              >
                <option value="">{t('admin.contests.teams.detail.flags.filterCategoryPlaceholder')}</option>
                {flagFilterOptions.categories.map((cat) => (
                  <option key={cat} value={cat}>
                    {cat}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-xs font-mono text-neutral-400 mb-1">
                {t('admin.contests.teams.detail.flags.filterSolved')}
              </label>
              <select
                value={flagFilters.solved}
                onChange={(e) => setFlagFilters((prev) => ({ ...prev, solved: e.target.value }))}
                className="select-custom select-custom-sm w-full"
              >
                <option value="">{t('admin.contests.teams.detail.flags.filterSolvedPlaceholder')}</option>
                <option value="true">{t('admin.contests.teams.detail.flags.filterSolvedYes')}</option>
                <option value="false">{t('admin.contests.teams.detail.flags.filterSolvedNo')}</option>
              </select>
            </div>
          </div>

          {detailFlagsLoading ? (
            <Card variant="default" padding="md" className="flex justify-center items-center h-32">
              <div className="animate-spin w-8 h-8 border-2 border-geek-500 rounded-full border-t-transparent"></div>
            </Card>
          ) : filteredFlags.length === 0 ? (
            <Card variant="default" padding="md">
              <EmptyState title={t('admin.contests.teamDetail.empty.flags')} />
            </Card>
          ) : (
            <div className="space-y-4">
              {filteredFlags.map((challenge) => (
                <Card key={challenge.id} variant="default" padding="none" className="overflow-hidden">
                  {/* Challenge header */}
                  <div className="px-4 py-3 bg-neutral-800/50 flex items-center gap-2 flex-wrap">
                    <span className="text-sm font-mono text-neutral-50 font-medium">{challenge.name}</span>
                    {challenge.category && (
                      <span className="px-2 py-0.5 rounded text-xs font-mono bg-geek-500/20 text-geek-400">
                        {challenge.category}
                      </span>
                    )}
                    {challenge.type && (
                      <span className="px-2 py-0.5 rounded text-xs font-mono bg-neutral-400/20 text-neutral-400">
                        {challenge.type}
                      </span>
                    )}
                    {challenge.hidden && (
                      <span className="px-2 py-0.5 rounded text-xs font-mono bg-yellow-400/20 text-yellow-400">
                        {t('admin.contests.teams.detail.flags.hidden')}
                      </span>
                    )}
                  </div>
                  {/* Flags table */}
                  <div className="overflow-x-auto">
                    <table className="w-full table-fixed">
                      <thead className="bg-neutral-800/30">
                        <tr>
                          <th className="w-[22%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.value')}
                          </th>
                          <th className="w-[22%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.template')}
                          </th>
                          <th className="w-[9%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.currentScore')}
                          </th>
                          <th className="w-[9%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.initScore')}
                          </th>
                          <th className="w-[9%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.minScore')}
                          </th>
                          <th className="w-[9%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.decay')}
                          </th>
                          <th className="w-[9%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.solvers')}
                          </th>
                          <th className="w-[11%] px-4 py-2 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                            {t('admin.contests.teams.detail.flags.solved')}
                          </th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-neutral-700">
                        {(challenge.flags || []).map((flag, idx) => (
                          <tr key={idx} className="hover:bg-neutral-800/30 transition-colors">
                            <td
                              className="px-4 py-2 text-sm font-mono text-neutral-50 break-all truncate"
                              title={flag.value || '-'}
                            >
                              {flag.value || '-'}
                            </td>
                            <td
                              className="px-4 py-2 text-sm font-mono text-neutral-400 break-all truncate"
                              title={flag.template || '-'}
                            >
                              {flag.template || '-'}
                            </td>
                            <td className="px-4 py-2 text-sm font-mono text-neutral-50">{flag.current_score ?? '-'}</td>
                            <td className="px-4 py-2 text-sm font-mono text-neutral-400">{flag.init_score ?? '-'}</td>
                            <td className="px-4 py-2 text-sm font-mono text-neutral-400">{flag.min_score ?? '-'}</td>
                            <td className="px-4 py-2 text-sm font-mono text-neutral-400">{flag.decay ?? '-'}</td>
                            <td className="px-4 py-2 text-sm font-mono text-neutral-400">{flag.solvers ?? '-'}</td>
                            <td className="px-4 py-2">
                              {flag.solved ? (
                                <span className="px-2 py-0.5 rounded text-xs font-mono bg-green-400/20 text-green-400">
                                  {t('admin.contests.teams.detail.flags.solvedYes')}
                                </span>
                              ) : (
                                <span className="px-2 py-0.5 rounded text-xs font-mono bg-neutral-400/20 text-neutral-400">
                                  {t('admin.contests.teams.detail.flags.solvedNo')}
                                </span>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </Card>
              ))}
            </div>
          )}
        </div>
      )}

      {/* 题解列表 */}
      {activeTab === 'writeups' && (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-mono text-neutral-50 flex items-center">
              <IconFile size={20} className="mr-2 text-geek-400" />
              {t('admin.contests.teamDetail.sections.writeups')}
            </h2>
          </div>

          {loading.writeups ? (
            <Card variant="default" padding="md" className="flex justify-center items-center h-32">
              <div className="animate-spin w-8 h-8 border-2 border-geek-500 rounded-full border-t-transparent"></div>
            </Card>
          ) : teamWriteups.length === 0 ? (
            <Card variant="default" padding="md">
              <EmptyState title={t('admin.contests.teamDetail.empty.writeups')} />
            </Card>
          ) : (
            <>
              <Card variant="default" padding="none" className="overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-neutral-800/50">
                      <tr>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.writeups.columns.submittedAt')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.writeups.columns.filename')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.writeups.columns.size')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.writeups.columns.hash')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.writeups.columns.uploader')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.writeups.columns.actions')}
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-neutral-700">
                      {teamWriteups.map((writeup) => (
                        <motion.tr
                          key={writeup.id}
                          className="hover:bg-neutral-800/30 transition-colors"
                          whileHover={{ backgroundColor: 'rgba(64, 64, 64, 0.3)' }}
                        >
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">{formatDate(writeup.date)}</td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">{writeup.filename}</td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {formatFileSize(writeup.size)}
                          </td>
                          <td
                            className="px-4 py-3 text-sm font-mono text-neutral-200 max-w-xs truncate"
                            title={writeup.hash}
                          >
                            {writeup.hash}
                          </td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {onUserClick ? (
                              <span
                                className="text-geek-400 hover:text-geek-300 cursor-pointer transition-colors"
                                onClick={() => onUserClick(writeup.user_id)}
                              >
                                {writeup.user_id}
                              </span>
                            ) : (
                              writeup.user_id
                            )}
                          </td>
                          <td className="px-4 py-3">
                            <Button
                              variant="ghost"
                              size="icon"
                              className="!text-geek-400 hover:!text-geek-300"
                              onClick={() => onDownloadWriteup(writeup)}
                              title={t('admin.contests.teamDetail.writeups.actions.download')}
                            >
                              <IconDownload size={18} />
                            </Button>
                          </td>
                        </motion.tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </Card>

              <div className="mt-6">
                <Pagination
                  total={calcTotalPages(writeupCount)}
                  current={currentWriteupPage}
                  pageSize={pageSize}
                  onChange={(page) => onPageChange('writeups', page)}
                  showTotal={true}
                  totalItems={writeupCount}
                />
              </div>
            </>
          )}
        </div>
      )}

      {/* 靶机流量 */}
      {activeTab === 'traffic' && (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-mono text-neutral-50 flex items-center">
              <IconServer size={20} className="mr-2 text-geek-400" />
              {t('admin.contests.teamDetail.sections.traffic')}
            </h2>
          </div>

          {loading.traffic ? (
            <Card variant="default" padding="md" className="flex justify-center items-center h-32">
              <div className="animate-spin w-8 h-8 border-2 border-geek-500 rounded-full border-t-transparent"></div>
            </Card>
          ) : containerTraffic.length === 0 ? (
            <Card variant="default" padding="md">
              <EmptyState title={t('admin.contests.teamDetail.empty.traffic')} />
            </Card>
          ) : (
            <>
              <Card variant="default" padding="none" className="overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-neutral-800/50">
                      <tr>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.traffic.columns.containerId')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.traffic.columns.userId')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.traffic.columns.challengeName')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.traffic.columns.startTime')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.traffic.columns.duration')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.traffic.columns.status')}
                        </th>
                        <th className="px-4 py-3 text-left text-xs font-mono text-neutral-300 uppercase tracking-wider">
                          {t('admin.contests.teamDetail.traffic.columns.actions')}
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-neutral-700">
                      {containerTraffic.map((container) => (
                        <motion.tr
                          key={container.id}
                          className="hover:bg-neutral-800/30 transition-colors"
                          whileHover={{ backgroundColor: 'rgba(64, 64, 64, 0.3)' }}
                        >
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">{container.id}</td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {onUserClick ? (
                              <span
                                className="text-geek-400 hover:text-geek-300 cursor-pointer transition-colors"
                                onClick={() => onUserClick(container.user_id)}
                              >
                                {container.user_id}
                              </span>
                            ) : (
                              container.user_id
                            )}
                          </td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {container.contest_challenge_name}
                          </td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {formatDate(container.start)}
                          </td>
                          <td className="px-4 py-3 text-sm font-mono text-neutral-200">
                            {formatDuration(container.duration)}
                          </td>
                          <td className="px-4 py-3">
                            <span
                              className={`px-2 py-1 rounded-md text-xs font-mono ${getContainerStatusClass(container.start, container.duration)}`}
                            >
                              {getContainerStatus(container.start, container.duration)}
                            </span>
                          </td>
                          <td className="px-4 py-3">
                            <div className="flex gap-2">
                              {onViewTrafficGraph && (
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  className="!text-geek-400 hover:!text-geek-300"
                                  onClick={() => onViewTrafficGraph(container)}
                                  title={t('admin.contests.teamDetail.traffic.actions.viewTraffic')}
                                >
                                  <IconGraph size={18} />
                                </Button>
                              )}
                              <Button
                                variant="ghost"
                                size="icon"
                                className="!text-geek-400 hover:!text-geek-300"
                                onClick={() => onDownloadTraffic(container)}
                                title={t('admin.contests.teamDetail.traffic.actions.downloadTraffic')}
                              >
                                <IconDownload size={18} />
                              </Button>
                            </div>
                          </td>
                        </motion.tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </Card>

              <div className="mt-6">
                <Pagination
                  total={calcTotalPages(trafficCount)}
                  current={currentTrafficPage}
                  pageSize={pageSize}
                  onChange={(page) => onPageChange('traffic', page)}
                  showTotal={true}
                  totalItems={trafficCount}
                />
              </div>
            </>
          )}
        </div>
      )}
    </div>
  );
}

// 获取容器状态样式
const getContainerStatusClass = (startTime, duration) => {
  const now = new Date();
  const start = new Date(startTime);
  const durationInMilliseconds = duration / 1000000; // 纳秒转毫秒
  const end = new Date(start.getTime() + durationInMilliseconds);

  if (now < start) return 'bg-geek-400/20 text-geek-400';
  if (now > end) return 'bg-red-400/20 text-red-400';
  return 'bg-green-400/20 text-green-400';
};

export default TeamDetail;
