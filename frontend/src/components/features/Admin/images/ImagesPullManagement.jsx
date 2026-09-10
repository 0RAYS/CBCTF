import { useEffect, useEffectEvent, useMemo, useReducer, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ImagesPullView from './ImagesPullView';
import { toast } from '../../../../utils/toast';
import {
  buildTargetKey,
  buildTargets,
  missingTargetKeys,
  normalizePayload,
  normalizeTargetImages,
  parseManualImages,
  parseTargetKey,
} from './imageModel';

function ImagesPullManagement({ scope = 'contest', fetchImages, pullImages }) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [nodes, setNodes] = useState([]);
  const [targetImages, setTargetImages] = useState([]);
  const [selectedTargetKeys, setSelectedTargetKeys] = useState([]);
  const [selectedNodes, setSelectedNodes] = useState([]);
  const [manualImagesText, setManualImagesText] = useState('');
  const [pullPolicy, setPullPolicy] = useState('IfNotPresent');
  const [revision, refresh] = useReducer((value) => value + 1, 0);
  const refreshTimer = useRef(null);
  const mounted = useRef(false);

  const allImages = useMemo(() => normalizeTargetImages(undefined, nodes), [nodes]);
  const availableTargetKeys = useMemo(() => missingTargetKeys(nodes, targetImages), [nodes, targetImages]);

  const loadImages = useEffectEvent(() => fetchImages());
  const reportFetchError = useEffectEvent((error) => {
    toast.danger({ description: error.message || t('admin.contests.imagesPull.toast.fetchFailed') });
  });

  useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
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
          const normalized = normalizePayload(response.data);
          setNodes(normalized.nodes);
          setTargetImages(normalized.targetImages);
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

    setSubmitting(true);
    try {
      const response = await pullImages({
        targets,
        pull_policy: pullPolicy,
      });

      if (mounted.current && response.code === 200) {
        toast.success({ description: t('admin.contests.imagesPull.toast.submitSuccess') });
        clearTimeout(refreshTimer.current);
        refreshTimer.current = setTimeout(refresh, 2000);
      }
    } catch (error) {
      if (mounted.current) {
        toast.danger({ description: error.message || t('admin.contests.imagesPull.toast.pullFailed') });
      }
    } finally {
      if (mounted.current) setSubmitting(false);
    }
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

  return (
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
      submitting={submitting}
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
  );
}

export default ImagesPullManagement;
