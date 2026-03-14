import { useEffect, useState } from 'react';
import { IconEdit } from '@tabler/icons-react';
import { Button, Input, List, Modal, Pagination } from '../../components/common';
import { getCronJobList, updateCronJob } from '../../api/admin/cronjob';
import { toast } from '../../utils/toast';
import { useTranslation } from 'react-i18next';

const PAGE_SIZE = 20;
const SECOND_NS = 1_000_000_000;
const MINUTE_NS = 60 * SECOND_NS;
const HOUR_NS = 60 * MINUTE_NS;

const DEFAULT_FORM = {
  hours: '0',
  minutes: '0',
  seconds: '0',
};

function parseDurationNs(value) {
  if (typeof value === 'number' && Number.isFinite(value)) return value;
  if (typeof value !== 'string') return 0;

  const hourMatch = value.match(/(\d+)h/);
  const minuteMatch = value.match(/(\d+)m/);
  const secondMatch = value.match(/(\d+)s/);
  const msMatch = value.match(/(\d+)ms/);

  let total = 0;
  if (hourMatch) total += Number(hourMatch[1]) * HOUR_NS;
  if (minuteMatch) total += Number(minuteMatch[1]) * MINUTE_NS;
  if (secondMatch) total += Number(secondMatch[1]) * SECOND_NS;
  if (msMatch) total += Number(msMatch[1]) * 1_000_000;
  return total;
}

function splitDuration(durationNs) {
  const totalSeconds = Math.max(1, Math.floor(durationNs / SECOND_NS));
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  return {
    hours: String(hours),
    minutes: String(minutes),
    seconds: String(seconds),
  };
}

function formatDuration(durationNs, t) {
  const totalSeconds = Math.max(1, Math.floor(durationNs / SECOND_NS));
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  const parts = [];
  if (hours > 0) parts.push(t('admin.cronjobs.duration.hours', { count: hours }));
  if (minutes > 0) parts.push(t('admin.cronjobs.duration.minutes', { count: minutes }));
  if (seconds > 0 || parts.length === 0) parts.push(t('admin.cronjobs.duration.seconds', { count: seconds }));
  return parts.join(' ');
}

function buildDurationNs(form) {
  const hours = Math.max(0, Number(form.hours) || 0);
  const minutes = Math.max(0, Number(form.minutes) || 0);
  const seconds = Math.max(0, Number(form.seconds) || 0);
  const total = hours * HOUR_NS + minutes * MINUTE_NS + seconds * SECOND_NS;
  return Math.max(SECOND_NS, total);
}

function getLatestRunTime(successLast, failureLast) {
  const successTime = successLast ? new Date(successLast) : null;
  const failureTime = failureLast ? new Date(failureLast) : null;
  const successTs = successTime && !Number.isNaN(successTime.getTime()) ? successTime.getTime() : null;
  const failureTs = failureTime && !Number.isNaN(failureTime.getTime()) ? failureTime.getTime() : null;

  if (successTs === null && failureTs === null) return null;
  if (successTs === null) return failureTime;
  if (failureTs === null) return successTime;
  return successTs >= failureTs ? successTime : failureTime;
}

