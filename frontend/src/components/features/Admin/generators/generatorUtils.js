export const GENERATOR_PAGE_SIZE = 20;

export const isGeneratorStoppable = (generator) => generator?.status === 'running';

export function normalizeStartCount(value) {
  const count = Number.parseInt(value, 10);
  return Number.isFinite(count) ? Math.max(0, count) : 0;
}

export function expandStartCounts(counts) {
  // Repeated challenge IDs request independent generator instances.
  return Object.entries(counts).flatMap(([id, count]) => Array(normalizeStartCount(count)).fill(id));
}

export function getGeneratorPageStats(generators) {
  return generators.reduce(
    (stats, generator) => ({
      successes: stats.successes + (generator.success ?? 0),
      failures: stats.failures + (generator.failure ?? 0),
    }),
    { successes: 0, failures: 0 }
  );
}

export function retainStoppableSelection(selectedIds, generators) {
  const stoppableIds = new Set(generators.filter(isGeneratorStoppable).map((generator) => generator.id));
  return selectedIds.filter((id) => stoppableIds.has(id));
}

export function getGeneratorResponseData(response) {
  if (response?.code !== 200) throw new Error(response?.msg || 'Generator request failed');
  return response.data ?? {};
}
