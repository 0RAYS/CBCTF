import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Card } from '../../../common';
import Checkbox from '../../../common/Checkbox.jsx';
import { ANALYSIS_TABS, evidenceOffset, filterAnalysisRows, victimOverlaps } from './trafficAnalysisModel.js';
import { formatBytes } from './trafficPresentation.js';

const PREFIX = 'admin.contests.trafficGraph.analysis';
const PAGE_SIZE = 20;

export default function TrafficAnalysisPanel({
  analysis,
  topology,
  victimId,
  contestId,
  selectedNodeId,
  selectedEdge,
  onSeek,
}) {
  const { t } = useTranslation();
  const [tab, setTab] = useState('flags');
  const [currentWindow, setCurrentWindow] = useState(false);
  const [selectedOnly, setSelectedOnly] = useState(false);
  const [page, setPage] = useState(0);
  const report = analysis.data?.report;
  const overlaps = useMemo(() => victimOverlaps(analysis.overlaps || [], victimId), [analysis.overlaps, victimId]);
  const source = tab === 'accesses' ? analysis.data?.accesses : tab === 'overlaps' ? overlaps : report?.[tab];
  const rows = useMemo(
    () =>
      tab === 'overlaps'
        ? source || []
        : filterAnalysisRows(source || [], {
            startedAt: topology?.started_at,
            window: topology?.window || { start: 0, end: 0 },
            currentWindow,
            nodeId: selectedOnly ? selectedNodeId : '',
            edge: selectedOnly && !selectedNodeId ? selectedEdge : null,
          }),
    [source, tab, topology, currentWindow, selectedOnly, selectedNodeId, selectedEdge]
  );
  useEffect(() => setPage(0), [tab, currentWindow, selectedOnly, selectedNodeId, selectedEdge?.id, victimId]);
  const lastPage = Math.max(0, Math.ceil(rows.length / PAGE_SIZE) - 1);
  const currentPage = Math.min(page, lastPage);
  const visible = rows.slice(currentPage * PAGE_SIZE, (currentPage + 1) * PAGE_SIZE);
  const label = (key, options) => t(`${PREFIX}.${key}`, options);

  const renderDetails = (row) => {
    switch (tab) {
      case 'flags':
        return (
          <>
            <code className="block whitespace-pre-wrap break-all text-geek-300">{row.value}</code>
            <span>
              {label(row.verified ? 'verified' : 'candidate')} · {row.encoding}
            </span>
          </>
        );
      case 'indicators':
        return (
          <>
            <strong>{label(`rules.${row.rule}`, { defaultValue: row.rule })}</strong>
            <span className="block break-all">{row.detail}</span>
          </>
        );
      case 'http':
        return (
          <>
            <strong>
              {row.method || row.status} {row.host}
            </strong>
            <span className="block break-all">{row.uri || row.content_type}</span>
          </>
        );
      case 'dns':
        return (
          <>
            <code className="break-all">{row.name}</code>
            <span>
              {' '}
              · {row.type} · {label(row.response ? 'response' : 'request')}
            </span>
          </>
        );
      case 'sessions':
        return (
          <>
            {formatBytes(row.bytes)} · {label('packets', { count: row.packets })} · SYN {row.syn} / RST {row.rst}
          </>
        );
      case 'accesses':
        return (
          <>
            <code>{row.ip}</code> · {label(`sources.${row.source}`, { defaultValue: row.source })}
          </>
        );
      case 'overlaps':
        return (
          <>
            <code>{row.ip}</code>
            <div className="mt-1 grid gap-1">
              {(row.teams || []).map((team) => (
                <span key={team.team_id}>
                  {label('team', { id: team.team_id })} · {label('victims', { ids: team.victim_ids?.join(', ') })}
                  <time className="block">{team.first_seen}</time>
                </span>
              ))}
            </div>
          </>
        );
      default:
        return null;
    }
  };

  return (
    <Card className="min-w-0 rounded-2xl border-neutral-600 bg-neutral-900" padding="sm">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 className="font-semibold">{label('title')}</h3>
          <p className="mt-1 text-xs text-neutral-400">{label('scope')}</p>
        </div>
        <Button size="sm" onClick={analysis.refresh} disabled={analysis.loading || analysis.overlapsLoading}>
          {label('refresh')}
        </Button>
      </div>
      <p className="mt-2 text-xs text-neutral-400">{label('hints')}</p>
      {analysis.loading ? (
        <p role="status" className="mt-3">
          {label('loading')}
        </p>
      ) : null}
      {analysis.error ? (
        <p role="alert" className="mt-3 text-red-300">
          {label('error')}
        </p>
      ) : null}
      {report ? (
        <>
          <div className="mt-3 flex flex-wrap gap-3 text-xs text-neutral-400">
            <span>{label(analysis.data.archived ? 'archived' : 'live')}</span>
            <span>
              {label('packets', { count: report.packets })} · {formatBytes(report.bytes)}
            </span>
            {report.truncated ? <span className="text-amber-300">{label('truncated')}</span> : null}
          </div>
          {report.warnings?.length ? (
            <details className="mt-3 rounded-lg border border-amber-500/30 p-2 text-xs text-amber-200">
              <summary className="cursor-pointer">{label('warnings', { count: report.warnings.length })}</summary>
              <ul className="mt-2 list-inside list-disc">
                {report.warnings.map((warning) => (
                  <li key={warning} className="break-all">
                    {label(`warningCodes.${warning}`, {
                      defaultValue: warning,
                    })}
                  </li>
                ))}
              </ul>
            </details>
          ) : null}
          <div className="mt-3 flex flex-wrap gap-2" role="tablist" aria-label={label('title')}>
            {ANALYSIS_TABS.filter((key) => key !== 'overlaps' || contestId).map((key) => (
              <button
                key={key}
                type="button"
                role="tab"
                aria-selected={tab === key}
                onClick={() => setTab(key)}
                className={`rounded-lg border px-3 py-2 text-xs ${tab === key ? 'border-geek-400 text-geek-300' : 'border-neutral-600 text-neutral-300'}`}
              >
                {label(`tabs.${key}`)} (
                {key === 'accesses'
                  ? analysis.data.accesses?.length || 0
                  : key === 'overlaps'
                    ? overlaps.length
                    : report[key]?.length || 0}
                )
              </button>
            ))}
          </div>
          {tab !== 'overlaps' ? (
            <div className="mt-3 flex flex-wrap gap-4 text-xs">
              <Checkbox
                label={label('currentWindow')}
                checked={currentWindow}
                disabled={!topology?.started_at}
                onChange={(event) => setCurrentWindow(event.target.checked)}
              />
              <Checkbox
                label={label('selectedOnly')}
                checked={selectedOnly}
                disabled={!selectedNodeId && !selectedEdge}
                onChange={(event) => setSelectedOnly(event.target.checked)}
              />
            </div>
          ) : (
            <p className="mt-3 text-xs text-amber-200">{label('overlapHint')}</p>
          )}
          {tab === 'overlaps' && analysis.overlapsLoading ? <p role="status">{label('loading')}</p> : null}
          {tab === 'overlaps' && analysis.overlapsError ? (
            <p role="alert" className="mt-2 text-red-300">
              {label('overlapError')}
            </p>
          ) : null}
          <div
            className="mt-3 grid max-h-[28rem] gap-2 overflow-y-auto"
            role="tabpanel"
            aria-label={label(`tabs.${tab}`)}
          >
            {visible.map((row, index) => {
              const evidence = row.evidence || row;
              const offset = evidenceOffset(evidence, topology?.started_at);
              return (
                <article
                  key={`${tab}-${currentPage}-${index}`}
                  className="min-w-0 rounded-lg border border-neutral-700 bg-black/20 p-3 text-xs"
                >
                  <div className="break-words">{renderDetails(row)}</div>
                  {tab !== 'overlaps' ? (
                    <div className="mt-2 flex flex-wrap items-end justify-between gap-2 text-neutral-400">
                      <div className="min-w-0 break-all">
                        {evidence.src_ip ? (
                          <div>
                            {evidence.src_ip}:{evidence.src_port} → {evidence.dst_ip}:{evidence.dst_port} ·{' '}
                            {evidence.protocol}
                          </div>
                        ) : null}
                        <div>{evidence.capture}</div>
                        <time>{evidence.time}</time>
                        {evidence.end_time && evidence.end_time !== evidence.time ? (
                          <time className="block">→ {evidence.end_time}</time>
                        ) : null}
                      </div>
                      <Button
                        size="sm"
                        disabled={offset === null || offset < 0 || offset >= (topology?.total_duration || 0)}
                        onClick={() => onSeek(offset)}
                      >
                        {label('seek')}
                      </Button>
                    </div>
                  ) : null}
                </article>
              );
            })}
            {!rows.length && !analysis.overlapsLoading && !(tab === 'overlaps' && analysis.overlapsError) ? (
              <p className="py-4 text-sm text-neutral-400">{label('empty')}</p>
            ) : null}
          </div>
          <div className="mt-3 flex items-center justify-end gap-3 text-xs">
            <Button size="sm" disabled={currentPage === 0} onClick={() => setPage(currentPage - 1)}>
              {label('previous')}
            </Button>
            <span>
              {currentPage + 1} / {lastPage + 1} · {rows.length}
            </span>
            <Button size="sm" disabled={currentPage === lastPage} onClick={() => setPage(currentPage + 1)}>
              {label('next')}
            </Button>
          </div>
        </>
      ) : null}
    </Card>
  );
}
