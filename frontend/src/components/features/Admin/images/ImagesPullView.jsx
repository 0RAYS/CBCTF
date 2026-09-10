import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconRefresh } from '@tabler/icons-react';
import { Button } from '../../../common';
import ImagePullSummary from './ImagePullSummary';
import ImageTargetSelection from './ImageTargetSelection';
import ManualImagePull from './ManualImagePull';
import NodeImageInventory from './NodeImageInventory';

export default function ImagesPullView({ scope, loading, onRefresh, ...props }) {
  const { t } = useTranslation();
  const [filterText, setFilterText] = useState('');
  const scopeKey = scope === 'global' ? 'admin.imagesPull' : 'admin.contests.imagesPull';
  const normalizedFilter = filterText.trim().toLowerCase();
  const filteredTargetImages = props.targetImages.filter((image) => image.toLowerCase().includes(normalizedFilter));

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
