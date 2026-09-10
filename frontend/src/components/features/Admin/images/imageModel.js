export function hasImageTag(image) {
  const value = String(image || '').trim();
  if (!value) return false;
  return value.lastIndexOf(':') > value.lastIndexOf('/');
}

export function normalizeNodes(data) {
  if (!Array.isArray(data)) return [];
  return data
    .map((item) => ({
      node: item?.node || '',
      images: Array.isArray(item?.images)
        ? Array.from(new Set(item.images.map((image) => String(image || '').trim()).filter(hasImageTag))).sort()
        : [],
    }))
    .filter((item) => item.node)
    .sort((a, b) => a.node.localeCompare(b.node));
}

export function normalizeTargetImages(data, nodes) {
  if (Array.isArray(data)) {
    return Array.from(new Set(data.map((image) => String(image || '').trim()).filter(hasImageTag))).sort();
  }
  return Array.from(new Set(nodes.flatMap((node) => node.images))).sort();
}

export function normalizePayload(data) {
  const nodes = normalizeNodes(data?.nodes ?? data);
  return { nodes, targetImages: normalizeTargetImages(data?.target_images, nodes) };
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

export function missingTargetKeys(nodes, images) {
  return images.flatMap((image) =>
    nodes.filter((node) => !node.images.includes(image)).map((node) => buildTargetKey(node.node, image))
  );
}
