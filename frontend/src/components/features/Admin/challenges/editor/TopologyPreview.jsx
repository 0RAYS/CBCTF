import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import { IconWorld } from '@tabler/icons-react';
import '@/components/features/Admin/challenges/editor/TopologyPreview.css';

const tt = (key) => `admin.challengeModal.topology.${key}`;

const connectionPath = (source, target) => {
  const dx = target.x - source.x;
  const dy = target.y - source.y;
  const distance = Math.hypot(dx, dy) || 1;
  // A perpendicular bend puts reverse connections on opposite sides.
  const controlX = (source.x + target.x) / 2 - (dy / distance) * 6;
  const controlY = (source.y + target.y) / 2 + (dx / distance) * 6;
  const startControlX = (source.x + controlX) / 2;
  const startControlY = (source.y + controlY) / 2;
  const endControlX = (controlX + target.x) / 2;
  const endControlY = (controlY + target.y) / 2;
  const midX = (startControlX + endControlX) / 2;
  const midY = (startControlY + endControlY) / 2;
  // Split the curve so its midpoint can carry an arrow outside the node cards.
  return `M ${source.x} ${source.y} Q ${startControlX} ${startControlY} ${midX} ${midY} Q ${endControlX} ${endControlY} ${target.x} ${target.y}`;
};

