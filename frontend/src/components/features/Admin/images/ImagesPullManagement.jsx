import { useEffect, useEffectEvent, useMemo, useReducer, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ImagesPullView from './ImagesPullView';
import { toast } from '../../../../utils/toast';
import useBatchAction from '../batch/useBatchAction.js';
import BatchResultPanel from '../batch/BatchResultPanel.jsx';
import { remainingBatchIds } from '../batch/batchModel.js';
import {
  buildTargetKey,
  buildTargets,
  selectableTargetKeys,
  normalizePayload,
  normalizeTargetImages,
  parseManualImages,
  parseTargetKey,
  completedImageTargetKeys,
} from './imageModel';

const DETECT_POLL_INTERVAL = 2000;
const DETECT_POLL_TIMEOUT = 5 * 60 * 1000;

function ImagesPullManagement({ scope = 'contest', fetchImages, pullImages, detectImages }) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const action = useBatchAction(scope);
  const [nodes, setNodes] = useState([]);
  const [targetImages, setTargetImages] = useState([]);
  const [selectedTargetKeys, setSelectedTargetKeys] = useState([]);
  const [selectedNodes, setSelectedNodes] = useState([]);
  const [manualImagesText, setManualImagesText] = useState('');
  const [pullPolicy, setPullPolicy] = useState('IfNotPresent');
  const [detectedAt, setDetectedAt] = useState('');
  const [detecting, setDetecting] = useState(false);
  const [revision, refresh] = useReducer((value) => value + 1, 0);
  const refreshTimer = useRef(null);
  const detectedAtRef = useRef('');

  const allImages = useMemo(() => normalizeTargetImages(undefined, nodes), [nodes]);
  const availableTargetKeys = useMemo(
    () => selectableTargetKeys(nodes, targetImages, pullPolicy),
    [nodes, targetImages, pullPolicy]
  );

  const loadImages = useEffectEvent(() => fetchImages());
  const reportFetchError = useEffectEvent((error) => {
    toast.danger({ description: error.message || t('admin.contests.imagesPull.toast.fetchFailed') });
  });
  const applyPayload = useEffectEvent((normalized) => {
    setNodes(normalized.nodes);
    setTargetImages(normalized.targetImages);
    setDetectedAt(normalized.detectedAt);
    detectedAtRef.current = normalized.detectedAt;
  });

  useEffect(() => {
    return () => {
      clearTimeout(refreshTimer.current);
    };
  }, []);

  useEffect(() => {
    let current = true;
    setLoading(true);
    Promise.resolve().then(async () => {
      if (!current) return;
      try {
        const response = await loadImages();
        if (current && response.code === 200) {
          applyPayload(normalizePayload(response.data));
        }
      } catch (error) {
        if (current) reportFetchError(error);
      } finally {
        if (current) setLoading(false);
      }
    });
    return () => {
      current = false;
    };
  }, [revision]);

  // Detection runs in the background, so the page waits for the new result.
  useEffect(() => {
    if (!detecting) return undefined;
    let cancelled = false;
    let busy = false;
    const deadline = Date.now() + DETECT_POLL_TIMEOUT;
    const timer = setInterval(async () => {
      if (cancelled || busy) return;
      if (Date.now() > deadline) {
        cancelled = true;
        setDetecting(false);
        return;
      }
      busy = true;
      try {
        const response = await loadImages();
        if (!cancelled && response?.code === 200) {
          const normalized = normalizePayload(response.data);
          if (normalized.detectedAt && normalized.detectedAt !== detectedAtRef.current) {
            applyPayload(normalized);
            setDetecting(false);
          }
        }
      } catch {
        // Keep polling until the deadline.
      } finally {
        busy = false;
      }
    }, DETECT_POLL_INTERVAL);
    return () => {
      cancelled = true;
      clearInterval(timer);
    };
  }, [detecting]);

  useEffect(() => {
    setSelectedNodes((prev) => prev.filter((node) => nodes.some((item) => item.node === node)));
  }, [nodes]);

  useEffect(() => {
    const availableSet = new Set(availableTargetKeys);
    setSelectedTargetKeys((prev) => prev.filter((key) => availableSet.has(key)));
  }, [availableTargetKeys]);

  const handleNodeToggle = (nodeName) => {
    setSelectedNodes((prev) =>
      prev.includes(nodeName) ? prev.filter((node) => node !== nodeName) : [...prev, nodeName]
    );
  };

  const handleToggleAllNodes = () => {
    if (selectedNodes.length === nodes.length) {
      setSelectedNodes([]);
    } else {
      setSelectedNodes(nodes.map((node) => node.node));
    }
  };

  const handleTargetToggle = (nodeName, imageName) => {
    const targetKey = buildTargetKey(nodeName, imageName);
    setSelectedTargetKeys((prev) =>
      prev.includes(targetKey) ? prev.filter((key) => key !== targetKey) : [...prev, targetKey]
    );
  };

  const handleToggleAllTargets = () => {
    if (availableTargetKeys.length > 0 && selectedTargetKeys.length === availableTargetKeys.length) {
      setSelectedTargetKeys([]);
    } else {
      setSelectedTargetKeys(availableTargetKeys);
    }
  };

  const submitTargets = async (targets, emptyMessageKey, requireNodeSelection = false) => {
    if (requireNodeSelection && selectedNodes.length === 0) {
      toast.warning({ description: t('admin.contests.imagesPull.toast.nodeRequired') });
      return;
    }

    if (targets.length === 0) {
      toast.warning({ description: t(emptyMessageKey) });
      return;
    }

    await action.run(() => pullImages({ targets, pull_policy: pullPolicy }), {
      successMessage: t(
        pullPolicy === 'Never'
          ? 'admin.contests.imagesPull.toast.skipped'
          : 'admin.contests.imagesPull.toast.submitSuccess'
      ),
      failureMessage: t('admin.contests.imagesPull.toast.pullFailed'),
      onResult: (batch) => {
        if (batch) {
          if (requireNodeSelection) {
            setSelectedNodes((current) => remainingBatchIds(current, batch));
          } else {
            const completed = new Set(completedImageTargetKeys(targets, batch));
            setSelectedTargetKeys((current) => current.filter((key) => !completed.has(key)));
          }
        }
        refresh();
        clearTimeout(refreshTimer.current);
        refreshTimer.current = setTimeout(refresh, 2000);
      },
    });
  };

  const handlePullFromSelection = async () => {
    const targets = selectedTargetKeys.map(parseTargetKey);
    await submitTargets(targets, 'admin.contests.imagesPull.toast.selectRequired');
  };

  const handlePullFromManualInput = async () => {
    const manualImages = parseManualImages(manualImagesText);
    const targets = buildTargets(selectedNodes, manualImages).map((target) => ({
      ...target,
      manual: true,
    }));
    await submitTargets(targets, 'admin.contests.imagesPull.toast.manualRequired', true);
  };

  const handleDetect = async () => {
    if (!detectImages) return;
    await action.run(() => detectImages(), {
      successMessage: t('admin.contests.imagesPull.toast.detectQueued'),
      failureMessage: t('admin.contests.imagesPull.toast.detectFailed'),
      onResult: (batch) => {
        if (batch && batch.failed > 0) return;
        setDetecting(true);
      },
    });
  };

  return (
    <div className="space-y-4">
      <BatchResultPanel result={action.result} error={action.error} queued />
      <ImagesPullView
        scope={scope}
        nodes={nodes}
        targetImages={targetImages}
        allImages={allImages}
        selectedTargetKeys={selectedTargetKeys}
        selectedNodes={selectedNodes}
        manualImagesText={manualImagesText}
        pullPolicy={pullPolicy}
        loading={loading}
        submitting={action.pending}
        detecting={detecting}
        detectedAt={detectedAt}
        onDetect={handleDetect}
        onTargetToggle={handleTargetToggle}
        onToggleAllTargets={handleToggleAllTargets}
        onNodeToggle={handleNodeToggle}
        onToggleAllNodes={handleToggleAllNodes}
        onManualImagesChange={setManualImagesText}
        onPullPolicyChange={setPullPolicy}
        onPullFromSelection={handlePullFromSelection}
        onPullFromManualInput={handlePullFromManualInput}
        onRefresh={refresh}
      />
    </div>
  );
}

export default ImagesPullManagement;
