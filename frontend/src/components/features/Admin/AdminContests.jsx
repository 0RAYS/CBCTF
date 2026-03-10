import { motion } from 'motion/react';
import { IconPlus, IconTrash } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { List, StatusTag } from '../../common';
import { useTranslation } from 'react-i18next';

/**
 * 比赛管理展示组件
 * @param {Object} props
 * @param {Array} props.contests - 比赛列表数据
 * @param {number} props.totalCount - 总比赛数量
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {function} props.onPageChange - 页码改变回调
 * @param {function} props.onCreateContest - 创建比赛回调
 * @param {function} props.onDeleteContest - 删除比赛回调
 * @param {function} props.onContestClick - 比赛点击回调
 * @param {function} props.onPictureUpload - 上传头像回调
 */
function AdminContests({
  contests = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 6,
  onPageChange,
  onCreateContest,
  onDeleteContest,
  onContestClick,
  onPictureUpload,
}) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language === 'zh-CN' ? 'zh-CN' : 'en-US';

  // 获取比赛状态
  const getContestStatus = (startTime, duration) => {
    const now = new Date().getTime() / 1000;
    const start = new Date(startTime).getTime() / 1000;
    const end = start + duration;

    if (now < start) return 'upcoming';
    if (now > end) return 'ended';
    return 'running';
  };

  const getStatusTag = (status) => {
    if (status === 'running') return { type: 'info', text: t('admin.contests.status.running') };
    if (status === 'upcoming') return { type: 'success', text: t('admin.contests.status.upcoming') };
    return { type: 'default', text: t('admin.contests.status.ended') };
  };

  // 格式化日期时间
  const formatDateTime = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleString(locale, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  // 格式化时长
  const formatDuration = (seconds) => {
    const hours = Math.floor(seconds / 3600);
    return t('admin.contests.units.hours', { count: hours });
  };

  const columns = [
    { key: 'cover', label: t('admin.contests.table.cover'), width: '10%' },
    { key: 'name', label: t('admin.contests.table.name'), width: '25%' },
    { key: 'status', label: t('admin.contests.table.status'), width: '7%' },
    { key: 'schedule', label: t('admin.contests.table.schedule'), width: '15%' },
    { key: 'metrics', label: t('admin.contests.table.metrics'), width: '10%' },
    { key: 'actions', label: t('admin.contests.table.actions'), width: '5%' },
  ];

  const renderCell = (contest, column) => {
    switch (column.key) {
      case 'cover':
        return (
          <div
            className="relative w-24 h-14 rounded-md overflow-hidden border border-neutral-300/20 group cursor-pointer"
            onClick={(e) => {
              e.stopPropagation();
              onPictureUpload?.(contest);
            }}
          >
            <img src={contest.picture} alt={contest.name} loading="lazy" className="w-full h-full object-cover" />
            <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
              <span className="text-neutral-300 text-xs">{t('admin.contests.actions.replaceCover')}</span>
            </div>
          </div>
        );

      case 'name':
        return (
          <div className="flex flex-col gap-1 min-w-0">
            <span className="text-neutral-50 font-mono truncate">{contest.name}</span>
            <span className="text-xs text-neutral-400 line-clamp-2">{contest.description}</span>
          </div>
        );

      case 'status': {
        const status = getContestStatus(contest.start, contest.duration);
        const { type, text } = getStatusTag(status);
        return (
          <div className="flex flex-wrap gap-2">
            <StatusTag type={type} text={text} />
            {contest.hidden && <StatusTag type="warning" text={t('admin.contests.status.hidden')} />}
          </div>
        );
      }

      case 'schedule':
        return (
          <div className="flex flex-col gap-1 text-xs font-mono text-neutral-400">
            <span>
              {t('admin.contests.labels.startTime')}: {formatDateTime(contest.start)}
            </span>
            <span>
              {t('admin.contests.labels.duration')}: {formatDuration(contest.duration)}
            </span>
          </div>
        );

      case 'metrics':
        return (
          <div className="flex flex-col gap-1 text-xs font-mono text-neutral-400">
            <span>{t('admin.contests.metrics.teamSize', { count: contest.size })}</span>
            <span>{t('admin.contests.metrics.registrations', { count: contest.users })}</span>
            <span>{t('admin.contests.metrics.teams', { count: contest.teams })}</span>
          </div>
        );

      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="!bg-transparent !text-red-400"
              onClick={(e) => {
                e.stopPropagation();
                onDeleteContest?.(contest);
              }}
            >
              <IconTrash size={18} />
            </Button>
          </div>
        );

      default:
        return contest[column.key];
    }
  };

  return (
    <div className="w-full mx-auto space-y-6">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <div className="flex justify-end items-center mb-6">
          <Button variant="primary" size="sm" align="icon-left" icon={<IconPlus size={16} />} onClick={onCreateContest}>
            {t('admin.contests.actions.add')}
          </Button>
        </div>

        <List
          columns={columns}
          data={contests}
          renderCell={renderCell}
          onRowClick={onContestClick}
          empty={contests.length === 0}
          emptyContent={t('admin.contests.empty')}
        />

        {totalCount > pageSize && (
          <div className="mt-6">
            <Pagination
              total={Math.ceil(totalCount / pageSize)}
              current={currentPage}
              onChange={onPageChange}
              showTotal
              totalItems={totalCount}
              showJumpTo
            />
          </div>
        )}
      </motion.div>
    </div>
  );
}

export default AdminContests;
