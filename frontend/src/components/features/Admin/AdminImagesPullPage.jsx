import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import AdminImagesPull from './Contests/AdminImagesPull.jsx';
import { toast } from '../../../utils/toast';

function hasImageTag(image) {
  const value = String(image || '').trim();
  if (!value) {
    return false;
  }
  const lastSlash = value.lastIndexOf('/');
  const lastColon = value.lastIndexOf(':');
  return lastColon > lastSlash;
}

function normalizeNodes(data) {
  if (!Array.isArray(data)) {
    return [];
  }

  return data
    .map((item) => ({
      node: item?.node || '',
      images: Array.isArray(item?.images)
        ? Array.from(
            new Set(item.images.map((image) => String(image || '').trim()).filter((image) => hasImageTag(image)))
          ).sort()
        : [],
    }))
    .filter((item) => item.node)
    .sort((a, b) => a.node.localeCompare(b.node));
}

function normalizeTargetImages(data, nodes) {
  if (Array.isArray(data)) {
    return Array.from(
      new Set(data.map((image) => String(image || '').trim()).filter((image) => hasImageTag(image)))
    ).sort();
  }

  const fallback = new Set();
  for (const node of nodes) {
    for (const image of node.images) {
      fallback.add(image);
    }
  }
  return Array.from(fallback).sort();
}

function parseManualImages(text) {
  return Array.from(
    new Set(
      text
        .split('\n')
        .map((item) => item.trim())
        .filter((item) => item && hasImageTag(item))
    )
  );
}

function buildTargetKey(node, image) {
  return `${node}\u0000${image}`;
}

function parseTargetKey(key) {
  const [node, image] = key.split('\u0000');
  return { node, image };
}

function buildTargets(nodes, images) {
  return nodes.flatMap((node) =>
    images.map((image) => ({
      node,
      image,
    }))
  );
}

function normalizePayload(data) {
  const nodes = normalizeNodes(data?.nodes ?? data);
  const targetImages = normalizeTargetImages(data?.target_images, nodes);
  return { nodes, targetImages };
}

function AdminImagesPullPage({ scope = 'contest', fetchImages, pullImages, refreshKey }) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [nodes, setNodes] = useState([]);
  const [targetImages, setTargetImages] = useState([]);
  const [selectedTargetKeys, setSelectedTargetKeys] = useState([]);
  const [selectedNodes, setSelectedNodes] = useState([]);
  const [manualImagesText, setManualImagesText] = useState('');
  const [pullPolicy, setPullPolicy] = useState('IfNotPresent');

  const allImages = useMemo(() => {
    const set = new Set();
    for (const node of nodes) {
      for (const image of node.images) {
        set.add(image);
      }
    }
    return Array.from(set).sort();
  }, [nodes]);

  const availableTargetKeys = useMemo(
    () =>
      targetImages.flatMap((image) =>
        nodes.filter((node) => !node.images.includes(image)).map((node) => buildTargetKey(node.node, image))
      ),
    [nodes, targetImages]
  );

  useEffect(() => {
    fetchPullImages();
  }, [refreshKey]);

  useEffect(() => {
    setSelectedNodes((prev) => prev.filter((node) => nodes.some((item) => item.node === node)));
  }, [nodes]);

  useEffect(() => {
    const availableSet = new Set(availableTargetKeys);
    setSelectedTargetKeys((prev) => prev.filter((key) => availableSet.has(key)));
  }, [availableTargetKeys]);

  const fetchPullImages = async () => {
    setLoading(true);
    try {
      const response = await fetchImages();
      if (response.code === 200) {
        const normalized = normalizePayload(response.data);
        setNodes(normalized.nodes);
        setTargetImages(normalized.targetImages);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.imagesPull.toast.fetchFailed') });
    } finally {
      setLoading(false);
    }
  };

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

      if (response.code === 200) {
        toast.success({ description: t('admin.contests.imagesPull.toast.submitSuccess') });
        setTimeout(() => {
          fetchPullImages();
        }, 2000);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.imagesPull.toast.pullFailed') });
    } finally {
      setSubmitting(false);
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
    <AdminImagesPull
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
      onRefresh={fetchPullImages}
    />
  );
}

export default AdminImagesPullPage;
