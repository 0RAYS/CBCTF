import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Chip } from '../../../common';
import { ellipsis, resolveVisibleInlineItems } from './trafficPresentation.js';

export default function TrafficHeader({ topology, container, demoMode }) {
  const { t } = useTranslation();
  const accessIPRowRef = useRef(null);
  const [accessIPRowWidth, setAccessIPRowWidth] = useState(0);
  const accessIPs = topology?.ips || [];
  const { visibleItems: visibleAccessIPs, hiddenItems: hiddenAccessIPs } = resolveVisibleInlineItems(
    accessIPs,
    accessIPRowWidth
  );
  const visibleExposed = (topology?.center?.exposed || []).slice(0, 1);

  useEffect(() => {
    const element = accessIPRowRef.current;
    if (!element) return;
    const updateWidth = () => {
      const nextWidth = Math.round(element.getBoundingClientRect().width);
      setAccessIPRowWidth((current) => (current === nextWidth ? current : nextWidth));
    };
    updateWidth();
    if (typeof ResizeObserver === 'undefined') return;
    const observer = new ResizeObserver(updateWidth);
    observer.observe(element);
    return () => observer.disconnect();
  }, []);

  return (
    <div className="rounded-2xl border border-neutral-600 bg-black/20 px-4 py-3">
      <div className="flex flex-col gap-3">
        <div className="min-w-0">
          <div className="text-[11px] uppercase tracking-[0.24em] text-geek-400/80">
            {t('admin.contests.trafficGraph.hero.kicker')}
          </div>
          <div className="mt-1 text-lg font-['Maple_UI'] text-white">
            {topology?.center?.label || `Victim #${container?.id || '-'}`}
          </div>
          <div className="mt-1 max-w-[64ch] text-xs leading-5 text-neutral-400">
            {t('admin.contests.trafficGraph.hero.subtitle', {
              challenge: container?.challenge || container?.contest_challenge_name || `#${container?.id}`,
            })}
          </div>
        </div>
        <div ref={accessIPRowRef} className="min-w-0">
          {accessIPs.length > 0 ? (
            <div className="flex flex-nowrap items-center gap-2 overflow-hidden">
              {visibleAccessIPs.map((ip) => (
                <Chip
                  key={ip}
                  label={ip}
                  title={ip}
                  variant="tag"
                  size="sm"
                  className="shrink-0"
                  colorClass="border-neutral-400/30 bg-neutral-400/10 text-neutral-300"
                />
              ))}
              {hiddenAccessIPs.length > 0 ? (
                <Chip
                  label={`+${hiddenAccessIPs.length}`}
                  title={hiddenAccessIPs.join('\n')}
                  variant="tag"
                  size="sm"
                  className="shrink-0 cursor-help"
                  colorClass="border-neutral-500/30 bg-neutral-500/10 text-neutral-400"
                />
              ) : null}
            </div>
          ) : null}
        </div>
        <div className="flex flex-wrap gap-2">
          {visibleExposed.map((item) => (
            <Chip
              key={item}
              label={ellipsis(item, 26)}
              variant="tag"
              size="sm"
              colorClass="border-neutral-500/30 bg-black/20 text-neutral-300"
            />
          ))}
          {demoMode ? (
            <Chip
              label={t('admin.contests.trafficGraph.hero.demoMode')}
              variant="tag"
              size="sm"
              colorClass="border-neutral-400/30 bg-neutral-400/10 text-neutral-300"
            />
          ) : null}
        </div>
      </div>
    </div>
  );
}
