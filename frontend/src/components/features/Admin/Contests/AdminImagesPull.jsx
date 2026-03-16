import {
  IconCheck,
  IconCopy,
  IconDownload,
  IconRefresh,
  IconSearch,
  IconServer,
  IconStack2,
  IconWriting,
} from '@tabler/icons-react';
import { Button, Card, Chip, EmptyState, Input, Textarea } from '../../../common';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

function AdminImagesPull({
  nodes = [],
  unionImages = [],
  allImages = [],
  selectedTargetKeys = [],
  selectedNodes = [],
  manualImagesText = '',
  pullPolicy = 'IfNotPresent',
  loading = false,
  submitting = false,
  onTargetToggle,
  onToggleAllTargets,
  onNodeToggle,
  onToggleAllNodes,
  onManualImagesChange,
  onPullPolicyChange,
  onPullFromSelection,
  onPullFromManualInput,
  onRefresh,
}) {
  const { t } = useTranslation();
  const [filterText, setFilterText] = useState('');

  const pullPolicyOptions = [
    { value: 'Always', label: t('admin.contests.imagesPull.pullPolicy.always') },
    { value: 'IfNotPresent', label: t('admin.contests.imagesPull.pullPolicy.ifNotPresent') },
    { value: 'Never', label: t('admin.contests.imagesPull.pullPolicy.never') },
  ];

  const manualImageCount = manualImagesText
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean).length;

  const allTargetCount = unionImages.reduce(
    (count, imageName) => count + nodes.filter((node) => !node.images.includes(imageName)).length,
    0
  );
  const isTargetSelected = (nodeName, imageName) => selectedTargetKeys.includes(`${nodeName}\u0000${imageName}`);
  const isMissingOnNode = (nodeName, imageName) => {
    const node = nodes.find((item) => item.node === nodeName);
    return node ? !node.images.includes(imageName) : false;
  };
  const normalizedFilter = filterText.trim().toLowerCase();
  const filteredUnionImages = unionImages.filter((imageName) => imageName.toLowerCase().includes(normalizedFilter));

  if (loading) {
    return (
      <div className="w-full mx-auto">
        <div className="flex justify-center items-center h-64">
          <div className="flex items-center gap-3 text-neutral-400">
            <div className="animate-spin">
              <IconRefresh size={20} />
            </div>
            <span className="font-mono">{t('common.loading')}</span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full mx-auto space-y-6">
      <div className="flex justify-end items-center">
        <Button variant="primary" size="sm" align="icon-left" icon={<IconRefresh size={16} />} onClick={onRefresh}>
          {t('common.refresh')}
        </Button>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <Card variant="default" padding="md" animate className="xl:col-span-2">
          <div className="flex items-start justify-between gap-4 mb-6">
            <div>
              <h2 className="text-lg font-mono text-neutral-50">{t('admin.contests.imagesPull.control.title')}</h2>
              <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">
                {t('admin.contests.imagesPull.control.subtitle')}
              </p>
            </div>
            <div className="flex min-w-fit shrink-0 items-center gap-2 whitespace-nowrap text-neutral-400 font-mono text-sm">
              <IconServer size={16} />
              <span>{t('admin.contests.imagesPull.summary.nodes', { count: nodes.length })}</span>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
              <div className="text-xs text-neutral-500 font-mono mb-2">
                {t('admin.contests.imagesPull.summary.nodesLabel')}
              </div>
              <div className="text-2xl text-neutral-50 font-mono">{nodes.length}</div>
            </div>
            <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
              <div className="text-xs text-neutral-500 font-mono mb-2">
                {t('admin.contests.imagesPull.summary.unionLabel')}
              </div>
              <div className="text-2xl text-neutral-50 font-mono">{unionImages.length}</div>
            </div>
            <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
              <div className="text-xs text-neutral-500 font-mono mb-2">
                {t('admin.contests.imagesPull.summary.totalLabel')}
              </div>
              <div className="text-2xl text-neutral-50 font-mono">{allImages.length}</div>
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

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 xl:items-stretch">
        <Card variant="default" padding="md" animate className="flex h-full flex-col">
          <div className="flex items-start justify-between gap-4 mb-5">
            <div>
              <h2 className="text-lg font-mono text-neutral-50">{t('admin.contests.imagesPull.selection.title')}</h2>
              <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">
                {t('admin.contests.imagesPull.selection.subtitle')}
              </p>
            </div>
            <Button variant="outline" size="sm" onClick={onToggleAllTargets} className="min-w-fit shrink-0 whitespace-nowrap">
              {allTargetCount > 0 && selectedTargetKeys.length === allTargetCount
                ? t('admin.contests.imagesPull.actions.deselectAllTargets')
                : t('admin.contests.imagesPull.actions.selectAllTargets')}
            </Button>
          </div>

          <div className="mb-4">
            <Input
              value={filterText}
              onChange={(e) => setFilterText(e.target.value)}
              placeholder={t('admin.contests.imagesPull.filters.placeholder')}
              icon={<IconSearch size={16} />}
              className="font-mono"
            />
          </div>

          <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
            <div className="flex items-start justify-between gap-4 mb-3">
              <div>
                <div className="text-sm font-mono text-neutral-300">
                  {t('admin.contests.imagesPull.intersection.title')}
                </div>
                <p className="mt-1 text-xs font-mono text-neutral-500 leading-5">
                  {t('admin.contests.imagesPull.intersection.subtitle')}
                </p>
              </div>
              <div className="text-xs font-mono text-neutral-500">
                {t('admin.contests.imagesPull.selection.selectedTargets', { count: selectedTargetKeys.length })}
              </div>
            </div>

            {filteredUnionImages.length === 0 ? (
              <EmptyState title={t('admin.contests.imagesPull.intersection.empty')} />
            ) : (
              <div className="max-h-130 overflow-y-auto pr-1 space-y-2">
                {filteredUnionImages.map((imageName) => (
                  <div key={imageName} className="rounded-md border border-neutral-300/20 bg-black/20 p-3">
                    <div className="break-all text-sm font-mono text-neutral-50 mb-3">{imageName}</div>
                    <div className="flex flex-wrap gap-2">
                      {nodes.map((node) => {
                        const isMissing = isMissingOnNode(node.node, imageName);
                        const isSelected = isTargetSelected(node.node, imageName);
                        return (
                          <button
                            key={`${node.node}-${imageName}`}
                            type="button"
                            onClick={isMissing ? () => onTargetToggle(node.node, imageName) : undefined}
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

            <div className="mt-4 flex items-center justify-between gap-4">
              <div className="text-xs font-mono text-neutral-500">
                {t('admin.contests.imagesPull.labels.comboHint')}
              </div>
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

        <Card variant="default" padding="md" animate>
          <div className="flex items-start gap-3 mb-5">
            <IconWriting size={18} className="mt-1 text-neutral-400" />
            <div>
              <h2 className="text-lg font-mono text-neutral-50">{t('admin.contests.imagesPull.manual.title')}</h2>
              <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">
                {t('admin.contests.imagesPull.manual.subtitle')}
              </p>
            </div>
          </div>

          <div className="rounded-md border border-neutral-300/20 bg-black/20 p-4 mb-4">
            <div className="flex items-center justify-between gap-4 mb-3">
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

          <div className="mt-3 flex items-center justify-between gap-4">
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
      </div>

      <Card variant="default" padding="md" animate>
        <div className="flex items-start justify-between gap-4 mb-5">
          <div>
            <h2 className="text-lg font-mono text-neutral-50">{t('admin.contests.imagesPull.statusTitle')}</h2>
            <p className="mt-2 text-sm text-neutral-400 font-mono leading-6">
              {t('admin.contests.imagesPull.statusSubtitle')}
            </p>
          </div>
          <Chip
            label={t('admin.contests.imagesPull.summary.total', { count: allImages.length })}
            variant="tag"
            colorClass="border-neutral-300/30 text-neutral-300"
          />
        </div>

        {nodes.length === 0 ? (
          <EmptyState title={t('admin.contests.imagesPull.empty')} />
        ) : (
          <div className="space-y-4">
            {nodes.map((node) => (
              <div key={node.node} className="rounded-md border border-neutral-300/20 bg-black/20 p-4">
                <div className="flex items-center justify-between gap-4 mb-3">
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
                    {node.images.filter((imageName) => imageName.toLowerCase().includes(normalizedFilter)).length === 0 ? (
                      <div className="text-sm font-mono text-neutral-500">
                        {t('admin.contests.imagesPull.node.empty')}
                      </div>
                    ) : (
                      <div className="flex flex-wrap gap-2 max-h-60 overflow-y-auto pr-1">
                        {node.images
                          .filter((imageName) => imageName.toLowerCase().includes(normalizedFilter))
                          .map((imageName) => (
                            <div
                              key={`${node.node}-${imageName}`}
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
                    <div className="text-xs font-mono text-neutral-500 mb-2">
                      {t('admin.contests.imagesPull.node.missingTitle')}
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {filteredUnionImages
                        .filter((imageName) => !node.images.includes(imageName))
                        .map((imageName) => (
                          <div
                            key={`${node.node}-${imageName}-missing`}
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
                      {filteredUnionImages.filter((imageName) => !node.images.includes(imageName)).length === 0 && (
                        <div className="text-sm font-mono text-neutral-500">
                          {t('admin.contests.imagesPull.node.noMissing')}
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}

export default AdminImagesPull;
