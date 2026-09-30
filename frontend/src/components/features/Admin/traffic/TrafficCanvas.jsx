import { useMemo } from 'react';
import { IconZoomIn, IconZoomOut } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Card, Chip } from '../../../common';
import { buildLines, buildPositions, CENTER_Y, resolveEdgeLabels, VIEWBOX_WIDTH, ZOOM_STEP } from './trafficLayout.js';
import { buildCompactNodeLabel, buildCompactNodeMeta, ellipsis } from './trafficPresentation.js';
import useTopologyViewport from './useTopologyViewport.js';

const TOPOLOGY_CANVAS_FIXED_HEIGHT = 'min(800px, calc(100vh - 48px))';

export default function TrafficCanvas({
  scopeKey,
  nodes,
  edges,
  universeNodes,
  universeProtocols,
  activeProtocolSet,
  protocolFilter,
  setProtocolFilter,
  selectedEdgeId,
  selectedEdge,
  selectedNode,
  onSelectEdge,
  onSelectNode,
}) {
  const { t } = useTranslation();
  const { canvasRef, zoom, viewBox, applyZoom, resetView, startDrag, onDrag, stopDrag, canSelect } =
    useTopologyViewport(scopeKey);
  const allLayoutNodes = universeNodes ?? nodes;
  const positions = useMemo(() => buildPositions(allLayoutNodes), [allLayoutNodes]);
  const activeNodeIds = useMemo(() => new Set(nodes.map((node) => node.id)), [nodes]);
  const lines = useMemo(() => buildLines(edges, positions), [edges, positions]);
  const labelPlacements = useMemo(
    () => resolveEdgeLabels(lines, positions, selectedEdgeId),
    [lines, positions, selectedEdgeId]
  );
  const orderedLines = useMemo(
    () =>
      [...lines].sort((a, b) => {
        const aSelected = a.id === selectedEdgeId ? 1 : 0;
        const bSelected = b.id === selectedEdgeId ? 1 : 0;
        if (aSelected !== bSelected) return aSelected - bSelected;
        if ((a.bytes || 0) !== (b.bytes || 0)) return (a.bytes || 0) - (b.bytes || 0);
        return String(a.id || '').localeCompare(String(b.id || ''));
      }),
    [lines, selectedEdgeId]
  );

  const toggleProtocol = (proto) => {
    setProtocolFilter((current) => {
      const next = new Set(current);
      if (next.has(proto)) next.delete(proto);
      else next.add(proto);
      return next;
    });
  };
  const handleSelectEdge = (id) => {
    if (canSelect()) onSelectEdge(id);
  };
  const handleSelectNode = (id) => {
    if (canSelect()) onSelectNode(id);
  };

  return (
    <Card padding="none" className="overflow-hidden rounded-2xl border-neutral-600 bg-neutral-900">
      <div className="flex flex-col">
        <div className="flex flex-wrap items-start justify-between gap-2 border-b border-neutral-600 px-4 py-3">
          <div>
            <div className="text-sm font-mono text-neutral-300">{t('admin.contests.trafficGraph.canvas.title')}</div>
            <div className="mt-1 text-[11px] text-neutral-500">{t('admin.contests.trafficGraph.canvas.subtitle')}</div>
          </div>
          <div className="flex flex-wrap items-center justify-end gap-2">
            <Chip
              label={t('admin.contests.trafficGraph.legend.ingress')}
              variant="tag"
              size="sm"
              colorClass="border-geek-400/30 bg-geek-400/10 text-geek-400"
            />
            <Chip
              label={t('admin.contests.trafficGraph.legend.egress')}
              variant="tag"
              size="sm"
              colorClass="border-neutral-400/30 bg-neutral-400/10 text-neutral-300"
            />
            <Chip
              label={t('admin.contests.trafficGraph.legend.internal')}
              variant="tag"
              size="sm"
              colorClass="border-neutral-500/30 bg-neutral-500/10 text-neutral-300"
            />
            <Chip
              label={`x${zoom.toFixed(2)}`}
              variant="tag"
              size="sm"
              colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
            />
            <Button
              variant="ghost"
              size="icon"
              className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
              onClick={() => applyZoom(zoom - ZOOM_STEP)}
            >
              <IconZoomOut size={16} />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="!h-8 !w-8 !text-neutral-300 hover:!text-neutral-100"
              onClick={() => applyZoom(zoom + ZOOM_STEP)}
            >
              <IconZoomIn size={16} />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="!h-8 !px-3 !text-neutral-300 hover:!text-neutral-100"
              onClick={resetView}
            >
              {t('common.reset')}
            </Button>
          </div>
        </div>
        {universeProtocols.length > 0 ? (
          <div className="flex flex-wrap items-center gap-2 border-b border-neutral-700/80 px-4 py-2">
            <span className="text-[11px] text-neutral-400 shrink-0">
              {t('admin.contests.trafficGraph.filter.label')}
            </span>
            {universeProtocols.map((proto) => {
              const active = protocolFilter.has(proto);
              const inCurrentFrame = activeProtocolSet.has(proto);
              return (
                <button
                  key={proto}
                  type="button"
                  onClick={() => toggleProtocol(proto)}
                  className={[
                    'inline-flex items-center rounded border px-2 py-0.5 text-[11px] font-mono transition-colors',
                    active
                      ? 'border-geek-400/60 bg-geek-400/20 text-geek-300'
                      : inCurrentFrame
                        ? 'border-neutral-600 bg-black/20 text-neutral-400 hover:border-neutral-500 hover:text-neutral-300'
                        : 'border-neutral-600/60 bg-black/15 text-neutral-500 hover:border-neutral-500 hover:text-neutral-400',
                  ].join(' ')}
                >
                  {proto}
                </button>
              );
            })}
            {protocolFilter.size > 0 ? (
              <button
                type="button"
                onClick={() => setProtocolFilter(new Set())}
                className="ml-1 text-[11px] text-neutral-500 hover:text-neutral-300 transition-colors"
              >
                {t('admin.contests.trafficGraph.filter.clear')}
              </button>
            ) : null}
          </div>
        ) : null}
        <div
          ref={canvasRef}
          style={{ height: TOPOLOGY_CANVAS_FIXED_HEIGHT }}
          className="relative overflow-hidden bg-[linear-gradient(180deg,rgba(0,0,0,.14),rgba(0,0,0,.05))] cursor-grab active:cursor-grabbing"
          onPointerMove={onDrag}
          onPointerUp={stopDrag}
          onPointerLeave={stopDrag}
        >
          <div className="pointer-events-none absolute inset-0 opacity-20 [background-image:linear-gradient(rgba(82,82,82,.45)_1px,transparent_1px),linear-gradient(90deg,rgba(82,82,82,.45)_1px,transparent_1px)] [background-size:34px_34px]" />
          <div className="absolute left-3 top-3 z-10 text-[11px] text-neutral-400">
            {t('admin.contests.trafficGraph.canvas.zoomHint')}
          </div>
          <svg
            viewBox={viewBox}
            className="absolute inset-0 h-full w-full select-none touch-none"
            shapeRendering="geometricPrecision"
            textRendering="geometricPrecision"
            onPointerDown={startDrag}
          >
            <g opacity="0.18">
              <path
                d={`M 34 ${CENTER_Y - 132} H ${VIEWBOX_WIDTH - 34}`}
                stroke="#404040"
                strokeWidth="1"
                vectorEffect="non-scaling-stroke"
              />
              <path
                d={`M 34 ${CENTER_Y} H ${VIEWBOX_WIDTH - 34}`}
                stroke="#404040"
                strokeWidth="1"
                vectorEffect="non-scaling-stroke"
              />
              <path
                d={`M 34 ${CENTER_Y + 132} H ${VIEWBOX_WIDTH - 34}`}
                stroke="#404040"
                strokeWidth="1"
                vectorEffect="non-scaling-stroke"
              />
            </g>
            {orderedLines.map((edge) => {
              const label = labelPlacements.get(edge.id);
              const selected = selectedEdge?.id === edge.id;
              return (
                <g key={edge.id}>
                  <path
                    d={edge.path}
                    fill="none"
                    stroke={edge.tone.hard}
                    strokeWidth={selected ? 8.5 : 7}
                    opacity={selected ? 0.22 : 0.08 + edge.intensity * 0.1}
                    vectorEffect="non-scaling-stroke"
                  />
                  <path
                    d={edge.path}
                    fill="none"
                    stroke={edge.tone.soft}
                    strokeWidth={2.2 + edge.intensity * 3.2}
                    className="traffic-line"
                    opacity={selected ? 0.95 : 0.78}
                    onClick={() => handleSelectEdge(edge.id)}
                    vectorEffect="non-scaling-stroke"
                  />
                  {label ? (
                    <g
                      transform={`translate(${label.x + label.w / 2}, ${label.y + label.h / 2})`}
                      onClick={() => handleSelectEdge(edge.id)}
                    >
                      <rect
                        x={-label.w / 2}
                        y={-label.h / 2}
                        width={label.w}
                        height={label.h}
                        rx={10}
                        fill="rgba(10,10,10,.95)"
                        stroke={selected ? edge.tone.hard : 'rgba(115,115,115,.35)'}
                        vectorEffect="non-scaling-stroke"
                      />
                      <text x="0" y="4" textAnchor="middle" fill="#d4d4d4" fontSize="10.2" fontFamily="Maple Mono">
                        {label.text}
                      </text>
                    </g>
                  ) : null}
                </g>
              );
            })}
            {allLayoutNodes.map((node) => {
              const pos = positions.get(node.id);
              if (!pos) return null;
              const isGhost = !activeNodeIds.has(node.id);
              const active = !isGhost && selectedNode?.id === node.id;
              const isVictim = node.kind === 'victim';
              const x = pos.x - pos.w / 2;
              const y = pos.y - pos.h / 2;
              const title = buildCompactNodeLabel(node);
              const meta = buildCompactNodeMeta(node);
              return (
                <g
                  key={node.id}
                  transform={`translate(${x}, ${y})`}
                  onClick={isGhost ? undefined : () => handleSelectNode(node.id)}
                  style={{ cursor: isGhost ? 'default' : 'pointer' }}
                  opacity={isGhost ? 0.45 : 1}
                >
                  <rect
                    x={active ? -2.5 : 0}
                    y={active ? -2.5 : 0}
                    width={active ? pos.w + 5 : pos.w}
                    height={active ? pos.h + 5 : pos.h}
                    rx="15"
                    fill={isGhost ? 'rgba(24,24,27,0.7)' : isVictim ? 'rgba(89,126,247,.12)' : 'rgba(24,24,27,.94)'}
                    stroke={active ? '#f5f5f5' : isGhost ? 'rgba(100,100,100,0.55)' : isVictim ? '#597ef7' : '#666666'}
                    strokeWidth={active ? 1.7 : 1.1}
                    vectorEffect="non-scaling-stroke"
                  />
                  <circle
                    cx="18"
                    cy={pos.h / 2}
                    r={isVictim ? 9.5 : 8}
                    fill={isGhost ? 'rgba(55,55,62,0.7)' : isVictim ? 'rgba(89,126,247,.14)' : 'rgba(82,82,91,.62)'}
                    stroke={isGhost ? 'rgba(130,130,130,0.5)' : isVictim ? '#597ef7' : '#8a8a8a'}
                    vectorEffect="non-scaling-stroke"
                  />
                  <text
                    x="34"
                    y={pos.h / 2 - 4}
                    fill={isGhost ? '#707070' : '#f5f5f5'}
                    fontSize="10.7"
                    fontFamily="Maple Mono"
                  >
                    {title}
                  </text>
                  <text
                    x="34"
                    y={pos.h / 2 + 11}
                    fill={isGhost ? '#565656' : '#a3a3a3'}
                    fontSize="9.4"
                    fontFamily="Maple Mono"
                  >
                    {ellipsis(meta, 20)}
                  </text>
                </g>
              );
            })}
          </svg>
        </div>
      </div>
    </Card>
  );
}
