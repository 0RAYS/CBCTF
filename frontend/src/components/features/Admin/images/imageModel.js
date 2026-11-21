export function hasImageTag(image) {
  const value = String(image || '').trim();
  if (!value) return false;
  return value.lastIndexOf(':') > value.lastIndexOf('/');
}

function normalizeImageList(images) {
  if (!Array.isArray(images)) return [];
  return Array.from(new Set(images.map((image) => String(image || '').trim()).filter(hasImageTag))).sort();
}

export function normalizeNodes(data) {
  if (!Array.isArray(data)) return [];
  return data
    .map((item) => ({
      node: item?.node || '',
      images: normalizeImageList(item?.images),
      missing: normalizeImageList(item?.missing),
    }))
    .filter((item) => item.node)
    .sort((a, b) => a.node.localeCompare(b.node));
}

export function normalizeTargetImages(data, nodes) {
  if (Array.isArray(data)) {
    return normalizeImageList(data);
  }
  return Array.from(new Set(nodes.flatMap((node) => node.images))).sort();
}

export function normalizePayload(data) {
  const nodes = normalizeNodes(data?.nodes ?? data);
  return {
    nodes,
    targetImages: normalizeTargetImages(data?.target_images, nodes),
    detectedAt: String(data?.detected_at || ''),
  };
}

export function parseManualImages(text) {
  return Array.from(
    new Set(
      text
        .split('\n')
        .map((item) => item.trim())
        .filter(hasImageTag)
    )
  );
}

// Every detected node is selected by default; the deselected list carries the explicit exclusions.
export function selectedNodeNames(nodes, deselected) {
  return nodes.map((node) => node.node).filter((name) => !deselected.includes(name));
}

export function toggleDeselectedNode(deselected, nodeName) {
  return deselected.includes(nodeName) ? deselected.filter((name) => name !== nodeName) : [...deselected, nodeName];
}

export function buildTargetKey(node, image) {
  return `${node}\u0000${image}`;
}

export function parseTargetKey(key) {
  const [node, image] = key.split('\u0000');
  return { node, image };
}

export function buildTargets(nodes, images) {
  return nodes.flatMap((node) => images.map((image) => ({ node, image })));
}

export const IMAGE_STATE_PRESENT = 'present';
export const IMAGE_STATE_MISSING = 'missing';
export const IMAGE_STATE_UNVERIFIED = 'unverified';

// A node holds an image only when the last detection proved it; the rest is proven absent or unverified.
export function imageState(node, image) {
  if (node.images.includes(image)) return IMAGE_STATE_PRESENT;
  if (node.missing.includes(image)) return IMAGE_STATE_MISSING;
  return IMAGE_STATE_UNVERIFIED;
}

export function unconfirmedTargetKeys(nodes, images) {
  return images.flatMap((image) =>
    nodes.filter((node) => !node.images.includes(image)).map((node) => buildTargetKey(node.node, image))
  );
}

export function selectableTargetKeys(nodes, images, pullPolicy) {
  if (pullPolicy === 'Never') return [];
  if (pullPolicy === 'Always') {
    return buildTargets(
      nodes.map((node) => node.node),
      images
    ).map(({ node, image }) => buildTargetKey(node, image));
  }
  return unconfirmedTargetKeys(nodes, images);
}

export function completedImageTargetKeys(targets, batch) {
  const completedNodes = new Set(
    batch.items.filter((item) => ['success', 'skipped'].includes(item.status)).map((item) => item.id)
  );
  return targets
    .filter((target) => completedNodes.has(target.node))
    .map(({ node, image }) => buildTargetKey(node, image));
}

export function formatDetectionTime(value, language) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString(language || undefined, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}
