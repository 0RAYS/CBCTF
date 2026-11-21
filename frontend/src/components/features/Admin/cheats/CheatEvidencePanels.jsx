import { useTranslation } from 'react-i18next';
import { formatDateTime } from '@/components/features/Admin/cheats/cheatFormat';

function EvidenceTable({ columns, rows, empty }) {
  if (!rows || rows.length === 0) {
    return <div className="text-neutral-500 text-xs">{empty}</div>;
  }
  return (
    <div className="overflow-x-auto">
      <table className="w-full table-fixed text-xs">
        <thead>
          <tr className="bg-black/40">
            {columns.map((column) => (
              <th key={column.key} scope="col" className="p-2 text-left text-neutral-400 font-mono">
                {column.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={index} className="border-t border-neutral-300/10 align-top">
              {columns.map((column) => (
                <td key={column.key} className="p-2 text-neutral-300 font-mono break-words">
                  {row[column.key] ?? '-'}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function Section({ title, children }) {
  return (
    <div className="space-y-2">
      <div className="text-sm font-mono text-neutral-400">{title}</div>
      {children}
    </div>
  );
}

// SameWebIPPanel 展示共享网页 IP 涉及的题目、请求与提交记录。
function SameWebIPPanel({ evidence, language }) {
  const { t } = useTranslation();
  const web = evidence.web;
  if (!web) return null;
  return (
    <div className="space-y-4">
      <Section title={t('admin.contests.cheats.evidence.challenges')}>
        <EvidenceTable
          empty={t('admin.contests.cheats.evidence.empty')}
          columns={[
            { key: 'name', label: t('admin.contests.cheats.evidence.challenge') },
            { key: 'actors', label: t('admin.contests.cheats.evidence.actor') },
          ]}
          rows={(web.challenges || []).map((challenge) => ({
            name: challenge.name || `#${challenge.contest_challenge_id}`,
            actors: (challenge.actor_ids || []).map((id) => `User-${id}`).join(', ') || '-',
          }))}
        />
      </Section>
      <Section title={t('admin.contests.cheats.evidence.requests')}>
        <EvidenceTable
          empty={t('admin.contests.cheats.evidence.empty')}
          columns={[
            { key: 'at', label: t('admin.contests.cheats.evidence.at') },
            { key: 'user', label: t('admin.contests.cheats.evidence.user') },
            { key: 'method', label: t('admin.contests.cheats.evidence.method') },
            { key: 'path', label: t('admin.contests.cheats.evidence.path') },
            { key: 'status', label: t('admin.contests.cheats.evidence.status') },
            { key: 'latency', label: t('admin.contests.cheats.evidence.latency') },
          ]}
          rows={(web.requests || []).map((request) => ({
            at: formatDateTime(request.at, language),
            user: `User-${request.user_id}`,
            method: request.method,
            path: request.path,
            status: request.status,
            latency: request.latency_ms ? `${request.latency_ms}ms` : '-',
          }))}
        />
      </Section>
      <Section title={t('admin.contests.cheats.evidence.submissions')}>
        <EvidenceTable
          empty={t('admin.contests.cheats.evidence.empty')}
          columns={[
            { key: 'at', label: t('admin.contests.cheats.evidence.at') },
            { key: 'team', label: t('admin.contests.cheats.evidence.team') },
            { key: 'user', label: t('admin.contests.cheats.evidence.user') },
            { key: 'challenge', label: t('admin.contests.cheats.evidence.challenge') },
            { key: 'value', label: t('admin.contests.cheats.evidence.value') },
          ]}
          rows={(web.submissions || []).map((submission) => ({
            at: formatDateTime(submission.at, language),
            team: `Team-${submission.team_id}`,
            user: `User-${submission.user_id}`,
            challenge: submission.challenge_name || `#${submission.contest_challenge_id}`,
            value: submission.value,
          }))}
        />
      </Section>
    </div>
  );
}

// SameVictimIPPanel 展示共享靶机 IP 命中的靶机与访问记录。
function SameVictimIPPanel({ evidence, language }) {
  const { t } = useTranslation();
  const victim = evidence.victim;
  if (!victim) return null;
  return (
    <div className="space-y-4">
      <Section title={t('admin.contests.cheats.evidence.victims')}>
        <EvidenceTable
          empty={t('admin.contests.cheats.evidence.empty')}
          columns={[
            { key: 'victim', label: t('admin.contests.cheats.evidence.victim') },
            { key: 'team', label: t('admin.contests.cheats.evidence.team') },
            { key: 'challenge', label: t('admin.contests.cheats.evidence.challenge') },
            { key: 'count', label: t('admin.contests.cheats.evidence.count') },
            { key: 'window', label: t('admin.contests.cheats.evidence.window') },
          ]}
          rows={(victim.victims || []).map((item) => ({
            victim: `Victim-${item.victim_id}`,
            team: `Team-${item.team_id}`,
            challenge: item.challenge_name || '#-',
            count: item.access_count,
            window: `${formatDateTime(item.first_seen, language)} → ${formatDateTime(item.last_seen, language)}`,
          }))}
        />
      </Section>
      <Section title={t('admin.contests.cheats.evidence.accesses')}>
        <EvidenceTable
          empty={t('admin.contests.cheats.evidence.empty')}
          columns={[
            { key: 'at', label: t('admin.contests.cheats.evidence.at') },
            { key: 'team', label: t('admin.contests.cheats.evidence.team') },
            { key: 'victim', label: t('admin.contests.cheats.evidence.victim') },
            { key: 'source', label: t('admin.contests.cheats.evidence.source') },
            { key: 'capture', label: t('admin.contests.cheats.evidence.capture') },
          ]}
          rows={(victim.accesses || []).map((access) => ({
            at: formatDateTime(access.at, language),
            team: `Team-${access.team_id}`,
            victim: `Victim-${access.victim_id}`,
            source: access.source,
            capture: access.capture,
          }))}
        />
      </Section>
    </div>
  );
}

// WrongFlagPanel 展示被冒用 flag 的合法持有者与冒用提交记录。
function WrongFlagPanel({ evidence, language }) {
  const { t } = useTranslation();
  const flag = evidence.flag;
  if (!flag) return null;
  return (
    <div className="space-y-4">
      <Section title={t('admin.contests.cheats.evidence.owners')}>
        <EvidenceTable
          empty={t('admin.contests.cheats.evidence.empty')}
          columns={[
            { key: 'team', label: t('admin.contests.cheats.evidence.team') },
            { key: 'user', label: t('admin.contests.cheats.evidence.user') },
            { key: 'obtained', label: t('admin.contests.cheats.evidence.obtainedAt') },
          ]}
          rows={(flag.owners || []).map((owner) => ({
            team: `Team-${owner.team_id}`,
            user: owner.user_id ? `User-${owner.user_id}` : '-',
            obtained: owner.obtained_at ? formatDateTime(owner.obtained_at, language) : '-',
          }))}
        />
      </Section>
      <Section title={t('admin.contests.cheats.evidence.submissions')}>
        <EvidenceTable
          empty={t('admin.contests.cheats.evidence.empty')}
          columns={[
            { key: 'at', label: t('admin.contests.cheats.evidence.at') },
            { key: 'team', label: t('admin.contests.cheats.evidence.team') },
            { key: 'user', label: t('admin.contests.cheats.evidence.user') },
            { key: 'ip', label: t('admin.contests.cheats.evidence.ip') },
            { key: 'value', label: t('admin.contests.cheats.evidence.value') },
          ]}
          rows={(flag.submissions || []).map((submission) => ({
            at: formatDateTime(submission.at, language),
            team: `Team-${submission.team_id}`,
            user: `User-${submission.user_id}`,
            ip: submission.ip,
            value: submission.value,
          }))}
        />
      </Section>
    </div>
  );
}

// CheatEvidencePanel 按判定类型分发到对应的证据面板。
export default function CheatEvidencePanel({ evidence, language }) {
  if (!evidence) return null;
  if (evidence.web) return <SameWebIPPanel evidence={evidence} language={language} />;
  if (evidence.victim) return <SameVictimIPPanel evidence={evidence} language={language} />;
  if (evidence.flag) return <WrongFlagPanel evidence={evidence} language={language} />;
  return null;
}
