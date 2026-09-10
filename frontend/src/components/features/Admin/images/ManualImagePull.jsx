import { useTranslation } from 'react-i18next';
import { IconCopy, IconWriting } from '@tabler/icons-react';
import { Button, Card, Chip, EmptyState, Textarea } from '../../../common';
import { parseManualImages } from './imageModel';

export default function ManualImagePull({
  scopeKey,
  nodes,
  selectedNodes,
  manualImagesText,
  onNodeToggle,
  onToggleAllNodes,
  onManualImagesChange,
  onPullFromManualInput,
  submitting,
}) {
  const { t } = useTranslation();
  const manualImageCount = parseManualImages(manualImagesText).length;

  return (
    <Card variant="default" padding="md" animate className="min-w-0">
      <div className="flex items-start gap-3 mb-5">
        <IconWriting size={18} className="mt-1 text-neutral-400 shrink-0" />
        <div>
          <h2 className="text-lg font-mono text-neutral-50">{t(`${scopeKey}.manual.title`)}</h2>
          <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">{t(`${scopeKey}.manual.subtitle`)}</p>
        </div>
      </div>
      <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4 mb-4">
        <div className="flex flex-wrap items-center justify-between gap-4 mb-3">
          <div className="text-sm font-mono text-neutral-300">{t('admin.contests.imagesPull.labels.nodes')}</div>
          <Button variant="outline" size="sm" onClick={onToggleAllNodes}>
            {nodes.length > 0 && selectedNodes.length === nodes.length
              ? t('admin.contests.imagesPull.actions.deselectAllNodes')
              : t('admin.contests.imagesPull.actions.selectAllNodes')}
          </Button>
        </div>
        {nodes.length === 0 ? (
          <EmptyState title={t('admin.contests.imagesPull.empty')} />
        ) : (
          <div className="flex flex-wrap gap-2">
            {nodes.map((node) => {
              const isSelected = selectedNodes.includes(node.node);
              return (
                <button
                  key={node.node}
                  type="button"
                  onClick={() => onNodeToggle(node.node)}
                  className={`flex items-center gap-2 rounded-md border px-3 py-2 font-mono text-sm transition-colors ${
                    isSelected
                      ? 'border-geek-400/60 bg-geek-400/10 text-geek-300'
                      : 'border-neutral-300/20 bg-black/20 text-neutral-300 hover:border-neutral-300/40'
                  }`}
                >
                  <span>{node.node}</span>
                  <Chip
                    label={t('admin.contests.imagesPull.selection.nodeImages', { count: node.images.length })}
                    variant="tag"
                    size="sm"
                    colorClass={isSelected ? 'border-geek-400/40 text-geek-300' : undefined}
                  />
                </button>
              );
            })}
          </div>
        )}
      </div>
      <Textarea
        value={manualImagesText}
        onChange={(e) => onManualImagesChange(e.target.value)}
        rows={10}
        resize="vertical"
        placeholder={t('admin.contests.imagesPull.manual.placeholder')}
        className="font-mono"
      />
      <div className="mt-3 flex flex-wrap items-center justify-between gap-4">
        <div className="text-xs font-mono text-neutral-500">
          {t('admin.contests.imagesPull.manual.count', { count: manualImageCount })}
        </div>
        <Button
          variant="primary"
          onClick={onPullFromManualInput}
          loading={submitting}
          disabled={selectedNodes.length === 0 || manualImageCount === 0}
          align="icon-left"
          icon={<IconCopy size={18} />}
        >
          {t('admin.contests.imagesPull.actions.pullManual', { count: selectedNodes.length * manualImageCount })}
        </Button>
      </div>
    </Card>
  );
}
