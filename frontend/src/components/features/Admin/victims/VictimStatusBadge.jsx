const styles = {
  waiting: 'bg-yellow-400/10 text-yellow-400 border-yellow-400/30',
  pending: 'bg-geek-400/10 text-geek-400 border-geek-400/30',
  terminating: 'bg-orange-400/10 text-orange-300 border-orange-400/30',
  running: 'bg-green-400/10 text-green-400 border-green-400/30',
  stopped: 'bg-neutral-500/10 text-neutral-400 border-neutral-500/30',
};

export default function VictimStatusBadge({ status, t, translationKey }) {
  return (
    <span className={`inline-block px-2 py-0.5 rounded border text-xs font-mono ${styles[status] ?? styles.stopped}`}>
      {t(`${translationKey}.statusBadge.${status}`, status)}
    </span>
  );
}
