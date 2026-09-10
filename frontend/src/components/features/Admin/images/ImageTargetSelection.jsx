import { useTranslation } from 'react-i18next';
import { IconCheck, IconDownload, IconSearch } from '@tabler/icons-react';
import { Button, Card, Chip, EmptyState, Input } from '../../../common';
import { buildTargetKey, missingTargetKeys } from './imageModel';

export default function ImageTargetSelection({
  scopeKey,
  nodes,
  targetImages,
  filteredTargetImages,
  selectedTargetKeys,
  filterText,
  onFilterChange,
  onTargetToggle,
  onToggleAllTargets,
  onPullFromSelection,
  submitting,
}) {
  const { t } = useTranslation();
  const allTargetCount = missingTargetKeys(nodes, targetImages).length;

  return (
    <Card variant="default" padding="md" animate className="flex h-full min-w-0 flex-col">
      <div className="flex flex-wrap items-start justify-between gap-4 mb-5">
        <div>
          <h2 className="text-lg font-mono text-neutral-50">{t(`${scopeKey}.selection.title`)}</h2>
          <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">{t(`${scopeKey}.selection.subtitle`)}</p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={onToggleAllTargets}
          className="min-w-fit shrink-0 whitespace-nowrap"
        >
          {allTargetCount > 0 && selectedTargetKeys.length === allTargetCount
            ? t('admin.contests.imagesPull.actions.deselectAllTargets')
            : t('admin.contests.imagesPull.actions.selectAllTargets')}
        </Button>
      </div>
      <div className="mb-4">
        <Input
          value={filterText}
          onChange={(e) => onFilterChange(e.target.value)}
          placeholder={t('admin.contests.imagesPull.filters.placeholder')}
          icon={<IconSearch size={16} />}
          className="font-mono"
        />
      </div>
      <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
        <div className="flex flex-wrap items-start justify-between gap-4 mb-3">
          <div>
            <div className="text-sm font-mono text-neutral-300">{t(`${scopeKey}.intersection.title`)}</div>
            <p className="mt-1 text-xs font-mono text-neutral-500 leading-5">
              {t(`${scopeKey}.intersection.subtitle`)}
            </p>
          </div>
          <div className="text-xs font-mono text-neutral-500">
            {t('admin.contests.imagesPull.selection.selectedTargets', { count: selectedTargetKeys.length })}
          </div>
        </div>
        {filteredTargetImages.length === 0 ? (
          <EmptyState title={t(`${scopeKey}.intersection.empty`)} />
        ) : (
          <div className="max-h-130 overflow-y-auto pr-1 space-y-2">
            {filteredTargetImages.map((imageName) => (
              <div key={imageName} className="rounded-md border border-neutral-300/20 bg-black/20 p-3">
                <div className="break-all text-sm font-mono text-neutral-50 mb-3">{imageName}</div>
                <div className="flex flex-wrap gap-2">
                  {nodes.map((node) => {
                    const isMissing = !node.images.includes(imageName);
                    const isSelected = selectedTargetKeys.includes(buildTargetKey(node.node, imageName));
                    return (
                      <button
                        key={buildTargetKey(node.node, imageName)}
                        type="button"
                        disabled={!isMissing}
                        onClick={() => onTargetToggle(node.node, imageName)}
                        className={`flex items-center gap-2 rounded-md border px-3 py-2 font-mono text-sm transition-colors ${
                          !isMissing
                            ? 'border-neutral-300/10 bg-black/10 text-neutral-500 cursor-not-allowed'
                            : isSelected
                              ? 'border-geek-400/60 bg-geek-400/10 text-geek-300'
                              : 'border-neutral-300/20 bg-black/20 text-neutral-300 hover:border-neutral-300/40'
                        }`}
                      >
                        <span>{node.node}</span>
                        <Chip
                          label={
                            isMissing
                              ? isSelected
                                ? t('admin.contests.imagesPull.status.selectedMissing')
                                : t('admin.contests.imagesPull.status.missing')
                              : t('admin.contests.imagesPull.status.present')
                          }
                          variant="tag"
                          size="sm"
                          colorClass={
                            !isMissing
                              ? 'border-neutral-300/20 text-neutral-500'
                              : isSelected
                                ? 'border-geek-400/40 text-geek-300'
                                : 'border-amber-400/40 text-amber-300'
                          }
                        />
                        {isSelected && <IconCheck size={14} />}
                      </button>
                    );
                  })}
                </div>
              </div>
            ))}
          </div>
        )}
        <div className="mt-4 flex flex-wrap items-center justify-between gap-4">
          <div className="text-xs font-mono text-neutral-500">{t('admin.contests.imagesPull.labels.comboHint')}</div>
          <Button
            variant="primary"
            onClick={onPullFromSelection}
            loading={submitting}
            disabled={selectedTargetKeys.length === 0}
            align="icon-left"
            icon={<IconDownload size={18} />}
          >
            {t('admin.contests.imagesPull.actions.pullSelection', { count: selectedTargetKeys.length })}
          </Button>
        </div>
      </div>
    </Card>
  );
}
