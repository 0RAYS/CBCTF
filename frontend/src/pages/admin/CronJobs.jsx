import { useEffect, useState } from 'react';
import { IconEdit } from '@tabler/icons-react';
import { Button, Input, List, Modal, Pagination, StatusTag } from '../../components/common';
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
    { key: 'name', label: t('admin.cronjobs.columns.name'), width: '20%' },
    { key: 'description', label: t('admin.cronjobs.columns.description'), width: '35%' },
    { key: 'schedule', label: t('admin.cronjobs.columns.schedule'), width: '20%' },
    { key: 'status', label: t('admin.cronjobs.columns.status'), width: '15%' },
    { key: 'actions', label: t('admin.cronjobs.columns.actions'), width: '10%' },
  ];

  const renderStatus = (status) => {
    if (status === 'enabled') {
      return <StatusTag type="success" text={t('admin.cronjobs.status.enabled')} />;
    }
    if (status === 'running') {
      return <StatusTag type="info" text={t('admin.cronjobs.status.running')} />;
    }
    return <StatusTag type="warning" text={t('admin.cronjobs.status.disabled')} />;
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
      case 'status':
        return renderStatus(cronJob.status);
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
              {t('admin.cronjobs.form.statusLabel')}
            </label>
            <Input
              type="text"
              value={t(`admin.cronjobs.status.${selectedCronJob?.status || 'disabled'}`)}
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
