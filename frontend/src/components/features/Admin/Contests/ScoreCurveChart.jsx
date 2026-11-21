import { lazy, Suspense, useCallback, useEffect, useMemo, useRef } from 'react';
import { useTranslation } from 'react-i18next';

const ReactECharts = lazy(() => import('@/components/common/EChart'));

// ECharts 需要字面量色值，取自 tailwind.config.js 的设计 token
const AXIS_LINE = '#2a2a2a';
const AXIS_LABEL = '#8a8a8a';
const SPLIT_LINE = '#1a1a1a';
const TOOLTIP_TEXT = '#d9d9d9';
const ACCENT = '#597ef7';
const HANDLE = 'rgba(251, 146, 60, 0.9)';

/**
 * 分数曲线预览图表（支持拖拽控制点调整 decay）
 * Static 无控制点；Linear 沿曲线拖动改变下降速率；Exponential 自由拖动并由 (x, y) 反算 decay；score 与 minScore 仅通过输入框调整
 */
function ScoreCurveChart({ scoreType = 0, score = 1000, decay = 50, minScore = 100, onChange }) {
  const { t } = useTranslation();
  const chartRef = useRef(null);
  const propsRef = useRef({ scoreType, score, decay, minScore });
  useEffect(() => {
    propsRef.current = { scoreType, score, decay, minScore };
  }, [scoreType, score, decay, minScore]);

  // 与后端 CalcScore 一致
  const calcScore = useCallback((solvers, s, d, ms, st) => {
    let calc;
    switch (st) {
      case 0:
        calc = s;
        break;
      case 1:
        calc = s - solvers * d;
        break;
      case 2:
        if (d > 0) {
          const k = 5.0 / d;
          calc = (s - ms) * Math.exp(-k * solvers) + ms;
        } else {
          calc = ms;
        }
        break;
      default:
        calc = s;
    }
    if (calc < ms) calc = ms;
    return Math.trunc(calc * 100) / 100;
  }, []);

  const maxSolvers = useMemo(() => {
    let m;
    if (scoreType === 0) {
      m = 20;
    } else if (scoreType === 1) {
      m = decay > 0 ? Math.ceil((score - minScore) / decay) + 10 : 20;
    } else {
      m = decay > 0 ? Math.ceil(decay * 1.4) : 20;
    }
    m = Math.max(m, 10);
    m = Math.min(m, 500);
    return m;
  }, [scoreType, score, decay, minScore]);

  const lineData = useMemo(() => {
    const step = maxSolvers <= 100 ? 1 : Math.ceil(maxSolvers / 100);
    const data = [];
    for (let i = 0; i <= maxSolvers; i += step) {
      data.push([i, calcScore(i, score, decay, minScore, scoreType)]);
    }
    if (data[data.length - 1][0] !== maxSolvers) {
      data.push([maxSolvers, calcScore(maxSolvers, score, decay, minScore, scoreType)]);
    }
    return data;
  }, [scoreType, score, decay, minScore, maxSolvers, calcScore]);

  const pixelToData = useCallback((px, py) => {
    const inst = chartRef.current?.getEchartsInstance();
    if (!inst) return null;
    return inst.convertFromPixel({ gridIndex: 0 }, [px, py]);
  }, []);

  const dataToPixel = useCallback((dx, dy) => {
    const inst = chartRef.current?.getEchartsInstance();
    if (!inst) return null;
    return inst.convertToPixel({ gridIndex: 0 }, [dx, dy]);
  }, []);

  const onHandleDrag = useCallback(
    (e) => {
      const pos = pixelToData(e.offsetX, e.offsetY);
      if (!pos) return;
      const { scoreType: st, score: s, minScore: ms } = propsRef.current;
      const x = Math.max(pos[0], 0.5);
      const y = pos[1];

      let newDecay;
      if (st === 1) {
        // Linear: y = s - x * decay  →  decay = (s - y) / x
        const clampedY = Math.max(Math.min(y, s - 0.1), ms);
        newDecay = (s - clampedY) / x;
      } else {
        // Exponential: y = (s - ms) * e^(-5x/decay) + ms
        // → (y - ms) / (s - ms) = e^(-5x/decay)
        // → decay = -5x / ln((y - ms) / (s - ms))
        const range = s - ms;
        if (range <= 0) return;
        const ratio = Math.max(Math.min((y - ms) / range, 0.99), 0.01);
        newDecay = (-5 * x) / Math.log(ratio);
      }

      newDecay = Math.round(Math.max(newDecay, 0.1) * 10) / 10;
      newDecay = Math.min(newDecay, 9999);
      if (onChange) onChange({ decay: newDecay });
    },
    [pixelToData, onChange]
  );

  // 控制点落在 score 与 minScore 中点对应的 solvers 处
  const handleDataPos = useMemo(() => {
    if (scoreType === 0) return null;
    const midY = (score + minScore) / 2;

    if (scoreType === 1) {
      // Linear: midY = score - x * decay → x = (score - midY) / decay
      if (decay <= 0) return null;
      const x = (score - midY) / decay;
      return [x, midY];
    } else {
      // Exponential: midY = (score - minScore) * e^(-5x/decay) + minScore
      // → x = -decay/5 * ln((midY - minScore) / (score - minScore))
      const range = score - minScore;
      if (range <= 0 || decay <= 0) return null;
      const ratio = (midY - minScore) / range;
      if (ratio <= 0 || ratio >= 1) return null;
      const x = (-decay / 5) * Math.log(ratio);
      const y = calcScore(x, score, decay, minScore, 2);
      return [x, y];
    }
  }, [scoreType, score, decay, minScore, calcScore]);

  const buildGraphicElements = useCallback(() => {
    if (!onChange || !handleDataPos) return [];

    const px = dataToPixel(handleDataPos[0], handleDataPos[1]);
    if (!px) return [];

    return [
      {
        type: 'circle',
        shape: { cx: 0, cy: 0, r: 4 },
        position: [px[0], px[1]],
        style: {
          fill: HANDLE,
          stroke: '#f9f9f9',
          lineWidth: 1.5,
          shadowBlur: 4,
          shadowColor: 'rgba(251, 146, 60, 0.5)',
        },
        cursor: 'grab',
        draggable: true,
        z: 100,
        ondrag: (e) => onHandleDrag(e),
      },
    ];
  }, [handleDataPos, dataToPixel, onHandleDrag, onChange]);

  const curveLabels = useMemo(
    () => ({
      0: t('admin.challengeModal.scoreCurve.static'),
      1: t('admin.challengeModal.scoreCurve.linear'),
      2: t('admin.challengeModal.scoreCurve.log'),
    }),
    [t]
  );

  const option = useMemo(() => {
    return {
      grid: {
        left: 8,
        right: 16,
        top: 24,
        bottom: 8,
        containLabel: true,
      },
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(0, 0, 0, 0.9)',
        borderColor: 'rgba(255, 255, 255, 0.2)',
        borderWidth: 1,
        borderRadius: 6,
        textStyle: {
          color: TOOLTIP_TEXT,
          fontSize: 11,
          fontFamily: '"Maple Mono", "Source Han Sans SC", ui-monospace, monospace',
        },
        formatter: (params) => {
          const p = params[0];
          return (
            `<div style="font-family:'Maple Mono','Source Han Sans SC',ui-monospace,monospace;font-size:11px">` +
            `<span style="color:${AXIS_LABEL}">Solvers:</span> ${p.value[0]}<br/>` +
            `<span style="color:${AXIS_LABEL}">Score:</span> <span style="color:${ACCENT};font-weight:bold">${p.value[1]}</span>` +
            `</div>`
          );
        },
      },
      xAxis: {
        type: 'value',
        name: t('admin.challengeModal.labels.solvers'),
        nameTextStyle: {
          color: AXIS_LABEL,
          fontSize: 10,
          fontFamily: '"Maple Mono", "Source Han Sans SC", ui-monospace, monospace',
        },
        min: 0,
        max: maxSolvers,
        axisLine: { lineStyle: { color: AXIS_LINE } },
        axisLabel: { color: AXIS_LABEL, fontSize: 10, fontFamily: 'monospace' },
        splitLine: { lineStyle: { color: SPLIT_LINE, type: 'dashed' } },
      },
      yAxis: {
        type: 'value',
        name: t('admin.challengeModal.labels.initialScore'),
        nameTextStyle: {
          color: AXIS_LABEL,
          fontSize: 10,
          fontFamily: '"Maple Mono", "Source Han Sans SC", ui-monospace, monospace',
        },
        min: (value) => Math.max(0, Math.floor(value.min * 0.9)),
        axisLine: { lineStyle: { color: AXIS_LINE } },
        axisLabel: { color: AXIS_LABEL, fontSize: 10, fontFamily: 'monospace' },
        splitLine: { lineStyle: { color: SPLIT_LINE, type: 'dashed' } },
      },
      series: [
        {
          name: curveLabels[scoreType] || '',
          type: 'line',
          data: lineData,
          smooth: scoreType === 2,
          showSymbol: false,
          lineStyle: {
            color: ACCENT,
            width: 2,
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(89, 126, 247, 0.25)' },
                { offset: 1, color: 'rgba(89, 126, 247, 0.02)' },
              ],
            },
          },
        },
      ],
    };
  }, [scoreType, lineData, maxSolvers, curveLabels, t]);

  const onChartReady = useCallback(
    (instance) => {
      setTimeout(() => {
        const elements = buildGraphicElements();
        if (elements.length > 0) {
          instance.setOption({ graphic: elements }, { replaceMerge: ['graphic'] });
        }
      }, 50);
    },
    [buildGraphicElements]
  );

  useEffect(() => {
    const inst = chartRef.current?.getEchartsInstance();
    if (!inst || !onChange) return;
    const timer = setTimeout(() => {
      const elements = buildGraphicElements();
      // replaceMerge 会清除旧控制点（scoreType 切到 static 时 elements 为空）
      inst.setOption({ graphic: elements }, { replaceMerge: ['graphic'] });
    }, 60);
    return () => clearTimeout(timer);
  }, [scoreType, score, decay, minScore, maxSolvers, buildGraphicElements, onChange]);

  return (
    <div className="mt-2 border border-neutral-700/50 rounded-md bg-black/30 p-2">
      <Suspense fallback={<div className="flex h-[180px] items-center justify-center text-xs text-neutral-500" />}>
        <ReactECharts
          ref={chartRef}
          option={option}
          style={{ height: '180px', width: '100%' }}
          opts={{ renderer: 'canvas' }}
          onChartReady={onChartReady}
        />
      </Suspense>
    </div>
  );
}

export default ScoreCurveChart;
