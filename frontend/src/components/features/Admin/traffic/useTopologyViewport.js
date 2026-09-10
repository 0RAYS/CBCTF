import { useEffect, useRef, useState } from 'react';
import { clampPanToViewport, getViewportMetrics, MAX_ZOOM, MIN_ZOOM, ZOOM_STEP } from './trafficLayout.js';
import { clamp } from './trafficPresentation.js';

export default function useTopologyViewport(scopeKey) {
  const canvasRef = useRef(null);
  const dragRef = useRef({ active: false, moved: false, x: 0, y: 0, panX: 0, panY: 0 });
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const viewport = getViewportMetrics(zoom);
  const viewBox = `${pan.x} ${pan.y} ${viewport.width} ${viewport.height}`;

  const resetView = () => {
    setZoom(1);
    setPan({ x: 0, y: 0 });
  };

  useEffect(() => {
    resetView();
    dragRef.current.active = false;
    dragRef.current.moved = false;
  }, [scopeKey]);

  const applyZoom = (targetZoom, anchor) => {
    const rect = canvasRef.current?.getBoundingClientRect();
    const nextZoom = clamp(targetZoom, MIN_ZOOM, MAX_ZOOM);
    if (!rect?.width || !rect?.height) {
      setZoom(nextZoom);
      setPan((current) => clampPanToViewport(current, nextZoom));
      return;
    }
    const anchorX = anchor?.x ?? rect.width / 2;
    const anchorY = anchor?.y ?? rect.height / 2;
    const currentViewport = getViewportMetrics(zoom);
    const nextViewport = getViewportMetrics(nextZoom);
    const worldX = pan.x + (anchorX / rect.width) * currentViewport.width;
    const worldY = pan.y + (anchorY / rect.height) * currentViewport.height;
    const nextPan = clampPanToViewport(
      {
        x: worldX - (anchorX / rect.width) * nextViewport.width,
        y: worldY - (anchorY / rect.height) * nextViewport.height,
      },
      nextZoom
    );
    setZoom(nextZoom);
    setPan(nextPan);
  };

  useEffect(() => {
    const element = canvasRef.current;
    if (!element) return;
    const handleNativeWheel = (event) => {
      event.preventDefault();
      event.stopPropagation();
      const rect = element.getBoundingClientRect();
      const nextZoom = clamp(zoom + (event.deltaY < 0 ? ZOOM_STEP : -ZOOM_STEP), MIN_ZOOM, MAX_ZOOM);
      applyZoom(nextZoom, { x: event.clientX - rect.left, y: event.clientY - rect.top });
    };
    element.addEventListener('wheel', handleNativeWheel, { passive: false });
    return () => element.removeEventListener('wheel', handleNativeWheel);
  }, [zoom, pan]);

  const startDrag = (event) => {
    if (event.button !== 0) return;
    dragRef.current = { active: true, moved: false, x: event.clientX, y: event.clientY, panX: pan.x, panY: pan.y };
  };
  const onDrag = (event) => {
    if (!dragRef.current.active) return;
    const rect = canvasRef.current?.getBoundingClientRect();
    if (!rect?.width || !rect?.height) return;
    const dx = event.clientX - dragRef.current.x;
    const dy = event.clientY - dragRef.current.y;
    if (Math.abs(dx) > 3 || Math.abs(dy) > 3) dragRef.current.moved = true;
    const currentViewport = getViewportMetrics(zoom);
    setPan(
      clampPanToViewport(
        {
          x: dragRef.current.panX - dx * (currentViewport.width / rect.width),
          y: dragRef.current.panY - dy * (currentViewport.height / rect.height),
        },
        zoom
      )
    );
  };
  const stopDrag = () => {
    dragRef.current.active = false;
  };
  const canSelect = () => !dragRef.current.moved;

  return { canvasRef, zoom, viewBox, applyZoom, resetView, startDrag, onDrag, stopDrag, canSelect };
}
