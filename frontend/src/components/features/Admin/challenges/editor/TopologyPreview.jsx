import { useEffect, useId, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconWorld } from '@tabler/icons-react';
import { layoutTopology } from '@/components/features/Admin/challenges/editor/networkTopologyLayout.js';
import '@/components/features/Admin/challenges/editor/TopologyPreview.css';

const tt = (key) => `admin.challengeModal.topology.${key}`;

export default function TopologyPreview({ topology }) {
  const { t } = useTranslation();
  const markerId = useId();
  const frameRef = useRef(null);
  const [width, setWidth] = useState(0);
  const hasNodes = topology.nodes.length > 0;

  useEffect(() => {
    const frame = frameRef.current;
    if (!frame) return undefined;
    const observer = new ResizeObserver((entries) => {
      const measured = Math.round(entries[0].contentRect.width);
      setWidth((current) => (current === measured ? current : measured));
    });
    observer.observe(frame);
    return () => observer.disconnect();
  }, [hasNodes]);

  const layout = useMemo(() => layoutTopology(topology, { width }), [topology, width]);

  if (!hasNodes) {
    return (
      <div className="rounded-md border border-neutral-700 bg-black/20 p-4 text-center font-mono text-sm text-neutral-500">
        {t(tt('empty'))}
      </div>
    );
  }

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
          {topology.routes?.length > 0 && (
            <span className="inline-flex items-center gap-1">
              <span className="w-3 border-t-2 border-dashed border-purple-400" /> {t(tt('route'))}
            </span>
          )}
        </div>
      </div>
      {topology.lanes.length > 0 && (
        <div className="mb-3 flex flex-wrap gap-x-4 gap-y-1">
          {topology.lanes.map((lane) => (
            <span key={lane.id} className="inline-flex items-center gap-1.5 text-[11px] font-mono">
              <span className="h-2 w-2 rounded-full" style={{ backgroundColor: lane.color }} />
              <span className="text-neutral-300">{lane.name}</span>
              <span className="text-neutral-600">
                {lane.members
                  .map((id) => topology.nodes.find((node) => node.id === id)?.label)
                  .filter(Boolean)
                  .join(', ')}
              </span>
            </span>
          ))}
        </div>
      )}
      <div ref={frameRef} className="overflow-x-auto rounded-md border border-neutral-800 bg-black/20">
        {width > 0 ? (
          <div className="relative" style={{ width: layout.width, height: layout.height }}>
            <svg
              className="absolute left-0 top-0"
              width={layout.width}
              height={layout.height}
              viewBox={`0 0 ${layout.width} ${layout.height}`}
              fill="none"
            >
              <defs>
                <marker
                  id={`${markerId}-allow`}
                  viewBox="0 0 12 12"
                  refX="6"
                  refY="6"
                  markerUnits="userSpaceOnUse"
                  markerWidth="12"
                  markerHeight="12"
                  orient="auto-start-reverse"
                >
                  <path
                    d="M 1 1 L 11 6 L 1 11 L 4 6 Z"
                    fill="#4ade80"
                    stroke="#0a0a0a"
                    strokeWidth="1"
                    strokeLinejoin="round"
                  />
                </marker>
                <marker
                  id={`${markerId}-deny`}
                  viewBox="0 0 12 12"
                  refX="6"
                  refY="6"
                  markerUnits="userSpaceOnUse"
                  markerWidth="12"
                  markerHeight="12"
                  orient="auto"
                >
                  <path
                    d="M 3.5 1.5 H 5.5 V 10.5 H 3.5 Z M 8.5 1.5 H 10.5 V 10.5 H 8.5 Z"
                    fill="#f87171"
                    stroke="#0a0a0a"
                    strokeWidth="0.8"
                  />
                </marker>
              </defs>
              {layout.edges.map((edge) => (
                <g
                  key={edge.id}
                  className={`topology-connection ${edge.routed ? 'topology-connection-route' : edge.allowed ? 'topology-connection-allow' : 'topology-connection-deny'}`}
                  fill="none"
                  pointerEvents="none"
                >
                  <path d={edge.d} stroke="transparent" strokeWidth="12" pointerEvents="stroke" />
                  <g className="topology-connection-visual">
                    <path
                      className="topology-connection-line"
                      d={edge.d}
                      stroke={edge.routed ? '#c084fc' : edge.allowed ? '#22c55e' : '#f87171'}
                      strokeWidth={edge.routed ? 1.6 : edge.allowed ? 1.6 : 0.9}
                      strokeDasharray={edge.routed ? '7 5' : edge.allowed ? 'none' : '4 4'}
                      strokeOpacity={edge.routed ? 0.85 : edge.allowed ? 0.8 : 0.32}
                    />
                    {edge.allowed && !edge.routed && (
                      <path
                        className="topology-connection-flow"
                        d={edge.d}
                        stroke="#bbf7d0"
                        strokeWidth="2.4"
                        strokeDasharray="3 6"
                        strokeLinecap="round"
                      />
                    )}
                    {!edge.routed && (
                      <path
                        d={edge.d}
                        markerMid={`url(#${markerId}-${edge.allowed ? 'allow' : 'deny'})`}
                        markerStart={edge.bidirectional ? `url(#${markerId}-allow)` : undefined}
                      />
                    )}
                  </g>
                </g>
              ))}
            </svg>
            <div
              className="absolute flex items-center gap-1 rounded-md border border-neutral-600 bg-neutral-900/95 px-2 py-1"
              style={{
                left: layout.internet.x - layout.internet.width / 2,
                top: layout.internet.y - layout.internet.height / 2,
                width: layout.internet.width,
                height: layout.internet.height,
              }}
            >
              <IconWorld size={14} className="shrink-0 text-geek-400" />
              <span className="truncate text-xs font-mono text-neutral-100">{t(tt('internet'))}</span>
            </div>
            {topology.nodes.map((node) => {
              const position = layout.positions[node.id];
              if (!position) return null;
              return (
                <div
                  key={node.id}
                  className="absolute rounded-md border border-neutral-700 bg-neutral-900/95 p-2"
                  style={{
                    left: position.x - position.width / 2,
                    top: position.y - position.height / 2,
                    width: position.width,
                  }}
                >
                  <div className="flex items-center gap-1.5">
                    <div className="truncate text-sm font-mono text-neutral-50" title={node.label}>
                      {node.label}
                    </div>
                    {node.role && (
                      <span className="shrink-0 rounded border border-geek-400/40 px-1 text-[10px] font-mono text-geek-300">
                        {t(`admin.challengeModal.composeGuide.range.roles.${node.role}`)}
                      </span>
                    )}
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
                          <div className="flex items-center gap-1">
                            <span
                              className="h-1.5 w-1.5 shrink-0 rounded-full"
                              style={{ backgroundColor: network.color }}
                            />
                            <span className="truncate text-[10px] font-mono text-neutral-500" title={network.name}>
                              {network.name}
                            </span>
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
              );
            })}
          </div>
        ) : (
          <div style={{ height: 320 }} />
        )}
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
        {(topology.routes || []).map((route) => (
          <div
            key={route.id}
            className="rounded border border-purple-400/25 bg-purple-400/10 px-2 py-1.5 font-mono text-xs text-purple-200"
          >
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-neutral-200">{route.from}</span>
              <span>→</span>
              <span className="text-neutral-200">{route.to}</span>
              <span className="text-neutral-500">{t(tt('route'))}</span>
            </div>
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
