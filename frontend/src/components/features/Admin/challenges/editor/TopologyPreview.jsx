import { useTranslation } from 'react-i18next';

export default function TopologyPreview({ topology }) {
  const { t } = useTranslation();
  if (topology.nodes.length === 0) {
    return (
      <div className="rounded-md border border-neutral-700 bg-black/20 p-4 text-center font-mono text-sm text-neutral-500">
        {t('admin.challengeModal.topology.empty')}
      </div>
    );
  }
  return (
    <div className="rounded-md border border-neutral-700 bg-black/20 p-3">
      <div className="mb-3 flex items-center justify-between gap-3">
        <div>
          <div className="text-sm font-mono text-neutral-100">{t('admin.challengeModal.topology.title')}</div>
          <div className="text-xs font-mono text-neutral-500">{t('admin.challengeModal.topology.subtitle')}</div>
        </div>
        <div className="flex items-center gap-3 text-xs font-mono text-neutral-400">
          <span className="inline-flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-geek-400" /> {t('admin.challengeModal.topology.allow')}
          </span>
          <span className="inline-flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-red-400" /> {t('admin.challengeModal.topology.deny')}
          </span>
        </div>
      </div>
      <div className="relative h-[420px] overflow-hidden rounded-md border border-neutral-800 bg-black/20">
        <svg className="absolute inset-0 h-full w-full" viewBox="0 0 100 100" preserveAspectRatio="none">
          <defs>
            <marker
              id="topology-arrow-allow"
              viewBox="0 0 10 10"
              refX="8"
              refY="5"
              markerWidth="3"
              markerHeight="3"
              orient="auto-start-reverse"
            >
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#22c55e" />
            </marker>
            <marker
              id="topology-arrow-deny"
              viewBox="0 0 10 10"
              refX="8"
              refY="5"
              markerWidth="3"
              markerHeight="3"
              orient="auto-start-reverse"
            >
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#f87171" />
            </marker>
          </defs>
          {topology.connections.map((connection, index) => {
            const offset = (index % 2 === 0 ? 1 : -1) * 2.5;
            const midX = (connection.source.x + connection.target.x) / 2 + offset;
            const midY = (connection.source.y + connection.target.y) / 2 - offset;
            return (
              <path
                key={connection.id}
                d={`M ${connection.source.x} ${connection.source.y} Q ${midX} ${midY} ${connection.target.x} ${connection.target.y}`}
                fill="none"
                stroke={connection.allowed ? '#22c55e' : '#f87171'}
                strokeWidth="0.35"
                strokeDasharray={connection.allowed ? 'none' : '1.2 1.2'}
                markerEnd={`url(#${connection.allowed ? 'topology-arrow-allow' : 'topology-arrow-deny'})`}
                opacity={connection.allowed ? 0.75 : 0.55}
              />
            );
          })}
        </svg>
        {topology.nodes.map((node) => (
          <div
            key={node.id}
            className="absolute w-40 -translate-x-1/2 -translate-y-1/2 rounded-md border border-neutral-700 bg-neutral-950/95 p-2"
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
                <div className="text-xs font-mono text-neutral-500">{t('admin.challengeModal.topology.noIp')}</div>
              )}
            </div>
          </div>
        ))}
      </div>
      <div className="mt-3 max-h-52 space-y-2 overflow-y-auto pr-1">
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
                {connection.networks.length
                  ? connection.networks.join(', ')
                  : t('admin.challengeModal.topology.crossNetwork')}
              </span>
            </div>
            <div className="mt-1 text-neutral-400">
              {t(`admin.challengeModal.topology.reasons.${connection.reasonKey}`)}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
