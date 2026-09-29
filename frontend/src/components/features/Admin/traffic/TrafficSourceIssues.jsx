import { useTranslation } from 'react-i18next';

export default function TrafficSourceIssues({ issues }) {
  const { t } = useTranslation();
  if (!issues?.length) return null;
  return (
    <details className="mt-3 rounded-lg border border-amber-500/30 p-3 text-xs text-amber-200">
      <summary className="cursor-pointer">
        {t('admin.contests.trafficGraph.analysis.sourceIssues', { count: issues.length })}
      </summary>
      <ul className="mt-2 grid max-h-64 gap-2 overflow-y-auto">
        {issues.map((issue, index) => (
          <li key={`${issue.file}-${issue.phase}-${index}`} className="break-all">
            <strong>{issue.file}</strong> · {issue.phase}
            <div className="whitespace-pre-wrap">{issue.error}</div>
          </li>
        ))}
      </ul>
    </details>
  );
}
