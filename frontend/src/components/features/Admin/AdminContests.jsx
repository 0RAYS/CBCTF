import { motion } from 'motion/react';
import { IconEdit, IconPlus, IconTrash } from '@tabler/icons-react';
import { Button, Pagination } from '../../../components/common';
import { List, StatusTag } from '../../common';
import { useTranslation } from 'react-i18next';

/** 比赛管理展示组件 */
function AdminContests({
  contests = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 6,
  loading = false,
  onPageChange,
  onCreateContest,
  onDeleteContest,
  onContestClick,
  onPictureUpload,
}) {
  const { t, i18n } = useTranslation();
  const locale = i18n.language === 'zh-CN' ? 'zh-CN' : 'en-US';

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

  const formatDuration = (seconds) => {
    const hours = Math.floor(seconds / 3600);
    return t('utils.time.units.hour', { count: hours });
  };

  const columns = [
    { key: 'cover', label: t('admin.contests.table.cover'), width: '10%' },
    { key: 'name', label: t('admin.contests.table.name'), width: '25%' },
    { key: 'status', label: t('admin.contests.table.status'), width: '7%' },
    {
      key: 'schedule',
      label: t('admin.contests.table.schedule'),
      width: '15%',
    },
    { key: 'metrics', label: t('admin.contests.table.metrics'), width: '10%' },
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
            <span>
              {t('admin.contests.metrics.registrations', {
                count: contest.users ?? t('common.notAvailable'),
              })}
            </span>
            <span>{t('admin.contests.metrics.teams', { count: contest.teams ?? t('common.notAvailable') })}</span>
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
          getRowActions={(contest) => [
            {
              key: 'edit',
              label: t('common.edit'),
              icon: <IconEdit size={18} />,
              hidden: !onContestClick,
              onClick: () => onContestClick(contest),
            },
            {
              key: 'delete',
              label: t('common.delete'),
              icon: <IconTrash size={18} />,
              danger: true,
              onClick: () => onDeleteContest?.(contest),
            },
          ]}
          onRowClick={onContestClick}
          loading={loading}
          empty={!loading && contests.length === 0}
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