export default function TopologyPreview({ topology }) {
  const { t } = useTranslation();
  const markerId = useId();
  if (topology.nodes.length === 0) {
    return (
      <div className="rounded-md border border-neutral-700 bg-black/20 p-4 text-center font-mono text-sm text-neutral-500">
        {t(tt('empty'))}
      </div>
    );
  }
  const cardWidth = topology.lanes.length > 2 ? 'w-28' : 'w-36';
  const laneHeight = Math.max(...topology.lanes.map((lane) => lane.nodes.length));
  const boxHeight = Math.max(460, 130 + laneHeight * 96);
  const edges = [...topology.connections, ...topology.egress];
  return (
    <div className="rounded-md border border-neutral-700 bg-black/20 p-3">
      <div className="mb-3 flex items-center justify-between gap-3">
        <div>
          <div className="text-sm font-mono text-neutral-100">{t(tt('title'))}</div>
          <div className="text-xs font-mono text-neutral-500">{t(tt('subtitle'))}</div>
        </div>
        <div className="flex items-center gap-3 text-xs font-mono text-neutral-400">
          <span className="inline-flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-geek-400" /> {t(tt('allow'))}
          </span>
          <span className="inline-flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-red-400" /> {t(tt('deny'))}
          </span>
        </div>
      </div>
      <div
        className="relative overflow-hidden rounded-md border border-neutral-800 bg-black/20"
        style={{ height: `${boxHeight}px` }}
      >
        <svg className="absolute inset-0 h-full w-full" viewBox="0 0 100 100" preserveAspectRatio="none">
          <defs>
            <marker
              id={`${markerId}-allow`}
              viewBox="0 0 16 16"
              refX="8"
              refY="8"
              markerUnits="userSpaceOnUse"
              markerWidth="4.5"
              markerHeight="4.5"
              orient="auto"
            >
              <path
                d="M 2 2 L 14 8 L 2 14 L 5 8 Z"
                fill="#4ade80"
                stroke="#0a0a0a"
                strokeWidth="1"
                strokeDasharray="none"
                strokeOpacity="1"
                strokeLinejoin="round"
              />
            </marker>
            <marker
              id={`${markerId}-deny`}
              viewBox="0 0 16 16"
              refX="8"
              refY="8"
              markerUnits="userSpaceOnUse"
              markerWidth="4.5"
              markerHeight="4.5"
              orient="auto"
            >
              <path
                d="M 4 2 H 6 V 14 H 4 Z M 10 2 H 12 V 14 H 10 Z"
                fill="#f87171"
                stroke="#0a0a0a"
                strokeWidth="1"
                strokeDasharray="none"
                strokeOpacity="1"
                strokeLinejoin="round"
              />
            </marker>
          </defs>
          {topology.lanes.map((lane) => (
            <g key={lane.id}>
              <rect
                x={lane.x - lane.width / 2 + 0.6}
                y={15}
                width={lane.width - 1.2}
                height={82}
                rx={2}
                fill="#ffffff"
                fillOpacity="0.02"
                stroke="#404040"
                strokeOpacity="0.5"
                strokeWidth="0.2"
                strokeDasharray="1.5 1.5"
              />
              <text
                x={lane.x}
                y={12.6}
                textAnchor="middle"
                fill="#a3a3a3"
                fontSize="2.6"
                fontFamily="monospace"
                stroke="#0a0a0a"
                strokeWidth="0.7"
                paintOrder="stroke"
              >
                {lane.name}
              </text>
            </g>
          ))}
          {edges.map((connection) => {
            const { source, target } = connection;
            const path = connectionPath(source, target);
            return (
              <g
                key={connection.id}
                className={`topology-connection ${connection.allowed ? 'topology-connection-allow' : 'topology-connection-deny'}`}
                fill="none"
                pointerEvents="none"
              >
                <path d={path} stroke="transparent" strokeWidth="2.5" pointerEvents="stroke" />
                <g className="topology-connection-visual">
                  <path
                    className="topology-connection-line"
                    d={path}
                    stroke={connection.allowed ? '#22c55e' : '#f87171'}
                    strokeWidth="0.35"
                    strokeDasharray={connection.allowed ? 'none' : '1.2 1.2'}
                    strokeOpacity={connection.allowed ? 0.75 : 0.65}
                  />
                  {connection.allowed && (
                    <path
                      className="topology-connection-flow"
                      d={path}
                      stroke="#bbf7d0"
                      strokeWidth="0.6"
                      strokeDasharray="0.9 1.5"
                      strokeLinecap="round"
                    />
                  )}
                  <path d={path} markerMid={`url(#${markerId}-${connection.allowed ? 'allow' : 'deny'})`} />
                </g>
              </g>
            );
          })}
        </svg>
        <div
          className="absolute flex -translate-x-1/2 -translate-y-1/2 items-center gap-1 rounded-md border border-neutral-600 bg-neutral-900/95 px-2 py-1"
          style={{ left: `${topology.internet.x}%`, top: `${topology.internet.y}%` }}
        >
          <IconWorld size={14} className="text-geek-400" />
          <span className="text-xs font-mono text-neutral-100">{t(tt('internet'))}</span>
        </div>
        {topology.nodes.map((node) => (
          <div
            key={node.id}
            className={`absolute ${cardWidth} -translate-x-1/2 -translate-y-1/2 rounded-md border border-neutral-700 bg-neutral-900/95 p-2`}
            style={{ left: `${node.x}%`, top: `${node.y}%` }}
          >
            <div className="truncate text-sm font-mono text-neutral-50" title={node.label}>
              {node.label}
            </div>
            {node.image && (
              <div className="mt-0.5 truncate text-[10px] font-mono text-neutral-500" title={node.image}>
                {node.image}
              </div>
            )}
            <div className="mt-2 space-y-1">
              {node.networks.length > 0 ? (
                node.networks.map((network, index) => (
                  <div
                    key={`${network.name}-${index}`}
                    className="rounded border border-neutral-700 bg-black/40 px-1.5 py-1"
                  >
                    <div className="truncate text-[10px] font-mono text-neutral-500" title={network.name}>
                      {network.name}
                    </div>
                    <div className="truncate text-xs font-mono text-geek-300" title={network.ip}>
                      {network.ip}
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-xs font-mono text-neutral-500">{t(tt('noIp'))}</div>
              )}
            </div>
          </div>
        ))}
      </div>
      <div className="mt-3 max-h-60 space-y-2 overflow-y-auto pr-1">
        {topology.egress.map((connection) => (
          <div
            key={connection.id}
            className={`rounded border px-2 py-1.5 font-mono text-xs ${connection.allowed ? 'border-geek-400/25 bg-geek-400/10 text-geek-200' : 'border-red-400/25 bg-red-400/10 text-red-200'}`}
          >
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-neutral-200">{connection.source.label}</span>
              <span>→</span>
              <span className="text-neutral-200">{t(tt('internet'))}</span>
              {connection.ranges.length > 0 && <span className="text-neutral-500">{connection.ranges.join(', ')}</span>}
            </div>
            <div className="mt-1 text-neutral-400">{t(tt(`reasons.${connection.reasonKey}`))}</div>
          </div>
        ))}
        {topology.connections.map((connection) => (
          <div
            key={connection.id}
            className={`rounded border px-2 py-1.5 font-mono text-xs ${connection.allowed ? 'border-geek-400/25 bg-geek-400/10 text-geek-200' : 'border-red-400/25 bg-red-400/10 text-red-200'}`}
          >
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-neutral-200">{connection.source.label}</span>
              <span>{connection.allowed ? '->' : '-/->'}</span>
              <span className="text-neutral-200">{connection.target.label}</span>
              <span className="text-neutral-500">
                {connection.networks.length ? connection.networks.join(', ') : t(tt('crossNetwork'))}
              </span>
            </div>
            <div className="mt-1 text-neutral-400">{t(tt(`reasons.${connection.reasonKey}`))}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
