export const ANALYSIS_TABS = ['flags', 'indicators', 'http', 'dns', 'sessions', 'accesses', 'overlaps'];

export function evidenceOffset(evidence, startedAt) {
  const start = Date.parse(startedAt);
  const time = Date.parse(evidence?.time);
  return Number.isFinite(start) && Number.isFinite(time) ? time - start : null;
}

export function filterAnalysisRows(rows, { startedAt, window, currentWindow, nodeId, edge }) {
  return rows.filter((row) => {
    const evidence = row.evidence || row;
    if (nodeId && ![evidence.src_ip, evidence.dst_ip, row.ip].includes(nodeId)) return false;
    if (
      edge &&
      !(
        (evidence.src_ip === edge.source && evidence.dst_ip === edge.target) ||
        (evidence.src_ip === edge.target && evidence.dst_ip === edge.source)
      )
    )
      return false;
    if (!currentWindow) return true;
    const offset = evidenceOffset(evidence, startedAt);
    if (offset === null) return false;
    const end = evidenceOffset({ time: evidence.end_time || evidence.time }, startedAt);
    return offset < window.end && Math.max(offset, end ?? offset) >= window.start;
  });
}

export function victimOverlaps(overlaps, victimId) {
  return overlaps.filter((overlap) => overlap.teams?.some((team) => team.victim_ids?.includes(Number(victimId))));
}
