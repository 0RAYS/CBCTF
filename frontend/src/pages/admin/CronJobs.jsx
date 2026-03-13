import { useEffect, useState } from 'react';
import { IconEdit } from '@tabler/icons-react';
import { Button, Input, List, Modal, Pagination } from '../../components/common';
import { getCronJobList, updateCronJob } from '../../api/admin/cronjob';
import { toast } from '../../utils/toast';
import { useTranslation } from 'react-i18next';

const PAGE_SIZE = 20;

const DEFAULT_FORM = {
  schedule: '',
};

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
      toast.danger({ description: error.message || t('admin.cronjobs.toast.fetchListFailed') });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCronJobs(currentPage);
  }, [currentPage]);

  const openEditModal = (cronJob) => {
    setSelectedCronJob(cronJob);
    setForm({
      schedule: cronJob.schedule,
    });
    setIsModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      if (selectedCronJob) {
        const payload = {};
        if (form.schedule !== selectedCronJob.schedule) payload.schedule = form.schedule;
        const response = await updateCronJob(selectedCronJob.id, payload);
        if (response.code === 200) {
          toast.success({ description: t('admin.cronjobs.toast.updateSuccess') });
        }
      }
      setIsModalOpen(false);
      fetchCronJobs(currentPage);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.cronjobs.toast.updateFailed') });
    }
  };

  const columns = [
    { key: 'id', label: t('admin.cronjobs.columns.id'), width: '10%' },
    { key: 'name', label: t('admin.cronjobs.columns.name'), width: '16%' },
    { key: 'description', label: t('admin.cronjobs.columns.description'), width: '24%' },
    { key: 'schedule', label: t('admin.cronjobs.columns.schedule'), width: '18%' },
    { key: 'last', label: t('admin.cronjobs.columns.last'), width: '16%' },
    { key: 'next', label: t('admin.cronjobs.columns.next'), width: '16%' },
    { key: 'actions', label: t('admin.cronjobs.columns.actions'), width: '10%' },
  ];

  const formatDateTime = (value) => {
    if (!value) return t('common.notAvailable');
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return t('admin.cronjobs.time.invalid');
    return date.toLocaleString();
  };

  const getEveryMs = (schedule) => {
    const match = /^@every\s+(\d+)(ms|s|m|h)$/i.exec(schedule?.trim() || '');
    if (!match) return null;
    const value = Number(match[1]);
    const unit = match[2].toLowerCase();
    const unitMs = { ms: 1, s: 1000, m: 60 * 1000, h: 60 * 60 * 1000 };
    return value * unitMs[unit];
  };

  const matchField = (field, value, min, max) => {
    if (field === '*') return true;
    if (/^\d+$/.test(field)) return Number(field) === value;
    if (/^\*\/\d+$/.test(field)) return value % Number(field.slice(2)) === 0;
    if (/^\d+-\d+$/.test(field)) {
      const [start, end] = field.split('-').map(Number);
      return value >= start && value <= end;
    }
    if (/^\d+(,\d+)+$/.test(field)) return field.split(',').map(Number).includes(value);
    if (/^\d+-\d+\/\d+$/.test(field)) {
      const [range, step] = field.split('/');
      const [start, end] = range.split('-').map(Number);
      return value >= start && value <= end && (value - start) % Number(step) === 0;
    }
    if (/^\*\/\d+(,\*\/\d+)*$/.test(field)) {
      return field.split(',').some((part) => value % Number(part.slice(2)) === 0);
    }
    return value >= min && value <= max && field.split(',').some((part) => matchField(part, value, min, max));
  };

  const matchesCron = (schedule, date) => {
    const parts = (schedule || '').trim().split(/\s+/);
    if (parts.length !== 5 && parts.length !== 6) return false;
    const normalized = parts.length === 5 ? ['0', ...parts] : parts;
    const [second, minute, hour, day, month, week] = normalized;
    return (
      matchField(second, date.getSeconds(), 0, 59) &&
      matchField(minute, date.getMinutes(), 0, 59) &&
      matchField(hour, date.getHours(), 0, 23) &&
      matchField(day, date.getDate(), 1, 31) &&
      matchField(month, date.getMonth() + 1, 1, 12) &&
      matchField(week, date.getDay(), 0, 6)
    );
  };

  const getNextRun = (schedule, last) => {
    const everyMs = getEveryMs(schedule);
    const base = last ? new Date(last) : new Date();
    if (!Number.isNaN(base.getTime()) && everyMs) {
      return new Date(base.getTime() + everyMs);
    }

    const start = new Date();
    start.setMilliseconds(0);
    const limit = 366 * 24 * 60 * 60;
    for (let offset = 1; offset <= limit; offset++) {
      const candidate = new Date(start.getTime() + offset * 1000);
      if (matchesCron(schedule, candidate)) {
        return candidate;
      }
    }
    return null;
  };

  const renderCell = (cronJob, column) => {
    switch (column.key) {
      case 'id':
        return <span className="text-neutral-50 font-medium">#{cronJob.id}</span>;
      case 'name':
        return <span className="text-neutral-50 font-medium">{cronJob.name}</span>;
      case 'description':
        return <span className="text-neutral-300 text-sm">{cronJob.description || t('common.none')}</span>;
      case 'schedule':
        return <span className="text-neutral-300 font-mono text-sm">{cronJob.schedule}</span>;
      case 'last':
        return <span className="text-neutral-300 text-sm">{formatDateTime(cronJob.last)}</span>;
      case 'next': {
        const nextRun = getNextRun(cronJob.schedule, cronJob.last);
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
              {t('admin.cronjobs.form.lastLabel')}
            </label>
            <Input type="text" value={formatDateTime(selectedCronJob?.last)} fullWidth disabled />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.nextLabel')}
            </label>
            <Input
              type="text"
              value={
                selectedCronJob
                  ? formatDateTime(getNextRun(form.schedule || selectedCronJob.schedule, selectedCronJob.last))
                  : t('admin.cronjobs.time.unknown')
              }
              fullWidth
              disabled
            />
          </div>
          <div>
            <label className="block text-neutral-300 text-sm font-medium mb-2">
              {t('admin.cronjobs.form.scheduleLabel')}
            </label>
            <Input
              type="text"
              value={form.schedule}
              onChange={(e) => setForm((prev) => ({ ...prev, schedule: e.target.value }))}
              placeholder={t('admin.cronjobs.form.schedulePlaceholder')}
              fullWidth
              required
            />
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default CronJobs;
