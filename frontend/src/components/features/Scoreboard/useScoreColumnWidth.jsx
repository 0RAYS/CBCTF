import { useLayoutEffect, useRef, useState } from 'react';

/** Measure rendered text, including locale separators, responsive type and loaded fonts. */
export default function useScoreColumnWidth({
  scores,
  label,
  scoreClassName,
  padding = 0,
  minWidth = 64,
  maxWidth = 240,
}) {
  const measureRef = useRef(null);
  const [contentWidth, setContentWidth] = useState(0);

  useLayoutEffect(() => {
    const element = measureRef.current;
    const measure = () => setContentWidth(Math.ceil(element.getBoundingClientRect().width));
    measure();
    if (typeof ResizeObserver === 'undefined') return;
    const observer = new ResizeObserver(measure);
    observer.observe(element);
    return () => observer.disconnect();
  }, [scores, label, scoreClassName, maxWidth, padding]);

  return {
    width: Math.max(minWidth, Math.min(maxWidth, contentWidth + padding)),
    measurement: (
      <div
        ref={measureRef}
        aria-hidden="true"
        className="pointer-events-none invisible fixed -left-[10000px] top-0 grid w-max overflow-hidden font-mono whitespace-nowrap"
        style={{ maxWidth: Math.max(0, maxWidth - padding) }}
      >
        <span className="col-start-1 row-start-1 min-w-0 truncate text-[10px] tracking-[0.18em] uppercase">
          {label}
        </span>
        {scores.map((score, index) => (
          <span key={index} className={`col-start-1 row-start-1 min-w-0 truncate tabular-nums ${scoreClassName}`}>
            {score}
          </span>
        ))}
      </div>
    ),
  };
}
