import { useState, useEffect } from 'react';
import { toast } from '../../utils/toast';
import { getGenerators, startGenerators, stopGenerators } from '../../api/admin/generators';
import { getChallengeList } from '../../api/admin/challenge';
import { Modal } from '../../components/common';
import { Button, Pagination, Card, EmptyState, StatCard } from '../../components/common';
import { motion } from 'motion/react';
import { IconPlayerPlay, IconBan, IconRefresh, IconCheck, IconX, IconTrash } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';

const STATUS_STYLES = {
  waiting: 'bg-yellow-400/10 text-yellow-400 border-yellow-400/30',
  pending: 'bg-geek-400/10 text-geek-400 border-geek-400/30',
  running: 'bg-green-400/10 text-green-400 border-green-400/30',
  stopped: 'bg-neutral-500/10 text-neutral-400 border-neutral-500/30',
};

function GeneratorStatusBadge({ status, t }) {
  const style = STATUS_STYLES[status] ?? STATUS_STYLES.stopped;
  return (
    <span className={`inline-block px-2 py-0.5 rounded border text-xs font-mono ${style}`}>
      {t(`admin.contests.generators.status.${status}`, status)}
    </span>
  );
}

const PAGE_SIZE = 20;

function AdminGenerators() {
  const { t } = useTranslation();

  const [generators, setGenerators] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedIds, setSelectedIds] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showDeleted, setShowDeleted] = useState(false);

  const [dynamicChallenges, setDynamicChallenges] = useState([]);
  const [startModalOpen, setStartModalOpen] = useState(false);
  const [startCounts, setStartCounts] = useState({});

  const fetchGenerators = async (page = 1, deleted = showDeleted) => {
    setLoading(true);
    try {
      const res = await getGenerators({
        limit: PAGE_SIZE,
        offset: (page - 1) * PAGE_SIZE,
        ...(deleted && { deleted: true }),
      });
      setGenerators(res.data?.generators ?? []);
      setTotalCount(res.data?.count ?? 0);
    } catch {
      toast.error(t('admin.generators.toast.fetchFailed'));
    } finally {
      setLoading(false);
    }
  };

  const fetchDynamicChallenges = async () => {
    try {
      const res = await getChallengeList({ limit: 100, offset: 0, type: 'dynamic' });
      setDynamicChallenges(res.data?.challenges ?? []);
    } catch {
      // silent — modal will show empty list
    }
  };

  useEffect(() => {
    fetchGenerators(1);
    fetchDynamicChallenges();
  }, []);

  const handlePageChange = (page) => {
    setCurrentPage(page);
    setSelectedIds([]);
    fetchGenerators(page);
  };

  const toggleShowDeleted = () => {
    const next = !showDeleted;
    setShowDeleted(next);
    setCurrentPage(1);
    setSelectedIds([]);
    fetchGenerators(1, next);
  };

  const toggleSelect = (id) => {
    setSelectedIds((prev) => (prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]));
  };

  const toggleSelectAll = () => {
    if (selectedIds.length === generators.length) {
      setSelectedIds([]);
    } else {
      setSelectedIds(generators.map((g) => g.id));
    }
  };

  const handleStop = async () => {
    if (selectedIds.length === 0) return;
    try {
      await stopGenerators(selectedIds);
      toast.success(t('admin.generators.toast.stopSuccess'));
      setSelectedIds([]);
      fetchGenerators(currentPage);
    } catch {
      toast.error(t('admin.generators.toast.stopFailed'));
    }
  };

  const openStartModal = () => {
    setStartCounts({});
    setStartModalOpen(true);
  };

  const handleCountChange = (id, value) => {
    const count = Math.max(0, parseInt(value, 10) || 0);
    setStartCounts((prev) => ({ ...prev, [id]: count }));
  };

  const handleStart = async () => {
    const challenges = Object.entries(startCounts)
      .filter(([, count]) => count > 0)
      .flatMap(([id, count]) => Array(count).fill(id));
    if (challenges.length === 0) {
      toast.error(t('admin.generators.toast.selectRequired'));
      return;
    }
    try {
      await startGenerators(challenges);
      setStartModalOpen(false);
      setCurrentPage(1);
      fetchGenerators(1);
    } catch {
      toast.error(t('admin.generators.toast.startFailed'));
    }
  };

  const totalSuccesses = generators.reduce((sum, g) => sum + (g.success ?? 0), 0);
  const totalFailures = generators.reduce((sum, g) => sum + (g.failure ?? 0), 0);

  const formatTime = (ts) => {
    if (!ts) return '—';
    return new Date(ts).toLocaleString();
  };

  return (
    <div className="flex flex-col gap-6 p-6">
      {/* Stat cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <StatCard title={t('admin.generators.stats.total')} value={totalCount} icon={<IconPlayerPlay size={20} />} />
        <StatCard title={t('admin.generators.stats.successes')} value={totalSuccesses} icon={<IconCheck size={20} />} />
        <StatCard title={t('admin.generators.stats.failures')} value={totalFailures} icon={<IconX size={20} />} />
      </div>

      {/* Toolbar */}
      <div className="flex flex-wrap gap-2 items-center">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => fetchGenerators(currentPage)}
          leftIcon={<IconRefresh size={14} />}
        >
          {t('common.refresh')}
        </Button>
        <Button variant="primary" size="sm" onClick={openStartModal} leftIcon={<IconPlayerPlay size={14} />}>
          {t('admin.generators.startButton')}
        </Button>
        <Button
          variant={showDeleted ? 'danger' : 'ghost'}
          size="sm"
          onClick={toggleShowDeleted}
          leftIcon={<IconTrash size={14} />}
        >
          {t('admin.generators.showDeleted')}
        </Button>
        {selectedIds.length > 0 && (
          <Button variant="danger" size="sm" onClick={handleStop} leftIcon={<IconBan size={14} />}>
            {t('admin.generators.stopButton')} ({selectedIds.length})
          </Button>
        )}
      </div>

      {/* List */}
      <Card>
        {loading ? (
          <div className="flex justify-center py-12 text-neutral-400 text-sm">{t('common.loading')}</div>
        ) : generators.length === 0 ? (
          <EmptyState title={t('admin.generators.noGenerators')} description={t('admin.generators.noGeneratorsDesc')} />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-neutral-300">
              <thead>
                <tr className="border-b border-neutral-700 text-neutral-400 text-xs uppercase tracking-wider">
                  <th className="py-3 px-4 text-left w-10" scope="col">
                    <input
                      type="checkbox"
                      className="accent-geek-400"
                      checked={generators.length > 0 && selectedIds.length === generators.length}
                      onChange={toggleSelectAll}
                    />
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.id')}
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.name')}
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.challengeId')}
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.success')}
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.successLast')}
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.failure')}
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.failureLast')}
                  </th>
                  <th className="py-3 px-4 text-left" scope="col">
                    {t('admin.generators.columns.status')}
                  </th>
                </tr>
              </thead>
              <tbody>
                {generators.map((gen) => (
                  <motion.tr
                    key={gen.id}
                    className="border-b border-neutral-800 hover:bg-neutral-800/40 transition-colors"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                  >
                    <td className="py-3 px-4">
                      <input
                        type="checkbox"
                        className="accent-geek-400"
                        checked={selectedIds.includes(gen.id)}
                        onChange={() => toggleSelect(gen.id)}
                      />
                    </td>
                    <td className="py-3 px-4 font-mono text-xs text-neutral-500">{gen.id}</td>
                    <td className="py-3 px-4 font-mono text-xs text-neutral-200">{gen.name}</td>
                    <td className="py-3 px-4 text-neutral-400">{gen.challenge_id}</td>
                    <td className="py-3 px-4 text-green-400">{gen.success ?? 0}</td>
                    <td className="py-3 px-4 text-neutral-400 text-xs">{formatTime(gen.success_last)}</td>
                    <td className="py-3 px-4 text-red-400">{gen.failure ?? 0}</td>
                    <td className="py-3 px-4 text-neutral-400 text-xs">{formatTime(gen.failure_last)}</td>
                    <td className="py-3 px-4">
                      <GeneratorStatusBadge status={gen.status} t={t} />
                    </td>
                  </motion.tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      {totalCount > PAGE_SIZE && (
        <Pagination
          current={currentPage}
          total={Math.ceil(totalCount / PAGE_SIZE)}
          totalItems={totalCount}
          showTotal
          onChange={handlePageChange}
        />
      )}

      {/* Start Modal */}
      <Modal
        isOpen={startModalOpen}
        onClose={() => setStartModalOpen(false)}
        title={t('admin.generators.selectChallenges')}
      >
        <div className="flex flex-col gap-4">
          {dynamicChallenges.length === 0 ? (
            <p className="text-neutral-400 text-sm py-4 text-center">{t('admin.generators.noDynamicChallenges')}</p>
          ) : (
            <div className="flex flex-col gap-2 max-h-72 overflow-y-auto">
              {dynamicChallenges.map((c) => {
                const id = c.rand_id ?? c.id;
                return (
                  <div
                    key={id}
                    className="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 transition-colors"
                  >
                    <span className="text-sm text-neutral-200 flex-1">{c.title ?? c.name}</span>
                    <span className="text-xs text-neutral-500 font-mono">{id}</span>
                    <input
                      type="number"
                      min={0}
                      step={1}
                      value={startCounts[id] ?? 0}
                      onChange={(e) => handleCountChange(id, e.target.value)}
                      className="w-16 bg-neutral-800 border border-neutral-700 rounded px-2 py-1 text-sm text-neutral-200 text-center focus:outline-none focus:border-geek-400"
                    />
                  </div>
                );
              })}
            </div>
          )}
          <div className="flex justify-end gap-2 pt-2">
            <Button variant="ghost" size="sm" onClick={() => setStartModalOpen(false)}>
              {t('common.cancel')}
            </Button>
            <Button
              variant="primary"
              size="sm"
              onClick={handleStart}
              disabled={!Object.values(startCounts).some((c) => c > 0)}
              leftIcon={<IconPlayerPlay size={14} />}
            >
              {t('admin.generators.startSelected')}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default AdminGenerators;
