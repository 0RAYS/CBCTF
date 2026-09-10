import { useTranslation } from 'react-i18next';
import { IconServer } from '@tabler/icons-react';
import { Card, Chip, EmptyState } from '../../../common';
import { buildTargetKey } from './imageModel';

export default function NodeImageInventory({ scopeKey, nodes, imageCount, filteredTargetImages, normalizedFilter }) {
  const { t } = useTranslation();

  return (
    <Card variant="default" padding="md" animate>
      <div className="flex flex-wrap items-start justify-between gap-4 mb-5">
        <div>
          <h2 className="text-lg font-mono text-neutral-50">{t(`${scopeKey}.statusTitle`)}</h2>
          <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">{t(`${scopeKey}.statusSubtitle`)}</p>
        </div>
        <Chip
          label={t('admin.contests.imagesPull.summary.total', { count: imageCount })}
          variant="tag"
          colorClass="border-neutral-300/30 text-neutral-300"
        />
      </div>
      {nodes.length === 0 ? (
        <EmptyState title={t('admin.contests.imagesPull.empty')} />
      ) : (
        <div className="space-y-4">
          {nodes.map((node) => {
            const presentImages = node.images.filter((image) => image.toLowerCase().includes(normalizedFilter));
            const missingImages = filteredTargetImages.filter((image) => !node.images.includes(image));
            return (
              <div key={node.node} className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
                <div className="flex flex-wrap items-center justify-between gap-4 mb-3">
                  <div className="flex items-center gap-2">
                    <IconServer size={16} className="text-neutral-400" />
                    <span className="text-sm font-mono text-neutral-50">{node.node}</span>
                  </div>
                  <Chip
                    label={t('admin.contests.imagesPull.selection.nodeImages', { count: node.images.length })}
                    variant="tag"
                  />
                </div>
                <div className="space-y-3">
                  <div>
                    {presentImages.length === 0 ? (
                      <div className="text-sm font-mono text-neutral-500">
                        {t('admin.contests.imagesPull.node.empty')}
                      </div>
                    ) : (
                      <div className="flex flex-wrap gap-2 max-h-60 overflow-y-auto pr-1">
                        {presentImages.map((imageName) => (
                          <div
                            key={buildTargetKey(node.node, imageName)}
                            className="flex items-center gap-2 rounded-md border border-neutral-300/20 bg-black/30 px-3 py-2"
                          >
                            <span className="break-all text-xs font-mono text-neutral-200">{imageName}</span>
                            <Chip
                              label={t('admin.contests.imagesPull.status.present')}
                              variant="tag"
                              size="sm"
                              colorClass="border-green-400/30 text-green-300"
                            />
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                  <div className="pt-3 border-t border-neutral-300/10">
                    <div className="text-xs font-mono text-neutral-500 mb-2">{t(`${scopeKey}.node.missingTitle`)}</div>
                    <div className="flex flex-wrap gap-2">
                      {missingImages.map((imageName) => (
                        <div
                          key={buildTargetKey(node.node, imageName)}
                          className="flex items-center gap-2 rounded-md border border-amber-400/20 bg-amber-400/5 px-3 py-2"
                        >
                          <span className="break-all text-xs font-mono text-amber-100">{imageName}</span>
                          <Chip
                            label={t('admin.contests.imagesPull.status.missing')}
                            variant="tag"
                            size="sm"
                            colorClass="border-amber-400/30 text-amber-300"
                          />
                        </div>
                      ))}
                      {missingImages.length === 0 && (
                        <div className="text-sm font-mono text-neutral-500">{t(`${scopeKey}.node.noMissing`)}</div>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </Card>
  );
}
