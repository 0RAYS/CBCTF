import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconRefresh, IconSearch } from '@tabler/icons-react';
import { Button } from '@/components/common';
import Loading from '@/components/common/Loading';
import ImagePullSummary from '@/components/features/Admin/images/ImagePullSummary';
import ImageTargetSelection from '@/components/features/Admin/images/ImageTargetSelection';
import ManualImagePull from '@/components/features/Admin/images/ManualImagePull';
import NodeImageInventory from '@/components/features/Admin/images/NodeImageInventory';
import { formatDetectionTime } from '@/components/features/Admin/images/imageModel';

export default function ImagesPullView({ scope, loading, onRefresh, ...props }) {
  const { t, i18n } = useTranslation();
  const [filterText, setFilterText] = useState('');
  const scopeKey = scope === 'global' ? 'admin.imagesPull' : 'admin.contests.imagesPull';
  const normalizedFilter = filterText.trim().toLowerCase();
  const filteredTargetImages = props.targetImages.filter((image) => image.toLowerCase().includes(normalizedFilter));

  if (loading) {
    return (
      <div className="w-full mx-auto">
        <Loading />
      </div>
    );
  }

  return (
    <div className="w-full mx-auto space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="text-xs font-mono text-neutral-500">
          {props.detectedAt
            ? t('admin.contests.imagesPull.detect.lastResult', {
                time: formatDetectionTime(props.detectedAt, i18n.language),
              })
            : t('admin.contests.imagesPull.detect.noResult')}
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            align="icon-left"
            icon={<IconSearch size={16} />}
            loading={props.detecting}
            disabled={props.detecting}
            onClick={props.onDetect}
          >
            {t('admin.contests.imagesPull.actions.detect')}
          </Button>
          <Button variant="primary" size="sm" align="icon-left" icon={<IconRefresh size={16} />} onClick={onRefresh}>
            {t('common.refresh')}
          </Button>
        </div>
      </div>
      <ImagePullSummary
        scopeKey={scopeKey}
        nodeCount={props.nodes.length}
        targetCount={props.targetImages.length}
        imageCount={props.allImages.length}
        pullPolicy={props.pullPolicy}
        onPullPolicyChange={props.onPullPolicyChange}
      />
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 xl:items-stretch">
        <ImageTargetSelection
          scopeKey={scopeKey}
          nodes={props.nodes}
          targetImages={props.targetImages}
          pullPolicy={props.pullPolicy}
          filteredTargetImages={filteredTargetImages}
          selectedTargetKeys={props.selectedTargetKeys}
          filterText={filterText}
          onFilterChange={setFilterText}
          onTargetToggle={props.onTargetToggle}
          onToggleAllTargets={props.onToggleAllTargets}
          onPullFromSelection={props.onPullFromSelection}
          submitting={props.submitting}
        />
        <ManualImagePull
          scopeKey={scopeKey}
          nodes={props.nodes}
          selectedNodes={props.selectedNodes}
          manualImagesText={props.manualImagesText}
          onNodeToggle={props.onNodeToggle}
          onToggleAllNodes={props.onToggleAllNodes}
          onManualImagesChange={props.onManualImagesChange}
          onPullFromManualInput={props.onPullFromManualInput}
          submitting={props.submitting}
        />
      </div>
      <NodeImageInventory
        scopeKey={scopeKey}
        nodes={props.nodes}
        imageCount={props.allImages.length}
        filteredTargetImages={filteredTargetImages}
        normalizedFilter={normalizedFilter}
      />
    </div>
  );
}
