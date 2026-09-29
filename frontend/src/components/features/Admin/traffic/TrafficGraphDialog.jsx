import { useEffect, useEffectEvent, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Modal } from '../../../common';
import TrafficCanvas from './TrafficCanvas.jsx';
import TrafficControls from './TrafficControls.jsx';
import TrafficHeader from './TrafficHeader.jsx';
import TrafficRankings from './TrafficRankings.jsx';
import TrafficSelectionPanel from './TrafficSelectionPanel.jsx';
import TrafficSummary from './TrafficSummary.jsx';
import TrafficTimeline from './TrafficTimeline.jsx';
import { filterTraffic, resolveTrafficSelection } from './trafficPresentation.js';
import useTrafficPlayback from './useTrafficPlayback.js';
import useTrafficSession from './useTrafficSession.js';
import useTrafficAnalysis from './useTrafficAnalysis.js';
import TrafficAnalysisPanel from './TrafficAnalysisPanel.jsx';

export default function TrafficGraphDialog({ isOpen, onClose, container, contestId, teamId, fetchTraffic }) {
  const { t } = useTranslation();
  const session = useTrafficSession({
    isOpen,
    container,
    contestId,
    teamId,
    fetchTraffic,
  });
  const analysis = useTrafficAnalysis({
    isOpen,
    victimId: container?.id,
    contestId,
    teamId,
  });
  const { topology, universeNodes, isFetching, scopeKey, selectionVersion } = session;
  const playback = useTrafficPlayback({
    isOpen,
    containerId: container?.id,
    scopeKey,
    topology,
    isFetching,
  });
  const { shift, slice } = playback;
  const [protocolFilter, setProtocolFilter] = useState(new Set());
  const [selectedEdgeId, setSelectedEdgeId] = useState('');
  const [selectedNodeId, setSelectedNodeId] = useState('');
  const [selectedFlow, setSelectedFlow] = useState(null);

  useEffect(() => {
    if (isOpen) setProtocolFilter(new Set());
  }, [isOpen, scopeKey]);

  useEffect(() => {
    setSelectedEdgeId('');
    setSelectedNodeId('');
    setSelectedFlow(null);
  }, [selectionVersion]);

  // Custom fetchers may be inline callbacks; read the latest without making their identity a request trigger.
  const fetchFrame = useEffectEvent(() => {
    void session.fetchData({ nextShift: shift, nextSlice: slice });
  });
  const invalidateFrame = useEffectEvent(() => session.invalidateFrame());
  useEffect(() => {
    if (!isOpen || !container?.id) return;
    fetchFrame();
    return () => invalidateFrame();
  }, [isOpen, scopeKey, shift, slice]);

  const graph = useMemo(
    () => filterTraffic(topology, universeNodes, protocolFilter),
    [topology, universeNodes, protocolFilter]
  );
  const selection = resolveTrafficSelection(graph.nodes, graph.edges, selectedNodeId, selectedEdgeId);
  if (!isOpen) return null;

  return (
    <>
      <style>{`
        @keyframes traffic-flow { from { stroke-dashoffset: 0; } to { stroke-dashoffset: -120; } }
        .traffic-line { stroke-dasharray: 11 10; animation: traffic-flow 10s linear infinite; }
      `}</style>
      <Modal
        isOpen={isOpen}
        onClose={onClose}
        title={t('admin.contests.trafficGraph.title')}
        size="full"
        className="!bg-neutral-900/95 !border-neutral-600"
        showHeader={false}
        bodyClassName="overflow-y-auto p-4 sm:p-5"
      >
        <div className="flex flex-col gap-3 text-neutral-100">
          <div className="grid grid-cols-1 gap-3 xl:grid-cols-[minmax(0,1fr)_360px]">
            <TrafficHeader topology={topology} container={container} />
            <TrafficControls
              playback={playback}
              isFetching={isFetching}
              onRefresh={() => session.fetchData({ nextShift: shift, nextSlice: slice })}
              onDownload={session.downloadTraffic}
            />
          </div>
          <TrafficSummary summary={topology?.summary || {}} />
          {session.fetchError ? (
            <p role="alert" className="text-sm text-red-300">
              {t('admin.contests.trafficGraph.toast.fetchFailed')}
            </p>
          ) : null}
          <div className="grid grid-cols-1 gap-3 xl:grid-cols-[minmax(0,1.55fr)_360px]">
            <TrafficCanvas
              scopeKey={scopeKey}
              {...graph}
              {...selection}
              universeNodes={universeNodes}
              protocolFilter={protocolFilter}
              setProtocolFilter={setProtocolFilter}
              selectedEdgeId={selectedEdgeId}
              onSelectEdge={(id) => {
                setSelectedEdgeId(id);
                setSelectedNodeId('');
                setSelectedFlow(graph.edges.find((edge) => edge.id === id));
              }}
              onSelectNode={(id) => {
                setSelectedNodeId(id);
                setSelectedEdgeId('');
                setSelectedFlow(null);
              }}
            />
            <div className="grid gap-3">
              <TrafficSelectionPanel {...selection} />
              <TrafficTimeline
                timeline={topology?.timeline || []}
                windowInfo={playback.windowInfo}
                setShift={playback.setShift}
              />
              <TrafficRankings topology={topology} playback={playback} />
            </div>
          </div>
          <TrafficAnalysisPanel
            key={scopeKey}
            analysis={analysis}
            topology={topology}
            victimId={container?.id}
            contestId={contestId}
            selectedNodeId={selectedNodeId}
            selectedEdge={selectedFlow}
            onSeek={(offset) => {
              playback.setIsPlaying(false);
              playback.setShift(Math.min(Math.max(offset, 0), playback.maxShift));
            }}
          />
        </div>
      </Modal>
    </>
  );
}
