import { useTranslation } from 'react-i18next';
import { IconServer, IconStack2 } from '@tabler/icons-react';
import { Card } from '../../../common';

export default function ImagePullSummary({
  scopeKey,
  nodeCount,
  targetCount,
  imageCount,
  pullPolicy,
  onPullPolicyChange,
}) {
  const { t } = useTranslation();
  const pullPolicyOptions = [
    { value: 'Always', label: t('admin.contests.imagesPull.pullPolicy.always') },
    { value: 'IfNotPresent', label: t('admin.contests.imagesPull.pullPolicy.ifNotPresent') },
    { value: 'Never', label: t('admin.contests.imagesPull.pullPolicy.never') },
  ];

  return (
    <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <Card variant="default" padding="md" animate className="xl:col-span-2">
        <div className="flex flex-wrap items-start justify-between gap-4 mb-6">
          <div>
            <h2 className="text-lg font-mono text-neutral-50">{t(`${scopeKey}.control.title`)}</h2>
            <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">{t(`${scopeKey}.control.subtitle`)}</p>
          </div>
          <div className="flex min-w-fit shrink-0 items-center gap-2 whitespace-nowrap text-neutral-400 font-mono text-sm">
            <IconServer size={16} />
            <span>{t('admin.contests.imagesPull.summary.nodes', { count: nodeCount })}</span>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
            <div className="text-xs text-neutral-500 font-mono mb-2">
              {t('admin.contests.imagesPull.summary.nodesLabel')}
            </div>
            <div className="text-2xl text-neutral-50 font-mono">{nodeCount}</div>
          </div>
          <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
            <div className="text-xs text-neutral-500 font-mono mb-2">{t(`${scopeKey}.summary.unionLabel`)}</div>
            <div className="text-2xl text-neutral-50 font-mono">{targetCount}</div>
          </div>
          <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
            <div className="text-xs text-neutral-500 font-mono mb-2">
              {t('admin.contests.imagesPull.summary.totalLabel')}
            </div>
            <div className="text-2xl text-neutral-50 font-mono">{imageCount}</div>
          </div>
        </div>
      </Card>
      <Card variant="default" padding="md" animate>
        <div className="flex items-center gap-2 mb-4 text-neutral-300">
          <IconStack2 size={18} />
          <h3 className="text-base font-mono text-neutral-50">{t('admin.contests.imagesPull.labels.pullPolicy')}</h3>
        </div>
        <select
          value={pullPolicy}
          onChange={(e) => onPullPolicyChange(e.target.value)}
          className="select-custom select-custom-md"
        >
          {pullPolicyOptions.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        <p className="mt-3 text-xs text-neutral-500 font-mono leading-5">
          {t('admin.contests.imagesPull.labels.pullPolicyHint')}
        </p>
      </Card>
    </div>
  );
}