function CronJobs() {
  const { t } = useTranslation();
  const [cronJobs, setCronJobs] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedCronJob, setSelectedCronJob] = useState(null);
  const [form, setForm] = useState(DEFAULT_FORM);

  const fetchCronJobs = async (page = currentPage) => {
    setLoading(true);
    try {
      const response = await getCronJobList({
        limit: PAGE_SIZE,
        offset: (page - 1) * PAGE_SIZE,
      });
      if (response.code === 200) {
        setCronJobs(response.data.cronjobs || []);
        setTotalCount(response.data.count || 0);
      }
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.cronjobs.toast.fetchListFailed'),
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCronJobs(currentPage);
  }, [currentPage]);

  const openEditModal = (cronJob) => {
    const durationNs = cronJob.schedule_ns ?? parseDurationNs(cronJob.schedule);
    setSelectedCronJob(cronJob);
    setForm(splitDuration(durationNs));
    setIsModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      if (selectedCronJob) {
        const durationNs = buildDurationNs(form);
        const currentNs = selectedCronJob.schedule_ns ?? parseDurationNs(selectedCronJob.schedule);
        const payload = {};
        if (durationNs !== currentNs) payload.schedule = durationNs;
        if (Object.keys(payload).length === 0) {
          setIsModalOpen(false);
          return;
        }
        const response = await updateCronJob(selectedCronJob.id, payload);
        if (response.code === 200) {
          toast.success({
            description: t('admin.cronjobs.toast.updateSuccess'),
          });
        }
      }
      setIsModalOpen(false);
      fetchCronJobs(currentPage);
    } catch (error) {
      toast.danger({
        description: error.message || t('admin.cronjobs.toast.updateFailed'),
      });
    }
  };

  const formatDateTime = (value) => {
    if (!value) return t('common.notAvailable');
    const date = new Date(value);
    if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) return t('common.notAvailable');
    return date.toLocaleString();
  };

  const getNextRun = (durationNs, successLast, failureLast) => {
    if (!durationNs) return null;
    const latestRun = getLatestRunTime(successLast, failureLast);
    const base = latestRun || new Date();
    if (Number.isNaN(base.getTime())) return null;
    return new Date(base.getTime() + durationNs / 1_000_000);
  };

  const columns = [
    { key: 'id', label: t('admin.cronjobs.columns.id'), width: '8%' },
    { key: 'name', label: t('admin.cronjobs.columns.name'), width: '14%' },
    {
      key: 'description',
      label: t('admin.cronjobs.columns.description'),
      width: '20%',
    },
    {
      key: 'schedule',
      label: t('admin.cronjobs.columns.schedule'),
      width: '14%',
    },
    {
      key: 'successCount',
      label: t('admin.cronjobs.columns.successCount'),
      width: '10%',
    },
    {
      key: 'failureCount',
      label: t('admin.cronjobs.columns.failureCount'),
      width: '10%',
    },
    {
      key: 'successLast',
      label: t('admin.cronjobs.columns.successLast'),
      width: '12%',
    },
    {
      key: 'failureLast',
      label: t('admin.cronjobs.columns.failureLast'),
      width: '12%',
    },
    { key: 'next', label: t('admin.cronjobs.columns.next'), width: '12%' },
    {
      key: 'actions',
      label: t('admin.cronjobs.columns.actions'),
      width: '8%',
    },
  ];

  const renderCell = (cronJob, column) => {
    const durationNs = cronJob.schedule_ns ?? parseDurationNs(cronJob.schedule);
    switch (column.key) {
      case 'id':
        return <span className="text-neutral-50 font-medium">#{cronJob.id}</span>;
      case 'name':
        return <span className="text-neutral-50 font-medium">{cronJob.name}</span>;
      case 'description':
        return <span className="text-neutral-300 text-sm">{cronJob.description || t('common.none')}</span>;
      case 'schedule':
        return <span className="text-neutral-300 font-mono text-sm">{formatDuration(durationNs, t)}</span>;
      case 'successCount':
        return <span className="text-emerald-300 text-sm font-medium">{cronJob.success ?? 0}</span>;
      case 'failureCount':
        return <span className="text-rose-300 text-sm font-medium">{cronJob.failure ?? 0}</span>;
      case 'successLast':
        return <span className="text-neutral-300 text-sm">{formatDateTime(cronJob.success_last)}</span>;
      case 'failureLast':
        return <span className="text-neutral-300 text-sm">{formatDateTime(cronJob.failure_last)}</span>;
      case 'next': {
        const nextRun = getNextRun(durationNs, cronJob.success_last, cronJob.failure_last);
        return (
          <span className="text-neutral-300 text-sm">
            {nextRun ? formatDateTime(nextRun) : t('admin.cronjobs.time.unknown')}
          </span>
        );
      }
      case 'actions':
        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              onClick={(e) => {
                e.stopPropagation();
                openEditModal(cronJob);
              }}
            >
              <IconEdit size={18} />
            </Button>
          </div>
        );
      default:
        return cronJob[column.key];
    }
  };

  const previewDurationNs = buildDurationNs(form);
  const previewNextRun = selectedCronJob
    ? getNextRun(previewDurationNs, selectedCronJob.success_last, selectedCronJob.failure_last)
    : null;

  return (
    <div className="w-full mx-auto">
      <div className="rounded-md bg-neutral-900 overflow-hidden p-6">
        <List
          data={cronJobs}
          columns={columns}
          renderCell={renderCell}
          onRowClick={openEditModal}
          loading={loading}
          empty={cronJobs.length === 0}
          emptyContent={t('admin.cronjobs.empty')}
        />

        {totalCount > PAGE_SIZE && (
          <Pagination
            className="mt-4"
            current={currentPage}
            total={Math.ceil(totalCount / PAGE_SIZE)}
            onChange={setCurrentPage}
            showTotal
            totalItems={totalCount}
            showJumpTo
          />
        )}
      </div>

      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={t('admin.cronjobs.modal.editTitle')}
        size="lg"
        footer={
          <div className="flex justify-end gap-2">
            <Button variant="ghost" onClick={() => setIsModalOpen(false)}>
              {t('common.cancel')}
            </Button>
            <Button variant="primary" onClick={handleSubmit}>
              {t('common.save')}
            </Button>
          </div>
        }
      >
        <div className="space-y-4">
          <p className="text-sm text-neutral-400">{t('admin.cronjobs.modal.scheduleOnlyHint')}</p>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.nameLabel')}
            </label>
            <Input type="text" value={selectedCronJob?.name || ''} fullWidth disabled />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.descriptionLabel')}
            </label>
            <Input type="text" value={selectedCronJob?.description || ''} fullWidth disabled />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.successCountLabel')}
            </label>
            <Input type="text" value={String(selectedCronJob?.success ?? 0)} fullWidth disabled />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.failureCountLabel')}
            </label>
            <Input type="text" value={String(selectedCronJob?.failure ?? 0)} fullWidth disabled />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.successLastLabel')}
            </label>
            <Input type="text" value={formatDateTime(selectedCronJob?.success_last)} fullWidth disabled />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.failureLastLabel')}
            </label>
            <Input type="text" value={formatDateTime(selectedCronJob?.failure_last)} fullWidth disabled />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.nextLabel')}
            </label>
            <Input
              type="text"
              value={previewNextRun ? formatDateTime(previewNextRun) : t('admin.cronjobs.time.unknown')}
              fullWidth
              disabled
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.scheduleLabel')}
            </label>
            <div className="grid grid-cols-3 gap-4">
              <Input
                type="number"
                min="0"
                value={form.hours}
                onChange={(e) => setForm((prev) => ({ ...prev, hours: e.target.value }))}
                placeholder={t('admin.cronjobs.form.hoursPlaceholder')}
                fullWidth
              />
              <Input
                type="number"
                min="0"
                value={form.minutes}
                onChange={(e) => setForm((prev) => ({ ...prev, minutes: e.target.value }))}
                placeholder={t('admin.cronjobs.form.minutesPlaceholder')}
                fullWidth
              />
              <Input
                type="number"
                min="0"
                value={form.seconds}
                onChange={(e) => setForm((prev) => ({ ...prev, seconds: e.target.value }))}
                placeholder={t('admin.cronjobs.form.secondsPlaceholder')}
                fullWidth
              />
            </div>
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.durationPreviewLabel')}
            </label>
            <Input type="text" value={formatDuration(previewDurationNs, t)} fullWidth disabled />
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default CronJobs;
