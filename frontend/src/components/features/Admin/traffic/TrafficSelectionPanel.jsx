import { useTranslation } from 'react-i18next';
import { Card, Chip } from '../../../common';
import { edgeChipClass, ellipsis, formatBytes, formatProcessLabel } from './trafficPresentation.js';

function ProcessList({ processes }) {
  const { t } = useTranslation();
  if (!processes?.length) return null;
  return (
    <div className="mt-3 rounded-lg border border-neutral-700 bg-black/20 px-2 py-2">
      <div className="text-[11px] text-neutral-500">{t('admin.contests.trafficGraph.panel.processes')}</div>
      <div className="mt-2 grid gap-1.5">
        {processes.slice(0, 3).map((process, index) => (
          <div
            key={`${formatProcessLabel(process)}-${index}`}
            className="flex items-center justify-between gap-2 text-[11px] font-mono"
          >
            <span className="min-w-0 truncate text-neutral-100">{formatProcessLabel(process) || '--'}</span>
            <span className="shrink-0 text-geek-400">{formatBytes(process.bytes)}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export default function TrafficSelectionPanel({ selectedEdge, selectedNode }) {
  const { t } = useTranslation();
  return (
    <Card padding="sm" className="rounded-2xl border-neutral-600 bg-neutral-900">
      <div className="grid gap-2">
        <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
          <div className="flex items-start justify-between gap-2">
            <div className="min-w-0">
              <div className="text-xs font-mono text-neutral-300">
                {t('admin.contests.trafficGraph.panel.selectedFlow')}
              </div>
              <div className="mt-1 truncate text-[11px] text-neutral-500">
                {selectedEdge
                  ? `${selectedEdge.source} -> ${selectedEdge.target}`
                  : t('admin.contests.trafficGraph.panel.selectedFlowHint')}
              </div>
            </div>
            {selectedEdge ? (
              <Chip
                label={t(`admin.contests.trafficGraph.direction.${selectedEdge.direction || 'internal'}`)}
                variant="tag"
                size="sm"
                colorClass={edgeChipClass(selectedEdge.direction)}
              />
            ) : null}
          </div>
          {selectedEdge ? (
            <>
              <div className="mt-3 grid grid-cols-3 gap-2 text-[11px] font-mono">
                <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                  <div className="text-neutral-500">{t('admin.contests.trafficGraph.panel.edgeBytes')}</div>
                  <div className="mt-1 text-geek-400">{formatBytes(selectedEdge.bytes)}</div>
                </div>
                <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                  <div className="text-neutral-500">{t('admin.contests.trafficGraph.panel.edgePackets')}</div>
                  <div className="mt-1 text-neutral-100">{selectedEdge.packets || 0}</div>
                </div>
                <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                  <div className="text-neutral-500">Proto</div>
                  <div className="mt-1 truncate text-neutral-100">
                    {selectedEdge.dominant_proto || selectedEdge.dominant_app || '--'}
                  </div>
                </div>
              </div>
              <ProcessList processes={selectedEdge.processes} />
            </>
          ) : null}
        </div>
        <div className="rounded-xl border border-neutral-600 bg-black/20 p-3">
          <div className="text-xs font-mono text-neutral-300">
            {t('admin.contests.trafficGraph.panel.selectedNode')}
          </div>
          {selectedNode ? (
            <>
              <div className="mt-1 truncate text-[11px] text-white">{selectedNode.label || selectedNode.ip}</div>
              <div className="mt-1 truncate text-[11px] text-neutral-500">
                {selectedNode.service ? `${selectedNode.service} / ${selectedNode.ip}` : selectedNode.ip}
              </div>
              {(selectedNode.services || []).length > 1 ? (
                <div className="mt-2 flex flex-wrap gap-1.5">
                  {(selectedNode.services || []).slice(0, 4).map((service) => (
                    <Chip
                      key={service}
                      label={ellipsis(service, 18)}
                      variant="tag"
                      size="sm"
                      colorClass="border-geek-400/30 bg-geek-400/10 text-geek-400"
                    />
                  ))}
                </div>
              ) : null}
              <div className="mt-3 grid grid-cols-3 gap-2 text-[11px] font-mono">
                <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                  <div className="text-neutral-500">{t('admin.contests.trafficGraph.panel.nodeTraffic')}</div>
                  <div className="mt-1 text-geek-400">{formatBytes(selectedNode.bytes)}</div>
                </div>
                <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                  <div className="text-neutral-500">{t('admin.contests.trafficGraph.panel.nodeConnections')}</div>
                  <div className="mt-1 text-neutral-100">{selectedNode.connections || 0}</div>
                </div>
                <div className="rounded-lg border border-neutral-700 bg-black/20 px-2 py-1.5">
                  <div className="text-neutral-500">Proto</div>
                  <div className="mt-1 truncate text-neutral-100">
                    {(selectedNode.protocols || []).slice(0, 2).join(' / ') || '--'}
                  </div>
                </div>
              </div>
              <ProcessList processes={selectedNode.processes} />
            </>
          ) : null}
        </div>
      </div>
    </Card>
  );
}
