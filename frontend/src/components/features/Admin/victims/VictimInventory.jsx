import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from '../../../../utils/toast';
import { downloadBlobResponse } from '../../../../utils/fileDownload';
import TrafficGraphDialog from '../traffic/TrafficGraphDialog';
import { VictimFilters } from './VictimFilters';
import { VictimStats } from './VictimStats';
import { VictimTable } from './VictimTable';
import VictimStopDialog from './VictimStopDialog';
import VictimLogDialog from './VictimLogDialog';
import VictimStatusBadge from './VictimStatusBadge';
import useVictimList from './useVictimList';
import { isVictimStoppable } from './victimPayload';

export default function VictimInventory({ scope, renderQuickActions }) {
  const { t, i18n } = useTranslation();
  const list = useVictimList(scope);
  const [stopOpen, setStopOpen] = useState(false);
  const [trafficVictim, setTrafficVictim] = useState(null);
  const [logVictim, setLogVictim] = useState(null);
  const { translationKey, contestId } = scope;

  const stopSelected = async () => {
    if (!list.selectedContainers.length) {
      toast.warning({ description: t(`${translationKey}.toast.selectStopRequired`) });
      return;
    }
    try {
      const response = await scope.stopVictims(list.selectedContainers);
      if (response.code === 200) {
        toast.success({ description: t(`${translationKey}.toast.taskDispatched`) });
        list.setSelectedContainers([]);
        list.refresh();
      }
    } catch (error) {
      toast.danger({ description: error.message || t(`${translationKey}.toast.taskDispatchFailed`) });
    }
    setStopOpen(false);
  };

  const downloadTraffic = async (victim) => {
    try {
      const response = await scope.downloadTraffic(victim);
      if (response.headers?.['file'] === 'true') downloadBlobResponse(response, `traffic_${victim.id}.zip`);
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.teamContainers.toast.downloadTrafficFailed') });
    }
  };

  const formatTime = (value) =>
    value
      ? new Date(value).toLocaleString(i18n.language || 'en-US', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
        })
      : '-';
  const formatRemaining = (remaining) => {
    if (!remaining || remaining <= 0) return t(`${translationKey}.status.stopped`);
    return `${Math.floor(remaining / 3600)}h ${Math.floor((remaining % 3600) / 60)}m ${Math.floor(remaining % 60)}s`;
  };
  const getContainerStatusStyle = (remaining) =>
    !remaining || remaining <= 0
      ? 'text-red-400 bg-red-400/10 border-red-400/30'
      : 'text-green-400 bg-green-400/10 border-green-400/30';

  return (
    <div className="w-full mx-auto space-y-6">
      <VictimStats stats={list.stats} t={t} translationKey={translationKey} />
      {renderQuickActions?.(list.refresh)}
      <VictimFilters
        scope={scope}
        filters={list.filters}
        onFilterChange={list.onFilterChange}
        onResetFilters={list.onResetFilters}
      />
      <VictimTable
        {...list}
        t={t}
        translationKey={translationKey}
        contestId={contestId}
        onRefreshIntervalChange={list.setRefreshInterval}
        onRefresh={list.refresh}
        onOpenStopModal={() => setStopOpen(true)}
        onViewTrafficGraph={setTrafficVictim}
        onDownloadTraffic={downloadTraffic}
        onViewLogs={setLogVictim}
        isVictimStoppable={isVictimStoppable}
        formatTime={formatTime}
        formatRemaining={formatRemaining}
        getContainerStatusStyle={getContainerStatusStyle}
        VictimStatusBadge={VictimStatusBadge}
      />
      <VictimStopDialog
        t={t}
        isOpen={stopOpen}
        onClose={() => setStopOpen(false)}
        onConfirm={stopSelected}
        selectedCount={list.selectedContainers.length}
        translationKey={translationKey}
        contestId={contestId}
      />
      <TrafficGraphDialog
        isOpen={!!trafficVictim}
        onClose={() => setTrafficVictim(null)}
        container={trafficVictim}
        contestId={contestId}
        teamId={trafficVictim?.team_id}
        fetchTraffic={scope.fetchTraffic}
      />
      <VictimLogDialog
        victim={logVictim}
        onClose={() => setLogVictim(null)}
        loadPods={scope.loadPods}
        loadLogs={scope.loadLogs}
        translationKey={translationKey}
      />
    </div>
  );
}
